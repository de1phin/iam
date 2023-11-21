
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
      image_id = module.globals.ubuntu_2204_image_id
    }
  }

  # ipv4 network interface
  network_interface {
    subnet_id = yandex_vpc_subnet.A.id
    ipv4 = true
    nat = true
    nat_ip_address = yandex_vpc_address.bastion.external_ipv4_address[0].address
  }

  metadata = {
    ssh-keys = "ubuntu:${file(join(".", [var.iam_bastion_ssh_key_file, "pub"]))}"
  }
}

resource "yandex_dns_recordset" "bastion" {
  name = "bastion.${yandex_dns_zone.iam.zone}"
  zone_id = yandex_dns_zone.iam.id
  type = "A"
  ttl = 600
  data = [ yandex_vpc_address.bastion.external_ipv4_address[0].address ]
}

data "external" "bastion_root_key" {
    program = ["bash", "-c", join(" ", [
               "BASTION_SSH_KEY=${var.iam_bastion_ssh_key_file}",
               "BASTION_HOST=ubuntu@${yandex_dns_recordset.bastion.name}",
               "BASTION_ROOT_SSH_KEY_FILE=/etc/ssh/ssh_host_rsa_key.pub",
               "../modules/iam/scripts/get_bastion_ssh_key.sh"])]
}

locals {
    bastion_key = "ubuntu:${data.external.bastion_root_key.result.key}"
}
