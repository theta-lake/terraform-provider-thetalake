package identity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccIdentityDataSource_basic(t *testing.T) {
	if acctest.IdentityEmail == "" {
		t.Skip("TF_ACC_TEST_IDENTITY_EMAIL not set")
	}

	dataSourceName := "data.thetalake_identity.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.TestProviderConfig + `
data "thetalake_identity" "test" {
  email = "` + acctest.IdentityEmail + `"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "email", acctest.IdentityEmail),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
				),
			},
		},
	})
}
