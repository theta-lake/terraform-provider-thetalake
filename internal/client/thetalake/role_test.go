package thetalake

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestCreateRole(t *testing.T) {
	roleRequest := Role{
		Name:        "Reviewer",
		Description: "Permissions for a user tasked solely with reviewing content",
		Permissions: []string{"cases:read", "cases:create"},
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedRole createRoleRequest
		err := json.Unmarshal(body, &receivedRole)
		assert.NoError(t, err)

		assert.Equal(t, roleRequest.Name, receivedRole.Name)
		assert.Equal(t, roleRequest.Description, receivedRole.Description)
		assert.Equal(t, roleRequest.Permissions, receivedRole.Permissions)

		w.WriteHeader(http.StatusCreated)
		w.Write(readTestData("create_role_response.json"))
	}

	client := newTestClient(t, http.MethodPost, "/roles", testHandler)

	createdRole, err := client.CreateRole(context.Background(), roleRequest)
	assert.NoError(t, err)
	assert.Equal(t, int64(369), createdRole.Id)
	assert.Equal(t, roleRequest.Name, createdRole.Name)
	assert.Equal(t, roleRequest.Description, createdRole.Description)
	assert.Equal(t, roleRequest.Permissions, createdRole.Permissions)
}

func TestGetRoleByName(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(readTestData("get_role_response.json")))
	}

	client := newTestClient(t, http.MethodGet, "/roles", testHandler)

	rl, err := client.GetRoleByName(context.TODO(), "Reviewer")
	assert.NoError(t, err)
	assert.Equal(t, "Reviewer", rl.Name)
	assert.Equal(t, int64(369), rl.Id)
	assert.Equal(t, []string{"cases:read", "cases:create"}, rl.Permissions)
}

func TestGetRolePermissions(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_role_permissions_response.json"))
	}

	client := newTestClient(t, http.MethodGet, "/roles/permissions", testHandler)

	permissions, err := client.GetRolePermissions(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, []string{"cases:add_records", "cases:create", "cases:read", "cases:update"}, permissions)
}

func TestGetRoleById(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_role_by_id_response.json"))
	}

	client := newTestClient(t, http.MethodGet, "/roles/369", testHandler)

	retrievedRole, err := client.GetRoleById(context.Background(), 369)
	assert.NoError(t, err)
	assert.Equal(t, int64(369), retrievedRole.Id)
	assert.Equal(t, "Reviewer", retrievedRole.Name)
	assert.Equal(t, "Permissions for a user tasked solely with reviewing content", retrievedRole.Description)
	assert.Equal(t, []string{"cases:read", "cases:create"}, retrievedRole.Permissions)
}

func TestUpdateRole(t *testing.T) {
	roleRequest := Role{
		Id:          369,
		Name:        "Reviewer Updated",
		Description: "Updated reviewer permissions",
		Permissions: []string{"cases:read", "cases:create", "cases:update"},
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedRole updateRoleRequest
		err := json.Unmarshal(body, &receivedRole)
		assert.NoError(t, err)

		assert.Equal(t, roleRequest.Name, receivedRole.Name)
		assert.Equal(t, roleRequest.Description, receivedRole.Description)
		assert.Equal(t, roleRequest.Permissions, receivedRole.Permissions)

		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("update_role_response.json"))
	}

	client := newTestClient(t, http.MethodPut, "/roles/369", testHandler)

	updatedRole, err := client.UpdateRole(context.Background(), roleRequest)
	assert.NoError(t, err)
	assert.Equal(t, int64(369), updatedRole.Id)
	assert.Equal(t, roleRequest.Name, updatedRole.Name)
	assert.Equal(t, roleRequest.Description, updatedRole.Description)
	assert.Equal(t, roleRequest.Permissions, updatedRole.Permissions)
	assert.Equal(t, int64(12), updatedRole.NumberOfUsers)

}

func TestDeleteRole(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status_code":200,"status_string":"OK","request_id":"13fc3dec-b271-4a9b-bd4f-ae22f072e130","status":"Role deleted successfully"}`))
	}

	client := newTestClient(t, http.MethodDelete, "/roles/369", testHandler)

	err := client.DeleteRole(context.Background(), 369)
	assert.NoError(t, err)
}
