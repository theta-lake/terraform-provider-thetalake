package acctest

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/theta-lake/terraform-provider-thetalake/internal/provider"
)

const (
	TestProviderConfig = `
provider "thetalake" {
  endpoint = "https://api.thetalake.com"
  client_id = "test123"
  client_secret = "secret123"
}
`
)

var (
	ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"thetalake": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
)
