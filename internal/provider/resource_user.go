package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type userResource struct {
	client *http.Client
}

type userModel struct {
	CreatedAt             string         `tfsdk:"created_at"`
	CurrentOrgUnit        currentOrgUnit `tfsdk:"current_org_unit"`
	Disabled              bool           `tfsdk:"disabled"`
	DisabledAt            string         `tfsdk:"disabled_at"`
	DefaultUserTimezone   string         `tfsdk:"default_user_timezone"`
	Email                 string         `tfsdk:"email"`
	ForceSso              bool           `tfsdk:"force_sso"`
	HasDatums             bool           `tfsdk:"has_datums"`
	HasMultipleWorkspaces bool           `tfsdk:"has_multiple_workspaces"`
	ID                    int            `tfsdk:"id"`
	LastLogin             string         `tfsdk:"last_login"`
	Name                  string         `tfsdk:"name"`
	OtpEnabled            bool           `tfsdk:"otp_enabled"`
	OtpEnabledAt          string         `tfsdk:"otp_enabled_at"`
	PasswordChangedAt     string         `tfsdk:"password_changed_at"`
	QueuePaused           bool           `tfsdk:"queue_paused"`
	Role                  role           `tfsdk:"role"`
	SecurityFilter        securityFilter `tfsdk:"security_filter"`
	UpdatedAt             string         `tfsdk:"updated_at"`
	UserInitials          string         `tfsdk:"user_initials"`
}

type currentOrgUnit struct {
	ArchiveOnly bool   `tfsdk:"archive_only"`
	ID          int    `tfsdk:"id"`
	Name        string `tfsdk:"name"`
}

type role struct {
	ID   int    `tfsdk:"id"`
	Name string `tfsdk:"name"`
}

type securityFilter struct {
	SearchID int    `tfsdk:"search_id"`
	Name     string `tfsdk:"name"`
}

type createUserRequest struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	RoleId               int    `json:"role_id"`
	SearchId             int    `json:"search_id"`
}

type createUserResponse struct {
	User userModel `json:"user"`
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
	r.client = req.ProviderData.(*http.Client)
}

func (r *userResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// TODO: add all model attributes
			"name": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &userModel{}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		// TODO: Add logging
		return
	}

	// TODO: add request values
	createUserReq := &createUserRequest{}

	body, err := json.Marshal(createUserReq)
	if err != nil {
		// TODO: Add logging
		return
	}

	// TODO: add token from provider
	token := fmt.Sprintf("Bearer %s", "insert-token")

	// TODO: add endpoint from provider
	userReq, err := http.NewRequest("POST", "insert-endpoint", bytes.NewReader(body))
	if err != nil {
		return
	}
	userReq.Header.Set("Authorization", token)
	userReq.Header.Set("Content-Type", "application/json")

	userResp, err := r.client.Do(userReq)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Status: %d", userResp.StatusCode))
		return
	}

	createUserResp := &createUserResponse{}

	json.NewDecoder(userResp.Body).Decode(createUserResp)

	// TODO: populate plan
	plan.Name = createUserResp.User.Name

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
