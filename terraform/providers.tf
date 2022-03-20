provider "google" {
  credentials = "${file(var.local_credential_path)}"

  project = var.project
  region  = var.region
  zone    = var.zone
}

provider "google-beta" {
  credentials = "${file(var.local_credential_path)}"

  project = var.project
  region  = var.region
  zone    = var.zone
}