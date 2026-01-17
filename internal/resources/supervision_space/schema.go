package supervisionspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *supervisionSpaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"all_participants": schema.BoolAttribute{
				Required:            true,
				Description:         "Indicates if the supervision space contains all the participants in the org",
				MarkdownDescription: "Indicates if the supervision space contains all the participants in the org",
			},
			"all_users": schema.BoolAttribute{
				Required:            true,
				Description:         "Indicates if the supervision space contains all the users in the org",
				MarkdownDescription: "Indicates if the supervision space contains all the users in the org",
			},
			"description": schema.StringAttribute{
				Required:            true,
				Description:         "The description of the supervision space given when created",
				MarkdownDescription: "The description of the supervision space given when created",
			},
			"directory_group_ids": schema.ListAttribute{
				ElementType:         types.Int64Type,
				Required:            true,
				Description:         "An array of directory group IDs to associate with this supervision space",
				MarkdownDescription: "An array of directory group IDs to associate with this supervision space",
			},
			"external_id": schema.StringAttribute{
				Required:            true,
				Description:         "An external ID for the supervision space",
				MarkdownDescription: "An external ID for the supervision space",
			},
			"hard_enforce": schema.BoolAttribute{
				Required:            true,
				Description:         "Indicates if users can be assigned records from other supervision spaces",
				MarkdownDescription: "Indicates if users can be assigned records from other supervision spaces",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "The supervision space ID",
				MarkdownDescription: "The supervision space ID",
			},
			"integration_ids": schema.ListAttribute{
				ElementType:         types.Int64Type,
				Required:            true,
				Description:         "An array of integration IDs to associate with this supervision space",
				MarkdownDescription: "An array of integration IDs to associate with this supervision space",
			},
			"media_type_ids": schema.ListAttribute{
				ElementType:         types.Int64Type,
				Required:            true,
				Description:         "An array of media type IDs to associate with this supervision space",
				MarkdownDescription: "An array of media type IDs to associate with this supervision space",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The supervision space name",
				MarkdownDescription: "The supervision space name",
			},
			"retention_library_ids": schema.ListAttribute{
				ElementType:         types.Int64Type,
				Required:            true,
				Description:         "An array of retention library IDs to associate with this supervision space",
				MarkdownDescription: "An array of retention library IDs to associate with this supervision space",
			},
			"supervision_space_priority": schema.Int64Attribute{
				Required:            true,
				Description:         "The priority of assigning records to a supervision space",
				MarkdownDescription: "The priority of assigning records to a supervision space",
			},
		},
	}
}
