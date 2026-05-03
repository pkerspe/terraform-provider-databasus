terraform {
  required_providers {
    databasus = {
      source = "registry.terraform.io/pkerspe/databasus"
    }
  }
}

provider "databasus" {
  baseurl = "http://localhost:4005/api/v1"
  email   = "admin"
  # NOTE: make sure this secret matches the one in /docker_compose/databasus_test.env if you test against the local docker image from the docker-compose file provided in this lib
  password = "supersecret123"
}

resource "databasus_workspace" "itest_generated_workspace" {
  name = "itest-workspace"
}

data "databasus_all_workspaces" "existing_workspaces" {}

data "databasus_workspace" "existing_workspace" {
  id = resource.databasus_workspace.itest_generated_workspace.id
}

data "databasus_users_settings" "current_settings" {
}

resource "databasus_users_settings" "new_settings" {
  allow_external_registrations        = false
  allow_member_invitations            = false
  member_allowed_to_create_workspaces = false
}

resource "databasus_storage_s3" "new_s3_storage" {
  name                        = "itest-s3-storage"
  workspace_id                = resource.databasus_workspace.itest_generated_workspace.id
  is_system                   = true
  s3_access_key               = "SKFHJSKJLHDF-SDFDFDFDF-DFDFDSFD"
  s3_secret_key               = "SECRET-SKFHJSKJLHDF-SDFDFDFDF-DFDFDSFD"
  s3_bucket                   = "bucketname"
  s3_endpoint                 = ""
  s3_prefix                   = ""
  s3_region                   = "eu-west-2"
  s3_storage_class            = ""
  s3_use_virtual_hosted_style = true
  skip_tls_verify             = true
}

resource "databasus_storage_local" "example" {
  name         = "my-local-storage"
  workspace_id = resource.databasus_workspace.itest_generated_workspace.id
}

resource "databasus_database_mariadb" "example" {
  name         = "my-maria-db"
  database     = "test_db"
  host         = "mariadb"
  port         = 3306
  is_https     = false
  username     = "admin"
  password     = "admin"
  workspace_id = resource.databasus_workspace.itest_generated_workspace.id
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
  workspace_id    = resource.databasus_workspace.itest_generated_workspace.id
}

resource "databasus_database_postgresql" "example-2" {
  name            = "my-postgres-db-2"
  database        = "test_db"
  host            = "db"
  port            = 5432
  is_https        = false
  username        = "admin"
  password        = "admin"
  include_schemas = ["public"]
  workspace_id    = resource.databasus_workspace.itest_generated_workspace.id
}

resource "databasus_notifier_webhook" "example" {
  name           = "my-webhook-notifier"
  body_template  = "{ \"title\": \"{{heading}}\", \"message\": \"{{message}}\" }"
  webhook_method = "POST"
  webhook_url    = "https://localhost:8088/webhooktest"
  headers = {
    Authorization = "Bearer myescuretoken"
  }
  workspace_id = resource.databasus_workspace.itest_generated_workspace.id
}

resource "databasus_backup_config" "example" {
  enabled                              = true
  interval                             = "DAILY"
  time_of_day                          = "08:00"
  weekday                              = 1
  day_of_month                         = 1
  cron_expression                      = "0 0 * * *"
  max_failed_retry_count               = 0
  encryption                           = false
  retention_policy_type                = "COUNT"
  retention_count                      = 30
  retention_time_period                = "WEEK"
  retention_gfs_hours                  = 24
  retention_gfs_days                   = 7
  retention_gfs_weeks                  = 14
  retention_gfs_months                 = 12
  retention_gfs_years                  = 3
  send_notifications_on_backup_success = true
  send_notifications_on_backup_failure = true
  storage_id                           = resource.databasus_storage_local.example.id
  database_id                          = resource.databasus_database_postgresql.example.id
}

resource "databasus_backup_config" "example-2" {
  enabled               = true
  interval              = "DAILY"
  time_of_day           = "08:00"
  retention_policy_type = "COUNT"
  storage_id            = resource.databasus_storage_local.example.id
  database_id           = resource.databasus_database_postgresql.example-2.id
}

output "all_workspaces" {
  value = data.databasus_all_workspaces.existing_workspaces
}

output "workspace" {
  value = data.databasus_workspace.existing_workspace
}

output "settings" {
  value = data.databasus_users_settings.current_settings
}

output "database" {
  value     = resource.databasus_database_postgresql.example
  sensitive = true
}

output "notifier-id" {
  value = resource.databasus_notifier_webhook.example.id
  //sensitive = true
}