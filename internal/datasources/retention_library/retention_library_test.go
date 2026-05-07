package retentionlibrary_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccRetentionLibraryDataSource_basic(t *testing.T) {
	if acctest.RetentionLibName == "" {
		t.Skip("TF_ACC_TEST_RETENTION_LIBRARY_NAME not set")
	}

	dataSourceName := "data.thetalake_retention_library.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.TestProviderConfig + `
data "thetalake_retention_library" "test" {
  name = "` + acctest.RetentionLibName + `"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", acctest.RetentionLibName),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}
