package motion

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/richseviora/huego/pkg/resources/motion"
	"regexp"
	"terraform-provider-philips-hue/internal/provider/device"
)

var (
	_ resource.Resource                = &MotionResource{}
	_ resource.ResourceWithImportState = &MotionResource{}
	_ resource.ResourceWithConfigure   = &MotionResource{}
)

type MotionResourceModel struct {
	Id        types.String `tfsdk:"id"`
	DeviceID  types.String `tfsdk:"device_id"`
	Reference types.Object `tfsdk:"reference"`
	Enabled   types.Bool   `tfsdk:"enabled"`
}

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
	var data MotionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resource, err := m.client.MotionService().GetMotion(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading light",
			"Could not read motion ID "+data.Id.ValueString()+": "+err.Error())
		return
	}
	data.Id = types.StringValue(resource.ID)
	data.Enabled = types.BoolValue(resource.Enabled)
	data.Reference, _ = types.ObjectValue(map[string]attr.Type{
		"rid":   types.StringType,
		"rtype": types.StringType,
	}, map[string]attr.Value{
		"rid":   types.StringValue(resource.ID),
		"rtype": types.StringValue("motion"),
	})
	data.DeviceID = types.StringValue(resource.Owner.RID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m MotionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MotionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	update := &motion.UpdateRequest{
		Enabled: data.Enabled.ValueBool(),
	}
	_, err := m.client.MotionService().UpdateMotion(ctx, data.Id.ValueString(), *update)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating motion",
			"Could not update motion ID "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m MotionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError("Not implemented", "Direct delete is not supported for this resource. Please remove the resource from the app instead.")
	return
}

func NewMotionResource() resource.Resource {
	return &MotionResource{}
}
