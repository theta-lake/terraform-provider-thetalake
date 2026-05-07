package user_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccUserDataSource_basic(t *testing.T) {
	if acctest.UserEmail == "" {
		t.Skip("TF_ACC_TEST_USER_EMAIL not set")
	}

	dataSourceName := "data.thetalake_user.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.TestProviderConfig + `
data "thetalake_user" "test" {
  email = "` + acctest.UserEmail + `"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "email", acctest.UserEmail),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
				),
			},
		},
	})
}
