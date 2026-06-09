package thetalake

import (
	"context"
	"encoding/json"
	"io"
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

func TestCreateRetentionLibrary(t *testing.T) {
	externalID := "123456-ABC"
	request := RetentionLibrary{
		Description:                "Test retention library description",
		ExternalId:                 &externalID,
		Name:                       "Test Retention Library",
		RetainInReview:             false,
		RetentionPeriodDays:        90,
		RetentionPeriodEnabled:     false,
		SecCompliantStorageEnabled: false,
		StorageAccountId:           1,
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var received RetentionLibrary
		err := json.Unmarshal(body, &received)
		assert.NoError(t, err)

		assert.Equal(t, request.Description, received.Description)
		if assert.NotNil(t, received.ExternalId) {
			assert.Equal(t, *request.ExternalId, *received.ExternalId)
		}
		assert.Equal(t, request.Name, received.Name)
		assert.Equal(t, request.RetainInReview, received.RetainInReview)
		assert.Equal(t, request.RetentionPeriodDays, received.RetentionPeriodDays)
		assert.Equal(t, request.RetentionPeriodEnabled, received.RetentionPeriodEnabled)
		assert.Equal(t, request.SecCompliantStorageEnabled, received.SecCompliantStorageEnabled)
		assert.Equal(t, request.StorageAccountId, received.StorageAccountId)

		w.WriteHeader(http.StatusCreated)
		w.Write(readTestData("create_retention_library_response.json"))
	}

	client := newTestClient(t, http.MethodPost, "/retention_libraries", testHandler)

	createdLibrary, err := client.CreateRetentionLibrary(context.Background(), request)
	assert.NoError(t, err)

	assert.Equal(t, int64(477), createdLibrary.Id)
	assert.Equal(t, "Test Retention Library", createdLibrary.Name)
	assert.Equal(t, int64(1), createdLibrary.StorageAccountId)
	assert.Equal(t, false, createdLibrary.RetainInReview)
}

func TestGetRetentionLibraryById(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_retention_library_by_id_response.json"))
	}

	client := newTestClient(t, http.MethodGet, "/retention_libraries/477", testHandler)

	retrievedLibrary, err := client.GetRetentionLibraryById(context.Background(), 477)
	assert.NoError(t, err)

	assert.Equal(t, int64(477), retrievedLibrary.Id)
	assert.Equal(t, "Test Retention Library", retrievedLibrary.Name)
	assert.Equal(t, int64(72), retrievedLibrary.DatumCount)
	assert.Equal(t, false, retrievedLibrary.RetainInReview)
	assert.Equal(t, int64(1), retrievedLibrary.StorageAccountId)
}

func TestUpdateRetentionLibrary(t *testing.T) {
	externalID := "123456-ABC"
	request := RetentionLibrary{
		Id:                         477,
		Description:                "Updated retention library description",
		ExternalId:                 &externalID,
		Name:                       "Updated Retention Library",
		RetainInReview:             true,
		RetentionPeriodDays:        180,
		RetentionPeriodEnabled:     true,
		SecCompliantStorageEnabled: false,
		StorageAccountId:           1,
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var received RetentionLibrary
		err := json.Unmarshal(body, &received)
		assert.NoError(t, err)

		assert.Equal(t, request.Id, received.Id)
		assert.Equal(t, request.Description, received.Description)
		if assert.NotNil(t, received.ExternalId) {
			assert.Equal(t, *request.ExternalId, *received.ExternalId)
		}
		assert.Equal(t, request.Name, received.Name)
		assert.Equal(t, request.RetainInReview, received.RetainInReview)
		assert.Equal(t, request.RetentionPeriodDays, received.RetentionPeriodDays)
		assert.Equal(t, request.RetentionPeriodEnabled, received.RetentionPeriodEnabled)

		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("update_retention_library_response.json"))
	}

	client := newTestClient(t, http.MethodPut, "/retention_libraries/477", testHandler)

	updatedLibrary, err := client.UpdateRetentionLibrary(context.Background(), request)
	assert.NoError(t, err)

	assert.Equal(t, int64(477), updatedLibrary.Id)
	assert.Equal(t, "Updated Retention Library", updatedLibrary.Name)
	assert.Equal(t, true, updatedLibrary.RetainInReview)
	assert.Equal(t, int64(180), updatedLibrary.RetentionPeriodDays)
}

func TestDeleteRetentionLibrary(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_retention_library_by_id_response.json"))
	}

	client := newTestClient(t, http.MethodDelete, "/retention_libraries/477", testHandler)

	err := client.DeleteRetentionLibrary(context.Background(), 477)
	assert.NoError(t, err)
}
