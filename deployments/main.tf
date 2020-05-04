provider "google" {
  credentials = "${file("~/.config/gcloud/legacy_credentials/lbfdeatq@gmail.com/adc.json")}"

  project = "tag-mng-243823"
  region  = "asia-northeast1"
  zone    = "asia-northeast1-c"
}

provider "google-beta" {
  credentials = "${file("~/.config/gcloud/legacy_credentials/lbfdeatq@gmail.com/adc.json")}"

  project = "tag-mng-243823"
  region  = "asia-northeast1"
  zone    = "asia-northeast1-c"
}


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

# bucket
resource "google_storage_bucket" "bucket" {
  name = "tag-mng"
}

# Cloud PubSub topic for mail-notificator
resource "google_pubsub_topic" "topic" {
  name = "mail-notificator-topic"
}

# Cloud Scheduler for mail-notificator
resource "google_cloud_scheduler_job" "job" {
  name        = "mail-notificator-scheduler-job"
  region      = "us-central1"
  description = "mail-notificator scheduler job"
  schedule    = "30 9 * * *"
  time_zone   = "Asia/Tokyo"
  pubsub_target {
    # topic.id is the topic's full resource name.
    topic_name = google_pubsub_topic.topic.id
    data       = base64encode("mail-notificator-publish!!")
  }
}

# Cloud Functions for mail-notificator
## Archive multiple files.
data "archive_file" "mail_notificator" {
  type = "zip"
  # source_dir配下に、goのファイル持ってこないと無理っぽい
  source_dir  = "${path.module}/files/mail-notificator"
  output_path = "${path.module}/files/mail-notificator.zip"
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
  region      = "asia-northeast1"
  description = ""
  runtime     = "go113"
  depends_on  = [google_storage_bucket_object.archive]

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  event_trigger {
    event_type = "google.pubsub.topic.publish"
    resource   = google_pubsub_topic.topic.id
  }
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

/* terraform からGAEへvueをデプロイ出来なかった。なので、makeコマンドでGAEへデプロイする。
# GAE for app
## vueをGAEへデプロイする参考：https://cloudpack.media/45462
data "archive_file" "app" {
  type        = "zip"
  source_dir  = "${path.module}/files/app/dist"
  output_path = "${path.module}/files/app/dist.zip"
}

resource "google_storage_bucket_object" "app" {
  name   = "dist.zip"
  bucket = google_storage_bucket.bucket.name
  source = "${path.module}/files/app/dist.zip"
}

resource "google_app_engine_standard_app_version" "app" {
  version_id = "v1"
  service    = "default"
  runtime    = "php55"
  #threadsafe = true

  deployment {
    zip {
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.app.name}"
    }
  }

  // 書き方参考：https://github.com/terraform-providers/terraform-provider-google/issues/5716#issuecomment-590886082
  handlers {
    url_regex = "/"
    static_files {
      path              = "/index.html"
      upload_path_regex = "/index.html"
    }
  }
  handlers {
    url_regex = "/(.*)"
    static_files {
      path              = "/"
      upload_path_regex = "/(.*)"
    }
  }
}
*/

# Cloud Run for api
resource "google_cloud_run_service" "api" {
  provider = google-beta
  name     = "tag-mng-api"
  location = "asia-northeast1"

  template {
    spec {
      containers {
        image = "gcr.io/tag-mng-243823/api"
        env {
          name  = "DBDSN"
          value = "host=/cloudsql/tag-mng-243823:asia-northeast1:tag-mng-cloud dbname=tag-mng user=${google_sql_user.user.name} password=${random_password.db_password.result} sslmode=disable"
        }
      }
    }
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "1000"
        # NOTE: Cloud Run -> Cloud SQLへ接続するために必要
        "run.googleapis.com/cloudsql-instances" = "tag-mng-243823:asia-northeast1:${google_sql_database_instance.instance.name}"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  autogenerate_revision_name = true
}

## FIXME: 以下、一時的に
data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  provider = google-beta
  location = google_cloud_run_service.api.location
  project  = google_cloud_run_service.api.project
  service  = google_cloud_run_service.api.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

output "api_status" {
  value = "${google_cloud_run_service.api.status}"
}
