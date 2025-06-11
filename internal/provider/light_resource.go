package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/richseviora/huego/pkg/resources/light"
	"regexp"
	"terraform-provider-philips/internal/provider/device"
)

var _ resource.Resource = &LightResource{}
var _ resource.ResourceWithImportState = &LightResource{}
var _ resource.ResourceWithConfigure = &LightResource{}

type LightResource struct {
	client device.ClientWithLightIDCache
}

type LightResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Function types.String `tfsdk:"function"`
	DeviceID types.String `tfsdk:"device_id"`
	// Archetype types.String `tfsdk:"archetype"`
	// TODO: Add power-on attributes

}

func NewLightResource() resource.Resource {
	return &LightResource{}
}

func (l *LightResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_light"
}

func (l *LightResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "A representation of a Philips Hue light.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The UUID of the Light device in the Hue Bridge. ",
				MarkdownDescription: "",
				DeprecationMessage:  "",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Default:   nil,
				WriteOnly: false,
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the Light device in the Hue Bridge. ",
				MarkdownDescription: "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"function": schema.StringAttribute{
				Optional:            true,
				Description:         "The function of the Light device in the Hue Bridge. ",
				MarkdownDescription: "",
				Validators: []validator.String{
					stringvalidator.OneOf("mixed", "decorative", "functional", "unknown"),
				},
			},
			"device_id": schema.StringAttribute{
				Computed:    true,
				Sensitive:   false,
				Description: "The device UUID of the light in the Hue Bridge. This ID is used to assign device membership in rooms and scenes.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (l *LightResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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
	l.client = client
}

func (l *LightResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	response.Diagnostics.AddError("Not implemented", "Direct create is not supported for this resource. Please import the resource instead.")
	return
}

func (l *LightResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data LightResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}
	light, err := l.client.LightService().GetLight(ctx, data.Id.ValueString())
	tflog.Info(ctx, "Returning Value", map[string]interface{}{"light": light, "err": err, "id": data.Id.ValueString()})
	if err != nil {
		response.Diagnostics.AddError(
			"Error reading light",
			"Could not read light ID "+data.Id.ValueString()+": "+err.Error())
		return
	}
	data.Name = types.StringValue(light.Metadata.Name)
	data.Function = types.StringValue(light.Metadata.Function)
	data.Id = types.StringValue(light.ID)
	data.DeviceID = types.StringValue(light.Owner.RID)
	if data.Id.ValueString() == "" {
		response.Diagnostics.AddError(
			"Error reading light",
			"Light ID not returned from object")
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (l *LightResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data LightResourceModel

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	update := struct {
		Name     *string `json:"name"`
		Function *string `json:"function"`
	}{Name: nil, Function: nil}
	if data.Name.ValueString() != "" {
		update.Name = data.Name.ValueStringPointer()
	}
	if data.Function.ValueString() != "" {
		update.Function = data.Function.ValueStringPointer()
	} else if data.Function.IsNull() || data.Function.IsUnknown() {
		s := "decorative"
		update.Function = &s
	}

	lightUpdate := light.LightUpdate{
		ID:       data.Id.ValueString(),
		Metadata: &update,
	}
	tflog.Info(ctx, "Updating Light", map[string]interface{}{"id": data.Id.ValueString(), "light": lightUpdate})
	err := l.client.LightService().UpdateLight(ctx, lightUpdate)

	tflog.Info(ctx, "Returning Updated Light", map[string]interface{}{"err": err, "id": data.Id.ValueString()})

	if err != nil {
		response.Diagnostics.AddError(
			"Error updating light",
			"Could not update light ID "+data.Id.ValueString()+": "+err.Error())
		return
	}
	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (l *LightResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	response.Diagnostics.AddError("Not implemented", "Delete is not supported for this resource. Please delete the light from the app instead.")
	return
}

func (l *LightResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	matched, err := regexp.MatchString(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`, request.ID)
	if err != nil {
		resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
	} else if !matched {
		resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
	}

	lightID, err := l.client.GetLightIDForMacAddress(request.ID)
	if err != nil {
		response.Diagnostics.AddError("Error importing light", "Could not find light with MAC address "+request.ID+": "+err.Error())
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), resource.ImportStateRequest{
		ID:                 lightID,
		ClientCapabilities: request.ClientCapabilities,
	}, response)

}
