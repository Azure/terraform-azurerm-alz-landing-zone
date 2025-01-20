<!-- BEGIN_TF_DOCS -->
# Landing zone Entra ID (AAD) Group submodule

## Overview

Creates groups in Entra ID and role assignments for resources.

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

```terraform
module "aadgroup" {
  source  = "Azure/lz-vending/azurerm/modules/aadgroup"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  aad_groups = {
    contributor_group = {
      name = "my-ad-group-name"

      # optional parameters
      description = "the description for my ad group"
      members = {
        object_ids = [
          "e64a9602-6a56-4d45-a4b0-7a7fe605f89d",
          "8c537ad4-0289-41f5-84b7-3d1450c04643",
        ]
      }
      owners = {
        object_ids = ["1f32f09d-bae9-4f02-8905-1ae0a5d97d2f"]
      }

      # optional role assignment
      role_assignments = {
        rg_contributor = {
          definition     = "Contributor"
          relative_scope = "/resourceGroups/rg-some-resource-group"
        }
      }

      # optionally tell Terraform to ignore changes to owners & members
      ignore_owner_and_member_changes = true

      # optionally add the deployment user to the owners to allow subsequent membership updates
      add_deployment_user_as_owner = true
    }
  }

  subscription_id = "00000000-0000-0000-0000-000000000000"
}
```

## Documentation
<!-- markdownlint-disable MD033 -->

## Requirements

The following requirements are needed by this module:

- <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) (>= 1.3.0)

- <a name="requirement_azuread"></a> [azuread](#requirement\_azuread) (~> 2.47)

## Modules

No modules.

<!-- markdownlint-disable MD013 -->
## Required Inputs

The following input variables are required:

### <a name="input_aad_groups"></a> [aad\_groups](#input\_aad\_groups)

Description: A map defining the configuration for an Entra ID (Azure Active Directory) group.

- `name` - The display name of the group.

**Optional Parameters:**

- `administrative_unit_ids` - (optional) A list of object IDs of administrative units for group membership.
- `assignable_to_role` - (optional) Whether the group can be assigned to an Azure AD role (default: false).
- `description` - (optional) The description for the group (default: "").
- `ignore_owner_and_member_changes` - (optional) If true, changes to ownership and membership will be ignored (default: false).
- `members` - (optional) A set of members (Users, Groups, or Service Principals).
- `owners` - (optional) A list of object IDs of owners (Users or Service Principals) (default: current user).
- `prevent_duplicate_names` - (optional) If true, throws an error on duplicate names (default: true).
- `add_deployment_user_as_owner` - (optional) If true, adds the current service principal the terraform deployment is running as to the owners, useful if further management by terraform is required (default: false).

- `role_assignments` - (optional) A map defining role assignments for the group.
  - `definition` - The name of the role to assign.
  - `relative_scope` - The scope of the role assignment relative to the subscription
  - `description` - (optional) Description for the role assignment.
  - `skip_service_principal_aad_check` - (optional) If true, skips the Azure AD check for service principal (default: false).
  - `condition` - (optional) The condition for the role assignment.
  - `condition_version` - (optional) The condition version for the role assignment.
  - `delegated_managed_identity_resource_id` - (optional) The resource ID of the delegated managed identity.

Type:

```hcl
map(object({
    name = string

    administrative_unit_ids         = optional(list(string), null)
    assignable_to_role              = optional(bool, false)
    description                     = optional(string, null)
    ignore_owner_and_member_changes = optional(bool, false)
    members                         = optional(map(list(string)), null)
    owners                          = optional(map(list(string)), null)
    prevent_duplicate_names         = optional(bool, true)
    add_deployment_user_as_owner    = optional(bool, false)
    role_assignments = optional(map(object({
      definition                             = string
      relative_scope                         = string
      description                            = optional(string, null)
      skip_service_principal_aad_check       = optional(bool, false)
      condition                              = optional(string, null)
      condition_version                      = optional(string, null)
      delegated_managed_identity_resource_id = optional(string, null)
    })), {})
  }))
```

### <a name="input_subscription_id"></a> [subscription\_id](#input\_subscription\_id)

Description: The subscription ID of the subscriptions where group role assignments are applied.

Type: `string`

## Optional Inputs

No optional inputs.

## Resources

The following resources are used by this module:

- [azuread_group.ignore_owner_and_member_changes](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/resources/group) (resource)
- [azuread_group.this](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/resources/group) (resource)
- [azurerm_role_assignment.groups](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/role_assignment) (resource)
- [azurerm_client_config.current](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/data-sources/client_config) (data source)

## Outputs

No outputs.

<!-- markdownlint-enable -->
<!-- END_TF_DOCS -->