package supervisionspace_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccUserResource_basic(t *testing.T) {
	resourceName := "thetalake_supervision_space.test_space"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_supervision_space" "test_space" {
  all_participants = false
  all_users = false
  name = "Terraform test space"
  description = "Supervision space used to test terraform provider development"
  directory_group_ids = []
  external_id = "space-accept-test-id"
  hard_enforce = false
  integration_ids = []
  media_type_ids = []
  retention_library_ids = []
  requested_supervision_space_priority = 100
  user_group_ids = []
  user_ids = []
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform test space"),
					resource.TestCheckResourceAttr(resourceName, "description", "Supervision space used to test terraform provider development"),
					resource.TestCheckResourceAttr(resourceName, "external_id", "space-accept-test-id"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"requested_supervision_space_priority"},
			},
			// Update and Read testing
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_supervision_space" "test_space" {
  all_participants = false
  all_users = false
  name = "Terraform test space update"
  description = "Supervision space used to test terraform provider development update"
  directory_group_ids = []
  external_id = "space-accept-test-id-update"
  hard_enforce = false
  integration_ids = []
  media_type_ids = []
  retention_library_ids = []
  requested_supervision_space_priority = 100
  user_group_ids = []
  user_ids = []
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform test space update"),
					resource.TestCheckResourceAttr(resourceName, "description", "Supervision space used to test terraform provider development update"),
					resource.TestCheckResourceAttr(resourceName, "external_id", "space-accept-test-id-update"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
