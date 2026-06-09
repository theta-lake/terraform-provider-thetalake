package role

import (
	"context"
	"testing"

	frameworkdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestNewRoleDataSource(t *testing.T) {
	if _, ok := NewRoleDataSource().(*roleDataSource); !ok {
		t.Fatal("expected NewRoleDataSource to return *roleDataSource")
	}
}

func TestRoleDataSourceMetadata(t *testing.T) {
	dataSource := &roleDataSource{}
	resp := &frameworkdatasource.MetadataResponse{}

	dataSource.Metadata(context.Background(), frameworkdatasource.MetadataRequest{ProviderTypeName: "thetalake"}, resp)

	if resp.TypeName != "thetalake_role" {
		t.Fatalf("expected type name thetalake_role, got %q", resp.TypeName)
	}
}

func TestRoleDataSourceConfigure(t *testing.T) {
	dataSource := &roleDataSource{}
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

func TestRoleDataSourceSchema(t *testing.T) {
	dataSource := &roleDataSource{}
	resp := &frameworkdatasource.SchemaResponse{}

	dataSource.Schema(context.Background(), frameworkdatasource.SchemaRequest{}, resp)

	if resp.Schema.MarkdownDescription == "" {
		t.Fatal("expected schema description to be populated")
	}
	if len(resp.Schema.Attributes) != 8 {
		t.Fatalf("expected 8 attributes, got %d", len(resp.Schema.Attributes))
	}

	nameAttr, ok := resp.Schema.Attributes["name"].(datasourceschema.StringAttribute)
	if !ok {
		t.Fatal("expected name to be a string attribute")
	}
	if !nameAttr.Required {
		t.Fatal("expected name to be required")
	}

	idAttr, ok := resp.Schema.Attributes["id"].(datasourceschema.Int64Attribute)
	if !ok {
		t.Fatal("expected id to be an int64 attribute")
	}
	if !idAttr.Computed {
		t.Fatal("expected id to be computed")
	}

	createdAtAttr, ok := resp.Schema.Attributes["created_at"].(datasourceschema.StringAttribute)
	if !ok {
		t.Fatal("expected created_at to be a string attribute")
	}
	if !createdAtAttr.Computed {
		t.Fatal("expected created_at to be computed")
	}

	defaultAttr, ok := resp.Schema.Attributes["default"].(datasourceschema.BoolAttribute)
	if !ok {
		t.Fatal("expected default to be a bool attribute")
	}
	if !defaultAttr.Computed {
		t.Fatal("expected default to be computed")
	}
}