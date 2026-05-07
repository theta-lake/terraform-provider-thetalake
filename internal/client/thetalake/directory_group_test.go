package thetalake

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestGetDirectoryGroupByName(t *testing.T) {
	pageCounter := 1
	handler := func(w http.ResponseWriter, r *http.Request) {
		if pageCounter > 1 {
			expectedPageToken := fmt.Sprintf("page-%d-token", pageCounter)
			actualPageToken := r.URL.Query().Get("page_token")
			assert.Equal(t, expectedPageToken, actualPageToken)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		responseFile := fmt.Sprintf("get_directory_group_page_%d_response.json", pageCounter)
		w.Write(readTestData(responseFile))

		pageCounter++
	}

	client := newTestClient(t, http.MethodGet, "/directory_groups", handler)

	dg, err := client.GetDirectoryGroupByName(context.Background(), "Test directory group")
	assert.NoError(t, err)

	assert.Equal(t, int64(2065), dg.Id)
	assert.Equal(t, "Test directory group", dg.Name)
}

func TestCreateDirectoryGroup(t *testing.T) {
	createHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var received DirectoryGroup
		assert.NoError(t, json.Unmarshal(body, &received))
		assert.Equal(t, "Test Directory Group", received.Name)

		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("create_directory_group_response.json"))
	}

	addIdentitiesHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var received []int64
		assert.NoError(t, json.Unmarshal(body, &received))
		assert.Equal(t, []int64{49944}, received)

		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("add_directory_group_identities_response.json"))
	}

	client := newTestClientWithRoutes(t,
		testRoute{http.MethodPost, "/directory_groups", createHandler},
		testRoute{http.MethodPost, "/directory_groups/1996/identities", addIdentitiesHandler},
	)

	dg, err := client.CreateDirectoryGroup(context.Background(), DirectoryGroup{
		Name:        "Test Directory Group",
		IdentityIds: []int64{49944},
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1996), dg.Id)
	assert.Equal(t, "Test Directory Group", dg.Name)
	assert.Equal(t, []int64{49944}, dg.IdentityIds)
}

func TestGetDirectoryGroupById(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_directory_group_by_id_response.json"))
	}

	client := newTestClient(t, http.MethodGet, "/directory_groups/1996", handler)

	dg, err := client.GetDirectoryGroupById(context.Background(), 1996)
	assert.NoError(t, err)
	assert.Equal(t, int64(1996), dg.Id)
	assert.Equal(t, "Test Directory Group", dg.Name)
	assert.Equal(t, []int64{49944}, dg.IdentityIds)
}

func TestUpdateDirectoryGroup(t *testing.T) {
	desc := "Updated description"
	extId := "dg-ext-002"

	updateHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("update_directory_group_response.json"))
	}

	getByIdHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_directory_group_by_id_response.json"))
	}

	client := newTestClientWithRoutes(t,
		testRoute{http.MethodPut, "/directory_groups/1996", updateHandler},
		testRoute{http.MethodGet, "/directory_groups/1996", getByIdHandler},
	)

	dg, err := client.UpdateDirectoryGroup(context.Background(), DirectoryGroup{
		Id:          1996,
		Name:        "Updated Directory Group",
		Description: &desc,
		ExternalId:  &extId,
		IdentityIds: []int64{49944}, // same as in get_by_id response, so no add/remove
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1996), dg.Id)
	assert.Equal(t, "Updated Directory Group", dg.Name)
	assert.Equal(t, []int64{49944}, dg.IdentityIds)
}

func TestUpdateDirectoryGroupMembershipDiff(t *testing.T) {
	// Current membership: identities 49944 and 49945.
	// Desired membership: identities 49945 and 50000.
	// Expected: POST add [50000], DELETE remove identity 49944.
	desc := "Updated description"
	extId := "dg-ext-002"

	var addedIds []int64
	removeWasCalled := false

	updateHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("update_directory_group_response.json"))
	}

	// GET returns current membership: identities 49944 and 49945
	getByIdHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status_code": 200,
			"directory_group": {
				"id": 1996,
				"name": "Updated Directory Group",
				"identities": [{"id": 49944, "name": "Identity A"}, {"id": 49945, "name": "Identity B"}]
			}
		}`))
	}

	addIdentitiesHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		json.Unmarshal(body, &addedIds)
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("add_directory_group_identities_response.json"))
	}

	removeIdentityHandler := func(w http.ResponseWriter, r *http.Request) {
		removeWasCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status_code": 200}`))
	}

	client := newTestClientWithRoutes(t,
		testRoute{http.MethodPut, "/directory_groups/1996", updateHandler},
		testRoute{http.MethodGet, "/directory_groups/1996", getByIdHandler},
		testRoute{http.MethodPost, "/directory_groups/1996/identities", addIdentitiesHandler},
		testRoute{http.MethodDelete, "/directory_groups/1996/identity/49944", removeIdentityHandler},
	)

	dg, err := client.UpdateDirectoryGroup(context.Background(), DirectoryGroup{
		Id:          1996,
		Name:        "Updated Directory Group",
		Description: &desc,
		ExternalId:  &extId,
		IdentityIds: []int64{49945, 50000},
	})
	assert.NoError(t, err)
	assert.Equal(t, []int64{49945, 50000}, dg.IdentityIds)
	assert.Equal(t, []int64{50000}, addedIds)
	assert.True(t, removeWasCalled)
}

func TestDeleteDirectoryGroup(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status_code": 200, "status_string": "OK"}`))
	}

	client := newTestClient(t, http.MethodDelete, "/directory_groups/1996", handler)

	err := client.DeleteDirectoryGroup(context.Background(), 1996)
	assert.NoError(t, err)
}
