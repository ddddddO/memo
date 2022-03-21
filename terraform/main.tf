module "cloud_sql" {
  source = "./modules/cloud_sql"
}

module "cloud_sql_user_app" {
  source = "./modules/cloud_sql_user"

  db_user_name            = "appuser"
  cloud_sql_instance_name = module.cloud_sql.instance_name
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
module "cloud_functions_notified_cnt_incrementer" {
  source = "./modules/cloud_functions"

  source_dir  = "${path.module}/files/notified-cnt-incrementer"
  output_path = "${path.module}/files/notified-cnt-incrementer.zip"

  db_user_name   = module.cloud_sql_user_app.db_user_name
  db_user_passwd = module.cloud_sql_user_app.db_user_passwd
  topic_id       = google_pubsub_topic.topic.id
}

# Cloud Run for api
resource "random_string" "session_key" {
  length  = 32
  special = false
}
module "cloud_run_api" {
  source = "./modules/cloud_run"

  db_user_name            = module.cloud_sql_user_app.db_user_name
  db_user_passwd          = module.cloud_sql_user_app.db_user_passwd
  cloud_sql_instance_name = module.cloud_sql.instance_name
  session_key             = random_string.session_key.result
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