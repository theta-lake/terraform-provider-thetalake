package identity

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type identityDataSource struct {
	client thetalake.Client
}

func NewIdentityDataSource() datasource.DataSource {
	return &identityDataSource{}
}

func (r *identityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_identity", req.ProviderTypeName)
}

func (r *identityDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = *req.ProviderData.(*thetalake.Client)
}

func (r *identityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model IdentityDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identity, err := r.client.GetIdentityByEmail(ctx, model.Email.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read Identity", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	model = FromApiModel(identity)

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
