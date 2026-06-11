package thetalake

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestCreateSwrvRule(t *testing.T) {
	description := "My SWRV rule"
	priority := int64(4)
	supervisionSpaceID := int64(10420)
	ruleRequest := SwrvRule{
		Description: &description,
		InputSources: []SwrvRuleInputSource{{
			Id:   2345,
			Type: "integration",
		}},
		Name:               "swrv-example",
		PolicyId:           147,
		Priority:           &priority,
		RetentionLibraryId: 1,
		SupervisionSpaceId: &supervisionSpaceID,
		WorkflowId:         14536,
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedRule swrvRuleRequest
		err := json.Unmarshal(body, &receivedRule)
		assert.NoError(t, err)

		if assert.NotNil(t, receivedRule.Description) {
			assert.Equal(t, description, *receivedRule.Description)
		}
		if assert.Len(t, receivedRule.InputSources, 1) {
			assert.Equal(t, int64(2345), receivedRule.InputSources[0].Id)
			assert.Equal(t, "integration", receivedRule.InputSources[0].Type)
		}
		assert.Equal(t, ruleRequest.Name, receivedRule.Name)
		assert.Equal(t, ruleRequest.PolicyId, receivedRule.PolicyId)
		if assert.NotNil(t, receivedRule.Priority) {
			assert.Equal(t, priority, *receivedRule.Priority)
		}
		assert.Equal(t, ruleRequest.RetentionLibraryId, receivedRule.RetentionLibraryId)
		if assert.NotNil(t, receivedRule.SupervisionSpaceId) {
			assert.Equal(t, supervisionSpaceID, *receivedRule.SupervisionSpaceId)
		}
		assert.Equal(t, ruleRequest.WorkflowId, receivedRule.WorkflowId)

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(readTestData("create_swrv_rule_response.json"))
	}

	client := newTestClient(t, http.MethodPost, "/workflows/swrv_rules", testHandler)

	createdRule, err := client.CreateSwrvRule(context.Background(), ruleRequest)
	assert.NoError(t, err)

	assert.Equal(t, int64(2337), createdRule.Id)
	assert.Equal(t, "swrv-example", createdRule.Name)
	if assert.NotNil(t, createdRule.Description) {
		assert.Equal(t, description, *createdRule.Description)
	}
	if assert.NotNil(t, createdRule.Policy) {
		assert.Equal(t, int64(147), createdRule.Policy.Id)
	}
	if assert.NotNil(t, createdRule.Priority) {
		assert.Equal(t, priority, *createdRule.Priority)
	}
	if assert.NotNil(t, createdRule.RetentionLibrary) {
		assert.Equal(t, int64(1), createdRule.RetentionLibrary.Id)
	}
	if assert.NotNil(t, createdRule.SupervisionSpace) {
		assert.Equal(t, int64(10420), createdRule.SupervisionSpace.Id)
	}
	if assert.NotNil(t, createdRule.Workflow) {
		assert.Equal(t, int64(14536), createdRule.Workflow.Id)
	}
	if assert.Len(t, createdRule.InputSource, 1) {
		assert.NotNil(t, createdRule.InputSource[0].Integration)
	}
}

func TestGetSwrvRuleById(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(readTestData("get_swrv_rule_by_id_response.json"))
	}

	client := newTestClient(t, http.MethodGet, "/workflows/swrv_rules/2337", testHandler)

	retrievedRule, err := client.GetSwrvRuleById(context.Background(), 2337)
	assert.NoError(t, err)

	assert.Equal(t, int64(2337), retrievedRule.Id)
	assert.Equal(t, "swrv-example", retrievedRule.Name)
	if assert.NotNil(t, retrievedRule.Description) {
		assert.Equal(t, "My SWRV rule", *retrievedRule.Description)
	}
	if assert.NotNil(t, retrievedRule.Policy) {
		assert.Equal(t, int64(147), retrievedRule.Policy.Id)
	}
	if assert.NotNil(t, retrievedRule.Priority) {
		assert.Equal(t, int64(4), *retrievedRule.Priority)
	}
	if assert.NotNil(t, retrievedRule.RetentionLibrary) {
		assert.Equal(t, int64(1), retrievedRule.RetentionLibrary.Id)
	}
	if assert.NotNil(t, retrievedRule.SupervisionSpace) {
		assert.Equal(t, int64(10420), retrievedRule.SupervisionSpace.Id)
	}
	if assert.NotNil(t, retrievedRule.Workflow) {
		assert.Equal(t, int64(14536), retrievedRule.Workflow.Id)
	}
	if assert.Len(t, retrievedRule.InputSource, 1) {
		assert.NotNil(t, retrievedRule.InputSource[0].Integration)
	}
}

func TestUpdateSwrvRule(t *testing.T) {
	description := "Updated SWRV rule"
	priority := int64(0)
	supervisionSpaceID := int64(10420)
	ruleRequest := SwrvRule{
		Id:          2337,
		Description: &description,
		InputSources: []SwrvRuleInputSource{{
			Type: "all_uploads",
		}},
		Name:               "swrv-example-updated",
		PolicyId:           147,
		Priority:           &priority,
		RetentionLibraryId: 1,
		SupervisionSpaceId: &supervisionSpaceID,
		WorkflowId:         14536,
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedRule swrvRuleRequest
		err := json.Unmarshal(body, &receivedRule)
		assert.NoError(t, err)

		if assert.NotNil(t, receivedRule.Description) {
			assert.Equal(t, description, *receivedRule.Description)
		}
		if assert.Len(t, receivedRule.InputSources, 1) {
			assert.Equal(t, "all_uploads", receivedRule.InputSources[0].Type)
		}
		assert.Equal(t, ruleRequest.Name, receivedRule.Name)
		assert.Equal(t, ruleRequest.PolicyId, receivedRule.PolicyId)
		if assert.NotNil(t, receivedRule.Priority) {
			assert.Equal(t, priority, *receivedRule.Priority)
		}
		assert.Equal(t, ruleRequest.RetentionLibraryId, receivedRule.RetentionLibraryId)
		if assert.NotNil(t, receivedRule.SupervisionSpaceId) {
			assert.Equal(t, supervisionSpaceID, *receivedRule.SupervisionSpaceId)
		}
		assert.Equal(t, ruleRequest.WorkflowId, receivedRule.WorkflowId)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(readTestData("update_swrv_rule_response.json"))
	}

	client := newTestClient(t, http.MethodPut, "/workflows/swrv_rules/2337", testHandler)

	updatedRule, err := client.UpdateSwrvRule(context.Background(), ruleRequest)
	assert.NoError(t, err)

	assert.Equal(t, int64(2337), updatedRule.Id)
	assert.Equal(t, "swrv-example-updated", updatedRule.Name)
	if assert.NotNil(t, updatedRule.Description) {
		assert.Equal(t, description, *updatedRule.Description)
	}
	if assert.NotNil(t, updatedRule.Priority) {
		assert.Equal(t, priority, *updatedRule.Priority)
	}
	if assert.Len(t, updatedRule.InputSource, 1) {
		assert.Equal(t, "All Uploads", updatedRule.InputSource[0].Name)
	}
}

func TestDeleteSwrvRule(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status_code":200,"status_string":"OK","request_id":"test-request-id","status":"Workflow deleted successfully"}`))
	}

	client := newTestClient(t, http.MethodDelete, "/workflows/swrv_rules/2337", testHandler)

	err := client.DeleteSwrvRule(context.Background(), 2337)
	assert.NoError(t, err)
}
