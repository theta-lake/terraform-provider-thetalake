package directorygroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type directoryGroupDataSource struct {
	client thetalake.Client
}

func NewDirectoryGroupDataSource() datasource.DataSource {
	return &directoryGroupDataSource{}
}

func (r *directoryGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_group", req.ProviderTypeName)
}

func (r *directoryGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = *req.ProviderData.(*thetalake.Client)
}

func (r *directoryGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	directoryGroupModel := DirectoryGroupDataSourceModel{}

	// Read Terraform configuration data
	resp.Diagnostics.Append(req.Config.Get(ctx, &directoryGroupModel)...)

	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to read Directory Group", "Read failed to read Directory Group data source configuration")
		return
	}

	directoryGroup, err := r.client.GetDirectoryGroupByName(ctx, directoryGroupModel.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read Directory Group", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	directoryGroupModel = FromApiModel(directoryGroup)
	// Write Terraform state data
	resp.Diagnostics.Append(resp.State.Set(ctx, directoryGroupModel)...)
}
