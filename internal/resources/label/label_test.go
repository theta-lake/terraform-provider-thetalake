package label_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccLabelResource_basic(t *testing.T) {
	resourceName := "thetalake_label.test"
	suffix := time.Now().UnixNano() % 1000000
	shortName := fmt.Sprintf("TF%06d", suffix)
	longName := fmt.Sprintf("Terraform test label %06d", suffix)
	updatedLongName := fmt.Sprintf("Terraform test label updated %06d", suffix)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_label" "test" {
  short_name       = "` + shortName + `"
  long_name        = "` + longName + `"
  background_color = "#03A71E"
  hidden           = false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "short_name", shortName),
					resource.TestCheckResourceAttr(resourceName, "long_name", longName),
					resource.TestCheckResourceAttr(resourceName, "background_color", "#03A71E"),
					resource.TestCheckResourceAttr(resourceName, "hidden", "false"),
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
			// Update
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_label" "test" {
  short_name       = "` + shortName + `"
  long_name        = "` + updatedLongName + `"
  background_color = "#03A71E"
  hidden           = false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "short_name", shortName),
					resource.TestCheckResourceAttr(resourceName, "long_name", updatedLongName),
					resource.TestCheckResourceAttr(resourceName, "background_color", "#03A71E"),
					resource.TestCheckResourceAttr(resourceName, "hidden", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
