package directorygroup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *directoryGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The directory group resource allows for the creation and management of directory groups within Theta Lake.",
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The created timestamp using the RFC3339 date-time format",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "A description for the directory group",
			},
			"external_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "An external identifier for the directory group",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The directory group ID",
			},
			"identity_ids": schema.ListAttribute{
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "An array of identity IDs to include in this directory group",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the directory group",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The updated timestamp using the RFC3339 date-time format",
			},
		},
	}
}
