terraform {
  required_providers {
    philips-hue = {
      source = "hashicorp.com/edu/philips-hue"
    }
  }
}

provider "philips-hue" {}

import {
  to = philips-hue_room.bathroom
  id = "0d960eab-68c6-4ed7-8c0d-a24ca756d58e"
}

resource philips-hue_light "first_light" {
  name = "Hallway Lamp"
  function = "decorative"
}

resource philips-hue_room "bathroom" {
  name = "Bathroom"
  archetype = "bathroom"
  device_ids = [philips-hue_light.first_light.device_id]
}
