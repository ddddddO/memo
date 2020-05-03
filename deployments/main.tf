provider "google" {
  credentials = "${file("~/.config/gcloud/legacy_credentials/lbfdeatq@gmail.com/adc.json")}"

  project = "tag-mng-243823"
  region  = "asia-northeast1"
  zone    = "asia-northeast1-c"
}

# DB
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

# Cloud Functions for mail-notificator
## Archive multiple files.
data "archive_file" "mail_notificator" {
  type = "zip"
  # source_dir配下に、goのファイル持ってこないと無理っぽい
  source_dir  = "${path.module}/files/mail-notificator"
  output_path = "${path.module}/files/mail-notificator.zip"
}

resource "google_storage_bucket" "bucket" {
  name = "tag-mng"
}

resource "google_storage_bucket_object" "archive" {
  name   = "mail-notificator.zip"
  bucket = google_storage_bucket.bucket.name
  source = "${path.module}/files/mail-notificator.zip"
}

variable "mail_password" {
  type = string
}

resource "google_cloudfunctions_function" "function" {
  name        = "mail-notificator-function"
  description = ""
  runtime     = "go113"
  depends_on  = [google_storage_bucket_object.archive]

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  trigger_http          = true
  # pubsub作った後で
  #   event_trigger {
  #     #event_type = "providers/cloud.pubsub/eventTypes/topic.publish"
  #     event_type = "google.pubsub.topic.publish"
  #     resource = # pubsub作らないとダメ
  #   }
  timeout     = 300
  entry_point = "Run"
  labels = {
    my-label = "mail-label"
  }

  environment_variables = {
    DBDSN         = "host=/cloudsql/tag-mng-243823:asia-northeast1:tag-mng-cloud dbname=tag-mng user=${google_sql_user.user.name} password=${random_password.db_password.result} sslmode=disable"
    MAIL_PASSWORD = var.mail_password
  }
}
