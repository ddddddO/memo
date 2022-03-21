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
          value = var.session_key
        }
        env {
          name  = "DBDSN"
          value = "host=/cloudsql/tag-mng-243823:asia-northeast1:tag-mng-cloud dbname=tag-mng user=${var.db_user_name} password=${var.db_user_passwd} sslmode=disable"
        }
      }
    }
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "1000"
        # NOTE: Cloud Run -> Cloud SQLへ接続するために必要
        "run.googleapis.com/cloudsql-instances" = "tag-mng-243823:asia-northeast1:${var.cloud_sql_instance_name}"
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
