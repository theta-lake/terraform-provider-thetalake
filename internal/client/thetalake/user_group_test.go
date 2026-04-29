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
	userGroupRequest := UserGroup{
		Name:        "Test User Group",
		Description: "Test user group description",
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
	assert.Equal(t, "Test user group description", createdUserGroup.Description)
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
	assert.Equal(t, "Test user group description", retrievedUserGroup.Description)
}

func TestUpdateUserGroup(t *testing.T) {
	externalId := "ug-ext-002"
	userGroupUpdateRequest := UserGroup{
		Id:          75,
		Name:        "Updated User Group",
		Description: "Updated user group description",
		ExternalId:  &externalId,
		CategoryIds: []int64{},
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
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

	client := newTestClient(t, http.MethodPut, "/user_groups/75", testHandler)

	updatedUserGroup, err := client.UpdateUserGroup(context.Background(), userGroupUpdateRequest)
	assert.NoError(t, err)

	assert.Equal(t, int64(75), updatedUserGroup.Id)
	assert.Equal(t, "Updated User Group", updatedUserGroup.Name)
	assert.Equal(t, "Updated user group description", updatedUserGroup.Description)
	assert.Equal(t, &externalId, updatedUserGroup.ExternalId)
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
