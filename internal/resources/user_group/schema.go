package usergroup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *userGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The user group resource allows for the creation and management of user groups within Theta Lake.",
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The created timestamp using the RFC3339 date-time format",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The description of the user group",
			},
			"external_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "An external ID for the user group",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The user group ID",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "**Required on resource creation.** The name of the user group",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The updated timestamp using the RFC3339 date-time format",
			},
			"user_ids": schema.ListAttribute{
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "An array of user IDs to include in this user group",
			},
		},
	}
}
