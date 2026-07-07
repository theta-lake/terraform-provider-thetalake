package workspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *workspaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The workspace data source enables lookup of workspace IDs by name.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A description of the workspace",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The workspace ID",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the workspace to look up",
			},
		},
	}
}
