package swrvrule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *swrvRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The SWRV rule resource allows for the creation and management of smart workflow review rules within Theta Lake.",
		Attributes: map[string]schema.Attribute{
			"default_rule": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if this is the default SWRV rule",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The description of the SWRV rule",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The SWRV rule ID",
			},
			"input_sources": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "The input sources for the SWRV rule",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Optional:            true,
							MarkdownDescription: "The integration ID. Required when `type` is `integration`.",
						},
						"type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The input source type",
							Validators: []validator.String{
								stringvalidator.OneOf("all_integration_uploads", "all_submission_portal_uploads", "all_uploads", "all_user_uploads", "integration"),
							},
						},
					},
				},
			},
			"is_built_in": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if this is a built-in SWRV rule",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the SWRV rule",
			},
			"policy_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The policy ID associated with this SWRV rule",
			},
			"policy_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The policy name associated with this SWRV rule",
			},
			"priority": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The priority of the SWRV rule. Lower numbers have higher priority.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"retention_library_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The retention library ID associated with this SWRV rule",
			},
			"retention_library_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The retention library name associated with this SWRV rule",
			},
			"search_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The saved search ID associated with this SWRV rule, if any",
			},
			"search_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The saved search name associated with this SWRV rule, if any",
			},
			"supervision_space_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The supervision space ID associated with this SWRV rule",
			},
			"supervision_space_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The supervision space name associated with this SWRV rule, if any",
			},
			"workflow_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The workflow ID associated with this SWRV rule",
			},
			"workflow_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The workflow name associated with this SWRV rule",
			},
		},
	}
}
