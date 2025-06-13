// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestSceneResource_Unit(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testUnitExampleResourceConfig(),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestSceneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleResourceConfig("one"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"philips_scene.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("example-id"),
					),
					statecheck.ExpectKnownValue(
						"philips_scene.test",
						tfjsonpath.New("defaulted"),
						knownvalue.StringExact("example value when not configured"),
					),
					statecheck.ExpectKnownValue(
						"philips_scene.test",
						tfjsonpath.New("configurable_attribute"),
						knownvalue.StringExact("one"),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "philips_scene.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"configurable_attribute", "defaulted"},
			},
			// Update and Read testing
			{
				Config: testAccExampleResourceConfig("two"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"philips_scene.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("example-id"),
					),
					statecheck.ExpectKnownValue(
						"philips_scene.test",
						tfjsonpath.New("defaulted"),
						knownvalue.StringExact("example value when not configured"),
					),
					statecheck.ExpectKnownValue(
						"philips_scene.test",
						tfjsonpath.New("configurable_attribute"),
						knownvalue.StringExact("two"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testUnitExampleResourceConfig() string {
	return `
resource "philips_scene" "test" {
  name = "test scene"
  group = { id: "abc", type: "room" }
  actions = [
  {
  target_id = "abc"
  target_type = "light"
  on = true
  brightness = 100
  color_temperature = 2500
  color = { x = 0.31271592, y = 0.3290015 }
  }
]
}
`
}

func testAccExampleResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "philips_light" "test" {
  name = %[1]q
  function = "functional"
}

resource "philips_room" "room" {
  name = "test room"
  archetype = "living_room"
device_ids = [philips_light.test.device_id]
}

resource "philips_scene" "test" {
  name = "test scene"
 actions = []
 group = philips_room.room.reference
}
`, configurableAttribute)
}
