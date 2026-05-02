// Copyright (c) pkerspe
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
  is_https        = false
  username        = "admin"
  password        = "admin"
  include_schemas = ["public"]
  workspace_id    = resource.databasus_workspace.example.id
}

resource "databasus_backup_config" "test" {
  interval              = "DAILY"
  time_of_day           = "08:00"
  retention_policy_type = "COUNT"
  retention_count       = 30
  storage_id            = resource.databasus_storage_local.example.id
  database_id           = resource.databasus_database_postgresql.example.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("databasus_backup_config.test", "interval", "DAILY"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "time_of_day", "08:00"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "retention_policy_type", "COUNT"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "retention_count", "30"),
					resource.TestCheckResourceAttrSet("databasus_backup_config.test", "storage_id"),
					resource.TestCheckResourceAttrSet("databasus_backup_config.test", "database_id"),

					// check optional default values
					resource.TestCheckResourceAttr("databasus_backup_config.test", "enabled", "true"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "weekday", "1"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "day_of_month", "1"),
					resource.TestCheckResourceAttr("databasus_backup_config.test", "cron_expression", "0 0 * * *"),
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
