locals {
    terraform_config = <<EOF
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
EOF

    yandex_provider_config = <<EOF
provider "yandex" {
  token = data.external.yc_token.result.iam_token
  endpoint = var.yc_endpoint
  storage_endpoint = var.storage_endpoint
  zone = "ru-central1"
}
EOF
}

remote_state {
    backend  = "s3"
    generate = {
        path      = "backend.tf"
        if_exists = "overwrite"
    }
    config = {
        endpoint = "https://storage.yandexcloud.net"
        bucket   = "dpp-iam-terraform"
        key      = "${path_relative_to_include()}/terraform.tfstate"
        region   = "us-east-1"

        skip_credentials_validation = true
        skip_metadata_api_check     = true
    }
}

generate "variables" {
    path      = "variables.tf"
    if_exists = "overwrite"
    contents  = <<EOF

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
EOF
}

generate "provider" {
    path      = "provider.tf"
    if_exists = "overwrite"
    contents  = join("\n\n", [local.terraform_config, local.yandex_provider_config])
}

generate "credentials" {
    path      = "credentials.tf"
    if_exists = "overwrite"
    contents  = <<EOF
data "external" "yc_token" {
  program = ["bash", "-c", format("yc --profile %s iam create-token --format json | jq '{iam_token, expires_at}'", var.yc_profile)]
}
EOF
}
