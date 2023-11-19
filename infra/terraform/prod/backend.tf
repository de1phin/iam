terraform {
  backend "s3" {
    bucket                      = "dpp-iam-terraform"
    endpoint                    = "https://storage.yandexcloud.net"
    key                         = "prod/terraform.tfstate"
    region                      = "us-east-1"
    skip_credentials_validation = true
    skip_metadata_api_check     = true
  }
}