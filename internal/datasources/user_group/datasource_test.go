package usergroup

import (
	"context"
	"testing"

	frameworkdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestNewUserGroupDataSource(t *testing.T) {
	if _, ok := NewUserGroupDataSource().(*userGroupDataSource); !ok {
		t.Fatal("expected NewUserGroupDataSource to return *userGroupDataSource")
	}
}

func TestUserGroupDataSourceMetadata(t *testing.T) {
	dataSource := &userGroupDataSource{}
	resp := &frameworkdatasource.MetadataResponse{}

	dataSource.Metadata(context.Background(), frameworkdatasource.MetadataRequest{ProviderTypeName: "thetalake"}, resp)

	if resp.TypeName != "thetalake_user_group" {
		t.Fatalf("expected type name thetalake_user_group, got %q", resp.TypeName)
	}
}

func TestUserGroupDataSourceConfigure(t *testing.T) {
	dataSource := &userGroupDataSource{}
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

func TestUserGroupDataSourceSchema(t *testing.T) {
	dataSource := &userGroupDataSource{}
	resp := &frameworkdatasource.SchemaResponse{}

	dataSource.Schema(context.Background(), frameworkdatasource.SchemaRequest{}, resp)

	if resp.Schema.MarkdownDescription == "" {
		t.Fatal("expected schema description to be populated")
	}
	if len(resp.Schema.Attributes) != 2 {
		t.Fatalf("expected 2 attributes, got %d", len(resp.Schema.Attributes))
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
}