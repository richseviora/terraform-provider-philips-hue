terraform {
  required_providers {
    philips = {
      source = "hashicorp.com/edu/philips-hue"
    }
  }
}

provider "philips" {
  output_imports = true
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
  ]
  living_room_lights = [
    philips_light.living_room_tv_wall,
    philips_light.desk_left,
    philips_light.desk_right,
    philips_light.living_room_window_strip,
    philips_light.living_room_orb,
    philips_light.living_room_left,
    philips_light.living_room_kitchen_wall,
    philips_light.living_room_bookshelf,
  ]
  hallway_lights = [
    philips_light.hallway_lamp,
    philips_light.hallway_strip
  ]
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

resource philips_light bedroom_left {
  name = "Bed Left"
}

resource philips_light bedroom_right {
  name = "Bed Right"
}

resource philips_light bedroom_valance {
  name = "Bedroom Valance"
}

resource philips_light bedroom_overhead {
  count    = 2
  name     = "Bedroom Overhead ${count.index + 1}"
  function = "decorative"
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
  name = "Kitchen Range"
}

resource philips_light "kitchen_shelves" {
  name = "Kitchen Shelves"
}

resource philips_light "kitchen_sink_strip" {
  name = "Sink Strip"
}
resource philips_light "kitchen_spot_outside" {
  name = "Kitchen Spot Outside"
}

resource philips_light "kitchen_table_spot" {
  name = "Kitchen Table Spot"
}

resource philips_light "kitchen_valance" {
  name = "Kitchen Valance"
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
  name = "Desk Left"
}

resource philips_light "desk_right" {
  name = "Desk Right"
}

resource philips_light "tv_strip" {
  name = "Desk and TV Strip"
}

resource philips_light "hallway_overhead_1" {
  name = "Hallway Overhead 1"
}

resource philips_light "hallway_overhead_2" {
  name = "Hallway Overhead 2"
}

resource philips_light "living_room_bookshelf" {
  name = "Living Room Bookshelf"
}

resource philips_light "living_room_kitchen_wall" {
  name = "Living Room Kitchen Wall"
}

resource philips_light "living_room_left" {
  name = "Living Room Left"
}

resource philips_light "living_room_orb" {
  name = "Living Room ORB"
}

resource philips_light "living_room_tv_wall" {
  name = "Living Room TV Wall"
}

resource philips_light "living_room_window_strip" {
  name = "Living Room Window Strip"
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
  name = "Hallway Lamp"
}

resource philips_light "hallway_strip" {
  name = "Hallway Strip"
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
  for_each = local.bedroom_overheads
  to       = philips_light.bedroom_overhead[each.key]
  id       = each.value
}

import {
  for_each = local.bathroom_lights
  to       = philips_light.bathroom[each.key]
  id       = each.value
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
  light_ids = [for light in philips_light.bathroom : light.id]
}