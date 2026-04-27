terraform {
  required_providers {
    databasus = {
      source = "registry.terraform.io/pkerspe/databasus"
    }
  }
}

provider "databasus" {
  baseurl = "http://localhost:4005/api/v1"
  email   = "admin"
  # NOTE: make sure this secret matches the one in /docker_compose/databasus_test.env if you test against the local docker image from the docker-compose file provided in this lib
  password = "supersecret123"
}

locals {
  workspace_ts = "itest-workspace_${formatdate("YYYY-MM-DD_hh-mm-ss", timestamp())}"
}

resource "databasus_workspace" "itest_generated_workspace" {
  name = local.workspace_ts
}

data "databasus_all_workspaces" "existing_workspaces" {}

data "databasus_workspace" "existing_workspace" {
  id = "c3c57774-920e-4c01-bafd-e27b3b51a0d7"
}

output "all_workspaces" {
  value = data.databasus_all_workspaces.existing_workspaces
}

output "workspace" {
  value = data.databasus_workspace.existing_workspace
}
