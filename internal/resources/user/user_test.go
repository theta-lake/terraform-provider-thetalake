package user_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
  email    = "test-user@example.com"
  password = "example-password"
  role_id  = data.thetalake_role.reviewer.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test User"),
					resource.TestCheckResourceAttr(resourceName, "email", "test-user@example.com"),
					resource.TestCheckResourceAttr(resourceName, "role", "Reviewer"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_user" "test" {
  name      = "Updated User"
  email     = "updated-user@example.com"
  password  = "example-password"
  role      = "Admin"
  search_id = 42
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Updated User"),
					resource.TestCheckResourceAttr(resourceName, "email", "updated-user@example.com"),
					resource.TestCheckResourceAttr(resourceName, "role_id", "2"),
					resource.TestCheckResourceAttr(resourceName, "search_id", "42"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
