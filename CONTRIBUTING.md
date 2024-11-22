# Contributing to Hyperstack Terraform Provider

We love your input! We want to make contributing to this project as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features

## We Develop with Github

We use github to host code, to track issues and feature requests, as well as accept pull requests.

## We Use [Github Flow](https://guides.github.com/introduction/flow/index.html), So All Code Changes Happen Through Pull Requests

Pull requests are the best way to propose changes to the codebase (we use [Github Flow](https://guides.github.com/introduction/flow/index.html)). We actively welcome your pull requests:

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

## Report bugs using Github's [issues](https://github.com/NexGenCloud/terraform-provider-hyperstack/issues)

We use GitHub issues to track public bugs. Report a bug by [opening a new issue](); it's that easy!

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can, includes sample code that _anyone_ can run to reproduce.
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

## Developing the Hyperstack Terraform Provider

Before you start, make sure you have the following tools installed:

- [Terraform 1.5](https://developer.hashicorp.com/terraform/install) - [tfenv](https://github.com/tfutils/tfenv) is suggested
- [Terragrunt](https://terragrunt.gruntwork.io/) - [tgswitch](https://tgswitch.warrensbox.com/) is suggested
- [Go 1.23](https://golang.org/dl/)
- [Task 3.25](https://taskfile.dev/installation/): A task runner for executing project tasks.
- [jq 1.6](https://jqlang.github.io/jq/download/): A command-line JSON processor.
- [yq v4.44](https://github.com/mikefarah/yq/): A command-line YAML processor.
- GPG and gpg-agent
- Python 3.11
- [canonicaljson](https://pypi.org/project/canonicaljson/) pip module
  - You can run installation via pip:
  ```bash
  python3 -m pip install -r requirements.txt
  ```

There are also CLI dependencies that are installed with Go:

- [GoReleaser](https://goreleaser.com/)
  ```bash
  go install github.com/goreleaser/goreleaser/v2@latest
  ```
- [OpenAPI Provider Spec Generator](https://developer.hashicorp.com/terraform/plugin/code-generation/openapi-generator): generates provider spec using OpenAPI definition
  ```bash
  go install github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi@latest
  ```
- [Framework Code Generator](https://developer.hashicorp.com/terraform/plugin/code-generation/framework-generator): generates Golang schemas using provider spec
  ```bash
  go install github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework@latest
  ```

### Building the Provider

The provider uses a `Taskfile.yaml` for task management. To build the provider, run the following command:

```bash
task build

# to run a test individually, see example below
# task TEST=core_vm test-run
```

This will compile the provider and output the binary in the `artifacts/provider` directory.

### Generating Schemas

The provider uses the OpenAPI generator to generate schemas. To generate the schemas, run the following command:

```bash
task gen
```

This will pull the latest API specification from the server, generate the schemas, and output them in the `artifacts/provider-spec.json` file.

### Testing the Provider

To test the provider, run the following command.

Make sure to set the following environment variables:

- `HYPERSTACK_API_KEY`: Your Hyperstack API key
- `HYPERSTACK_API_REGION`: Your Hyperstack region (e.g. CANADA-1)

```bash
task test
```

This will build the provider and run the tests.

Please note:

1. The tests will create and destroy many resources in your Hyperstack account.
2. The Kubernetes API is still in beta so it might be unstable. Therefore, you might need to re-run the tests if they fail during the Kubernetes tests.

### Documentation

For more information about the features of the Hyperstack API, visit the [Hyperstack Documentation](https://infrahub-doc.nexgencloud.com/docs/features/).

For more information about the OpenAPI generator used in this project, visit the [HashiCorp Developer Guide](https://developer.hashicorp.com/terraform/plugin/code-generation/openapi-generator).
