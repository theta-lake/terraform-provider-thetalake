package role

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *roleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The role resource allows for the creation and management of roles within Theta Lake.",
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "Role creation timestamp using the RFC3339 date-time format. Null for built-in roles.",
			},
			"default": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if this role is the default role.",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A description of the role.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The role ID.",
			},
			"is_built_in": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the role is built-in or custom.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The role's name. Must be unique.",
			},
			"number_of_users": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The number of users currently assigned to this role.",
			},
			"permissions": schema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "An array of permissions enabled for this role. Must be a subset of the permissions returned by the GET /roles/permissions endpoint.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "Role updated timestamp using the RFC3339 date-time format. Null for built-in roles.",
			},
		},
	}
}
