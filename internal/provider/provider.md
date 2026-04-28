A Provider to manage your Databasus (<https://databasus.com/>) configuration using Terraform.

This provider is still in an early stage and only supports a limited set of resources.

Additional resource will be added over time, but in a demand drive way due to limited development capacity.
If you need a specific resource please open a feature request ticket in the github repository: <https://github.com/pkerspe/terraform-provider-databasus/issues>

Please note that this provider is not developed by the Databasus team and not supported by them, in fact not even encouraged to use an IaC approach for configuration since it is not in line with their goals for the tool. As such, provider support could break with any new release of Databasus.

To limit risk please make sure to match the Provider release version with the Databasus version you are using.

Provider Support Map:

| Provider Version | Tested Databasus Version |
|------------------|--------------------------|
| v0.1.x           | v3.32.2                  |