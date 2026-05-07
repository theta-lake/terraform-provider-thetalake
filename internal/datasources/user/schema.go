package user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (r *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The user data source enables lookup of user IDs by email address.",
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The email address of the user to look up",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The user ID",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the user",
			},
		},
	}
}
