package user

import (
	"context"
	"testing"

	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestNewUserResource(t *testing.T) {
	if _, ok := NewUserResource().(*userResource); !ok {
		t.Fatal("expected NewUserResource to return *userResource")
	}
}

func TestUserMetadata(t *testing.T) {
	r := &userResource{}
	resp := &frameworkresource.MetadataResponse{}

	r.Metadata(context.Background(), frameworkresource.MetadataRequest{ProviderTypeName: "thetalake"}, resp)

	if resp.TypeName != "thetalake_user" {
		t.Fatalf("expected type name thetalake_user, got %q", resp.TypeName)
	}
}

func TestUserConfigure(t *testing.T) {
	r := &userResource{}
	resp := &frameworkresource.ConfigureResponse{}
	client := &thetalake.Client{}

	r.Configure(context.Background(), frameworkresource.ConfigureRequest{}, resp)
	if r.client != nil {
		t.Fatal("expected client to remain nil when provider data is absent")
	}

	r.Configure(context.Background(), frameworkresource.ConfigureRequest{ProviderData: client}, resp)
	if r.client != client {
		t.Fatal("expected client to be assigned from provider data")
	}
}

func TestUserSchema(t *testing.T) {
	r := &userResource{}
	resp := &frameworkresource.SchemaResponse{}

	r.Schema(context.Background(), frameworkresource.SchemaRequest{}, resp)

	if resp.Schema.MarkdownDescription == "" {
		t.Fatal("expected schema description to be populated")
	}
	if len(resp.Schema.Attributes) != 13 {
		t.Fatalf("expected 13 attributes, got %d", len(resp.Schema.Attributes))
	}

	emailAttr, ok := resp.Schema.Attributes["email"].(resourceschema.StringAttribute)
	if !ok {
		t.Fatal("expected email to be a string attribute")
	}
	if !emailAttr.Required {
		t.Fatal("expected email to be required")
	}

	passwordAttr, ok := resp.Schema.Attributes["password"].(resourceschema.StringAttribute)
	if !ok {
		t.Fatal("expected password to be a string attribute")
	}
	if !passwordAttr.Required || !passwordAttr.Sensitive {
		t.Fatal("expected password to be required and sensitive")
	}

	disabledAttr, ok := resp.Schema.Attributes["disabled"].(resourceschema.BoolAttribute)
	if !ok {
		t.Fatal("expected disabled to be a bool attribute")
	}
	if !disabledAttr.Optional || !disabledAttr.Computed {
		t.Fatal("expected disabled to be optional and computed")
	}

	securityFilterNameAttr, ok := resp.Schema.Attributes["security_filter_name"].(resourceschema.StringAttribute)
	if !ok {
		t.Fatal("expected security_filter_name to be a string attribute")
	}
	if !securityFilterNameAttr.Optional || !securityFilterNameAttr.Computed {
		t.Fatal("expected security_filter_name to be optional and computed")
	}
}

func TestUserImportStateInvalidID(t *testing.T) {
	r := &userResource{}
	resp := &frameworkresource.ImportStateResponse{}

	r.ImportState(context.Background(), frameworkresource.ImportStateRequest{ID: "not-a-number"}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected invalid import ID to add an error diagnostic")
	}
}