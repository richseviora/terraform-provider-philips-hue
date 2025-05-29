package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/richseviora/huego/pkg/resources"
)

var (
	_ function.Function = RGBToColorFunction{}
)

func NewRGBToColorFunction() function.Function { return RGBToColorFunction{} }

type RGBToColorFunction struct{}

func (R RGBToColorFunction) Metadata(_ context.Context, request function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "rgb_to_color"
}

func (R RGBToColorFunction) Definition(_ context.Context, request function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Summary: "Generates color value from RGB value",
		Parameters: []function.Parameter{
			function.Int32Parameter{
				Name:                "red",
				MarkdownDescription: "Red value from 0 to 255",
			},
			function.Int32Parameter{
				Name:                "green",
				MarkdownDescription: "Green value from 0 to 255",
			},
			function.Int32Parameter{
				Name:                "blue",
				MarkdownDescription: "Blue value from 0 to 255",
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: map[string]attr.Type{
				"x":          types.Float32Type,
				"y":          types.Float32Type,
				"brightness": types.Int32Type,
			},
			CustomType: nil,
		},
	}
}

func (R RGBToColorFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var (
		red, green, blue int32
	)
	response.Error = function.ConcatFuncErrors(request.Arguments.Get(ctx, &red, &green, &blue))
	if response.Error != nil {
		return
	}
	coord, lum := resources.RGBToXY(resources.RGBColor{
		R: int(red),
		G: int(green),
		B: int(blue),
	})
	output := struct {
		X          types.Float32 `tfsdk:"x"`
		Y          types.Float32 `tfsdk:"y"`
		Brightness types.Int32   `tfsdk:"brightness"`
	}{
		X:          types.Float32Value(float32(coord.X)),
		Y:          types.Float32Value(float32(coord.Y)),
		Brightness: types.Int32Value(int32(lum)),
	}
	response.Error = function.ConcatFuncErrors(response.Result.Set(ctx, output))
}
