output "session_key" {
  value = "${random_string.session_key}"
}

output "api_status" {
  value = "${google_cloud_run_service.api.status}"
}
