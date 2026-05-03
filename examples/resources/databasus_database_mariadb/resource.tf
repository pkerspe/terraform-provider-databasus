resource "databasus_workspace" "example" {
  name = "my-workspace-local-storage"
}

resource "databasus_database_mariadb" "example" {
  name           = "my-maria-db"
  database       = "my-test-db"
  auth_database  = "admin"
  host           = "my-db-host.local"
  port           = 27017
  is_https       = true
  username       = "test-user"
  password       = "test-pwd"
  workspace_id   = resource.databasus_workspace.example.id
}