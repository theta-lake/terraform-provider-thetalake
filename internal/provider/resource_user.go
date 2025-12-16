package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type userResource struct {
	apiClient *apiClient
}

type userModel struct {
	CreatedAt             types.String `tfsdk:"created_at"`
	CurrentOrgUnit        types.Object `tfsdk:"current_org_unit"`
	Disabled              types.Bool   `tfsdk:"disabled"`
	DisabledAt            types.String `tfsdk:"disabled_at"`
	DefaultUserTimezone   types.String `tfsdk:"default_user_timezone"`
	Email                 types.String `tfsdk:"email"`
	ForceSso              types.Bool   `tfsdk:"force_sso"`
	HasDatums             types.Bool   `tfsdk:"has_datums"`
	HasMultipleWorkspaces types.Bool   `tfsdk:"has_multiple_workspaces"`
	ID                    types.Int64  `tfsdk:"id"`
	LastLogin             types.String `tfsdk:"last_login"`
	Name                  types.String `tfsdk:"name"`
	OtpEnabled            types.Bool   `tfsdk:"otp_enabled"`
	OtpEnabledAt          types.String `tfsdk:"otp_enabled_at"`
	Password              types.String `tfsdk:"password"`
	PasswordChangedAt     types.String `tfsdk:"password_changed_at"`
	PasswordConfirmation  types.String `tfsdk:"password_confirmation"`
	QueuePaused           types.Bool   `tfsdk:"queue_paused"`
	Role                  types.Object `tfsdk:"role"`
	RoleId                types.Int64  `tfsdk:"role_id"`
	SearchId              types.Int64  `tfsdk:"search_id"`
	SecurityFilter        types.Object `tfsdk:"security_filter"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
	UserInitials          types.String `tfsdk:"user_initials"`
}

type userModelCurrentOrgUnit struct {
	ArchiveOnly types.Bool   `tfsdk:"archive_only"`
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
}

type userModelRole struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type userModelSecurityFilter struct {
	SearchID types.Int64  `tfsdk:"search_id"`
	Name     types.String `tfsdk:"name"`
}

type apiCreateUserRequest struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	RoleId               int    `json:"role_id"`
	SearchId             *int   `json:"search_id,omitempty"`
}

type apiCreateUserResponse struct {
	User apiUser `json:"user"`
}

type apiUser struct {
	CreatedAt             string                `json:"created_at"`
	CurrentOrgUnit        apiUserCurrentOrgUnit `json:"current_org_unit"`
	Disabled              bool                  `json:"disabled"`
	DisabledAt            string                `json:"disabled_at"`
	DefaultUserTimezone   string                `json:"default_user_timezone"`
	Email                 string                `json:"email"`
	ForceSso              bool                  `json:"force_sso"`
	HasDatums             bool                  `json:"has_datums"`
	HasMultipleWorkspaces bool                  `json:"has_multiple_workspaces"`
	ID                    int                   `json:"id"`
	LastLogin             string                `json:"last_login"`
	Name                  string                `json:"name"`
	OtpEnabled            bool                  `json:"otp_enabled"`
	OtpEnabledAt          string                `json:"otp_enabled_at"`
	PasswordChangedAt     string                `json:"password_changed_at"`
	QueuePaused           bool                  `json:"queue_paused"`
	Role                  apiUserRole           `json:"role"`
	SecurityFilter        apiUserSecurityFilter `json:"security_filter"`
	UpdatedAt             string                `json:"updated_at"`
	UserInitials          string                `json:"user_initials"`
}

type apiUserCurrentOrgUnit struct {
	ArchiveOnly bool   `json:"archive_only"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
}
type apiUserRole struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type apiUserSecurityFilter struct {
	Name     *string `json:"name,omitempty"`
	SearchID *int    `json:"search_id"`
}

func NewUserResource() resource.Resource {
	return &userResource{}
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_user", req.ProviderTypeName)
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		// TODO: add logging
		return
	}
	r.apiClient = req.ProviderData.(*apiClient)
}

