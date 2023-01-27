package subscription

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/azureutils"
	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/matt-FFFFFF/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var billingScope = os.Getenv("AZURE_BILLING_SCOPE")

// TestDeploySubscriptionAliasValid tests the deployment of a subscription alias
// with valid input variables.
func TestDeploySubscriptionAliasValid(t *testing.T) {
	t.Parallel()

	utils.PreCheckDeployTests(t)

	v, err := getValidInputVariables(billingScope)
	require.NoError(t, err)
	test := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, test.Err)
	defer test.Cleanup()
	require.NoError(t, test.Err)

	// Defer the cleanup of the subscription alias to the end of the test.
	// Should be run after the Terraform destroy.
	// We don't know the sub ID yet, so use zeros for now and then
	// update it after the apply.
	u := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	defer func() {
		err := azureutils.CancelSubscription(t, &u)
		t.Logf("cannot cancel subscription: %v", err)
	}()

	defer test.DestroyWithRetry(t, setuptest.DefaultRetry)
	err = test.ApplyIdempotent(t)
	assert.NoError(t, err)

	sid, err := terraform.OutputE(t, test.Options, "subscription_id")
	assert.NoError(t, err)
	u, err = uuid.Parse(sid)
	require.NoErrorf(t, err, "subscription id %s is not a valid uuid", sid)
}

// TestDeploySubscriptionAliasManagementGroupValid tests the deployment of a subscription alias
// with valid input variables.
func TestDeploySubscriptionAliasManagementGroupValid(t *testing.T) {
	t.Parallel()
	utils.PreCheckDeployTests(t)

	v, err := getValidInputVariables(billingScope)
	require.NoError(t, err)
	v["subscription_billing_scope"] = billingScope
	v["subscription_management_group_id"] = v["subscription_alias_name"]
	v["subscription_management_group_association_enabled"] = true

	testDir := filepath.Join("testdata", t.Name())
	test := setuptest.Dirs(moduleDir, testDir).WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, test.Err)
	defer test.Cleanup()
	require.NoError(t, test.Err)

	// Defer the cleanup of the subscription alias to the end of the test.
	// Should be run after the Terraform destroy.
	// We don't know the sub ID yet, so use zeros for now and then
	// update it after the apply.
	u := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	defer func() {
		err := azureutils.CancelSubscription(t, &u)
		t.Logf("cannot cancel subscription: %v", err)
	}()

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	// due to eventual consistency of the subscription aliases API
	defer test.DestroyWithRetry(t, setuptest.DefaultRetry)
	err = test.ApplyIdempotent(t)
	assert.NoError(t, err)

	sid, err := terraform.OutputE(t, test.Options, "subscription_id")
	assert.NoError(t, err)

	u, err = uuid.Parse(sid)
	assert.NoErrorf(t, err, "subscription id %s is not a valid uuid", sid)

	err = azureutils.IsSubscriptionInManagementGroup(t, u, v["subscription_management_group_id"].(string))
	assert.NoErrorf(t, err, "subscription %s is not in management group %s", sid, v["subscription_management_group_id"].(string))

	// removed as azurerm_management_group_subscription_association handles this for us
	// tid := os.Getenv("AZURE_TENANT_ID")
	// if err := setSubscriptionManagementGroup(u, tid); err != nil {
	// 	t.Logf("could not move subscription to management group %s: %s", tid, err)
	// }
}

// getValidInputVariables returns a set of valid input variables that can be used and modified for testing scenarios.
func getValidInputVariables(billingScope string) (map[string]interface{}, error) {
	r, err := utils.RandomHex(4)
	if err != nil {
		return nil, fmt.Errorf("cannot generate random hex, %s", err)
	}
	name := fmt.Sprintf("testdeploy-%s", r)
	return map[string]interface{}{
		"subscription_alias_name":    name,
		"subscription_display_name":  name,
		"subscription_billing_scope": billingScope,
		"subscription_workload":      "DevTest",
		"subscription_alias_enabled": true,
	}, nil
}
