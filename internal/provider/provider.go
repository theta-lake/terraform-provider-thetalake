package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
	"github.com/theta-lake/terraform-provider-thetalake/internal/datasources/role"
	supervisionspace "github.com/theta-lake/terraform-provider-thetalake/internal/resources/supervision_space"
	"github.com/theta-lake/terraform-provider-thetalake/internal/resources/user"
)

var _ provider.Provider = &ThetalakeProvider{}

// ThetalakeProvider defines the provider implementation
type ThetalakeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing
	version string
}

// ThetalakeProviderModel describes the provider data model
type ThetalakeProviderModel struct {
	ApiServerUrl types.String `tfsdk:"api_server"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
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

func (p *ThetalakeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	providerModel := &ThetalakeProviderModel{}

	resp.Diagnostics.Append(req.Config.Get(ctx, &providerModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if providerModel.ApiServerUrl.IsUnknown() || providerModel.ApiServerUrl.IsNull() {
		resp.Diagnostics.AddError("Missing Theta Lake API endpoint", "The provider requires a Theta Lake API endpoint. This can be found at https://developer.thetalake.ai for your data center.")
		return
	}

	if providerModel.ClientId.IsUnknown() || providerModel.ClientId.IsNull() {
		resp.Diagnostics.AddError("Missing Client ID", "The provider requires a API token client ID. API tokens are managed on the Theta Lake Developer site at https://developer.thetalake.ai.")
		return

	}
	if providerModel.ClientSecret.IsUnknown() || providerModel.ClientSecret.IsNull() {
		resp.Diagnostics.AddError("Missing Client Secret", "The provider requires a API token client secret. API tokens are managed on the Theta Lake Developer site at https://developer.thetalake.ai.")
		return
	}

	client := thetalake.NewClient(providerModel.ApiServerUrl.ValueString(), providerModel.ClientId.ValueString(), providerModel.ClientSecret.ValueString())

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ThetalakeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		user.NewUserResource,
		supervisionspace.NewSupervisionSpaceResource,
	}
}

func (p *ThetalakeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		role.NewRoleDataSource,
	}
}
