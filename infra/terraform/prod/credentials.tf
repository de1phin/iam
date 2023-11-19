data "external" "yc_token" {
  program = ["bash", "-c", format("yc --profile %s iam create-token --format json | jq '{iam_token, expires_at}'", var.yc_profile)]
}
