
resource "yandex_iam_service_account" "iam_instances_sa" {
    name = "iam-instances-sa"
    description = "service account for iam services instances"
    folder_id = var.folder_id
}
