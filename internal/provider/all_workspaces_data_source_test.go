// Copyright pkerspe 2026
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAllWorkspacesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test reading workspaces via http
			{
				Config: ProviderConfig + `
resource "databasus_workspace" "example" {
  name = "my-workspace"
}

data "databasus_all_workspaces" "test" {
  depends_on = [databasus_workspace.example]
}

output "workspace_names" {
  value = [for w in data.databasus_all_workspaces.test.workspaces : w.name]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.databasus_all_workspaces.test", "workspaces.0.name", "my-workspace"),
				),
			},
			// Test reading workspaces via https
			{
				Config: ProviderConfigHttps + `
resource "databasus_workspace" "example" {
  name = "my-workspace"
}

data "databasus_all_workspaces" "test" {
  depends_on = [databasus_workspace.example]
}

output "workspace_names" {
  value = [for w in data.databasus_all_workspaces.test.workspaces : w.name]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.databasus_all_workspaces.test", "workspaces.0.name", "my-workspace"),
				),
			},
		},
	})
}
