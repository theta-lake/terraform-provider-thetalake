package customlexicon_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

// customLexiconConfig renders the acceptance test lexicon with name,
// description and any extra attributes the step under test needs, so each step
// only has to spell out what it is changing.
func customLexiconConfig(name string, description string, extra string) string {
	return acctest.TestProviderConfig + `
resource "thetalake_custom_lexicon" "test" {
  name        = "` + name + `"
  description = "` + description + `"
  risk_type   = "risk"
  rules       = ["accept-test-word-1", "accept-test-word-2"]

  rule_scope              = ["chat", "email", "doc"]
  communication_direction = ["inbound", "outbound", "internal"]

  attachments_enabled           = false
  boilerplate_enabled           = false
  chatroom_name_analyzed        = false
  count_proximity_by_characters = false
  email_smart_body              = false
  email_subject_analyzed        = false
  filename_analyzed             = false
` + extra + `
}
`
}

func TestAccCustomLexiconResource(t *testing.T) {
	resourceName := "thetalake_custom_lexicon.test"
	const updatedName = "Terraform acceptance test custom lexicon updated"

	policyConfig := ""
	policyCheck := resource.TestCheckNoResourceAttr(resourceName, "policy_ids.#")
	if acctest.CustomLexiconPolicyID != "" {
		policyConfig = `
  policy_ids = [` + acctest.CustomLexiconPolicyID + `]`
		policyCheck = resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resourceName, "policy_ids.#", "1"),
			resource.TestCheckTypeSetElemAttr(resourceName, "policy_ids.*", acctest.CustomLexiconPolicyID),
		)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read. Nothing optional is configured, so the optional
			// attributes must all read back as null rather than as whatever the
			// API defaults them to.
			{
				Config: customLexiconConfig(
					"Terraform acceptance test custom lexicon",
					"Custom lexicon used to test Terraform provider",
					`  disabled = false`,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test custom lexicon"),
					resource.TestCheckResourceAttr(resourceName, "description", "Custom lexicon used to test Terraform provider"),
					resource.TestCheckResourceAttr(resourceName, "risk_type", "risk"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "rules.*", "accept-test-word-1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "rules.*", "accept-test-word-2"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckNoResourceAttr(resourceName, "start_date"),
					resource.TestCheckNoResourceAttr(resourceName, "end_date"),
					resource.TestCheckNoResourceAttr(resourceName, "max_participants"),
					resource.TestCheckNoResourceAttr(resourceName, "min_num_rules_with_hits"),
					resource.TestCheckNoResourceAttr(resourceName, "policy_ids.#"),
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
				Config: customLexiconConfig(updatedName, "Updated description", `  disabled   = false
  start_date = "2024-01-01"
  end_date   = "2024-12-31"`+policyConfig),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "start_date", "2024-01-01"),
					resource.TestCheckResourceAttr(resourceName, "end_date", "2024-12-31"),
					policyCheck,
				),
			},
			// Update — remove start_date, end_date and policy_ids from the
			// configuration. Those attributes are optional but not computed, so
			// dropping them must plan a real change and clear the values rather
			// than silently re-using the prior state.
			{
				Config: customLexiconConfig(updatedName, "Updated description", `  disabled = false`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "start_date"),
					resource.TestCheckNoResourceAttr(resourceName, "end_date"),
					resource.TestCheckNoResourceAttr(resourceName, "policy_ids.#"),
				),
			},
			// Add a create-only optional attribute: this must force
			// replacement, since the API cannot change it in place.
			{
				Config: customLexiconConfig(updatedName, "Updated description", `  disabled         = false
  max_participants = 5`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.TestCheckResourceAttr(resourceName, "max_participants", "5"),
			},
			// Remove that create-only attribute again. Going back to "no limit"
			// is also a change the API cannot make in place, so it must force
			// replacement instead of being planned as a no-op.
			{
				Config: customLexiconConfig(updatedName, "Updated description", `  disabled = false`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.TestCheckNoResourceAttr(resourceName, "max_participants"),
			},
			// Update — disable the lexicon. Disabling must be an in-place
			// update, not a replacement, and must carry start_date/end_date
			// forward rather than nulling them out.
			{
				Config: customLexiconConfig(updatedName, "Updated description", `  disabled   = true
  start_date = "2024-01-01"
  end_date   = "2024-12-31"`),
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
					resource.TestCheckNoResourceAttr(resourceName, "policy_ids.#"),
				),
			},
			// Delete testing automatically occurs in TestCase. The lexicon
			// is already disabled from the prior step, so Delete is a no-op.
		},
	})
}
