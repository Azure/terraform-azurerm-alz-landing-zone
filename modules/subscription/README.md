<!-- BEGIN_TF_DOCS -->
# Landing zone subscription submodule

## Overview

Creates a subscription alias, and optionally manages management group association for the resulting subscription.

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

```terraform
module "subscription" {
  source  = "Azure/lz-vending/azurerm/modules/subscription"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  subscription_alias_billing_scope       = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
  subscription_alias_display_name        = "my-subscription-display-name"
  subscription_alias_name                = "my-subscription-alias"
  subscription_alias_workload            = "Production"
  subscription_alias_management_group_id = "mymg"
}
```

## Documentation
<!-- markdownlint-disable MD033 -->

## Requirements

The following requirements are needed by this module:

- <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) (~> 1.4)

- <a name="requirement_azapi"></a> [azapi](#requirement\_azapi) (~> 1.11)

- <a name="requirement_azurerm"></a> [azurerm](#requirement\_azurerm) (~> 3.7)

- <a name="requirement_time"></a> [time](#requirement\_time) (~> 0.9)

## Modules

No modules.

<!-- markdownlint-disable MD013 -->
## Required Inputs

No required inputs.

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_subscription_alias_enabled"></a> [subscription\_alias\_enabled](#input\_subscription\_alias\_enabled)

Description: Whether to create a new subscription using the subscription alias resource.

If enabled, the following must also be supplied:

- `subscription_alias_name`
- `subscription_display_name`
- `subscription_billing_scope`
- `subscription_workload`

Optionally, supply the following to enable the placement of the subscription into a management group:

- `subscription_management_group_id`
- `subscription_management_group_association_enabled`

If disabled, supply the `subscription_id` variable to use an existing subscription instead.

> **Note**: When the subscription is destroyed, this module will try to remove the NetworkWatcherRG resource group using `az cli`.
> This requires the `az cli` tool be installed and authenticated.
> If the command fails for any reason, the provider will attempt to cancel the subscription anyway.

Type: `bool`

Default: `false`

### <a name="input_subscription_alias_name"></a> [subscription\_alias\_name](#input\_subscription\_alias\_name)

Description: The name of the subscription alias.

The string must be comprised of a-z, A-Z, 0-9, - and \_.  
The maximum length is 63 characters.

You may also supply an empty string if you do not want to create a new subscription alias.  
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.

Type: `string`

Default: `""`

### <a name="input_subscription_billing_scope"></a> [subscription\_billing\_scope](#input\_subscription\_billing\_scope)

Description: The billing scope for the new subscription alias.

A valid billing scope starts with `/providers/Microsoft.Billing/billingAccounts/` and is case sensitive.

E.g.

- For CustomerLed and FieldLed, e.g. MCA - `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/billingProfiles/{billingProfileName}/invoiceSections/{invoiceSectionName}`
- For PartnerLed, e.g. MPA - `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/customers/{customerName}`
- For Legacy EA - `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/enrollmentAccounts/{enrollmentAccountName}`

You may also supply an empty string if you do not want to create a new subscription alias.  
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.

Type: `string`

Default: `""`

### <a name="input_subscription_display_name"></a> [subscription\_display\_name](#input\_subscription\_display\_name)

Description: The display name of the subscription alias.

The string must be comprised of a-z, A-Z, 0-9, -, \_ and space.  
The maximum length is 64 characters.

You may also supply an empty string if you do not want to create a new subscription alias.  
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.

Type: `string`

Default: `""`

### <a name="input_subscription_id"></a> [subscription\_id](#input\_subscription\_id)

Description: n/a

Type: `string`

Default: `""`

### <a name="input_subscription_management_group_association_enabled"></a> [subscription\_management\_group\_association\_enabled](#input\_subscription\_management\_group\_association\_enabled)

Description: Whether to create the `azurerm_management_group_subscription_association` resource.

If enabled, the `subscription_management_group_id` must also be supplied.

Type: `bool`

Default: `false`

### <a name="input_subscription_management_group_id"></a> [subscription\_management\_group\_id](#input\_subscription\_management\_group\_id)

Description: The destination management group ID for the new subscription.

**Note:** Do not supply the display name.  
The management group ID forms part of the Azure resource ID. E.g.,
`/providers/Microsoft.Management/managementGroups/{managementGroupId}`.

Type: `string`

Default: `""`

### <a name="input_subscription_tags"></a> [subscription\_tags](#input\_subscription\_tags)

Description: A map of tags to assign to the newly created subscription.  
Only valid when `subsciption_alias_enabled` is set to `true`.

Example value:

```terraform
subscription_tags = {
  mytag  = "myvalue"
  mytag2 = "myvalue2"
}
```

Type: `map(string)`

Default: `{}`

### <a name="input_subscription_update_existing"></a> [subscription\_update\_existing](#input\_subscription\_update\_existing)

Description: Whether to update an existing subscription with the supplied tags and display name.  
If enabled, the following must also be supplied:
- `subscription_id`

Type: `bool`

Default: `false`

### <a name="input_subscription_use_azapi"></a> [subscription\_use\_azapi](#input\_subscription\_use\_azapi)

Description: Whether to use the azapi\_resource resource to create the subscription alias. This includes the subscription alias in the management group.

Type: `bool`

Default: `false`

### <a name="input_subscription_workload"></a> [subscription\_workload](#input\_subscription\_workload)

Description: The billing scope for the new subscription alias.

The workload type can be either `Production` or `DevTest` and is case sensitive.

You may also supply an empty string if you do not want to create a new subscription alias.  
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.

Type: `string`

Default: `""`

### <a name="input_wait_for_subscription_before_subscription_operations"></a> [wait\_for\_subscription\_before\_subscription\_operations](#input\_wait\_for\_subscription\_before\_subscription\_operations)

Description: The duration to wait after vending a subscription before performing subscription operations.

Type:

```hcl
object({
    create  = optional(string, "30s")
    destroy = optional(string, "0s")
  })
```

Default: `{}`

## Resources

The following resources are used by this module:

- [azapi_resource.subscription](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource_action.subscription_association](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource_action) (resource)
- [azapi_resource_action.subscription_cancel](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource_action) (resource)
- [azapi_resource_action.subscription_rename](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource_action) (resource)
- [azapi_update_resource.subscription_tags](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/update_resource) (resource)
- [azurerm_management_group_subscription_association.this](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/management_group_subscription_association) (resource)
- [azurerm_subscription.this](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/subscription) (resource)
- [terraform_data.replacement](https://registry.terraform.io/providers/hashicorp/terraform/latest/docs/resources/data) (resource)
- [time_sleep.wait_for_subscription_before_subscription_operations](https://registry.terraform.io/providers/hashicorp/time/latest/docs/resources/sleep) (resource)
- [azapi_resource_list.subscription_management_group_association](https://registry.terraform.io/providers/Azure/azapi/latest/docs/data-sources/resource_list) (data source)
- [azapi_resource_list.subscriptions](https://registry.terraform.io/providers/Azure/azapi/latest/docs/data-sources/resource_list) (data source)

## Outputs

The following outputs are exported:

### <a name="output_management_group_subscription_association_id"></a> [management\_group\_subscription\_association\_id](#output\_management\_group\_subscription\_association\_id)

Description: The management\_group\_subscription\_association\_id output is the ID of the management group subscription association.  
Value will be null if `var.subscription_management_group_association_enabled` is false.

### <a name="output_subscription_id"></a> [subscription\_id](#output\_subscription\_id)

Description: The subscription\_id is the id of the newly created subscription, or that of the supplied var.subscription\_id.  
Value will be null if `var.subscription_id` is blank and `var.subscription_alias_enabled` is false.

### <a name="output_subscription_resource_id"></a> [subscription\_resource\_id](#output\_subscription\_resource\_id)

Description: The subscription\_resource\_id output is the Azure resource id for the newly created subscription.  
Value will be null if `var.subscription_id` is blank and `var.subscription_alias_enabled` is false.

<!-- markdownlint-enable -->
<!-- END_TF_DOCS -->