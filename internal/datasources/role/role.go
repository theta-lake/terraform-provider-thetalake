package role

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type roleDataSource struct {
	client thetalake.Client
}

func NewRoleDataSource() datasource.DataSource {
	return &roleDataSource{}
}

func (r *roleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_role", req.ProviderTypeName)
}

func (r *roleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = *req.ProviderData.(*thetalake.Client)
}

func (r *roleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	roleModel := RoleDataSourceModel{}

	// Read Terraform configuration data
	resp.Diagnostics.Append(req.Config.Get(ctx, &roleModel)...)

	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to read Role", "Read failed to read Role data source configuration")
		return
	}

	role, err := r.client.GetRoleByName(ctx, roleModel.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read Role", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	roleModel = FromApiModel(role)

	// Write Terraform state data
	resp.Diagnostics.Append(resp.State.Set(ctx, roleModel)...)
}
