variable "resource_group_creation_enabled" {
  type        = bool
  description = "Whether to create additional resource groups in the target subscription. Requires `var.resource_groups`."
  default     = false
}

variable "resource_groups" {
  type = map(object({
    name     = string
    location = string
    tags     = optional(map(string), {})
  }))
  description = <<DESCRIPTION
A map of the resource groups to create. The value is an object with the following attributes:

- `name` - the name of the resource group
- `location` - the location of the resource group
- `tags` - (optional) a map of type string

We recommend that you include an entry to create the NetworkWatcherRG resource group so that this is managed by Terraform.
DESCRIPTION
  nullable    = false
  default     = {}
}
