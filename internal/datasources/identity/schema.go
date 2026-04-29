package identity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (r *identityDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The identity data source enables lookup of identity IDs by email address.",
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The email address of the identity to look up",
			},
			"external_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The external ID of the identity",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The identity ID",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the identity",
			},
		},
	}
}
