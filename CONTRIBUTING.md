# SitecoreAI Terraform provider

The project is mainly a Terraform provider built using
* Go
* Terraform Plugin Framework

There are two packages in the `pkg` folder.
* `pkg/apiclient/` which is a client to interact with the SitecoreAI Deploy API
* `pkg/provider/` which is the Terraform provider that exposes resources and datasources and uses the apiclient to call the api.

Additionally there are some terraform examples in the `examples/` folder to show how the provider can be used in terraform modules.

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