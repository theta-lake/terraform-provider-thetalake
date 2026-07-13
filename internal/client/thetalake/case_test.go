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

func TestCreateCase(t *testing.T) {
	createHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var received Case
		assert.NoError(t, json.Unmarshal(body, &received))
		assert.Equal(t, "Test Case", received.Name)

		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("create_case_response.json"))
	}

	addManagerHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var received struct {
			UserId int64 `json:"user_id"`
		}
		assert.NoError(t, json.Unmarshal(body, &received))
		assert.Equal(t, int64(422), received.UserId)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status_code": 201, "status_string": "Created"}`))
	}

	client := newTestClientWithRoutes(t,
		testRoute{http.MethodPost, "/cases", createHandler},
		testRoute{http.MethodPut, "/cases/628/managers", addManagerHandler},
	)

	c, err := client.CreateCase(context.Background(), Case{
		Name:       "Test Case",
		Number:     "CASE-001",
		ManagerIds: []int64{422},
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(628), c.Id)
	assert.Equal(t, "Test Case", c.Name)
	assert.Equal(t, []int64{422}, c.ManagerIds)
}

func TestGetCaseById(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_case_by_id_response.json"))
	}

	client := newTestClient(t, http.MethodGet, "/cases/628", handler)

	c, err := client.GetCaseById(context.Background(), 628)
	assert.NoError(t, err)
	assert.Equal(t, int64(628), c.Id)
	assert.Equal(t, "Test Case", c.Name)
	assert.Equal(t, "OPEN", c.Status)
	assert.Equal(t, []int64{422}, c.ManagerIds)
}

func TestUpdateCase(t *testing.T) {
	updateHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer func() { _ = r.Body.Close() }()

		var received map[string]json.RawMessage
		assert.NoError(t, json.Unmarshal(body, &received))
		_, hasCloseDate := received["close_date"]
		assert.False(t, hasCloseDate, "update request body must not include close_date")
		assert.Contains(t, string(body), `"name":"Updated Case"`)

		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("update_case_response.json"))
	}

	getByIdHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_case_by_id_response.json"))
	}

	client := newTestClientWithRoutes(t,
		testRoute{http.MethodPut, "/cases/628", updateHandler},
		testRoute{http.MethodGet, "/cases/628", getByIdHandler},
	)

	c, err := client.UpdateCase(context.Background(), Case{
		Id:         628,
		Name:       "Updated Case",
		Number:     "CASE-001",
		ManagerIds: []int64{422}, // same as in get_by_id response, so no add/remove
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(628), c.Id)
	assert.Equal(t, "Updated Case", c.Name)
	assert.Equal(t, []int64{422}, c.ManagerIds)
}

func TestUpdateCaseManagerDiff(t *testing.T) {
	// Current managers: users 422.
	// Desired managers: users 500.
	// Expected: PUT add manager 500, DELETE remove manager 422.
	var addedUserId int64
	removeWasCalled := false

	updateHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("update_case_response.json"))
	}

	getByIdHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_case_by_id_response.json"))
	}

	addManagerHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var received struct {
			UserId int64 `json:"user_id"`
		}
		json.Unmarshal(body, &received)
		addedUserId = received.UserId

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status_code": 201, "status_string": "Created"}`))
	}

	removeManagerHandler := func(w http.ResponseWriter, r *http.Request) {
		removeWasCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status_code": 200, "status_string": "OK"}`))
	}

	client := newTestClientWithRoutes(t,
		testRoute{http.MethodPut, "/cases/628", updateHandler},
		testRoute{http.MethodGet, "/cases/628", getByIdHandler},
		testRoute{http.MethodPut, "/cases/628/managers", addManagerHandler},
		testRoute{http.MethodDelete, "/cases/628/managers/422", removeManagerHandler},
	)

	c, err := client.UpdateCase(context.Background(), Case{
		Id:         628,
		Name:       "Updated Case",
		ManagerIds: []int64{500},
	})
	assert.NoError(t, err)
	assert.Equal(t, []int64{500}, c.ManagerIds)
	assert.Equal(t, int64(500), addedUserId)
	assert.True(t, removeWasCalled)
}

func TestCloseCase(t *testing.T) {
	closeHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var received struct {
			CloseDate string `json:"close_date"`
		}
		assert.NoError(t, json.Unmarshal(body, &received))
		assert.Equal(t, "2024-02-01T10:00:00Z", received.CloseDate)

		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("close_case_response.json"))
	}

	client := newTestClient(t, http.MethodPut, "/cases/628/close", closeHandler)

	closeDate, err := time.Parse(time.RFC3339, "2024-02-01T10:00:00Z")
	assert.NoError(t, err)

	c, err := client.CloseCase(context.Background(), 628, closeDate)
	assert.NoError(t, err)
	assert.Equal(t, int64(628), c.Id)
	assert.Equal(t, "CLOSED", c.Status)
	assert.NotNil(t, c.CloseDate)
	assert.True(t, closeDate.Equal(*c.CloseDate))
}

func TestReopenCase(t *testing.T) {
	reopenHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("reopen_case_response.json"))
	}

	client := newTestClient(t, http.MethodPut, "/cases/628/open", reopenHandler)

	c, err := client.ReopenCase(context.Background(), 628)
	assert.NoError(t, err)
	assert.Equal(t, int64(628), c.Id)
	assert.Equal(t, "OPEN", c.Status)
	assert.Nil(t, c.CloseDate)
}

func TestDeleteCase(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status_code": 200, "status_string": "OK"}`))
	}

	client := newTestClient(t, http.MethodDelete, "/cases/628", handler)

	err := client.DeleteCase(context.Background(), 628)
	assert.NoError(t, err)
}
