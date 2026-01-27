package role

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (r *roleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The role data source enables look up from role name to ID",
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "Only included in custom roles. Role creation timestamp using the RFC3339 date-time format",
			},
			"default": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if this role is the default role",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A description of the role",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The role ID",
			},
			"is_built_in": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the role is built in or custom",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the role",
			},
			"number_of_users": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The number of users currently assigned to this role",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "Only included in custom roles. Role updated timestamp using the RFC3339 date-time format",
			},
		},
	}
}
