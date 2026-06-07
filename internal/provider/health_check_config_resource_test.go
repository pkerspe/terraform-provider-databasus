// Copyright pkerspe 2026
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	HealthCheckBaseConfig = `
resource "databasus_workspace" "test" {
  name = "test-workspace"
}

resource "databasus_database_postgresql" "test" {
  name            = "test-postgres-db"
  database        = "test_db"
  host            = "db"
  port            = 5432
  is_https        = false
  username        = "admin"
  password        = "admin"
  include_schemas = ["public"]
  workspace_id    = resource.databasus_workspace.test.id
}
resource "databasus_database_postgresql" "test_2" {
  name            = "test-postgres-db"
  database        = "test_db"
  host            = "db"
  port            = 5432
  is_https        = false
  username        = "admin"
  password        = "admin"
  include_schemas = ["public"]
  workspace_id    = resource.databasus_workspace.test.id
}
resource "databasus_database_postgresql" "test_3" {
  name            = "test-postgres-db"
  database        = "test_db"
  host            = "db"
  port            = 5432
  is_https        = false
  username        = "admin"
  password        = "admin"
  include_schemas = ["public"]
  workspace_id    = resource.databasus_workspace.test.id
}
`
)

func TestHealthCheckConfigResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + HealthCheckBaseConfig + `
resource "databasus_health_check_config" "test" {
  database_id           = resource.databasus_database_postgresql.test.id
  health_check_enabled	= true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("databasus_health_check_config.test", "database_id"),
					resource.TestCheckResourceAttr("databasus_health_check_config.test", "health_check_enabled", "true"),
					// check default values
					resource.TestCheckResourceAttr("databasus_health_check_config.test", "sent_notification_when_unavailable", "false"),
					resource.TestCheckResourceAttr("databasus_health_check_config.test", "attempts_before_considered_down", "3"),
					resource.TestCheckResourceAttr("databasus_health_check_config.test", "store_attempts_days", "7"),
					resource.TestCheckResourceAttr("databasus_health_check_config.test", "interval_minutes", "1"),
				),
			},
			{
				Config: ProviderConfig + HealthCheckBaseConfig + `
resource "databasus_health_check_config" "test_2" {
  	database_id           				= resource.databasus_database_postgresql.test_2.id
  	health_check_enabled				= true
	sent_notification_when_unavailable 	= true
	attempts_before_considered_down 	= 1
	store_attempts_days 				= 1
	interval_minutes 					= 2
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("databasus_health_check_config.test_2", "database_id"),
					resource.TestCheckResourceAttr("databasus_health_check_config.test_2", "health_check_enabled", "true"),
					resource.TestCheckResourceAttr("databasus_health_check_config.test_2", "sent_notification_when_unavailable", "true"),
					resource.TestCheckResourceAttr("databasus_health_check_config.test_2", "attempts_before_considered_down", "1"),
					resource.TestCheckResourceAttr("databasus_health_check_config.test_2", "store_attempts_days", "1"),
					resource.TestCheckResourceAttr("databasus_health_check_config.test_2", "interval_minutes", "2"),
				),
			},
			{
				Config: ProviderConfig + HealthCheckBaseConfig + `
resource "databasus_health_check_config" "test_3" {
  database_id           = resource.databasus_database_postgresql.test_3.id
  health_check_enabled	= false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("databasus_health_check_config.test_3", "database_id"),
					resource.TestCheckResourceAttr("databasus_health_check_config.test_3", "health_check_enabled", "false"),
				),
			},
		},
	})
}
