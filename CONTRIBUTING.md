# SitecoreAI Terraform provider

The project is a Terraform provider built using
* Go
* Terraform Plugin Framework

The folder structure:
* `pkg/apiclient/` which is a client to interact with the SitecoreAI Deploy API, see [SitecoreAI Deploy API: OpenAPI Specification](https://xmclouddeploy-api.sitecorecloud.io/swagger/v1/swagger.json)
* `pkg/provider/` which is the Terraform provider that exposes resources and datasources and uses the apiclient to call the api.
* `examples/` with several terraform examples to show how the provider can be used in terraform modules.

## Testing api client

To run the authentication test:

```bash
# Set your Sitecore credentials as environment variables
export SITECORE_CLIENT_ID=your_client_id
export SITECORE_CLIENT_SECRET=your_client_secret

# Run a specific test
go test ./pkg/apiclient/... -v -run TestClientAuthentication

# Run all tests
go test ./pkg/apiclient/... -v
```

## Testing terraform provider

Build the provider:

```sh
go build -o out/terraform-provider-sitecoreai
```

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
