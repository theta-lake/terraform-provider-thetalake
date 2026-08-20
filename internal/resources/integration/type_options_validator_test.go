package integration

import (
	"context"
	"maps"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func genericJournalingObjectValue(values map[string]attr.Value) types.Object {
	base := map[string]attr.Value{
		"download_o365_onedrive_links": types.BoolNull(),
		"download_salesforce_doclinks": types.BoolNull(),
		"index_headers":                types.StringNull(),
		"sender_spf_override":          types.StringNull(),
		"undeliverable_disabled":       types.BoolNull(),
		"undeliverable_email_address":  types.StringNull(),
		"undeliverable_email_password": types.StringNull(),
		"undeliverable_email_port":     types.Int64Null(),
		"undeliverable_email_server":   types.StringNull(),
		"undeliverable_email_user":     types.StringNull(),
	}
	maps.Copy(base, values)

	return types.ObjectValueMust(genericJournalingAttrTypes, base)
}

func TestGenericJournalingUndeliverableMailboxValidator_NullObject(t *testing.T) {
	req := validator.ObjectRequest{
		Path:        path.Root("generic_journaling"),
		ConfigValue: types.ObjectNull(genericJournalingAttrTypes),
	}
	resp := &validator.ObjectResponse{}

	genericJournalingUndeliverableMailboxValidator{}.ValidateObject(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no diagnostics for null object, got: %v", resp.Diagnostics)
	}
}

func TestGenericJournalingUndeliverableMailboxValidator_UnknownObject(t *testing.T) {
	req := validator.ObjectRequest{
		Path:        path.Root("generic_journaling"),
		ConfigValue: types.ObjectUnknown(genericJournalingAttrTypes),
	}
	resp := &validator.ObjectResponse{}

	genericJournalingUndeliverableMailboxValidator{}.ValidateObject(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no diagnostics for unknown object, got: %v", resp.Diagnostics)
	}
}

func TestGenericJournalingUndeliverableMailboxValidator_NoMailboxFieldSet(t *testing.T) {
	req := validator.ObjectRequest{
		Path:        path.Root("generic_journaling"),
		ConfigValue: genericJournalingObjectValue(nil),
	}
	resp := &validator.ObjectResponse{}

	genericJournalingUndeliverableMailboxValidator{}.ValidateObject(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no diagnostics when no mailbox field is set, got: %v", resp.Diagnostics)
	}
}

func TestGenericJournalingUndeliverableMailboxValidator_OnlyEmailAddressSet(t *testing.T) {
	req := validator.ObjectRequest{
		Path: path.Root("generic_journaling"),
		ConfigValue: genericJournalingObjectValue(map[string]attr.Value{
			"undeliverable_email_address": types.StringValue("bounce@example.com"),
		}),
	}
	resp := &validator.ObjectResponse{}

	genericJournalingUndeliverableMailboxValidator{}.ValidateObject(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no diagnostics when only undeliverable_email_address is set, got: %v", resp.Diagnostics)
	}
}

func TestGenericJournalingUndeliverableMailboxValidator_PortOnlySet(t *testing.T) {
	req := validator.ObjectRequest{
		Path: path.Root("generic_journaling"),
		ConfigValue: genericJournalingObjectValue(map[string]attr.Value{
			"undeliverable_email_port": types.Int64Value(993),
		}),
	}
	resp := &validator.ObjectResponse{}

	genericJournalingUndeliverableMailboxValidator{}.ValidateObject(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected diagnostics when only undeliverable_email_port is set")
	}
	if len(resp.Diagnostics) != 3 {
		t.Fatalf("expected 3 diagnostics, got %d: %v", len(resp.Diagnostics), resp.Diagnostics)
	}

	wantPaths := map[string]bool{
		"undeliverable_email_server":   false,
		"undeliverable_email_user":     false,
		"undeliverable_email_password": false,
	}
	for _, d := range resp.Diagnostics {
		attrDiag, ok := d.(interface{ Path() path.Path })
		if !ok {
			t.Fatalf("expected diagnostic to implement Path(), got: %v", d)
		}
		lastStep := attrDiag.Path().String()
		found := false
		for wantPath := range wantPaths {
			if lastStep == "generic_journaling."+wantPath {
				wantPaths[wantPath] = true
				found = true
			}
		}
		if !found {
			t.Fatalf("unexpected diagnostic path: %s", lastStep)
		}
	}
	for wantPath, seen := range wantPaths {
		if !seen {
			t.Fatalf("expected diagnostic on path %q", wantPath)
		}
	}
}

func TestGenericJournalingUndeliverableMailboxValidator_ServerAndUserSetWithoutPassword(t *testing.T) {
	req := validator.ObjectRequest{
		Path: path.Root("generic_journaling"),
		ConfigValue: genericJournalingObjectValue(map[string]attr.Value{
			"undeliverable_email_server": types.StringValue("email.example.com"),
			"undeliverable_email_user":   types.StringValue("mailbox-user"),
		}),
	}
	resp := &validator.ObjectResponse{}

	genericJournalingUndeliverableMailboxValidator{}.ValidateObject(context.Background(), req, resp)

	if len(resp.Diagnostics) != 1 {
		t.Fatalf("expected exactly 1 diagnostic, got %d: %v", len(resp.Diagnostics), resp.Diagnostics)
	}

	attrDiag, ok := resp.Diagnostics[0].(interface{ Path() path.Path })
	if !ok {
		t.Fatalf("expected diagnostic to implement Path(), got: %v", resp.Diagnostics[0])
	}
	wantPath := "generic_journaling.undeliverable_email_password"
	if got := attrDiag.Path().String(); got != wantPath {
		t.Fatalf("expected diagnostic on path %q, got %q", wantPath, got)
	}
}
