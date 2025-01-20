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
