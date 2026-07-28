package thetalake

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
)

func TestCreateCustomLexicon(t *testing.T) {
	request := CreateCustomLexiconRequest{
		Description: "My Lexicon description",
		Name:        "My Lexicon",
		Policies:    []int64{1, 2, 3},
		RiskType:    "risk",
		Rules:       []string{"word1", "word2", "word3"},
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var received map[string]any
		err := json.Unmarshal(body, &received)
		assert.NoError(t, err)

		assert.Equal(t, "My Lexicon description", received["description"])
		assert.Equal(t, "My Lexicon", received["name"])
		assert.Equal(t, "risk", received["risk_type"])

		rules, ok := received["rules"].([]any)
		if assert.True(t, ok) {
			assert.Equal(t, []any{"word1", "word2", "word3"}, rules)
		}

		policies, ok := received["policies"].([]any)
		if assert.True(t, ok) {
			assert.Equal(t, []any{float64(1), float64(2), float64(3)}, policies)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(readTestData("create_custom_lexicon_response.json"))
	}

	client := newTestClient(t, http.MethodPost, "/analysis/lexicons", testHandler)

	created, err := client.CreateCustomLexicon(context.Background(), request)
	assert.NoError(t, err)

	assert.Equal(t, int64(1235), created.Id)
	assert.Equal(t, "My Lexicon", created.Name)
	assert.Equal(t, "risk", created.RiskType)
	assert.Equal(t, 3, len(created.Rules))
	assert.Equal(t, "word1", created.Rules["4OyHtJthWjOcfdIPVFRga"])
	assert.Nil(t, created.DisabledAt)
	if assert.NotNil(t, created.MaxParticipants) {
		assert.Equal(t, int64(10), *created.MaxParticipants)
	}
	assert.Equal(t, []int64{1, 2, 3}, created.PolicyIds)
}

func TestGetCustomLexiconById(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_custom_lexicon_by_id_response.json"))
	}

	client := newTestClient(t, http.MethodGet, "/analysis/lexicons/1235", testHandler)

	lexicon, err := client.GetCustomLexiconById(context.Background(), 1235)
	assert.NoError(t, err)

	assert.Equal(t, int64(1235), lexicon.Id)
	assert.Equal(t, "My Lexicon", lexicon.Name)
	assert.Equal(t, 3, len(lexicon.Rules))
	assert.Equal(t, []int64{1, 2, 3}, lexicon.PolicyIds)
	if assert.NotNil(t, lexicon.StartDate) {
		assert.Equal(t, 2021, lexicon.StartDate.Year())
	}
	assert.Nil(t, lexicon.DisabledAt)
}

func TestGetCustomLexiconById_NotFound(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}

	client := newTestClient(t, http.MethodGet, "/analysis/lexicons/9999", testHandler)

	_, err := client.GetCustomLexiconById(context.Background(), 9999)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestUpdateCustomLexicon(t *testing.T) {
	name := "Updated Lexicon"
	description := "Updated Lexicon description"
	policyIds := []int64{1, 2}

	request := UpdateCustomLexiconRequest{
		Description: &description,
		Name:        &name,
		PolicyIds:   &policyIds,
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var received map[string]any
		err := json.Unmarshal(body, &received)
		assert.NoError(t, err)

		assert.Equal(t, "Updated Lexicon", received["name"])
		assert.Equal(t, "Updated Lexicon description", received["description"])

		policyIdsReceived, ok := received["policy_ids"].([]any)
		if assert.True(t, ok) {
			assert.Equal(t, []any{float64(1), float64(2)}, policyIdsReceived)
		}

		_, hasDisabled := received["disabled"]
		assert.False(t, hasDisabled)

		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("update_custom_lexicon_response.json"))
	}

	client := newTestClient(t, http.MethodPut, "/analysis/lexicons/1235", testHandler)

	updated, err := client.UpdateCustomLexicon(context.Background(), 1235, request)
	assert.NoError(t, err)

	assert.Equal(t, int64(1235), updated.Id)
	assert.Equal(t, "Updated Lexicon", updated.Name)
	assert.Equal(t, []int64{1, 2}, updated.PolicyIds)
}

func TestUpdateCustomLexicon_Disable(t *testing.T) {
	disabled := true
	request := UpdateCustomLexiconRequest{
		Disabled: &disabled,
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var received map[string]any
		err := json.Unmarshal(body, &received)
		assert.NoError(t, err)

		assert.Equal(t, true, received["disabled"])

		disabledAt := time.Now().UTC().Format(time.RFC3339)
		responseBody := `{
			"status_code": 200,
			"status_string": "OK",
			"request_id": "ca116f96-bbd5-11ef-9468-53af98260bba",
			"lexicon": {
				"id": 1235,
				"name": "My Lexicon",
				"disabled_at": "` + disabledAt + `"
			}
		}`

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseBody))
	}

	client := newTestClient(t, http.MethodPut, "/analysis/lexicons/1235", testHandler)

	updated, err := client.UpdateCustomLexicon(context.Background(), 1235, request)
	assert.NoError(t, err)

	assert.NotNil(t, updated.DisabledAt)
}
