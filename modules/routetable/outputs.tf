output "route_table_resource_id" {
  description = "The created route table ID."
  value = {
    route_table = azapi_resource.route_table.id
  }
}
