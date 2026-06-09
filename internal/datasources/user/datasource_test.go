package user

import (
	"context"
	"testing"

	frameworkdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestNewUserDataSource(t *testing.T) {
	if _, ok := NewUserDataSource().(*userDataSource); !ok {
		t.Fatal("expected NewUserDataSource to return *userDataSource")
	}
}

func TestUserDataSourceMetadata(t *testing.T) {
	dataSource := &userDataSource{}
	resp := &frameworkdatasource.MetadataResponse{}

	dataSource.Metadata(context.Background(), frameworkdatasource.MetadataRequest{ProviderTypeName: "thetalake"}, resp)

	if resp.TypeName != "thetalake_user" {
		t.Fatalf("expected type name thetalake_user, got %q", resp.TypeName)
	}
}

func TestUserDataSourceConfigure(t *testing.T) {
	dataSource := &userDataSource{}
	resp := &frameworkdatasource.ConfigureResponse{}
	client := &thetalake.Client{}

	dataSource.Configure(context.Background(), frameworkdatasource.ConfigureRequest{}, resp)

	if dataSource.client != (thetalake.Client{}) {
		t.Fatal("expected client to remain zero value when provider data is absent")
	}

	dataSource.Configure(context.Background(), frameworkdatasource.ConfigureRequest{ProviderData: client}, resp)

	if dataSource.client != *client {
		t.Fatal("expected client to be copied from provider data")
	}
}

func TestUserDataSourceSchema(t *testing.T) {
	dataSource := &userDataSource{}
	resp := &frameworkdatasource.SchemaResponse{}

	dataSource.Schema(context.Background(), frameworkdatasource.SchemaRequest{}, resp)

	if resp.Schema.MarkdownDescription == "" {
		t.Fatal("expected schema description to be populated")
	}
	if len(resp.Schema.Attributes) != 3 {
		t.Fatalf("expected 3 attributes, got %d", len(resp.Schema.Attributes))
	}

	emailAttr, ok := resp.Schema.Attributes["email"].(datasourceschema.StringAttribute)
	if !ok {
		t.Fatal("expected email to be a string attribute")
	}
	if !emailAttr.Required {
		t.Fatal("expected email to be required")
	}

	idAttr, ok := resp.Schema.Attributes["id"].(datasourceschema.Int64Attribute)
	if !ok {
		t.Fatal("expected id to be an int64 attribute")
	}
	if !idAttr.Computed {
		t.Fatal("expected id to be computed")
	}

	nameAttr, ok := resp.Schema.Attributes["name"].(datasourceschema.StringAttribute)
	if !ok {
		t.Fatal("expected name to be a string attribute")
	}
	if !nameAttr.Computed {
		t.Fatal("expected name to be computed")
	}
}
