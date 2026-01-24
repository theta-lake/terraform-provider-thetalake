package thetalake

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestGetUserGroupByName(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		// Don't need to check method here; it's done in test client
		// Don't need to check URL path here; it's done in test client

		// Return test data
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(readTestData("get_user_group_response.json")))
	}

	client := newTestClient(t, http.MethodGet, "/user_groups", testHandler)

	// Example test: Fetch retention libraries (will get empty response)
	rl, err := client.GetUserGroupByName(context.TODO(), "testing-fabio-group-nil")
	assert.NoError(t, err)
	assert.Equal(t, "testing-fabio-group-nil", rl.Name)
	assert.Equal(t, int64(2660), rl.Id)
}
