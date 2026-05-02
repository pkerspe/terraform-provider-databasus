// Copyright (c) pkerspe
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDatabasePostgresqlResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
resource "databasus_workspace" "test" {
  name = "test-workspace"
}

resource "databasus_database_postgresql" "test" {
  name            = "test-postgres-db"
  database        = "test_db"
  host            = "localhost"
  port            = 5432
  is_https        = false
  username        = "admin"
  password        = "admin"
  include_schemas = ["public"]
  workspace_id    = resource.databasus_workspace.test.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("databasus_database_postgresql.test", "name", "test-postgres-db"),
					resource.TestCheckResourceAttr("databasus_database_postgresql.test", "database", "test_db"),
					resource.TestCheckResourceAttr("databasus_database_postgresql.test", "host", "localhost"),
					resource.TestCheckResourceAttr("databasus_database_postgresql.test", "port", "5432"),
					resource.TestCheckResourceAttr("databasus_database_postgresql.test", "is_https", "false"),

					resource.TestCheckResourceAttrSet("databasus_database_postgresql.test", "username"),
					resource.TestCheckResourceAttrSet("databasus_database_postgresql.test", "password"),

					resource.TestCheckResourceAttr("databasus_database_postgresql.test", "include_schemas.0", "public"),
					resource.TestCheckResourceAttrSet("databasus_database_postgresql.test", "workspace_id"),
				),
			},
		},
	})
}
