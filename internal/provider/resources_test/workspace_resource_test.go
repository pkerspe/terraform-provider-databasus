// Copyright kerspep
// SPDX-License-Identifier: MPL-2.0

package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/pkerspe/terraform-provider-databasus/internal/provider"
)

func TestAccWorkspaceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWorkspaceResourceConfig("one"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"databasus_workspace.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("example-id"),
					),
					statecheck.ExpectKnownValue(
						"databasus_workspace.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("one"),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "databasus_workspace.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"name", "defaulted"},
			},
			// Update and Read testing
			{
				Config: testAccWorkspaceResourceConfig("two"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"databasus_workspace.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("example-id"),
					),
					statecheck.ExpectKnownValue(
						"databasus_workspace.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("two"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccWorkspaceResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "databasus_workspace" "test" {
  name = %[1]q
}
`, name)
}
