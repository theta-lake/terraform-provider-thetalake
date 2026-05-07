package integration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccIntegrationDataSource_basic(t *testing.T) {
	if acctest.IntegrationName == "" {
		t.Skip("TF_ACC_TEST_INTEGRATION_NAME not set")
	}

	dataSourceName := "data.thetalake_integration.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.TestProviderConfig + `
data "thetalake_integration" "test" {
  name = "` + acctest.IntegrationName + `"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", acctest.IntegrationName),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}
