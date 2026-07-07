package thetalake

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type WorkspaceSupervisionSpace struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type WorkspaceLanguage struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

type WorkspaceUser struct {
	Email string `json:"email"`
	Id    int64  `json:"id"`
	Name  string `json:"name"`
}

type Workspace struct {
	AllowAnonymousViaSharedLinks    bool                        `json:"allow_anonymous_via_shared_links"`
	AnalysisSupervisionSpaceIds     []int64                     `json:"analysis_supervision_space_ids,omitempty"`
	AnalysisSupervisionSpaces       []WorkspaceSupervisionSpace `json:"analysis_supervision_spaces,omitempty"`
	AuditLogRetentionPeriod         *int64                      `json:"audit_log_retention_period"`
	CaseManagementManagerAssignment bool                        `json:"case_management_manager_assignment"`
	CreatedAt                       time.Time                   `json:"created_at"`
	DefaultTranscriptionLanguage    string                      `json:"default_transcription_language"`
	DefaultWorkspaceTimezone        string                      `json:"default_workspace_timezone"`
	DeleteOnExpiration              bool                        `json:"delete_on_expiration"`
	Description                     string                      `json:"description"`
	Disabled                        bool                        `json:"disabled"`
	DisabledAt                      *time.Time                  `json:"disabled_at"`
	HideAttachmentsFromSearch       bool                        `json:"hide_attachments_from_search"`
	Id                              int64                       `json:"id"`
	Name                            string                      `json:"name"`
	PreferredLanguageList           []WorkspaceLanguage         `json:"preferred_language_list,omitempty"`
	PreferredLanguages              []string                    `json:"preferred_languages"`
	ReauthenticateOnNetworkChange   bool                        `json:"reauthenticate_on_network_change"`
	SharedLinksExpirationPeriod     *int64                      `json:"shared_links_expiration_period"`
	ShowSystemMessagesInChat        bool                        `json:"show_system_messages_in_chat"`
	UpdatedAt                       time.Time                   `json:"updated_at"`
	UseNameMatcher                  bool                        `json:"use_name_matcher"`
	UseOwnerOnlySpaceMatcher        bool                        `json:"use_owner_only_space_matcher"`
	Users                           []WorkspaceUser             `json:"users,omitempty"`
}

type workspaceUpdateRequest struct {
	AllowAnonymousViaSharedLinks    bool     `json:"allow_anonymous_via_shared_links"`
	AnalysisSupervisionSpaceIds     []int64  `json:"analysis_supervision_space_ids"`
	AuditLogRetentionPeriod         *int64   `json:"audit_log_retention_period"`
	CaseManagementManagerAssignment bool     `json:"case_management_manager_assignment"`
	DefaultTranscriptionLanguage    string   `json:"default_transcription_language"`
	DefaultWorkspaceTimezone        string   `json:"default_workspace_timezone"`
	DeleteOnExpiration              bool     `json:"delete_on_expiration"`
	HideAttachmentsFromSearch       bool     `json:"hide_attachments_from_search"`
	PreferredLanguages              []string `json:"preferred_languages"`
	ReauthenticateOnNetworkChange   bool     `json:"reauthenticate_on_network_change"`
	SharedLinksExpirationPeriod     *int64   `json:"shared_links_expiration_period"`
	ShowSystemMessagesInChat        bool     `json:"show_system_messages_in_chat"`
	UseNameMatcher                  bool     `json:"use_name_matcher"`
	UseOwnerOnlySpaceMatcher        bool     `json:"use_owner_only_space_matcher"`
}

func (c *Client) GetWorkspaceByName(ctx context.Context, name string) (Workspace, error) {
	var workspaces []Workspace

	err := c.doRequest(http.MethodGet, "/workspaces", nil, "workspaces", &workspaces)
	if err != nil {
		return Workspace{}, err
	}

	for _, w := range workspaces {
		if w.Name == name {
			return w, nil
		}
	}

	return Workspace{}, ErrNotFound
}

func (c *Client) GetWorkspaceById(ctx context.Context, id int64) (Workspace, error) {
	var responseWorkspaces []Workspace
	err := c.doRequest(http.MethodGet, "/workspaces", nil, "workspaces", &responseWorkspaces)
	if err != nil {
		return Workspace{}, err
	}

	var workspace *Workspace
	for _, w := range responseWorkspaces {
		if w.Id == id {
			workspace = &w
			break
		}
	}
	if workspace == nil {
		return Workspace{}, ErrNotFound
	}

	workspace.AnalysisSupervisionSpaceIds = make([]int64, 0, len(workspace.AnalysisSupervisionSpaces))
	for _, ss := range workspace.AnalysisSupervisionSpaces {
		workspace.AnalysisSupervisionSpaceIds = append(workspace.AnalysisSupervisionSpaceIds, ss.Id)
	}
	return *workspace, nil
}

type addUserToWorkspaceRequest struct {
	UserId int64 `json:"user_id"`
}

func (c *Client) AddUserToWorkspace(ctx context.Context, userId int64) error {
	return c.doRequest(http.MethodPut, "/workspaces/users", addUserToWorkspaceRequest{UserId: userId}, "", nil)
}

func (c *Client) RemoveUserFromWorkspace(ctx context.Context, userId int64) error {
	endpoint := fmt.Sprintf("/workspaces/users/%d", userId)
	return c.doRequest(http.MethodDelete, endpoint, nil, "", nil)
}

func (c *Client) UpdateWorkspace(ctx context.Context, workspace Workspace) (Workspace, error) {
	var responseWorkspace Workspace

	updateReq := workspaceUpdateRequest{
		AllowAnonymousViaSharedLinks:    workspace.AllowAnonymousViaSharedLinks,
		AnalysisSupervisionSpaceIds:     workspace.AnalysisSupervisionSpaceIds,
		AuditLogRetentionPeriod:         workspace.AuditLogRetentionPeriod,
		CaseManagementManagerAssignment: workspace.CaseManagementManagerAssignment,
		DefaultTranscriptionLanguage:    workspace.DefaultTranscriptionLanguage,
		DefaultWorkspaceTimezone:        workspace.DefaultWorkspaceTimezone,
		DeleteOnExpiration:              workspace.DeleteOnExpiration,
		HideAttachmentsFromSearch:       workspace.HideAttachmentsFromSearch,
		PreferredLanguages:              workspace.PreferredLanguages,
		ReauthenticateOnNetworkChange:   workspace.ReauthenticateOnNetworkChange,
		SharedLinksExpirationPeriod:     workspace.SharedLinksExpirationPeriod,
		ShowSystemMessagesInChat:        workspace.ShowSystemMessagesInChat,
		UseNameMatcher:                  workspace.UseNameMatcher,
		UseOwnerOnlySpaceMatcher:        workspace.UseOwnerOnlySpaceMatcher,
	}

	if updateReq.AnalysisSupervisionSpaceIds == nil {
		updateReq.AnalysisSupervisionSpaceIds = []int64{}
	}

	if updateReq.PreferredLanguages == nil {
		updateReq.PreferredLanguages = []string{}
	}

	err := c.doRequest(http.MethodPut, "/workspaces", updateReq, "workspace", &responseWorkspace)
	if err != nil {
		return Workspace{}, err
	}

	responseWorkspace.AnalysisSupervisionSpaceIds = make([]int64, 0, len(responseWorkspace.AnalysisSupervisionSpaces))
	for _, ss := range responseWorkspace.AnalysisSupervisionSpaces {
		responseWorkspace.AnalysisSupervisionSpaceIds = append(responseWorkspace.AnalysisSupervisionSpaceIds, ss.Id)
	}

	return responseWorkspace, nil
}
