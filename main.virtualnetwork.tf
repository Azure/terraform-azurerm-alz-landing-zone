# Virtual networking submodule, disabled by default
# Will create a vnet, and optionally peerings and a virtual hub connection
module "virtualnetwork" {
  source           = "./modules/virtualnetwork"
  count            = var.virtual_network_enabled ? 1 : 0
  subscription_id  = local.subscription_id
  virtual_networks = var.virtual_networks
  location         = var.location

  depends_on = [
    module.resourcegroup
  ]
}
