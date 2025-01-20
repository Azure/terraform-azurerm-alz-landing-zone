# Resource groups module
module "resourcegroup" {
  source              = "./modules/resourcegroup"
  for_each            = var.resource_group_creation_enabled ? var.resource_groups : {}
  subscription_id     = local.subscription_id
  location            = each.value.location
  resource_group_name = each.value.name
  tags                = each.value.tags
}

# v3.3.0 introduced networkwatcherrg support,
# this was then moved into a more general resourcegroups module in later versions
moved {
  from = module.networkwatcherrg[0].azapi_resource.network_watcher_rg
  to   = module.resourcegroup_networkwatcherrg[0].azapi_resource.rg
}

# v5.0.0 consolidated the resource groups into a single module call
moved {
  from = module.resourcegroup_networkwatcherrg[0].azapi_resource.rg
  to   = module.resourcegroup["NetworkWatcherRG"].azapi_resource.rg
}
