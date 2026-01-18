package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type userResource struct {
	client *thetalake.Client
}

func NewUserResource() resource.Resource {
	return &userResource{}
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_user", req.ProviderTypeName)
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*thetalake.Client)
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	userPlan := userPlanModel{}

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &userPlan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("email"), &userPlan.Email)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("password"), &userPlan.Password)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("role_id"), &userPlan.RoleId)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("security_filter_id"), &userPlan.SecurityFilterId)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to create User", "Create failed to read User resource plan data")
		return
	}

	// Validate required fields
	if userPlan.Email.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid User Configuration", "Email is required to create a User")
	}

	if userPlan.Name.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid User Configuration", "Name is required to create a User")
	}

	if userPlan.Password.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid User Configuration", fmt.Sprintf("Password is required to create a User %q", userPlan.Password.ValueString()))
	}

	if userPlan.RoleId.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid User Configuration", "Role name is required to create a User")
	}

	// Convert to API model
	newUser := toApiModel(&userPlan)

	// Call API to create user
	createdUser, err := r.client.CreateUser(ctx, newUser)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create User", fmt.Sprintf("Create failed with error: %s", err.Error()))
		return
	}

	// Read back created user to get all fields
	state := fromApiModel(createdUser)
	// The API does not return the original password, so preserve the
	// configured value in state to keep Terraform satisfied and allow
	// future plans to compare consistently.
	state.Password = userPlan.Password

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := userStateModel{}

	// Read Terraform state data
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	user, err := r.client.GetUserById(ctx, state.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read User", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	// Map API model to state model
	updatedState := fromApiModel(user)

	// The API does not return the password field, so preserve the
	// existing sensitive value from state to avoid inconsistencies.
	updatedState.Password = state.Password

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userPlanModel
	var state userStateModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("email"), &plan.Email)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("password"), &plan.Password)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("role_id"), &plan.RoleId)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("security_filter_id"), &plan.SecurityFilterId)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to create User", "Create failed to read User resource plan data")
		return
	}

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	updatedUser := thetalake.User{}

	updatedUser.Id = state.Id.ValueInt64()
	updatedUser.Email = plan.Email.ValueString() // Required in the update
	updatedUser.Name = plan.Name.ValueString()
	updatedUser.SearchId = plan.SecurityFilterId.ValueInt64()
	updatedUser.RoleId = plan.RoleId.ValueInt64()
	updatedUser.SearchId = plan.SecurityFilterId.ValueInt64()

	updatedUser, err := r.client.UpdateUser(ctx, updatedUser)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update User", fmt.Sprintf("Update failed with error: %s", err.Error()))
		return
	}

	// Map API model to state model
	updatedState := fromApiModel(updatedUser)

	// Preserve the configured password in state since the API does not
	// echo it back. Prefer the planned value, falling back to prior state.
	if !plan.Password.IsNull() && !plan.Password.IsUnknown() {
		updatedState.Password = plan.Password
	} else {
		updatedState.Password = state.Password
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userStateModel

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	// Call API to delete user
	err := r.client.DeleteUser(ctx, state.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete User", fmt.Sprintf("Delete failed with error: %s", err.Error()))
		return
	}
}

// ImportState allows existing users to be brought under Terraform
// management by specifying their ID. The ID from the import command is
// mapped directly to the "id" attribute, after which Read will
// populate the rest of the state.
func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Parse the string ID provided by Terraform into an int64 so it
	// matches the Int64 "id" attribute in the schema.
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected numeric user ID, got %q: %s", req.ID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
