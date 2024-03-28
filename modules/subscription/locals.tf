locals {
  # subscription_id_alias is the id of the newly created subscription, if it exists.
  subscription_id_alias = try(azurerm_subscription.this[0].subscription_id, jsondecode(azapi_resource.subscription[0].output).properties.subscriptionId, null)

  # subscription_id is the id of the newly created subscription, or the id supplied by var.subscription_id.
  subscription_id = coalesce(local.subscription_id_alias, var.subscription_id)
}

locals {
  # Check if subscription is vended.
  is_subscription_vended = (var.subscription_management_group_association_enabled && var.subscription_use_azapi) ? contains(jsondecode(data.azapi_resource_list.subscriptions[0].output).value[*].subscriptionId, local.subscription_id) : true
  # Check for drift between subscription and target management group.
  is_subscription_associated_to_management_group = (var.subscription_management_group_association_enabled && var.subscription_use_azapi) && local.is_subscription_vended ? contains(jsondecode(data.azapi_resource_list.subscription_management_group_association[0].output).value[*].id, "/providers/Microsoft.Management/managementGroups/${var.subscription_management_group_id}/subscriptions/${local.subscription_id}") : true
}

locals {
  # Transform subscription budgets to be able to use them with the API.
  transformed_budgets = {
    for key, value in var.subscription_budgets :
    key => {
      amount            = value.amount
      time_grain        = value.time_grain
      time_period_start = value.time_period_start
      time_period_end   = value.time_period_end
      notifications = {
        for n_key, n_value in value.notifications :
        n_key => {
          enabled       = lookup(n_value, "enabled", false)
          operator      = lookup(n_value, "operator", "GreaterThan")
          threshold     = lookup(n_value, "threshold", 1000)
          thresholdType = lookup(n_value, "threshold_type", "Forecasted")
          contactEmails = lookup(n_value, "contact_emails", [])
          contactRoles  = lookup(n_value, "contact_roles", [])
          contactGroups = lookup(n_value, "contact_groups", [])
          locale        = lookup(n_value, "locale", "en-us")
        }
      }
    }
  }
}
