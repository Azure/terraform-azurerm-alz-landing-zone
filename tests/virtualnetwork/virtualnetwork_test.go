package virtualnetwork

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../modules/virtualnetwork"
)

// TestVirtualNetworkCreateValid tests the creation of a plan that
// creates two virtual networks in the specified resource groups with custom DNS servers.
func TestVirtualNetworkCreateValidWithCustomDns(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	primaryvnet["dns_servers"] = []any{"1.2.3.4", "4.3.2.1"}
	secondaryvnet["dns_servers"] = []any{}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// want 8 resources, like TestVirtualNetworkCreateValid
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(8).ErrorIsNilFatal(t)

	// Loop through each virtual network and check the values
	vns := v["virtual_networks"].(map[string]map[string]any)
	for k, v := range vns {
		res := fmt.Sprintf("azapi_resource.vnet[\"%s\"]", k)
		check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.dhcpOptions.dnsServers").HasValue(v["dns_servers"]).ErrorIsNil(t)
	}
}

// TestVirtualNetworkCreateValidWithTags tests the creation of a plan that
// creates two virtual networks in the specified resource groups with tags on vnet and rg.
func TestVirtualNetworkCreateValidWithTags(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["tags"] = map[string]any{
		"tag1": "value1",
		"tag2": "2",
	}
	primaryvnet["resource_group_tags"] = map[string]any{
		"tag1": "value1",
		"tag2": "2",
	}
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 8 resources here, same as TestVirtualNetworkCreateValid test
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(8).ErrorIsNilFatal(t)

	check.InPlan(test.PlanStruct).That("azapi_resource.vnet[\"primary\"]").Key("tags").HasValue(primaryvnet["tags"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That("azapi_resource.rg[\"primary-rg\"]").Key("tags").HasValue(primaryvnet["resource_group_tags"]).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithMeshPeering tests the creation of a plan that
// creates two virtual networks in the specified resource groups with mesh peering.
func TestVirtualNetworkCreateValidWithMeshPeering(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	secondaryvnet := v["virtual_networks"].(map[string]map[string]any)["secondary"]
	primaryvnet["mesh_peering_enabled"] = true
	secondaryvnet["mesh_peering_enabled"] = true
	secondaryvnet["mesh_peering_allow_forwarded_traffic"] = true

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 10 resources here, 2 more than the TestVirtualNetworkCreateValid test
	// The additional two are the inbound & outbound peering
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(10).ErrorIsNilFatal(t)

	peer1 := "azapi_resource.peering_mesh[\"primary-secondary\"]"
	check.InPlan(test.PlanStruct).That(peer1).Key("body").Query("properties.allowForwardedTraffic").HasValue(false).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(peer1).Key("body").Query("properties.allowVirtualNetworkAccess").HasValue(true).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(peer1).Key("body").Query("properties.allowGatewayTransit").HasValue(false).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(peer1).Key("body").Query("properties.useRemoteGateways").HasValue(false).ErrorIsNil(t)
	peer1Remote := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/secondary-rg/providers/Microsoft.Network/virtualNetworks/secondary-vnet"
	check.InPlan(test.PlanStruct).That(peer1).Key("body").Query("properties.remoteVirtualNetwork.id").HasValue(peer1Remote).ErrorIsNil(t)

	peer2 := "azapi_resource.peering_mesh[\"secondary-primary\"]"
	check.InPlan(test.PlanStruct).That(peer2).Key("body").Query("properties.allowForwardedTraffic").HasValue(true).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(peer2).Key("body").Query("properties.allowVirtualNetworkAccess").HasValue(true).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(peer2).Key("body").Query("properties.allowGatewayTransit").HasValue(false).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(peer2).Key("body").Query("properties.useRemoteGateways").HasValue(false).ErrorIsNil(t)
	peer2Remote := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/primary-rg/providers/Microsoft.Network/virtualNetworks/primary-vnet"
	check.InPlan(test.PlanStruct).That(peer2).Key("body").Query("properties.remoteVirtualNetwork.id").HasValue(peer2Remote).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidInvalidMeshPeering tests the creation of a plan that
// creates two virtual networks in the specified resource groups with mesh peering
// enabled on only one of the two vnets.
func TestVirtualNetworkCreateValidInvalidMeshPeering(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["mesh_peering_enabled"] = true
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 8 resources here, as only one of the two vnets has mesh peering enabled, then no peerings should be created
	resources := []string{
		"azapi_resource.rg[\"primary-rg\"]",
		"azapi_resource.rg[\"secondary-rg\"]",
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.rg_lock[\"primary-rg\"]",
		"azapi_resource.rg_lock[\"secondary-rg\"]",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}
}

// TestVirtualNetworkCreateValidSameRg tests the creation of a plan that
// creates two virtual networks in the same resource group.
func TestVirtualNetworkCreateValidSameRg(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["resource_group_name"] = "secondary-rg"
	primaryvnet["resource_group_creation_enabled"] = false

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 6 resources here, as the two vnets have the same rg, then 2 fewer resources than
	// TestVirtualNetworkCreateValid (rg + rg lock)
	resources := []string{
		"azapi_resource.rg[\"secondary-rg\"]",
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.rg_lock[\"secondary-rg\"]",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}
}

// TestVirtualNetworkCreateValidSameRgSameLocation tests the creation of a plan that
// creates two virtual networks in the same resource group in the same location.
func TestVirtualNetworkCreateValidSameRgSameLocation(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["resource_group_name"] = "secondary-rg"
	primaryvnet["location"] = "northeurope"

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 6 resources here, as the two vnets have the same rg, then 2 fewer resources than
	// TestVirtualNetworkCreateValid (rg + rg lock)
	resources := []string{
		"azapi_resource.rg[\"secondary-rg\"]",
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.rg_lock[\"secondary-rg\"]",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}
}

// TestVirtualNetworkCreateValidWithPeering tests the creation of a plan that
// creates a virtual network with bidirectional peering to a hub.
func TestVirtualNetworkCreateValidWithHubPeering(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()

	// Enable hub network peering to primary vnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet2"
	primaryvnet["hub_peering_enabled"] = true

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 10 resources here, 2 more than the TestVirtualNetworkCreateValid test
	// The additional two are the inbound & outbound peering
	resources := []string{
		"azapi_resource.rg[\"primary-rg\"]",
		"azapi_resource.rg[\"secondary-rg\"]",
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.rg_lock[\"primary-rg\"]",
		"azapi_resource.rg_lock[\"secondary-rg\"]",
		"azapi_resource.peering_hub_outbound[\"primary\"]",
		"azapi_resource.peering_hub_inbound[\"primary\"]",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	// We can only check the body of the outbound peering as the inbound values
	// are not known until apply
	outbound := "azapi_resource.peering_hub_outbound[\"primary\"]"
	terraform.RequirePlannedValuesMapKeyExists(t, test.PlanStruct, outbound)

	check.InPlan(test.PlanStruct).That(outbound).Key("body").Query("properties.allowForwardedTraffic").HasValue(true).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(outbound).Key("body").Query("properties.allowVirtualNetworkAccess").HasValue(true).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(outbound).Key("body").Query("properties.allowGatewayTransit").HasValue(false).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(outbound).Key("body").Query("properties.useRemoteGateways").HasValue(true).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(outbound).Key("body").Query("properties.remoteVirtualNetwork.id").HasValue(primaryvnet["hub_network_resource_id"]).ErrorIsNil(t)

	// More limited checks on the inbound peering
	inbound := "azapi_resource.peering_hub_inbound[\"primary\"]"
	check.InPlan(test.PlanStruct).That(inbound).Key("parent_id").HasValue(primaryvnet["hub_network_resource_id"]).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithPeeringCustomNames tests the creation of a plan that
// creates a virtual network with bidirectional peering to a hub, with custom names for peers.
func TestVirtualNetworkCreateValidWithHubPeeringCustomNames(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()

	// Enable hub network peering to primary vnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet2"
	primaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_peering_name_tohub"] = "test-tohub"
	primaryvnet["hub_peering_name_fromhub"] = "test-fromhub"

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 10 resources here, 2 more than the TestVirtualNetworkCreateValid test
	// The additional two are the inbound & outbound peering
	resources := []string{
		"azapi_resource.rg[\"primary-rg\"]",
		"azapi_resource.rg[\"secondary-rg\"]",
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.rg_lock[\"primary-rg\"]",
		"azapi_resource.rg_lock[\"secondary-rg\"]",
		"azapi_resource.peering_hub_outbound[\"primary\"]",
		"azapi_resource.peering_hub_inbound[\"primary\"]",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	// Check outbound peering name
	outbound := "azapi_resource.peering_hub_outbound[\"primary\"]"
	check.InPlan(test.PlanStruct).That(outbound).Key("name").HasValue(primaryvnet["hub_peering_name_tohub"]).ErrorIsNil(t)

	// Check inbound peering name
	inbound := "azapi_resource.peering_hub_inbound[\"primary\"]"
	check.InPlan(test.PlanStruct).That(inbound).Key("name").HasValue(primaryvnet["hub_peering_name_fromhub"]).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithOnlyToHubPeering tests the creation of a plan that
// creates a virtual network with unidirectional peering to a hub, with custom names for peers.
func TestVirtualNetworkCreateValidWithOnlyToHubPeering(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()

	// Enable hub network peering to primary vnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet2"
	primaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_peering_direction"] = "tohub"
	primaryvnet["hub_peering_name_tohub"] = "test-tohub"
	primaryvnet["hub_peering_name_fromhub"] = "test-fromhub"

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 9 resources here, 1 more than the TestVirtualNetworkCreateValid test
	// The additional one is the outbound peering
	resources := []string{
		"azapi_resource.rg[\"primary-rg\"]",
		"azapi_resource.rg[\"secondary-rg\"]",
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.rg_lock[\"primary-rg\"]",
		"azapi_resource.rg_lock[\"secondary-rg\"]",
		"azapi_resource.peering_hub_outbound[\"primary\"]",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}
}

// TestVirtualNetworkCreateValidWithOnlyFromHubPeering tests the creation of a plan that
// creates a virtual network with unidirectional peering from a hub, with custom names for peers.
func TestVirtualNetworkCreateValidWithOnlyFromHubPeering(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()

	// Enable hub network peering to primary vnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet2"
	primaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_peering_direction"] = "fromhub"
	primaryvnet["hub_peering_name_tohub"] = "test-tohub"
	primaryvnet["hub_peering_name_fromhub"] = "test-fromhub"

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 10 resources here, 2 more than the TestVirtualNetworkCreateValid test
	// The additional two are the inbound & outbound peering
	resources := []string{
		"azapi_resource.rg[\"primary-rg\"]",
		"azapi_resource.rg[\"secondary-rg\"]",
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.rg_lock[\"primary-rg\"]",
		"azapi_resource.rg_lock[\"secondary-rg\"]",
		"azapi_resource.peering_hub_inbound[\"primary\"]",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}
}

// TestVirtualNetworkCreateValidWithPeeringUseRemoteGatewaysDisabled
// tests the creation of a plan that configured the outbound peering
// with useRemoteGateways disabled.
func TestVirtualNetworkCreateValidWithPeeringUseRemoteGatewaysDisabled(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	// Enable hub network peering to primary vnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet2"
	primaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_peering_use_remote_gateways"] = false

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 10 resources here, 2 more than the TestVirtualNetworkCreateValid test
	// The additional two are the inbound & outbound peering
	resources := []string{
		"azapi_resource.rg[\"primary-rg\"]",
		"azapi_resource.rg[\"secondary-rg\"]",
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.rg_lock[\"primary-rg\"]",
		"azapi_resource.rg_lock[\"secondary-rg\"]",
		"azapi_resource.peering_hub_outbound[\"primary\"]",
		"azapi_resource.peering_hub_inbound[\"primary\"]",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	// We can only check the body of the outbound peering as the inbound values
	// not known until apply
	res := "azapi_resource.peering_hub_outbound[\"primary\"]"
	check.InPlan(test.PlanStruct).That(res).Key("body").Query("properties.useRemoteGateways").HasValue(false).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithVhub tests the creation of a plan that
// creates a virtual network with a vhub connection.
func TestVirtualNetworkCreateValidWithVhub(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()

	// Enable vhub connection to primary vnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	primaryvnet["vwan_connection_enabled"] = true

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 9 resources here, 1 more than the TestVirtualNetworkCreateValid test
	// The additional resource is the vhub connection
	resources := []string{
		"azapi_resource.rg[\"primary-rg\"]",
		"azapi_resource.rg[\"secondary-rg\"]",
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.rg_lock[\"primary-rg\"]",
		"azapi_resource.rg_lock[\"secondary-rg\"]",
		"azapi_resource.vhubconnection[\"primary\"]",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vhcres := "azapi_resource.vhubconnection[\"primary\"]"
	check.InPlan(test.PlanStruct).That(vhcres).Key("parent_id").HasValue(primaryvnet["vwan_hub_resource_id"]).ErrorIsNil(t)

	drt := primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/defaultRouteTable"
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.routingConfiguration.associatedRouteTable.id").HasValue(drt).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.routingConfiguration.propagatedRouteTables.labels").HasValue([]any{"default"}).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.routingConfiguration.propagatedRouteTables.ids.#").HasValue(1).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.routingConfiguration.propagatedRouteTables.ids.0.id").HasValue(drt).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithVhubCustomRouting tests the creation of a plan that
// creates a virtual network with a vhub connection with custom routing.
func TestVirtualNetworkCreateValidWithVhubCustomRouting(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()

	// Enable vhub connection to primary vnet in test mock input variables
	// & add custom routing
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	primaryvnet["vwan_connection_enabled"] = true
	primaryvnet["vwan_propagated_routetables_labels"] = []any{"testlabel", "testlabel2"}
	primaryvnet["vwan_propagated_routetables_resource_ids"] = []any{
		primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/testRouteTable",
		primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/testRouteTable2",
	}
	primaryvnet["vwan_associated_routetable_resource_id"] = primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/testRouteTable3"

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 9 resources here, 1 more than the TestVirtualNetworkCreateValid test
	// The additional resource is the vhub connection
	resources := []string{
		"azapi_resource.rg[\"primary-rg\"]",
		"azapi_resource.rg[\"secondary-rg\"]",
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.rg_lock[\"primary-rg\"]",
		"azapi_resource.rg_lock[\"secondary-rg\"]",
		"azapi_resource.vhubconnection[\"primary\"]",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vhcres := "azapi_resource.vhubconnection[\"primary\"]"
	check.InPlan(test.PlanStruct).That(vhcres).Key("parent_id").HasValue(primaryvnet["vwan_hub_resource_id"]).ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.routingConfiguration.associatedRouteTable.id").HasValue(primaryvnet["vwan_associated_routetable_resource_id"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.routingConfiguration.propagatedRouteTables.labels").HasValue(primaryvnet["vwan_propagated_routetables_labels"]).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.routingConfiguration.propagatedRouteTables.ids.#.id").HasValue(primaryvnet["vwan_propagated_routetables_resource_ids"]).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithVhubSecureInternetTraffic tests that secure_internet_traffic == true
func TestVirtualNetworkCreateValidWithVhubSecureInternetTraffic(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()

	// Enable vhub connection to primary vnet in test mock input variables
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	primaryvnet["vwan_connection_enabled"] = true
	primaryvnet["vwan_security_configuration"] = map[string]any{
		"secure_internet_traffic": true,
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 9 resources here, 1 more than the TestVirtualNetworkCreateValid test
	// The additional resource is the vhub connection
	resources := []string{
		"azapi_resource.rg[\"primary-rg\"]",
		"azapi_resource.rg[\"secondary-rg\"]",
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.rg_lock[\"primary-rg\"]",
		"azapi_resource.rg_lock[\"secondary-rg\"]",
		"azapi_resource.vhubconnection[\"primary\"]",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vhcres := "azapi_resource.vhubconnection[\"primary\"]"
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.enableInternetSecurity").HasValue(true).ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithVhubSecurePrivateTraffic that managed vnets propagate to "noneRouteTable" with labels "none"
func TestVirtualNetworkCreateValidWithVhubSecurePrivateTraffic(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()

	// Enable vhub connection to primary vnet in test mock input variables
	// & add custom routing
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	primaryvnet["vwan_connection_enabled"] = true
	primaryvnet["vwan_security_configuration"] = map[string]any{
		"secure_private_traffic": true,
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 9 resources here, 1 more than the TestVirtualNetworkCreateValid test
	// The additional resource is the vhub connection
	resources := []string{
		"azapi_resource.rg[\"primary-rg\"]",
		"azapi_resource.rg[\"secondary-rg\"]",
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.rg_lock[\"primary-rg\"]",
		"azapi_resource.rg_lock[\"secondary-rg\"]",
		"azapi_resource.vhubconnection[\"primary\"]",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vhcres := "azapi_resource.vhubconnection[\"primary\"]"

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").
		Query("properties.routingConfiguration.associatedRouteTable.id").
		HasValue(primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/defaultRouteTable").
		ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").
		Query("properties.routingConfiguration.propagatedRouteTables.labels").
		HasValue([]any{"none"}).
		ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").
		Query("properties.routingConfiguration.propagatedRouteTables.ids.0.id").
		HasValue(primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/noneRouteTable").
		ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").
		Query("properties.routingConfiguration.propagatedRouteTables.ids.#").
		HasValue(1).
		ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithVhubSecureInternetAndPrivateTraffic tests secure_internet_traffic == true
// and that managed vnets propagate to "noneRouteTable" with labels "none"
func TestVirtualNetworkCreateValidWithVhubSecureInternetAndPrivateTraffic(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()

	// Enable vhub connection to primary vnet in test mock input variables
	// & add custom routing
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	primaryvnet["vwan_connection_enabled"] = true
	primaryvnet["vwan_security_configuration"] = map[string]any{
		"secure_internet_traffic": true,
		"secure_private_traffic":  true,
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 9 resources here, 1 more than the TestVirtualNetworkCreateValid test
	// The additional resource is the vhub connection
	resources := []string{
		"azapi_resource.rg[\"primary-rg\"]",
		"azapi_resource.rg[\"secondary-rg\"]",
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.rg_lock[\"primary-rg\"]",
		"azapi_resource.rg_lock[\"secondary-rg\"]",
		"azapi_resource.vhubconnection[\"primary\"]",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vhcres := "azapi_resource.vhubconnection[\"primary\"]"
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.enableInternetSecurity").HasValue(true).ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").
		Query("properties.routingConfiguration.associatedRouteTable.id").
		HasValue(primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/defaultRouteTable").
		ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").
		Query("properties.routingConfiguration.propagatedRouteTables.labels").
		HasValue([]any{"none"}).
		ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").
		Query("properties.routingConfiguration.propagatedRouteTables.ids.0.id").
		HasValue(primaryvnet["vwan_hub_resource_id"].(string) + "/hubRouteTables/noneRouteTable").
		ErrorIsNil(t)

	check.InPlan(test.PlanStruct).That(vhcres).Key("body").
		Query("properties.routingConfiguration.propagatedRouteTables.ids.#").
		HasValue(1).
		ErrorIsNil(t)
}

// TestVirtualNetworkCreateValidWithVhubRoutingIntentEnabled tests that routingConfiguration is null when
// routing intent is enabled
func TestVirtualNetworkCreateValidWithVhubRoutingIntentEnabled(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()

	// Enable vhub connection to primary vnet in test mock input variables
	// & add custom routing
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"
	primaryvnet["vwan_connection_enabled"] = true
	primaryvnet["vwan_security_configuration"] = map[string]any{
		"routing_intent_enabled": true,
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// We want 9 resources here, 1 more than the TestVirtualNetworkCreateValid test
	// The additional resource is the vhub connection
	resources := []string{
		"azapi_resource.rg[\"primary-rg\"]",
		"azapi_resource.rg[\"secondary-rg\"]",
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.rg_lock[\"primary-rg\"]",
		"azapi_resource.rg_lock[\"secondary-rg\"]",
		"azapi_resource.vhubconnection[\"primary\"]",
	}
	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)

	for _, r := range resources {
		check.InPlan(test.PlanStruct).That(r).Exists().ErrorIsNil(t)
	}

	vhcres := "azapi_resource.vhubconnection[\"primary\"]"
	check.InPlan(test.PlanStruct).That(vhcres).Key("body").Query("properties.routingConfiguration").DoesNotExist().ErrorIsNil(t)
}

// TestVirtualNetworkCreateInvalidHubNetResId tests the regex of the
// hub_network_resource_id variable.
func TestVirtualNetworkCreateInvalidHubNetResId(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["hub_peering_enabled"] = true
	primaryvnet["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroup/testrg/providers/Microsoft.Network/virtualNetworks/tes.-tvnet2"

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	assert.ErrorContains(t, err, "Hub network resource id must be an Azure virtual network resource id")
}

// TestVirtualNetworkCreateInvalidVhubResId tests the regex of the
// hub_network_resource_id variable.
func TestVirtualNetworkCreateInvalidVhubResId(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["vwan_connection_enabled"] = true
	primaryvnet["vwan_hub_resource_id"] = "/subscription/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/virtualHubs/te.st-hub"

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	assert.ErrorContains(t, err, "vWAN hub resource id must be an Azure vWAN hub network resource id")
}

// TestVirtualNetworkCreateZeroLengthAddressSpace tests the length of address_space > 0
func TestVirtualNetworkCreateZeroLengthAddressSpace(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["address_space"] = []string{}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	assert.ErrorContains(t, err, "At least 1 address space must be specified")
}

// TestVirtualNetworkCreateInvalidAddressSpace tests a valid CIDR address space is used
func TestVirtualNetworkCreateInvalidAddressSpace(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["address_space"] = []string{"10.37.242/35"}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	assert.ErrorContains(t, err, "Address space entries must be specified in CIDR notation")
}

// TestVirtualNetworkCreateInvalidResourceGroupCreation tests that resource group naming is unique
// when using vnets in multiple locaitons that share a resoruce group.
// NOTE - this is not a recommended deployment pattern.
func TestVirtualNetworkCreateInvalidResourceGroupCreation(t *testing.T) {
	t.Parallel()

	v := getMockInputVariables()
	primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
	primaryvnet["resource_group_name"] = "secondary-rg"

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	defer test.Cleanup()
	assert.Containsf(t, utils.SanitiseErrorMessage(err), "Resource group names with creation enabled must be unique. Virtual networks deployed into the same resource group must have only one enabled for resource group creation.", "Expected error message not found")
}

func TestVirtualNetworkDdosProtection(t *testing.T) {
	t.Parallel()

	// We want 8 resources here
	resources := []string{
		"azapi_resource.rg[\"primary-rg\"]",
		"azapi_resource.rg[\"secondary-rg\"]",
		"azapi_resource.vnet[\"primary\"]",
		"azapi_resource.vnet[\"secondary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"secondary\"]",
		"azapi_resource.rg_lock[\"primary-rg\"]",
		"azapi_resource.rg_lock[\"secondary-rg\"]",
	}

	vnetresources := []string{
		"azapi_resource.vnet[\"primary\"]",
		"azapi_update_resource.vnet[\"primary\"]",
	}

	t.Run("Enabled", func(t *testing.T) {
		v := getMockInputVariables()
		primaryvnet := v["virtual_networks"].(map[string]map[string]any)["primary"]
		primaryvnet["ddos_protection_enabled"] = true
		primaryvnet["ddos_protection_plan_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/ddosProtectionPlans/test-ddos-plan"

		test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
		defer test.Cleanup()
		require.NoError(t, err)

		check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)
		for _, r := range vnetresources {
			check.InPlan(test.PlanStruct).That(r).Key("body").Query("properties.enableDdosProtection").HasValue(true).ErrorIsNil(t)
			check.InPlan(test.PlanStruct).That(r).Key("body").Query("properties.ddosProtectionPlan.id").HasValue(primaryvnet["ddos_protection_plan_id"]).ErrorIsNil(t)
		}
	})

	t.Run("Disabled", func(t *testing.T) {
		v := getMockInputVariables()

		test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
		defer test.Cleanup()
		require.NoError(t, err)

		check.InPlan(test.PlanStruct).NumberOfResourcesEquals(len(resources)).ErrorIsNilFatal(t)
		for _, r := range vnetresources {
			check.InPlan(test.PlanStruct).That(r).Key("body").Query("properties.enableDdosProtection").DoesNotExist().ErrorIsNil(t)
			check.InPlan(test.PlanStruct).That(r).Key("body").Query("properties.ddosProtectionPlan.id").DoesNotExist().ErrorIsNil(t)
		}
	})
}

// getMockInputVariables returns a set of mock input variables that can be used and modified for testing scenarios.
func getMockInputVariables() map[string]any {
	return map[string]any{
		"subscription_id": "00000000-0000-0000-0000-000000000000",
		"virtual_networks": map[string]map[string]any{
			"primary": {
				"name":                "primary-vnet",
				"address_space":       []any{"192.168.0.0/24"},
				"location":            "westeurope",
				"resource_group_name": "primary-rg",
			},
			"secondary": {
				"name":                "secondary-vnet",
				"address_space":       []any{"192.168.1.0/24"},
				"location":            "northeurope",
				"resource_group_name": "secondary-rg",
			},
		},
	}
}
