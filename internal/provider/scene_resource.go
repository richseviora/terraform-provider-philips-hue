package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/richseviora/huego/pkg/resources"
)

var _ resource.Resource = &SceneResource{}
var _ resource.ResourceWithImportState = &SceneResource{}
var _ resource.ResourceWithConfigure = &SceneResource{}

type SceneResource struct {
	client *resources.APIClient
}

type SceneResourceModel struct {
	Id      types.String     `tfsdk:"id"`
	Name    types.String     `tfsdk:"name"`
	Actions []types.String   `tfsdk:"actions"`
	Group   types.ObjectType `tfsdk:"group"`
}

func (s *SceneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scene"
}

func (s *SceneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A representation of a Philips Hue scene.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Scene in the Hue Bridge.",
			},
			"group": schema.ObjectAttribute{
				Required:    true,
				Description: "The group this scene belongs to.",
				AttributeTypes: map[string]attr.Type{
					"rid":   types.StringType,
					"rtype": types.StringType,
				},
			},
		},
	}
}

func (s *SceneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*resources.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *resources.APIClient, got: %T.", req.ProviderData),
		)
		return
	}
	s.client = client
}

func (s *SceneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SceneResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implementation depends on the API client capabilities
	resp.Diagnostics.AddError("Not implemented", "Implementation pending")
}

func (s *SceneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SceneResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scene, err := s.client.SceneService.GetScene(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading scene",
			"Could not read scene ID "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	data.Name = types.StringValue(scene.Metadata.Name)
	data.Group = types.StringValue(scene.Group.Rid)

	deviceIDs := make([]types.String, len(scene.Actions))
	for i, action := range scene.Actions {
		deviceIDs[i] = types.StringValue(action.Target.RID)
	}
	data.DeviceIDs = deviceIDs

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (s *SceneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SceneResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	update := resources.SceneUpdate{
		ID: data.Id.ValueString(),
		Metadata: &struct {
			Name *string `json:"name"`
		}{
			Name: data.Name.ValueStringPointer(),
		},
	}

	err := s.client.SceneService.UpdateScene(ctx, update)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating scene",
			"Could not update scene ID "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (s *SceneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError("Not implemented", "Delete is not supported for this resource. Please delete the scene from the app instead.")
}

func (s *SceneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
