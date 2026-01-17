package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
)

func (p *ThetalakeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_server": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Theta Lake API server URL",
			},
			"client_id": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Client ID for Theta Lake API authentication",
			},
			"client_secret": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Client secret for Theta Lake API authentication",
			},
		},
	}
}
