package workspace_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccWorkspaceResource_basic(t *testing.T) {
	workspaceId := acctest.WorkspaceId
	if workspaceId == "" {
		t.Skip("TF_ACC_TEST_WORKSPACE_ID not set, skipping workspace acceptance test")
	}

	resourceName := "thetalake_workspace.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import an existing workspace
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_workspace" "test" {}
`,
				ResourceName:       resourceName,
				ImportState:        true,
				ImportStateId:      workspaceId,
				ImportStatePersist: true,
			},
			// Refresh state and verify key attributes are populated
			{
				RefreshState: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "default_workspace_timezone"),
					resource.TestCheckResourceAttrSet(resourceName, "default_transcription_language"),
					resource.TestCheckResourceAttrSet(resourceName, "shared_links_expiration_period"),
				),
			},
			// Update: toggle allow_anonymous_via_shared_links and verify
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_workspace" "test" {
  allow_anonymous_via_shared_links = false
  delete_on_expiration             = false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "allow_anonymous_via_shared_links", "false"),
					resource.TestCheckResourceAttr(resourceName, "delete_on_expiration", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
		},
	})
}
