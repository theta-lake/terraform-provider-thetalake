package customlexicon_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccCustomLexiconResource(t *testing.T) {
	resourceName := "thetalake_custom_lexicon.test"

	policyConfig := ""
	policyCheck := resource.TestCheckResourceAttr(resourceName, "policy_ids.#", "0")
	if acctest.CustomLexiconPolicyID != "" {
		policyConfig = `
  policy_ids  = [` + acctest.CustomLexiconPolicyID + `]`
		policyCheck = resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resourceName, "policy_ids.#", "1"),
			resource.TestCheckTypeSetElemAttr(resourceName, "policy_ids.*", acctest.CustomLexiconPolicyID),
		)
	}

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

  rule_scope              = ["chat", "email", "doc"]
  communication_direction = ["inbound", "outbound", "internal"]

  attachments_enabled           = false
  boilerplate_enabled           = false
  chatroom_name_analyzed        = false
  count_proximity_by_characters = false
  disabled                      = false
  email_smart_body              = false
  email_subject_analyzed        = false
  filename_analyzed             = false
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
			// Update — change name/description, add start/end dates and
			// (if configured) policy associations.
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_custom_lexicon" "test" {
  name        = "Terraform acceptance test custom lexicon updated"
  description = "Updated description"
  risk_type   = "risk"
  rules       = ["accept-test-word-1", "accept-test-word-2"]
  start_date  = "2024-01-01"
  end_date    = "2024-12-31"

  rule_scope              = ["chat", "email", "doc"]
  communication_direction = ["inbound", "outbound", "internal"]

  attachments_enabled           = false
  boilerplate_enabled           = false
  chatroom_name_analyzed        = false
  count_proximity_by_characters = false
  disabled                      = false
  email_smart_body              = false
  email_subject_analyzed        = false
  filename_analyzed             = false` + policyConfig + `
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test custom lexicon updated"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "start_date", "2024-01-01"),
					resource.TestCheckResourceAttr(resourceName, "end_date", "2024-12-31"),
					policyCheck,
				),
			},
			// Update — remove policy associations, disable the lexicon.
			// Disabling must be an in-place update, not a replacement, and
			// must carry start_date/end_date forward rather than nulling
			// them out.
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_custom_lexicon" "test" {
  name        = "Terraform acceptance test custom lexicon updated"
  description = "Updated description"
  risk_type   = "risk"
  rules       = ["accept-test-word-1", "accept-test-word-2"]
  start_date  = "2024-01-01"
  end_date    = "2024-12-31"
  disabled    = true
  policy_ids  = []

  rule_scope              = ["chat", "email", "doc"]
  communication_direction = ["inbound", "outbound", "internal"]

  attachments_enabled           = false
  boilerplate_enabled           = false
  chatroom_name_analyzed        = false
  count_proximity_by_characters = false
  email_smart_body              = false
  email_subject_analyzed        = false
  filename_analyzed             = false
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "disabled_at"),
					resource.TestCheckResourceAttr(resourceName, "start_date", "2024-01-01"),
					resource.TestCheckResourceAttr(resourceName, "end_date", "2024-12-31"),
					resource.TestCheckResourceAttr(resourceName, "policy_ids.#", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase. The lexicon
			// is already disabled from the prior step, so Delete is a no-op.
		},
	})
}
