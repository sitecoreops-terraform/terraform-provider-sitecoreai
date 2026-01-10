# SitecoreAI Provider Documentation

This directory contains the documentation for the SitecoreAI Terraform provider in the format required by the Terraform Registry.

## Documentation Structure

```
docs/
├── README.md                    # This file
├── index.md                     # Main provider documentation
├── guides/                      # Additional guides and tutorials
│   └── getting-started.md       # Getting started guide
├── resources/                   # Resource documentation
│   ├── project.md
│   ├── environment.md
│   ├── cm_environment.md
│   ├── cm_client.md
│   ├── edge_client.md
│   ├── deploy_client.md
│   └── editing_host_build_client.md
└── data-sources/                # Data source documentation
    ├── project.md
    ├── environment.md
    └── editing_secret.md
```

## Documentation Generation

The documentation is automatically generated from the provider's Go code using the `tfplugindocs` tool. To regenerate the documentation:

```bash
go generate ./pkg/provider/
```

This will update all the documentation files based on the current provider schema.

## Manual Updates

After generating documentation, you should manually enhance the generated files by:

1. Adding practical examples
2. Including usage notes and best practices
3. Adding import instructions where applicable
4. Including warnings and important notes

## Documentation Format

All documentation follows the Terraform Registry format:

### Provider Documentation (index.md)
- Overview of the provider
- Authentication instructions
- Example usage
- Argument reference

### Resource/Data Source Documentation
- Description of what the resource/data source manages/retrieves
- Example usage with HCL code
- Argument reference (required/optional attributes)
- Attribute reference (exported attributes)
- Import instructions (for resources)
- Notes and best practices

### Guides
- Step-by-step tutorials
- Common patterns and use cases
- Troubleshooting guides

## Adding New Documentation

When adding new resources or data sources:

1. Implement the resource/data source in Go code
2. Run `go generate ./pkg/provider/` to generate the initial documentation
3. Enhance the generated documentation with examples and explanations
4. Add any necessary guides or tutorials

## Publishing to Terraform Registry

The documentation in this directory will be automatically published to the Terraform Registry when a new provider version is released. The registry will display the documentation from the Git commit associated with each version.

## Best Practices

- Keep examples practical and realistic
- Use consistent formatting and style
- Mark sensitive data appropriately
- Include warnings for destructive operations
- Reference official SitecoreAI documentation where applicable