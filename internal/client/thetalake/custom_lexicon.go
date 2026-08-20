package thetalake

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// CustomLexicon corresponds to the custom_detection_rule_detail schema
// returned by the create/read/update custom lexicon endpoints.
type CustomLexicon struct {
	AcceptsInput               bool              `json:"accepts_input"`
	AttachmentsEnabled         bool              `json:"attachments_enabled"`
	BoilerplateEnabled         bool              `json:"boilerplate_enabled"`
	ChatroomNameAnalyzed       bool              `json:"chatroom_name_analyzed"`
	CommunicationDirection     []string          `json:"communication_direction"`
	CountProximityByCharacters bool              `json:"count_proximity_by_characters"`
	CreatedAt                  time.Time         `json:"created_at"`
	Description                string            `json:"description"`
	DisabledAt                 *time.Time        `json:"disabled_at"`
	EmailSmartBody             bool              `json:"email_smart_body"`
	EmailSubjectAnalyzed       bool              `json:"email_subject_analyzed"`
	EndDate                    *time.Time        `json:"end_date"`
	FilenameAnalyzed           bool              `json:"filename_analyzed"`
	Id                         int64             `json:"id"`
	MaxParticipants            *int64            `json:"max_participants"`
	MinNumRulesWithHits        *int64            `json:"min_num_rules_with_hits"`
	Name                       string            `json:"name"`
	PolicyIds                  []int64           `json:"policy_ids"`
	RiskType                   string            `json:"risk_type"`
	Rules                      map[string]string `json:"rules"`
	RuleScope                  []string          `json:"rule_scope"`
	StartDate                  *time.Time        `json:"start_date"`
	UpdatedAt                  time.Time         `json:"updated_at"`
}

// CreateCustomLexiconRequest corresponds to the CreateLexiconRequest schema.
type CreateCustomLexiconRequest struct {
	AttachmentsEnabled         *bool    `json:"attachments_enabled,omitempty"`
	BoilerplateEnabled         *bool    `json:"boilerplate_enabled,omitempty"`
	ChatroomNameAnalyzed       *bool    `json:"chatroom_name_analyzed,omitempty"`
	CommunicationDirection     []string `json:"communication_direction,omitempty"`
	CountProximityByCharacters *bool    `json:"count_proximity_by_characters,omitempty"`
	Description                string   `json:"description"`
	EmailSmartBody             *bool    `json:"email_smart_body,omitempty"`
	EmailSubjectAnalyzed       *bool    `json:"email_subject_analyzed,omitempty"`
	EndDate                    *string  `json:"end_date,omitempty"`
	FilenameAnalyzed           *bool    `json:"filename_analyzed,omitempty"`
	MaxParticipants            *int64   `json:"max_participants,omitempty"`
	MinNumRulesWithHits        *int64   `json:"min_num_rules_with_hits,omitempty"`
	Name                       string   `json:"name"`
	Policies                   []int64  `json:"policies,omitempty"`
	RiskType                   string   `json:"risk_type"`
	RuleScope                  []string `json:"rule_scope,omitempty"`
	Rules                      []string `json:"rules"`
	StartDate                  *string  `json:"start_date,omitempty"`
}

// UpdateCustomLexiconRequest corresponds to the UpdateLexiconRequest schema.
// PolicyIds is a pointer to a slice so that an explicit empty slice ([]) can
// be distinguished from an omitted field (nil), per the API's semantics.
type UpdateCustomLexiconRequest struct {
	Description *string  `json:"description,omitempty"`
	Disabled    *bool    `json:"disabled,omitempty"`
	EndDate     *string  `json:"end_date"`
	Name        *string  `json:"name,omitempty"`
	PolicyIds   *[]int64 `json:"policy_ids,omitempty"`
	StartDate   *string  `json:"start_date"`
}

func (c *Client) CreateCustomLexicon(ctx context.Context, request CreateCustomLexiconRequest) (CustomLexicon, error) {
	var response CustomLexicon
	err := c.doRequest(ctx, http.MethodPost, "/analysis/lexicons", request, "lexicon", &response)
	if err != nil {
		return CustomLexicon{}, err
	}
	return response, nil
}

func (c *Client) GetCustomLexiconById(ctx context.Context, id int64) (CustomLexicon, error) {
	var response CustomLexicon
	endpoint := fmt.Sprintf("/analysis/lexicons/%d", id)
	err := c.doRequest(ctx, http.MethodGet, endpoint, nil, "lexicon", &response)
	if err != nil {
		return CustomLexicon{}, err
	}
	return response, nil
}

func (c *Client) UpdateCustomLexicon(ctx context.Context, id int64, request UpdateCustomLexiconRequest) (CustomLexicon, error) {
	var response CustomLexicon
	endpoint := fmt.Sprintf("/analysis/lexicons/%d", id)
	err := c.doRequest(ctx, http.MethodPut, endpoint, request, "lexicon", &response)
	if err != nil {
		return CustomLexicon{}, err
	}
	return response, nil
}
