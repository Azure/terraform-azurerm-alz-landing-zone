data "azapi_client_config" "telemetry" {
  count = var.disable_telemetry ? 0 : 1
}

data "modtm_module_source" "telemetry" {
  count       = var.disable_telemetry ? 0 : 1
  module_path = path.module
}

resource "random_uuid" "telemetry" {
  count = var.disable_telemetry ? 0 : 1
}

resource "modtm_telemetry" "telemetry" {
  count = var.disable_telemetry ? 0 : 1

  tags = {
    subscription_id = local.subscription_id
    tenant_id       = one(data.azurerm_client_config.telemetry).tenant_id
    module_source   = one(data.modtm_module_source.telemetry).module_source
    module_version  = one(data.modtm_module_source.telemetry).module_version
    random_id       = one(random_uuid.telemetry).result
  }
}
