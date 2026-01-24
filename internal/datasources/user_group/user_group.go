package usergroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type userGroupDataSource struct {
	client thetalake.Client
}

func NewUserGroupDataSource() datasource.DataSource {
	return &userGroupDataSource{}
}

func (r *userGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_user_group", req.ProviderTypeName)
}

func (r *userGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = *req.ProviderData.(*thetalake.Client)
}

func (r *userGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	userGroupModel := UserGroupDataSourceModel{}

	// Read Terraform configuration data
	resp.Diagnostics.Append(req.Config.Get(ctx, &userGroupModel)...)

	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to read User Group", "Read failed to read User Group data source configuration")
		return
	}

	userGroup, err := r.client.GetUserGroupByName(ctx, userGroupModel.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read User Group", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	userGroupModel = FromApiModel(userGroup)
	// Write Terraform state data
	resp.Diagnostics.Append(resp.State.Set(ctx, userGroupModel)...)
}