func (r *userResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"current_org_unit": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"archive_only": schema.BoolAttribute{
						Computed: true,
					},
					"id": schema.Int64Attribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"disabled": schema.BoolAttribute{
				Computed: true,
			},
			"disabled_at": schema.StringAttribute{
				Computed: true,
			},
			"default_user_timezone": schema.StringAttribute{
				Computed: true,
			},
			"email": schema.StringAttribute{
				Required: true,
			},
			"force_sso": schema.BoolAttribute{
				Computed: true,
			},
			"has_datums": schema.BoolAttribute{
				Computed: true,
			},
			"has_multiple_workspaces": schema.BoolAttribute{
				Computed: true,
			},
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"last_login": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"otp_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"otp_enabled_at": schema.StringAttribute{
				Computed: true,
			},
			"password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				WriteOnly: true,
			},
			"password_changed_at": schema.StringAttribute{
				Computed: true,
			},
			"password_confirmation": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				WriteOnly: true,
			},
			"queue_paused": schema.BoolAttribute{
				Computed: true,
			},
			"role": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"role_id": schema.Int64Attribute{
				Required: true,
			},
			"search_id": schema.Int64Attribute{
				Optional: true,
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
				Computed: true,
			},
			"user_initials": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &userModel{}

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read plan data")
		return
	}

	createUserReq := &apiCreateUserRequest{
		Name:                 plan.Name.ValueString(),
		Email:                plan.Email.ValueString(),
		Password:             plan.Password.ValueString(),
		PasswordConfirmation: plan.PasswordConfirmation.ValueString(),
		RoleId:               int(plan.RoleId.ValueInt64()),
	}

	if !plan.SearchId.IsNull() {
		searchId := int(plan.RoleId.ValueInt64())
		createUserReq.SearchId = &searchId
	}

	body, err := json.Marshal(createUserReq)
	if err != nil {
		resp.Diagnostics.AddError("Bad Request", fmt.Sprintf("Failed to marshal request: %s", err.Error()))
		return
	}

	token := fmt.Sprintf("Bearer %s", r.apiClient.token)

	userReq, err := http.NewRequest("POST", r.apiClient.endpoint, bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Failed to create request: %s", err.Error()))
		return
	}
	userReq.Header.Set("Authorization", token)
	userReq.Header.Set("Content-Type", "application/json")

	userResp, err := r.apiClient.client.Do(userReq)
	if err != nil {
		resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Status: %d, Error: %s", userResp.StatusCode, err.Error()))
		return
	}

	if userResp.StatusCode != http.StatusCreated {
		resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Status: %d", userResp.StatusCode))
		return
	}

	createUserResp := &apiCreateUserResponse{}

	if err := json.NewDecoder(userResp.Body).Decode(createUserResp); err != nil {
		resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Failed to decode response: %s", err.Error()))
		return
	}

	currOrgUnit, diags := basetypes.NewObjectValue(
		map[string]attr.Type{
			"archive_only": types.BoolType,
			"id":           types.Int64Type,
			"name":         types.StringType,
		},
		map[string]attr.Value{
			"archive_only": types.BoolValue(createUserResp.User.CurrentOrgUnit.ArchiveOnly),
			"id":           types.Int64Value(int64(createUserResp.User.Role.ID)),
			"name":         types.StringValue(createUserResp.User.Role.Name),
		},
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read object current_org_unit data")
		return
	}

	role, diags := basetypes.NewObjectValue(
		map[string]attr.Type{
			"id":   types.Int64Type,
			"name": types.StringType,
		},
		map[string]attr.Value{
			"id":   types.Int64Value(int64(createUserResp.User.Role.ID)),
			"name": types.StringValue(createUserResp.User.Role.Name),
		},
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read object role data")
		return
	}

	securityFilterName := types.StringNull()
	if createUserResp.User.SecurityFilter.Name != nil {
		securityFilterName = types.StringValue(*createUserResp.User.SecurityFilter.Name)
	}

	securityFilterSearchId := types.Int64Null()
	if createUserResp.User.SecurityFilter.SearchID != nil {
		securityFilterSearchId = types.Int64Value(int64(*createUserResp.User.SecurityFilter.SearchID))
	}

	securityFilter, diags := basetypes.NewObjectValue(
		map[string]attr.Type{
			"name":      types.StringType,
			"search_id": types.Int64Type,
		},
		map[string]attr.Value{
			"name":      securityFilterName,
			"search_id": securityFilterSearchId,
		},
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read object role data")
		return
	}

	state := &userModel{
		CreatedAt:             types.StringValue(createUserResp.User.CreatedAt),
		CurrentOrgUnit:        currOrgUnit,
		Disabled:              types.BoolValue(createUserResp.User.Disabled),
		DisabledAt:            types.StringValue(createUserResp.User.DisabledAt),
		DefaultUserTimezone:   types.StringValue(createUserResp.User.DefaultUserTimezone),
		Email:                 types.StringValue(createUserResp.User.Email),
		ForceSso:              types.BoolValue(createUserResp.User.ForceSso),
		HasDatums:             types.BoolValue(createUserResp.User.HasDatums),
		HasMultipleWorkspaces: types.BoolValue(createUserResp.User.HasMultipleWorkspaces),
		ID:                    types.Int64Value(int64(createUserResp.User.ID)),
		LastLogin:             types.StringValue(createUserResp.User.LastLogin),
		Name:                  types.StringValue(createUserResp.User.Name),
		OtpEnabled:            types.BoolValue(createUserResp.User.OtpEnabled),
		OtpEnabledAt:          types.StringValue(createUserResp.User.OtpEnabledAt),
		PasswordChangedAt:     types.StringValue(createUserResp.User.PasswordChangedAt),
		QueuePaused:           types.BoolValue(createUserResp.User.QueuePaused),
		Role:                  role,
		RoleId:                plan.RoleId,
		SearchId:              plan.SearchId,
		SecurityFilter:        securityFilter,
		UpdatedAt:             types.StringValue(createUserResp.User.UpdatedAt),
		UserInitials:          types.StringValue(createUserResp.User.UserInitials),
	}

	if createUserResp.User.SecurityFilter.Name != nil {
		//state.SecurityFilter.Name = types.StringValue(*createUserResp.User.SecurityFilter.Name)
	}

	if createUserResp.User.SecurityFilter.SearchID != nil {
		//state.SecurityFilter.SearchID = types.Int64Value(int64(*createUserResp.User.SecurityFilter.SearchID))
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &userModel{}

	// Read Terraform plan data
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read plan data")
		return
	}

	if state.ID.IsNull() || state.ID.IsUnknown() {
		resp.Diagnostics.AddError("Internal Error", "Failed to delete the user: the id is missing in state")
		return
	}

	token := fmt.Sprintf("Bearer %s", r.apiClient.token)

	endpoint := fmt.Sprintf("%s/%v", r.apiClient.endpoint, state.ID.ValueInt64())

	userReq, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Failed to create request: %s", err.Error()))
		return
	}
	userReq.Header.Set("Authorization", token)
	userReq.Header.Set("Content-Type", "application/json")

	userResp, err := r.apiClient.client.Do(userReq)
	if err != nil {
		resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Status: %d, Error: %s", userResp.StatusCode, err.Error()))
		return
	}

	if userResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Status: %d", userResp.StatusCode))
		return
	}

	resp.State.RemoveResource(ctx)
}
