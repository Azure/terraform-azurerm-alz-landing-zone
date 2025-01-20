terraform {
  required_version = ">= 1.3.0"
  required_providers {
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.47"
    }
  }
}
