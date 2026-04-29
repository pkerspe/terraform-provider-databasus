resource "databasus_workspace" "example" {
  name = "my-workspace-local-storage"
}

resource "databasus_storage_local" "example" {
  name         = "my-local-storage"
  workspace_id = resource.databasus_workspace.example.id
}