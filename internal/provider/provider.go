package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
	directorygroupds "github.com/theta-lake/terraform-provider-thetalake/internal/datasources/directory_group"
	"github.com/theta-lake/terraform-provider-thetalake/internal/datasources/identity"
	"github.com/theta-lake/terraform-provider-thetalake/internal/datasources/integration"
	retentionlibrary "github.com/theta-lake/terraform-provider-thetalake/internal/datasources/retention_library"
	"github.com/theta-lake/terraform-provider-thetalake/internal/datasources/role"
	userds "github.com/theta-lake/terraform-provider-thetalake/internal/datasources/user"
	usergroup "github.com/theta-lake/terraform-provider-thetalake/internal/datasources/user_group"
	directorygroup "github.com/theta-lake/terraform-provider-thetalake/internal/resources/directory_group"
	"github.com/theta-lake/terraform-provider-thetalake/internal/resources/label"
	retentionlibraryresource "github.com/theta-lake/terraform-provider-thetalake/internal/resources/retention_library"
	supervisionspace "github.com/theta-lake/terraform-provider-thetalake/internal/resources/supervision_space"
	"github.com/theta-lake/terraform-provider-thetalake/internal/resources/user"
	usergroupresource "github.com/theta-lake/terraform-provider-thetalake/internal/resources/user_group"
)

var _ provider.Provider = &ThetalakeProvider{}

type ThetalakeProvider struct {
	version string
	client  *thetalake.Client // optional pre-built client; used by tests to avoid repeated token fetches
}

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

// NewWithClient returns a provider factory that uses the supplied pre-built
// client instead of constructing one in Configure. Intended for acceptance
// tests so a single authenticated client is shared across all test steps.
func NewWithClient(version string, client *thetalake.Client) func() provider.Provider {
	return func() provider.Provider {
		return &ThetalakeProvider{
			version: version,
			client:  client,
		}
	}
}

func (p *ThetalakeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "thetalake"
	resp.Version = p.version
}

func (p *ThetalakeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// If a client was injected (e.g. by tests), skip configuration and use it directly.
	if p.client != nil {
		resp.DataSourceData = p.client
		resp.ResourceData = p.client
		return
	}

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
	client.SetVersion(p.version)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ThetalakeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		user.NewUserResource,
		supervisionspace.NewSupervisionSpaceResource,
		label.NewLabelResource,
		directorygroup.NewDirectoryGroupResource,
		usergroupresource.NewUserGroupResource,
		retentionlibraryresource.NewRetentionLibraryResource,
	}
}

func (p *ThetalakeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		role.NewRoleDataSource,
		integration.NewIntegrationDataSource,
		retentionlibrary.NewRetentionLibraryDataSource,
		directorygroupds.NewDirectoryGroupDataSource,
		usergroup.NewUserGroupDataSource,
		identity.NewIdentityDataSource,
		userds.NewUserDataSource,
	}
}
