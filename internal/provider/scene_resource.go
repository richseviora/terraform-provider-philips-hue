package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
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
	"github.com/hashicorp/terraform-plugin-log/tfsdklog"
	"github.com/richseviora/huego/pkg/resources"
)

var _ resource.Resource = &SceneResource{}
var _ resource.ResourceWithImportState = &SceneResource{}
var _ resource.ResourceWithConfigure = &SceneResource{}

func NewSceneResource() resource.Resource {
	return &SceneResource{}
}

type SceneResource struct {
	client *resources.APIClient
}

type SceneActionColorModel struct {
	X types.Float64 `tfsdk:"x"`
	Y types.Float64 `tfsdk:"y"`
}

type ResourceReference struct {
	Rid   types.String `tfsdk:"rid"`
	Rtype types.String `tfsdk:"rtype"`
}

type SceneActionModel struct {
	TargetId         types.String           `tfsdk:"target_id"`
	TargetType       types.String           `tfsdk:"target_type"`
	Brightness       types.Float64          `tfsdk:"brightness"`
	On               types.Bool             `tfsdk:"on"`
	Color            *SceneActionColorModel `tfsdk:"color"`
	ColorTemperature types.Int32            `tfsdk:"color_temperature"`
}

type SceneResourceModel struct {
	Id      types.String       `tfsdk:"id"`
	Name    types.String       `tfsdk:"name"`
	Actions []SceneActionModel `tfsdk:"actions"`
	Group   ResourceReference  `tfsdk:"group"`
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
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"actions": schema.ListNestedAttribute{
				Required:    true,
				Description: "The actions and targets to perform when the scene is triggered.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"target_id": schema.StringAttribute{
							Required:    true,
							Description: "The target ID to apply the action to.",
						},
						"target_type": schema.StringAttribute{
							Required:    true,
							Description: "The target type to apply the action to.",
						},
						"brightness": schema.Float64Attribute{
							Required:    true,
							Description: "The brightness to apply to the target from 0 to 100",
							Validators: []validator.Float64{
								float64validator.Between(0, 100),
							},
						},
						"on": schema.BoolAttribute{
							Required:    true,
							Description: "Whether the target should be turned on or off.",
						},
						"color": schema.SingleNestedAttribute{
							Optional: true,
							Validators: []validator.Object{
								objectvalidator.ExactlyOneOf(
									path.Expressions{
										path.MatchRelative().AtParent().AtName("color_temperature"),
									}...,
								),
							},
							Attributes: map[string]schema.Attribute{
								"x": schema.Float64Attribute{
									Required:    true,
									Description: "The x value of the color to apply to the target.",
									Validators: []validator.Float64{
										float64validator.Between(0, 1),
									},
								},
								"y": schema.Float64Attribute{
									Required:    true,
									Description: "The y value of the color to apply to the target.",
									Validators: []validator.Float64{
										float64validator.Between(0, 1),
									},
								},
							},
						},
						"color_temperature": schema.Int32Attribute{
							Optional:    true,
							Description: "The color temperature to apply to the target from 2000 K to 6500 K",
							Validators: []validator.Int32{
								int32validator.Between(2000, 6500),
							},
						},
					},
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

	createObj := s.createSceneCreateObj(data)
	newObj, err := s.client.SceneService.CreateScene(ctx, createObj)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating scene",
			"Could not create scene: "+err.Error(),
		)
		return
	}
	data.Id = types.StringValue(newObj.RID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (s *SceneResource) createSceneCreateObj(data SceneResourceModel) resources.SceneCreate {
	actionTargets := s.createSceneActionObj(data)

	createObj := resources.SceneCreate{
		Metadata: resources.SceneMetadata{
			Name: data.Name.ValueString(),
		},
		Actions: actionTargets,
		Group: resources.Group{
			Rid:   data.Group.Rid.ValueString(),
			Rtype: data.Group.Rtype.ValueString(),
		},
	}
	return createObj
}

func (s *SceneResource) createSceneUpdateObj(data SceneResourceModel) resources.SceneUpdate {
	actionTargets := s.createSceneActionObj(data)

	return resources.SceneUpdate{
		Metadata: resources.SceneMetadata{
			Name: data.Name.ValueString(),
		},
		Actions: actionTargets,
	}
}

func (s *SceneResource) createSceneActionObj(data SceneResourceModel) []resources.ActionTarget {
	actionTargets := make([]resources.ActionTarget, len(data.Actions))
	for i, action := range data.Actions {
		newAction := resources.Action{
			On: &resources.On{
				On: action.On.ValueBool(),
			},
			Dimming: &resources.Dimming{Brightness: action.Brightness.ValueFloat64()},
		}
		if action.Color != nil {
			newAction.Color = &resources.Color{
				XY: resources.XYCoord{
					X: action.Color.X.ValueFloat64(),
					Y: action.Color.Y.ValueFloat64(),
				},
			}
		}
		if !action.ColorTemperature.IsNull() && !action.ColorTemperature.IsUnknown() {
			newAction.ColorTemperature = &resources.ColorTemperature{
				Mirek: int(resources.KelvinToMirekRounded(int32(action.ColorTemperature.ValueInt32()))),
			}
		}
		actionTarget := resources.ActionTarget{
			Target: resources.Target{
				Rid:   action.TargetId.ValueString(),
				Rtype: action.TargetType.ValueString(),
			},
			Action: newAction,
		}
		actionTargets[i] = actionTarget
	}
	return actionTargets
}

func (s *SceneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SceneResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scene, err := s.client.SceneService.GetScene(ctx, data.Id.ValueString())
	tfsdklog.Info(ctx, "scene:", map[string]interface{}{"scene": scene})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading scene",
			"Could not read scene ID "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	data.Name = types.StringValue(scene.Metadata.Name)
	data.Group = ResourceReference{
		Rid:   types.StringValue(scene.Group.Rid),
		Rtype: types.StringValue(scene.Group.Rtype),
	}
	actions := make([]SceneActionModel, len(scene.Actions))
	for i, action := range scene.Actions {
		tflog.Info(ctx, "action:", map[string]interface{}{"action": action})
		var onValue types.Bool
		if action.Action.On != nil {
			onValue = types.BoolValue(action.Action.On.On)
			tflog.Info(ctx, "onValue:", map[string]interface{}{"onValue": onValue})
		}
		var colorTemp types.Int32
		if action.Action.ColorTemperature != nil {
			mirek := action.Action.ColorTemperature.Mirek
			kelvin := resources.MirekToKelvinRounded(int32(mirek))
			colorTemp = types.Int32Value(int32(kelvin))
			tflog.Info(ctx, "colorTemp:", map[string]interface{}{"colorTemp": colorTemp.String(), "originalColorTemp": kelvin})
		}
		var color *SceneActionColorModel
		if action.Action.Color != nil {
			color = &SceneActionColorModel{
				X: types.Float64Value(action.Action.Color.XY.X),
				Y: types.Float64Value(action.Action.Color.XY.Y),
			}
			tflog.Info(ctx, "color:", map[string]interface{}{"color": color})
			fmt.Println(color)
		}
		model := SceneActionModel{
			TargetId:         types.StringValue(action.Target.Rid),
			TargetType:       types.StringValue(action.Target.Rtype),
			On:               onValue,
			Brightness:       types.Float64Value(action.Action.Dimming.Brightness),
			Color:            color,
			ColorTemperature: colorTemp,
		}
		tflog.Info(ctx, "writing action:", map[string]interface{}{"action": model})
		actions[i] = model
	}

	tflog.Info(ctx, "actions:", map[string]interface{}{"actions": actions})

	data.Actions = actions
	for _, action := range data.Actions {
		tflog.Info(ctx, "RETURN action:", map[string]interface{}{"target_id": action.TargetId.ValueString()})
	}
	// TODO: Populate actions

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (s *SceneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SceneResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	update := s.createSceneUpdateObj(data)

	_, err := s.client.SceneService.UpdateScene(ctx, data.Id.ValueString(), update)
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
