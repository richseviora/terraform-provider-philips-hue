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
	"github.com/richseviora/huego/pkg/resources"
)

var _ resource.Resource = &RoomResource{}
var _ resource.ResourceWithImportState = &RoomResource{}
var _ resource.ResourceWithConfigure = &RoomResource{}

func NewRoomResource() resource.Resource {
	return &RoomResource{}
}

type RoomResource struct {
	client *resources.APIClient
}

type RoomResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	DeviceIds types.Set    `tfsdk:"device_ids"`
	Archetype types.String `tfsdk:"archetype"`
	Reference types.Object `tfsdk:"reference"`
}

func (r *RoomResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_room"
}

func (r *RoomResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "A representation of a Philips Hue room.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The UUID of the Room in the Hue Bridge.",
				MarkdownDescription: "",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the Room in the Hue Bridge.",
				MarkdownDescription: "",
			},
			"device_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Required:            true,
				Description:         "The Device IDs to assign to the Room in the Hue Bridge.",
				MarkdownDescription: "",
				DeprecationMessage:  "",
				Validators:          []validator.Set{},
				PlanModifiers:       nil,
				Default:             nil,
			},
			"archetype": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(resources.AreaNames[:]...),
				},
			},
			"reference": schema.ObjectAttribute{
				Computed:            true,
				Description:         "The reference of the Room in the Hue Bridge.",
				MarkdownDescription: "",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				AttributeTypes: map[string]attr.Type{
					"rid":  types.StringType,
					"type": types.StringType,
				},
			},
		},
	}
}

func (r *RoomResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	client, ok := request.ProviderData.(*resources.APIClient)
	if !ok {
		response.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *resources.APIClient, got: %T. Please report this issue to the provider developers.", request.ProviderData))
	}
	r.client = client
}

func (r *RoomResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data RoomResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	area, err := resources.ParseArea(data.Archetype.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Error parsing room archetype",
			"Could not parse room archetype: "+err.Error())
		return
	}

	children := make([]resources.Child, len(data.DeviceIds.Elements()))
	for _, deviceId := range data.DeviceIds.Elements() {
		children = append(children, resources.Child{
			Rid:   deviceId.String(),
			Rtype: "light",
		})
	}

	room := resources.RoomCreate{
		Metadata: resources.RoomMetadata{
			Name:      data.Name.ValueString(),
			Archetype: area,
		},
		Children: children,
	}

	createdRoom, err := r.client.RoomService.CreateRoom(ctx, room)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating room",
			"Could not create room: "+err.Error())
		return
	}

	data.Reference, _ = types.ObjectValue(map[string]attr.Type{
		"rid":   types.StringType,
		"rtype": types.StringType,
	}, map[string]attr.Value{
		"rid":   types.StringValue(createdRoom.RID),
		"rtype": types.StringValue("room"),
	})
	data.Id = types.StringValue(createdRoom.RID)
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *RoomResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data RoomResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	room, err := r.client.RoomService.GetRoom(ctx, data.Id.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Error reading room",
			"Could not read room ID "+data.Id.ValueString()+": "+err.Error())
		return
	}

	data.Name = types.StringValue(room.Metadata.Name)
	data.Archetype = types.StringValue(room.Metadata.Archetype.String())
	data.Id = types.StringValue(room.ID)
	data.Reference, _ = types.ObjectValue(map[string]attr.Type{
		"rid":   types.StringType,
		"rtype": types.StringType,
	}, map[string]attr.Value{
		"rid":   types.StringValue(room.ID),
		"rtype": types.StringValue("room"),
	})
	deviceIds := make([]string, len(room.Children))
	for i, child := range room.Children {
		deviceIds[i] = child.Rid
	}
	data.DeviceIds, _ = types.SetValueFrom(ctx, types.StringType, deviceIds)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *RoomResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data RoomResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	area, err := resources.ParseArea(data.Archetype.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Error parsing room archetype",
			"Could not parse room archetype: "+err.Error())
		return
	}

	children := make([]resources.Child, len(data.DeviceIds.Elements()))
	for _, deviceId := range data.DeviceIds.Elements() {
		children = append(children, resources.Child{
			Rid:   deviceId.String(),
			Rtype: "light",
		})
	}

	update := resources.RoomUpdate{
		ID: data.Id.ValueString(),
		Metadata: &resources.RoomMetadata{
			Name:      data.Name.ValueString(),
			Archetype: area,
		},
		Children: &children,
	}

	err = r.client.RoomService.UpdateRoom(ctx, update)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating room",
			"Could not update room ID "+data.Id.ValueString()+": "+err.Error())
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *RoomResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data RoomResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := fmt.Errorf("Not implemented")
	//err := r.client.RoomService.DeleteRoom(ctx, data.Id.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting room",
			"Could not delete room ID "+data.Id.ValueString()+": "+err.Error())
		return
	}
}

func (r *RoomResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}
