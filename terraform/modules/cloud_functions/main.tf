## bucket
resource "google_storage_bucket" "bucket" {
  name = "tag-mng"
}
## Archive multiple files.
data "archive_file" "notified-cnt-incrementer" {
  type = "zip"
  # source_dir配下に、goのファイル持ってこないと無理っぽい
  source_dir  = var.source_dir
  output_path = var.output_path
}
resource "google_storage_bucket_object" "archive" {
  name   = "notified-cnt-incrementer.zip"
  bucket = google_storage_bucket.bucket.name
  source = var.output_path
}
resource "google_cloudfunctions_function" "function" {
  name        = "notified-cnt-incrementer-function"
  region      = "asia-northeast1"
  description = ""
  runtime     = "go116"
  depends_on  = [google_storage_bucket_object.archive]

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  event_trigger {
    event_type = "google.pubsub.topic.publish"
    resource   = var.topic_id
  }
  timeout     = 300
  entry_point = "Run"
  labels = {
    my-label = "notified-label"
  }

  environment_variables = {
    DBDSN = "host=/cloudsql/tag-mng-243823:asia-northeast1:tag-mng-cloud dbname=tag-mng user=${var.db_user_name} password=${var.db_user_passwd} sslmode=disable"
  }
}