package acctest

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/theta-lake/terraform-provider-thetalake/internal/provider"
)

const (
	TestProviderConfig = `
provider "thetalake" {
  api_server    = "https://api-dev1.thetalake.com"
  client_id     = "6YET9xip3t0H2fRmVI5Mlo5kP3lBeIQG"
  client_secret = "bluBQtLlLBjeaDE1kDxc1qYF268ibTUzRczlKjVw3L2VMcfA6lrSN6H7CpICjVQJkPNgCR8I55k3H-ieBvqgT042Hn0fA4tFXRHfLNsDFsZbsOEcyzb1zYr2s-9A-jkJ"
}
`
)

var (
	ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"thetalake": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
)
