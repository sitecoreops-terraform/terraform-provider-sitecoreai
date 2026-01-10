variable "sitecore_client_id" {
  description = "The client ID for authentication"
  type        = string
  sensitive   = true
  default     = ""
}

variable "sitecore_client_secret" {
  description = "The client secret for authentication"
  type        = string
  sensitive   = true
  default     = ""
}

variable "project_name" {
  description = "The name of the SitecoreAI project"
  type        = string
  default     = "XMC"
}

variable "environment_name" {
  description = "The name of the SitecoreAI environment"
  type        = string
  default     = "dev2"
}

variable "environment_is_prod" {
  description = "Environment is using production SLA"
  type        = bool
  default     = false
}
