
resource "yandex_vpc_address" "bastion" {
    name = "bastion"
    description = "reserved ip address for bastion"
    folder_id = var.folder_id
    external_ipv4_address {
      zone_id = module.globals.zones.A
    }
}

resource "yandex_compute_instance" "bastion" {
  name = "bastion"
  hostname = "bastion"
  folder_id = var.folder_id
  zone = module.globals.zones.A
    
  resources {
    cores = 2
    core_fraction = 20
    memory = 4
  }
  
  boot_disk {
    auto_delete = true
    initialize_params {
      size = 32
      type = "network-hdd"
      image_id = "fd8m3j9ott9u69hks0gg"
    }
  }

  # ipv4 network interface
  network_interface {
    subnet_id = yandex_vpc_subnet.A.id
    nat = true
    nat_ip_address = yandex_vpc_address.bastion.external_ipv4_address[0].address
  }

  metadata = {
    # this only grants access for the creator of the bastion host
    # another ssh keys will have to be added manually
    ssh-keys = "ubuntu:${file("~/.ssh/id_rsa.pub")}"
  }
}

resource "yandex_dns_recordset" "bastion" {
  name = "bastion.${yandex_dns_zone.iam.zone}"
  zone_id = yandex_dns_zone.iam.id
  type = "A"
  ttl = 600
  data = [ yandex_vpc_address.bastion.external_ipv4_address[0].address ]
}
