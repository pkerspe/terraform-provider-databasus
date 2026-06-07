// Copyright pkerspe 2026
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestWorkspaceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test reading workspace by name via http
			{
				Config: ProviderConfig + `
resource "databasus_workspace" "example" {
  name = "test-workspace"
}

data "databasus_workspace" "test" {
  id       = databasus_workspace.example.id
  depends_on = [databasus_workspace.example]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.databasus_workspace.test", "id"),
					resource.TestCheckResourceAttr("data.databasus_workspace.test", "name", "test-workspace"),
					resource.TestCheckResourceAttrSet("data.databasus_workspace.test", "created_at"),
				),
			},
			// Test reading workspace by name via https
			{
				Config: ProviderConfigHttps + `
resource "databasus_workspace" "example" {
  name = "test-workspace"
}

data "databasus_workspace" "test" {
  id       = databasus_workspace.example.id
  depends_on = [databasus_workspace.example]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.databasus_workspace.test", "id"),
					resource.TestCheckResourceAttr("data.databasus_workspace.test", "name", "test-workspace"),
					resource.TestCheckResourceAttrSet("data.databasus_workspace.test", "created_at"),
				),
			},
		},
	})
}
