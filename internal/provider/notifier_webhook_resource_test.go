// Copyright (c) pkerspe
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestNotifierWebhookResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `

resource "databasus_workspace" "test" {
  name = "test-workspace"
}

resource "databasus_notifier_webhook" "test" {
  name           = "test-webhook-notifier"
  body_template  = "{ \"title\": \"{{heading}}\", \"message\": \"{{message}}\" }"
  webhook_method = "POST"
  webhook_url    = "https://localhost:8088/webhooktest"
  headers = {
    Authorization = "Bearer myescuretoken"
		X-User				= "test"
  }
  workspace_id = resource.databasus_workspace.test.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("databasus_notifier_webhook.test", "name", "test-webhook-notifier"),
					resource.TestCheckResourceAttr("databasus_notifier_webhook.test", "body_template", "{ \"title\": \"{{heading}}\", \"message\": \"{{message}}\" }"),
					resource.TestCheckResourceAttr("databasus_notifier_webhook.test", "webhook_method", "POST"),
					resource.TestCheckResourceAttr("databasus_notifier_webhook.test", "webhook_url", "https://localhost:8088/webhooktest"),
					// headers are encoded by databasus, we just check here if the count is correct
					resource.TestCheckResourceAttr("databasus_notifier_webhook.test", "headers.%", "2"),
					resource.TestCheckResourceAttrSet("databasus_notifier_webhook.test", "workspace_id"),
				),
			},
		},
	})
}
