// Copyright kerspep
// SPDX-License-Identifier: MPL-2.0

package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/pkerspe/terraform-provider-databasus/internal/provider"
)

func TestAccWorkspaceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccWorkspaceDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.databasus_workspace.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("example-id"),
					),
				},
			},
		},
	})
}

const testAccWorkspaceDataSourceConfig = `
data "databasus_workspace" "test" {
  id = "example-id"
}
`
