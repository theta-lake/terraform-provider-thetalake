package thetalake

import (
	"context"
	"encoding/json"
	"io"
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

func TestCreateUserGroup(t *testing.T) {
	externalId := "ug-ext-001"
	desc := "Test user group description"
	userGroupRequest := UserGroup{
		Name:        "Test User Group",
		Description: &desc,
		ExternalId:  &externalId,
		CategoryIds: []int64{},
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedUserGroup UserGroup
		err := json.Unmarshal(body, &receivedUserGroup)
		assert.NoError(t, err)

		assert.Equal(t, userGroupRequest.Name, receivedUserGroup.Name)
		assert.Equal(t, userGroupRequest.Description, receivedUserGroup.Description)
		assert.Equal(t, userGroupRequest.ExternalId, receivedUserGroup.ExternalId)

		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("create_user_group_response.json"))
	}

	client := newTestClient(t, http.MethodPost, "/user_groups", testHandler)

	createdUserGroup, err := client.CreateUserGroup(context.Background(), userGroupRequest)
	assert.NoError(t, err)

	assert.Equal(t, int64(75), createdUserGroup.Id)
	assert.Equal(t, "Test User Group", createdUserGroup.Name)
	assert.Equal(t, &desc, createdUserGroup.Description)
	assert.Equal(t, &externalId, createdUserGroup.ExternalId)
}

func TestGetUserGroupById(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_user_group_by_id_response.json"))
	}

	client := newTestClient(t, http.MethodGet, "/user_groups/75", testHandler)

	retrievedUserGroup, err := client.GetUserGroupById(context.Background(), 75)
	assert.NoError(t, err)

	assert.Equal(t, int64(75), retrievedUserGroup.Id)
	assert.Equal(t, "Test User Group", retrievedUserGroup.Name)
	assert.NotNil(t, retrievedUserGroup.Description)
	assert.Equal(t, "Test user group description", *retrievedUserGroup.Description)
}

func TestUpdateUserGroup(t *testing.T) {
	externalId := "ug-ext-002"
	desc := "Updated user group description"
	userGroupUpdateRequest := UserGroup{
		Id:          75,
		Name:        "Updated User Group",
		Description: &desc,
		ExternalId:  &externalId,
		CategoryIds: []int64{},
	}

	putHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedUserGroup UserGroup
		err := json.Unmarshal(body, &receivedUserGroup)
		assert.NoError(t, err)

		assert.Equal(t, userGroupUpdateRequest.Name, receivedUserGroup.Name)
		assert.Equal(t, userGroupUpdateRequest.Description, receivedUserGroup.Description)
		assert.Equal(t, userGroupUpdateRequest.ExternalId, receivedUserGroup.ExternalId)

		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("update_user_group_response.json"))
	}

	getHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_user_group_by_id_response.json"))
	}

	client := newTestClientWithRoutes(t,
		testRoute{Method: http.MethodPut, Path: "/user_groups/75", Handler: putHandler},
		testRoute{Method: http.MethodGet, Path: "/user_groups/75", Handler: getHandler},
	)

	updatedUserGroup, err := client.UpdateUserGroup(context.Background(), userGroupUpdateRequest)
	assert.NoError(t, err)

	assert.Equal(t, int64(75), updatedUserGroup.Id)
	assert.Equal(t, "Updated User Group", updatedUserGroup.Name)
	assert.Equal(t, &desc, updatedUserGroup.Description)
	assert.Equal(t, &externalId, updatedUserGroup.ExternalId)
}

func TestUpdateUserGroupMembershipDiff(t *testing.T) {
	// Current membership: users 1 and 2. Desired membership: users 2 and 3.
	// Expected: add [3], remove [1].
	externalId := "ug-ext-002"
	diffDesc := "Updated user group description"
	desiredUserGroup := UserGroup{
		Id:          75,
		Name:        "Updated User Group",
		Description: &diffDesc,
		ExternalId:  &externalId,
		UserIds:     []int64{2, 3},
	}

	var addedIds []int64
	var removedIds []int64

	putHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("update_user_group_response.json"))
	}

	// GET returns current membership: users 1 and 2
	getHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status_code": 200,
			"user_group": {
				"id": 75,
				"name": "Updated User Group",
				"description": "Updated user group description",
				"users": [{"id": 1}, {"id": 2}]
			}
		}`))
	}

	addHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		json.Unmarshal(body, &addedIds)
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("update_user_group_response.json"))
	}

	removeHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		json.Unmarshal(body, &removedIds)
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("update_user_group_response.json"))
	}

	client := newTestClientWithRoutes(t,
		testRoute{Method: http.MethodPut, Path: "/user_groups/75", Handler: putHandler},
		testRoute{Method: http.MethodGet, Path: "/user_groups/75", Handler: getHandler},
		testRoute{Method: http.MethodPut, Path: "/user_groups/75/add_users", Handler: addHandler},
		testRoute{Method: http.MethodPut, Path: "/user_groups/75/remove_users", Handler: removeHandler},
	)

	result, err := client.UpdateUserGroup(context.Background(), desiredUserGroup)
	assert.NoError(t, err)
	assert.Equal(t, []int64{2, 3}, result.UserIds)
	assert.Equal(t, []int64{3}, addedIds)
	assert.Equal(t, []int64{1}, removedIds)
}

func TestDeleteUserGroup(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_user_group_by_id_response.json"))
	}

	client := newTestClient(t, http.MethodDelete, "/user_groups/75", testHandler)

	err := client.DeleteUserGroup(context.Background(), 75)
	assert.NoError(t, err)
}
