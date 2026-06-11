package role

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
			"permissions": schema.SetAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "A set of permissions enabled for this role. See the [Role Permissions Reference](../guides/role_permissions.md) for the current catalog. The provider validates requested permissions against the live `/roles/permissions` endpoint.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "Role updated timestamp using the RFC3339 date-time format. Null for built-in roles.",
			},
		},
	}
}
