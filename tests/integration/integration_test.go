package integration

// Integration tests apply to the root module and test the deployment of resources
// in specific scenarios.

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../"
)

// TestIntegrationHubAndSpoke tests the resource plan when creating a new subscription,
// with a new virtual network with peerings to a supplied hub network.
func TestIntegrationHubAndSpoke(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet"
	v["subscription_alias_enabled"] = true
	v["virtual_network_enabled"] = true
	v["virtual_network_peering_enabled"] = true
	v["virtual_network_resource_lock_enabled"] = true
	terraformOptions.Vars = v

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	resources := []string{
		"azapi_resource.telemetry_root[0]",
		"module.subscription[0].azurerm_subscription.this[0]",
		"module.virtualnetwork[0].azapi_resource.peering[\"inbound\"]",
		"module.virtualnetwork[0].azapi_resource.peering[\"outbound\"]",
		"module.virtualnetwork[0].azapi_resource.rg_lock[0]",
		"module.virtualnetwork[0].azapi_resource.rg",
		"module.virtualnetwork[0].azapi_resource.vnet",
		"module.virtualnetwork[0].azapi_update_resource.vnet",
	}
	assert.Lenf(t, plan.ResourcePlannedValuesMap, len(resources), "expected %d resources to be created, but got %d", len(resources), len(plan.ResourcePlannedValuesMap))
	for _, v := range resources {
		terraform.AssertPlannedValuesMapKeyExists(t, plan, v)
	}
	// check bit field is correct
	telem := plan.ResourcePlannedValuesMap["azapi_resource.telemetry_root[0]"]
	require.Contains(t, telem.AttributeValues, "name")
	telemName := telem.AttributeValues["name"].(string)
	telemBf := strings.Split(telemName, "_")[2]
	expectBf := "00000b05"
	assert.Equalf(t, expectBf, telemBf, "expected bit field to be %s, but got %s", expectBf, telemBf)
}

// TestIntegrationVwan tests the resource plan when creating a new subscription,
// with a new virtual network and vwan connection to a supplied vhub.
// RG resource lock is disabled
func TestIntegrationVwan(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualHubs/testhub"
	v["subscription_alias_enabled"] = true
	v["virtual_network_enabled"] = true
	v["virtual_network_vwan_connection_enabled"] = true
	terraformOptions.Vars = v

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	resources := []string{
		"azapi_resource.telemetry_root[0]",
		"module.subscription[0].azurerm_subscription.this[0]",
		"module.virtualnetwork[0].azapi_resource.vhubconnection[\"this\"]",
		"module.virtualnetwork[0].azapi_resource.rg",
		"module.virtualnetwork[0].azapi_resource.vnet",
		"module.virtualnetwork[0].azapi_update_resource.vnet",
	}
	assert.Lenf(t, plan.ResourcePlannedValuesMap, len(resources), "expected %d resources to be created, but got %d", len(resources), len(plan.ResourcePlannedValuesMap))
	for _, v := range resources {
		terraform.AssertPlannedValuesMapKeyExists(t, plan, v)
	}
	// check bit field is correct
	telem := plan.ResourcePlannedValuesMap["azapi_resource.telemetry_root[0]"]
	require.Contains(t, telem.AttributeValues, "name")
	telemName := telem.AttributeValues["name"].(string)
	telemBf := strings.Split(telemName, "_")[2]
	expectBf := "00000505"
	assert.Equalf(t, expectBf, telemBf, "expected bit field to be %s, but got %s", expectBf, telemBf)
}

// TestIntegrationSubscriptionAndRoleAssignmentOnly tests the resource plan when creating a new subscription,
// with a role assignments, but no networking.
// This tests that the depends_on property of the roleassignments module is working
// when a dependent resource is disabled through the use of count.
func TestIntegrationSubscriptionAndRoleAssignmentOnly(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["subscription_alias_enabled"] = true
	v["virtual_network_enabled"] = false
	v["role_assignment_enabled"] = true
	v["role_assignments"] = []interface{}{
		map[string]interface{}{
			"principal_id":   "00000000-0000-0000-0000-000000000000",
			"definition":     "Owner",
			"relative_scope": "",
		},
	}
	terraformOptions.Vars = v

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	resources := []string{
		"azapi_resource.telemetry_root[0]",
		"module.subscription[0].azurerm_subscription.this[0]",
		"module.roleassignment[\"7f69efa3-575a-5f8b-a989-c3978b92b58a\"].azurerm_role_assignment.this",
	}
	assert.Lenf(t, plan.ResourcePlannedValuesMap, len(resources), "expected %d resources to be created, but got %d", len(resources), len(plan.ResourcePlannedValuesMap))
	for _, v := range resources {
		terraform.AssertPlannedValuesMapKeyExists(t, plan, v)
	}
	// check bit field is correct
	telem := plan.ResourcePlannedValuesMap["azapi_resource.telemetry_root[0]"]
	require.Contains(t, telem.AttributeValues, "name")
	telemName := telem.AttributeValues["name"].(string)
	telemBf := strings.Split(telemName, "_")[2]
	expectBf := "00010005"
	assert.Equalf(t, expectBf, telemBf, "expected bit field to be %s, but got %s", expectBf, telemBf)
}

