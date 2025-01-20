variable "aadgroup_enabled" {
  type        = bool
  description = <<DESCRIPTION
Whether to create Entra ID (Azure AD) groups.
If enabled, supply the list of aadgroups in `var.aadgroups`.
DESCRIPTION
  default     = false
}

variable "aad_groups" {
  type = map(object({
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
  default     = {}
  nullable    = false
  description = <<DESCRIPTION
A map defining the configuration for Entra ID (AAD) groups. 

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
DESCRIPTION
}
