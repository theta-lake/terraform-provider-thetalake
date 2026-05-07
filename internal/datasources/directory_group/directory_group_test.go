package directorygroup_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccDirectoryGroupDataSource_basic(t *testing.T) {
	dataSourceName := "data.thetalake_directory_group.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_directory_group" "test" {
  name        = "Terraform DS acceptance test directory group"
  description = "Directory group used to test datasource lookup"
  external_id = "dg-ds-accept-test-id"
}

data "thetalake_directory_group" "test" {
  name       = thetalake_directory_group.test.name
  depends_on = [thetalake_directory_group.test]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", "Terraform DS acceptance test directory group"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}
