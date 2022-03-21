module "cloud_sql" {
  source = "./modules/cloud_sql"
}

# NOTE: 以下３つでmoduleかなあ
resource "random_password" "db_password" {
  length  = 16
  special = false
}
resource "google_sql_user" "user" {
  name     = "appuser"
  instance = module.cloud_sql.instance_name
  password = "${random_password.db_password.result}"
}
output "db_appuser_passwd" {
  value = "${random_password.db_password.result}"
}

# bucket
resource "google_storage_bucket" "bucket" {
  name = "tag-mng"
}

# Cloud PubSub topic for notified-cnt-incrementer
resource "google_pubsub_topic" "topic" {
  name = "notified-cnt-incrementer-topic"
}

# Cloud Scheduler for notified-cnt-incrementer
resource "google_cloud_scheduler_job" "job" {
  name        = "notified-cnt-incrementer-scheduler-job"
  region      = "us-central1"
  description = "notified-cnt-incrementer scheduler job"
  schedule    = "00 00 * * *"
  time_zone   = "Asia/Tokyo"
  pubsub_target {
    # topic.id is the topic's full resource name.
    topic_name = google_pubsub_topic.topic.id
    data       = base64encode("notified-cnt-incrementer-publish!!")
  }
}

# Cloud Functions for notified-cnt-incrementer
## Archive multiple files.
data "archive_file" "notified-cnt-incrementer" {
  type = "zip"
  # source_dir配下に、goのファイル持ってこないと無理っぽい
  source_dir  = "${path.module}/files/notified-cnt-incrementer"
  output_path = "${path.module}/files/notified-cnt-incrementer.zip"
}

resource "google_storage_bucket_object" "archive" {
  name   = "notified-cnt-incrementer.zip"
  bucket = google_storage_bucket.bucket.name
  source = "${path.module}/files/notified-cnt-incrementer.zip"
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
    resource   = google_pubsub_topic.topic.id
  }
  timeout     = 300
  entry_point = "Run"
  labels = {
    my-label = "notified-label"
  }

  environment_variables = {
    DBDSN = "host=/cloudsql/tag-mng-243823:asia-northeast1:tag-mng-cloud dbname=tag-mng user=${google_sql_user.user.name} password=${random_password.db_password.result} sslmode=disable"
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
resource "random_string" "session_key" {
  length  = 32
  special = false
}

## NOTE: apiに変更があった場合は、make buildapiでイメージを更新&GCRへpushする。で、cloud runをdestroy -> applyする
resource "google_cloud_run_service" "api" {
  provider = google-beta
  name     = "tag-mng-api"
  location = "asia-northeast1"

  template {
    spec {
      containers {
        image = "gcr.io/tag-mng-243823/api"
        env {
          name  = "SESSION_KEY"
          value = "${random_string.session_key.result}"
        }
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
        "run.googleapis.com/cloudsql-instances" = "tag-mng-243823:asia-northeast1:${module.cloud_sql.instance_name}"
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

# サイトジェネレーター/メモ公開フラグポーリングプログラム用GCE
resource "google_compute_instance" "hugo-generator" {
  zone         = "us-central1-a"
  name         = "hugo-generator"
  machine_type = "f1-micro"
  boot_disk {
    auto_delete = true
    device_name = "hugo-generator"
    mode        = "READ_WRITE"

    initialize_params {
      image  = "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-10-buster-v20200618"
      labels = {}
      size   = 10
      type   = "pd-standard"
    }
  }
  network_interface {
    network            = "https://www.googleapis.com/compute/v1/projects/tag-mng-243823/global/networks/default"
    network_ip         = "10.128.0.5"
    subnetwork         = "https://www.googleapis.com/compute/v1/projects/tag-mng-243823/regions/us-central1/subnetworks/default"
    subnetwork_project = "tag-mng-243823"

    access_config {
      network_tier = "PREMIUM"
    }
  }
  service_account {
    email = "154979913991-compute@developer.gserviceaccount.com"
    scopes = [
      "https://www.googleapis.com/auth/devstorage.read_write",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring.write",
      "https://www.googleapis.com/auth/service.management.readonly",
      "https://www.googleapis.com/auth/servicecontrol",
      "https://www.googleapis.com/auth/sqlservice.admin",
      "https://www.googleapis.com/auth/trace.append",
    ]
  }
}