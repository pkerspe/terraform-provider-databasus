resource "databasus_workspace" "example" {
  name = "my-workspace"
}

resource "databasus_storage_local" "example" {
  name         = "my-local-storage"
  workspace_id = resource.databasus_workspace.example.id
}

resource "databasus_database_postgresql" "example" {
  name            = "my-postgres-db"
  database        = "test_db"
  host            = "db"
  port            = 5432
  is_https        = false
  username        = "admin"
  password        = "admin"
  include_schemas = ["public"]
  workspace_id    = resource.databasus_workspace.example.id
}

resource "databasus_backup_config" "example" {
  enabled               = true
  interval              = "DAILY"
  time_of_day           = "08:00"
  retention_policy_type = "COUNT"
  retention_count       = 30
  storage_id            = resource.databasus_storage_local.example.id
  database_id           = resource.databasus_database_postgresql.example.id
}