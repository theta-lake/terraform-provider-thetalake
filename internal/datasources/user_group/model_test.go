package usergroup

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestFromApiModel(t *testing.T) {
	ug := thetalake.UserGroup{
		Id:   2660,
		Name: "Test User Group",
	}

	model := FromApiModel(ug)

	assert.Equal(t, int64(2660), model.Id.ValueInt64())
	assert.Equal(t, "Test User Group", model.Name.ValueString())
}
