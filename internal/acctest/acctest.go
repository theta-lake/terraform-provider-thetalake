package acctest

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/theta-lake/terraform-provider-thetalake/internal/provider"
)

var (
	clientId           = os.Getenv("TF_ACC_TEST_CLIENT_ID")
	clientSecret       = os.Getenv("TF_ACC_TEST_CLIENT_SECRET")
	apiServer          = os.Getenv("TF_ACC_TEST_API_SERVER")
	TestProviderConfig = fmt.Sprintf(`
provider "thetalake" {
  api_server    = "%s"
  client_id     = "%s"
  client_secret = "%s"
}
`, apiServer, clientId, clientSecret)
)

var (
	ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"thetalake": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
)
