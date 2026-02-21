# Example: Environment variable with no target (applies to all targets)
resource "sitecore_environment" "example" {
  name       = "example-environment"
  project_id = "example-project-id"
  # ... other required fields ...
}

resource "sitecore_environment_variable" "legacy" {
  environment_id = sitecore_environment.example.id
  name           = "SHARED_VAR"
  value          = "hello_world" # Non-sensitive value
  # No target specified (applies to all targets)
}

# Example: CM environment variable
resource "sitecore_environment_variable" "cm_non_secret" {
  environment_id = sitecore_environment.example.id
  name           = "Sitecore_GraphQL_ExposePlayground"
  value          = "true" # Non-sensitive value
  target         = "CM"   # CM target
}

resource "sitecore_environment_variable" "cm_secret" {
  environment_id = sitecore_environment.example.id
  name           = "CM_SECRET_VAR"
  secret_value   = "s3cr3t_p@ssw0rd" # Sensitive value (masked in logs/state)
  target         = "CM"              # CM target
}

# Example: EH environment variable
resource "sitecore_environment_variable" "eh_non_secret" {
  environment_id = sitecore_environment.example.id
  name           = "EH_NON_SECRET_VAR"
  value          = "hello_world" # Non-sensitive value
  target         = "website"     # Target editing host named website in xmcloud.deploy.json
}
