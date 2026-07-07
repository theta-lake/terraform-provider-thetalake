package thetalake

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestGetWorkspaceByName(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_workspaces_response.json"))
	}

	client := newTestClient(t, http.MethodGet, "/workspaces", testHandler)

	workspace, err := client.GetWorkspaceByName(context.Background(), "East Region Sales")
	assert.NoError(t, err)

	// Only used by the data source so these are the only attributes it has
	assert.Equal(t, int64(108), workspace.Id)
	assert.Equal(t, "East Region Sales", workspace.Name)
	assert.Equal(t, "East Asia Regional", workspace.Description)
}

func TestGetWorkspaceById(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_workspaces_response.json"))
	}

	client := newTestClient(t, http.MethodGet, "/workspaces", testHandler)

	workspace, err := client.GetWorkspaceById(context.Background(), 108)
	assert.NoError(t, err)

	assert.Equal(t, int64(108), workspace.Id)
	assert.Equal(t, "East Region Sales", workspace.Name)
	assert.Equal(t, "East Asia Regional", workspace.Description)
	assert.Equal(t, "Etc/UTC", workspace.DefaultWorkspaceTimezone)
	assert.Equal(t, false, workspace.AllowAnonymousViaSharedLinks)
	assert.Equal(t, "en", workspace.DefaultTranscriptionLanguage)
	assert.Equal(t, []string{"en", "es"}, workspace.PreferredLanguages)
	assert.Equal(t, int64(730), *workspace.SharedLinksExpirationPeriod)
	assert.Equal(t, []int64{1}, workspace.AnalysisSupervisionSpaceIds)
	assert.Nil(t, workspace.AuditLogRetentionPeriod)
	assert.Nil(t, workspace.DisabledAt)
	assert.Equal(t, false, workspace.CaseManagementManagerAssignment)
	assert.Equal(t, false, workspace.HideAttachmentsFromSearch)
	assert.Equal(t, false, workspace.ReauthenticateOnNetworkChange)
	assert.Equal(t, false, workspace.ShowSystemMessagesInChat)
	assert.Equal(t, 1, len(workspace.Users))
	assert.Equal(t, "jane.doe@example.com", workspace.Users[0].Email)
}

func TestAddUserToWorkspace(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var req addUserToWorkspaceRequest
		err := json.Unmarshal(body, &req)
		assert.NoError(t, err)
		assert.Equal(t, int64(380), req.UserId)

		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("add_user_to_workspace_response.json"))
	}

	client := newTestClient(t, http.MethodPut, "/workspaces/users", testHandler)

	err := client.AddUserToWorkspace(context.Background(), 380)
	assert.NoError(t, err)
}

func TestRemoveUserFromWorkspace(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("remove_user_from_workspace_response.json"))
	}

	client := newTestClient(t, http.MethodDelete, "/workspaces/users/380", testHandler)

	err := client.RemoveUserFromWorkspace(context.Background(), 380)
	assert.NoError(t, err)
}

func TestUpdateWorkspace(t *testing.T) {
	retentionPeriod := int64(365)
	updateRequest := Workspace{
		Id:                              108,
		AllowAnonymousViaSharedLinks:    true,
		AnalysisSupervisionSpaceIds:     []int64{},
		AuditLogRetentionPeriod:         &retentionPeriod,
		CaseManagementManagerAssignment: true,
		DefaultTranscriptionLanguage:    "fr",
		DefaultWorkspaceTimezone:        "America/New_York",
		DeleteOnExpiration:              true,
		HideAttachmentsFromSearch:       true,
		PreferredLanguages:              []string{"fr"},
		ReauthenticateOnNetworkChange:   true,
		ShowSystemMessagesInChat:        true,
		UseNameMatcher:                  false,
		UseOwnerOnlySpaceMatcher:        true,
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedReq workspaceUpdateRequest
		err := json.Unmarshal(body, &receivedReq)
		assert.NoError(t, err)

		assert.Equal(t, true, receivedReq.AllowAnonymousViaSharedLinks)
		assert.Equal(t, []int64{}, receivedReq.AnalysisSupervisionSpaceIds)
		assert.Equal(t, "America/New_York", receivedReq.DefaultWorkspaceTimezone)
		assert.Equal(t, "fr", receivedReq.DefaultTranscriptionLanguage)
		assert.Equal(t, []string{"fr"}, receivedReq.PreferredLanguages)
		assert.NotNil(t, receivedReq.AuditLogRetentionPeriod)
		assert.Equal(t, int64(365), *receivedReq.AuditLogRetentionPeriod)
		assert.Equal(t, true, receivedReq.CaseManagementManagerAssignment)
		assert.Equal(t, true, receivedReq.HideAttachmentsFromSearch)
		assert.Equal(t, true, receivedReq.ReauthenticateOnNetworkChange)
		assert.Equal(t, true, receivedReq.ShowSystemMessagesInChat)

		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("update_workspace_response.json"))
	}

	client := newTestClient(t, http.MethodPut, "/workspaces", testHandler)

	updated, err := client.UpdateWorkspace(context.Background(), updateRequest)
	assert.NoError(t, err)

	assert.Equal(t, int64(108), updated.Id)
	assert.Equal(t, true, updated.AllowAnonymousViaSharedLinks)
	assert.Equal(t, "America/New_York", updated.DefaultWorkspaceTimezone)
	assert.Equal(t, "fr", updated.DefaultTranscriptionLanguage)
	assert.Equal(t, []string{"fr"}, updated.PreferredLanguages)
	assert.Equal(t, int64(365), *updated.AuditLogRetentionPeriod)
	assert.Equal(t, []int64{}, updated.AnalysisSupervisionSpaceIds)
	assert.Equal(t, true, updated.CaseManagementManagerAssignment)
	assert.Equal(t, true, updated.HideAttachmentsFromSearch)
	assert.Equal(t, true, updated.ReauthenticateOnNetworkChange)
	assert.Equal(t, true, updated.ShowSystemMessagesInChat)
}
