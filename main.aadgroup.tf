module "aadgroup" {
  source = "./modules/aadgroup"
  count  = var.aadgroup_enabled ? 1 : 0
  depends_on = [
    module.resourcegroup_networkwatcherrg,
    module.resourcegroup,
    module.subscription,
    module.usermanagedidentity,
    module.virtualnetwork,
  ]
  aad_groups      = var.aad_groups
  subscription_id = local.subscription_id
}
