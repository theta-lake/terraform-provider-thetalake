package label

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
)

func (r *labelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The label resource allows for the creation and management of labels within Theta Lake.",
		Attributes: map[string]schema.Attribute{
			"background_color": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "RGB hex color value representing the background color of the label",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The created timestamp using the RFC3339 date-time format",
			},
			"hidden": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates if the label is hidden on the UI",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The label ID",
			},
			"long_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The description of the label",
			},
			"org_unit_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The org unit associated with the label",
			},
			"short_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the label shown on records",
			},
			"tagged_datums_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The number of records with this label assigned",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The updated timestamp using the RFC3339 date-time format",
			},
			"user_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The ID of the user who created the label",
			},
		},
	}
}
