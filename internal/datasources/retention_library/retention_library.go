package retentionlibrary

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type retentionLibraryDataSource struct {
	client thetalake.Client
}

func NewRetentionLibraryDataSource() datasource.DataSource {
	return &retentionLibraryDataSource{}
}

func (r *retentionLibraryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_retention_library", req.ProviderTypeName)
}

func (r *retentionLibraryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = *req.ProviderData.(*thetalake.Client)
}

func (r *retentionLibraryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	retentionLibraryModel := RetentionLibraryDataSourceModel{}

	// Read Terraform configuration data
	resp.Diagnostics.Append(req.Config.Get(ctx, &retentionLibraryModel)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to read Retention Library", "Read failed to read Retention Library data source configuration")
		return
	}

	retentionLibrary, err := r.client.GetRetentionLibraryByName(ctx, retentionLibraryModel.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read Retention Library", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	retentionLibraryModel = FromApiModel(retentionLibrary)
	// Write Terraform state data
	resp.Diagnostics.Append(resp.State.Set(ctx, retentionLibraryModel)...)
}
