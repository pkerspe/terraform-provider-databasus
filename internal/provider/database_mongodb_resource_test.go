// Copyright (c) pkerspe
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDatabaseMongoDbResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
resource "databasus_workspace" "test" {
  name = "test-workspace"
}

resource "databasus_database_mongodb" "test" {
  name            = "test-mongo-db"
  database        = "test_db"
	auth_database   = "admin"
  host            = "mongodb"
  port            = 27017
  is_https        = false
  username        = "admin"
  password        = "admin"
  workspace_id    = resource.databasus_workspace.test.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("databasus_database_mongodb.test", "name", "test-mongo-db"),
					resource.TestCheckResourceAttr("databasus_database_mongodb.test", "database", "test_db"),
					resource.TestCheckResourceAttr("databasus_database_mongodb.test", "auth_database", "admin"),

					resource.TestCheckResourceAttr("databasus_database_mongodb.test", "host", "mongodb"),
					resource.TestCheckResourceAttr("databasus_database_mongodb.test", "port", "27017"),

					resource.TestCheckResourceAttrSet("databasus_database_mongodb.test", "username"),
					resource.TestCheckResourceAttrSet("databasus_database_mongodb.test", "password"),

					resource.TestCheckResourceAttrSet("databasus_database_mongodb.test", "workspace_id"),
					//test default values
					resource.TestCheckResourceAttr("databasus_database_mongodb.test", "is_direct_connection", "false"),
					resource.TestCheckResourceAttr("databasus_database_mongodb.test", "is_srv", "false"),
					resource.TestCheckResourceAttr("databasus_database_mongodb.test", "cpu_count", "1"),
				),
			},
		},
	})
}
