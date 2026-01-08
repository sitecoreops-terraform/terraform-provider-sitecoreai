variable "sitecore_client_id" {
  description = "The client ID for authentication"
  type        = string
  sensitive   = true
}

variable "sitecore_client_secret" {
  description = "The client secret for authentication"
  type        = string
  sensitive   = true
}

variable "project_name" {
  description = "The name of the SitecoreAI project"
  type        = string
  default     = "XMC"
}
