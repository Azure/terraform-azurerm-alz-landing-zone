module "usermanagedidentity" {
  source   = "./modules/usermanagedidentity"
  count    = var.umi_enabled ? 1 : 0
  name     = var.umi_name
  location = var.location
  tags     = var.umi_tags

  resource_group_creation_enabled = var.umi_resource_group_creation_enabled
  resource_group_name             = var.umi_resource_group_name
  resource_group_lock_enabled     = var.umi_resource_group_lock_enabled
  resource_group_lock_name        = var.umi_resource_group_lock_name
  resource_group_tags             = var.umi_resource_group_tags

  subscription_id = local.subscription_id

  federated_credentials_advanced        = var.umi_federated_credentials_advanced
  federated_credentials_github          = var.umi_federated_credentials_github
  federated_credentials_terraform_cloud = var.umi_federated_credentials_terraform_cloud

  depends_on = [
    module.resourcegroup
  ]
}
