terraform {
  required_providers {
    philips = {
      source = "hashicorp.com/edu/philips"
    }
  }
}

variable "application_key" {
  type      = string
  sensitive = true
}

provider "philips" {
  bridge = {
    ip_address      = "192.168.50.209"
    application_key = var.application_key
  }
}

locals {
  bathroom_lights = [
    "00:17:88:01:0c:e6:56:86",
    "00:17:88:01:0c:e3:c7:5d",
    "00:17:88:01:0c:e3:ba:09",
    "00:17:88:01:0c:e6:f4:21",
    "00:17:88:01:0c:e3:d1:2f",
    "00:17:88:01:0c:e3:d1:29",
  ]
  bedroom_overheads = [
    "00:17:88:01:09:4d:9a:29",
    "00:17:88:01:09:51:29:ae"
  ]
  kitchen_lights = [
    philips_light.kitchen_range,
    philips_light.kitchen_shelves,
    philips_light.kitchen_spot_outside,
    philips_light.kitchen_table_spot,
    philips_light.kitchen_shelves,
    philips_light.kitchen_sink_strip,
    philips_light.kitchen_valance,
  ]
  living_room_lights = [
    philips_light.living_room_tv_wall,
    philips_light.tv_strip,
    philips_light.desk_left,
    philips_light.desk_right,
    philips_light.living_room_window_strip,
    philips_light.living_room_orb,
    philips_light.living_room_left,
    philips_light.living_room_kitchen_wall,
    philips_light.living_room_bookshelf,
    philips_light.hallway_overhead_1,
    philips_light.hallway_overhead_2,
  ]
  hallway_lights = [
    philips_light.hallway_lamp,
    philips_light.hallway_strip
  ]
  bedroom_lights = concat(philips_light.bedroom_overhead, [
    philips_light.bedroom_left, philips_light.bedroom_right, philips_light.bedroom_valance
  ])
  all_lights = concat(local.kitchen_lights, local.living_room_lights, local.hallway_lights, local.bedroom_lights, philips_light.bathroom)
}


resource philips_zone "everything" {
  name      = "EVERYTHING"
  archetype = "home"
  light_ids = [for light in local.all_lights : light.id]
}