// Copyright pkerspe 2026
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestBackupConfigResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
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
  ssl_mode        = "disable"
  username        = "admin"
  password        = "admin"
  include_schemas = ["public"]
  workspace_id    = resource.databasus_workspace.example.id
	
	# needed for proper cleanup after tests, so that TF does not destroy the storage before the database resource
	depends_on = [databasus_storage_local.example]
}

resource "databasus_database_postgresql" "example_2" {
  name            = "my-postgres-db_2"
  database        = "test_db"
  host            = "db"
  port            = 5432
  ssl_mode        = "disable"
  username        = "admin"
  password        = "admin"
  include_schemas = ["public"]
  workspace_id    = resource.databasus_workspace.example.id
	
	# needed for proper cleanup after tests, so that TF does not destroy the storage before the database resource
	depends_on = [databasus_storage_local.example]
}

resource "databasus_database_postgresql" "example_3" {
  name            = "my-postgres-db_3"
  database        = "test_db"
  host            = "db"
  port            = 5432
  ssl_mode        = "disable"
  username        = "admin"
  password        = "admin"
  include_schemas = ["public"]
  workspace_id    = resource.databasus_workspace.example.id
	
	# needed for proper cleanup after tests, so that TF does not destroy the storage before the database resource
	depends_on = [databasus_storage_local.example]
}

resource "databasus_database_postgresql" "example_4" {
  name            = "my-postgres-db_4"
  database        = "test_db"
  host            = "db"
  port            = 5432
  ssl_mode        = "disable"
  username        = "admin"
  password        = "admin"
  include_schemas = ["public"]
  workspace_id    = resource.databasus_workspace.example.id
	
	# needed for proper cleanup after tests, so that TF does not destroy the storage before the database resource
	depends_on = [databasus_storage_local.example]
}

resource "databasus_backup_config" "test" {
  interval              = "DAILY"
  time_of_day           = "08:00"
  retention_policy_type = "COUNT"
  retention_count       = 30
  storage_id            = resource.databasus_storage_local.example.id
  database_id           = resource.databasus_database_postgresql.example.id
}

resource "databasus_backup_config" "test_2" {
  interval              = "CRON"
  cron_expression       = "0 0 * * *"
  retention_policy_type = "COUNT"
  retention_count       = 30
  storage_id            = resource.databasus_storage_local.example.id
  database_id           = resource.databasus_database_postgresql.example_2.id
}

resource "databasus_backup_config" "test_3" {
  interval              = "WEEKLY"
  time_of_day           = "10:00"
  weekday       				= "3"
  retention_policy_type = "COUNT"
  retention_count       = 30
  storage_id            = resource.databasus_storage_local.example.id
  database_id           = resource.databasus_database_postgresql.example_3.id
}

resource "databasus_backup_config" "test_4" {
  interval              = "MONTHLY"
  day_of_month       		= "15"
  retention_policy_type = "COUNT"
  retention_count       = 30
  storage_id            = resource.databasus_storage_local.example.id
  database_id           = resource.databasus_database_postgresql.example_4.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("databasus_backup_config.test", "interval", "DAILY"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "time_of_day", "08:00"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "retention_policy_type", "COUNT"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "retention_count", "30"),
					resource.TestCheckResourceAttrSet("databasus_backup_config.test", "storage_id"),
					resource.TestCheckResourceAttrSet("databasus_backup_config.test", "database_id"),

					// check individual interval configs
					resource.TestCheckResourceAttr("databasus_backup_config.test_2", "cron_expression", "0 0 * * *"),
					resource.TestCheckResourceAttr("databasus_backup_config.test_3", "weekday", "3"),
					resource.TestCheckResourceAttr("databasus_backup_config.test_4", "day_of_month", "15"),

					// check optional default values
					resource.TestCheckResourceAttr("databasus_backup_config.test", "enabled", "true"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "max_failed_retry_count", "0"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "encryption", "true"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "retention_time_period", "MONTH"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "retention_gfs_hours", "24"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "retention_gfs_days", "14"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "retention_gfs_weeks", "4"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "retention_gfs_months", "12"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "retention_gfs_years", "3"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "send_notifications_on_backup_success", "false"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "send_notifications_on_backup_failure", "true"),
				),
			},
		},
	})
}
