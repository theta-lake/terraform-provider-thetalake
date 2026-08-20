package integration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// typeBlockPaths lists the three type-selecting nested attributes. Exactly one must be
// set; which one determines the integration's type and cannot change without replacing
// the resource.
var typeBlockPaths = []path.Expression{
	path.MatchRoot("generic_journaling"),
	path.MatchRoot("google_workspace_email"),
	path.MatchRoot("theta_lake_api"),
}

// requiresReplaceOnTypeChange forces replacement when a type block is added or removed
// (i.e. the integration's type changes), but allows in-place updates when values inside
// an already-set block change.
func requiresReplaceOnTypeChange() planmodifier.Object {
	return objectplanmodifier.RequiresReplaceIf(
		func(_ context.Context, req planmodifier.ObjectRequest, resp *objectplanmodifier.RequiresReplaceIfFuncResponse) {
			resp.RequiresReplace = req.StateValue.IsNull() != req.PlanValue.IsNull()
		},
		"an integration's type cannot be changed after creation",
		"an integration's `type` cannot be changed after creation",
	)
}

func (r *integrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The integration resource allows for the creation and management of integrations within Theta Lake. " +
			"Exactly one of `generic_journaling`, `google_workspace_email`, or `theta_lake_api` must be set; whichever is set " +
			"determines the integration's type, which cannot be changed after creation.",
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The created timestamp using the RFC3339 date-time format",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"generic_journaling": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Configuration options for a Generic Journaling integration",
				PlanModifiers: []planmodifier.Object{
					requiresReplaceOnTypeChange(),
				},
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(typeBlockPaths...),
					validateGenericJournalingUndeliverableMailbox(),
				},
				Attributes: map[string]schema.Attribute{
					"download_o365_onedrive_links": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "Indicates if links to Microsoft OneDrive/O365 documents found in email should be downloaded and included in the archived record",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"download_salesforce_doclinks": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "Indicates if Salesforce document links found in email should be downloaded and included in the archived record",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"index_headers": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "A comma separated list of email headers that should be indexed for search",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"sender_spf_override": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "An SPF DNS TXT record value used to override sender validation for the integration's journaling domain",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"undeliverable_disabled": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "Indicates if handling of undeliverable (bounce) email is disabled for this integration",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"undeliverable_email_address": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The email address used to receive undeliverable (bounce) notifications for this integration",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"undeliverable_email_password": schema.StringAttribute{
						Optional:  true,
						Computed:  true,
						Sensitive: true,
						// Deliberately has no UseStateForUnknown plan modifier, unlike its companion
						// mailbox attributes, so that an out-of-band password change is always
						// surfaced as drift on refresh rather than masked by the prior state value.
						MarkdownDescription: "The password used to authenticate with the undeliverable email mailbox. On create and update the configured value is authoritative; a password changed outside Terraform is detected on refresh.",
					},
					"undeliverable_email_port": schema.Int64Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The port used to connect to the undeliverable email mailbox's server",
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"undeliverable_email_server": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The hostname of the mail server used to check for undeliverable (bounce) email",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"undeliverable_email_user": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The username used to authenticate with the undeliverable email mailbox",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"google_workspace_email": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Configuration options for a Google Workspace Email integration",
				PlanModifiers: []planmodifier.Object{
					requiresReplaceOnTypeChange(),
				},
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(typeBlockPaths...),
				},
				Attributes: map[string]schema.Attribute{
					"download_o365_onedrive_links": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "Indicates if links to Microsoft OneDrive/O365 documents found in email should be downloaded and included in the archived record",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"index_headers": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "A comma separated list of email headers that should be indexed for search",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"sender_spf_override": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "An SPF DNS TXT record value used to override sender validation for the integration's journaling domain",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"undeliverable_email_address": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The email address used to receive undeliverable (bounce) notifications for this integration",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The integration ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"integration_group": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The category the integration belongs to (e.g. \"Collaboration Platform\")",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"integration_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the service integrated with",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"integration_type_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Theta Lake's internal integration type ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the integration",
			},
			"paused": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates if the integration is paused. Changing this after creation calls the pause/start endpoints rather than updating the integration directly.",
			},
			"theta_lake_api": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Selects the Theta Lake API integration type. This integration type has no configuration options, so the object must be set but is always empty (`{}`).",
				PlanModifiers: []planmodifier.Object{
					requiresReplaceOnTypeChange(),
				},
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(typeBlockPaths...),
				},
				Attributes: map[string]schema.Attribute{},
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The type of integration, derived from whichever of `generic_journaling`, `google_workspace_email`, or `theta_lake_api` is set",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
