// Copyright (c) pkerspe
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDatabaseMySqlResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
resource "databasus_workspace" "test" {
  name = "test-workspace"
}

resource "databasus_database_mysql" "test" {
  name            = "test-mysql-db"
  database        = "test_db"
  host            = "mysql"
  port            = 3306
  is_https        = false
  username        = "admin"
  password        = "admin"
  workspace_id    = resource.databasus_workspace.test.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("databasus_database_mysql.test", "name", "test-mysql-db"),
					resource.TestCheckResourceAttr("databasus_database_mysql.test", "database", "test_db"),
					resource.TestCheckResourceAttr("databasus_database_mysql.test", "host", "mysql"),
					resource.TestCheckResourceAttr("databasus_database_mysql.test", "port", "3306"),

					resource.TestCheckResourceAttrSet("databasus_database_mysql.test", "username"),
					resource.TestCheckResourceAttrSet("databasus_database_mysql.test", "password"),

					resource.TestCheckResourceAttrSet("databasus_database_mysql.test", "workspace_id"),
				),
			},
		},
	})
}
