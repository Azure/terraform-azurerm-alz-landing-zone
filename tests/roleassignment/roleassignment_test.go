package roleassignment

import (
	"testing"

	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../modules/roleassignment"
)

func TestRoleAssignmentValidWithRoleName(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	defer cleanup()
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	terraformOptions.Vars = v
	// Create plan and ensure only one resource is created.

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equal(t, 1, len(plan.ResourcePlannedValuesMap))
	require.Contains(t, plan.ResourcePlannedValuesMap, "azurerm_role_assignment.this")
	ra := plan.ResourcePlannedValuesMap["azurerm_role_assignment.this"]
	assert.Equalf(t, v["role_assignment_definition"], ra.AttributeValues["role_definition_name"], "role_definition_name incorrect")
	assert.Nilf(t, ra.AttributeValues["role_definition_id"], "role_definition_id should be nil")
	assert.Equalf(t, v["role_assignment_principal_id"], ra.AttributeValues["principal_id"], "role_definition_principal_id incorrect")
}

func TestRoleAssignmentValidWithRoleDefId(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	defer cleanup()
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["role_assignment_definition"] = "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/00000000-0000-0000-0000-000000000000"
	terraformOptions.Vars = v
	// Create plan and ensure only one resource is created.

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	require.Equal(t, 1, len(plan.ResourcePlannedValuesMap))
	require.Contains(t, plan.ResourcePlannedValuesMap, "azurerm_role_assignment.this")
	ra := plan.ResourcePlannedValuesMap["azurerm_role_assignment.this"]
	assert.Equalf(t, v["role_assignment_definition"], ra.AttributeValues["role_definition_id"], "role_definition_id incorrect")
	assert.Nilf(t, ra.AttributeValues["role_definition_name"], "role_definition_name should be nil")
	assert.Equalf(t, v["role_assignment_principal_id"], ra.AttributeValues["principal_id"], "role_definition_principal_id incorrect")
}

func getMockInputVariables() map[string]interface{} {
	return map[string]interface{}{
		"role_assignment_principal_id": "00000000-0000-0000-0000-000000000000",
		"role_assignment_scope":        "/subscriptions/00000000-0000-0000-0000-000000000000",
		"role_assignment_definition":   "Owner",
	}
}
