package role

import (
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestFromApiModel(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	updated := ts.Add(time.Hour * 24)

	role := thetalake.Role{
		CreatedAt:     &ts,
		Default:       true,
		Description:   "Administrator role",
		Id:            3,
		IsBuiltIn:     true,
		Name:          "Administrator",
		NumberOfUsers: 5,
		UpdatedAt:     &updated,
	}

	model := FromApiModel(role)

	assert.Equal(t, ts.Format(time.RFC3339), model.CreatedAt.ValueString())
	assert.Equal(t, true, model.Default.ValueBool())
	assert.Equal(t, "Administrator role", model.Description.ValueString())
	assert.Equal(t, int64(3), model.Id.ValueInt64())
	assert.Equal(t, true, model.IsBuiltIn.ValueBool())
	assert.Equal(t, "Administrator", model.Name.ValueString())
	assert.Equal(t, int64(5), model.NumberOfUsers.ValueInt64())
	assert.Equal(t, updated.Format(time.RFC3339), model.UpdatedAt.ValueString())
}

func TestFromApiModel_NilUpdatedAt(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")

	role := thetalake.Role{
		CreatedAt: &ts,
		Id:        1,
		Name:      "Reviewer",
		UpdatedAt: nil,
	}

	model := FromApiModel(role)

	assert.True(t, model.UpdatedAt.IsNull())
}
