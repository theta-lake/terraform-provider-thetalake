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
	IntegrationIds           types.List   `tfsdk:"integration_ids"`
	MediaTypeIds             types.List   `tfsdk:"media_type_ids"`
	Name                     types.String `tfsdk:"name"`
	RetentionLibraryIds      types.List   `tfsdk:"retention_library_ids"`
	SupervisionSpacePriority types.Int64  `tfsdk:"supervision_space_priority"`
}
type apiSupervisionSpace struct {
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

	createSupervisionSpaceReq := &apiSupervisionSpace{
		AllParticipants:          plan.AllParticipants.ValueBool(),
		AllUsers:                 plan.AllParticipants.ValueBool(),
		Description:              plan.Description.ValueString(),
		DirectoryGroupIds:        []int{},
		ExternalId:               plan.ExternalId.ValueString(),
		HardEnforce:              plan.HardEnforce.ValueBool(),
		IntegrationIds:           []int{},
		MediaTypeIds:             []int{},
		Name:                     plan.Name.ValueString(),
		RetentionLibraryIds:      []int{},
		SupervisionSpacePriority: int(plan.SupervisionSpacePriority.ValueInt64()),
	}

	body, err := json.Marshal(createSupervisionSpaceReq)
	if err != nil {
		resp.Diagnostics.AddError("Bad Request", fmt.Sprintf("Failed to marshal request: %s", err.Error()))
		return
	}

	token := fmt.Sprintf("Bearer %s", r.apiClient.token)

	supervisionSpaceReq, err := http.NewRequest("POST", r.apiClient.endpoint, bytes.NewReader(body))
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

	createSupervisionSpaceResp := &apiCreateSupervisionSpaceResponse{}

	if err := json.NewDecoder(supervisionSpaceResp.Body).Decode(createSupervisionSpaceResp); err != nil {
		resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Failed to decode response: %s", err.Error()))
		return
	}

	state := &supervisionSpaceModel{
		AllParticipants: types.BoolValue(createSupervisionSpaceResp.SupervisionSpace.AllParticipants),
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *supervisionSpaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *supervisionSpaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *supervisionSpaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
