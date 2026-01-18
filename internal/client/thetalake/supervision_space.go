package thetalake

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type SupervisionSpace struct {
	AllParticipants          bool               `json:"all_participants"`
	AllUsers                 bool               `json:"all_users"`
	CanDelete                bool               `json:"can_delete"`
	CanEnableAllParticipants bool               `json:"can_enable_all_participants"`
	CompiledUserList         []User             `json:"compiled_user_list,omitempty"`
	CreatedAt                time.Time          `json:"created_at"`
	Description              string             `json:"description"`
	DirectoryGroupIds        []int64            `json:"directory_group_ids,omitempty"`
	DirectoryGroups          []DirectoryGroup   `json:"directory_groups,omitempty"`
	Disabled                 bool               `json:"disabled"`
	ExternalId               string             `json:"external_id"`
	HardEnforce              bool               `json:"hard_enforce"`
	ID                       int                `json:"id"`
	IntegrationIds           []int64            `json:"integration_ids,omitempty"`
	Integrations             []Integration      `json:"integrations,omitempty"`
	MediaTypeIds             []int64            `json:"media_type_ids,omitempty"`
	MediaTypes               []MediaType        `json:"media_types,omitempty"`
	Name                     string             `json:"name"`
	RetentionLibraryIds      []int64            `json:"retention_library_ids,omitempty"`
	RetentionLibraries       []RetentionLibrary `json:"retention_libraries,omitempty"`
	SupervisionSpacePriority int                `json:"supervision_space_priority"`
	UpdatedAt                time.Time          `json:"updated_at"`
	UserGroupIds             []int64            `json:"user_group_ids,omitempty"`
	UserGroups               []UserGroup        `json:"user_groups,omitempty"`
	UserIds                  []int64            `json:"user_ids,omitempty"`
	Users                    []User             `json:"users,omitempty"`
}

func (s *Client) CreateSupervisionSpace(ctx context.Context, space SupervisionSpace) (SupervisionSpace, error) {
	var responseSpace SupervisionSpace
	err := s.doRequest(http.MethodPost, "/supervision_spaces", space, "supervision_space", &responseSpace)
	if err != nil {
		return SupervisionSpace{}, err
	}

	return responseSpace, nil
}

func (s *Client) GetSupervisionSpaceById(ctx context.Context, spaceId int64) (SupervisionSpace, error) {
	var responseSupervisionSpace SupervisionSpace
	endpoint := fmt.Sprintf("/supervision_spaces/%v", spaceId)

	err := s.doRequest(http.MethodGet, endpoint, nil, "supervision_space", &responseSupervisionSpace)
	if err != nil {
		return SupervisionSpace{}, err
	}

	return responseSupervisionSpace, nil
}

func (s *Client) UpdateSupervisionSpace(ctx context.Context, space SupervisionSpace) (SupervisionSpace, error) {
	var responseSpace SupervisionSpace
	endpoint := fmt.Sprintf("/supervision_spaces/%v", space.ID)

	err := s.doRequest(http.MethodPut, endpoint, space, "supervision_space", &responseSpace)
	if err != nil {
		return SupervisionSpace{}, err
	}

	return responseSpace, nil
}

func (s *Client) DeleteSupervisionSpace(ctx context.Context, spaceId int64) error {
	endpoint := fmt.Sprintf("/supervision_spaces/%v", spaceId)
	err := s.doRequest(http.MethodDelete, endpoint, nil, "", nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *Client) AddDirectoryGroupToSupervisionSpace(ctx context.Context, spaceId int64, groupId int64) error {

	return errors.New("not implemented")
}

func (s *Client) RemoveDirectoryGroupFromSupervisionSpace(ctx context.Context, spaceId int64, groupId int64) error {
	return errors.New("not implemented")
}

func (s *Client) AddUserToSupervisionSpace(ctx context.Context, spaceId int64, userId int64) error {
	return errors.New("not implemented")
}

func (s *Client) RemoveUserFromSupervisionSpace(ctx context.Context, spaceId int64, userId int64) error {
	return errors.New("not implemented")
}

func (s *Client) AddUserGroupToSupervisionSpace(ctx context.Context, spaceId int64, groupId int64) error {
	return errors.New("not implemented")
}

func (s *Client) RemoveUserGroupFromSupervisionSpace(ctx context.Context, spaceId int64, groupId int64) error {
	return errors.New("not implemented")
}
