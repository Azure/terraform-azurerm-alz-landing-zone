locals {
  aad_groups_role_assignments = { for ra in flatten([
    for k_group, v_group in var.aad_groups : [
      for k_role, v_role in v_group.role_assignments : {
        group_key       = k_group
        ra_key          = k_role
        role_assignment = v_role
        ignore_changes  = try(v_group.ignore_owner_and_member_changes, false)
      }
    ]
  ]) : "${ra.group_key}-${ra.ra_key}" => ra }

  role_definition_resource_substring = "/providers/Microsoft.Authorization/roleDefinitions"
}
