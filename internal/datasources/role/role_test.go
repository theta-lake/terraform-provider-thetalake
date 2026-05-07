package role_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccRoleDataSource_basic(t *testing.T) {
	dataSourceName := "data.thetalake_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.TestProviderConfig + `
data "thetalake_role" "test" {
  name = "Reviewer"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", "Reviewer"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "is_built_in"),
				),
			},
		},
	})
}
