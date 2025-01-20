<!-- BEGIN_TF_DOCS -->
# Landing zone route table submodule

## Overview

Creates multiple route tables in the supplied subscription.
Optionally:

- Creates routes within the route tables

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

See documentation for optional parameters.

```terraform
module "routetable" {
  source  = "Azure/lz-vending/azurerm/modules/routetable"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  subscription_id = "00000000-0000-0000-0000-000000000000"
  route_tables = {
    rt1 = {
      name                   = "myroutetable"
      address_prefix         = ["192.168.0.0/24"]
      next_hop_in_ip_address = "192.168.0.5"
      next_hop_type          = "VirtualAppliance"
    },
    rt2 = {
      name           = "myroutetable2"
      address_prefix = "GatewayManager"
      next_hop_type  = "Internet"
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

### <a name="input_location"></a> [location](#input\_location)

Description: The location of the route table.

Type: `string`

### <a name="input_name"></a> [name](#input\_name)

Description: The name of the route table to create.

Type: `string`

### <a name="input_resource_group_name"></a> [resource\_group\_name](#input\_resource\_group\_name)

Description: The name of the resource group to create the virtual network in.  
The resource group must exist, this module will not create it.

Type: `string`

### <a name="input_subscription_id"></a> [subscription\_id](#input\_subscription\_id)

Description: The subscription ID of the subscription to create the virtual network in.

Type: `string`

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_bgp_route_propagation_enabled"></a> [bgp\_route\_propagation\_enabled](#input\_bgp\_route\_propagation\_enabled)

Description: Whether BGP route propagation is enabled.

Type: `bool`

Default: `true`

### <a name="input_routes"></a> [routes](#input\_routes)

Description: A list of objects defining route tables and their associated routes to be created:

- `name` (required): The name of the route.
- `address_prefix` (required): The address prefix for the route.
- `next_hop_type` (required): The type of next hop for the route.
- `next_hop_in_ip_address` (required): The next hop IP address for the route.

Type:

```hcl
list(object({
    name                   = string
    address_prefix         = string
    next_hop_type          = string
    next_hop_in_ip_address = string
  }))
```

Default: `[]`

### <a name="input_tags"></a> [tags](#input\_tags)

Description: A map of tags to assign to the route table.

Type: `map(string)`

Default: `{}`

## Resources

The following resources are used by this module:

- [azapi_resource.route_table](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) (resource)

## Outputs

The following outputs are exported:

### <a name="output_route_table_resource_id"></a> [route\_table\_resource\_id](#output\_route\_table\_resource\_id)

Description: The created route table ID.

<!-- markdownlint-enable -->
<!-- END_TF_DOCS -->