// TestIntegrationHubAndSpokeExistingSubscription tests the resource plan when supplying an existing subscription,
// with a new virtual network with peerings to a supplied hub network.
func TestIntegrationHubAndSpokeExistingSubscription(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet"
	v["subscription_alias_enabled"] = false
	v["subscription_id"] = "00000000-0000-0000-0000-000000000000"
	v["virtual_network_enabled"] = true
	v["virtual_network_peering_enabled"] = true
	delete(v, "subscription_tags")
	terraformOptions.Vars = v

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	resources := []string{
		"azapi_resource.telemetry_root[0]",
		"module.virtualnetwork[0].azapi_resource.peering[\"inbound\"]",
		"module.virtualnetwork[0].azapi_resource.peering[\"outbound\"]",
		"module.virtualnetwork[0].azapi_resource.rg",
		"module.virtualnetwork[0].azapi_resource.vnet",
		"module.virtualnetwork[0].azapi_update_resource.vnet",
	}
	assert.Lenf(t, plan.ResourcePlannedValuesMap, len(resources), "expected %d resources to be created, but got %d", len(resources), len(plan.ResourcePlannedValuesMap))
	for _, v := range resources {
		terraform.AssertPlannedValuesMapKeyExists(t, plan, v)
	}
	// check bit field is correct
	telem := plan.ResourcePlannedValuesMap["azapi_resource.telemetry_root[0]"]
	require.Contains(t, telem.AttributeValues, "name")
	telemName := telem.AttributeValues["name"].(string)
	telemBf := strings.Split(telemName, "_")[2]
	expectBf := "00000300"
	assert.Equalf(t, expectBf, telemBf, "expected bit field to be %s, but got %s", expectBf, telemBf)
}

// TestIntegrationHubAndSpokeExistingSubscriptionWithMgAssoc tests the resource plan when supplying an existing subscription,
// with a new virtual network with peerings to a supplied hub network.
func TestIntegrationHubAndSpokeExistingSubscriptionWithMgAssoc(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet"
	v["subscription_alias_enabled"] = false
	v["subscription_id"] = "00000000-0000-0000-0000-000000000000"
	v["virtual_network_enabled"] = true
	v["virtual_network_peering_enabled"] = true
	v["subscription_management_group_association_enabled"] = true
	v["subscription_management_group_id"] = "Test"
	delete(v, "subscription_tags")
	terraformOptions.Vars = v

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	resources := []string{
		"azapi_resource.telemetry_root[0]",
		"module.virtualnetwork[0].azapi_resource.peering[\"inbound\"]",
		"module.virtualnetwork[0].azapi_resource.peering[\"outbound\"]",
		"module.virtualnetwork[0].azapi_resource.rg",
		"module.virtualnetwork[0].azapi_resource.vnet",
		"module.virtualnetwork[0].azapi_update_resource.vnet",
		"module.subscription[0].azurerm_management_group_subscription_association.this[0]",
	}
	assert.Lenf(t, plan.ResourcePlannedValuesMap, len(resources), "expected %d resources to be created, but got %d", len(resources), len(plan.ResourcePlannedValuesMap))
	for _, v := range resources {
		terraform.AssertPlannedValuesMapKeyExists(t, plan, v)
	}
	// check bit field is correct
	telem := plan.ResourcePlannedValuesMap["azapi_resource.telemetry_root[0]"]
	require.Contains(t, telem.AttributeValues, "name")
	telemName := telem.AttributeValues["name"].(string)
	telemBf := strings.Split(telemName, "_")[2]
	expectBf := "00000302"
	assert.Equalf(t, expectBf, telemBf, "expected bit field to be %s, but got %s", expectBf, telemBf)
}

