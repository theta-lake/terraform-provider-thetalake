package retentionlibrary

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (r *retentionLibraryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The retention library ID",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the retention library",
			},
		},
	}
}
