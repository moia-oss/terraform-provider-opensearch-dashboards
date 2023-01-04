<a href="https://terraform.io">
    <img src="https://raw.githubusercontent.com/hashicorp/terraform-website/master/public/img/logo-text.svg" alt="Terraform logo" title="Terraform" height="50" />
</a>

# Terraform Provider for OpenSearch Dashboards

The Terraform OpenSearch provider is a plugin for Terraform that allows for the full lifecycle management of OpenSearch resources. Manage Kibana saved objects, including dashboards, visualizations, and more.
This provider is maintained internally by the [MOIA GmbH](https://moia.io) team.

## Examples

All the resources and data sources has [one or more examples](./examples) to give you an idea of how to use this provider to build your own OpenSearch SavedObjects infrastructure.

# Development Environment Setup

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.15.0 or newer.
- [Go](https://golang.org/doc/install) 1.19 (to build the provider plugin)

## Quick Start

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (please check the [requirements](#requirements) before proceeding).

_Note:_ This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside
your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your
home directory outside the standard GOPATH (i.e `$HOME/development/terraform-providers/`).

```sh
$ git clone git@github.com:moia-oss/terraform-provider-opensearch-dashboards.git
$ cd terraform-provider-opensearch-dashboards
...
```

`make install` will install the needed tools for the provider.

```sh
$ make install
```

To compile the provider, run `make release`. This will build the provider and put the provider binary in the `./bin`
directory. The provider will be compiled for different supported architectures and operation systems (check [release config](./.goreleaser.yml) for supported OSs and archs). Please check your system before executing the provider locally.

```sh
$ make release
...
$ ./bin/terraform-provider-opensearch-dashboards_<version>_<os>_<arch>
...
```

## Testing the Provider

In order to run unit tests, you can run `make test`.

```sh
$ make test
```

## Using the Provider

To use a released provider in your Terraform environment,run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the
provider. To specify a particular provider version when installing released providers, see the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions) .

To instead use a custom-built provider in your Terraform environment (e.g. the provider binary from the build instructions above), follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-plugins). After placing the custom-built provider into your plugins' directory, run `terraform init` to initialize it.

## Contributing

We really appreciate your help!

To contribute, simply make a PR and a maintainer will review it shortly.

Issues on GitHub are intended to be related to the bugs or feature requests with provider codebase.
See [Plugin SDK Community](https://www.terraform.io/community)
and [Discuss forum](https://discuss.hashicorp.com/c/terraform-providers/31/none) for a list of community resources to
ask questions about Terraform.
