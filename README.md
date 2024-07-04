# terraform-provider-hyperstack

This project is a Terraform provider for Hyperstack. It uses the Hyperstack SDK to interact with the Hyperstack API and manage resources. The provider is written in Go and uses the Terraform Plugin SDK for compatibility with Terraform.

## Getting Started

Before you start, make sure you have the following tools installed:

- [Go 1.22](https://golang.org/dl/)
- [Task 3.25](https://taskfile.dev/installation/): A task runner for executing project tasks.
- [jq 1.6](https://jqlang.github.io/jq/download/): A command-line JSON processor.
- [yq v4](https://github.com/mikefarah/yq/): A command-line YAML processor.
- GPG and gpg-agent
- Python 3.11
- [canonicaljson](https://pypi.org/project/canonicaljson/) pip module

There are also CLI dependencies that are installed vith Go:

- [GoReleaser](https://goreleaser.com/)
  ````bash
  go install github.com/goreleaser/goreleaser@latest
  ````
- [OpenAPI Provider Spec Generator](https://developer.hashicorp.com/terraform/plugin/code-generation/openapi-generator): generates provider spec using OpenAPI definition
  ````bash
  go install github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi@latest
  ````
- [Framework Code Generator](https://developer.hashicorp.com/terraform/plugin/code-generation/framework-generator): generates Golang schemas using provider spec
  ````bash
  go install github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework@latest
  ````

## Building the Provider

The provider uses a `Taskfile.yaml` for task management. To build the provider, run the following command:

```bash
task build
```

This will compile the provider and output the binary in the `artifacts/provider` directory.

## Generating Schemas

The provider uses the OpenAPI generator to generate schemas. To generate the schemas, run the following command:

```bash
task gen
```

This will pull the latest API specification from the server, generate the schemas, and output them in the `artifacts/provider-spec.json` file.

## Testing the Provider

To test the provider, run the following command:

```bash
task test-examples
```

This will build the provider and run the tests.

## Publishing the Provider

To publish a new version of the provider, run the following command:

```bash
task provider-publish
```

This will build the provider, create a new release, and upload it to the Terraform registry.

## Documentation

For more information about the features of the Hyperstack API, visit the [Hyperstack Documentation](https://infrahub-doc.nexgencloud.com/docs/features/).

For more information about the OpenAPI generator used in this project, visit the [HashiCorp Developer Guide](https://developer.hashicorp.com/terraform/plugin/code-generation/openapi-generator).

## Contributing

Contributions to this project are welcome. Please make sure to test your changes before submitting a pull request.
