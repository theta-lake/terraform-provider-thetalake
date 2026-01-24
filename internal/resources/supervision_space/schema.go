package supervisionspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *supervisionSpaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"all_participants": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Indicates if the supervision space contains all the participants in the org",
			},
			"all_users": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Indicates if the supervision space contains all the users in the org",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The description of the supervision space given when created",
			},
			"directory_group_ids": schema.ListAttribute{
				ElementType:         types.Int64Type,
				Required:            true,
				MarkdownDescription: "An array of directory group IDs to associate with this supervision space",
			},
			"external_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "An external ID for the supervision space",
			},
			"hard_enforce": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Indicates if users can be assigned records from other supervision spaces",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The supervision space ID",
			},
			"integration_ids": schema.ListAttribute{
				ElementType:         types.Int64Type,
				Required:            true,
				MarkdownDescription: "An array of integration IDs to associate with this supervision space",
			},
			"media_types": schema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "An array of media types to associate with this supervision space. Allowed values are: video, audio, chat, attachment, email, image",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOf("video", "audio", "chat", "attachment", "email", "image"),
					),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The supervision space name",
			},
			"retention_library_ids": schema.ListAttribute{
				ElementType:         types.Int64Type,
				Required:            true,
				MarkdownDescription: "An array of retention library IDs to associate with this supervision space",
			},
			"requested_supervision_space_priority": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The priority of assigning records to a supervision space",
			},
			"assigned_supervision_space_priority": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The priority of assigning records to a supervision space",
			},
			"can_delete": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the supervision space can be deleted",
			},
			"can_enable_all_participants": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if all participants can be enabled for the supervision space",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The created timestamp using the RFC3339 date-time format",
			},
			"disabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the supervision space has been disabled",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The updated timestamp using the RFC3339 date-time format",
			},
			"user_group_ids": schema.ListAttribute{
				ElementType:         types.Int64Type,
				Required:            true,
				MarkdownDescription: "An array of user group IDs to associate with this supervision space",
			},
			"user_ids": schema.ListAttribute{
				ElementType:         types.Int64Type,
				Required:            true,
				MarkdownDescription: "An array of user IDs to associate with this supervision space",
			},
		},
	}
}
