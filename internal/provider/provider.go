// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/richseviora/huego/pkg"
	"github.com/richseviora/huego/pkg/resources/client"
	"terraform-provider-philips/internal/provider/device"
	"terraform-provider-philips/internal/provider/logger"
	"terraform-provider-philips/internal/provider/motion"
)

// Ensure PhilipsHueProvider satisfies various provider interfaces.
var _ provider.Provider = &PhilipsHueProvider{}
var _ provider.ProviderWithFunctions = &PhilipsHueProvider{}
var _ provider.ProviderWithEphemeralResources = &PhilipsHueProvider{}

// PhilipsHueProvider defines the provider implementation.
type PhilipsHueProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type PhilipsHueBridge struct {
	IPAddress      types.String `tfsdk:"ip_address"`
	ApplicationKey types.String `tfsdk:"application_key"`
}

type PhilipsHueClient struct {
	FilePath types.String `tfsdk:"file_path"`
	ID       types.String `tfsdk:"id"`
}

// PhilipsHueProviderModel describes the provider data model.
type PhilipsHueProviderModel struct {
	Bridge PhilipsHueBridge `tfsdk:"bridge"`
	Output types.String     `tfsdk:"output"`
	Client PhilipsHueClient `tfsdk:"auto_connect"`
}

func (p *PhilipsHueProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "philips"
	resp.Version = p.version
}

func (p *PhilipsHueProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bridge": schema.SingleNestedAttribute{
				MarkdownDescription: "The Philips Hue Bridge to connect to.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"ip_address": schema.StringAttribute{
						Required:    true,
						Description: "The IP address of the Philips Hue Bridge.",
					},
					"application_key": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "The application key for the Philips Hue Bridge.",
					},
				},
			},
			"output": schema.StringAttribute{
				MarkdownDescription: "If set, the location of the output file to write the import data to. Example: `/tmp/import.tf`. If set to \"STDOUT\", the output will be written as a warning.",
				Optional:            true,
			},
			"client": schema.SingleNestedAttribute{
				MarkdownDescription: "The client configuration to use for connecting to the bridge.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"file_path": schema.StringAttribute{
						Required:    true,
						Description: "The path to the client configuration file.",
					},
					"id": schema.StringAttribute{
						Optional:    true,
						Description: "The ID of the bridge to connect to.",
					},
				},
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.Expressions{
							path.MatchRelative().AtParent().AtName("bridge"),
						}...,
					),
				},
			},
		},
	}
}

func (p *PhilipsHueProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PhilipsHueProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, err := p.generateClient(ctx, resp, &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create client, got error: %s", err))
		return
	}

	clientWithCache := device.NewClientWithCache(c)
	p.output(data, clientWithCache, resp)

	resp.DataSourceData = clientWithCache
	resp.ResourceData = clientWithCache
}

func (p *PhilipsHueProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewLightResource,
		NewRoomResource,
		NewSceneResource,
		NewZoneResource,
		motion.NewMotionResource,
		motion.NewMotionAutomationResource,
	}
}

func (p *PhilipsHueProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		//NewExampleEphemeralResource,
	}
}

func (p *PhilipsHueProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		//NewExampleDataSource,
	}
}

func (p *PhilipsHueProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewRGBToColorFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PhilipsHueProvider{
			version: version,
		}
	}
}

func (p *PhilipsHueProvider) generateClient(ctx context.Context, resp *provider.ConfigureResponse, data *PhilipsHueProviderModel) (client.HueServiceClient, error) {
	if data.Bridge.IPAddress.ValueString() != "" && data.Bridge.ApplicationKey.ValueString() != "" {
		return pkg.NewClientWithoutPath(data.Bridge.IPAddress.ValueString(), data.Bridge.ApplicationKey.ValueString(), logger.NewContextLogger(ctx))
	} else {
		pv, err := pkg.NewClientProviderWithPath(data.Client.FilePath.ValueString(), logger.NewContextLogger(ctx))
		if err != nil {
			return nil, err
		}
		if data.Client.ID.ValueString() != "" {
			// if not directly configuring, check if bridge ID is set. If it is set, attempt to connect to that bridge ID.
			return pv.NewClientWithExistingBridge(data.Client.ID.ValueString())
		} else {
			// If it is not set, attempt to acquire a bridge ID.
			bridgeId, client, err := pv.NewClientWithNewBridge()
			if err != nil {
				return nil, err
			}
			data.Client.ID = types.StringValue(bridgeId)
			return client, nil
		}
	}
}

func (p *PhilipsHueProvider) output(data PhilipsHueProviderModel, clientWithCache *device.ClientWithCache, resp *provider.ConfigureResponse) {
	if data.Output.IsNull() || data.Output.IsUnknown() {
		return
	}
	devices, zigbeeErrors, err := clientWithCache.GetAllDevices()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get devices, got error: %s", err))
		return
	}
	if data.Output.ValueString() == "STDOUT" {
		resp.Diagnostics.AddWarning("Imports", device.GenerateImportOutput(devices, zigbeeErrors))
		return
	}
}
