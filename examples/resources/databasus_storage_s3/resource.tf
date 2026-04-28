resource "databasus_workspace" "example" {
  name = "my-workspace"
}

resource "databasus_storage_s3" "example" {
  name                        = "my-s3-storage"
  workspace_id                = resource.databasus_workspace.example.id
  is_system                   = true
  s3_access_key               = "YOUR S3 ACCESS KEY"
  s3_secret_key               = "YOUR S3 SECRET KEY"
  s3_bucket                   = "bucketname"
  s3_endpoint                 = "<your endpoint here or blank for standard AWS S3>"
  s3_prefix                   = ""
  s3_region                   = "your-region"
  s3_storage_class            = ""
  s3_use_virtual_hosted_style = true
  skip_tls_verify             = false
}