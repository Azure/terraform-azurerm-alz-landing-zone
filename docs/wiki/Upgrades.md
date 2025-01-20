<!-- markdownlint-disable MD041 -->
This module uses [Semantic Versioning (SemVer)](https://semver.org/s) versioning.

Given a version number `MAJOR.MINOR.PATCH`, we increment the:

* `MAJOR` version when we make incompatible/breaking changes,
* `MINOR` version when we add functionality in a backwards compatible manner, and
* `PATCH` version when we make backwards compatible bug fixes.

## Upgrade process

If you are upgrading to a new `MINOR` or `PATCH` release, you will not see any breaking changes.
If you are using the Terraform registry, you can update the version number for the module.

For a new `MAJOR` release, you will see breaking changes.
We will publish guidance in the release notes on GitHub.

See the [release notes](https://github.com/Azure/terraform-azurerm-lz-vending/releases) for more information.

## Upgrading from v1.x to v2.x

v2 of the module makes large-scale changes to the virtual networking capabilities of the module.
We therefore recommend that you keep any existing instances of the module at v1, and use v2 going forward for new instances.
If you would like multiple vnets in the same subscription using v1 of the module you can use the pattern [described here](https://github.com/Azure/terraform-azurerm-lz-vending/issues/97#issuecomment-1240712419)

## Upgrading from v2.x to v3.x

v3 of the module makes changes the the `role_assignments` variable, changing the format of the variable from a list of objects `list(object({...}))` to a map of objects `map(object({...}))`.
This change fixes [#153](https://github.com/Azure/terraform-azurerm-lz-vending/issues/153).

> **Due to the map key changing, all role assignments will be deleted and re-created.**

By way of explanation, in order to run the `for_each` loop on the role assignments, we need to use either a set or a map.
If using a map, Terraform needs to know all keys at plan time, they cannot be 'known after apply'.
Previously we converted the list of objects into a map of objects, using the `uuidv5()` function to generate predictable map keys from the inputs.
Unfortunately this caused issues when any of the inputs to the `uuidv5()` function were not known at plan time, in this case it was the principal id.

Rather than revert to a set, where ordering can be an issue, we decided to change the input variable to be a map from the outset.
This does mean a small change is required, you must specify a map key. This can be anything but do not use a reference to other object.

### v2.x `role_assignments` syntax

```terraform
module "lz_vending" {
  source  = "..."
  version = "..."

  # (other input variables hidden)

  role_assignments = [
    {
      principal_id   = "..."
      definition     = "contributor"
      relative_scope = ""
    }
  ]
}
```

### v3.x `role_assignments` syntax

```terraform
module "lz_vending" {
  source  = "..."
  version = "..."

  # (other input variables hidden)

  role_assignments = {
    contrib_to_group = {
      principal_id   = "..."
      definition     = "contributor"
      relative_scope = ""
    }
  }
}
```

## Upgrading from v4.x to v5.x

###  Terraform version

We now require a minimum of Terraform version 1.8.

### Provider Versions

We now require a minimum of AzureRM version 4.0 and AzAPI version 2.2.

### Resource Groups

We have removed the boolean input variable to create the network watcher resource group.
Instead, use `var.resource_groups` to specify the resource groups to create.

We have used the `moved {}` block to move the resource in state.
If you previously deployed the network watcher resource group, please modify the value of `var.resource_groups` to include the existing resource group. The key ***must*** be `NetworkWatcherRG` You can use the following example:

```hcl
resource_groups = {
  NetworkWatcherRG = {
    name     = "NetworkWatcherRG"
    location = "your-location"
    tags     = {} # add tags here
  }
}
```

## Virtual WAN

When joining virtual networks to a Virtual WAN hub, the behaviour with routing intent has changed.
Previously the AzAPI provider allowed us to use `ignore_body_properties` to dynamically ignore parts of the resource body
With AzAPI v2 this is no longer possible, so we have to use the `lifecycle` block to ignore changes.
However, as ignore changes is not able to be user configurable, we have had to split the virtual hub connections into two separate resources.

In order to avoid destroying and re-creating the virtual hub connections, you will have to use the `moved {}` block to move the resource in state.
We are unable to do this for you because we do not know the specific instances of the resources that require moving.

```hcl
moved {
  from = module.YOUR_MODULE_ALIAS.module.virtualnetwork.azapi_resource.vhubconnection["instance_name"]
  to   = module.YOUR_MODULE_ALIAS.module.virtualnetwork.azapi_resource.vhubconnection_routing_intent["instance_name"]
}
```
