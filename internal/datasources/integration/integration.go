package integration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type integrationDataSource struct {
	client thetalake.Client
}

func NewIntegrationDataSource() datasource.DataSource {
	return &integrationDataSource{}
}

func (r *integrationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_integration", req.ProviderTypeName)
}

func (r *integrationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = *req.ProviderData.(*thetalake.Client)
}

func (r *integrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	integrationModel := IntegrationDataSourceModel{}

	// Read Terraform configuration data
	resp.Diagnostics.Append(req.Config.Get(ctx, &integrationModel)...)

	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to read Integration", "Read failed to read Integration data source configuration")
		return
	}

	integration, err := r.client.GetIntegrationByName(ctx, integrationModel.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read Integration", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	integrationModel = FromApiModel(integration)
	// Write Terraform state data
	resp.Diagnostics.Append(resp.State.Set(ctx, integrationModel)...)
}
