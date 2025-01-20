<!-- BEGIN_TF_DOCS -->
# Landing zone resource provider submodule

## Overview

Registers resource providers and features. Must be performed after all other provisioning is complete.

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

```terraform
module "resourceproviders" {
  source  = "Azure/lz-vending/azurerm/modules/resourceprovider"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  resource_provider = "Microsoft.PowerBI"
  features          = ["DailyPrivateLinkServicesForPowerBI"]
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

### <a name="input_resource_provider"></a> [resource\_provider](#input\_resource\_provider)

Description: The resource provider namespace, e.g. `Microsoft.Compute`.

Type: `string`

### <a name="input_subscription_id"></a> [subscription\_id](#input\_subscription\_id)

Description: The subscription id to register the resource providers in.

Type: `string`

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_features"></a> [features](#input\_features)

Description: The resource provider features to register, e.g. [`MyFeature`]

Type: `set(string)`

Default: `[]`

## Resources

The following resources are used by this module:

- [azapi_resource_action.resource_provider_feature_registration](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource_action) (resource)
- [azapi_resource_action.resource_provider_registration](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource_action) (resource)

## Outputs

No outputs.

<!-- markdownlint-enable -->
<!-- END_TF_DOCS -->