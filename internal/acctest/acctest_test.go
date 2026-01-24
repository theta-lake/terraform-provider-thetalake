package acctest

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestAccUserResource_basic(t *testing.T) {
	assert.NotEmpty(t, clientId)     // These will be set via environment variables, and we want to ensure they are present
	assert.NotEmpty(t, clientSecret) // These will be set via environment variables, and we want to ensure they are present
}
