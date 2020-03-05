# Terraform Deployment for GCP Cloud Functions

This is a demo of using Terraform to deploy multiple Cloud Functions and other related resources. It includes:
* storage bucket to hold zipped code for Cloud Functions
* storage bucket to use as trigger for Cloud Function
* PubSub topic to use as trigger for Cloud Function
* Cloud Function using bucket trigger to publish a PubSub message
* Cloud Function using pubsub trigger to post a message to Slack

## Setup

1. Create a file "terraform.tfvars" in the root folder
2. Add the following variables
    * `project-name` = name of GCP project for relevant resources
    * `gcp-region` = GCP region, e.g. us-central1
    * `creds-file` = file containing GCP credentials for creating resources (typically .json)
    * `slack-url` = Slack webhook URL for posting messages, e.g. "https://hooks.slack.com/services/XXXXXX/YYYYYY/ZZZZZZZZ"
3. Change the names and descriptions of the GCP resources.
4. From the command-line, type `terraform init`
5. Type `terraform plan` to confirm the setup
6. Type `terraform apply` to deploy the resources

## Cloud Functions

Both of the Cloud Functions here use Go.

### Bucket-Trigger function

This function responds to new files uploaded in the trigger bucket. By default, it does not process the file. Additional logic could be added to do so. If the environment variable TOPIC_ID is set, it will send a message on the topic. The SendMessage function shows an example of adding arbitrary attributes to the message, which can be useful to subscribers.

### Pubsub-Function

This function responds to the pubsub topic "demo-topic." It will read the "status" attribute from the message, if available, and generate an emoji to accompany the message text. This message will be sent to the Slack webhook endpoint.