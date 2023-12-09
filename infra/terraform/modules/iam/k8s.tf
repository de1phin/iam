
resource "yandex_container_registry" "iam-docker-registry" {
  name = "iam-docker-registry"
  folder_id = var.folder_id
}

resource "yandex_kubernetes_cluster" "k8s-regional" {
  name = "iam-k8s-cluster"
  description = "k8s cluster for iam services"

  network_id = yandex_vpc_network.default.id
  folder_id = var.folder_id
  
  master {
    public_ip = true

    regional {

      region = "ru-central1"
      location {
        zone      = yandex_vpc_subnet.A.zone
        subnet_id = yandex_vpc_subnet.A.id
      }
      location {
        zone      = yandex_vpc_subnet.B.zone
        subnet_id = yandex_vpc_subnet.B.id
      }
      location {
        zone      = yandex_vpc_subnet.D.zone
        subnet_id = yandex_vpc_subnet.D.id
      }
    }
  }
  service_account_id      = yandex_iam_service_account.iam_instances_sa.id
  node_service_account_id = yandex_iam_service_account.iam_instances_sa.id

  depends_on = [
    yandex_resourcemanager_folder_iam_member.k8s-clusters-agent,
    yandex_resourcemanager_folder_iam_member.vpc-public-admin,
    yandex_resourcemanager_folder_iam_member.images-puller
  ]
}

resource "yandex_resourcemanager_folder_iam_member" "k8s-clusters-agent" {
  # The service account is assigned the k8s.clusters.agent role.
  folder_id = var.folder_id
  role      = "k8s.clusters.agent"
  member    = "serviceAccount:${yandex_iam_service_account.iam_instances_sa.id}"
}

resource "yandex_resourcemanager_folder_iam_member" "vpc-public-admin" {
  # The service account is assigned the vpc.publicAdmin role.
  folder_id = var.folder_id
  role      = "vpc.publicAdmin"
  member    = "serviceAccount:${yandex_iam_service_account.iam_instances_sa.id}"
}

resource "yandex_resourcemanager_folder_iam_member" "images-puller" {
  # The service account is assigned the container-registry.images.puller role.
  folder_id = var.folder_id
  role      = "container-registry.images.puller"
  member    = "serviceAccount:${yandex_iam_service_account.iam_instances_sa.id}"
}

resource "yandex_kubernetes_node_group" "k8s-node-group" {
  cluster_id = yandex_kubernetes_cluster.k8s-regional.id
  name       = "k8s-node-group"

  instance_template {
    platform_id = "standard-v2"
    
    network_interface {
        nat = false
        subnet_ids = [yandex_vpc_subnet.A.id, yandex_vpc_subnet.B.id, yandex_vpc_subnet.D.id]
    }

    resources {
        memory = 4
        cores = 2
    }

    boot_disk {
        type = "network-hdd"
        size = 64
    }

    metadata = {
        ssh-keys = local.bastion_key
    }
  }

  scale_policy {
    fixed_scale {
        size = 3
    }
  }

  allocation_policy {
    location {
      zone = yandex_vpc_subnet.A.zone
    }
    
    location {
      zone = yandex_vpc_subnet.B.zone
    }

    location {
      zone = yandex_vpc_subnet.D.zone
    }
  }
}
