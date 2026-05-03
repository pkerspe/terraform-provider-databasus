A Provider to manage your Databasus (<https://databasus.com/>) configuration using Terraform.

This provider is still in an early stage and currently supports a reduced set of resources only.

Additional resources will be added over time, but in a demand drive way due to limited development capacity.
If you need a specific resource please open a feature request ticket in the github repository: <https://github.com/pkerspe/terraform-provider-databasus/issues>

Please note that this provider is not developed by the Databasus team and not supported by them, in fact not even encouraged to use an IaC approach for configuration since it is not in line with their goals for the tool. As such, provider support could break with any new release of Databasus.

To limit risk please make sure to match the Provider release version with the Databasus version you are using.

Provider Support Map:

| Provider Version | Tested Databasus Version |
|------------------|--------------------------|
| v0.1.x - v0.5.x  | v3.32.2                  |

## Getting Started

To configure a backup for a Database you need to configure at least 4 resources:

1.) The workspace where the database and backup lives in

2.) the storage where the backups are to be stored

3.) the Database configuration itself of the database to be backed up

4.) the Backup Configuration that links all the pieces together and defines the schedule for backup execution and retention settings

Here is an example for a minimal setup using a local storage and a minimalist backup configuration mostly using the default values:

````
# create workspace where all configuration lives in
resource "databasus_workspace" "example" {
  name = "my-workspace"
}

# create a simple local storage (storing files on the databasus host directly)
resource "databasus_storage_local" "example" {
  name         = "my-local-storage"
  workspace_id = resource.databasus_workspace.example.id
}

# configuring the database to be backed up
resource "databasus_database_postgresql" "example" {
  name            = "my-postgres-db"
  database        = "test_db"
  host            = "db"
  port            = 5432
  is_https        = false
  username        = "admin"
  password        = "admin"
  include_schemas = ["public"]
  workspace_id    = resource.databasus_workspace.example.id
}

# creating the actual backup configuration with backup settings and linking the database and storage together
resource "databasus_backup_config" "example" {
  interval              = "DAILY"
  time_of_day           = "08:00"
  retention_policy_type = "COUNT"
  retention_count       = 30
  storage_id            = resource.databasus_storage_local.example.id
  database_id           = resource.databasus_database_postgresql.example.id
}
````
