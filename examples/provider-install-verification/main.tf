terraform {
  required_providers {
    philips-hue = {
      source = "hashicorp.com/edu/philips-hue"
    }
  }
}

provider "philips-hue" {}

resource philips-hue_light "first_light" {
  name = "hallway"
  function = "decorative"
}