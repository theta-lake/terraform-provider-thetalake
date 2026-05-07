package usergroup_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccUserGroupDataSource_basic(t *testing.T) {
	dataSourceName := "data.thetalake_user_group.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_user_group" "test" {
  name        = "Terraform DS acceptance test user group"
  description = "User group used to test datasource lookup"
  external_id = "ug-ds-accept-test-id"
}

data "thetalake_user_group" "test" {
  name       = thetalake_user_group.test.name
  depends_on = [thetalake_user_group.test]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", "Terraform DS acceptance test user group"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}
