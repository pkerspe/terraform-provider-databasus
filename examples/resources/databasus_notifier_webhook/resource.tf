resource "databasus_workspace" "example" {
  name = "my-workspace-local-storage"
}

resource "databasus_notifier_webhook" "example" {
  name           = "my-webhook-notifier"
  body_template  = "{ \"title\": \"{{heading}}\", \"message\": \"{{message}}\" }"
  webhook_method = "POST"
  webhook_url    = "https://localhost:8088/webhooktest"
  # please note that the headers field is marked as sensitive, yet the values will be stored unencrypted in your TF State file
  headers = {
    Authorization = "Bearer myescuretoken"
  }
  workspace_id = resource.databasus_workspace.example.id
}