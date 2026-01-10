#!/bin/bash

# Script to generate and update documentation for the SitecoreAI Terraform provider

set -e

echo "Generating documentation for SitecoreAI Terraform provider..."

# Generate documentation using tfplugindocs
go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name sitecore

echo "Documentation generated successfully!"

# Restore custom files that were overwritten
echo "Restoring custom documentation files..."

# Check if we have a backup of the getting started guide
if [ -f "docs/guides/getting-started.md.backup" ]; then
    mv docs/guides/getting-started.md.backup docs/guides/getting-started.md
fi

# Check if we have a backup of the main index.md
if [ -f "docs/index.md.backup" ]; then
    mv docs/index.md.backup docs/index.md
fi

echo "Custom files restored!"

echo "Documentation generation complete. Please review and enhance the generated files."
echo "Run 'go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name sitecore' to regenerate."