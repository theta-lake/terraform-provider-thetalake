package user_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccUserResource_basic(t *testing.T) {
	resourceName := "thetalake_user.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.TestProviderConfig + `
data "thetalake_role" "reviewer" {
  name = "Reviewer"
}

resource "thetalake_user" "test" {
  name     = "Test User"
  email    = "jacob+tf-acctest@thetalake.com"
  password = "example-password"
  role_id  = data.thetalake_role.reviewer.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test User"),
					resource.TestCheckResourceAttr(resourceName, "email", "jacob+tf-acctest@thetalake.com"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Update and Read testing
			{
				Config: acctest.TestProviderConfig + `
data "thetalake_role" "reviewer" {
  name = "Reviewer"
}

resource "thetalake_user" "test" {
  name      = "Updated User"
  email     = "jacob+tf-acctest-updated@thetalake.com"
  password  = "example-password"
  role_id  = data.thetalake_role.reviewer.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Updated User"),
					resource.TestCheckResourceAttr(resourceName, "email", "jacob+tf-acctest-updated@thetalake.com"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
