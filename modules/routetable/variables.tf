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

variable "resource_group_name" {
  type        = string
  description = <<DESCRIPTION
The name of the resource group to create the virtual network in.
The resource group must exist, this module will not create it.
DESCRIPTION
  nullable    = false
}

variable "bgp_route_propagation_enabled" {
  type        = bool
  description = <<DESCRIPTION
Whether BGP route propagation is enabled.
DESCRIPTION
  default     = true
}

variable "name" {
  type        = string
  description = <<DESCRIPTION
The name of the route table to create.
DESCRIPTION
  nullable    = false
}

variable "location" {
  type        = string
  description = <<DESCRIPTION
The location of the route table.
DESCRIPTION
  nullable    = false
}

variable "tags" {
  type        = map(string)
  description = <<DESCRIPTION
A map of tags to assign to the route table.
DESCRIPTION
  default     = {}
  nullable    = false
}

variable "routes" {
  type = list(object({
    name                   = string
    address_prefix         = string
    next_hop_type          = string
    next_hop_in_ip_address = string
  }))

  description = <<DESCRIPTION
A list of objects defining route tables and their associated routes to be created:

- `name` (required): The name of the route.
- `address_prefix` (required): The address prefix for the route.
- `next_hop_type` (required): The type of next hop for the route.
- `next_hop_in_ip_address` (required): The next hop IP address for the route.
DESCRIPTION
  nullable    = false
  default     = []
}
