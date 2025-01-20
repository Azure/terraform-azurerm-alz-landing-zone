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
