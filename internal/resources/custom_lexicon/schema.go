package customlexicon

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *customLexiconResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The custom lexicon resource allows for the creation and management of custom lexicons within Theta Lake. " +
			"The Theta Lake API has no delete endpoint for custom lexicons — destroying this resource (or setting `disabled = true`) disables " +
			"the lexicon via `PUT`, which also removes all policies associated with it. A disabled lexicon cannot be re-enabled or re-disabled. " +
			"Changing any of the create-only attributes (`risk_type`, `rules`, `rule_scope`, `communication_direction`, the `*_enabled`/`*_analyzed` " +
			"flags, `max_participants`, `min_num_rules_with_hits`) forces replacement, which disables the old lexicon (a permanent, disabled " +
			"tombstone that loses its policies) and creates a new one.",
		Attributes: map[string]schema.Attribute{
			"accepts_input": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the custom lexicon accepts custom input",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"attachments_enabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Indicates if the custom lexicon applies to attachments. Changing this value forces replacement.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"boilerplate_enabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Indicates if the boilerplate classifier is enabled. Changing this value forces replacement.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"chatroom_name_analyzed": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Indicates if the custom lexicon applies to the chatroom name. Changing this value forces replacement.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"communication_direction": schema.SetAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "Which directions to apply the custom lexicon to. Changing this value forces replacement.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf("inbound", "outbound", "internal")),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"count_proximity_by_characters": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Indicates if the proximity is counted by characters (`true`) or words (`false`). Changing this value forces replacement.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The created timestamp using the RFC3339 date-time format",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The description of the lexicon",
			},
			"disabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Indicates if the lexicon should be disabled. Disabling a lexicon removes all policies associated with it.",
			},
			"disabled_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The disabled timestamp using the RFC3339 date-time format, if the lexicon has been disabled",
			},
			"email_smart_body": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Indicates if the custom lexicon uses smart body analysis for emails. If `true`, `rule_scope` must include `email`. Changing this value forces replacement.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"email_subject_analyzed": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Indicates if the custom lexicon applies to the email subject. If `true`, `rule_scope` must include `email`. Changing this value forces replacement.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"end_date": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The end date of the lexicon in 'YYYY-MM-DD' format",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"filename_analyzed": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Indicates if the custom lexicon applies to the document title. Changing this value forces replacement.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The custom lexicon ID",
			},
			"max_participants": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The max number of participants the custom lexicon applies to. If there are more than `max_participants`, the custom lexicon will not apply to those participants. Omit for no limit. Changing this value forces replacement.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"min_num_rules_with_hits": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Minimum number of rule hits required to create a policy hit. Omit for no minimum. Changing this value forces replacement.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the lexicon",
			},
			"policy_ids": schema.SetAttribute{
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The IDs of the policies to associate with the lexicon",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"risk_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The risk type of the lexicon. Changing this value forces replacement.",
				Validators: []validator.String{
					stringvalidator.OneOf("risk", "information", "validation"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rule_scope": schema.SetAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "The scope of the custom lexicon. Changing this value forces replacement.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf("ai_interaction", "transcript", "chat", "doc", "email", "image", "ocr")),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"rules": schema.SetAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "The custom lexicon rules (keywords/terms) for the lexicon. Up to 1000 rules, each up to 1500 characters. Changing this value forces replacement.",
				Validators: []validator.Set{
					setvalidator.SizeAtMost(1000),
					setvalidator.ValueStringsAre(stringvalidator.LengthBetween(1, 1500)),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"start_date": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "The start date of the lexicon in 'YYYY-MM-DD' format",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The updated timestamp using the RFC3339 date-time format",
			},
		},
	}
}
