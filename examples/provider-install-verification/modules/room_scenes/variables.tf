variable "name" {
  type = string
  default = "Room"
  description = "The name of the scene to be created. Will also be used when generating scene names."
}

variable "lights" {
  type = list(object({
    id = string
    type = string
  }))
  description = "All the lights in the room."
}

variable "lights_to_turn_off" {
  type = list(object({
    id = string
    type = string
  }))
  description = "All the lights in the room you'd like to turn off with the scene"
}

variable "light_setting" {
  type = object({
    brightness = number
    color_temperature = number
  })
  description = "the light settings to apply"
}

variable "target" {
  type = object({
    id   = string
    type = string
  })
  description = "The target room or zone for the scene"
}