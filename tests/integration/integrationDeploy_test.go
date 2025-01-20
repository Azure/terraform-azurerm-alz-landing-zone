package integration

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/azureutils"
	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeployIntegrationHubAndSpoke(t *testing.T) {
	t.Parallel()

	utils.PreCheckDeployTests(t)
	testDir := "testdata/" + t.Name()
	v, err := getValidInputVariables()
	require.NoErrorf(t, err, "could not generate valid input variables")
	test, err := setuptest.Dirs(moduleDir, testDir).WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	// get the random hex name from vars
	name := v["subscription_alias_name"].(string)

	// List of resources to find in the plan, excluding the role assignment
	resources := []string{
		"azurerm_resource_group.hub",
		"azurerm_virtual_network.hub",
		"module.lz_vending.azapi_resource.telemetry_root[0]",
		"module.lz_vending.module.subscription[0].azurerm_subscription.this[0]",
		"module.lz_vending.module.virtualnetwork[0].azapi_resource.peering_hub_inbound[\"primary\"]",
		"module.lz_vending.module.virtualnetwork[0].azapi_resource.peering_hub_outbound[\"primary\"]",
		fmt.Sprintf("module.lz_vending.module.virtualnetwork[0].azapi_resource.rg_lock[\"%s\"]", name),
		fmt.Sprintf("module.lz_vending.module.virtualnetwork[0].azapi_resource.rg[\"%s\"]", name),
		"module.lz_vending.module.virtualnetwork[0].azapi_resource.vnet[\"primary\"]",
	}

	for _, v := range resources {
		check.InPlan(test.PlanStruct).That(v).Exists().ErrorIsNil(t)
	}

	// As the map key of the role assignment is a predictable UUID based on the object ID
	// of the calling identity, we cannot search for the literal value of the role assignment.
	// Instead, we search for the role assignment prefix in the ResourcePlannedValuesMap.
	i := 0
	for k := range test.PlanStruct.ResourceChangesMap {
		if !strings.Contains(k, "module.lz_vending.module.roleassignment[") {
			continue
		}
		i++
	}
	require.Equal(t, 1, i, "expected 1 role assignment to be planned, got %d", i)

	// Defer the cleanup of the subscription alias to the end of the test.
	// Should be run after the Terraform destroy.
	// We don't know the sub ID yet, so use zeros for now and then
	// update it after the apply.
	u := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	defer func() {
		err := azureutils.CancelSubscription(t, &u)
		if err != nil {
			t.Logf("failed to cancel subscription: %v", err)
		}
	}()

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	// due to eventual consistency issues
	rty := setuptest.Retry{
		Max:  3,
		Wait: 10 * time.Minute,
	}
	defer test.DestroyRetry(rty) //nolint:errcheck
	test.ApplyIdempotent().ErrorIsNil(t)

	id, err := test.Output("subscription_id").GetValue()
	assert.NoErrorf(t, err, "failed to get subscription id output")
	ids, ok := id.(string)
	assert.Truef(t, ok, "failed to cast subscription id output to string")
	u, err = uuid.Parse(ids)
	assert.NoErrorf(t, err, "cannot parse subscription id as uuid: %s", id)
}

func TestDeployIntegrationResourceGroupsRpRegUmiAndRoleAssignments(t *testing.T) {
	t.Parallel()

	utils.PreCheckDeployTests(t)
	testDir := "testdata/" + t.Name()
	r, err := utils.RandomHex(4)
	require.NoError(t, err)
	v := map[string]any{
		"random_hex":      r,
		"subscription_id": os.Getenv("AZURE_SUBSCRIPTION_ID"),
	}
	test, err := setuptest.Dirs(moduleDir, testDir).WithVars(v).InitPlanShow(t)
	require.NoError(t, err)
	defer test.Cleanup()

	defer test.DestroyRetry(setuptest.FastRetry) //nolint:errcheck
	test.ApplyIdempotent().ErrorIsNil(t)
}

func getValidInputVariables() (map[string]any, error) {
	r, err := utils.RandomHex(4)
	if err != nil {
		return nil, fmt.Errorf("cannot generate random hex, %s", err)
	}
	name := fmt.Sprintf("testdeploy-%s", r)
	return map[string]any{
		"location":                                         "northeurope",
		"subscription_alias_name":                          name,
		"subscription_display_name":                        name,
		"subscription_billing_scope":                       os.Getenv("AZURE_BILLING_SCOPE"),
		"subscription_workload":                            "DevTest",
		"subscription_alias_enabled":                       true,
		"subscription_register_resource_providers_enabled": true,
		"virtual_network_enabled":                          true,
		"virtual_networks": map[string]map[string]any{
			"primary": {
				"name":                            name,
				"resource_group_name":             name,
				"location":                        "northeurope",
				"address_space":                   []string{"10.1.0.0/24", "172.16.1.0/24"},
				"hub_peering_enabled":             true,
				"hub_peering_use_remote_gateways": false,
			},
		},
		"role_assignment_enabled": true,
	}, nil
}
