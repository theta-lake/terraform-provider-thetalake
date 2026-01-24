package thetalake

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestGetRetentionLibraryByName(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		// Don't need to check method here; it's done in test client
		// Don't need to check URL path here; it's done in test client

		// Return test data
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(readTestData("get_retention_library_response.json")))
	}

	client := newTestClient(t, http.MethodGet, "/retention_libraries", testHandler)

	// Example test: Fetch retention libraries (will get empty response)
	rl, err := client.GetRetentionLibraryByName(context.TODO(), "Test Retention Library")
	assert.NoError(t, err)
	assert.Equal(t, "Test Retention Library", rl.Name)
}
