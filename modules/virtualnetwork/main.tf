# azapi_resource.rg is the resource group that the virtual network will be created in
# the module will create as many as is required by the var.virtual_networks input variable
resource "azapi_resource" "rg" {
  for_each  = { for i in local.resource_group_data : i.name => i }
  parent_id = local.subscription_resource_id
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  name      = each.key
  location  = each.value.location
  tags      = each.value.tags
}

# azapi_resource.rg_lock is an optional resource group lock that can be used
# to prevent accidental deletion.
resource "azapi_resource" "rg_lock" {
  for_each  = { for i in local.resource_group_data : i.name => i if i.lock }
  type      = "Microsoft.Authorization/locks@2017-04-01"
  parent_id = azapi_resource.rg[each.key].id
  name      = coalesce(each.value.lock_name, substr("lock-${each.key}", 0, 90))
  body = {
    properties = {
      level = "CanNotDelete"
    }
  }
  depends_on = [
    azapi_resource.vnet,
    azapi_resource.peering_hub_outbound,
    azapi_resource.peering_hub_inbound,
    azapi_resource.peering_mesh,
    azapi_resource.vhubconnection,
    azapi_resource.vhubconnection_routing_intent,
  ]
}

# azapi_resource.vnet are the virtual networks that will be created
# lifecycle ignore changes to the body to prevent subnets being deleted
# see #45 for more information
resource "azapi_resource" "vnet" {
  for_each  = var.virtual_networks
  parent_id = "${local.subscription_resource_id}/resourceGroups/${each.value.resource_group_name}"
  type      = "Microsoft.Network/virtualNetworks@2023-09-01"
  name      = each.value.name
  location  = coalesce(each.value.location, var.location)
  body = {
    properties = local.vnet_body_properties[each.key]
  }
  tags = each.value.tags
  depends_on = [
    azapi_resource.rg,
  ]
}

# azapi_resource.peering_hub_outbound creates one-way peering from the spoke to the supplied hub virtual network.
# They are not created if the hub virtual network is an empty string.
resource "azapi_resource" "peering_hub_outbound" {
  for_each  = { for k, v in local.hub_peering_map : k => v if v.peering_direction != local.peering_direction_fromhub }
  type      = "Microsoft.Network/virtualNetworks/virtualNetworkPeerings@2021-08-01"
  parent_id = each.value["outbound"].this_resource_id
  name      = each.value["outbound"].name
  body = {
    properties = {
      remoteVirtualNetwork = {
        id = each.value["outbound"].remote_resource_id
      }
      allowVirtualNetworkAccess = true
      allowForwardedTraffic     = true
      allowGatewayTransit       = false
      useRemoteGateways         = each.value.use_remote_gateways
    }
  }
}

# azapi_resource.peering_hub_inbound creates one-way peering from the supplied hub network to the spoke.
# They are not created if the hub virtual network is an empty string.
resource "azapi_resource" "peering_hub_inbound" {
  for_each  = { for k, v in local.hub_peering_map : k => v if v.peering_direction != local.peering_direction_tohub }
  type      = "Microsoft.Network/virtualNetworks/virtualNetworkPeerings@2021-08-01"
  parent_id = each.value["inbound"].this_resource_id
  name      = each.value["inbound"].name
  body = {
    properties = {
      remoteVirtualNetwork = {
        id = each.value["inbound"].remote_resource_id
      }
      allowVirtualNetworkAccess = true
      allowForwardedTraffic     = true
      allowGatewayTransit       = true
      useRemoteGateways         = false
    }
  }
}

# azapi_resource.peering_mesh creates mesh peerings between the supplied var.virtual_networks.
# They are created if the boolean mesh_peering_enabled is set to true on more than one network.
resource "azapi_resource" "peering_mesh" {
  for_each  = { for i in local.virtual_networks_mesh_peering_list : "${i.source_key}-${i.destination_key}" => i }
  type      = "Microsoft.Network/virtualNetworks/virtualNetworkPeerings@2022-05-01"
  parent_id = each.value.this_resource_id
  name      = each.value.name
  body = {
    properties = {
      remoteVirtualNetwork = {
        id = each.value.remote_resource_id
      }
      allowVirtualNetworkAccess = true
      allowForwardedTraffic     = each.value.allow_forwarded_traffic
      allowGatewayTransit       = false
      useRemoteGateways         = false
    }
  }
}

# azapi_resource.vhubconnection creates a virtual wan hub connection between the spoke and the supplied vwan hub.
resource "azapi_resource" "vhubconnection" {
  for_each  = { for k, v in var.virtual_networks : k => v if v.vwan_connection_enabled && !v.vwan_security_configuration.routing_intent_enabled }
  type      = "Microsoft.Network/virtualHubs/hubVirtualNetworkConnections@2022-07-01"
  parent_id = each.value.vwan_hub_resource_id
  name      = coalesce(each.value.vwan_connection_name, "vhc-${uuidv5("url", azapi_resource.vnet[each.key].id)}")
  body = {
    properties = local.vhubconnection_body_properties[each.key]
  }
}

# azapi_resource.vhubconnection creates a virtual wan hub connection between the spoke and the supplied vwan hub.
# This resource is used when routing intent is enabled on the vwan security configuration,
# as the routing configuration is then ignored.
resource "azapi_resource" "vhubconnection_routing_intent" {
  for_each  = { for k, v in var.virtual_networks : k => v if v.vwan_connection_enabled && v.vwan_security_configuration.routing_intent_enabled }
  type      = "Microsoft.Network/virtualHubs/hubVirtualNetworkConnections@2022-07-01"
  parent_id = each.value.vwan_hub_resource_id
  name      = coalesce(each.value.vwan_connection_name, "vhc-${uuidv5("url", azapi_resource.vnet[each.key].id)}")
  body = {
    properties = local.vhubconnection_body_properties[each.key]
  }

  lifecycle {
    ignore_changes = [
      body.properties.routingConfiguration,
    ]
  }
}
