
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
  group   = philips_room.bathroom
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
  group   = philips_room.bathroom
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
