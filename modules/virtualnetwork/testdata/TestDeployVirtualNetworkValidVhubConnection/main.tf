resource "azapi_resource" "rg" {
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  parent_id = "/subscriptions/${var.subscription_id}"
  name      = "${var.virtual_networks["primary"].name}-hub"
  location  = var.virtual_networks["primary"].location
}

resource "azapi_resource" "vwan" {
  type      = "Microsoft.Network/virtualWans@2021-08-01"
  name      = "${var.virtual_networks["primary"].name}-vwan"
  location  = azapi_resource.rg.location
  parent_id = azapi_resource.rg.id
  body = {
    properties = {
      type                       = "Standard"
      allowBranchToBranchTraffic = true
      disableVpnEncryption       = false
    }
  }
}

resource "azapi_resource" "vhub" {
  type      = "Microsoft.Network/virtualHubs@2021-08-01"
  name      = "${var.virtual_networks["primary"].name}-vhub"
  location  = azapi_resource.vwan.location
  parent_id = azapi_resource.rg.id
  body = {
    properties = {
      addressPrefix = "192.168.100.0/23"
      sku           = "Standard"
      virtualWan = {
        id = azapi_resource.vwan.id
      }
    }
  }
}

locals {
  virtual_network_primary_merged = merge(var.virtual_networks["primary"], {
    vwan_hub_resource_id = azapi_resource.vhub.id
  })
  virtual_network_secondary_merged = merge(var.virtual_networks["secondary"], {
    vwan_hub_resource_id = azapi_resource.vhub.id
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
