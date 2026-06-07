// Copyright pkerspe 2026
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestWorkspaceByNameDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test reading workspace by name via http
			{
				Config: ProviderConfig + `
resource "databasus_workspace" "example" {
  name = "test-workspace-by-name"
}

data "databasus_workspace_by_name" "test" {
  name       = databasus_workspace.example.name
  depends_on = [databasus_workspace.example]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.databasus_workspace_by_name.test", "id"),
					resource.TestCheckResourceAttr("data.databasus_workspace_by_name.test", "name", "test-workspace-by-name"),
					resource.TestCheckResourceAttrSet("data.databasus_workspace_by_name.test", "created_at"),
				),
			},
			// Test reading workspace by name via https
			{
				Config: ProviderConfigHttps + `
resource "databasus_workspace" "example" {
  name = "test-workspace-by-name"
}

data "databasus_workspace_by_name" "test" {
  name       = databasus_workspace.example.name
  depends_on = [databasus_workspace.example]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.databasus_workspace_by_name.test", "id"),
					resource.TestCheckResourceAttr("data.databasus_workspace_by_name.test", "name", "test-workspace-by-name"),
					resource.TestCheckResourceAttrSet("data.databasus_workspace_by_name.test", "created_at"),
				),
			},
		},
	})
}
