# Cloud SQL
resource "google_sql_database" "database" {
  name     = "tag-mng"
  project  = "tag-mng-243823"
  instance = google_sql_database_instance.instance.name
}

resource "google_sql_database_instance" "instance" {
  name             = "tag-mng-cloud"
  database_version = "POSTGRES_11"
  region           = "asia-northeast1"
  settings {
    tier = "db-f1-micro"
    #tier = "db-custom-1-4096"
    maintenance_window {
      day  = 1
      hour = 0
    }
    ip_configuration {
      ipv4_enabled = true
      require_ssl  = false

      authorized_networks {
        name  = "hugo-generator"
        value = "35.226.69.24/32"
      }
    }
  }
}
