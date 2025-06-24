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

resource "philips_light" "hallway_lamp" {
  name     = "Hallway Lamp"
  function = "decorative"
}

resource "philips_light" "hallway_strip" {
  name     = "Hallway Strip"
  function = "decorative"
}


resource "philips_room" "hallway" {
  name       = "Hallway"
  archetype  = "hallway"
  device_ids = [for light in local.hallway_lights : light.device_id]
}
