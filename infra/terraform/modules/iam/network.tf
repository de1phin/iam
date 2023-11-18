

resource "yandex_vpc_network" "default" {
    name = "default-network"
    description = "Default network for IAM"
    folder_id = var.folder_id
}

resource "yandex_vpc_subnet" "A" {
  name           = "iam-nets-A"
  network_id     = yandex_vpc_network.default.id
  folder_id = var.folder_id
  v4_cidr_blocks = [
    "10.127.0.0/24"
  ]
  zone = module.globals.zones.A
}

resource "yandex_vpc_subnet" "B" {
  name           = "iam-nets-B"
  network_id     = yandex_vpc_network.default.id
  folder_id = var.folder_id
  v4_cidr_blocks = [
    "10.128.0.0/24"
  ]
  zone = module.globals.zones.B
}

# ignore zone C as it is out of service

resource "yandex_vpc_subnet" "D" {
  name           = "iam-nets-D"
  network_id     = yandex_vpc_network.default.id
  folder_id = var.folder_id
  v4_cidr_blocks = [
    "10.129.0.0/24"
  ]
  zone = module.globals.zones.D
}

resource "yandex_dns_zone" "iam" {
    name = "iam-dns"
    description = "default IAM dns zone"
    zone = "${var.dns_domain}."
    folder_id = var.folder_id
    public = true
}
