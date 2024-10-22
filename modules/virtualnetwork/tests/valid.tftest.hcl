mock_provider "azapi" {}

variables {
  subscription_id = "00000000-0000-0000-0000-000000000000"
  virtual_networks = {
    primary = {
      name                = "primary-vnet",
      address_space       = ["192.168.0.0/24"],
      location            = "westeurope",
      resource_group_name = "primary-rg",
    }
    secondary = {
      name                = "secondary-vnet",
      address_space       = ["192.168.1.0/24"],
      location            = "northeurope",
      resource_group_name = "secondary-rg",
    }
  }
}

run "valid" {
  command = plan

  assert {
    condition     = alltrue([for k, v in var.virtual_networks : azapi_resource.vnet[k].name == v.name])
    error_message = "Virtual network names do not match input"
  }
  assert {
    condition     = alltrue([for k, v in var.virtual_networks : azapi_resource.vnet[k].location == v.location])
    error_message = "Virtual network locations do not match input"
  }
  assert {
    condition     = alltrue([for k, v in var.virtual_networks : azapi_resource.vnet[k].body.properties.addressSpace.addressPrefixes == v.address_space])
    error_message = "Virtual network locations do not match input"
  }
  assert {
    condition     = alltrue([for k, v in var.virtual_networks : azapi_resource.rg["${k}-rg"].name == v.resource_group_name])
    error_message = "Resource group names do not match input"
  }
  assert {
    condition     = alltrue([for k, v in var.virtual_networks : azapi_resource.rg["${k}-rg"].location == v.location])
    error_message = "Resource group locations do not match input"
  }
}
