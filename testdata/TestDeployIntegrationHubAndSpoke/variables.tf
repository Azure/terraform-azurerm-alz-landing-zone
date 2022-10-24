variable "subscription_alias_enabled" {
  type = bool
}

variable "subscription_billing_scope" {
  type = string
}

variable "subscription_display_name" {
  type = string
}

variable "subscription_alias_name" {
  type = string
}

variable "subscription_workload" {
  type = string
}

variable "virtual_network_enabled" {
  type = string
}

variable "virtual_networks" {
  type = any
}
variable "location" {
  type = string
}

variable "role_assignment_enabled" {
  type = bool
}
