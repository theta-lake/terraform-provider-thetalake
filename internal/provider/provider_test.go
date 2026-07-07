package provider

import (
	"context"
	"testing"

	frameworkdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestNew(t *testing.T) {
	providerFactory := New("1.2.3")
	providerInstance, ok := providerFactory().(*ThetalakeProvider)
	if !ok {
		t.Fatal("expected New to return *ThetalakeProvider")
	}

	if providerInstance.version != "1.2.3" {
		t.Fatalf("expected version 1.2.3, got %q", providerInstance.version)
	}
	if providerInstance.client != nil {
		t.Fatal("expected New to create provider without injected client")
	}
}

func TestNewWithClient(t *testing.T) {
	client := &thetalake.Client{}
	providerFactory := NewWithClient("1.2.3", client)
	providerInstance, ok := providerFactory().(*ThetalakeProvider)
	if !ok {
		t.Fatal("expected NewWithClient to return *ThetalakeProvider")
	}

	if providerInstance.version != "1.2.3" {
		t.Fatalf("expected version 1.2.3, got %q", providerInstance.version)
	}
	if providerInstance.client != client {
		t.Fatal("expected injected client to be preserved")
	}
}

func TestMetadata(t *testing.T) {
	p := &ThetalakeProvider{version: "9.9.9"}
	resp := &frameworkprovider.MetadataResponse{}

	p.Metadata(context.Background(), frameworkprovider.MetadataRequest{}, resp)

	if resp.TypeName != "thetalake" {
		t.Fatalf("expected type name thetalake, got %q", resp.TypeName)
	}
	if resp.Version != "9.9.9" {
		t.Fatalf("expected version 9.9.9, got %q", resp.Version)
	}
}

func TestSchema(t *testing.T) {
	p := &ThetalakeProvider{}
	resp := &frameworkprovider.SchemaResponse{}

	p.Schema(context.Background(), frameworkprovider.SchemaRequest{}, resp)

	if len(resp.Schema.Attributes) != 3 {
		t.Fatalf("expected 3 provider attributes, got %d", len(resp.Schema.Attributes))
	}

	apiServerAttr, ok := resp.Schema.Attributes["api_server"].(providerschema.StringAttribute)
	if !ok {
		t.Fatal("expected api_server to be a string attribute")
	}
	if !apiServerAttr.Required {
		t.Fatal("expected api_server to be required")
	}

	clientIDAttr, ok := resp.Schema.Attributes["client_id"].(providerschema.StringAttribute)
	if !ok {
		t.Fatal("expected client_id to be a string attribute")
	}
	if !clientIDAttr.Required || !clientIDAttr.Sensitive {
		t.Fatal("expected client_id to be required and sensitive")
	}

	clientSecretAttr, ok := resp.Schema.Attributes["client_secret"].(providerschema.StringAttribute)
	if !ok {
		t.Fatal("expected client_secret to be a string attribute")
	}
	if !clientSecretAttr.Required || !clientSecretAttr.Sensitive {
		t.Fatal("expected client_secret to be required and sensitive")
	}
}

func TestConfigureWithInjectedClient(t *testing.T) {
	client := &thetalake.Client{}
	p := &ThetalakeProvider{client: client}
	resp := &frameworkprovider.ConfigureResponse{}

	p.Configure(context.Background(), frameworkprovider.ConfigureRequest{}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatal("expected injected client configure path to succeed without diagnostics")
	}
	if resp.DataSourceData != client {
		t.Fatal("expected injected client to be propagated to data sources")
	}
	if resp.ResourceData != client {
		t.Fatal("expected injected client to be propagated to resources")
	}
}

func TestResources(t *testing.T) {
	p := &ThetalakeProvider{}
	resourceFactories := p.Resources(context.Background())

	expectedTypes := []string{
		"thetalake_directory_group",
		"thetalake_label",
		"thetalake_retention_library",
		"thetalake_role",
		"thetalake_supervision_space",
		"thetalake_user",
		"thetalake_user_group",
		"thetalake_swrv_rule",
		"thetalake_workspace",
	}

	if len(resourceFactories) != len(expectedTypes) {
		t.Fatalf("expected %d resource factories, got %d", len(expectedTypes), len(resourceFactories))
	}

	for index, factory := range resourceFactories {
		resp := &frameworkresource.MetadataResponse{}
		factory().Metadata(context.Background(), frameworkresource.MetadataRequest{ProviderTypeName: "thetalake"}, resp)

		if resp.TypeName != expectedTypes[index] {
			t.Fatalf("expected resource type %q at index %d, got %q", expectedTypes[index], index, resp.TypeName)
		}
	}
}

func TestDataSources(t *testing.T) {
	p := &ThetalakeProvider{}
	dataSourceFactories := p.DataSources(context.Background())

	expectedTypes := []string{
		"thetalake_role",
		"thetalake_integration",
		"thetalake_retention_library",
		"thetalake_directory_group",
		"thetalake_user_group",
		"thetalake_identity",
		"thetalake_user",
		"thetalake_workspace",
	}

	if len(dataSourceFactories) != len(expectedTypes) {
		t.Fatalf("expected %d data source factories, got %d", len(expectedTypes), len(dataSourceFactories))
	}

	for index, factory := range dataSourceFactories {
		resp := &frameworkdatasource.MetadataResponse{}
		factory().Metadata(context.Background(), frameworkdatasource.MetadataRequest{ProviderTypeName: "thetalake"}, resp)

		if resp.TypeName != expectedTypes[index] {
			t.Fatalf("expected data source type %q at index %d, got %q", expectedTypes[index], index, resp.TypeName)
		}
	}
}
