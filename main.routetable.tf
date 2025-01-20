# route table submodule, disabled by default
# Will create a route table, and optionally routes
module "routetable" {
  for_each        = var.route_table_enabled ? var.route_tables : {}
  source          = "./modules/routetable"
  subscription_id = local.subscription_id

  resource_group_name           = each.value.resource_group_name
  bgp_route_propagation_enabled = each.value.bgp_route_propagation_enabled
  name                          = each.value.name
  location                      = each.value.location
  routes                        = each.value.routes
  tags                          = each.value.tags

  depends_on = [
    module.resourcegroup,
  ]
}
