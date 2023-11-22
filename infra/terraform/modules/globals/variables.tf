
variable "zones" {
  type = map(string)
  description = "Yandex Cloud Zone IDs"
  default = {
    A = "ru-central1-a"
    B = "ru-central1-b"
    # ignore zone C as it is out of serivce
    #C = "ru-central1-c"
    D = "ru-central1-d"
  }
}

variable "ubuntu_2204_image_id" {
  type = string
  default = "fd8m3j9ott9u69hks0gg"
}
