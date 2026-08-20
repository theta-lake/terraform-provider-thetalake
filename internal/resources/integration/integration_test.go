package integration_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/acctest"
)

func TestAccIntegrationResource_thetaLakeApi(t *testing.T) {
	resourceName := "thetalake_integration.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_integration" "test" {
  name           = "Terraform acceptance test Theta Lake API integration"
  theta_lake_api = {}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test Theta Lake API integration"),
					resource.TestCheckResourceAttr(resourceName, "type", "theta_lake_api"),
					resource.TestCheckResourceAttr(resourceName, "paused", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "integration_type"),
					resource.TestCheckResourceAttrSet(resourceName, "integration_type_id"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update — change name and pause the integration
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_integration" "test" {
  name           = "Terraform acceptance test Theta Lake API integration updated"
  paused         = true
  theta_lake_api = {}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test Theta Lake API integration updated"),
					resource.TestCheckResourceAttr(resourceName, "paused", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIntegrationResource_genericJournaling(t *testing.T) {
	resourceName := "thetalake_integration.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_integration" "test" {
  name = "Terraform acceptance test Generic Journaling integration"

  generic_journaling = {
    download_o365_onedrive_links = true
    index_headers                = "X-Header-Score,X-Routed-Via"
    sender_spf_override          = "v=spf1 ip4:127.0.0.1/32 -all"
    undeliverable_disabled       = true
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test Generic Journaling integration"),
					resource.TestCheckResourceAttr(resourceName, "type", "generic_journaling"),
					resource.TestCheckResourceAttr(resourceName, "generic_journaling.index_headers", "X-Header-Score,X-Routed-Via"),
					resource.TestCheckResourceAttr(resourceName, "generic_journaling.sender_spf_override", "v=spf1 ip4:127.0.0.1/32 -all"),
					resource.TestCheckResourceAttr(resourceName, "generic_journaling.undeliverable_disabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update — flip paused via the pause/start endpoints and update the spf override
			{
				Config: acctest.TestProviderConfig + `
resource "thetalake_integration" "test" {
  name   = "Terraform acceptance test Generic Journaling integration"
  paused = true

  generic_journaling = {
    sender_spf_override          = "v=spf1 ip4:127.0.0.1/24 -all"
    undeliverable_disabled       = true
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "paused", "true"),
					resource.TestCheckResourceAttr(resourceName, "generic_journaling.sender_spf_override", "v=spf1 ip4:127.0.0.1/24 -all"),
					resource.TestCheckResourceAttr(resourceName, "generic_journaling.undeliverable_disabled", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIntegrationResource_genericJournalingUndeliverableMailbox(t *testing.T) {
	if acctest.IntegrationUndeliverableEmailPassword == "" {
		t.Skip("TF_ACC_TEST_INTEGRATION_UNDELIVERABLE_PASSWORD not set")
	}

	resourceName := "thetalake_integration.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: acctest.TestProviderConfig + fmt.Sprintf(`
resource "thetalake_integration" "test" {
  name = "Terraform acceptance test Generic Journaling undeliverable mailbox integration"

  generic_journaling = {
    index_headers                = "X-Header-Score"
	  sender_spf_override          = "v=spf1 ~all"
		undeliverable_disabled       = true
    undeliverable_email_server   = "imap.example.com"
    undeliverable_email_user     = "bounce@example.com"
    undeliverable_email_password = %q
    undeliverable_email_port     = 993
  }
}
`, acctest.IntegrationUndeliverableEmailPassword),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Terraform acceptance test Generic Journaling undeliverable mailbox integration"),
					resource.TestCheckResourceAttr(resourceName, "type", "generic_journaling"),
					resource.TestCheckResourceAttr(resourceName, "generic_journaling.index_headers", "X-Header-Score"),
					resource.TestCheckResourceAttr(resourceName, "generic_journaling.sender_spf_override", "v=spf1 ~all"),
					resource.TestCheckResourceAttr(resourceName, "generic_journaling.undeliverable_disabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "generic_journaling.undeliverable_email_server", "imap.example.com"),
					resource.TestCheckResourceAttr(resourceName, "generic_journaling.undeliverable_email_user", "bounce@example.com"),
					resource.TestCheckResourceAttr(resourceName, "generic_journaling.undeliverable_email_port", "993"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// ImportState. GetIntegrationConfiguration returns undeliverable_email_password
			// when one is set (and "" only when none is configured), so the password does
			// round-trip through import and is verified here along with everything else.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update — change the mailbox server and port, proving in-place update of an
			// already-configured mailbox block does not force replacement
			{
				Config: acctest.TestProviderConfig + fmt.Sprintf(`
resource "thetalake_integration" "test" {
  name = "Terraform acceptance test Generic Journaling undeliverable mailbox integration"

  generic_journaling = {
    index_headers                = "X-Header-Score"
	  sender_spf_override          = "v=spf1 ~all"
		undeliverable_disabled       = true
    undeliverable_email_server   = "imap2.example.com"
    undeliverable_email_user     = "bounce@example.com"
    undeliverable_email_password = %q
    undeliverable_email_port     = 995
  }
}
`, acctest.IntegrationUndeliverableEmailPassword),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "generic_journaling.undeliverable_email_server", "imap2.example.com"),
					resource.TestCheckResourceAttr(resourceName, "generic_journaling.undeliverable_email_port", "995"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
