

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
  function = "decorative"
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
