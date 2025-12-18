package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type supervisionSpaceResource struct {
	apiClient *apiClient
}

type supervisionSpaceModel struct {
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
type apiSupervisionSpace struct {
	AllParticipants          bool                                  `json:"all_participants"`
	AllUsers                 bool                                  `json:"all_users"`
	CanDelete                bool                                  `json:"can_delete"`
	CanEnableAllParticipants bool                                  `json:"can_enable_all_participants"`
	CompiledParticipantList  []apiCompiledParticipant              `json:"compiled_participant_list"`
	CompiledUserList         []apiCompiledUser                     `json:"compiled_user_list"`
	CreatedAt                string                                `json:"created_at"`
	Description              string                                `json:"description"`
	DirectoryGroups          []apiSupervisionSpaceDirectoryGroup   `json:"directory_groups"`
	Disabled                 bool                                  `json:"disabled"`
	EntryStrategiesCount     int                                   `json:"entry_strategies_count"`
	ExternalId               string                                `json:"external_id"`
	HardEnforce              bool                                  `json:"hard_enforce"`
	ID                       int                                   `json:"id"`
	Integrations             []apiSupervisionSpaceIntegration      `json:"integrations"`
	MediaTypes               []apiSupervisionSpaceMediaType        `json:"media_types"`
	Name                     string                                `json:"name"`
	ParticipantCount         int                                   `json:"participant_count"`
	Participants             []apiSupervisionSpaceParticipant      `json:"participants"`
	RetentionLibraries       []apiSupervisionSpaceRetentionLibrary `json:"retention_libraries"`
	ReviewerCount            int                                   `json:"reviewer_count"`
	SupervisionSpacePriority int                                   `json:"supervision_space_priority"`
	UpdatedAt                string                                `json:"updated_at"`
	UserGroups               []apiSupervisionSpaceUserGroup        `json:"user_groups"`
	Users                    []apiSupervisionSpaceUser             `json:"users"`
}

type apiCompiledParticipant struct {
}

type apiCompiledUser struct {
}

type apiSupervisionSpaceDirectoryGroup struct {
	ID int `json:"id"`
}

type apiSupervisionSpaceIntegration struct {
	ID int `json:"id"`
}

type apiSupervisionSpaceMediaType struct {
	ID int `json:"id"`
}

type apiSupervisionSpaceParticipant struct {
}

type apiSupervisionSpaceRetentionLibrary struct {
	ID int `json:"id"`
}

type apiSupervisionSpaceUser struct {
}
type apiSupervisionSpaceUserGroup struct {
}

type apiCreateSupervisionSpaceRequest struct {
	AllParticipants          bool   `json:"all_participants"`
	AllUsers                 bool   `json:"all_users"`
	Description              string `json:"description"`
	DirectoryGroupIds        []int  `json:"directory_group_ids"`
	ExternalId               string `json:"external_id"`
	HardEnforce              bool   `json:"hard_enforce"`
	IntegrationIds           []int  `json:"integration_ids"`
	MediaTypeIds             []int  `json:"media_type_ids"`
	Name                     string `json:"name"`
	RetentionLibraryIds      []int  `json:"retention_library_ids"`
	SupervisionSpacePriority int    `json:"supervision_space_priority"`
}

type apiCreateSupervisionSpaceResponse struct {
	SupervisionSpace apiSupervisionSpace `json:"supervision_space"`
}

func NewSupervisionSpaceResource() resource.Resource {
	return &supervisionSpaceResource{}
}

func (r *supervisionSpaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_supervision_space", req.ProviderTypeName)
}

func (r *supervisionSpaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		// TODO: add logging
		return
	}
	r.apiClient = req.ProviderData.(*apiClient)
}

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

