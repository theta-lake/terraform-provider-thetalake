package thetalake

import (
	"context"
	"encoding/json"
	"errors"
	"io"
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

func TestGetIntegrationByName_NotFound(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(readTestData("get_integration_response.json")))
	}

	client := newTestClient(t, http.MethodGet, "/integrations", testHandler)

	_, err := client.GetIntegrationByName(context.TODO(), "no such integration")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestCreateIntegration(t *testing.T) {
	paused := false
	indexHeaders := "X-Header-Score,X-Routed-Via"
	request := Integration{
		Name:   "Custom Generic Journaling Integration",
		Type:   IntegrationTypeGenericJournaling,
		Paused: &paused,
		Options: &IntegrationOptions{
			IndexHeaders: &indexHeaders,
		},
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var received map[string]any
		err := json.Unmarshal(body, &received)
		assert.NoError(t, err)

		// The wire body must use "options" and "paused", not "service_params"/"service_paused".
		assert.Equal(t, "Custom Generic Journaling Integration", received["name"])
		assert.Equal(t, "generic_journaling", received["type"])
		assert.Equal(t, false, received["paused"])
		if assert.NotNil(t, received["options"]) {
			options := received["options"].(map[string]any)
			assert.Equal(t, "X-Header-Score,X-Routed-Via", options["index_headers"])
		}
		_, hasServiceParams := received["service_params"]
		assert.False(t, hasServiceParams)

		w.WriteHeader(http.StatusCreated)
		w.Write(readTestData("create_integration_response.json"))
	}

	client := newTestClient(t, http.MethodPost, "/integrations", testHandler)

	created, err := client.CreateIntegration(context.Background(), request)
	assert.NoError(t, err)

	assert.Equal(t, int64(302), created.Id)
	assert.Equal(t, "Custom Generic Journaling Integration", created.Name)
	assert.Equal(t, "Generic Journaling", created.IntegrationType)
	assert.Equal(t, int64(41), created.IntegrationTypeId)
	assert.Equal(t, false, created.ServicePaused)
	if assert.NotNil(t, created.ServiceParams) {
		if assert.NotNil(t, created.ServiceParams.UndeliverableEmailServer) {
			assert.Equal(t, "email.realbank.com", *created.ServiceParams.UndeliverableEmailServer)
		}
	}
}

func TestGetIntegrationById(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_integration_by_id_response.json"))
	}

	client := newTestClient(t, http.MethodGet, "/integrations/302", testHandler)

	integration, err := client.GetIntegrationById(context.Background(), 302)
	assert.NoError(t, err)

	assert.Equal(t, int64(302), integration.Id)
	assert.Equal(t, "Custom Generic Journaling Integration", integration.Name)
	assert.Equal(t, "Generic Journaling", integration.IntegrationType)
	assert.Equal(t, int64(41), integration.IntegrationTypeId)
	assert.Nil(t, integration.ServiceParams)
}

func TestGetIntegrationConfiguration(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_integration_configuration_response.json"))
	}

	client := newTestClient(t, http.MethodGet, "/integrations/302/configuration", testHandler)

	configuration, err := client.GetIntegrationConfiguration(context.Background(), 302)
	assert.NoError(t, err)

	assert.Equal(t, "Generic Journaling", configuration.IntegrationType)
	assert.Equal(t, int64(41), configuration.IntegrationTypeId)
	if assert.NotNil(t, configuration.Options.UndeliverableEmailServer) {
		assert.Equal(t, "email.realbank.com", *configuration.Options.UndeliverableEmailServer)
	}
}

func TestUpdateIntegration(t *testing.T) {
	indexHeaders := "X-Header-Score,X-Routed-Via"
	request := Integration{
		Id:   302,
		Name: "Updated Generic Journaling Integration",
		Type: IntegrationTypeGenericJournaling,
		Options: &IntegrationOptions{
			IndexHeaders: &indexHeaders,
		},
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var received map[string]any
		err := json.Unmarshal(body, &received)
		assert.NoError(t, err)

		assert.Equal(t, "Updated Generic Journaling Integration", received["name"])
		assert.Equal(t, "generic_journaling", received["type"])

		// UpdateIntegrationRequest has no paused field in the spec; it must be omitted entirely.
		_, hasPaused := received["paused"]
		assert.False(t, hasPaused)

		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("update_integration_response.json"))
	}

	client := newTestClient(t, http.MethodPut, "/integrations/302", testHandler)

	updated, err := client.UpdateIntegration(context.Background(), request)
	assert.NoError(t, err)

	assert.Equal(t, int64(302), updated.Id)
	assert.Equal(t, "Updated Generic Journaling Integration", updated.Name)
}

func TestPauseIntegration(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status_code":200,"status_string":"OK","request_id":"r","status":"The integration has been paused"}`))
	}

	client := newTestClient(t, http.MethodPut, "/integrations/302/pause", testHandler)

	err := client.PauseIntegration(context.Background(), 302)
	assert.NoError(t, err)
}

func TestStartIntegration(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status_code":200,"status_string":"OK","request_id":"r","status":"The integration has been started"}`))
	}

	client := newTestClient(t, http.MethodPut, "/integrations/302/start", testHandler)

	err := client.StartIntegration(context.Background(), 302)
	assert.NoError(t, err)
}

func TestDeleteIntegration(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status_code":200,"status_string":"OK","request_id":"r","status":"The integration has been removed"}`))
	}

	client := newTestClient(t, http.MethodDelete, "/integrations/302", testHandler)

	err := client.DeleteIntegration(context.Background(), 302)
	assert.NoError(t, err)
}

func TestIntegrationTypeSlug(t *testing.T) {
	tests := []struct {
		name     string
		typeId   int64
		typeName string
		want     string
	}{
		{
			name:     "id and name both known and agree",
			typeId:   41,
			typeName: "Generic Journaling",
			want:     IntegrationTypeGenericJournaling,
		},
		{
			name:     "id and name both known and agree (google workspace email)",
			typeId:   57,
			typeName: "Google Workspace Email",
			want:     IntegrationTypeGoogleWorkspaceEmail,
		},
		{
			name:     "id and name both known and agree (theta lake api)",
			typeId:   80,
			typeName: "Theta Lake API",
			want:     IntegrationTypeThetaLakeApi,
		},
		{
			name:     "name only match: unrecognized id falls back to name",
			typeId:   9999,
			typeName: "Generic Journaling",
			want:     IntegrationTypeGenericJournaling,
		},
		{
			name:     "id only match: unrecognized name falls back to id",
			typeId:   57,
			typeName: "Some New Display Name",
			want:     IntegrationTypeGoogleWorkspaceEmail,
		},
		{
			name:     "name and id disagree: name wins",
			typeId:   41, // generic_journaling
			typeName: "Theta Lake API",
			want:     IntegrationTypeThetaLakeApi,
		},
		{
			name:     "neither recognized",
			typeId:   9999,
			typeName: "Zoom",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IntegrationTypeSlug(tt.typeId, tt.typeName))
		})
	}
}
