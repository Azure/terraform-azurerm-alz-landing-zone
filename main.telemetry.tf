# This is the root module telemetry deployment that is only created if telemetry is enabled.
# It is deployed to the created or supplied subscription
resource "azapi_resource" "telemetry_root" {
  count     = var.disable_telemetry ? 0 : 1
  parent_id = local.subscription_resource_id
  name      = local.telem_root_arm_deployment_name
  type      = "Microsoft.Resources/deployments@2021-04-01"
  location  = var.location
  body = {
    properties = {
      mode     = "Incremental"
      template = local.telem_arm_subscription_template
    }
  }
}
