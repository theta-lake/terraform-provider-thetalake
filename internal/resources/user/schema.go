package user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
)

func (r *userResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The created timestamp using the RFC3339 date-time format",
			},
			"disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates if the user has been disabled",
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "**Required on resource creation.** The users's email address",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The user ID",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "**Required on resource creation.** The name of the user",
			},
			"otp_enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the use has OTP enabled",
			},
			"password": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "**Required on resource creation.** The user's password",
			},
			"queue_paused": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the users queue has been paused",
			},
			"role": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the role that will be assigned to the new user",
			},
			"role_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "**Required on resource creation.**  The ID of the role assigned to the user",
			},
			"security_filter_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "**Optional for resource creation.** The ID of the search used for the user's security filter",
			},
			"security_filter_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the search used for the user's security filter",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The updated timestamp using the RFC3339 date-time format",
			},
		},
	}
}
