output "budget_resource_id" {
  description = "The Azure resource id of the created budget."
  value       = azapi_resource.budget.id
}
