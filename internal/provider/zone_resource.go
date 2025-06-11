package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/richseviora/huego/pkg/resources/client"
	"terraform-provider-philips/internal/provider/device"

	"github.com/richseviora/huego/pkg/resources/common"
	"github.com/richseviora/huego/pkg/resources/zone"

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
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &ZoneResource{}
var _ resource.ResourceWithConfigure = &ZoneResource{}
var _ resource.ResourceWithImportState = &ZoneResource{}

type ZoneResource struct {
	client device.ClientWithLightIDCache
}

type Reference struct {
	RID   types.String `tfsdk:"rid"`
	RType types.String `tfsdk:"rtype"`
}

type ZoneResourceModel struct {
	ID        types.String   `tfsdk:"id"`
	Name      types.String   `tfsdk:"name"`
	LightIDs  []types.String `tfsdk:"light_ids"`
	Type      types.String   `tfsdk:"type"`
	Reference types.Object   `tfsdk:"reference"`
}

func NewZoneResource() resource.Resource {
	return &ZoneResource{}
}

func (z *ZoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(device.ClientWithLightIDCache)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected device.ClientWithLightIDCache, got: %T.", req.ProviderData),
		)
		return
	}
	z.client = client
}

func (z *ZoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone"
}

func (z *ZoneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}},
			"name": schema.StringAttribute{Required: true, Description: "The name of the zone. 1-32 characters long.", Validators: []validator.String{
				stringvalidator.LengthBetween(1, 32),
			}},
			"light_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "The IDs of the lights to be assigned to the zone.", PlanModifiers: []planmodifier.Set{},
			},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(common.AreaNames[:]...),
				},
			},
			"reference": schema.ObjectAttribute{
				Computed:            true,
				Description:         "The reference of the Zone in the Hue Bridge.",
				MarkdownDescription: "",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				AttributeTypes: map[string]attr.Type{
					"rid":   types.StringType,
					"rtype": types.StringType,
				},
			},
		},
		Description: "Represents a Philips Hue zone. Lights can belong to multiple zones at once.",
		Version:     0,
	}
}

func (z *ZoneResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (z *ZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ZoneResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "failed to populate record")
		return
	}

	body := createZoneBodyFromModel(data)

	createdBody, err := z.client.ZoneService().CreateZone(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating zone", err.Error())
		return
	}
	tflog.Info(ctx, "Creation Update", map[string]interface{}{
		"createdBody": createdBody,
	})
	data.ID = types.StringValue(createdBody.Data[0].RID)
	data.Reference, _ = types.ObjectValue(map[string]attr.Type{
		"rid":   types.StringType,
		"rtype": types.StringType,
	}, map[string]attr.Value{
		"rid":   types.StringValue(createdBody.Data[0].RID),
		"rtype": types.StringValue("room"),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (z *ZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ZoneResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	zone, err := z.client.ZoneService().GetZone(ctx, data.ID.ValueString())

	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading zone", err.Error())
		return
	}

	data = createZoneModelFromData(zone)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (z *ZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ZoneResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := createZoneBodyFromModel(data)
	tflog.Info(ctx, "Update", map[string]interface{}{"data": data, "number_of_children(body)": len(body.Children), "number_children(data)": len(data.LightIDs)})
	_, err := z.client.ZoneService().UpdateZone(ctx, data.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Zone", err.Error())
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (z *ZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ZoneResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := data.ID.ValueString()
	err := z.client.ZoneService().DeleteZone(ctx, id)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return
		}
		resp.Diagnostics.AddError("Error deleting zone", err.Error())
	}
}

func createZoneModelFromData(data *zone.ZoneData) ZoneResourceModel {
	lightIds := make([]types.String, len(data.Children))
	for i, child := range data.Children {
		lightIds[i] = types.StringValue(child.RID)
	}
	reference, _ := types.ObjectValue(map[string]attr.Type{
		"rid":   types.StringType,
		"rtype": types.StringType,
	}, map[string]attr.Value{
		"rid":   types.StringValue(data.ID),
		"rtype": types.StringValue("room"),
	})
	return ZoneResourceModel{
		ID:        types.StringValue(data.ID),
		Name:      types.StringValue(data.Metadata.Name),
		LightIDs:  lightIds,
		Type:      types.StringValue(data.Metadata.Archetype),
		Reference: reference,
	}
}

func createZoneBodyFromModel(model ZoneResourceModel) *zone.ZoneCreateOrUpdate {
	children := make([]common.Reference, len(model.LightIDs))
	for i, id := range model.LightIDs {
		children[i] = common.Reference{
			RID:   id.ValueString(),
			RType: "light",
		}
	}
	return &zone.ZoneCreateOrUpdate{
		Children: children,
		Metadata: zone.ZoneMetadata{
			Name:      model.Name.ValueString(),
			Archetype: model.Type.ValueString(),
		},
	}
}
