package workspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var supportedLanguages = []string{"en", "es", "nl", "de", "fr", "it", "ja", "cmn", "pt"}

func (r *workspaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The workspace resource allows for the management of workspace settings within Theta Lake. Workspaces cannot be created or deleted via the API — use `terraform import thetalake_workspace.<name> <id>` to bring an existing workspace under Terraform management.",
		Attributes: map[string]schema.Attribute{
			"allow_anonymous_via_shared_links": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates the ability to share a link to records publicly",
			},
			"analysis_supervision_space_ids": schema.SetAttribute{
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "IDs of supervision spaces whose matching records will be analyzed for risky behaviors",
			},
			"analysis_supervision_spaces": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Supervision spaces whose matching records will be analyzed for risky behaviors",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The supervision space ID",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The supervision space name",
						},
					},
				},
			},
			"audit_log_retention_period": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The duration in days that audit logs for this workspace will be retained; null means the platform default is used",
			},
			"case_management_manager_assignment": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if cases are managed by users who are manually or automatically assigned to them. If `false`, cases are managed by the supervision space manager for the participants in the case.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The created timestamp using the RFC3339 date-time format",
			},
			"default_transcription_language": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The default language used for transcription if a record's language is unknown or not in the preferred language list",
				Validators: []validator.String{
					stringvalidator.OneOf(supportedLanguages...),
				},
			},
			"default_workspace_timezone": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The workspace's default timezone (e.g. `Etc/UTC`, `America/New_York`)",
			},
			"delete_on_expiration": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if records in this workspace should be deleted when they expire from the archive",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A description of the workspace",
			},
			"disabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the workspace is disabled",
			},
			"disabled_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The timestamp when the workspace was disabled, or null if not disabled, using the RFC3339 date-time format",
			},
			"hide_attachments_from_search": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if attachments will be shown in search results",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The workspace ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The workspace name",
			},
			"preferred_languages": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "An array of selected preferred language codes",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOf(supportedLanguages...),
					),
				},
			},
			"reauthenticate_on_network_change": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether users will be forced to re-authenticate after changing networks",
			},
			"shared_links_expiration_period": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The duration in days that shared links for this workspace will be valid",
			},
			"show_system_messages_in_chat": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if system generated messages are shown in chat",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The updated timestamp using the RFC3339 date-time format",
			},
			"use_name_matcher": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the workspace identity matching algorithm uses names",
			},
			"use_owner_only_space_matcher": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the workspace identity matching algorithm should stop after owner and then fall back to default",
			},
			"user_ids": schema.SetAttribute{
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "IDs of users assigned to this workspace. When set, Terraform will add or remove users to match the set. Omit this attribute to leave workspace membership unmanaged.",
			},
			"users": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Users currently assigned to this workspace",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"email": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user's email address",
						},
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The user ID",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user's display name",
						},
					},
				},
			},
		},
	}
}
