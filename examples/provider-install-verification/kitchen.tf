
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
