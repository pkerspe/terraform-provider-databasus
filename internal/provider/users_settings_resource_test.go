// Copyright (c) pkerspe
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUsersSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
resource "databasus_users_settings" "test" {
  allow_external_registrations        = false
  allow_member_invitations            = false
  member_allowed_to_create_workspaces = false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("databasus_users_settings.test", "allow_external_registrations", "false"),
					resource.TestCheckResourceAttr("databasus_users_settings.test", "allow_member_invitations", "false"),
					resource.TestCheckResourceAttr("databasus_users_settings.test", "member_allowed_to_create_workspaces", "false"),
				),
			},
			// Test update
			{
				Config: ProviderConfig + `
resource "databasus_users_settings" "test" {
  allow_external_registrations        = true
  allow_member_invitations            = true
  member_allowed_to_create_workspaces = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("databasus_users_settings.test", "allow_external_registrations", "true"),
					resource.TestCheckResourceAttr("databasus_users_settings.test", "allow_member_invitations", "true"),
					resource.TestCheckResourceAttr("databasus_users_settings.test", "member_allowed_to_create_workspaces", "true"),
				),
			},
		},
	})
}
