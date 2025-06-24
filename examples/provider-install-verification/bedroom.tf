import {
  for_each = local.bedroom_overheads
  to       = philips_light.bedroom_overhead[each.key]
  id       = each.value
}

import {
  # Name = Bed Left
  id = "00:17:88:01:09:4f:0e:7f"
  to = philips_light.bedroom_left
}

import {
  # Name = Bed Right
  id = "00:17:88:01:09:4f:0a:b9"
  to = philips_light.bedroom_right
}

import {
  # Name = Bedroom Valance
  id = "00:17:88:01:09:b7:f9:99"
  to = philips_light.bedroom_valance
}

resource "philips_room" "bedroom" {
  name       = "Bedroom"
  archetype  = "bedroom"
  device_ids = [for l in local.bedroom_lights : l.device_id]
}

resource "philips_light" "bedroom_left" {
  name     = "Bed Left"
  function = "decorative"
}

resource "philips_light" "bedroom_right" {
  name     = "Bed Right"
  function = "decorative"
}

resource "philips_light" "bedroom_valance" {
  name     = "Bedroom Valance"
  function = "decorative"
}

resource "philips_light" "bedroom_overhead" {
  count    = 2
  name     = "Bedroom Overhead ${count.index + 1}"
  function = "decorative"
}

module "bedroom_reading" {
  source             = "./modules/room_scenes"
  name               = "Bedroom Reading"
  target             = philips_room.bedroom
  lights             = local.bedroom_lights
  lights_to_turn_off = [philips_light.bedroom_overhead[0], philips_light.bedroom_overhead[1]]
  light_setting = {
    brightness        = 100
    color_temperature = 2200
  }
}

module "bedroom_day" {
  source = "./modules/room_scenes"
  name   = "Bedroom Day"
  target = philips_room.bedroom
  lights = local.bedroom_lights
  lights_to_turn_off = [
    philips_light.bedroom_overhead[0], philips_light.bedroom_overhead[1], philips_light.bedroom_right,
    philips_light.bedroom_left
  ]
  light_setting = {
    brightness        = 100
    color_temperature = 6000
  }
}