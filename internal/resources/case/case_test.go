package legalcase_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccCaseResource_basic(t *testing.T) {
	if acctest.CaseManagerUserID == "" {
		t.Skip("TF_ACC_TEST_CASE_MANAGER_USER_ID not set")
	}

	resourceName := "thetalake_case.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_case" "test" {
  name        = "Terraform acceptance test case"
  number      = "TF-ACC-001"
  description = "Case used to test Terraform provider"
  open_date   = "2024-01-15T10:00:00Z"
  visibility  = "PRIVATE"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test case"),
					resource.TestCheckResourceAttr(resourceName, "number", "TF-ACC-001"),
					resource.TestCheckResourceAttr(resourceName, "description", "Case used to test Terraform provider"),
					resource.TestCheckResourceAttr(resourceName, "open_date", "2024-01-15T10:00:00Z"),
					resource.TestCheckResourceAttr(resourceName, "visibility", "PRIVATE"),
					resource.TestCheckResourceAttr(resourceName, "manager_ids.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
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
			// Update — change description/visibility, add a manager
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_case" "test" {
  name        = "Terraform acceptance test case updated"
  number      = "TF-ACC-001"
  description = "Updated description"
  open_date   = "2024-01-15T10:00:00Z"
  visibility  = "PUBLIC"
  manager_ids = [` + acctest.CaseManagerUserID + `]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test case updated"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "visibility", "PUBLIC"),
					resource.TestCheckResourceAttr(resourceName, "manager_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "manager_ids.*", acctest.CaseManagerUserID),
				),
			},
			// Close the case
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_case" "test" {
  name        = "Terraform acceptance test case updated"
  number      = "TF-ACC-001"
  description = "Updated description"
  open_date   = "2024-01-15T10:00:00Z"
  close_date  = "2024-02-01T10:00:00Z"
  visibility  = "PUBLIC"
  manager_ids = [` + acctest.CaseManagerUserID + `]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "CLOSED"),
					resource.TestCheckResourceAttr(resourceName, "close_date", "2024-02-01T10:00:00Z"),
				),
			},
			// Reopen the case by removing close_date
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_case" "test" {
  name        = "Terraform acceptance test case updated"
  number      = "TF-ACC-001"
  description = "Updated description"
  open_date   = "2024-01-15T10:00:00Z"
  visibility  = "PUBLIC"
  manager_ids = [` + acctest.CaseManagerUserID + `]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "OPEN"),
					resource.TestCheckNoResourceAttr(resourceName, "close_date"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
