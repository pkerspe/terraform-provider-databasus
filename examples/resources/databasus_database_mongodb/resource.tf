resource "databasus_workspace" "example" {
  name = "my-workspace-local-storage"
}

resource "databasus_database_mongodb" "example" {
  name          = "my-mongo-db"
  auth_database = "admin"
  database      = "my-test-db"
  host          = "my-db-host.local"
  port          = 27017
  is_https      = true
  username      = "test-user"
  password      = "test-pwd"
  workspace_id  = resource.databasus_workspace.example.id
}