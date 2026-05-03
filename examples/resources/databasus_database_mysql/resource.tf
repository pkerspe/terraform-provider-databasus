resource "databasus_workspace" "example" {
  name = "my-workspace-local-storage"
}

resource "databasus_database_mysql" "example" {
  name           = "my-mysql-db"
  database       = "my-test-db"
  host           = "my-db-host.local"
  port           = 3306
  is_https       = true
  username       = "test-user"
  password       = "test-pwd"
  exclude_events = false
  workspace_id   = resource.databasus_workspace.example.id
}
