provider "google" {
    project = var.gcp_project_id
}

data "google_project" "current" {}