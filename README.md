<a href="https://terraform.io">
    <img src="https://raw.githubusercontent.com/hashicorp/terraform-website/master/public/img/logo-text.svg" alt="Terraform logo" title="Terraform" height="50" />
</a>

# Terraform Provider for OpenSearch Dashboards

The Terraform OpenSearch provider is a plugin for Terraform that allows for the full lifecycle management of OpenSearch resources. Manage Kibana saved objects, including dashboards, visualizations, and more.
This provider is maintained internally by the [MOIA GmbH](https://moia.io) team.

## Examples

All the resources and data sources have [one or more examples](./smoketest) to give you an idea of how to use this provider to build your own OpenSearch SavedObjects infrastructure.

# Development Environment Setup

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.15.0 or newer.
- [Go](https://golang.org/doc/install) 1.24 (to build the provider plugin)

## Opensearch version

Currently this provider is only tested with Opensearch version 1.3.6
To update the smoke-test to a new Opensearch version, changes need to be made
* in the Makefile (env-var VERSION)
* in the github-actions (tag of the opensearch-docker-image)

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

## Testing with Unit Tests
In order to run unit tests, you can run `make test`.

```sh
$ make test
```

## Smoketest
This provider contains a smoketest which can test that the apply works for a few default examples
against a local opensearch-instance.

### Prereqesites
To run the smoketest you need a local opensearch. If you do not already have one running on your
machine, follow these steps to start it:

1. install podman or docker (and terraform if it's not already installed ;) )
2. make sure that ports 9200, 9600 and 5601 are currently not in use 
3. run `make start_opensearch` (when using docker instead of podman: `make start opensearch CONTAINER_RUNTIME=docker`)

If you get an error like `max virtual memory areas vm.max_map_count [65530] likely too low, increase to at least [262144]` 
in your container logs, running `sysctl -w vm.max_map_count=262144` can help (this resets after reboot so add to .bashrc or other startup-file if needed)

### Running the smoketest

`make smoke_test`

While actively developing this plugin if you need to run the smoke_test often you can also use
`make smoke_test_fast` but this is not as stable, so if you run into errors, fallback to `make smoke_test`.

### Cleanup

When you finished testing and want to remove your local opensearch again, execute `make remove_opensearch`
(or when using docker instead of podman: `make remove_opensearch CONTAINER_RUNTIME=docker`)

## Using the Provider

To use a released provider in your Terraform environment,run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the
provider. To specify a particular provider version when installing released providers, see the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions) .

To instead use a custom-built provider in your Terraform environment (e.g. the provider binary from the build instructions above), follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-plugins). After placing the custom-built provider into your plugins' directory, run `terraform init` to initialize it.

## Releasing a new provider version

The release of the provider is done automatically with a [GitHub action](.github/workflows/release.yaml). The workflow gets triggered if a new tag is created.
Tags can be created with the GitHub web UI or the command line.

The release is signed with a GPG key that is stored as a repository secret and in [terraform](https://registry.terraform.io/providers/moia-oss/opensearch-dashboards/). 

## Contributing

We really appreciate your help!

To contribute, simply make a PR and a maintainer will review it shortly.

Issues on GitHub are intended to be related to the bugs or feature requests with provider codebase.
See [Plugin SDK Community](https://www.terraform.io/community)
and [Discuss forum](https://discuss.hashicorp.com/c/terraform-providers/31/none) for a list of community resources to
ask questions about Terraform.
