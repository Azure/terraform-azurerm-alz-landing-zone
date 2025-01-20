<!-- BEGIN_TF_DOCS -->
# Terraform landing zone vending module for Azure

[![Average time to resolve an issue](http://isitmaintained.com/badge/resolution/Azure/terraform-azurerm-lz-vending.svg)](http://isitmaintained.com/project/Azure/terraform-azurerm-lz-vending "Average time to resolve an issue")
[![Percentage of issues still open](http://isitmaintained.com/badge/open/Azure/terraform-azurerm-lz-vending.svg)](http://isitmaintained.com/project/Azure/terraform-azurerm-lz-vending "Percentage of issues still open")
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/Azure/terraform-azurerm-lz-vending/badge)](https://scorecard.dev/viewer/?uri=github.com/Azure/terraform-azurerm-lz-vending)

## Overview

The landing zone Terraform module is designed to accelerate deployment of individual landing zones within an Azure tenant.
We use the [AzureRM](https://registry.terraform.io/providers/hashicorp/azurerm/latest) and [AzAPI](https://registry.terraform.io/providers/azure/azapi/latest) providers to create the subscription and deploy the resources in a single `terraform apply` step.

The module is designed to be instantiated many times, once for each desired landing zone.

This is currently split logically into the following capabilities:

- Subscription creation and management group placement
- Networking - deploy multiple vnets with:
  - Hub & spoke connectivity (peering to a hub network)
  - vWAN connectivity
  - Mesh peering (peering between spokes)
- Role assignments
- Resource provider (and feature) registration
- Resource group creation
- User assigned managed identity creation
  - Federated credential configuration for GitHub Actions, Terraform Cloud, and other providers.

> When creating virtual network peerings, be aware of the [limit of peerings per virtual network](https://learn.microsoft.com/azure/azure-resource-manager/management/azure-subscription-service-limits?toc=%2Fazure%2Fvirtual-network%2Ftoc.json#azure-resource-manager-virtual-networking-limits).

We would like feedback on what's missing in the module.
Please raise an [issue](https://github.com/Azure/terraform-azurerm-lz-vending/issues) if you have any suggestions.

## Change log

Please see the [GitHub releases pages](https://github.com/Azure/terraform-azurerm-lz-vending/releases) for change log information.

## Notes

Please see the content in the [wiki](https://github.com/Azure/terraform-azurerm-lz-vending/wiki) for more detailed information.

## Example

The below example created a landing zone subscription with two virtual networks.
One virtual network is in the default location of the subscription, the other is in a different location.

The virtual networks are peered with the supplied hub network resource ids, they are also peered with each other using the mesh peering option.

```terraform
module "lz_vending" {
  source  = "Azure/lz-vending/azurerm"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  # Set the default location for resources
  location = "westeurope"

  # subscription variables
  subscription_alias_enabled = true
  subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
  subscription_display_name  = "my-subscription-display-name"
  subscription_alias_name    = "my-subscription-alias"
  subscription_workload      = "Production"

  network_watcher_resource_group_enabled = true

  # management group association variables
  subscription_management_group_association_enabled = true
  subscription_management_group_id                  = "Corp"

  # virtual network variables
  virtual_network_enabled = true
  virtual_networks = {
    one = {
      name                    = "my-vnet"
      address_space           = ["192.168.1.0/24"]
      hub_peering_enabled     = true
      hub_network_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-hub-network-rg/providers/Microsoft.Network/virtualNetworks/my-hub-network"
      mesh_peering_enabled    = true
    }
    two = {
      name                    = "my-vnet2"
      location                = "northeurope"
      address_space           = ["192.168.2.0/24"]
      hub_peering_enabled     = true
      hub_network_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-hub-network-rg/providers/Microsoft.Network/virtualNetworks/my-hub-network2"
      mesh_peering_enabled    = true
    }
  }

  umi_enabled             = true
  umi_name                = "umi"
  umi_resource_group_name = "rg-identity"
  umi_role_assignments = {
    myrg-contrib = {
      definition     = "Contributor"
      relative_scope = "/resourceGroups/MyRg"
    }
  }

  resource_group_creation_enabled = true
  resource_groups = {
    myrg = {
      name     = "MyRg"
      location = "westeurope"
    }
  }

  # role assignments
  role_assignment_enabled = true
  role_assignments = {
    # using role definition name, created at subscription scope
    contrib_user_sub = {
      principal_id   = "00000000-0000-0000-0000-000000000000"
      definition     = "Contributor"
      relative_scope = ""
    },
    # using a custom role definition
    custdef_sub_scope = {
      principal_id   = "11111111-1111-1111-1111-111111111111"
      definition     = "/providers/Microsoft.Management/MyMg/providers/Microsoft.Authorization/roleDefinitions/ffffffff-ffff-ffff-ffff-ffffffffffff"
      relative_scope = ""
    },
    # using relative scope (to the created or supplied subscription)
    rg_owner = {
      principal_id   = "00000000-0000-0000-0000-000000000000"
      definition     = "Owner"
      relative_scope = "/resourceGroups/MyRg"
    },
  }
}
```

## Documentation
<!-- markdownlint-disable MD033 -->

## Requirements

The following requirements are needed by this module:

- <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) (~> 1.8)

- <a name="requirement_azapi"></a> [azapi](#requirement\_azapi) (~> 2.2)

## Modules

The following Modules are called:

### <a name="module_budget"></a> [budget](#module\_budget)

Source: ./modules/budget

Version:

### <a name="module_resourcegroup"></a> [resourcegroup](#module\_resourcegroup)

Source: ./modules/resourcegroup

Version:

### <a name="module_resourceproviders"></a> [resourceproviders](#module\_resourceproviders)

Source: ./modules/resourceprovider

Version:

### <a name="module_roleassignment"></a> [roleassignment](#module\_roleassignment)

Source: ./modules/roleassignment

Version:

### <a name="module_roleassignment_umi"></a> [roleassignment\_umi](#module\_roleassignment\_umi)

Source: ./modules/roleassignment

Version:

### <a name="module_routetable"></a> [routetable](#module\_routetable)

Source: ./modules/routetable

Version:

### <a name="module_subscription"></a> [subscription](#module\_subscription)

Source: ./modules/subscription

Version:

### <a name="module_usermanagedidentity"></a> [usermanagedidentity](#module\_usermanagedidentity)

Source: ./modules/usermanagedidentity

Version:

### <a name="module_virtualnetwork"></a> [virtualnetwork](#module\_virtualnetwork)

Source: ./modules/virtualnetwork

Version:

<!-- markdownlint-disable MD013 -->
## Required Inputs

The following input variables are required:

### <a name="input_location"></a> [location](#input\_location)

Description: The default location of resources created by this module.  
Virtual networks will be created in this location unless overridden by the `location` attribute.

Type: `string`

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_budget_enabled"></a> [budget\_enabled](#input\_budget\_enabled)

Description: Whether to create budgets.  
If enabled, supply the list of budgets in `var.budgets`.

Type: `bool`

Default: `false`

### <a name="input_budgets"></a> [budgets](#input\_budgets)

Description: Map of budgets to create for the subscription.

- `amount` - The total amount of cost to track with the budget.
- `time_grain` - The time grain for the budget. Must be one of Annually, BillingAnnual, BillingMonth, BillingQuarter, Monthly, or Quarterly.
- `time_period_start` - The start date for the budget.
- `time_period_end` - The end date for the budget.
- `relative_scope` - (optional) Scope relative to the created subscription. Omit, or leave blank for subscription scope.
- `notifications` - (optional) The notifications to create for the budget.
  - `enabled` - Whether the notification is enabled.
  - `operator` - The operator for the notification. Must be one of GreaterThan or GreaterThanOrEqualTo.
  - `threshold` - The threshold for the notification. Must be between 0 and 1000.
  - `threshold_type` - The threshold type for the notification. Must be one of Actual or Forecasted.
  - `contact_emails` - The contact emails for the notification.
  - `contact_roles` - The contact roles for the notification.
  - `contact_groups` - The contact groups for the notification.
  - `locale` - The locale for the notification. Must be in the format xx-xx.

time\_period\_start and time\_period\_end must be UTC in RFC3339 format, e.g. 2018-05-13T07:44:12Z.

Example value:

```terraform
subscription_budgets = {
  budget1 = {
    amount            = 150
    time_grain        = "Monthly"
    time_period_start = "2024-01-01T00:00:00Z"
    time_period_end   = "2027-12-31T23:59:59Z"
    notifications = {
      eightypercent = {
        enabled        = true
        operator       = "GreaterThan"
        threshold      = 80
        threshold_type = "Actual"
        contact_emails = ["john@contoso.com"]
      }
      budgetexceeded = {
        enabled        = true
        operator       = "GreaterThan"
        threshold      = 120
        threshold_type = "Forecasted"
        contact_roles  = ["Owner"]
      }
    }
  }
}
```

Type:

```hcl
map(object({
    amount            = number
    time_grain        = string
    time_period_start = string
    time_period_end   = string
    relative_scope    = optional(string, "")
    notifications = optional(map(object({
      enabled        = bool
      operator       = string
      threshold      = number
      threshold_type = optional(string, "Actual")
      contact_emails = optional(list(string), [])
      contact_roles  = optional(list(string), [])
      contact_groups = optional(list(string), [])
      locale         = optional(string, "en-us")
    })), {})
  }))
```

Default: `{}`

### <a name="input_disable_telemetry"></a> [disable\_telemetry](#input\_disable\_telemetry)

Description: To disable tracking, we have included this variable with a simple boolean flag.  
The default value is `false` which does not disable the telemetry.  
If you would like to disable this tracking, then simply set this value to true and this module will not create the telemetry tracking resources and therefore telemetry tracking will be disabled.

For more information, see the [wiki](https://aka.ms/lz-vending/tf/telemetry)

E.g.

```terraform
module "lz_vending" {
  source  = "Azure/lz-vending/azurerm"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  # ... other module variables

  disable_telemetry = true
}
```

Type: `bool`

Default: `false`

### <a name="input_resource_group_creation_enabled"></a> [resource\_group\_creation\_enabled](#input\_resource\_group\_creation\_enabled)

Description: Whether to create additional resource groups in the target subscription. Requires `var.resource_groups`.

Type: `bool`

Default: `false`

### <a name="input_resource_groups"></a> [resource\_groups](#input\_resource\_groups)

Description: A map of the resource groups to create. The value is an object with the following attributes:

- `name` - the name of the resource group
- `location` - the location of the resource group
- `tags` - (optional) a map of type string

We recommend that you include an entry to create the NetworkWatcherRG resource group so that this is managed by Terraform.

Type:

```hcl
map(object({
    name     = string
    location = string
    tags     = optional(map(string), {})
  }))
```

Default: `{}`

### <a name="input_role_assignment_enabled"></a> [role\_assignment\_enabled](#input\_role\_assignment\_enabled)

Description: Whether to create role assignments.  
If enabled, supply the list of role assignments in `var.role_assignments`.

Type: `bool`

Default: `false`

### <a name="input_role_assignments"></a> [role\_assignments](#input\_role\_assignments)

Description: Supply a map of objects containing the details of the role assignments to create.

Object fields:

- `principal_id`: The directory/object id of the principal to assign the role to.
- `definition`: The role definition to assign. Either use the name or the role definition resource id.
- `relative_scope`: (optional) Scope relative to the created subscription. Omit, or leave blank for subscription scope.
- `condition`: (optional) A condition to apply to the role assignment. See [Conditions Custom Security Attributes](https://learn.microsoft.com/azure/role-based-access-control/conditions-custom-security-attributes) for more details.
- `condition_version`: (optional) The version of the condition syntax. See [Conditions Custom Security Attributes](https://learn.microsoft.com/azure/role-based-access-control/conditions-custom-security-attributes) for more details.

E.g.

```terraform
role_assignments = {
  # Example using role definition name:
  contributor_user = {
    principal_id      = "00000000-0000-0000-0000-000000000000",
    definition        = "Contributor",
    relative_scope    = "",
    condition         = "(!(ActionMatches{'Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read'} AND NOT SubOperationMatches{'Blob.List'})",
    condition_version = "2.0",
  },
  # Example using role definition id and RG scope:
  myrg_custom_role = {
    principal_id   = "11111111-1111-1111-1111-111111111111",
    definition     = "/providers/Microsoft.Management/managementGroups/mymg/providers/Microsoft.Authorization/roleDefinitions/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    relative_scope = "/resourceGroups/MyRg",
  }
}
```

Type:

```hcl
map(object({
    principal_id      = string,
    definition        = string,
    relative_scope    = optional(string, ""),
    condition         = optional(string, ""),
    condition_version = optional(string, ""),
  }))
```

Default: `{}`

### <a name="input_route_table_enabled"></a> [route\_table\_enabled](#input\_route\_table\_enabled)

Description: Whether to create route tables and routes in the target subscription. Requires `var.route_tables`.

Type: `bool`

Default: `false`

### <a name="input_route_tables"></a> [route\_tables](#input\_route\_tables)

Description: A map defining route tables and their associated routes to be created:

- `name` (required): The name of the route table.
- `location` (required): The location of the resource group.
- `resource_group_name` (required): The name of the resource group.
- `bgp_route_propagation_enabled` (optional): Boolean that controls whether routes learned by BGP are propagated to the route table. Default is `true`.
- `tags` (optional): A map of key-value pairs for tags associated with the route table.
- `routes` (optional): A map defining routes for the route table. Each route object has the following properties:
  - `name` (required): The name of the route.
  - `address_prefix` (required): The address prefix for the route.
  - `next_hop_type` (required): The type of next hop for the route.
  - `next_hop_in_ip_address` (required): The next hop IP address for the route.

Type:

```hcl
map(object({
    name                          = string
    location                      = string
    resource_group_name           = string
    bgp_route_propagation_enabled = optional(bool, true)
    tags                          = optional(map(string))

    routes = optional(map(object({
      name                   = string
      address_prefix         = string
      next_hop_type          = string
      next_hop_in_ip_address = string
    })))
  }))
```

Default: `{}`

### <a name="input_subscription_alias_enabled"></a> [subscription\_alias\_enabled](#input\_subscription\_alias\_enabled)

Description: Whether to create a new subscription using the subscription alias resource.

If enabled, the following must also be supplied:

- `subscription_alias_name`
- `subscription_display_name`
- `subscription_billing_scope`
- `subscription_workload`

Optionally, supply the following to enable the placement of the subscription into a management group:

- `subscription_management_group_id`
- `subscription_management_group_association_enabled`

If disabled, supply the `subscription_id` variable to use an existing subscription instead.

> **Note**: When the subscription is destroyed, this module will try to remove the NetworkWatcherRG resource group using `az cli`.
> This requires the `az cli` tool be installed and authenticated.
> If the command fails for any reason, the provider will attempt to cancel the subscription anyway.

Type: `bool`

Default: `false`

### <a name="input_subscription_alias_name"></a> [subscription\_alias\_name](#input\_subscription\_alias\_name)

Description: The name of the subscription alias.

The string must be comprised of a-z, A-Z, 0-9, - and \_.  
The maximum length is 63 characters.

You may also supply an empty string if you do not want to create a new subscription alias.  
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.

Type: `string`

Default: `""`

### <a name="input_subscription_billing_scope"></a> [subscription\_billing\_scope](#input\_subscription\_billing\_scope)

Description: The billing scope for the new subscription alias.

A valid billing scope starts with `/providers/Microsoft.Billing/billingAccounts/` and is case sensitive.

E.g.

- For CustomerLed and FieldLed, e.g. MCA - `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/billingProfiles/{billingProfileName}/invoiceSections/{invoiceSectionName}`
- For PartnerLed, e.g. MPA - `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/customers/{customerName}`
- For Legacy EA - `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/enrollmentAccounts/{enrollmentAccountName}`

You may also supply an empty string if you do not want to create a new subscription alias.  
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.

Type: `string`

Default: `""`

### <a name="input_subscription_display_name"></a> [subscription\_display\_name](#input\_subscription\_display\_name)

Description: The display name of the subscription alias.

The string must be comprised of a-z, A-Z, 0-9, -, \_ and space.  
The maximum length is 63 characters.

You may also supply an empty string if you do not want to create a new subscription alias.  
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.

Type: `string`

Default: `""`

### <a name="input_subscription_id"></a> [subscription\_id](#input\_subscription\_id)

Description: An existing subscription id.

Use this when you do not want the module to create a new subscription.  
But do want to manage the management group membership.

A GUID should be supplied in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.  
All letters must be lowercase.

When using this, `subscription_management_group_association_enabled` should be enabled,  
and `subscription_management_group_id` should be supplied.

You may also supply an empty string if you want to create a new subscription alias.  
In this scenario, `subscription_alias_enabled` should be set to `true` and the following other variables must be supplied:

- `subscription_alias_name`
- `subscription_alias_display_name`
- `subscription_alias_billing_scope`
- `subscription_alias_workload`

Type: `string`

Default: `""`

### <a name="input_subscription_management_group_association_enabled"></a> [subscription\_management\_group\_association\_enabled](#input\_subscription\_management\_group\_association\_enabled)

Description: Whether to create the `azurerm_management_group_subscription_association` resource.

If enabled, the `subscription_management_group_id` must also be supplied.

Type: `bool`

Default: `false`

### <a name="input_subscription_management_group_id"></a> [subscription\_management\_group\_id](#input\_subscription\_management\_group\_id)

Description:   The destination management group ID for the new subscription.

**Note:** Do not supply the display name.  
The management group ID forms part of the Azure resource ID. E.g.,
`/providers/Microsoft.Management/managementGroups/{managementGroupId}`.

Type: `string`

Default: `""`

### <a name="input_subscription_register_resource_providers_and_features"></a> [subscription\_register\_resource\_providers\_and\_features](#input\_subscription\_register\_resource\_providers\_and\_features)

Description: The map of resource providers to register.  
The map keys are the resource provider namespace, e.g. `Microsoft.Compute`.  
The map values are a list of provider features to enable.  
Leave the value empty to not register any resource provider features.

The default values are taken from [Hashicorp's AzureRM provider](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/internal/resourceproviders/required.go).

Type: `map(set(string))`

Default:

```json
{
  "Microsoft.AVS": [],
  "Microsoft.ApiManagement": [],
  "Microsoft.AppPlatform": [],
  "Microsoft.Authorization": [],
  "Microsoft.Automation": [],
  "Microsoft.Blueprint": [],
  "Microsoft.BotService": [],
  "Microsoft.Cache": [],
  "Microsoft.Cdn": [],
  "Microsoft.CognitiveServices": [],
  "Microsoft.Compute": [],
  "Microsoft.ContainerInstance": [],
  "Microsoft.ContainerRegistry": [],
  "Microsoft.ContainerService": [],
  "Microsoft.CostManagement": [],
  "Microsoft.CustomProviders": [],
  "Microsoft.DBforMariaDB": [],
  "Microsoft.DBforMySQL": [],
  "Microsoft.DBforPostgreSQL": [],
  "Microsoft.DataLakeAnalytics": [],
  "Microsoft.DataLakeStore": [],
  "Microsoft.DataMigration": [],
  "Microsoft.DataProtection": [],
  "Microsoft.Databricks": [],
  "Microsoft.DesktopVirtualization": [],
  "Microsoft.DevTestLab": [],
  "Microsoft.Devices": [],
  "Microsoft.DocumentDB": [],
  "Microsoft.EventGrid": [],
  "Microsoft.EventHub": [],
  "Microsoft.GuestConfiguration": [],
  "Microsoft.HDInsight": [],
  "Microsoft.HealthcareApis": [],
  "Microsoft.KeyVault": [],
  "Microsoft.Kusto": [],
  "Microsoft.Logic": [],
  "Microsoft.MachineLearningServices": [],
  "Microsoft.Maintenance": [],
  "Microsoft.ManagedIdentity": [],
  "Microsoft.ManagedServices": [],
  "Microsoft.Management": [],
  "Microsoft.Maps": [],
  "Microsoft.MarketplaceOrdering": [],
  "Microsoft.Media": [],
  "Microsoft.MixedReality": [],
  "Microsoft.Network": [],
  "Microsoft.NotificationHubs": [],
  "Microsoft.OperationalInsights": [],
  "Microsoft.OperationsManagement": [],
  "Microsoft.PolicyInsights": [],
  "Microsoft.PowerBIDedicated": [],
  "Microsoft.RecoveryServices": [],
  "Microsoft.Relay": [],
  "Microsoft.Resources": [],
  "Microsoft.Search": [],
  "Microsoft.Security": [],
  "Microsoft.SecurityInsights": [],
  "Microsoft.ServiceBus": [],
  "Microsoft.ServiceFabric": [],
  "Microsoft.Sql": [],
  "Microsoft.Storage": [],
  "Microsoft.StreamAnalytics": [],
  "Microsoft.Web": [],
  "microsoft.insights": []
}
```

### <a name="input_subscription_register_resource_providers_enabled"></a> [subscription\_register\_resource\_providers\_enabled](#input\_subscription\_register\_resource\_providers\_enabled)

Description: Whether to register resource providers for the subscription.  
Use `var.subscription_register_resource_providers_and_features` to customize registration.

Type: `bool`

Default: `false`

### <a name="input_subscription_tags"></a> [subscription\_tags](#input\_subscription\_tags)

Description: A map of tags to assign to the newly created subscription.  
Only valid when `subsciption_alias_enabled` is set to `true`.

Example value:

```terraform
subscription_tags = {
  mytag  = "myvalue"
  mytag2 = "myvalue2"
}
```

Type: `map(string)`

Default: `{}`

### <a name="input_subscription_update_existing"></a> [subscription\_update\_existing](#input\_subscription\_update\_existing)

Description: Whether to update an existing subscription with the supplied tags and display name.  
If enabled, the following must also be supplied:
- `subscription_id`

Type: `bool`

Default: `false`

### <a name="input_subscription_use_azapi"></a> [subscription\_use\_azapi](#input\_subscription\_use\_azapi)

Description: Whether to create a new subscription using the azapi provider. This may be required if the principal running  
terraform does not have the required permissions to create a subscription under the default management group.  
If enabled, the following must also be supplied:
- `subscription_alias_name`
- `subscription_display_name`
- `subscription_billing_scope`
- `subscription_workload`  
Optionally, supply the following to enable the placement of the subscription into a management group:
- `subscription_management_group_id`
- `subscription_management_group_association_enabled`  
If disabled, supply the `subscription_id` variable to use an existing subscription instead.
> **Note**: When the subscription is destroyed, this module will try to remove the NetworkWatcherRG resource group using `az cli`.
> This requires the `az cli` tool be installed and authenticated.
> If the command fails for any reason, the provider will attempt to cancel the subscription anyway.

Type: `bool`

Default: `false`

### <a name="input_subscription_workload"></a> [subscription\_workload](#input\_subscription\_workload)

Description: The billing scope for the new subscription alias.

The workload type can be either `Production` or `DevTest` and is case sensitive.

You may also supply an empty string if you do not want to create a new subscription alias.  
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.

Type: `string`

Default: `""`

### <a name="input_umi_enabled"></a> [umi\_enabled](#input\_umi\_enabled)

Description: Whether to enable the creation of a user-assigned managed identity.

Requires `umi_name` and `umi_resosurce_group_name` to be non-empty.

Type: `bool`

Default: `false`

### <a name="input_umi_federated_credentials_advanced"></a> [umi\_federated\_credentials\_advanced](#input\_umi\_federated\_credentials\_advanced)

Description: Configure federated identity credentials, using OpenID Connect, for use scenarios outside GitHub Actions and Terraform Cloud.

The may key is arbitrary and only used for the `for_each` in the resource declaration.

The map value is an object with the following attributes:

- `name`: The name of the federated credential resource, the last segment of the Azure resource id.
- `subject_identifier`: The subject of the token.
- `issuer_url`: The URL of the token issuer, should begin with `https://`
- `audience`: (optional) The token audience, defaults to `api://AzureADTokenExchange`.

Type:

```hcl
map(object({
    name               = string
    subject_identifier = string
    issuer_url         = string
    audiences          = optional(set(string), ["api://AzureADTokenExchange"])
  }))
```

Default: `{}`

### <a name="input_umi_federated_credentials_github"></a> [umi\_federated\_credentials\_github](#input\_umi\_federated\_credentials\_github)

Description: Configure federated identity credentials, using OpenID Connect, for use in GitHub actions.

The may key is arbitrary and only used for the `for_each` in the resource declaration.

The map value is an object with the following attributes:

- `name` - the name of the federated credential resource, the last segment of the Azure resource id.
- `organization` - the name of the GitHub organization, e.g. `Azure` in `https://github.com/Azure/terraform-azurerm-lz-vending`.
- `repository` - the name of the GitHub respository, e.g. `terraform-azurerm-lz-vending` in `https://github.com/Azure/terraform-azurerm-lz-vending`.
- `entity` - one of 'environment', 'pull\_request', 'tag', or 'branch'
- `value` - identifies the `entity` type, e.g. `main` when using entity is `branch`. Should be blank when `entity` is `pull_request`.

Type:

```hcl
map(object({
    name         = optional(string, "")
    organization = string
    repository   = string
    entity       = string
    value        = optional(string, "")
  }))
```

Default: `{}`

### <a name="input_umi_federated_credentials_terraform_cloud"></a> [umi\_federated\_credentials\_terraform\_cloud](#input\_umi\_federated\_credentials\_terraform\_cloud)

Description: Configure federated identity credentials, using OpenID Connect, for use in Terraform Cloud.

The may key is arbitrary and only used for the `for_each` in the resource declaration.

The map value is an object with the following attributes:

- `name` - the name of the federated credential resource, the last segment of the Azure resource id.
- `organization` - the name of the Terraform Cloud organization.
- `project` - the name of the Terraform Cloud project.
- `workspace` - the name of the Terraform Cloud workspace.
- `run_phase` - one of `plan`, or `apply`.

Type:

```hcl
map(object({
    name         = optional(string, "")
    organization = string
    project      = string
    workspace    = string
    run_phase    = string
  }))
```

Default: `{}`

### <a name="input_umi_name"></a> [umi\_name](#input\_umi\_name)

Description: The name of the user-assigned managed identity

Type: `string`

Default: `""`

### <a name="input_umi_resource_group_creation_enabled"></a> [umi\_resource\_group\_creation\_enabled](#input\_umi\_resource\_group\_creation\_enabled)

Description: Whether to create the supplied resource group for the user-assigned managed identity

Type: `bool`

Default: `true`

### <a name="input_umi_resource_group_lock_enabled"></a> [umi\_resource\_group\_lock\_enabled](#input\_umi\_resource\_group\_lock\_enabled)

Description: Whether to enable resource group lock for the user-assigned managed identity resource group

Type: `bool`

Default: `true`

### <a name="input_umi_resource_group_lock_name"></a> [umi\_resource\_group\_lock\_name](#input\_umi\_resource\_group\_lock\_name)

Description: The name of the resource group lock for the user-assigned managed identity resource group, if blank will be set to `lock-<resource_group_name>`

Type: `string`

Default: `""`

### <a name="input_umi_resource_group_name"></a> [umi\_resource\_group\_name](#input\_umi\_resource\_group\_name)

Description: The name of the resource group in which to create the user-assigned managed identity

Type: `string`

Default: `""`

### <a name="input_umi_resource_group_tags"></a> [umi\_resource\_group\_tags](#input\_umi\_resource\_group\_tags)

Description: The tags to apply to the user-assigned managed identity resource group, if we create it.

Type: `map(string)`

Default: `{}`

### <a name="input_umi_role_assignments"></a> [umi\_role\_assignments](#input\_umi\_role\_assignments)

Description: Supply a map of objects containing the details of the role assignments to create for the user-assigned managed identity.  
This will be merged with the other role assignments specified in `var.role_assignments`.

The role assignments can be used resource groups created by the `var.resource_groups` map.

Requires both `var.umi_enabled` and `var.role_assignment_enabled` to be `true`.

Object fields:

- `definition`: The role definition to assign. Either use the name or the role definition resource id.
- `relative_scope`: Scope relative to the created subscription. Leave blank for subscription scope.

Type:

```hcl
map(object({
    definition        = string
    relative_scope    = optional(string, "")
    condition         = optional(string, "")
    condition_version = optional(string, "")
  }))
```

Default: `{}`

### <a name="input_umi_tags"></a> [umi\_tags](#input\_umi\_tags)

Description: The tags to apply to the user-assigned managed identity

Type: `map(string)`

Default: `{}`

### <a name="input_virtual_network_enabled"></a> [virtual\_network\_enabled](#input\_virtual\_network\_enabled)

Description: Enables and disables the virtual network submodule.

Type: `bool`

Default: `false`

### <a name="input_virtual_networks"></a> [virtual\_networks](#input\_virtual\_networks)

Description: A map of the virtual networks to create. The map key must be known at the plan stage, e.g. must not be calculated and known only after apply.

### Required fields

- `name`: The name of the virtual network. [required]
- `address_space`: The address space of the virtual network as a list of strings in CIDR format, e.g. `["192.168.0.0/24", "10.0.0.0/24"]`. [required]
- `resource_group_name`: The name of the resource group to create the virtual network in. [required]

### DNS servers

- `dns_servers`: A list of DNS servers to use for the virtual network, e.g. `["192.168.0.1", "10.0.0.1"]`. If empty will use the Azure default DNS. [optional - default empty list]

### DDOS protection plan

- `ddos_protection_enabled`: Whether to enable ddos protection. [optional]
- `ddos_protection_plan_id`: The resource ID of the protection plan to attach the vnet. [optional - but required if ddos\_protection\_enabled is `true`]

### Location

- `location`: The location of the virtual network (and resource group if creation is enabled). [optional, will use `var.location` if not specified or empty string]

> Note at least one of `location` or `var.location` must be specified.
> If both are empty then the module will fail.

### Hub network peering values

The following values configure bi-directional hub & spoke peering for the given virtual network.

- `hub_peering_enabled`: Whether to enable hub peering. [optional]
- `hub_peering_direction`: The direction of the peering. [optional - allowed values are: `tohub`, `fromhub` or `both` - default `both`]
- `hub_network_resource_id`: The resource ID of the hub network to peer with. [optional - but required if hub\_peering\_enabled is `true`]
- `hub_peering_name_tohub`: The name of the peering to the hub network. [optional - leave empty to use calculated name]
- `hub_peering_name_fromhub`: The name of the peering from the hub network. [optional - leave empty to use calculated name]
- `hub_peering_use_remote_gateways`: Whether to use remote gateways for the hub peering. [optional - default true]

### Mesh peering values

Mesh peering is the capability to create a bi-directional peerings between all supplied virtual networks in `var.virtual_networks`.  
Peerings will only be created between virtual networks with the `mesh_peering_enabled` value set to `true`.

- `mesh_peering_enabled`: Whether to enable mesh peering for this virtual network. Must be enabled on more than one virtual network for any peerings to be created. [optional]
- `mesh_peering_allow_forwarded_traffic`: Whether to allow forwarded traffic for the mesh peering. [optional - default false]

### Resource group values

The default is that a resource group will be created for each resource\_group\_name specified in the `var.virtual_networks` map.  
It is possible to use a pre-existing resource group by setting `resource_group_creation_enabled` to `false`.  
We recommend using resource groups aligned to the region of the virtual network,  
however if you want multiple virtual networks in more than one location to share a resource group,  
only one of the virtual networks should have `resource_group_creation_enabled` set to `true`.

- `resource_group_creation_enabled`: Whether to create a resource group for the virtual network. [optional - default `true`]
- `resource_group_lock_enabled`: Whether to create a `CanNotDelete` resource lock on the resource group. [optional - default `true`]
- `resource_group_lock_name`: The name of the resource lock. [optional - leave empty to use calculated name]
- `resource_group_tags`: A map of tags to apply to the resource group, e.g. `{ mytag = "myvalue", mytag2 = "myvalue2" }`. [optional - default empty]

### Virtual WAN values

- `vwan_associated_routetable_resource_id`: The resource ID of the route table to associate with the virtual network. [optional - leave empty to use `defaultRouteTable` on hub]
- `vwan_connection_enabled`: Whether to create a connection to a Virtual WAN. [optional - default false]
- `vwan_connection_name`: The name of the connection to the Virtual WAN. [optional - leave empty to use calculated name]
- `vwan_hub_resource_id`: The resource ID of the hub to connect to. [optional - but required if vwan\_connection\_enabled is `true`]
- `vwan_propagated_routetables_labels`: A list of labels of route tables to propagate to the virtual network. [optional - leave empty to use `["default"]`]
- `vwan_propagated_routetables_resource_ids`: A list of resource IDs of route tables to propagate to the virtual network. [optional - leave empty to use `defaultRouteTable` on hub]
- `vwan_security_configuration`: A map of security configuration values for VWAN hub connection - see below. [optional - default empty]
  - `secure_internet_traffic`: Whether to forward internet-bound traffic to the destination specified in the routing policy. [optional - default `false`]
  - `secure_private_traffic`: Whether to all internal traffic to the destination specified in the routing policy. Not compatible with `routing_intent_enabled`. [optional - default `false`]
  - `routing_intent_enabled`: Enable to use with a Virtual WAN hub with routing intent enabled. Routing intent on hub is configured outside this module. [optional - default `false`]

### Tags

- `tags`: A map of tags to apply to the virtual network. [optional - default empty]

Type:

```hcl
map(object({
    name                = string
    address_space       = list(string)
    resource_group_name = string

    location = optional(string, "")

    dns_servers = optional(list(string), [])

    ddos_protection_enabled = optional(bool, false)
    ddos_protection_plan_id = optional(string, "")

    hub_network_resource_id         = optional(string, "")
    hub_peering_enabled             = optional(bool, false)
    hub_peering_direction           = optional(string, "both")
    hub_peering_name_tohub          = optional(string, "")
    hub_peering_name_fromhub        = optional(string, "")
    hub_peering_use_remote_gateways = optional(bool, true)

    mesh_peering_enabled                 = optional(bool, false)
    mesh_peering_allow_forwarded_traffic = optional(bool, false)

    resource_group_creation_enabled = optional(bool, true)
    resource_group_lock_enabled     = optional(bool, true)
    resource_group_lock_name        = optional(string, "")
    resource_group_tags             = optional(map(string), {})

    vwan_associated_routetable_resource_id   = optional(string, "")
    vwan_connection_enabled                  = optional(bool, false)
    vwan_connection_name                     = optional(string, "")
    vwan_hub_resource_id                     = optional(string, "")
    vwan_propagated_routetables_labels       = optional(list(string), [])
    vwan_propagated_routetables_resource_ids = optional(list(string), [])
    vwan_security_configuration = optional(object({
      secure_internet_traffic = optional(bool, false)
      secure_private_traffic  = optional(bool, false)
      routing_intent_enabled  = optional(bool, false)
    }), {})

    tags = optional(map(string), {})
  }))
```

Default: `{}`

### <a name="input_wait_for_subscription_before_subscription_operations"></a> [wait\_for\_subscription\_before\_subscription\_operations](#input\_wait\_for\_subscription\_before\_subscription\_operations)

Description: The duration to wait after vending a subscription before performing subscription operations.

Type:

```hcl
object({
    create  = optional(string, "30s")
    destroy = optional(string, "0s")
  })
```

Default: `{}`

## Resources

The following resources are used by this module:

- [azapi_resource.telemetry_root](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)

## Outputs

The following outputs are exported:

### <a name="output_budget_resource_id"></a> [budget\_resource\_id](#output\_budget\_resource\_id)

Description: The created budget resource IDs, expressed as a map.

### <a name="output_management_group_subscription_association_id"></a> [management\_group\_subscription\_association\_id](#output\_management\_group\_subscription\_association\_id)

Description: The management\_group\_subscription\_association\_id output is the ID of the management group subscription association.  
Value will be null if `var.subscription_management_group_association_enabled` is false.

### <a name="output_resource_group_resource_ids"></a> [resource\_group\_resource\_ids](#output\_resource\_group\_resource\_ids)

Description: The created resource group IDs, expressed as a map.

### <a name="output_route_table_resource_ids"></a> [route\_table\_resource\_ids](#output\_route\_table\_resource\_ids)

Description: The created route table resource IDs, expressed as a map.

### <a name="output_subscription_id"></a> [subscription\_id](#output\_subscription\_id)

Description: The subscription\_id is the Azure subscription id that resources have been deployed into.

### <a name="output_subscription_resource_id"></a> [subscription\_resource\_id](#output\_subscription\_resource\_id)

Description: The subscription\_resource\_id is the Azure subscription resource id that resources have been deployed into

### <a name="output_umi_client_id"></a> [umi\_client\_id](#output\_umi\_client\_id)

Description: The client id of the user managed identity.  
Value will be null if `var.umi_enabled` is false.

### <a name="output_umi_id"></a> [umi\_id](#output\_umi\_id)

Description: The Azure resource id of the user managed identity.  
Value will be null if `var.umi_enabled` is false.

### <a name="output_umi_principal_id"></a> [umi\_principal\_id](#output\_umi\_principal\_id)

Description: The principal id of the user managed identity, sometimes known as the object id.  
Value will be null if `var.umi_enabled` is false.

### <a name="output_umi_tenant_id"></a> [umi\_tenant\_id](#output\_umi\_tenant\_id)

Description: The tenant id of the user managed identity.  
Value will be null if `var.umi_enabled` is false.

### <a name="output_virtual_network_resource_group_ids"></a> [virtual\_network\_resource\_group\_ids](#output\_virtual\_network\_resource\_group\_ids)

Description: A map of resource group ids, keyed by the var.virtual\_networks input map. Only populated if the virtualnetwork submodule is enabled.

### <a name="output_virtual_network_resource_ids"></a> [virtual\_network\_resource\_ids](#output\_virtual\_network\_resource\_ids)

Description: A map of virtual network resource ids, keyed by the var.virtual\_networks input map. Only populated if the virtualnetwork submodule is enabled.

<!-- markdownlint-enable -->
<!-- markdownlint-disable MD041 -->
## Telemetry
<!-- markdownlint-enable -->

When you deploy one or more modules using the landing zone vending module, Microsoft can identify the installation of said module with the deployed Azure resources.
Microsoft can correlate these resources used to support the software.
Microsoft collects this information to provide the best experiences with their products and to operate their business.
The telemetry is collected through customer usage attribution.
The data is collected and governed by Microsoft's privacy policies.

If you don't wish to send usage data to Microsoft, details on how to turn it off can be found [here](https://github.com/Azure/terraform-azurerm-lz-vending/wiki/Telemetry).

## Contributing

This project welcomes contributions and suggestions.
Most contributions require you to agree to a Contributor License Agreement (CLA)
declaring that you have the right to, and actually do, grant us the rights to use your contribution.
For details, visit [https://cla.opensource.microsoft.com](https://cla.opensource.microsoft.com).

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment).
Simply follow the instructions provided by the bot.
You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

## Developing the Module

See [DEVELOPER.md](https://github.com/Azure/terraform-azurerm-lz-vending/blob/main/DEVELOPER.md).

## Trademarks

This project may contain trademarks or logos for projects, products, or services.
Authorized use of Microsoft trademarks or logos is subject to and must follow
[Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/legal/intellectualproperty/trademarks/usage/general).
Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship.
Any use of third-party trademarks or logos are subject to those third-party's policies.
<!-- END_TF_DOCS -->