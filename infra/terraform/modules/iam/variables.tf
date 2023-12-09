variable "folder_id" {
    type = string
    description = "Folder ID for resources to be created in"
    default = ""
}

variable "dns_domain" {
    type = string
}

variable "internal_dns_domain" {
  type = string
}

variable "iam_bastion_ssh_key_file" {
  type = string
  default = "~/.ssh/iam_bastion"
}

variable "database" {
  type = list(map(string))
}

variable "dns_endpoints" {
  type = list(map(string))
}
