# Terraform provider for SitecoreAI

- CONTRIBUTING.md

The OpenAPI specification is already downloaded to `/pkg/api-source/sitecore-api-deploy-v1-swagger.json` and `/pkg/api-source/sitecore-api-deploy-v2-swagger.json`.

## Commands

### Linting
```bash
golangci-lint run
```

### Formatting
```bash
gofmt -s -w -e .
```

### Building
```bash
go build
```

### Testing
#### API Client
```bash
export SITECOREAI_CLIENT_ID=your_client_id
export SITECOREAI_CLIENT_SECRET=your_client_secret
go test ./pkg/apiclient/... -v
```

#### Provider
```bash
go test ./pkg/provider/... -v
```

### Documentation Generation
```bash
cd tools && go generate
```

### Local Development Setup
```bash
go install -v ./...
```

### Debugging
```bash
export TF_LOG=DEBUG
```

### Proxy Setup for Diagnostics
```bash
export SITECOREAI_PROXY="http://host.docker.internal:8081"
```

### Authentication
#### Environment Variables
```bash
export SITECOREAI_CLIENT_ID=your_client_id
export SITECOREAI_CLIENT_SECRET=your_client_secret
```

#### Sitecore CLI
```bash
export SITECOREAI_USE_CLI=1
```

