terraform {
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
      version = "0.92.0"
    }

    external = {
      source = "hashicorp/external"
      version = "2.1.0"
    }
  }
}


provider "yandex" {
  token = data.external.yc_token.result.iam_token
  endpoint = var.yc_endpoint
  storage_endpoint = var.storage_endpoint
  zone = "ru-central1"
}