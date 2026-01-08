# SitecoreAI Terraform provider

## Testing api client

To run the authentication test:

```bash
# Set your Sitecore credentials as environment variables
export SITECORE_CLIENT_ID=your_client_id
export SITECORE_CLIENT_SECRET=your_client_secret

# Run the test
cd api && go test -v -run TestClientAuthentication
```

## Test terraform provider