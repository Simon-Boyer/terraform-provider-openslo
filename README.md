# OpenSLO Terraform provider

This terraform provider allows you to ingest OpenSLO (https://github.com/OpenSLO/OpenSLO) definitions (as yaml) using a data_source.
It provides the complete OpenSLO specification a terraform/HCL objects that can then be used
to manage your observabilty tooling.

To use it, put all the yaml files in a single string, delimited by `\n---\n` and pass it to the `openslo` data_source.
Ideally everything is in a single data_source as it allows it to resolve links between the definitions.

## How to use

```hcl
terraform {
  required_providers {
    openslo = {
      source = "registry.terraform.io/arctiq/openslo"
    }
  }
}

data "openslo_openslo" "definition" {
  yaml_input = <<EOF
  <YOUR YAML DEFINITIONS>
EOF
}

...

something = openslo_openslo.definition.object_kind["object_name"].object_property
```

## Contributing

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

### Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

### Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

### Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.