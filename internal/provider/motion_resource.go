package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"terraform-provider-philips-hue/internal/provider/device"
)

var (
	_ resource.Resource                = &MotionResource{}
	_ resource.ResourceWithImportState = &MotionResource{}
	_ resource.ResourceWithConfigure   = &MotionResource{}
)

type MotionResource struct {
	client device.ClientWithLightIDCache
}

func (m MotionResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	client, ok := request.ProviderData.(device.ClientWithLightIDCache)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected device.ClientWithLightIDCache, got: %T.", request.ProviderData),
		)
		return
	}
	m.client = client
}

func (m MotionResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	matched, err := regexp.MatchString(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`, request.ID)
	if err != nil {
		resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
	} else if !matched {
		resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
	}

	id, err := m.client.GetMotionIDForMacAddress(request.ID)
	if err != nil {
		response.Diagnostics.AddError("Error importing motion sensor", "Could not find motion with MAC address "+request.ID+": "+err.Error())
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), resource.ImportStateRequest{
		ID:                 id,
		ClientCapabilities: request.ClientCapabilities,
	}, response)
}

type MotionResourceModel struct {
	Id        string       `tfsdk:"id"`
	Name      string       `tfsdk:"name"`
	DeviceID  string       `tfsdk:"device_id"`
	Reference types.Object `tfsdk:"reference"`
	Enabled   bool         `tfsdk:"enabled"`
}

func (m MotionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_motion"
}

func (m MotionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A representation of a Philips Hue motion sensor.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the motion sensor in the Hue Bridge.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the motion sensor. 1-32 characters long.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
				},
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Whether the motion sensor is enabled.",
			},
			"reference": schema.ObjectAttribute{
				Computed:            true,
				Description:         "The reference of the Motion in the Hue Bridge.",
				MarkdownDescription: "",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				AttributeTypes: map[string]attr.Type{
					"rid":   types.StringType,
					"rtype": types.StringType,
				},
			},
			"device_id": schema.StringAttribute{
				Computed:    true,
				Description: "The Device ID of the motion sensor.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (m MotionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError("Not implemented", "Direct create is not supported for this resource. Please import the resource instead.")
	return
}

func (m MotionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (m MotionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (m MotionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError("Not implemented", "Direct delete is not supported for this resource. Please remove the resource from the app instead.")
	return
}

func NewMotionResource() resource.Resource {
	return &MotionResource{}
}
