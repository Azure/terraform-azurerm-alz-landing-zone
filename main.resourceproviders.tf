module "resourceproviders" {
  source            = "./modules/resourceprovider"
  for_each          = { for k, v in var.subscription_register_resource_providers_and_features : k => v if var.subscription_register_resource_providers_enabled }
  resource_provider = each.key
  features          = each.value
  subscription_id   = local.subscription_id

  depends_on = [
    module.resourcegroup,
    module.roleassignment,
    module.roleassignment_umi,
    module.subscription,
    module.usermanagedidentity,
    module.virtualnetwork,
  ]
}
