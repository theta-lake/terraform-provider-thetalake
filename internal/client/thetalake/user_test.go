package thetalake

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestCreateUser(t *testing.T) {
	// Example test: Create user
	userRequest := User{}
	userRequest.Name = "Test User"
	userRequest.Email = "test@test.com"
	userRequest.Password = "TestPassword123!"
	userRequest.PasswordConfirmation = "TestPassword123!"
	userRequest.RoleId = 1
	userRequest.SearchId = 10

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		// Parse and validate request body if needed

		body, _ := io.ReadAll(r.Body)
		requestedUser := User{}
		err := json.Unmarshal(body, &requestedUser)
		assert.NoError(t, err)

		assert.Equal(t, userRequest, requestedUser)

		// Return test data
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(readTestData("create_user_response.json")))
	}

	client := newTestClient(t, http.MethodPost, "/users", testHandler)

	createdUser, err := client.CreateUser(context.Background(), userRequest)
	assert.NoError(t, err)

	assert.Equal(t, int64(5), createdUser.Id)
	assert.Equal(t, "Test User", createdUser.Name)
	assert.Equal(t, "test@test.com", createdUser.Email)
	assert.Equal(t, int64(1), createdUser.Role.Id)
	assert.Equal(t, int64(10), createdUser.SearchId)
}

func TestGetUserById(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		// Don't need to check method here; it's done in test client
		// Don't need to check URL path here; it's done in test client

		// Return test data
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(readTestData("get_user_by_id_response.json")))
	}

	client := newTestClient(t, http.MethodGet, "/users/5", testHandler)

	retrievedUser, err := client.GetUserById(context.Background(), 5)
	assert.NoError(t, err)

	assert.Equal(t, int64(422), retrievedUser.Id)
	assert.Equal(t, "Jacob Christensen", retrievedUser.Name)
	assert.Equal(t, "jacob@thetalake.com", retrievedUser.Email)
	assert.Equal(t, int64(3), retrievedUser.Role.Id)
	assert.Equal(t, int64(0), retrievedUser.SearchId)
}

func TestUpdateUser(t *testing.T) {
	// Example test: Update user
	userUpdateRequest := User{}
	userUpdateRequest.Id = 422
	userUpdateRequest.Name = "Updated User"
	userUpdateRequest.Email = "updated@test.com"
	userUpdateRequest.RoleId = 2

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		// Parse and validate request body if needed

		body, _ := io.ReadAll(r.Body)
		updatedUser := User{}
		err := json.Unmarshal(body, &updatedUser)
		assert.NoError(t, err)

		assert.Equal(t, userUpdateRequest, updatedUser)

		// Return test data
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(readTestData("update_user_response.json")))
	}

	client := newTestClient(t, http.MethodPut, "/users/422", testHandler)

	updatedUser, err := client.UpdateUser(context.Background(), userUpdateRequest)
	assert.NoError(t, err)

	assert.Equal(t, int64(422), updatedUser.Id)
	assert.Equal(t, "Updated User", updatedUser.Name)
	assert.Equal(t, "updated@test.com", updatedUser.Email)
	assert.Equal(t, int64(2), updatedUser.Role.Id)
}

func TestDeleteUser(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		// Don't need to check method here; it's done in test client
		// Don't need to check URL path here; it's done in test client

		// Return no content
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_user_by_id_response.json")) // Response body is the user that was just deleted
	}

	client := newTestClient(t, http.MethodDelete, "/users/422", testHandler)

	err := client.DeleteUser(context.Background(), 422)
	assert.NoError(t, err)
}
