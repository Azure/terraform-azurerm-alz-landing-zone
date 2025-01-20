package virtualnetwork

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/azureutils"
	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDeployVirtualNetworkValid tests the deployment of virtual networks
// with valid input variables.
func TestDeployVirtualNetworkValid(t *testing.T) {
	t.Parallel()

	utils.PreCheckDeployTests(t)
	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(6).ErrorIsNil(t)

	resources := []string{
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
	}
	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	// defer terraform destroy with retry
	defer test.DestroyRetry(setuptest.DefaultRetry) //nolint:errcheck
	test.ApplyIdempotent().ErrorIsNil(t)

	// check there two outputs for the virtual network resource ids
	test.Output("virtual_network_resource_ids").Query("primary").Exists().ErrorIsNil(t)
	test.Output("virtual_network_resource_ids").Query("secondary").Exists().ErrorIsNil(t)
}

// TestDeployVirtualNetworkValidCustomDns tests the deployment of virtual networks
// with valid input variables and custom DNS servers.
func TestDeployVirtualNetworkValidCustomDns(t *testing.T) {
	t.Parallel()

	utils.PreCheckDeployTests(t)
	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	primaryvnet["dns_servers"] = []string{"192.168.0.250", "192.168.0.251"}
	secondaryvnet["dns_servers"] = []string{"192.168.1.250", "192.168.1.251"}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(6).ErrorIsNil(t)

	resources := []string{
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
	}
	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	// defer terraform destroy with retry
	defer test.DestroyRetry(setuptest.DefaultRetry) //nolint:errcheck
	test.ApplyIdempotent().ErrorIsNil(t)

	// check there two outputs for the virtual network resource ids
	test.Output("virtual_network_resource_ids").Query("primary").Exists().ErrorIsNil(t)
	test.Output("virtual_network_resource_ids").Query("secondary").Exists().ErrorIsNil(t)
}

// TestDeployVirtualNetworkValidVnetPeering tests the deployment of a virtual network
// with bidirectional peering to a hub virtual network.
func TestDeployVirtualNetworkValidVnetPeering(t *testing.T) {
	t.Parallel()

	utils.PreCheckDeployTests(t)
	testDir := "testdata/" + t.Name()
	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	primaryvnet["hub_peering_enabled"] = true
	secondaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_peering_use_remote_gateways"] = false
	secondaryvnet["hub_peering_use_remote_gateways"] = false

	test, err := setuptest.Dirs(moduleDir, testDir).WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(12).ErrorIsNil(t)

	resources := []string{
		"module.virtualnetwork_test.azapi_resource.vnet[\"primary\"]",
		"module.virtualnetwork_test.azapi_resource.vnet[\"secondary\"]",
		"module.virtualnetwork_test.azapi_resource.peering_hub_inbound[\"primary\"]",
		"module.virtualnetwork_test.azapi_resource.peering_hub_inbound[\"secondary\"]",
		"module.virtualnetwork_test.azapi_resource.peering_hub_outbound[\"primary\"]",
		"module.virtualnetwork_test.azapi_resource.peering_hub_outbound[\"secondary\"]",
	}
	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	// defer terraform destroy with retry
	defer test.DestroyRetry(setuptest.DefaultRetry) //nolint:errcheck
	test.ApplyIdempotent().ErrorIsNil(t)
}

// TestDeployVirtualNetworkValidUniDirectionalVnetPeering tests the deployment of a virtual network
// with unidirectional peering to a hub virtual network.
func TestDeployVirtualNetworkValidUniDirectionalVnetPeering(t *testing.T) {
	t.Parallel()

	utils.PreCheckDeployTests(t)
	testDir := "testdata/" + t.Name()
	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	primaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_peering_direction"] = "fromhub"
	secondaryvnet["hub_peering_enabled"] = true
	secondaryvnet["hub_peering_direction"] = "tohub"
	primaryvnet["hub_peering_use_remote_gateways"] = false
	secondaryvnet["hub_peering_use_remote_gateways"] = false

	test, err := setuptest.Dirs(moduleDir, testDir).WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(10).ErrorIsNil(t)

	resources := []string{
		"module.virtualnetwork_test.azapi_resource.vnet[\"primary\"]",
		"module.virtualnetwork_test.azapi_resource.vnet[\"secondary\"]",
		"module.virtualnetwork_test.azapi_resource.peering_hub_inbound[\"primary\"]",
		"module.virtualnetwork_test.azapi_resource.peering_hub_outbound[\"secondary\"]",
	}
	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	// defer terraform destroy with retry
	defer test.DestroyRetry(setuptest.DefaultRetry) //nolint:errcheck
	test.ApplyIdempotent().ErrorIsNil(t)
}

