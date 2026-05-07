package directorygroup_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccDirectoryGroupResource_basic(t *testing.T) {
	resourceName := "thetalake_directory_group.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_directory_group" "test" {
  name        = "Terraform acceptance test directory group"
  description = "Directory group used to test Terraform provider"
  external_id = "dg-accept-test-id"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test directory group"),
					resource.TestCheckResourceAttr(resourceName, "description", "Directory group used to test Terraform provider"),
					resource.TestCheckResourceAttr(resourceName, "external_id", "dg-accept-test-id"),
					resource.TestCheckResourceAttr(resourceName, "identity_ids.#", "0"),
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
			// Update — change name/description, omit optional fields (should not wipe them)
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_directory_group" "test" {
  name        = "Terraform acceptance test directory group updated"
  description = "Updated description"
  external_id = "dg-accept-test-id-updated"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test directory group updated"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "external_id", "dg-accept-test-id-updated"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
