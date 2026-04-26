// Copyright (c) KerspeP
// SPDX-License-Identifier: Apache-2.0

package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pkerspe/terraform-provider-databasus/internal/provider"
)

func TestAccAllWorkspacesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccAllWorkspacesDataSourceConfig,
			},
		},
	})
}

const testAccAllWorkspacesDataSourceConfig = `
data "databasus_all_workspaces" "test" {
}
`
