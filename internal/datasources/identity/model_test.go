package identity

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestFromApiModel(t *testing.T) {
	email := "test@example.com"
	extId := "ext-123"

	identity := thetalake.Identity{
		Id:         42,
		Name:       "Test Identity",
		Email:      &email,
		ExternalId: &extId,
	}

	model := FromApiModel(identity)

	assert.Equal(t, int64(42), model.Id.ValueInt64())
	assert.Equal(t, "Test Identity", model.Name.ValueString())
	assert.Equal(t, email, model.Email.ValueString())
	assert.Equal(t, extId, model.ExternalId.ValueString())
}

func TestFromApiModel_NilOptionalFields(t *testing.T) {
	identity := thetalake.Identity{
		Id:         1,
		Name:       "Minimal Identity",
		Email:      nil,
		ExternalId: nil,
	}

	model := FromApiModel(identity)

	assert.Equal(t, int64(1), model.Id.ValueInt64())
	assert.True(t, model.Email.IsNull())
	assert.True(t, model.ExternalId.IsNull())
}
