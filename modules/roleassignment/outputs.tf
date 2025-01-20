output "role_assignment_resource_id" {
  description = "The Azure resource id of the created role assignment."
  value       = azurerm_role_assignment.this.id
}

output "role_assignment_name" {
  description = "The Azure name (uuid) of the created role assignment."
  value       = azurerm_role_assignment.this.name
}
