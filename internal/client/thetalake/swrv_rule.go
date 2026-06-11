package thetalake

import (
	"context"
	"fmt"
	"net/http"
)

type SwrvInputSource struct {
	Name        string                      `json:"name,omitempty"`
	Integration *SwrvInputSourceIntegration `json:"integration,omitempty"`
}

type SwrvInputSourceIntegration struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type SwrvRuleInputSource struct {
	Id   int64  `json:"id,omitempty"`
	Type string `json:"type"`
}

type SwrvPolicy struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type SwrvRetentionLibrary struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type SwrvSearch struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type SwrvSupervisionSpace struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type SwrvWorkflow struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type SwrvRule struct {
	DefaultRule        bool                  `json:"default_rule"`
	Description        *string               `json:"description"`
	Id                 int64                 `json:"id"`
	InputSource        []SwrvInputSource     `json:"input_source,omitempty"`
	InputSources       []SwrvRuleInputSource `json:"input_sources,omitempty"`
	IsBuiltIn          bool                  `json:"is_built_in"`
	Name               string                `json:"name"`
	Policy             *SwrvPolicy           `json:"policy,omitempty"`
	PolicyId           int64                 `json:"policy_id,omitempty"`
	Priority           *int64                `json:"priority,omitempty"`
	RetentionLibrary   *SwrvRetentionLibrary `json:"retention_library,omitempty"`
	RetentionLibraryId int64                 `json:"retention_library_id,omitempty"`
	Search             *SwrvSearch           `json:"search,omitempty"`
	SupervisionSpace   *SwrvSupervisionSpace `json:"supervision_space,omitempty"`
	SupervisionSpaceId *int64                `json:"supervision_space_id,omitempty"`
	Workflow           *SwrvWorkflow         `json:"workflow,omitempty"`
	WorkflowId         int64                 `json:"workflow_id,omitempty"`
}

type swrvRuleRequest struct {
	Description        *string               `json:"description"`
	InputSources       []SwrvRuleInputSource `json:"input_sources"`
	Name               string                `json:"name"`
	PolicyId           int64                 `json:"policy_id"`
	Priority           *int64                `json:"priority,omitempty"`
	RetentionLibraryId int64                 `json:"retention_library_id"`
	SupervisionSpaceId *int64                `json:"supervision_space_id,omitempty"`
	WorkflowId         int64                 `json:"workflow_id"`
}

func (c *Client) CreateSwrvRule(ctx context.Context, rule SwrvRule) (SwrvRule, error) {
	return c.doSwrvRuleWrite(http.MethodPost, "/workflows/swrv_rules", swrvRuleRequestFromRule(rule))
}

func (c *Client) GetSwrvRuleById(ctx context.Context, ruleId int64) (SwrvRule, error) {
	var responseRule SwrvRule
	endpoint := fmt.Sprintf("/workflows/swrv_rules/%d", ruleId)
	err := c.doRequest(http.MethodGet, endpoint, nil, "swrv_rule", &responseRule)
	if err != nil {
		return SwrvRule{}, err
	}

	return responseRule, nil
}

func (c *Client) UpdateSwrvRule(ctx context.Context, rule SwrvRule) (SwrvRule, error) {
	endpoint := fmt.Sprintf("/workflows/swrv_rules/%d", rule.Id)
	return c.doSwrvRuleWrite(http.MethodPut, endpoint, swrvRuleRequestFromRule(rule))
}

func (c *Client) doSwrvRuleWrite(method, endpoint string, body any) (SwrvRule, error) {
	var responseRule SwrvRule
	err := c.doRequest(method, endpoint, body, "swrv_rule", &responseRule)
	if err != nil {
		return SwrvRule{}, err
	}

	return responseRule, nil
}

func swrvRuleRequestFromRule(rule SwrvRule) swrvRuleRequest {
	return swrvRuleRequest{
		Description:        rule.Description,
		InputSources:       rule.InputSources,
		Name:               rule.Name,
		PolicyId:           rule.PolicyId,
		Priority:           rule.Priority,
		RetentionLibraryId: rule.RetentionLibraryId,
		SupervisionSpaceId: rule.SupervisionSpaceId,
		WorkflowId:         rule.WorkflowId,
	}
}

func (c *Client) DeleteSwrvRule(ctx context.Context, ruleId int64) error {
	endpoint := fmt.Sprintf("/workflows/swrv_rules/%d", ruleId)
	return c.doRequest(http.MethodDelete, endpoint, nil, "", nil)
}
