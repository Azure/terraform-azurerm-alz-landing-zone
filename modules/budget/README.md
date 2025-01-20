<!-- BEGIN_TF_DOCS -->
# Landing zone budget submodule

## Overview

Creates a budget in Azure. Designed to be instantiated multiple times to create multiple budgets.

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

```terraform
module "budget" {
  source  = "Azure/lz-vending/azurerm/modules/roleassignment"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  budget_name       = "budget1"
  budget_amount     = 100
  budget_scope      = "/subscriptions/00000000-0000-0000-0000-000000000000"
  budget_time_grain = "Monthly"
  budget_time_period = {
    start_date = "2024-01-01"
    end_date   = "2025-01-01"
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

### <a name="input_budget_amount"></a> [budget\_amount](#input\_budget\_amount)

Description: The total amount of cost to track with the budget.

Type: `number`

### <a name="input_budget_name"></a> [budget\_name](#input\_budget\_name)

Description: The name of the budget.

Type: `string`

### <a name="input_budget_scope"></a> [budget\_scope](#input\_budget\_scope)

Description: The scope of the budget.

Type: `string`

### <a name="input_budget_time_grain"></a> [budget\_time\_grain](#input\_budget\_time\_grain)

Description: The time grain of the budget.

Type: `string`

### <a name="input_budget_time_period"></a> [budget\_time\_period](#input\_budget\_time\_period)

Description: The time period of the budget.

Type:

```hcl
object({
    start_date = string
    end_date   = string
  })
```

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_budget_notifications"></a> [budget\_notifications](#input\_budget\_notifications)

Description: The notifications for the budget.

Type:

```hcl
map(object({
    enabled        = bool
    operator       = string
    threshold      = number
    threshold_type = optional(string, "Actual")
    contact_emails = optional(list(string), [])
    contact_roles  = optional(list(string), [])
    contact_groups = optional(list(string), [])
    locale         = optional(string, "en-us")
  }))
```

Default: `{}`

## Resources

The following resources are used by this module:

- [azapi_resource.budget](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)

## Outputs

The following outputs are exported:

### <a name="output_budget_resource_id"></a> [budget\_resource\_id](#output\_budget\_resource\_id)

Description: The Azure resource id of the created budget.

<!-- markdownlint-enable -->
<!-- END_TF_DOCS -->