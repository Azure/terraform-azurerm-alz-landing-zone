<!-- BEGIN_TF_DOCS -->
# Landing zone virtual network submodule

## Overview

Creates multiple virtual networks in the supplied subscription.
Optionally:

- Creates bi-directional peering and/or a virtual WAN connection
- Creates peerings between the virtual networks (mesh peering)

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

See documentation for optional parameters.

```terraform
module "virtualnetwork" {
  source  = "Azure/lz-vending/azurerm/modules/virtualnetwork"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  subscription_id = "00000000-0000-0000-0000-000000000000"
  virtual_networks = {
    vnet1 = {
      name                = "myvnet"
      address_space       = ["192.168.0.0/24", "10.0.0.0/24"]
      location            = "westeurope"
      resource_group_name = "myrg"
    },
    vnet2 = {
      name                = "myvnet2"
      address_space       = ["192.168.1.0/24", "10.0.1.0/24"]
      location            = "northeurope"
      resource_group_name = "myrg2"
    }
  }
}
```

## Documentation
<!-- markdownlint-disable MD033 -->

## Requirements

The following requirements are needed by this module:

- <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) (~> 1.8)

- <a name="requirement_azapi"></a> [azapi](#requirement\_azapi) (~> 2.2)

## Modules

No modules.

<!-- markdownlint-disable MD013 -->
## Required Inputs

The following input variables are required:

### <a name="input_subscription_id"></a> [subscription\_id](#input\_subscription\_id)

Description: The subscription ID of the subscription to create the virtual network in.

Type: `string`

### <a name="input_virtual_networks"></a> [virtual\_networks](#input\_virtual\_networks)

Description: A map of the virtual networks to create. The map key must be known at the plan stage, e.g. must not be calculated and known only after apply.

### Required fields

- `name`: The name of the virtual network. [required]
- `address_space`: The address space of the virtual network as a list of strings in CIDR format, e.g. `["192.168.0.0/24", "10.0.0.0/24"]`. [required]
- `resource_group_name`: The name of the resource group to create the virtual network in. The default is that the resource group will be created by this module. [required]

### DNS servers

- `dns_servers`: A list of DNS servers to use for the virtual network, e.g. `["192.168.0.1", "10.0.0.1]`. If empty will use the Azure default DNS. [optional - default empty list]  
DNS. [optional - default empty list]

### DDOS protection plan

- `ddos_protection_enabled`: Whether to enable ddos protection. [optional]
- `ddos_protection_plan_id`: The resource ID of the protection plan to attach the vnet. [optional - but required if ddos\_protection\_enabled is `true`]

### Location

- `location`: The location of the virtual network (and resource group if creation is enabled). [optional, will use `var.location` if not specified or empty string]

> Note at least one of `location` or `var.location` must be specified.
> If both are empty then the module will fail.

### Hub network peering values

The following values configure bi-directional hub & spoke peering for the given virtual network.

- `hub_peering_enabled`: Whether to enable hub peering. [optional - default `false`]
- `hub_network_resource_id`: The resource ID of the hub network to peer with. [optional - but required if hub\_peering\_enabled is `true`]
- `hub_peering_name_tohub`: The name of the peering to the hub network. [optional - leave empty to use calculated name]
- `hub_peering_name_fromhub`: The name of the peering from the hub network. [optional - leave empty to use calculated name]
- `hub_peering_use_remote_gateways`: Whether to use remote gateways for the hub peering. [optional - default true]

### Mesh peering values

Mesh peering is the capability to create a bi-directional peerings between all supplied virtual networks in `var.virtual_networks`.  
Peerings will only be created between virtual networks with the `mesh_peering_enabled` value set to `true`.

- `mesh_peering_enabled`: Whether to enable mesh peering for this virtual network. Must be enabled on more than one virtual network for any peerings to be created. [optional]
- `mesh_peering_allow_forwarded_traffic`: Whether to allow forwarded traffic for the mesh peering. [optional - default false]

### Resource group values

The default is that a resource group will be created for each resource\_group\_name specified in the `var.virtual_networks` map.  
It is possible to use a pre-existing resource group by setting `resource_group_creation_enabled` to `false`.  
We recommend using resource groups aligned to the region of the virtual network,  
however if you want multiple virtual networks in more than one location to share a resource group,  
only one of the virtual networks should have `resource_group_creation_enabled` set to `true`.

- `resource_group_creation_enabled`: Whether to create a resource group for the virtual network. [optional - default `true`]
- `resource_group_lock_enabled`: Whether to create a `CanNotDelete` resource lock on the resource group. [optional - default `true`]
- `resource_group_lock_name`: The name of the resource lock. [optional - leave empty to use calculated name]
- `resource_group_tags`: A map of tags to apply to the resource group, e.g. `{ mytag = "myvalue", mytag2 = "myvalue2" }`. [optional - default empty]

### Virtual WAN values

- `vwan_associated_routetable_resource_id`: The resource ID of the route table to associate with the virtual network. [optional - leave empty to use `defaultRouteTable` on hub]
- `vwan_connection_enabled`: Whether to create a connection to a Virtual WAN. [optional - default false]
- `vwan_connection_name`: The name of the connection to the Virtual WAN. [optional - leave empty to use calculated name]
- `vwan_hub_resource_id`: The resource ID of the hub to connect to. [optional - but required if vwan\_connection\_enabled is `true`]
- `vwan_propagated_routetables_labels`: A list of labels of route tables to propagate to the virtual network. [optional - leave empty to use `["default"]`]
- `vwan_propagated_routetables_resource_ids`: A list of resource IDs of route tables to propagate to the virtual network. [optional - leave empty to use `defaultRouteTable` on hub]
- `vwan_security_configuration`: A map of security configuration values for VWAN hub connection - see below. [optional - default empty]
  - `secure_internet_traffic`: Whether to forward internet-bound traffic to the destination specified in the routing policy. [optional - default `false`]
  - `secure_private_traffic`: Whether to all internal traffic to the destination specified in the routing policy. Not compatible with `routing_intent_enabled`. [optional - default `false`]
  - `routing_intent_enabled`: Enable to use with a Virtual WAN hub with routing intent enabled. Routing intent on hub is configured outside this module. [optional - default `false`]

### Tags

- `tags`: A map of tags to apply to the virtual network. [optional - default empty]

Type:

```hcl
map(object({
    name                = string
    address_space       = list(string)
    resource_group_name = string

    location = optional(string, "")

    dns_servers = optional(list(string), [])

    ddos_protection_enabled = optional(bool, false)
    ddos_protection_plan_id = optional(string, "")

    hub_network_resource_id         = optional(string, "")
    hub_peering_enabled             = optional(bool, false)
    hub_peering_direction           = optional(string, "both")
    hub_peering_name_tohub          = optional(string, "")
    hub_peering_name_fromhub        = optional(string, "")
    hub_peering_use_remote_gateways = optional(bool, true)

    mesh_peering_enabled                 = optional(bool, false)
    mesh_peering_allow_forwarded_traffic = optional(bool, false)

    resource_group_creation_enabled = optional(bool, true)
    resource_group_lock_enabled     = optional(bool, true)
    resource_group_lock_name        = optional(string, "")
    resource_group_tags             = optional(map(string), {})

    vwan_associated_routetable_resource_id   = optional(string, "")
    vwan_connection_enabled                  = optional(bool, false)
    vwan_connection_name                     = optional(string, "")
    vwan_hub_resource_id                     = optional(string, "")
    vwan_propagated_routetables_labels       = optional(list(string), [])
    vwan_propagated_routetables_resource_ids = optional(list(string), [])
    vwan_security_configuration = optional(object({
      secure_internet_traffic = optional(bool, false)
      secure_private_traffic  = optional(bool, false)
      routing_intent_enabled  = optional(bool, false)
    }), {})

    tags = optional(map(string), {})
  }))
```

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_location"></a> [location](#input\_location)

Description: The default location of resources created by this module.  
Virtual networks will be created in this location unless overridden by the `location` attribute.

Type: `string`

Default: `""`

## Resources

The following resources are used by this module:

- [azapi_resource.peering_hub_inbound](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.peering_hub_outbound](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.peering_mesh](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.rg](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.rg_lock](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.vhubconnection](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.vhubconnection_routing_intent](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.vnet](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) (resource)

## Outputs

The following outputs are exported:

### <a name="output_resource_group_resource_ids"></a> [resource\_group\_resource\_ids](#output\_resource\_group\_resource\_ids)

Description: The created resource group IDs, expressed as a map.

### <a name="output_virtual_network_resource_ids"></a> [virtual\_network\_resource\_ids](#output\_virtual\_network\_resource\_ids)

Description: The created virtual network resource IDs, expressed as a map.

<!-- markdownlint-enable -->
<!-- END_TF_DOCS -->