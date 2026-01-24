package thetalake

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestGetIntegrationByName(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		// Don't need to check method here; it's done in test client
		// Don't need to check URL path here; it's done in test client

		// Return test data
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(readTestData("get_integration_response.json")))
	}

	client := newTestClient(t, http.MethodGet, "/integrations", testHandler)

	// Example test: Fetch retention libraries (will get empty response)
	rl, err := client.GetIntegrationByName(context.TODO(), "Jacob's test")
	assert.NoError(t, err)
	assert.Equal(t, "Jacob's test", rl.Name)
	assert.Equal(t, int64(1779), rl.Id)
}
