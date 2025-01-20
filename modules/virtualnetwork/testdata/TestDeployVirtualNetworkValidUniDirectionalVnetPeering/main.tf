resource "azapi_resource" "rg" {
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  parent_id = "/subscriptions/${var.subscription_id}"
  name      = "${var.virtual_networks["primary"].name}-hub"
  location  = var.virtual_networks["primary"].location
}

resource "azapi_resource" "hub" {
  type      = "Microsoft.Network/virtualNetworks@2021-08-01"
  name      = "${var.virtual_networks["primary"].name}-hub"
  parent_id = azapi_resource.rg.id
  location  = azapi_resource.rg.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "192.168.10.0/23"
        ]
      }
      subnets = [
        {
          name = "default"
          properties = {
            addressPrefix = "192.168.10.0/24"
          }
        },
        {
          name = "GatewaySubnet"
          properties = {
            addressPrefix = "192.168.11.0/24"
          }
        }
      ]
    }
  }
}

locals {
  virtual_network_primary_merged = merge(var.virtual_networks["primary"], {
    hub_network_resource_id = azapi_resource.hub.id
  })
  virtual_network_secondary_merged = merge(var.virtual_networks["secondary"], {
    hub_network_resource_id = azapi_resource.hub.id
  })
  virtual_networks_merged = {
    primary   = local.virtual_network_primary_merged
    secondary = local.virtual_network_secondary_merged
  }
}

module "virtualnetwork_test" {
  source           = "../../"
  subscription_id  = var.subscription_id
  virtual_networks = local.virtual_networks_merged
}
