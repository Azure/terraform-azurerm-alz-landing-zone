terraform {
  required_version = ">= 1.2.0"
  required_providers {
    azapi = {
      source  = "azure/azapi"
      version = ">= 0.3.0"
    }
  }
}

provider "azurerm" {
  features {}
}

variable "subscription_alias_billing_scope" {
  type = string
}

variable "subscription_alias_management_group_id" {
  type = string
}

variable "subscription_alias_name" {
  type = string
}

variable "subscription_alias_display_name" {
  type = string
}

variable "subscription_alias_workload" {
  type = string
}

resource "azapi_resource" "mg" {
  type      = "Microsoft.Management/managementGroups@2021-04-01"
  parent_id = "/"
  name      = var.subscription_alias_management_group_id
}

module "lz_test" {
  source                                 = "../../../../modules/subscription"
  subscription_alias_name                = var.subscription_alias_name
  subscription_alias_display_name        = var.subscription_alias_display_name
  subscription_alias_workload            = "DevTest"
  subscription_alias_management_group_id = azapi_resource.mg.name
  subscription_alias_billing_scope       = var.subscription_alias_billing_scope
}

output "subscription_id" {
  value = module.lz_test.subscription_id
}
