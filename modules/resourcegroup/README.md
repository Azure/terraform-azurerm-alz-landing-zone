<!-- BEGIN_TF_DOCS -->
# Landing zone resource group submodule

## Overview

Creates the supplied resource group in the specified location and subscription.

Useful in a subscription vending scenario to pre-create the `NetworkWatcherRG` so we may delete it prior to subscription cancellation. Also useful for pre-creating resource groups for the purposes of RBAC delegation.

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

See documentation for optional parameters.

```terraform
module "resourcegroups" {
  source  = "Azure/lz-vending/azurerm/modules/resourcegroups"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  subscription_id = "00000000-0000-0000-0000-000000000000"
  location        = "eastus"
  name            = "rg-test"
}
```

## Documentation
<!-- markdownlint-disable MD033 -->

## Requirements

The following requirements are needed by this module:

- <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) (~> 1.3)

- <a name="requirement_azapi"></a> [azapi](#requirement\_azapi) (~> 1.13)

## Modules

No modules.

<!-- markdownlint-disable MD013 -->
## Required Inputs

The following input variables are required:

### <a name="input_location"></a> [location](#input\_location)

Description: The Azure region to deploy resources into. E.g. `eastus`

Type: `string`

### <a name="input_resource_group_name"></a> [resource\_group\_name](#input\_resource\_group\_name)

Description: The name of the resource group E.g. `rg-test`

Type: `string`

### <a name="input_subscription_id"></a> [subscription\_id](#input\_subscription\_id)

Description: The ID of the subscription to deploy resources into. E.g. `00000000-0000-0000-0000-000000000000`

Type: `string`

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_tags"></a> [tags](#input\_tags)

Description: Map of tags to be applied to the resource group

Type: `map(string)`

Default: `{}`

## Resources

The following resources are used by this module:

- [azapi_resource.rg](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)

## Outputs

No outputs.

<!-- markdownlint-enable -->
<!-- END_TF_DOCS -->