// TestIntegrationWithYaml tests the use of the module with a for_each loop
// using YAML files as input.
func TestIntegrationWithYaml(t *testing.T) {
	testDir := "testdata/" + t.Name()
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, testDir)
	defer cleanup()
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	err = utils.GenerateRequiredProvidersFile(utils.NewRequiredProvidersData(), filepath.Clean(tmp+"/terraform.tf"))
	require.NoErrorf(t, err, "failed to create terraform.tf: %v", err)
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	require.NoErrorf(t, err, "failed to generate plan: %v", err)

	resources := []string{
		"module.alz_landing_zone[\"%s\"].azapi_resource.telemetry_root[0]",
		"module.alz_landing_zone[\"%s\"].module.virtualnetwork[0].azapi_update_resource.vnet",
		"module.alz_landing_zone[\"%s\"].module.virtualnetwork[0].azapi_resource.vnet",
		"module.alz_landing_zone[\"%s\"].module.virtualnetwork[0].azapi_resource.rg_lock[0]",
		"module.alz_landing_zone[\"%s\"].module.virtualnetwork[0].azapi_resource.rg",
		"module.alz_landing_zone[\"%s\"].module.subscription[0].azurerm_subscription.this[0]",
		"module.alz_landing_zone[\"%s\"].module.subscription[0].azurerm_management_group_subscription_association.this[0]",
		"module.alz_landing_zone[\"%s\"].module.roleassignment[\"8a6eec3e-78d9-5ff3-89cc-b144ae761a9a\"].azurerm_role_assignment.this",
		"module.alz_landing_zone[\"%s\"].module.roleassignment[\"7f69efa3-575a-5f8b-a989-c3978b92b58a\"].azurerm_role_assignment.this",
	}
	assert.Lenf(t, plan.ResourcePlannedValuesMap, len(resources)*3, "expected %d resources to be created, but got %d", len(resources)*3, len(plan.ResourcePlannedValuesMap))
	lzs := []string{
		"landing_zone_1.yaml",
		"landing_zone_2.yaml",
		"landing_zone_3.yaml",
	}
	for _, r := range resources {
		for _, lz := range lzs {
			res := fmt.Sprintf(r, lz)
			terraform.AssertPlannedValuesMapKeyExists(t, plan, res)
		}
	}
}

// TestIntegrationHubAndSpoke tests the resource plan when creating a new subscription,
// with a new virtual network with peerings to a supplied hub network.
func TestIntegrationDisableTelemetry(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	defer cleanup()
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet"
	v["subscription_alias_enabled"] = true
	v["virtual_network_enabled"] = true
	v["virtual_network_peering_enabled"] = true
	v["virtual_network_resource_lock_enabled"] = true
	v["disable_telemetry"] = true
	terraformOptions.Vars = v

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	resources := []string{
		"module.subscription[0].azurerm_subscription.this[0]",
		"module.virtualnetwork[0].azapi_resource.peering[\"inbound\"]",
		"module.virtualnetwork[0].azapi_resource.peering[\"outbound\"]",
		"module.virtualnetwork[0].azapi_resource.rg_lock[0]",
		"module.virtualnetwork[0].azapi_resource.rg",
		"module.virtualnetwork[0].azapi_resource.vnet",
		"module.virtualnetwork[0].azapi_update_resource.vnet",
	}
	assert.Lenf(t, plan.ResourcePlannedValuesMap, len(resources), "expected %d resources to be created, but got %d", len(resources), len(plan.ResourcePlannedValuesMap))
	for _, v := range resources {
		terraform.AssertPlannedValuesMapKeyExists(t, plan, v)
	}
}

func getMockInputVariables() map[string]interface{} {
	return map[string]interface{}{
		// subscription variables
		"subscription_billing_scope": "/providers/Microsoft.Billing/billingAccounts/0000000/enrollmentAccounts/000000",
		"subscription_display_name":  "test-subscription-alias",
		"subscription_alias_name":    "test-subscription-alias",
		"subscription_workload":      "Production",
		"subscription_tags": map[string]interface{}{
			"test-tag":   "test-value",
			"test-tag-2": "test-value-2",
		},

		// virtualnetwork variables
		"virtual_network_address_space":         []string{"10.1.0.0/24", "172.16.1.0/24"},
		"virtual_network_location":              "northeurope",
		"virtual_network_name":                  "testvnet",
		"virtual_network_resource_group_name":   "testrg",
		"virtual_network_resource_lock_enabled": false,
	}
}
