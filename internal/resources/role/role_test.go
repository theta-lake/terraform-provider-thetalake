package role_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccRoleResource_basic(t *testing.T) {
	resourceName := "thetalake_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_role" "test" {
  name        = "Terraform acceptance test role"
  description = "Role used to test Terraform provider"
  permissions = [
    "cases:read",
    "cases:create",
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test role"),
					resource.TestCheckResourceAttr(resourceName, "description", "Role used to test Terraform provider"),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "permissions.0", "cases:read"),
					resource.TestCheckResourceAttr(resourceName, "permissions.1", "cases:create"),
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
  name        = "Terraform acceptance test role updated"
  description = "Updated role description"
  permissions = [
    "cases:read",
    "cases:create",
    "cases:update",
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test role updated"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated role description"),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "3"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
