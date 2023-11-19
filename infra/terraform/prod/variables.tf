
variable "yc_profile" {
    type = string
    description = "YC CLI profile"
    default = "dpp"
}

variable "yc_endpoint" {
  type = string
  description = "Endpoint for public API (for yandex provider)"
  default = "api.cloud.yandex.net:443"
}

variable "storage_endpoint" {
  type = string
  description = "Storage endpoint for yandex_storage_* resources"
  default = "https://storage.yandexcloud.net"
}
