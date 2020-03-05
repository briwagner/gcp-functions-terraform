provider "google" {
  project = var.project-name
  region = var.gcp-region
  credentials = file(var.creds-file)
}

locals {
  runtime = "go111"
}

variable "project-name" {}
variable "gcp-region" {}
variable "creds-file" {}
variable "slack-url" {}

resource "google_pubsub_topic" "demo-topic" {
  name = "demo-topic"
}

resource "google_storage_bucket" "demo_functions_trigger_bucket" {
  name = "demo-functions-trigger-bucket"
  location = var.gcp-region
}

resource "google_storage_bucket" "functions_store_bucket" {
  name = "functions-store-bucket"
  location = var.gcp-region
}

data "archive_file" "bucket-trigger-dist" {
  type = "zip"
  source_dir = "./src/bucket_trigger"
  output_path = "dist/demo-bucket-trigger-function.zip"
}

resource "google_storage_bucket_object" "bucket-trigger-archive" {
  name = "bucket-trigger-archive.zip"
  bucket = google_storage_bucket.functions_store_bucket.name
  source = data.archive_file.bucket-trigger-dist.output_path
}

resource "google_cloudfunctions_function" "bucket-trigger-function" {
  name = "demo-bucket-trigger-function"
  description = "Demo function to respond to bucket upload and send pubsub message."
  runtime = local.runtime
  available_memory_mb = 128
  event_trigger {
    event_type = "google.storage.object.finalize"
    resource = google_storage_bucket.demo_functions_trigger_bucket.name
  }
  entry_point = "BucketTrigger"

  source_archive_bucket = google_storage_bucket.functions_store_bucket.name
  source_archive_object = google_storage_bucket_object.bucket-trigger-archive.name

  environment_variables = {
    TOPIC_ID = google_pubsub_topic.demo-topic.name
    PROJECT_ID = var.project-name
  }
}

data "archive_file" "pubsub-function-dist" {
  type = "zip"
  source_dir = "./src/pubsub_function"
  output_path = "dist/demo-pubsub-function.zip"
}

resource "google_storage_bucket_object" "pubsub-function-archive" {
  name = "demo-pubsub-function.zip"
  bucket = google_storage_bucket.functions_store_bucket.name
  source = data.archive_file.pubsub-function-dist.output_path
}

resource "google_cloudfunctions_function" "pubsub-function" {
  name = "demo-pubsub-function"
  description = "Demo function to respond to pubsub function and send to slack."
  runtime = local.runtime
  available_memory_mb = 128
  event_trigger {
    event_type = "google.pubsub.topic.publish"
    resource = google_pubsub_topic.demo-topic.name
  }
  entry_point = "SendMessage"

  source_archive_bucket = google_storage_bucket.functions_store_bucket.name
  source_archive_object = google_storage_bucket_object.pubsub-function-archive.name

  environment_variables = {
    SLACK_URL = var.slack-url
  }
}