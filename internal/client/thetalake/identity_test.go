package thetalake

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestGetIdentityByEmail(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test@thetalake.com", r.URL.Query().Get("query"))
		assert.Equal(t, "email", r.URL.Query().Get("field_name"))

		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_identity_by_email_response.json"))
	}

	client := newTestClient(t, http.MethodGet, "/identities", handler)

	identity, err := client.GetIdentityByEmail(context.Background(), "test@thetalake.com")
	assert.NoError(t, err)
	assert.Equal(t, int64(49944), identity.Id)
	assert.Equal(t, "Test Identity", identity.Name)
	assert.NotNil(t, identity.Email)
	assert.Equal(t, "test@thetalake.com", *identity.Email)
}
