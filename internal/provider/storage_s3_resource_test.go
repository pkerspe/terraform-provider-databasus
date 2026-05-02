// Copyright (c) pkerspe
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestS3StorageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `

resource "databasus_workspace" "test" {
  name = "test-workspace"
}

resource "databasus_storage_s3" "test" {
  name                        = "test-s3-storage"
  workspace_id                = resource.databasus_workspace.test.id
  is_system                   = true
  s3_access_key               = "YOUR S3 ACCESS KEY"
  s3_secret_key               = "YOUR S3 SECRET KEY"
  s3_bucket                   = "testbucket"
  s3_endpoint                 = "http://localhost"
  s3_prefix                   = ""
  s3_region                   = "test-region"
  s3_storage_class            = ""
  s3_use_virtual_hosted_style = true
  skip_tls_verify             = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("databasus_storage_s3.test", "name", "test-s3-storage"),
					resource.TestCheckResourceAttrSet("databasus_storage_s3.test", "workspace_id"),
					resource.TestCheckResourceAttr("databasus_storage_s3.test", "is_system", "true"),
					resource.TestCheckResourceAttr("databasus_storage_s3.test", "s3_bucket", "testbucket"),
					resource.TestCheckResourceAttr("databasus_storage_s3.test", "s3_endpoint", "http://localhost"),
					resource.TestCheckResourceAttr("databasus_storage_s3.test", "s3_prefix", ""),
					resource.TestCheckResourceAttr("databasus_storage_s3.test", "s3_region", "test-region"),
					resource.TestCheckResourceAttr("databasus_storage_s3.test", "s3_storage_class", ""),
					resource.TestCheckResourceAttr("databasus_storage_s3.test", "s3_use_virtual_hosted_style", "true"),
					resource.TestCheckResourceAttr("databasus_storage_s3.test", "skip_tls_verify", "true"),
				),
			},
		},
	})
}
