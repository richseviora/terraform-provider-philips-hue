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
)
import "github.com/richseviora/huego/pkg/resources"

var _ resource.Resource = &ZoneResource{}
var _ resource.ResourceWithConfigure = &ZoneResource{}
var _ resource.ResourceWithImportState = &ZoneResource{}

type ZoneResource struct {
	client *resources.APIClient
}

type ZoneResourceModel struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	LightIDs []string `json:"light_ids"`
	Type     string   `json:"type"`
}

func NewZoneResource() resource.Resource {
	return &ZoneResource{}
}

func (z *ZoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	z.client = client
}

func (z *ZoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scene"
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
					stringvalidator.OneOf(resources.AreaNames[:]...),
				},
			},
		},
		Blocks:              nil,
		Description:         "Represents a Philips Hue zone. Lights can be long to multiple zones at once.",
		MarkdownDescription: "",
		DeprecationMessage:  "",
		Version:             0,
	}
}

func (z *ZoneResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (z *ZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//TODO implement me
	panic("implement me")
}

func (z *ZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (z *ZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (z *ZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}
