package role

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type roleResource struct {
	client *thetalake.Client
}

func NewRoleResource() resource.Resource {
	return &roleResource{}
}

func (r *roleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_role", req.ProviderTypeName)
}

func (r *roleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*thetalake.Client)
}

func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := rolePlanModel{}

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("permissions"), &plan.Permissions)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to create Role", "Create failed to read Role resource plan data")
		return
	}

	// Map plan model to API model
	apiModel := toApiModel(&plan)

	role, err := r.client.CreateRole(ctx, apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Role", fmt.Sprintf("Create failed with error: %s", err.Error()))
		return
	}

	// Map API model to state model
	state := fromApiModel(role)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := roleStateModel{}

	// Read Terraform state data
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	role, err := r.client.GetRoleById(ctx, state.Id.ValueInt64())
	if err != nil {
		if errors.Is(err, thetalake.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read Role", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	// Map API model to state model
	updatedState := fromApiModel(role)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan rolePlanModel
	var state roleStateModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("permissions"), &plan.Permissions)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to update Role", "Update failed to read Role resource plan data")
		return
	}

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	// Map plan model to API model
	apiModel := toApiModel(&plan)
	apiModel.Id = state.Id.ValueInt64()

	updatedRole, err := r.client.UpdateRole(ctx, apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update Role", fmt.Sprintf("Update failed with error: %s", err.Error()))
		return
	}

	// Map API model to state model
	updatedState := fromApiModel(updatedRole)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state roleStateModel

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	err := r.client.DeleteRole(ctx, state.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete Role", fmt.Sprintf("Delete failed with error: %s", err.Error()))
		return
	}
}

// ImportState allows existing roles to be brought under Terraform management
// by specifying their ID. The ID from the import command is mapped directly
// to the "id" attribute, after which Read will populate the rest of the state.
func (r *roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected numeric role ID, got %q: %s", req.ID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
