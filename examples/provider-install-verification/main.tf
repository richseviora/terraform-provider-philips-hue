terraform {
  required_providers {
    philips = {
      source = "hashicorp.com/edu/philips"
    }
  }
}

provider "philips" {
  # output = "STDOUT"
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

resource philips_room bedroom {
  name       = "Bedroom"
  archetype   = "bedroom"
  device_ids = [for l in local.bedroom_lights : l.device_id]
}

resource philips_light bedroom_left {
  name     = "Bed Left"
  function = "decorative"
}

resource philips_light bedroom_right {
  name     = "Bed Right"
  function = "decorative"
}

resource philips_light bedroom_valance {
  name     = "Bedroom Valance"
  function = "decorative"
}

resource philips_light bedroom_overhead {
  count    = 2
  name     = "Bedroom Overhead ${count.index + 1}"
  function = "decorative"
}

module "bedroom_reading" {
  source    = "./modules/room_scenes"
  name      = "Bedroom Reading"
  target    = {
    rid = philips_room.bedroom.id
    rtype = "room"
  }
  light_ids = [for l in local.bedroom_lights : l.id]
  light_ids_to_turn_off = [for l in [philips_light.bedroom_overhead[0], philips_light.bedroom_overhead[1]] : l.id]
  light_setting = {
    brightness        = 100
    color_temperature = 2200
  }
}

import {
  # Name = Kitchen Range
  id = "00:17:88:01:0b:5d:d3:e1"
  to = philips_light.kitchen_range
}

import {
  # Name = Kitchen Shelves
  id = "00:17:88:01:0d:4b:ed:33"
  to = philips_light.kitchen_shelves
}

import {
  # Name = Kitchen Spot Outside
  id = "00:17:88:01:0d:4b:ed:1d"
  to = philips_light.kitchen_spot_outside
}

import {
  # Name = Kitchen Table Spot
  id = "00:17:88:01:0d:4b:ea:05"
  to = philips_light.kitchen_table_spot
}

import {
  # Name = Kitchen Valance
  id = "00:17:88:01:0b:44:c8:42"
  to = philips_light.kitchen_valance
}

import {
  # Name = Sink Strip
  id = "00:17:88:01:09:be:d7:cd"
  to = philips_light.kitchen_sink_strip
}

resource philips_light "kitchen_range" {
  name     = "Kitchen Range"
  function = "functional"
}

resource philips_light "kitchen_shelves" {
  name     = "Kitchen Shelves"
  function = "decorative"
}

resource philips_light "kitchen_sink_strip" {
  name     = "Sink Strip"
  function = "functional"
}
resource philips_light "kitchen_spot_outside" {
  name     = "Kitchen Spot Outside"
  function = "decorative"
}

resource philips_light "kitchen_table_spot" {
  name     = "Kitchen Table Spot"
  function = "decorative"
}

resource philips_light "kitchen_valance" {
  name     = "Kitchen Valance"
  function = "decorative"
}

resource philips_room "kitchen" {
  name       = "Kitchen"
  archetype  = "kitchen"
  device_ids = [for light in local.kitchen_lights : light.device_id]
}

import {
  # Name = Hallway Overhead 1
  id = "00:17:88:01:09:51:2c:0c"
  to = philips_light.hallway_overhead_1
}

import {
  # Name = Hallway Overhead 2
  id = "00:17:88:01:0b:5d:d8:ea"
  to = philips_light.hallway_overhead_2
}

import {
  # Name = Desk Left
  id = "00:17:88:01:08:9a:90:07"
  to = philips_light.desk_left
}


import {
  # Name = Desk Right
  id = "00:17:88:01:08:7c:8d:38"
  to = philips_light.desk_right
}

import {
  # Name = Desk and TV Strip
  id = "00:17:88:01:0b:44:d0:cf"
  to = philips_light.tv_strip
}


import {
  # Name = Living Room Bookshelf
  id = "00:17:88:01:0b:89:d3:d8"
  to = philips_light.living_room_bookshelf
}

import {
  # Name = Living Room Kitchen Wall
  id = "00:17:88:01:09:51:4e:03"
  to = philips_light.living_room_kitchen_wall
}

import {
  # Name = Living Room Left
  id = "00:17:88:01:09:4f:0f:28"
  to = philips_light.living_room_left
}

import {
  # Name = Living Room ORB
  id = "00:17:88:01:08:7c:92:8c"
  to = philips_light.living_room_orb
}

import {
  # Name = Living Room Window Strip
  id = "00:17:88:01:09:bf:18:7b"
  to = philips_light.living_room_window_strip
}

import {
  # Name = TV Wall
  id = "00:17:88:01:09:4d:99:49"
  to = philips_light.living_room_tv_wall
}


resource philips_light "desk_left" {
  name     = "Desk Left"
  function = "decorative"
}

resource philips_light "desk_right" {
  name     = "Desk Right"
  function = "decorative"
}

resource philips_light "tv_strip" {
  name     = "Desk and TV Strip"
  function = "decorative"
}

resource philips_light "hallway_overhead_1" {
  name     = "Hallway Overhead 1"
  function = "decorative"
}

resource philips_light "hallway_overhead_2" {
  name     = "Hallway Overhead 2"
  function = "decorative"
}

resource philips_light "living_room_bookshelf" {
  name     = "Living Room Bookshelf"
  function = "decorative"
}

resource philips_light "living_room_kitchen_wall" {
  name     = "Living Room Kitchen Wall"
  function = "decorative"
}

resource philips_light "living_room_left" {
  name     = "Living Room Left"
  function = "decorative"
}

resource philips_light "living_room_orb" {
  name     = "Living Room ORB"
  function = "functional"
}

resource philips_light "living_room_tv_wall" {
  name     = "Living Room TV Wall"
  function = "decorative"
}

resource philips_light "living_room_window_strip" {
  name     = "Living Room Window Strip"
  function = "decorative"
}

resource philips_room "living_room" {
  name       = "Living Room"
  archetype  = "living_room"
  device_ids = [for light in local.living_room_lights : light.device_id]
}


import {
  # Name = Hallway Lamp
  id = "00:17:88:01:0c:4a:3b:45"
  to = philips_light.hallway_lamp
}

import {
  # Name = Hallway Strip
  id = "00:17:88:01:0c:61:0f:29"
  to = philips_light.hallway_strip
}

resource philips_light "hallway_lamp" {
  name     = "Hallway Lamp"
  function = "decorative"
}

resource philips_light "hallway_strip" {
  name     = "Hallway Strip"
  function = "decorative"
}


resource philips_room "hallway" {
  name       = "Hallway"
  archetype  = "hallway"
  device_ids = [for light in local.hallway_lights : light.device_id]
}


import {
  to = philips_room.bathroom
  id = "0d960eab-68c6-4ed7-8c0d-a24ca756d58e"
}

import {
  to = philips_scene.bathroom_bright
  id = "68b39f81-1c15-4c82-bd0b-ab28606f3d2e"
}


import {
  for_each = local.bathroom_lights
  to       = philips_light.bathroom[each.key]
  id       = each.value
}

import {
  # Name = Bathroom
  id = "00:17:88:01:0c:d7:89:f9"
  to = philips_motion.bathroom
}

resource philips_light "bathroom" {
  count    = 6
  name     = "Bathroom ${count.index + 1}"
  function = "functional"
}


resource philips_room "bathroom" {
  name       = "Bathroom"
  archetype  = "bathroom"
  device_ids = philips_light.bathroom[*].device_id
}

resource philips_motion "bathroom" {
  enabled = true
}


resource philips_scene "bathroom_bright" {
  group   = philips_room.bathroom.reference
  name    = "Bright"
  actions = [
    for light in philips_light.bathroom : {
      target_id         = light.id
      target_type       = "light"
      on                = true
      color_temperature = 2700
      brightness        = 100
    }
  ]
}

resource philips_scene "bathroom_cool" {
  group   = philips_room.bathroom.reference
  name    = "Bathroom Cold"
  actions = [
    for light in philips_light.bathroom : {
      target_id         = light.id
      target_type       = "light"
      on                = true
      color_temperature = 6500
      brightness        = 100
    }
  ]
}

resource philips_zone "everything" {
  name      = "EVERYTHING"
  type      = "home"
  light_ids = [for light in local.all_lights : light.id]
}