package usergroup_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccUserGroupResource_basic(t *testing.T) {
	resourceName := "thetalake_user_group.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_user_group" "test" {
  name        = "Terraform acceptance test user group"
  description = "User group used to test Terraform provider"
  external_id = "ug-accept-test-id"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test user group"),
					resource.TestCheckResourceAttr(resourceName, "description", "User group used to test Terraform provider"),
					resource.TestCheckResourceAttr(resourceName, "external_id", "ug-accept-test-id"),
					resource.TestCheckResourceAttr(resourceName, "user_ids.#", "0"),
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
			// Update — change name/description/external_id
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_user_group" "test" {
  name        = "Terraform acceptance test user group updated"
  description = "Updated description"
  external_id = "ug-accept-test-id-updated"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test user group updated"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "external_id", "ug-accept-test-id-updated"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
