package mediatypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type mediaTypeDataSource struct {
	// Intentionally left blank
}

func NewMediaTypeDataSource() datasource.DataSource {
	return &mediaTypeDataSource{}
}

func (r *mediaTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_media_type", req.ProviderTypeName)
}

func (r *mediaTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Intentionally left blank
}

func (r *mediaTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	mediaTypeModel := MediaTypeDataSourceModel{}

	// Read Terraform configuration data
	resp.Diagnostics.Append(req.Config.Get(ctx, &mediaTypeModel)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to read Media Type", "Read failed to read Media Type data source configuration")
		return
	}

	name := mediaTypeModel.Name.ValueString()
	if name == "" {
		resp.Diagnostics.AddError("Invalid Media Type Configuration", "Name is required to read a Media Type")
		return
	}
	// This is a static map since Theta Lake API does not provide an endpoint to list media types
	id, ok := mediaTypeMap[name]
	if !ok {
		resp.Diagnostics.AddError("Failed to read Media Type", fmt.Sprintf("Media Type with name %q not found", name))
		return
	}

	mediaTypeModel.Id = types.Int64Value(id)

	// Write Terraform state data
	resp.Diagnostics.Append(resp.State.Set(ctx, mediaTypeModel)...)
}
