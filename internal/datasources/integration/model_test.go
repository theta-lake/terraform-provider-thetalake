package integration

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestFromApiModel(t *testing.T) {
	integration := thetalake.Integration{
		Id:   7,
		Name: "Test Integration",
	}

	model := FromApiModel(integration)

	assert.Equal(t, int64(7), model.Id.ValueInt64())
	assert.Equal(t, "Test Integration", model.Name.ValueString())
}
