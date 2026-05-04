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

resource "databasus_health_check_config" "example" {
  database_id                        = resource.databasus_database_postgresql.example.id
  health_check_enabled               = true
  sent_notification_when_unavailable = true
  attempts_before_considered_down    = 1
  store_attempts_days                = 1
  interval_minutes                   = 2
}