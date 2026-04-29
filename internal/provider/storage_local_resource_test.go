// Copyright (c) pkerspe
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestStorageLocalResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
resource "databasus_storage_local" "test" {
	name = "test_local_storage"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("databasus_storage_local.test", "name", "test_workspace"),
					resource.TestCheckResourceAttrSet("databasus_storage_local.test", "id"),
					// TODO: add more checks and also tests for resource update
				),
			},
		},
	})
}
