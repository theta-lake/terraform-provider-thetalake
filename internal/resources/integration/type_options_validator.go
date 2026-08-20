package integration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// genericJournalingUndeliverableMailboxValidator enforces that the undeliverable
// mailbox fields are only meaningful together: if any of server/user/password/port is
// set, server, user, and password must also be set. Note that undeliverable_email_address
// (the bounce notification address) is intentionally excluded from this check, since it is
// independent of the mailbox credentials. ExactlyOneOf (in schema.go) already guarantees at
// most/exactly one type block is set, so this only needs to check the fields within
// generic_journaling.
type genericJournalingUndeliverableMailboxValidator struct{}

func (v genericJournalingUndeliverableMailboxValidator) Description(context.Context) string {
	return "if any of undeliverable_email_server, undeliverable_email_user, undeliverable_email_password, or undeliverable_email_port is set, undeliverable_email_server, undeliverable_email_user, and undeliverable_email_password must also be set"
}

func (v genericJournalingUndeliverableMailboxValidator) MarkdownDescription(context.Context) string {
	return "if any of `undeliverable_email_server`, `undeliverable_email_user`, `undeliverable_email_password`, or `undeliverable_email_port` is set, `undeliverable_email_server`, `undeliverable_email_user`, and `undeliverable_email_password` must also be set"
}

func (v genericJournalingUndeliverableMailboxValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var options genericJournalingModel
	resp.Diagnostics.Append(req.ConfigValue.As(ctx, &options, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: false, UnhandledUnknownAsEmpty: false})...)
	if resp.Diagnostics.HasError() {
		return
	}

	anyMailboxFieldSet := isConfigured(options.UndeliverableEmailServer) ||
		isConfigured(options.UndeliverableEmailUser) ||
		isConfigured(options.UndeliverableEmailPassword) ||
		isConfigured(options.UndeliverableEmailPort)
	if !anyMailboxFieldSet {
		return
	}

	if !isConfigured(options.UndeliverableEmailServer) {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			req.Path.AtName("undeliverable_email_server"),
			"Missing undeliverable mailbox server",
			"undeliverable_email_server must be set when any of undeliverable_email_server, undeliverable_email_user, undeliverable_email_password, or undeliverable_email_port is set.",
		))
	}

	if !isConfigured(options.UndeliverableEmailUser) {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			req.Path.AtName("undeliverable_email_user"),
			"Missing undeliverable mailbox user",
			"undeliverable_email_user must be set when any of undeliverable_email_server, undeliverable_email_user, undeliverable_email_password, or undeliverable_email_port is set.",
		))
	}

	if !isConfigured(options.UndeliverableEmailPassword) {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			req.Path.AtName("undeliverable_email_password"),
			"Missing undeliverable mailbox password",
			"undeliverable_email_password must be set when any of undeliverable_email_server, undeliverable_email_user, undeliverable_email_password, or undeliverable_email_port is set.",
		))
	}
}

type nullableValue interface {
	IsNull() bool
	IsUnknown() bool
}

func isConfigured[T nullableValue](value T) bool {
	return !value.IsNull() && !value.IsUnknown()
}

func validateGenericJournalingUndeliverableMailbox() validator.Object {
	return genericJournalingUndeliverableMailboxValidator{}
}

var _ validator.Object = genericJournalingUndeliverableMailboxValidator{}
