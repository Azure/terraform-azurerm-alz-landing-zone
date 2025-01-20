resource "azapi_resource" "subnet" {
  parent_id = module.virtualnetwork_test.virtual_network_resource_ids["primary"]
  type      = "Microsoft.Network/virtualNetworks/subnets@2021-08-01"
  name      = "vnet-subnet"
  body = {
    properties = {
      addressPrefix = "192.168.0.0/26"
    }
  }
}

module "virtualnetwork_test" {
  source           = "../../"
  subscription_id  = var.subscription_id
  virtual_networks = var.virtual_networks
}
