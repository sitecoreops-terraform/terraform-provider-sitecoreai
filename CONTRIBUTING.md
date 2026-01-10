# SitecoreAI Terraform provider

The project is a Terraform provider built using
* Go
* Terraform Plugin Framework

The folder structure:
* `pkg/apiclient/` which is a client to interact with the SitecoreAI Deploy API, see [SitecoreAI Deploy API v1: OpenAPI Specification](https://xmclouddeploy-api.sitecorecloud.io/swagger/v1/swagger.json) and [SitecoreAI Deploy API v2: OpenAPI Specification](https://xmclouddeploy-api.sitecorecloud.io/swagger/v2/swagger.json)
* `pkg/provider/` which is the Terraform provider that exposes resources and datasources and uses the apiclient to call the api.
* `examples/` with several terraform examples to show how the provider can be used in terraform modules.
* `docs/` contains the provider documentation in the format required by the Terraform Registry.

## Linting

```bash
golangci-lint run
```

## Testing api client

To run the test, those are integration tests and will make actual calls to the SitecoreAI Deploy API:

```bash
# Set your Sitecore credentials as environment variables
export SITECORE_CLIENT_ID=your_client_id
export SITECORE_CLIENT_SECRET=your_client_secret

# Run a specific test, here client authentication
go test ./pkg/apiclient/... -v -run TestClientAuthentication

# Run all tests
go test ./pkg/apiclient/... -v
```

## Testing terraform provider

Build the provider:

```sh
# Build the provider
go build -o out/terraform-provider-sitecoreai

# Run a specific test, here client authentication
go test ./pkg/provider/... -v -run TestProviderMetadata

# Run all tests
go test ./pkg/provider/... -v
```

## Documentation

The provider documentation is located in the `docs/` directory and follows the Terraform Registry documentation format.

### Generating Documentation

The documentation can be automatically generated and updated using the `tfplugindocs` tool:

```bash
# Generate documentation from the provider schema
go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name sitecore
```

This will update the documentation files in the `docs/` directory based on the provider's Go code.

### Updating Documentation

After generating documentation, you should:

1. Review the generated files
2. Add examples and additional explanations where needed
3. Ensure all documentation follows the Terraform Registry format
4. Update the `guides/` directory with any new guides or tutorials

### Documentation Structure

```
docs/
├── index.md                  # Main provider documentation
├── guides/                   # Additional guides and tutorials
│   └── getting-started.md    # Getting started guide
├── resources/                # Resource documentation
│   ├── project.md
│   ├── environment.md
│   └── ...
└── data-sources/             # Data source documentation
    ├── project.md
    ├── environment.md
    └── ...
```

### Documentation Requirements

* All documentation must be in Markdown format
* Follow the Terraform Registry documentation guidelines
* Include examples for all resources and data sources
* Mark sensitive attributes appropriately
* Use proper YAML frontmatter for guides

## Diagnostics

The communication can be captured by specifying a proxy server eg. Burp

```bash
# When running in devcontainers and proxy server on host
export HTTPS_PROXY="http://host.docker.internal:8081"

# When running both project and proxy server on host
export HTTPS_PROXY="http://localhost:8081"
```

It is necessary to trust the certificate, eg. Burp can export the certificate in DER format.

```bash
# Convert DER public certificate to .pem format
openssl x509 -inform der -in burp.der -out burp.pem

# Set REQUESTS_CA_BUNDLE to trust the certificate(s)
export REQUESTS_CA_BUNDLE="$(pwd)/burp.pem"
```
