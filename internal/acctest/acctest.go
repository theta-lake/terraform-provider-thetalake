package acctest

import (
	"fmt"
	"os"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
	"github.com/theta-lake/terraform-provider-thetalake/internal/provider"
)

var (
	clientId                     = os.Getenv("TF_ACC_TEST_CLIENT_ID")
	clientSecret                 = os.Getenv("TF_ACC_TEST_CLIENT_SECRET")
	apiServer                    = os.Getenv("TF_ACC_TEST_API_SERVER")
	IntegrationName              = os.Getenv("TF_ACC_TEST_INTEGRATION_NAME")
	RetentionLibName             = os.Getenv("TF_ACC_TEST_RETENTION_LIBRARY_NAME")
	RetentionLibStorageAccountID = os.Getenv("TF_ACC_TEST_RETENTION_LIBRARY_STORAGE_ACCOUNT_ID")
	SwrvRulePolicyID             = os.Getenv("TF_ACC_TEST_SWRV_RULE_POLICY_ID")
	SwrvRuleRetentionLibraryID   = os.Getenv("TF_ACC_TEST_SWRV_RULE_RETENTION_LIBRARY_ID")
	SwrvRuleWorkflowID           = os.Getenv("TF_ACC_TEST_SWRV_RULE_WORKFLOW_ID")
	IdentityEmail                = os.Getenv("TF_ACC_TEST_IDENTITY_EMAIL")
	UserEmail                    = os.Getenv("TF_ACC_TEST_USER_EMAIL")
	TestProviderConfig           = fmt.Sprintf(`
provider "thetalake" {
  api_server    = "%s"
  client_id     = "%s"
  client_secret = "%s"
}
`, apiServer, clientId, clientSecret)
)

// sharedClient is a singleton authenticated client shared across all acceptance
// tests in a package. It is initialised once so that the token endpoint is only
// called once per test binary, regardless of how many test steps run.
var (
	sharedClient     *thetalake.Client
	sharedClientOnce sync.Once
)

func getSharedClient() *thetalake.Client {
	sharedClientOnce.Do(func() {
		sharedClient = thetalake.NewClient(apiServer, clientId, clientSecret)
	})
	return sharedClient
}

var (
	ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"thetalake": providerserver.NewProtocol6WithError(provider.NewWithClient("test", getSharedClient())()),
	}
)
