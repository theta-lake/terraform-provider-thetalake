package user

import (
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestToApiModel(t *testing.T) {
	plan := &userPlanModel{
		Email:            types.StringValue("test@email.com"),
		Name:             types.StringValue("Test User"),
		Password:         types.StringValue("SecurePassword123"),
		RoleId:           types.Int64Value(5),
		SecurityFilterId: types.Int64Value(42),
	}

	apiModel := toApiModel(plan)

	assert.Equal(t, plan.Email.ValueString(), apiModel.Email)
	assert.Equal(t, plan.Name.ValueString(), apiModel.Name)
	assert.Equal(t, plan.Password.ValueString(), apiModel.Password)
	assert.Equal(t, plan.Password.ValueString(), apiModel.PasswordConfirmation)
	assert.Equal(t, plan.RoleId.ValueInt64(), apiModel.RoleId)
	assert.Equal(t, plan.SecurityFilterId.ValueInt64(), apiModel.SearchId)
}

func TestToApiModelWithoutSecurityFilter(t *testing.T) {
	plan := &userPlanModel{
		Email:            types.StringValue("test@email.com"),
		Name:             types.StringValue("Test User"),
		Password:         types.StringValue("SecurePassword123"),
		RoleId:           types.Int64Value(5),
		SecurityFilterId: types.Int64Value(0),
	}

	apiModel := toApiModel(plan)

	assert.Equal(t, int64(0), apiModel.SearchId)
}

func TestFromApiModel(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
	apiUser := thetalake.User{}
	apiUser.CreatedAt = ts
	apiUser.Disabled = true
	apiUser.Email = "test@email.com"
	apiUser.Id = 10
	apiUser.Name = "Test User"
	apiUser.OtpEnabled = false
	apiUser.QueuePaused = true
	apiUser.Role = thetalake.EmbeddedRole{
		Name: "Admin",
		Id:   5,
	}
	apiUser.UpdatedAt = &ts

	stateModel := fromApiModel(apiUser)

	assert.Equal(t, apiUser.CreatedAt.Format(time.RFC3339), stateModel.CreatedAt.ValueString())
	assert.Equal(t, apiUser.Disabled, stateModel.Disabled.ValueBool())
	assert.Equal(t, apiUser.Email, stateModel.Email.ValueString())
	assert.Equal(t, apiUser.Id, stateModel.Id.ValueInt64())
	assert.Equal(t, apiUser.Name, stateModel.Name.ValueString())
	assert.Equal(t, apiUser.OtpEnabled, stateModel.OtpEnabled.ValueBool())
	assert.Equal(t, apiUser.QueuePaused, stateModel.QueuePaused.ValueBool())
	assert.Equal(t, apiUser.Role.Name, stateModel.Role.ValueString())
	assert.Equal(t, apiUser.Role.Id, stateModel.RoleId.ValueInt64())
	assert.Equal(t, apiUser.UpdatedAt.Format(time.RFC3339), stateModel.UpdatedAt.ValueString())
}

func TestFromApiModelWithoutSecurityFilterOrUpdatedAt(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
	apiUser := thetalake.User{}
	apiUser.CreatedAt = ts
	apiUser.Email = "test@email.com"
	apiUser.Id = 10
	apiUser.Name = "Test User"
	apiUser.Role = thetalake.EmbeddedRole{
		Name: "Admin",
		Id:   5,
	}
	apiUser.SecurityFilter = nil
	apiUser.UpdatedAt = nil

	stateModel := fromApiModel(apiUser)

	assert.True(t, stateModel.SecurityFilterId.IsNull())
	assert.True(t, stateModel.SecurityFilterName.IsNull())
	assert.True(t, stateModel.UpdatedAt.IsNull())
	assert.Equal(t, apiUser.Role.Name, stateModel.Role.ValueString())
	assert.Equal(t, apiUser.Role.Id, stateModel.RoleId.ValueInt64())
}
