package label

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type labelResource struct {
	client *thetalake.Client
}

func NewLabelResource() resource.Resource {
	return &labelResource{}
}

func (r *labelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_label", req.ProviderTypeName)
}

func (r *labelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*thetalake.Client)
}

func (r *labelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := labelPlanModel{}

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("background_color"), &plan.BackgroundColor)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("hidden"), &plan.Hidden)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("long_name"), &plan.LongName)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("short_name"), &plan.ShortName)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to create Label", "Create failed to read Label resource plan data")
		return
	}

	// Map plan model to API model
	apiModel := toApiModel(&plan)

	label, err := r.client.CreateLabel(ctx, apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Label", fmt.Sprintf("Create failed with error: %s", err.Error()))
		return
	}

	// Map API model to state model
	state := fromApiModel(label)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *labelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := labelStateModel{}

	// Read Terraform state data
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	label, err := r.client.GetLabelById(ctx, state.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read Label", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	// Map API model to state model
	updatedState := fromApiModel(label)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *labelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan labelPlanModel
	var state labelStateModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("background_color"), &plan.BackgroundColor)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("hidden"), &plan.Hidden)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("long_name"), &plan.LongName)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("short_name"), &plan.ShortName)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to update Label", "Update failed to read Label resource plan data")
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

	updatedLabel, err := r.client.UpdateLabel(ctx, apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update Label", fmt.Sprintf("Update failed with error: %s", err.Error()))
		return
	}

	// Map API model to state model
	updatedState := fromApiModel(updatedLabel)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *labelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state labelStateModel

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	err := r.client.DeleteLabel(ctx, state.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete Label", fmt.Sprintf("Delete failed with error: %s", err.Error()))
		return
	}
}

// ImportState allows existing labels to be brought under Terraform
// management by specifying their ID. The ID from the import command is
// mapped directly to the "id" attribute, after which Read will
// populate the rest of the state.
func (r *labelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Parse the string ID provided by Terraform into an int64 so it
	// matches the Int64 "id" attribute in the schema.
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected numeric label ID, got %q: %s", req.ID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