func (r *supervisionSpaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &supervisionSpaceModel{}

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read plan data")
		return
	}

	apiDirectoryGroupIds := []int{}
	resp.Diagnostics.Append(
		plan.DirectoryGroupIds.ElementsAs(ctx, &apiDirectoryGroupIds, false)...,
	)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to iterate over directory_group_ids")
		return
	}

	apiIntegrationIds := []int{}
	resp.Diagnostics.Append(
		plan.IntegrationIds.ElementsAs(ctx, &apiIntegrationIds, false)...,
	)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to iterate over integration_ids")
		return
	}

	apiMediaTypeIds := []int{}
	resp.Diagnostics.Append(
		plan.MediaTypeIds.ElementsAs(ctx, &apiMediaTypeIds, false)...,
	)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to iterate over media_type_ids")
		return
	}

	apiRetentionLibraryIds := []int{}
	resp.Diagnostics.Append(
		plan.RetentionLibraryIds.ElementsAs(ctx, &apiRetentionLibraryIds, false)...,
	)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to iterate over retention_library_ids")
		return
	}

	createSupervisionSpaceReq := &apiCreateSupervisionSpaceRequest{
		AllParticipants:          plan.AllParticipants.ValueBool(),
		AllUsers:                 plan.AllParticipants.ValueBool(),
		Description:              plan.Description.ValueString(),
		DirectoryGroupIds:        apiDirectoryGroupIds,
		ExternalId:               plan.ExternalId.ValueString(),
		HardEnforce:              plan.HardEnforce.ValueBool(),
		IntegrationIds:           apiIntegrationIds,
		MediaTypeIds:             apiMediaTypeIds,
		Name:                     plan.Name.ValueString(),
		RetentionLibraryIds:      apiRetentionLibraryIds,
		SupervisionSpacePriority: int(plan.SupervisionSpacePriority.ValueInt64()),
	}

	body, err := json.Marshal(createSupervisionSpaceReq)
	if err != nil {
		resp.Diagnostics.AddError("Bad Request", fmt.Sprintf("Failed to marshal request: %s", err.Error()))
		return
	}

	endpoint := fmt.Sprintf("%s/supervision_spaces", r.apiClient.endpoint)
	token := fmt.Sprintf("Bearer %s", r.apiClient.token)

	supervisionSpaceReq, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Failed to create request: %s", err.Error()))
		return
	}
	supervisionSpaceReq.Header.Set("Authorization", token)
	supervisionSpaceReq.Header.Set("Content-Type", "application/json")

	supervisionSpaceResp, err := r.apiClient.client.Do(supervisionSpaceReq)
	if err != nil {
		resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Status: %d, Error: %s", supervisionSpaceResp.StatusCode, err.Error()))
		return
	}

	if supervisionSpaceResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Status: %d", supervisionSpaceResp.StatusCode))
		return
	}

	createSupervisionSpaceResp := &apiCreateSupervisionSpaceResponse{}

	if err := json.NewDecoder(supervisionSpaceResp.Body).Decode(createSupervisionSpaceResp); err != nil {
		resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Failed to decode response: %s", err.Error()))
		return
	}

	state := &supervisionSpaceModel{
		AllParticipants:          types.BoolValue(createSupervisionSpaceResp.SupervisionSpace.AllParticipants),
		AllUsers:                 types.BoolValue(createSupervisionSpaceResp.SupervisionSpace.AllUsers),
		Description:              types.StringValue(createSupervisionSpaceResp.SupervisionSpace.Description),
		DirectoryGroupIds:        intSlicetoInt64List(GetSupervisionSpaceDirectoryGroupIds(createSupervisionSpaceResp.SupervisionSpace.DirectoryGroups)),
		ExternalId:               types.StringValue(createSupervisionSpaceResp.SupervisionSpace.ExternalId),
		HardEnforce:              types.BoolValue(createSupervisionSpaceResp.SupervisionSpace.HardEnforce),
		ID:                       types.Int64Value(int64(createSupervisionSpaceResp.SupervisionSpace.ID)),
		IntegrationIds:           intSlicetoInt64List(GetSupervisionSpaceIntegrationIds(createSupervisionSpaceResp.SupervisionSpace.Integrations)),
		MediaTypeIds:             intSlicetoInt64List(GetSupervisionSpaceMediaTypeIds(createSupervisionSpaceResp.SupervisionSpace.MediaTypes)),
		Name:                     types.StringValue(createSupervisionSpaceResp.SupervisionSpace.Name),
		RetentionLibraryIds:      intSlicetoInt64List(GetSupervisionSpaceRetentionLibraryIds(createSupervisionSpaceResp.SupervisionSpace.RetentionLibraries)),
		SupervisionSpacePriority: types.Int64Value(int64(createSupervisionSpaceResp.SupervisionSpace.SupervisionSpacePriority)),
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

}

func (r *supervisionSpaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *supervisionSpaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *supervisionSpaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &supervisionSpaceModel{}

	// Read Terraform plan data
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read plan data")
		return
	}

	if state.ID.IsNull() || state.ID.IsUnknown() {
		resp.Diagnostics.AddError("Internal Error", "Failed to delete the supervision_space: the id is missing in state")
		return
	}

	endpoint := fmt.Sprintf("%s/supervision_spaces/%v", r.apiClient.endpoint, state.ID.ValueInt64())
	token := fmt.Sprintf("Bearer %s", r.apiClient.token)

	supervisionSpaceReq, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Failed to create request: %s", err.Error()))
		return
	}
	supervisionSpaceReq.Header.Set("Authorization", token)
	supervisionSpaceReq.Header.Set("Content-Type", "application/json")

	supervisionSpaceResp, err := r.apiClient.client.Do(supervisionSpaceReq)
	if err != nil {
		resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Status: %d, Error: %s", supervisionSpaceResp.StatusCode, err.Error()))
		return
	}

	if supervisionSpaceResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Status: %d", supervisionSpaceResp.StatusCode))
		return
	}

	resp.State.RemoveResource(ctx)
}

func GetSupervisionSpaceDirectoryGroupIds(groups []apiSupervisionSpaceDirectoryGroup) []int {
	ids := make([]int, len(groups))
	for i, group := range groups {
		ids[i] = group.ID
	}
	return ids
}

func GetSupervisionSpaceIntegrationIds(groups []apiSupervisionSpaceIntegration) []int {
	ids := make([]int, len(groups))
	for i, group := range groups {
		ids[i] = group.ID
	}
	return ids
}

func GetSupervisionSpaceMediaTypeIds(groups []apiSupervisionSpaceMediaType) []int {
	ids := make([]int, len(groups))
	for i, group := range groups {
		ids[i] = group.ID
	}
	return ids
}

func GetSupervisionSpaceRetentionLibraryIds(groups []apiSupervisionSpaceRetentionLibrary) []int {
	ids := make([]int, len(groups))
	for i, group := range groups {
		ids[i] = group.ID
	}
	return ids
}
