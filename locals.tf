locals {
  # subscription_module_output_subscription_id is either the output of the subscription module,
  # or if disabled, an empty string.
  # Needed to avoid errors in local.subscription_id when referencing a module instance that does not exists.
  subscription_module_output_subscription_id = try(module.subscription[0].subscription_id, "")

  # subscription_module_output_subscription_id is either the output of the subscription module,
  # or if disabled, an empty string.
  # Needed to avoid errors in local.subscription_id when referencing a module instance that does not exists.
  subscription_module_output_subscription_resource_id = try(module.subscription[0].subscription_resource_id, "")

  # subscription_id is the id of the subscription into which resources will be created.
  # We pick the created sub id first, if it exists, otherwise we pick the subscription_id variable.
  subscription_id = coalesce(local.subscription_module_output_subscription_id, var.subscription_id)

  # subscription_resource_id is the Azure resource id of the subscription id that was supplied in the input variables.
  # If var.subscription_id is empty, then we will return en empty string so that we can correctly coalesce the subscription_resource_id output.
  supplied_subscription_resource_id = var.subscription_id == "" ? "" : "/subscriptions/${var.subscription_id}"

  # subscription_resource_id is the Azure resource id of the subscription into which resources will be created.
  # We use the created sub resource id first, if it exists, otherwise we pick the subscription_id variable.
  # If this is blank then the subscription submodule is disabled an no subscription id has been supplied as an input variable.
  subscription_resource_id = coalesce(local.subscription_module_output_subscription_resource_id, local.supplied_subscription_resource_id)

  # virtual_networks_merged is a map of virtual networks created, if the module has been enabled.
  # This is used in the outputs.tf file to return the virtual network resource ids.
  virtual_network_resource_ids = var.virtual_network_enabled ? module.virtualnetwork[0].virtual_network_resource_ids : {}

  # resource_group_ids is a map of resource groups created, if the module has been enabled.
  # This is used in the outputs.tf file to return the resource group ids.
  virtual_network_resource_group_ids = var.virtual_network_enabled ? module.virtualnetwork[0].resource_group_resource_ids : {}
}
