resource "sitecoreai_environment_variable" "env" {
  environment_id = "environment-1234"

  name  = "SXA_ENVIRONMENT_NAME"
  value = "dev"
}
