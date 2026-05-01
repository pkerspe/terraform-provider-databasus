resource "databasus_workspace" "example" {
  name = "my-workspace-local-storage"
}

resource "databasus_database_postgresql" "example" {
  name            = "my-postgres-db"
  database        = "my-test-db"
  host            = "my-db-host.local"
  port            = 5432
  is_https        = true
  username        = "test-user"
  password        = "test-pwd"
  include_schemas = ["public"]
  workspace_id    = resource.databasus_workspace.example.id
}