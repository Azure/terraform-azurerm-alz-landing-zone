data "azurerm_client_config" "current" {}

resource "azuread_group" "this" {
  for_each = { for key, value in var.aad_groups : key => value if !value.ignore_owner_and_member_changes }

  display_name = each.value.name

  administrative_unit_ids = each.value.administrative_unit_ids
  assignable_to_role      = each.value.assignable_to_role
  description             = each.value.description

  security_enabled        = true
  members                 = each.value.members.object_ids
  prevent_duplicate_names = each.value.prevent_duplicate_names
  owners                  = try(each.value.add_deployment_user_as_owner, false) ? setunion(each.value.owners.object_ids, [data.azurerm_client_config.current.object_id]) : each.value.owners.object_ids
  visibility              = "Private"
}

resource "azuread_group" "ignore_owner_and_member_changes" {
  for_each = { for key, value in var.aad_groups : key => value if value.ignore_owner_and_member_changes }

  display_name = each.value.name

  administrative_unit_ids = each.value.administrative_unit_ids
  assignable_to_role      = each.value.assignable_to_role
  description             = each.value.description

  security_enabled        = true
  members                 = each.value.members.object_ids
  prevent_duplicate_names = each.value.prevent_duplicate_names
  owners                  = try(each.value.add_deployment_user_as_owner, false) ? setunion(each.value.owners.object_ids, [data.azurerm_client_config.current.object_id]) : each.value.owners.object_ids
  visibility              = "Private"

  lifecycle {
    ignore_changes = [
      members,
      owners
    ]
  }
}

resource "azurerm_role_assignment" "groups" {
  for_each = local.aad_groups_role_assignments

  principal_id                           = each.value.ignore_changes ? azuread_group.ignore_owner_and_member_changes[each.value.group_key].object_id : azuread_group.this[each.value.group_key].object_id
  scope                                  = "/subscriptions/${var.subscription_id}${each.value.role_assignment.relative_scope}"
  condition                              = each.value.role_assignment.condition
  condition_version                      = each.value.role_assignment.condition_version
  delegated_managed_identity_resource_id = each.value.role_assignment.delegated_managed_identity_resource_id
  role_definition_id                     = strcontains(lower(each.value.role_assignment.definition), lower(local.role_definition_resource_substring)) ? each.value.role_assignment.definition : null
  role_definition_name                   = strcontains(lower(each.value.role_assignment.definition), lower(local.role_definition_resource_substring)) ? null : each.value.role_assignment.definition
  skip_service_principal_aad_check       = each.value.role_assignment.skip_service_principal_aad_check
}
