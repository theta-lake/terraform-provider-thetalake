package customlexicon_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccCustomLexiconResource_basic(t *testing.T) {
	resourceName := "thetalake_custom_lexicon.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_custom_lexicon" "test" {
  name        = "Terraform acceptance test custom lexicon"
  description = "Custom lexicon used to test Terraform provider"
  risk_type   = "risk"
  rules       = ["accept-test-word-1", "accept-test-word-2"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test custom lexicon"),
					resource.TestCheckResourceAttr(resourceName, "description", "Custom lexicon used to test Terraform provider"),
					resource.TestCheckResourceAttr(resourceName, "risk_type", "risk"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "rules.*", "accept-test-word-1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "rules.*", "accept-test-word-2"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update — change name/description, toggle disabled
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_custom_lexicon" "test" {
  name        = "Terraform acceptance test custom lexicon updated"
  description = "Updated description"
  risk_type   = "risk"
  rules       = ["accept-test-word-1", "accept-test-word-2"]
  disabled    = true
}
`,
				// Disabling must be an in-place update, not a replacement — the
				// create-only attributes are unconfigured here, so they must not
				// trigger RequiresReplace.
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test custom lexicon updated"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "disabled_at"),
				),
			},
			// Delete testing automatically occurs in TestCase. The lexicon is
			// already disabled from the prior step, so Delete is a no-op.
		},
	})
}

// TestAccCustomLexiconResource_disabledOnCreate covers creating a lexicon
// with disabled = true from the start. Create issues a follow-up disable
// call (the create endpoint cannot set disabled=true directly) that must
// carry start_date/end_date forward, or it would null them out.
func TestAccCustomLexiconResource_disabledOnCreate(t *testing.T) {
	resourceName := "thetalake_custom_lexicon.test_disabled"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_custom_lexicon" "test_disabled" {
  name        = "Terraform acceptance test disabled-on-create custom lexicon"
  description = "Custom lexicon created disabled to test Terraform provider"
  risk_type   = "risk"
  rules       = ["accept-test-disabled-word"]
  disabled    = true
  start_date  = "2024-01-01"
  end_date    = "2024-12-31"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "disabled_at"),
					resource.TestCheckResourceAttr(resourceName, "start_date", "2024-01-01"),
					resource.TestCheckResourceAttr(resourceName, "end_date", "2024-12-31"),
				),
			},
			// Delete testing automatically occurs in TestCase. The lexicon is
			// already disabled, so Delete is a no-op.
		},
	})
}

func TestAccCustomLexiconResource_withPolicies(t *testing.T) {
	if acctest.CustomLexiconPolicyID == "" {
		t.Skip("TF_ACC_TEST_CUSTOM_LEXICON_POLICY_ID not set")
	}

	resourceName := "thetalake_custom_lexicon.test_policies"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_custom_lexicon" "test_policies" {
  name        = "Terraform acceptance test custom lexicon with policies"
  description = "Custom lexicon used to test Terraform provider policy associations"
  risk_type   = "risk"
  rules       = ["accept-test-policy-word"]
  policy_ids  = [` + acctest.CustomLexiconPolicyID + `]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "policy_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "policy_ids.*", acctest.CustomLexiconPolicyID),
				),
			},
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_custom_lexicon" "test_policies" {
  name        = "Terraform acceptance test custom lexicon with policies"
  description = "Custom lexicon used to test Terraform provider policy associations"
  risk_type   = "risk"
  rules       = ["accept-test-policy-word"]
  policy_ids  = []
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "policy_ids.#", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
