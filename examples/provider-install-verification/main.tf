terraform {
  required_providers {
    philips-hue = {
      source = "hashicorp.com/edu/philips-hue"
    }
  }
}

provider "philips-hue" {}

locals {
  bathroom_lights = [
    "8c995774-d012-46cf-9e0b-4c769f4f951f",
    "8779601e-c192-4bfe-b28e-84eaec316705",
    "646a60ec-ead5-4306-988b-7732223e978f",
    "64ad8457-e15a-459e-8559-29a134b8fb3e",
    "86599ae8-2354-488a-80a2-45e737de2a55",
    "e4806fac-e561-4823-90fe-b8340f663840",
  ]
}

import {
  to = philips-hue_room.bathroom
  id = "0d960eab-68c6-4ed7-8c0d-a24ca756d58e"
}

import {
  to = philips-hue_scene.bathroom_bright
  id = "68b39f81-1c15-4c82-bd0b-ab28606f3d2e"
}

import {
  for_each = local.bathroom_lights
  to = philips-hue_light.bathroom[each.key]
  id = each.value
}

resource philips-hue_light "first_light" {
  name     = "Hallway Lamp"
  function = "decorative"
}

resource philips-hue_light "bathroom" {
  count    = 6
  name     = "Bathroom ${count.index + 1}"
  function = "functional"
}

resource philips-hue_room "bathroom" {
  name      = "Bathroom"
  archetype = "bathroom"
  device_ids = philips-hue_light.bathroom[*].device_id
}

resource philips-hue_scene "bathroom_bright" {
  group = philips-hue_room.bathroom.reference
  name = "Bright"
  actions = [for light in philips-hue_light.bathroom : {
    target_id = light.id
    target_type = "light"
    on = true
    color_temperature = 2702
    brightness = 100
  }]
}

resource philips-hue_scene "bathroom_cool" {
  group = philips-hue_room.bathroom.reference
  name = "Bathroom Cold"
  actions = [for light in philips-hue_light.bathroom : {
    target_id = light.id
    target_type = "light"
    on = true
    color_temperature = 6500
    brightness = 100
  }]
}