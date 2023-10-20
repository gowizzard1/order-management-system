provider "google" {
  credentials = file("path/to/your/credentials.json")
  project    = "your-gcp-project-id"
  region     = "us-central1"
}

resource "google_cloud_run_service" "orders_service" {
  name     = "orders-service"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "gcr.io/${google_project.project.project_id}/orders-service:latest"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_project" "project" {}

output "service_url" {
  value = google_cloud_run_service.orders_service.status[0].url
}
