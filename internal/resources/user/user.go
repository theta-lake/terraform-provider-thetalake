package user

import (
	"context"
	"fmt"

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
	userModel := &UserResourceModel{}

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, userModel)...)

	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to create User", "Create failed to read User resource plan data")
		return
	}

	// Validate required fields
	if userModel.Email.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid User Configuration", "Email is required to create a User")
	}

	if userModel.Name.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid User Configuration", "Name is required to create a User")
	}

	if userModel.Password.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid User Configuration", "Password is required to create a User")
	}

	if userModel.PasswordConfirmation.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid User Configuration", "Password Confirmation is required to create a User")
	}

	if userModel.Password.ValueString() != userModel.PasswordConfirmation.ValueString() {
		resp.Diagnostics.AddError("Invalid User Configuration", "Password and Password Confirmation must match to create a User")
	}

	if userModel.RoleId.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid User Configuration", "Role ID is required to create a User")
	}

	// Convert to API model
	newUser := userModel.ToApiModel()

	// Call API to create user
	createdUser, err := r.client.CreateUser(ctx, newUser)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create User", fmt.Sprintf("Create failed with error: %s", err.Error()))
		return
	}

	// Read back created user to get all fields
	state := FromApiModel(createdUser)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := UserResourceModel{}

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
	updatedState := FromApiModel(user)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state UserResourceModel

	// Read Terraform plan and state data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read plan data")
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	updatedUser := thetalake.User{}

	// Check for changes and update fields accordingly
	if !plan.Name.Equal(state.Name) {
		updatedUser.Name = plan.Name.ValueString()
	}

	if !plan.Email.Equal(state.Email) {
		updatedUser.Email = plan.Email.ValueString()
	}

	if !plan.RoleId.Equal(state.RoleId) {
		updatedUser.RoleId = plan.RoleId.ValueInt64()
	}

	if !plan.SearchId.Equal(state.SearchId) {
		updatedUser.SearchId = plan.SearchId.ValueInt64()
	}

	updatedUser, err := r.client.UpdateUser(ctx, updatedUser)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update User", fmt.Sprintf("Update failed with error: %s", err.Error()))
		return
	}

	// Map API model to state model
	updatedState := FromApiModel(updatedUser)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserResourceModel

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
