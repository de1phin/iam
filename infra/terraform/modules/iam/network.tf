

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
  route_table_id = yandex_vpc_route_table.nat-gateway.id
}

resource "yandex_vpc_subnet" "B" {
  name           = "iam-nets-B"
  network_id     = yandex_vpc_network.default.id
  folder_id = var.folder_id
  v4_cidr_blocks = [
    "10.128.0.0/24"
  ]
  zone = module.globals.zones.B
  route_table_id = yandex_vpc_route_table.nat-gateway.id
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
  route_table_id = yandex_vpc_route_table.nat-gateway.id
}

resource "yandex_vpc_gateway" "nat-gateway" {
  name = "iam-nat-gateway"
  description = "iam gateway to access container-registry from k8s nodes"
  folder_id = var.folder_id

  shared_egress_gateway {}
}

resource "yandex_vpc_route_table" "nat-gateway" {
  name = "nat-gateway-route-table"
  folder_id = var.folder_id
  network_id = yandex_vpc_network.default.id

  static_route {
    destination_prefix = "0.0.0.0/0"
    gateway_id = yandex_vpc_gateway.nat-gateway.id
  }
}

resource "yandex_dns_zone" "iam" {
    name = "iam-dns"
    description = "default IAM dns zone"
    zone = "${var.dns_domain}."
    folder_id = var.folder_id
    public = true
}

resource "yandex_dns_zone" "internal" {
  name = "iam-dns-internal"
  description = "internal IAM dns zone"
  zone = "${var.internal_dns_domain}."
  folder_id = var.folder_id
  public = false
}

resource "yandex_dns_recordset" "nlb" {
  count = length(var.dns_endpoints)
  name = "${var.dns_endpoints[count.index].hostname}.${var.dns_endpoints[count.index].public ? yandex_dns_zone.iam.zone : yandex_dns_zone.internal.zone}"
  zone_id = var.dns_endpoints[count.index].public ? yandex_dns_zone.iam.id : yandex_dns_zone.internal.id
  type = "A"
  ttl = 600
  data = [ "${var.dns_endpoints[count.index].ip}" ]
}