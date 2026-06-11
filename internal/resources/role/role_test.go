package role_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccRoleResource_basic(t *testing.T) {
	if acctest.ApiServer == "" || acctest.ClientId == "" || acctest.ClientSecret == "" {
		t.Skip("TF_ACC_TEST_API_SERVER, TF_ACC_TEST_CLIENT_ID, and TF_ACC_TEST_CLIENT_SECRET must be set")
	}

	resourceName := "thetalake_role.test"
	roleName := fmt.Sprintf("Terraform acceptance test role %d", time.Now().UnixNano())
	updatedRoleName := roleName + " updated"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_role" "test" {
	name        = "` + roleName + `"
  description = "Role used to test Terraform provider"
  permissions = [
    "cases:read",
    "cases:create",
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", roleName),
					resource.TestCheckResourceAttr(resourceName, "description", "Role used to test Terraform provider"),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "is_built_in", "false"),
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
			// Update — change name, description, and permissions
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_role" "test" {
	name        = "` + updatedRoleName + `"
  description = "Updated role description"
  permissions = [
    "cases:read",
    "cases:create",
    "cases:update",
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedRoleName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated role description"),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "3"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
