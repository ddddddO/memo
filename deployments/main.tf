provider "google" {
  credentials = "${file("~/.config/gcloud/legacy_credentials/lbfdeatq@gmail.com/adc.json")}"

  project = "tag-mng-243823"
  region  = "us-central1"
  zone    = "us-central1-c"
}

# DB
resource "google_sql_database" "database" {
  name     = "tag-mng"
  project  = "tag-mng-243823"
  instance = google_sql_database_instance.instance.name
}

resource "google_sql_database_instance" "instance" {
  name             = "tag-mng"
  database_version = "POSTGRES_11"
  region           = "us-central1"
  settings {
    tier = "db-f1-micro"
  }
}

resource "random_password" "db_password" {
  length  = 16
  special = false
}

resource "google_sql_user" "user" {
  name     = "appuser"
  instance = google_sql_database_instance.instance.name
  password = random_password.db_password.result
}

output "db_appuser_passwd" {
  value = "${random_password.db_password.result}"
}