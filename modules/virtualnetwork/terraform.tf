terraform {
  required_version = ">= 1.3.0"
  required_providers {
    azapi = {
      source  = "azure/azapi"
      version = "~> 1.14"
    }
  }
}
