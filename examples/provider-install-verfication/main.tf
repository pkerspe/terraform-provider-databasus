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
  id = resource.databasus_workspace.itest_generated_workspace.id
}

data "databasus_users_settings" "current_settings" {
}

resource "databasus_users_settings" "new_settings" {
  allow_external_registrations        = false
  allow_member_invitations            = false
  member_allowed_to_create_workspaces = false
}

resource "databasus_storage_s3" "new_s3_storage" {
  name                        = "itest-s3-storage"
  workspace_id                = resource.databasus_workspace.itest_generated_workspace.id
  is_system                   = true
  s3_access_key               = "SKFHJSKJLHDF-SDFDFDFDF-DFDFDSFD"
  s3_secret_key               = "SECRET-SKFHJSKJLHDF-SDFDFDFDF-DFDFDSFD"
  s3_bucket                   = "bucketname"
  s3_endpoint                 = ""
  s3_prefix                   = ""
  s3_region                   = "eu-west-2"
  s3_storage_class            = ""
  s3_use_virtual_hosted_style = true
  skip_tls_verify             = true
}

resource "databasus_storage_local" "example" {
  name         = "my-local-storage"
  workspace_id = resource.databasus_workspace.itest_generated_workspace.id
}

output "all_workspaces" {
  value = data.databasus_all_workspaces.existing_workspaces
}

output "workspace" {
  value = data.databasus_workspace.existing_workspace
}

output "settings" {
  value = data.databasus_users_settings.current_settings
}
