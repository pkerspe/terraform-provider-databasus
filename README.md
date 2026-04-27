# terraform-provider-databasus

A terraform Provider for DATABASUS (<https://databasus.com/>) to manage resources like workspaces, users, DB connections, storages and schedules.

## Note

This provider is still a DRAFT version.
Use at your own risk.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24

## Building the Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## testing the provider locally

This will install the provider locally in your go bin folder (depending on OS).
You can use the provider locally by following this guide: <https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider#prepare-terraform-for-local-provider-install>

Main steps are:

- find out where go stores the compiled binary, excute `go env GOBIN` or if empty use default path (under Windows usually `/Users/<Username>/go/bin`)
- update or create `.terraformrc` file and add the path to the provider binary created by `go install` to a override for the databasus provider:

````json
provider_installation {

  dev_overrides {
      "registry.terraform.io/pkerspe/databasus" = "<PATH to bin folder here>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
````  

- start the databasus docker image using the docker-compose file under /docker_compose
- now you can run a terraform plan or apply e.g. using the prepared `main.tf` file under `/example/provider-install-verification`

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) and for testing terraform installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

You can also just run terraform locally against a local test instance of databasus. You will find the needed docker compose file under `/docker_compose/docker-compose.yml` and you find a test terraform script under `/examples/provider-install-verification`
