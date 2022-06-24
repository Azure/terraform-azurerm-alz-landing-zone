variable "virtual_network_name" {
  type        = string
  description = <<DESCRIPTION
    The name of the virtual network.
  DESCRIPTION
  validation {
    condition     = can(regex("^[\\w-_.]{2,64}$", var.virtual_network_name))
    error_message = "The string must consist of a-z, A-Z, 0-9, -, _, and . (period) and be between 2 and 64 characters in length."
  }
}

variable "virtual_network_address_space" {
  type        = list(string)
  description = <<DESCRIPTION
    The address space of the virtual network, supplied as multiple CIDR blocks, e.g. `["10.0.0.0/16","172.16.0.0/12"]`.
  DESCRIPTION
}

variable "hub_network_resource_id" {
  type        = string
  description = <<DESCRIPTION
    The resource ID of the virtual network in the hub to which the created virtual network will be peered.

    E.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet`

    Leave blank to create the virtual network without peering.
  DESCRIPTION
  default     = ""
  validation {
    condition     = can(regex("^$|^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w-._]{1,89}[^\\s.]/providers/Microsoft.Network/virtualNetworks/[\\w-_.]{2,64}$", var.hub_network_resource_id))
    error_message = "Value must be an Azure virtual network resource id, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet."
  }
}

variable "vwan_hub_resource_id" {
  type        = string
  description = <<DESCRIPTION
    The resource ID of the vwan hub to which the virtual network will be connected.

    E.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-hub`

    Leave blank to create a virtual network without a vwan hub connection.
  DESCRIPTION
  default     = ""
  validation {
    condition     = can(regex("^$|^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w-._]{1,89}[^\\s.]/providers/Microsoft.Network/virtualHubs/[\\w-_.]{1,80}$", var.vwan_hub_resource_id))
    error_message = "Value must be an Azure vwan hub resource id, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-hub."
  }
}

variable "subscription_id" {
  type        = string
  description = <<DESCRIPTION
    The subscription ID of the subscription to create the virtual network in.
  DESCRIPTION
  validation {
    condition     = can(regex("^[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", var.subscription_id))
    error_message = "Must a GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx. All letters must be lowercase."
  }
}

variable "virtual_network_resource_group_name" {
  type        = string
  description = <<DESCRIPTION
    The name of the resource group to create the virtual network in.
  DESCRIPTION
  validation {
    condition     = can(regex("^[\\w-_.]{1,89}[^\\s.]$", var.virtual_network_resource_group_name))
    error_message = "Value must be between 1 and 90 characters in length and start with a letter or number, and end with a letter or number."
  }
}

variable "virtual_network_location" {
  type        = string
  description = <<DESCRIPTION
    The location of the virtual network.
  DESCRIPTION
}

variable "virtual_network_use_remote_gateways" {
  type        = bool
  description = <<DESCRIPTION
    Enables the use of remote gateways for the virtual network.

    Applies to both hub and spoke (vnet peerings) as well as virtual WAN connections.
  DESCRIPTION
  default     = true
}

variable "virtual_network_vwan_routetable_resource_id" {
  type = string
  description = <<DESCRIPTION
    The resource ID of the virtual network route table to use for the virtual network.

    Leave blank to use the `defaultRouteTable`.

    E.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-vhub/hubRouteTables/defaultRouteTable
  DESCRIPTION
  default = ""
  validation {
    condition     = can(regex("^$|^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w]{1,89}[^\\s.]/providers/Microsoft.Network/virtualHubs/[\\w-_.]{1,80}/hubRouteTables/[\\w-_.]{1,80}$", var.virtual_network_vwan_routetable_resource_id))
    error_message = "Value must be an Azure vwan hub resource id, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-vhub/hubRouteTables/defaultRouteTable."
  }
}

variable "virtual_network_vwan_propagated_routetables_resource_ids" {
  type = list(string)
  description = <<DESCRIPTION
    The list of route table resource ids to advertise routes to.

    Leave blank to use the `defaultRouteTable.
  DESCRIPTION
  default = []
}

variable "virtual_network_vwan_propagated_routetables_labels" {
  type = list(string)
  description = <<DESCRIPTION
    The list of virtual WAN labels to advertise the routes to.

    Leave blank to use the `default` label.
  DESCRIPTION
  default = []
}

# variable "virtual_network_subnets" {
#   type = map(object({
#     address_prefix = string
#     }))
#   description = <<DESCRIPTION
#     The subnets of the virtual network, supplied as multiple objects.

#     e.g.

#     ```terraform
#     virtual_network_subnets = {
#       subnet0 = {
#         address_prefix = "10.0.0.0/24"
#       },
#       subnet1 = {
#         address_prefix = "10.0.1.0/24"
#     } }
#     ```
#   DESCRIPTION
# }
