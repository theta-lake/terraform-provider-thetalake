// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &ThetalakeProvider{}
var _ provider.ProviderWithFunctions = &ThetalakeProvider{}
var _ provider.ProviderWithEphemeralResources = &ThetalakeProvider{}

// ThetalakeProvider defines the provider implementation
type ThetalakeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing
	version string
}

// ThetalakeProviderModel describes the provider data model
type ThetalakeProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

type apiClient struct {
	endpoint string
	client   *http.Client
	token    string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ThetalakeProvider{
			version: version,
		}
	}
}

func (p *ThetalakeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "thetalake"
	resp.Version = p.version
}

func (p *ThetalakeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Required:            true,
				Description:         "API endpoint",
				MarkdownDescription: "API endpoint",
			},
			"token": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				Description:         "Bearer API token",
				MarkdownDescription: "Bearer API token",
			},
		},
	}
}

func (p *ThetalakeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	providerModel := &ThetalakeProviderModel{}

	resp.Diagnostics.Append(req.Config.Get(ctx, &providerModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if providerModel.Token.IsUnknown() || providerModel.Token.IsNull() || strings.TrimSpace(providerModel.Endpoint.ValueString()) == "" {
		resp.Diagnostics.AddError("Missing API token", "The provider requires an API token")
		return
	}

	client := &apiClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		endpoint: providerModel.Endpoint.ValueString(),
		token:    providerModel.Token.ValueString(),
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ThetalakeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
		NewSupervisionSpaceResource,
	}
}

func (p *ThetalakeProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return nil
}

func (p *ThetalakeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *ThetalakeProvider) Functions(ctx context.Context) []func() function.Function {
	return nil
}
