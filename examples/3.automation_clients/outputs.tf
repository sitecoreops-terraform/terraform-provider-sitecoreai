output "deploy_id" {
  description = "The internal id of the SitecoreAI Deploy automation client"
  value       = sitecore_deploy_client.default.id
}

output "deploy_client_id" {
  description = "The client id for the SitecoreAI Deploy automation client"
  value       = sitecore_deploy_client.default.client_id
}

output "deploy_client_secret" {
  description = "The client secret for the SitecoreAI Deploy  automation client"
  value       = sitecore_deploy_client.default.client_secret
  sensitive   = true
}

output "cm_client_id" {
  description = "The client id for the CM automation client"
  value       = sitecore_cm_client.default.client_id
}

output "cm_client_secret" {
  description = "The client secret for the CM automation client"
  value       = sitecore_cm_client.default.client_secret
  sensitive   = true
}

output "edge_client_id" {
  description = "The client id for the edge automation client"
  value       = sitecore_edge_client.default.client_id
}

output "edge_client_secret" {
  description = "The client secret for the edge automation client"
  value       = sitecore_edge_client.default.client_secret
  sensitive   = true
}
