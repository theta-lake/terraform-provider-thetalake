package retentionlibrary

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestFromApiModel(t *testing.T) {
	rl := thetalake.RetentionLibrary{
		Id:   99,
		Name: "Test Retention Library",
	}

	model := FromApiModel(rl)

	assert.Equal(t, int64(99), model.Id.ValueInt64())
	assert.Equal(t, "Test Retention Library", model.Name.ValueString())
}
