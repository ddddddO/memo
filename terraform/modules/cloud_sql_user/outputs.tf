output "db_user_name" {
  value = google_sql_user.user.name
}
output "db_user_passwd" {
  value = "${random_password.db_password.result}"
}