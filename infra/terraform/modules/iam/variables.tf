variable "folder_id" {
    type = string
    description = "Folder ID for resources to be created in"
    default = ""
}

variable "dns_domain" {
    type = string
    default = "iam.de1phin.ru"
}

variable "iam_bastion_ssh_key_file" {
  type = string
  default = "~/.ssh/iam_bastion"
}
