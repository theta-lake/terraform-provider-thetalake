package supervisionspace_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccSupervisionSpaceResource_basic(t *testing.T) {
	resourceName := "thetalake_supervision_space.test_space"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.TestProviderConfig + `
data "thetalake_role" "reviewer" {
  name = "Reviewer"
}

data "thetalake_role" "api_only" {
  name = "API Only"
}

data "thetalake_integration" "test_integration" {
  name = "Jacob's test"
}

data "thetalake_integration" "test_integration2" {
  name = "Ryans Test Miro"
}

data "thetalake_retention_library" "test_retention_library" {
  name = "API Test Bucket"
}

data "thetalake_directory_group" "test_directory_group" {
  name = "QA Chat"
}

data "thetalake_user_group" "test_user_group" {
  name = "ytesty"
}

resource "thetalake_user" "jacob_test" {
  name = "Jacob TF Test Working! Demo!"
  email = "jacob+tf-acceptance-test-1@thetalake.com"
  password = "Testtesttest123"
  disabled = false
  role_id = data.thetalake_role.reviewer.id
}

resource "thetalake_user" "tf_test" {
  name = "Jacob TF Test Working!"
  email = "jacob+tf-acceptance-test-2@thetalake.com"
  password = "Testtesttest123"
  disabled = false
  role_id = data.thetalake_role.api_only.id
}

resource "thetalake_supervision_space" "test_space" {
  all_participants = false
  all_users = false
  name = "Terraform acceptance test space"
  description = "Supervision space used to test terraform provider development"
  directory_group_ids = [data.thetalake_directory_group.test_directory_group.id]
  external_id = "space-accept-test-id"
  hard_enforce = false
  integration_ids = [data.thetalake_integration.test_integration.id]
  media_types = ["chat", "email"]
  retention_library_ids = [data.thetalake_retention_library.test_retention_library.id]
  requested_supervision_space_priority = 100
  user_group_ids = [data.thetalake_user_group.test_user_group.id]
  user_ids = [resource.thetalake_user.jacob_test.id]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Required attributes
					resource.TestCheckResourceAttr(resourceName, "all_participants", "false"),
					resource.TestCheckResourceAttr(resourceName, "all_users", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test space"),
					resource.TestCheckResourceAttr(resourceName, "description", "Supervision space used to test terraform provider development"),
					resource.TestCheckResourceAttr(resourceName, "external_id", "space-accept-test-id"),
					resource.TestCheckResourceAttr(resourceName, "hard_enforce", "false"),
					resource.TestCheckResourceAttr(resourceName, "requested_supervision_space_priority", "100"),

					// List attributes: ensure expected counts
					resource.TestCheckResourceAttr(resourceName, "directory_group_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "integration_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "media_types.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "retention_library_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_group_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_ids.#", "1"),

					// Specific media_types values
					resource.TestCheckResourceAttr(resourceName, "media_types.0", "chat"),
					resource.TestCheckResourceAttr(resourceName, "media_types.1", "email"),

					// Computed attributes should be set
					resource.TestCheckResourceAttr(resourceName, "media_types.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "assigned_supervision_space_priority"),
					resource.TestCheckResourceAttrSet(resourceName, "can_delete"),
					resource.TestCheckResourceAttrSet(resourceName, "can_enable_all_participants"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "disabled"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
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
  all_users = true
  name = "Terraform test space update"
  description = "Supervision space used to test terraform provider development update"
  directory_group_ids = []
  external_id = "space-accept-test-id-update"
  hard_enforce = true
  integration_ids = []
  media_types = []
  retention_library_ids = []
  requested_supervision_space_priority = 200
  user_group_ids = []
  user_ids = []
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Required attributes updated
					resource.TestCheckResourceAttr(resourceName, "all_participants", "false"),
					resource.TestCheckResourceAttr(resourceName, "all_users", "true"),
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform test space update"),
					resource.TestCheckResourceAttr(resourceName, "description", "Supervision space used to test terraform provider development update"),
					resource.TestCheckResourceAttr(resourceName, "external_id", "space-accept-test-id-update"),
					resource.TestCheckResourceAttr(resourceName, "hard_enforce", "true"),
					resource.TestCheckResourceAttr(resourceName, "requested_supervision_space_priority", "200"),

					// Lists should now be empty
					resource.TestCheckResourceAttr(resourceName, "directory_group_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "integration_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "media_types.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "retention_library_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "user_group_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "user_ids.#", "0"),

					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "assigned_supervision_space_priority"),
					resource.TestCheckResourceAttrSet(resourceName, "can_delete"),
					resource.TestCheckResourceAttrSet(resourceName, "can_enable_all_participants"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "disabled"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
