variable "name" {
  type = string
  default = "Room"
  description = "The name of the scene to be created. Will also be used when generating scene names."
}

variable "light_ids" {
  type = list(string)
  description = "All the light IDs in the room."
}

variable "light_ids_to_turn_off" {
  type = list(string)
  description = "All the light IDs in the room you'd like to turn off with the scene"
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
    rid   = string
    rtype = string
  })
  description = "The target room or zone for the scene"
}