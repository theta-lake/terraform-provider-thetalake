package user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *userResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The created timestamp using the RFC3339 date-time format",
			},
			"current_org_unit": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "The ID of the current org unit",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The name of the current org unit",
					},
				},
			},
			"disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the user has been disabled",
			},
			"disabled_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "When the user was disabled timestamp using the RFC3339 date-time format",
			},
			"default_user_timezone": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timezone the user is assigned to by default",
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "**Required on resource creation.** The users's email address",
			},
			"force_sso": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the user has to use SSO to sign in",
			},
			"has_multiple_workspaces": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the user is assigned to multiple workspaces",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The user ID",
			},
			"last_login": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user's last login timestamp using the RFC3339 date-time format",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "**Required on resource creation.** The name of the user",
			},
			"otp_enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the use has OTP enabled",
			},
			"otp_enabled_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp of when OTP was enabled using the RFC3339 date-time format",
			},
			"password": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				WriteOnly:           true,
				MarkdownDescription: "**Required on resource creation.** The user's password",
			},
			"password_changed_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp of the last time the user's password was changed using the RFC3339 date-time format",
			},
			"queue_paused": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the users queue has been paused",
			},
			// "role": schema.SingleNestedAttribute{
			// 	Computed: true,
			// 	Attributes: map[string]schema.Attribute{
			// 		"id": schema.Int64Attribute{
			// 			Computed:            true,
			// 			MarkdownDescription: "The role ID",
			// 		},
			// 		"name": schema.StringAttribute{
			// 			Computed:            true,
			// 			MarkdownDescription: "The role name",
			// 		},
			// 	},
			// },
			"role": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "**Required on resource creation.** The name of the role that will be assigned to the new user",
			},
			"search_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "**Optional for resource creation.** The ID of the search used for the user's security filter",
			},
			"security_filter": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"search_id": schema.Int64Attribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The updated timestamp using the RFC3339 date-time format",
			},
		},
	}
}
