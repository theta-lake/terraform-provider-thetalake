package user

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type userDataSource struct {
	client thetalake.Client
}

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

func (r *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_user", req.ProviderTypeName)
}

func (r *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = *req.ProviderData.(*thetalake.Client)
}

func (r *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model UserDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	u, err := r.client.GetUserByEmail(ctx, model.Email.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read User", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	model = FromApiModel(u)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