// TestDeployVirtualNetworkValidVhubConnection tests the deployment of a virtual network
// with a virtual WAN connection.
func TestDeployVirtualNetworkValidVhubConnection(t *testing.T) {
	t.Parallel()

	utils.PreCheckDeployTests(t)
	testDir := "testdata/" + t.Name()
	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	primaryvnet["vwan_connection_enabled"] = true
	secondaryvnet["vwan_connection_enabled"] = true

	test, err := setuptest.Dirs(moduleDir, testDir).WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(11).ErrorIsNil(t)

	resources := []string{
		"module.virtualnetwork_test.azapi_resource.vnet[\"primary\"]",
		"module.virtualnetwork_test.azapi_resource.vnet[\"secondary\"]",
		"module.virtualnetwork_test.azapi_resource.vhubconnection[\"primary\"]",
		"module.virtualnetwork_test.azapi_resource.vhubconnection[\"secondary\"]",
	}
	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	// defer terraform destroy with retry
	rty := setuptest.Retry{
		Max:  3,
		Wait: 10 * time.Minute,
	}
	defer test.DestroyRetry(rty) //nolint:errcheck
	test.ApplyIdempotent().ErrorIsNil(t)
}

// TestDeployVirtualNetworkValidVhubConnectionAndRoutingIntent tests the deployment of a virtual network
// with a virtual WAN connection and routing intent.
func TestDeployVirtualNetworkValidVhubConnectionAndRoutingIntent(t *testing.T) {
	t.Parallel()

	utils.PreCheckDeployTests(t)
	testDir := "testdata/" + t.Name()
	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	primaryvnet["vwan_connection_enabled"] = true
	secondaryvnet["vwan_connection_enabled"] = true
	primaryvnet["vwan_security_configuration"] = map[string]any{
		"routing_intent_enabled": true,
	}
	secondaryvnet["vwan_security_configuration"] = map[string]any{
		"routing_intent_enabled": true,
	}

	test, err := setuptest.Dirs(moduleDir, testDir).WithVars(v).Init(t)
	require.NoError(t, utils.AzureRmAndRequiredProviders(test))

	require.NoError(t, err)
	defer test.Cleanup()

	// defer terraform destroy with retry
	rtyDestroy := setuptest.Retry{
		Max:  3,
		Wait: 10 * time.Minute,
	}
	rtyApply := setuptest.Retry{
		Max:  5,
		Wait: 5 * time.Minute,
	}
	defer test.DestroyRetry(rtyDestroy) //nolint:errcheck
	test.ApplyIdempotentRetry(rtyApply).ErrorIsNil(t)
}

// TestDeployVirtualNetworkSubnetIdempotency tests that we can make changes
// to the subnet configuration outside the module and that subsequent runs of terraform apply
// are idempotent. See main.tf file in the testdata directory for more details.
func TestDeployVirtualNetworkSubnetIdempotency(t *testing.T) {
	t.Parallel()

	utils.PreCheckDeployTests(t)
	testDir := "testdata/" + t.Name()
	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)

	test, err := setuptest.Dirs(moduleDir, testDir).WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// defer terraform destroy with retry
	defer test.DestroyRetry(setuptest.DefaultRetry) //nolint:errcheck
	test.ApplyIdempotent().ErrorIsNil(t)

	// test an update to vnet address space, then check for subnet still existing
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["address_space"] = []string{"192.168.0.0/23"}
	_, err = terraform.PlanE(t, test.Options)
	require.NoError(t, err)
	_, err = terraform.ApplyAndIdempotentE(t, test.Options)
	assert.NoError(t, err)
	name := primaryvnet["name"].(string)
	subnets, err := azureutils.ListSubnets(name, name, uuid.MustParse(os.Getenv("AZURE_SUBSCRIPTION_ID")))
	require.NoErrorf(t, err, "failed to list subnets")
	assert.Lenf(t, subnets, 1, "expected 1 subnet, got %d", len(subnets))
}

// TestDeployVirtualNetworkValidMeshPeering tests the deployment of virtual networks
// with mesh peering enables.
func TestDeployVirtualNetworkValidMeshPeering(t *testing.T) {
	t.Parallel()

	utils.PreCheckDeployTests(t)
	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables, %s", err)
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	primaryvnet["mesh_peering_enabled"] = true
	secondaryvnet["mesh_peering_enabled"] = true

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(8).ErrorIsNil(t)

	resources := []string{
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_resource.peering_mesh[\"primary-secondary\"]",
		"azapi_resource.peering_mesh[\"secondary-primary\"]",
	}
	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	// defer terraform destroy with retry
	defer test.DestroyRetry(setuptest.DefaultRetry) //nolint:errcheck
	test.ApplyIdempotent().ErrorIsNil(t)
}

func getValidInputVariables() (map[string]any, error) {
	r, err := utils.RandomHex(4)
	if err != nil {
		return nil, fmt.Errorf("cannot generate random hex, %s", err)
	}
	name := fmt.Sprintf("testdeploy-%s", r)
	name2 := name + "-2"

	return map[string]any{
		"subscription_id": os.Getenv("AZURE_SUBSCRIPTION_ID"),
		"virtual_networks": map[string]map[string]any{
			"primary": {
				"name":                        name,
				"address_space":               []string{"192.168.0.0/24"},
				"location":                    "westeurope",
				"resource_group_name":         name,
				"resource_group_lock_enabled": false,
			},
			"secondary": {
				"name":                        name2,
				"address_space":               []string{"192.168.1.0/24"},
				"location":                    "northeurope",
				"resource_group_name":         name2,
				"resource_group_lock_enabled": false,
			},
		},
	}, nil
}
