package swrvrule_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccSwrvRuleResource_basic(t *testing.T) {
	if acctest.SwrvRulePolicyID == "" {
		t.Skip("TF_ACC_TEST_SWRV_RULE_POLICY_ID not set")
	}
	if acctest.SwrvRuleRetentionLibraryID == "" {
		t.Skip("TF_ACC_TEST_SWRV_RULE_RETENTION_LIBRARY_ID not set")
	}
	if acctest.SwrvRuleWorkflowID == "" {
		t.Skip("TF_ACC_TEST_SWRV_RULE_WORKFLOW_ID not set")
	}

	resourceName := "thetalake_swrv_rule.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_swrv_rule" "test" {
  name                 = "Terraform acceptance test SWRV rule"
  description          = "SWRV rule used to test Terraform provider"
  policy_id            = ` + acctest.SwrvRulePolicyID + `
  retention_library_id = ` + acctest.SwrvRuleRetentionLibraryID + `
  workflow_id          = ` + acctest.SwrvRuleWorkflowID + `
  priority             = 1
  input_sources = [{
	type = "all_uploads"
  }]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test SWRV rule"),
					resource.TestCheckResourceAttr(resourceName, "description", "SWRV rule used to test Terraform provider"),
					resource.TestCheckResourceAttr(resourceName, "policy_id", acctest.SwrvRulePolicyID),
					resource.TestCheckResourceAttr(resourceName, "retention_library_id", acctest.SwrvRuleRetentionLibraryID),
					resource.TestCheckResourceAttr(resourceName, "workflow_id", acctest.SwrvRuleWorkflowID),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "input_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "input_sources.0.type", "all_uploads"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"input_sources",
				},
			},
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_swrv_rule" "test" {
  name                 = "Terraform acceptance test SWRV rule updated"
  description          = "Updated SWRV rule description"
  policy_id            = ` + acctest.SwrvRulePolicyID + `
  retention_library_id = ` + acctest.SwrvRuleRetentionLibraryID + `
  workflow_id          = ` + acctest.SwrvRuleWorkflowID + `
  priority             = 0
  input_sources = [{
	type = "all_uploads"
  }]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test SWRV rule updated"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated SWRV rule description"),
					resource.TestCheckResourceAttr(resourceName, "priority", "0"),
					resource.TestCheckResourceAttr(resourceName, "input_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "input_sources.0.type", "all_uploads"),
				),
			},
		},
	})
}
