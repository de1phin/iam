
resource "yandex_mdb_postgresql_cluster" "shared" {
  name                = "shared"
  environment         = "PRODUCTION"
  network_id          = yandex_vpc_network.default.id
  folder_id = var.folder_id
  deletion_protection = true

  config {
    version = "15"
    resources {
      resource_preset_id = "s2.micro"
      disk_type_id       = "network-hdd"
      disk_size          = 64
    }
  }

  host {
    zone             = yandex_vpc_subnet.A.zone
    name             = "iam-postgresql"
    subnet_id        = yandex_vpc_subnet.A.id
    assign_public_ip = false
  }
}

data "external" "password" {
    count   = length(var.database)
    program = ["bash", "-c", join(" ", [
               "DATABASE_NAME=${var.database[count.index].dbname}",
               "../modules/iam/scripts/get_db_password.sh"])]
}

resource "yandex_mdb_postgresql_database" "psql_database" {
  count      = length(var.database)
  cluster_id = yandex_mdb_postgresql_cluster.shared.id
  name       = var.database[count.index].dbname
  owner      = yandex_mdb_postgresql_user.psql_user[count.index].name
  depends_on = [
    yandex_mdb_postgresql_user.psql_user
  ]
}

resource "yandex_mdb_postgresql_user" "psql_user" {
  count      = length(var.database)
  cluster_id = yandex_mdb_postgresql_cluster.shared.id
  name       = var.database[count.index].user
  password   = data.external.password[count.index].result.password
}

