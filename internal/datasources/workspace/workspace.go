package workspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type workspaceDataSource struct {
	client thetalake.Client
}

func NewWorkspaceDataSource() datasource.DataSource {
	return &workspaceDataSource{}
}

func (d *workspaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_workspace", req.ProviderTypeName)
}

func (d *workspaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = *req.ProviderData.(*thetalake.Client)
}

func (d *workspaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model WorkspaceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspace, err := d.client.GetWorkspaceByName(ctx, model.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read Workspace", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	model = FromApiModel(workspace)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
