package supervisionspace

import "github.com/hashicorp/terraform-plugin-framework/types"

type supervisionSpaceResourceModel struct {
	AllParticipants          types.Bool   `tfsdk:"all_participants"`
	AllUsers                 types.Bool   `tfsdk:"all_users"`
	Description              types.String `tfsdk:"description"`
	DirectoryGroupIds        types.List   `tfsdk:"directory_group_ids"`
	ExternalId               types.String `tfsdk:"external_id"`
	HardEnforce              types.Bool   `tfsdk:"hard_enforce"`
	ID                       types.Int64  `tfsdk:"id"`
	IntegrationIds           types.List   `tfsdk:"integration_ids"`
	MediaTypeIds             types.List   `tfsdk:"media_type_ids"`
	Name                     types.String `tfsdk:"name"`
	RetentionLibraryIds      types.List   `tfsdk:"retention_library_ids"`
	SupervisionSpacePriority types.Int64  `tfsdk:"supervision_space_priority"`
}
