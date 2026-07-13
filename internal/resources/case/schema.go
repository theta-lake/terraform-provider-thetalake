package legalcase

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *caseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The case resource allows for the creation and management of cases within Theta Lake.",
		Attributes: map[string]schema.Attribute{
			"close_date": schema.StringAttribute{
				Optional:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The close timestamp using the RFC3339 date-time format. Setting this attribute closes the case as of the given timestamp. Removing it from the configuration reopens the case. Null while the case is open.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The created timestamp using the RFC3339 date-time format",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "A description for the case",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The case ID",
			},
			"manager_ids": schema.SetAttribute{
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "An array of user IDs assigned as managers of this case",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the case",
			},
			"number": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A usually-unique identifier for the case",
			},
			"open_date": schema.StringAttribute{
				Required:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The open timestamp using the RFC3339 date-time format",
			},
			"records_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The number of records attached to this case",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the case. One of `OPEN` or `CLOSED`.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The updated timestamp using the RFC3339 date-time format",
			},
			"visibility": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The visibility of the case. One of `PUBLIC` or `PRIVATE`.",
				Validators: []validator.String{
					stringvalidator.OneOf("PUBLIC", "PRIVATE"),
				},
			},
		},
	}
}
