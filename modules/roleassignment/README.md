<!-- BEGIN_TF_DOCS -->
# Landing zone role assignment submodule

## Overview

Creates a role assignment at subscription or lower scope.
Module is designed to be instantiated many times, once per role assignment.

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

```terraform
module "roleassignment" {
  source          = "Azure/lz-vending/azurerm/modules/roleassignment"
  version         = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints
  role_definition = "Owner"
  scope           = "/subscriptions/00000000-0000-0000-0000-000000000000"
  principal_id    = "00000000-0000-0000-0000-000000000000"
}
```

## Documentation
<!-- markdownlint-disable MD033 -->

## Requirements

The following requirements are needed by this module:

- <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) (~> 1.3)

- <a name="requirement_azurerm"></a> [azurerm](#requirement\_azurerm) (~> 3.7)

## Modules

No modules.

<!-- markdownlint-disable MD013 -->
## Required Inputs

The following input variables are required:

### <a name="input_role_assignment_condition"></a> [role\_assignment\_condition](#input\_role\_assignment\_condition)

Description: (Optional) The condition that limits the resources that the role can be assigned to.

Type: `string`

### <a name="input_role_assignment_condition_version"></a> [role\_assignment\_condition\_version](#input\_role\_assignment\_condition\_version)

Description: The version of the condition. Possible values are `""`, 1.0 or 2.0. If `""`, null will be set in role\_assignment\_condition and role\_assignment\_condition\_version.

Type: `string`

### <a name="input_role_assignment_definition"></a> [role\_assignment\_definition](#input\_role\_assignment\_definition)

Description: Either the role definition resource id, e.g. `/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c`.  
Or, the role definition name, e.g. `Contributor`.

Type: `string`

### <a name="input_role_assignment_principal_id"></a> [role\_assignment\_principal\_id](#input\_role\_assignment\_principal\_id)

Description: The principal (object) ID of the role assignment.  
Note, for a service principal, this is not the application id.

Can be user, group or service principal.

Type: `string`

### <a name="input_role_assignment_scope"></a> [role\_assignment\_scope](#input\_role\_assignment\_scope)

Description: The scope of the role assignment.

Must begin with `/subscriptions/{subscription-id}` to avoid accidentally creating a role assignment at higher scopes.

Type: `string`

## Optional Inputs

No optional inputs.

## Resources

The following resources are used by this module:

- [azurerm_role_assignment.this](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/role_assignment) (resource)

## Outputs

The following outputs are exported:

### <a name="output_role_assignment_id"></a> [role\_assignment\_id](#output\_role\_assignment\_id)

Description: The Azure resource id of the created role assignment.

### <a name="output_role_assignment_name"></a> [role\_assignment\_name](#output\_role\_assignment\_name)

Description: The Azure name (uuid) of the created role assignment.

<!-- markdownlint-enable -->
<!-- END_TF_DOCS -->