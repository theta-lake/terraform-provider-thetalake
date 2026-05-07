package user

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestFromApiModel(t *testing.T) {
	u := thetalake.User{}
	u.Id = 422
	u.Email = "jacob@thetalake.com"
	u.Name = "Jacob Christensen"

	model := FromApiModel(u)

	assert.Equal(t, int64(422), model.Id.ValueInt64())
	assert.Equal(t, "jacob@thetalake.com", model.Email.ValueString())
	assert.Equal(t, "Jacob Christensen", model.Name.ValueString())
}
