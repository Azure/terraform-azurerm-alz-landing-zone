resource "azapi_resource" "rg" {
  parent_id = "/subscriptions/${var.subscription_id}"
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  name      = var.resource_group_name
  location  = var.location
  body      = {}
  tags      = var.tags
}
