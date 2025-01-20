variable "route_table_enabled" {
  type        = bool
  description = "Whether to create route tables and routes in the target subscription. Requires `var.route_tables`."
  default     = false
}

variable "route_tables" {
  type = map(object({
    name                          = string
    location                      = string
    resource_group_name           = string
    bgp_route_propagation_enabled = optional(bool, true)
    tags                          = optional(map(string))

    routes = optional(map(object({
      name                   = string
      address_prefix         = string
      next_hop_type          = string
      next_hop_in_ip_address = string
    })))
  }))
  description = <<DESCRIPTION
A map defining route tables and their associated routes to be created:

- `name` (required): The name of the route table.
- `location` (required): The location of the resource group.
- `resource_group_name` (required): The name of the resource group.
- `bgp_route_propagation_enabled` (optional): Boolean that controls whether routes learned by BGP are propagated to the route table. Default is `true`.
- `tags` (optional): A map of key-value pairs for tags associated with the route table.
- `routes` (optional): A map defining routes for the route table. Each route object has the following properties:
  - `name` (required): The name of the route.
  - `address_prefix` (required): The address prefix for the route.
  - `next_hop_type` (required): The type of next hop for the route.
  - `next_hop_in_ip_address` (required): The next hop IP address for the route.
DESCRIPTION
  nullable    = false
  default     = {}
}
