package thetalake

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestFindIdsToRemove(t *testing.T) {
	existingIds := []int64{1, 2, 3, 4, 5}
	newIds := []int64{2, 4, 6}

	expected := []int64{1, 3, 5}
	result := diffIdSets(existingIds, newIds)

	assert.Equal(t, expected, result, "The IDs to remove do not match the expected values")
}

func TestCreateSupervisionSpace(t *testing.T) {
	requestedSupervisionSpace := SupervisionSpace{
		AllParticipants:          false,
		AllUsers:                 false,
		Description:              "Test Supervision Space",
		ExternalId:               "ext-123",
		HardEnforce:              true,
		Name:                     "Test Space",
		SupervisionSpacePriority: 1,
		IntegrationIds:           []int64{101, 102},
		MediaTypeIds:             []int64{1, 2},
		RetentionLibraryIds:      []int64{201},
		UserGroupIds:             []int64{301, 302},
		UserIds:                  []int64{401, 402},
		DirectoryGroupIds:        []int64{501},
	}

	createSupervisionSpaceHandler := func(w http.ResponseWriter, r *http.Request) {

		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedSpace SupervisionSpace
		err := json.Unmarshal(body, &receivedSpace)
		assert.NoError(t, err)

		assert.Equal(t, requestedSupervisionSpace, receivedSpace)

		// Return test data
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(readTestData("create_supervision_space_response.json")))
	}

	addUsersHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedUserIds []int64
		err := json.Unmarshal(body, &receivedUserIds)
		assert.NoError(t, err)

		assert.Equal(t, requestedSupervisionSpace.UserIds, receivedUserIds)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(readTestData("add_supervision_space_users_response.json")))
	}

	addUserGroupsHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedUserIds []int64
		err := json.Unmarshal(body, &receivedUserIds)
		assert.NoError(t, err)

		assert.Equal(t, requestedSupervisionSpace.UserGroupIds, receivedUserIds)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(readTestData("add_supervision_space_users_response.json")))
	}

	client := newTestClientWithRoutes(t,
		testRoute{http.MethodPost, "/supervision_spaces", createSupervisionSpaceHandler},
		testRoute{http.MethodPost, "/supervision_spaces/12/users", addUsersHandler},
		testRoute{http.MethodPost, "/supervision_spaces/12/user_groups", addUserGroupsHandler},
	)

	// Example test: Fetch retention libraries (will get empty response)
	ss, err := client.CreateSupervisionSpace(context.TODO(), requestedSupervisionSpace)
	assert.NoError(t, err)
	assert.Equal(t, "Central Time", ss.Name)
	assert.Equal(t, int64(12), ss.Id)
}

func TestGetSupervisionSpaceById(t *testing.T) {
	getSupervisionSpaceHandler := func(w http.ResponseWriter, r *http.Request) {
		// Return test data
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(readTestData("get_supervision_space_by_id_response.json")))
	}

	client := newTestClient(t, http.MethodGet, "/supervision_spaces/34", getSupervisionSpaceHandler)

	ss, err := client.GetSupervisionSpaceById(context.TODO(), 34)
	assert.NoError(t, err)
	assert.Equal(t, int64(62), ss.Id)
	assert.Equal(t, "Jeff's Space", ss.Name)
	assert.Equal(t, []int64{269, 311, 268}, ss.UserIds)
}

func TestUpdateSupervisionSpace(t *testing.T) {
	requestedSupervisionSpace := SupervisionSpace{
		AllParticipants:          false,
		AllUsers:                 false,
		Description:              "Test Supervision Space",
		ExternalId:               "ext-123",
		HardEnforce:              true,
		Name:                     "Test Space",
		SupervisionSpacePriority: 1,
		Id:                       12,
		IntegrationIds:           []int64{101, 102},
		MediaTypeIds:             []int64{1, 2},
		RetentionLibraryIds:      []int64{201},
		UserGroupIds:             []int64{301, 302},
		UserIds:                  []int64{401, 402},
		DirectoryGroupIds:        []int64{501},
	}

	updateSupervisionSpaceHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedSpace SupervisionSpace
		err := json.Unmarshal(body, &receivedSpace)
		assert.NoError(t, err)

		assert.Equal(t, requestedSupervisionSpace, receivedSpace)

		// Return test data
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(readTestData("create_supervision_space_response.json")))
	}

	addUsersHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedUserIds []int64
		err := json.Unmarshal(body, &receivedUserIds)
		assert.NoError(t, err)

		assert.Equal(t, requestedSupervisionSpace.UserIds, receivedUserIds)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(readTestData("add_supervision_space_users_response.json")))
	}

	deleteUsersHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedUserIds []int64
		err := json.Unmarshal(body, &receivedUserIds)
		assert.NoError(t, err)

		expectedDeletedUserIds := []int64{422, 422} // Assuming these were the IDs to be removed
		assert.Equal(t, expectedDeletedUserIds, receivedUserIds)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(readTestData("add_supervision_space_users_response.json")))
	}

	addUserGroupsHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedUserIds []int64
		err := json.Unmarshal(body, &receivedUserIds)
		assert.NoError(t, err)

		assert.Equal(t, requestedSupervisionSpace.UserGroupIds, receivedUserIds)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(readTestData("add_supervision_space_users_response.json")))
	}

	deleteUserGroupsHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedUserIds []int64
		err := json.Unmarshal(body, &receivedUserIds)
		assert.NoError(t, err)

		expectedDeletedUserGroupIds := []int64{64} // Assuming these were the IDs to be removed
		assert.Equal(t, expectedDeletedUserGroupIds, receivedUserIds)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(readTestData("add_supervision_space_users_response.json")))
	}

	client := newTestClientWithRoutes(t,
		testRoute{http.MethodPut, "/supervision_spaces/12", updateSupervisionSpaceHandler},
		testRoute{http.MethodPost, "/supervision_spaces/12/users", addUsersHandler},
		testRoute{http.MethodPost, "/supervision_spaces/12/user_groups", addUserGroupsHandler},
		testRoute{http.MethodDelete, "/supervision_spaces/12/users", deleteUsersHandler},
		testRoute{http.MethodDelete, "/supervision_spaces/12/user_groups", deleteUserGroupsHandler},
	)

	// Example test: Fetch retention libraries (will get empty response)
	ss, err := client.UpdateSupervisionSpace(context.TODO(), requestedSupervisionSpace)
	assert.NoError(t, err)
	assert.Equal(t, "Central Time", ss.Name)
	assert.Equal(t, int64(12), ss.Id)
}

func TestDeleteSupervisionSpace(t *testing.T) {
	deleteSupervisionSpaceHandler := func(w http.ResponseWriter, r *http.Request) {
		// Return test data
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_supervision_space_by_id_response.json"))

	}

	client := newTestClient(t, http.MethodDelete, "/supervision_spaces/45", deleteSupervisionSpaceHandler)

	err := client.DeleteSupervisionSpace(context.TODO(), 45)
	assert.NoError(t, err)
}
