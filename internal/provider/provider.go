// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/richseviora/huego/pkg"
	"github.com/richseviora/huego/pkg/resources/zigbee_connectivity"
	"terraform-provider-philips-hue/internal/provider/device"
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

// PhilipsHueProviderModel describes the provider data model.
type PhilipsHueProviderModel struct {
	Endpoint      types.String `tfsdk:"endpoint"`
	OutputImports types.Bool   `tfsdk:"output_imports"`
}

func (p *PhilipsHueProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "philips"
	resp.Version = p.version
}

func (p *PhilipsHueProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "NOT IMPLEMENTED YET - The Philips Hue bridge IP address. Example: `192.168.1.100` or `philips-hue.local`.",
				Optional:            true,
			},
			"output_imports": schema.BoolAttribute{
				MarkdownDescription: "Whether the provider will output the import state for resources. Defaults to `false`.",
				Optional:            true,
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

	client, err := pkg.NewClientFromMDNS()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create client, got error: %s", err))
		return
	}

	clientWithCache := device.NewClientWithCache(client)
	if data.OutputImports.ValueBool() {
		devices, zigbeeErrors, err := clientWithCache.GetAllDevices()
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get devices, got error: %s", err))
		} else {
			resp.Diagnostics.AddWarning("Imports", generateImportOutput(devices, zigbeeErrors))
		}
	}

	resp.DataSourceData = clientWithCache
	resp.ResourceData = clientWithCache
}

func (p *PhilipsHueProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewLightResource,
		NewRoomResource,
		NewSceneResource,
		NewZoneResource,
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

func generateImportOutput(entries []device.DeviceMappingEntry, missingEntries []zigbee_connectivity.Data) string {
	result := ""

	for _, entry := range entries {
		if !entry.IsLight() {
			continue
		}
		result += fmt.Sprintf(`
import {
  # Name = %s
  id = "%s"
  to = philips_light.REPLACE_ME
}
`, entry.Name, entry.MacAddress)
	}

	for _, entry := range missingEntries {
		result += fmt.Sprintf(`
/* 
Could not resolve MAC address:
%+v
*/
`, entry)
	}
	return result
}
