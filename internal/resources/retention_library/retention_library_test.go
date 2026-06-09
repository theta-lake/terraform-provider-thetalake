package retentionlibrary_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccRetentionLibraryResource_basic(t *testing.T) {
	if acctest.RetentionLibStorageAccountID == "" {
		t.Skip("TF_ACC_TEST_RETENTION_LIBRARY_STORAGE_ACCOUNT_ID not set")
	}

	resourceName := "thetalake_retention_library.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_retention_library" "test" {
  name                         = "Terraform acceptance test retention library"
  storage_account_id           = ` + acctest.RetentionLibStorageAccountID + `
  description                  = "Retention library used to test Terraform provider"
  external_id                  = "rl-accept-test-id"
  retain_in_review             = false
  retention_period_days        = 30
  retention_period_enabled     = true
  sec_compliant_storage_enabled = false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test retention library"),
					resource.TestCheckResourceAttr(resourceName, "storage_account_id", acctest.RetentionLibStorageAccountID),
					resource.TestCheckResourceAttr(resourceName, "description", "Retention library used to test Terraform provider"),
					resource.TestCheckResourceAttr(resourceName, "external_id", "rl-accept-test-id"),
					resource.TestCheckResourceAttr(resourceName, "retain_in_review", "false"),
					resource.TestCheckResourceAttr(resourceName, "retention_period_days", "30"),
					resource.TestCheckResourceAttr(resourceName, "retention_period_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "sec_compliant_storage_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			// ImportState
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"retain_in_review"},
			},
			// Update
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_retention_library" "test" {
  name                         = "Terraform acceptance test retention library updated"
  storage_account_id           = ` + acctest.RetentionLibStorageAccountID + `
  description                  = "Updated description"
  external_id                  = "rl-accept-test-id-updated"
  retain_in_review             = true
  retention_period_days        = 90
  retention_period_enabled     = true
  sec_compliant_storage_enabled = false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test retention library updated"),
					resource.TestCheckResourceAttr(resourceName, "storage_account_id", acctest.RetentionLibStorageAccountID),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "external_id", "rl-accept-test-id-updated"),
					resource.TestCheckResourceAttr(resourceName, "retain_in_review", "true"),
					resource.TestCheckResourceAttr(resourceName, "retention_period_days", "90"),
					resource.TestCheckResourceAttr(resourceName, "retention_period_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
