# The azurerm_subscription resource represents the subscription alias that is being created.
resource "azurerm_subscription" "this" {
  count             = var.subscription_alias_enabled && !var.subscription_use_azapi ? 1 : 0
  subscription_name = var.subscription_display_name
  alias             = var.subscription_alias_name
  billing_scope_id  = var.subscription_billing_scope
  workload          = var.subscription_workload
  tags              = var.subscription_tags
}

# This resource ensures that we can manage the management group for the subscription
# throughout its lifecycle.
resource "azurerm_management_group_subscription_association" "this" {
  count               = var.subscription_management_group_association_enabled && !var.subscription_use_azapi ? 1 : 0
  management_group_id = "/providers/Microsoft.Management/managementGroups/${var.subscription_management_group_id}"
  subscription_id     = "/subscriptions/${local.subscription_id}"
}

resource "azapi_resource" "subscription" {
  count = (var.subscription_alias_enabled && var.subscription_use_azapi) ? 1 : 0

  type      = "Microsoft.Subscription/aliases@2021-10-01"
  name      = var.subscription_alias_name
  parent_id = "/"

  body = {
    properties = {
      displayName  = var.subscription_display_name
      workload     = var.subscription_workload
      billingScope = var.subscription_billing_scope
      additionalProperties = {
        managementGroupId = var.subscription_management_group_association_enabled ? "/providers/Microsoft.Management/managementGroups/${var.subscription_management_group_id}" : null
        tags              = var.subscription_tags
      }
    }
  }
  response_export_values = ["properties.subscriptionId"]
  lifecycle {
    ignore_changes = [
      body,
      name
    ]
  }
}

resource "terraform_data" "replacement" {
  count = (var.subscription_management_group_association_enabled && var.subscription_use_azapi) ? 1 : 0

  input = local.is_subscription_associated_to_management_group
}

resource "time_sleep" "wait_for_subscription_before_subscription_operations" {
  count = (var.subscription_alias_enabled && var.subscription_use_azapi) ? 1 : 0

  create_duration  = var.wait_for_subscription_before_subscription_operations.create
  destroy_duration = var.wait_for_subscription_before_subscription_operations.destroy

  depends_on = [
    azapi_resource.subscription
  ]
}

resource "azapi_resource_action" "subscription_association" {
  count = (var.subscription_management_group_association_enabled && var.subscription_use_azapi) ? 1 : 0

  type        = "Microsoft.Management/managementGroups/subscriptions@2021-04-01"
  resource_id = "/providers/Microsoft.Management/managementGroups/${var.subscription_management_group_id}/subscriptions/${local.subscription_id}"
  method      = "PUT"

  lifecycle {
    replace_triggered_by = [terraform_data.replacement]
  }

  depends_on = [
    time_sleep.wait_for_subscription_before_subscription_operations
  ]
}

resource "azapi_update_resource" "subscription_tags" {
  count = (var.subscription_alias_enabled && var.subscription_use_azapi) || (var.subscription_id != "" && var.subscription_update_existing) ? 1 : 0

  type        = "Microsoft.Resources/tags@2022-09-01"
  resource_id = "/subscriptions/${local.subscription_id}/providers/Microsoft.Resources/tags/default"
  body = {
    properties = {
      tags = var.subscription_tags
    }
  }

  depends_on = [
    time_sleep.wait_for_subscription_before_subscription_operations
  ]
}

resource "azapi_resource_action" "subscription_rename" {
  count = (var.subscription_alias_enabled && var.subscription_use_azapi) || (var.subscription_id != "" && var.subscription_update_existing) ? 1 : 0

  type        = "Microsoft.Resources/subscriptions@2021-10-01"
  resource_id = "/subscriptions/${local.subscription_id}"
  method      = "POST"
  action      = "providers/Microsoft.Subscription/rename"
  body = {
    subscriptionName = var.subscription_display_name
  }

  depends_on = [
    time_sleep.wait_for_subscription_before_subscription_operations
  ]
}

resource "azapi_resource_action" "subscription_cancel" {
  count = (var.subscription_alias_enabled && var.subscription_use_azapi) ? 1 : 0

  type        = "Microsoft.Resources/subscriptions@2021-10-01"
  resource_id = "/subscriptions/${local.subscription_id}"
  method      = "POST"
  action      = "providers/Microsoft.Subscription/cancel"
  when        = "destroy"

  depends_on = [
    time_sleep.wait_for_subscription_before_subscription_operations
  ]
}
