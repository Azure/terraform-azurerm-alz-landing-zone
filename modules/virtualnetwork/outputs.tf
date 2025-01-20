output "virtual_network_resource_ids" {
  description = "The created virtual network resource IDs, expressed as a map."
  value = {
    for k, v in azapi_resource.vnet : k => v.id
  }
}

output "resource_group_resource_ids" {
  description = "The created resource group IDs, expressed as a map."
  value = {
    for k, v in azapi_resource.rg : k => v.id
  }
}
