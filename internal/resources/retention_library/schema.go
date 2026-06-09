package retentionlibrary

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func (r *retentionLibraryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The retention library resource allows for the creation and management of retention libraries within Theta Lake.",
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The created timestamp using the RFC3339 date-time format",
			},
			"datum_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The number of datums stored in the retention library",
			},
			"datum_size": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The total size of all the datums in the retention library",
			},
			"delete_on_expiration": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Denotes if this retention library should delete records associated with this space when they expire",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "A description of the retention library",
			},
			"display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The retention library display name including the region",
			},
			"external_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "An external identifier for the retention library",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The retention library ID",
			},
			"legal_hold_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The number of datums in legal hold in this retention library",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the retention library",
			},
			"retain_in_review": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates if the retention library should be retained while review is open",
			},
			"retention_period_days": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				MarkdownDescription: "Retention period in days",
			},
			"retention_period_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates if a retention period policy is currently enabled",
			},
			"retention_summary_text": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Summary description of the retention policy",
			},
			"sec_compliant_storage_confirmed": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the SEC Rule 17a-4 compliance has been confirmed",
			},
			"sec_compliant_storage_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates if the storage has SEC Rule 17a-4 compliance enabled",
			},
			"storage_account_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The storage account ID to associate with this retention library",
			},
			"swrv_rule_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Smart workflow review (SWRV) rules currently applied to this retention library",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The updated timestamp using the RFC3339 date-time format",
			},
		},
	}
}
