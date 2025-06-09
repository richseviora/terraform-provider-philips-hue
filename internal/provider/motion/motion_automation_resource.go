package motion

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/richseviora/huego/pkg/resources/behavior_instance"
	"github.com/richseviora/huego/pkg/resources/client"
	"github.com/richseviora/huego/pkg/resources/common"
	"terraform-provider-philips-hue/internal/provider/device"
)

var (
	_ resource.ResourceWithImportState = &MotionAutomationResource{}
	_ resource.ResourceWithConfigure   = &MotionAutomationResource{}
	_ resource.Resource                = &MotionAutomationResource{}
)

type MotionAutomationResource struct {
	client device.ClientWithLightIDCache
}

type Reference struct {
	Id   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

type TimeSlot struct {
	Hour   types.Int32 `tfsdk:"hour"`
	Minute types.Int32 `tfsdk:"minute"`
	Scenes []Reference `tfsdk:"scenes"`
	// The delay period in minutes. 0 to 60.
	AfterDelay types.Int32  `tfsdk:"after_delay"`
	AfterState types.String `tfsdk:"after_state"`
}

type MotionAutomationResourceModel struct {
	ID            types.String `tfsdk:"id"`
	SensorID      types.String `tfsdk:"sensor_id"`
	Targets       []Reference  `tfsdk:"targets"`
	DarkThreshold types.Int32  `tfsdk:"dark_threshold"`
	TimeSlots     []TimeSlot   `tfsdk:"time_slots"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Name          types.String `tfsdk:"name"`
}

func NewMotionAutomationResource() resource.Resource {
	return &MotionAutomationResource{}
}

func (m MotionAutomationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	m.client = client
}

func (m MotionAutomationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_motion_automation"
}

func (m MotionAutomationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "",
				MarkdownDescription: "",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:          []validator.String{stringvalidator.OneOf("previous_state", "all_of")},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the motion automation.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
				},
			},
			"sensor_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the motion sensor.",
			},
			"targets": schema.ListNestedAttribute{
				Required:    true,
				Description: "The IDs of the targets (rooms or zones) to target when the sensor is triggered.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:    true,
							Description: "The ID of the target.",
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "The type of the target.",
							Validators:  []validator.String{stringvalidator.OneOf("room", "zone")},
						},
					},
				},
			},
			"dark_threshold": schema.Int32Attribute{
				Required:   true,
				Validators: []validator.Int32{int32validator.Between(0, 65535)},
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Whether the automation is enabled.",
			},
			"time_slots": schema.ListNestedAttribute{
				Required:    true,
				Description: "The time slots to trigger the automation.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"hour": schema.Int32Attribute{
							Required:    true,
							Description: "The hour of the time slot.",
							Validators:  []validator.Int32{int32validator.Between(0, 23)},
						},
						"minute": schema.Int32Attribute{
							Required:    true,
							Description: "The minute of the time slot.",
							Validators:  []validator.Int32{int32validator.Between(0, 59)},
						},
						"after_delay": schema.Int32Attribute{
							Required:    true,
							Description: "The delay in minutes after the time slot has been triggered.",
							Validators: []validator.Int32{
								int32validator.Between(0, 60),
							},
						},
						"after_state": schema.StringAttribute{
							Required:    true,
							Description: "The state to return the lights to after the delay period.",
							Validators:  []validator.String{stringvalidator.OneOf("previous_state", "all_of")},
						},
						"scenes": schema.ListNestedAttribute{
							Required:    true,
							Description: "The IDs of the target scenes to activate when the sensor is triggered.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Required:    true,
										Description: "The ID of the target.",
									},
									"type": schema.StringAttribute{
										Required:    true,
										Description: "The type of the target.",
										Validators:  []validator.String{stringvalidator.OneOf("scene")},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (m MotionAutomationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MotionAutomationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	scriptId, err := m.client.GetBehaviorScriptIDForMetadataName("Motion Sensor")
	create := SetCreateFromBody(data, scriptId)

	response, err := m.client.BehaviorInstanceService().CreateBehaviorInstance(ctx, create)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating motion automation",
			"Could not create motion automation ID: "+err.Error(),
		)
		return
	}
	data.ID = types.StringValue(response.RID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m MotionAutomationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MotionAutomationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resource, err := m.client.BehaviorInstanceService().GetBehaviorInstance(ctx, data.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading motion automation",
			"Could not read motion automation ID "+data.ID.ValueString()+": "+err.Error(),
		)
	}
	data = *SetModelFromBody(*resource)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m MotionAutomationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MotionAutomationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	update := SetUpdateFromBody(data)
	_, err := m.client.BehaviorInstanceService().UpdateBehaviorInstance(ctx, data.ID.ValueString(), update)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating motion automation",
			"Could not update motion automation ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m MotionAutomationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MotionAutomationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := m.client.BehaviorInstanceService().DeleteBehaviorInstance(ctx, data.ID.ValueString())
	if errors.Is(err, client.ErrNotFound) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting motion automation",
			"Could not delete motion automation ID "+data.ID.ValueString()+": "+err.Error(),
		)
	}
}

func (m MotionAutomationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func Map[T any, R any](input []T, mapper func(T) R) []R {
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = mapper(v)
	}
	return result
}

func SetModelFromBody(bi behavior_instance.Data) *MotionAutomationResourceModel {
	m := &MotionAutomationResourceModel{
		ID:       types.StringValue(bi.ID),
		SensorID: types.StringValue(bi.Configuration.Source.RID),
		Targets: Map(bi.Configuration.Where, func(t behavior_instance.Where) Reference {
			return Reference{
				Id:   types.StringValue(t.Group.RID),
				Type: types.StringValue(t.Group.RType),
			}
		}),
		DarkThreshold: types.Int32Value(int32(bi.Configuration.Settings.DaylightSensitivity.DarkThreshold)),
		TimeSlots: Map(bi.Configuration.When.Timeslots, func(t behavior_instance.TimeSlots) TimeSlot {
			return TimeSlot{
				Hour:   types.Int32Value(int32(t.StartTime.Time.Hour)),
				Minute: types.Int32Value(int32(t.StartTime.Time.Minute)),
				Scenes: Map(t.OnMotion.RecallSingle, func(t behavior_instance.RecallSingle) Reference {
					return Reference{
						Id:   types.StringValue(t.Action.Recall.RID),
						Type: types.StringValue(t.Action.Recall.RType),
					}
				}),
				AfterDelay: types.Int32Value(int32(t.OnNoMotion.After.Minutes)),
				AfterState: types.StringValue(t.OnNoMotion.RecallSingle[0].Action),
			}
		}),
	}

	return m
}

func SetCreateFromBody(model MotionAutomationResourceModel, scriptId string) behavior_instance.CreateRequest {
	return behavior_instance.CreateRequest{
		ScriptID: scriptId,
		Configuration: behavior_instance.Configuration{
			Settings: behavior_instance.Settings{
				DaylightSensitivity: behavior_instance.DaylightSensitivity{
					DarkThreshold: int(model.DarkThreshold.ValueInt32()),
					Offset:        7000,
				},
			},
			Source: common.Reference{
				RID:   model.SensorID.ValueString(),
				RType: "sensor",
			},
			When: behavior_instance.When{
				Timeslots: Map[TimeSlot](model.TimeSlots, func(t TimeSlot) behavior_instance.TimeSlots {
					return behavior_instance.TimeSlots{
						StartTime: behavior_instance.StartTime{
							Time: behavior_instance.Time{
								Hour:   int(t.Hour.ValueInt32()),
								Minute: int(t.Minute.ValueInt32()),
							},
							Type: "time",
						},
						OnMotion: behavior_instance.OnMotion{
							RecallSingle: Map[Reference, behavior_instance.RecallSingle](t.Scenes, func(s Reference) behavior_instance.RecallSingle {
								return behavior_instance.RecallSingle{
									Action: behavior_instance.Action{Recall: common.Reference{
										RID:   s.Id.ValueString(),
										RType: s.Type.ValueString(),
									}},
								}
							}),
						},
						OnNoMotion: behavior_instance.OnNoMotion{
							After: behavior_instance.After{
								Minutes: int(t.AfterDelay.ValueInt32()),
							},
							RecallSingle: []behavior_instance.RecallSingleNoMotion{
								{
									Action: t.AfterState.ValueString(),
								},
							},
						},
					}
				}),
			},
			Where: Map[Reference](model.Targets, func(t Reference) behavior_instance.Where {
				return behavior_instance.Where{
					Group: common.Reference{
						RID:   t.Id.ValueString(),
						RType: t.Type.ValueString(),
					},
				}
			}),
		},
		Enabled:  model.Enabled.ValueBool(),
		Metadata: &behavior_instance.Metadata{Name: model.Name.ValueString()},
	}
}

func SetUpdateFromBody(model MotionAutomationResourceModel) behavior_instance.UpdateRequest {
	return behavior_instance.UpdateRequest{
		Configuration: &behavior_instance.Configuration{
			Settings: behavior_instance.Settings{
				DaylightSensitivity: behavior_instance.DaylightSensitivity{
					DarkThreshold: int(model.DarkThreshold.ValueInt32()),
					Offset:        7000,
				},
			},
			Source: common.Reference{
				RID:   model.SensorID.ValueString(),
				RType: "sensor",
			},
			When: behavior_instance.When{
				Timeslots: Map[TimeSlot](model.TimeSlots, func(t TimeSlot) behavior_instance.TimeSlots {
					return behavior_instance.TimeSlots{
						StartTime: behavior_instance.StartTime{
							Time: behavior_instance.Time{
								Hour:   int(t.Hour.ValueInt32()),
								Minute: int(t.Minute.ValueInt32()),
							},
							Type: "time",
						},
						OnMotion: behavior_instance.OnMotion{
							RecallSingle: Map[Reference, behavior_instance.RecallSingle](t.Scenes, func(s Reference) behavior_instance.RecallSingle {
								return behavior_instance.RecallSingle{
									Action: behavior_instance.Action{Recall: common.Reference{
										RID:   s.Id.ValueString(),
										RType: s.Type.ValueString(),
									}},
								}
							}),
						},
						OnNoMotion: behavior_instance.OnNoMotion{
							After: behavior_instance.After{
								Minutes: int(t.AfterDelay.ValueInt32()),
							},
							RecallSingle: []behavior_instance.RecallSingleNoMotion{
								{
									Action: t.AfterState.ValueString(),
								},
							},
						},
					}
				}),
			},
			Where: Map[Reference](model.Targets, func(t Reference) behavior_instance.Where {
				return behavior_instance.Where{
					Group: common.Reference{
						RID:   t.Id.ValueString(),
						RType: t.Type.ValueString(),
					},
				}
			}),
		},
		Enabled:  model.Enabled.ValueBoolPointer(),
		Metadata: &behavior_instance.Metadata{Name: model.Name.ValueString()},
	}
}
