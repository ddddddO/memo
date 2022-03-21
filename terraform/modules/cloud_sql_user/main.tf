resource "random_password" "db_password" {
  length  = 16
  special = false
}
resource "google_sql_user" "user" {
  name     = var.db_user_name
  instance = var.cloud_sql_instance_name
  password = "${random_password.db_password.result}"
}
