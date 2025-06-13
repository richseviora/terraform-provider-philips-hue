terraform {
  required_version = ">= 1.0.0" # Ensure that the Terraform version is 1.0.0 or higher

  required_providers {
    philips = {
      source = "hashicorp.com/edu/philips"
    }
  }
}

locals {
  enabled_act = setsubtract(var.lights, var.lights_to_turn_off)
  disabled_lights = var.lights_to_turn_off
}


resource philips_scene "scene" {
  name  = var.name
  group = var.target
  actions = concat([
    for light in local.enabled_act : {
      target_id         = light.id
      target_type       = light.type
      on                = true
      brightness        = var.light_setting.brightness
      color_temperature = var.light_setting.color_temperature
    }
  ], [
    for light in local.disabled_lights : {
      target_id         = light.id
      target_type       = light.type
      on                = false
      brightness        = var.light_setting.brightness
      color_temperature = var.light_setting.color_temperature
    }
  ])
}