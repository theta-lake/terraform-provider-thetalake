package thetalake

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type SupervisionSpace struct {
	AllParticipants          bool                   `json:"all_participants"`
	AllUsers                 bool                   `json:"all_users"`
	CanDelete                bool                   `json:"can_delete"`
	CanEnableAllParticipants bool                   `json:"can_enable_all_participants"`
	CreatedAt                time.Time              `json:"created_at"`
	Description              string                 `json:"description"`
	DirectoryGroupIds        []int64                `json:"directory_group_ids"`
	DirectoryGroups          []DirectoryGroup       `json:"directory_groups,omitempty"`
	Disabled                 bool                   `json:"disabled"`
	ExternalId               string                 `json:"external_id"`
	HardEnforce              bool                   `json:"hard_enforce"`
	Id                       int64                  `json:"id"`
	IntegrationIds           []int64                `json:"integration_ids"`
	Integrations             []Integration          `json:"integrations,omitempty"`
	MediaTypeIds             []int64                `json:"media_type_ids"`
	MediaTypes               []MediaType            `json:"media_types,omitempty"`
	Name                     string                 `json:"name"`
	RetentionLibraryIds      []int64                `json:"retention_library_ids"`
	RetentionLibraries       []RetentionLibrary     `json:"retention_libraries,omitempty"`
	SupervisionSpacePriority int                    `json:"supervision_space_priority"`
	UpdatedAt                time.Time              `json:"updated_at"`
	UserGroupIds             []int64                `json:"user_group_ids,omitempty"`
	UserGroups               []UserGroup            `json:"user_groups,omitempty"`
	UserIds                  []int64                `json:"user_ids,omitempty"`
	Users                    []SupervisionSpaceUser `json:"users,omitempty"`
}

func (s *Client) CreateSupervisionSpace(ctx context.Context, space SupervisionSpace) (SupervisionSpace, error) {
	var responseSpace SupervisionSpace
	err := s.doRequest(http.MethodPost, "/supervision_spaces", space, "supervision_space", &responseSpace)
	if err != nil {
		return SupervisionSpace{}, err
	}

	// Add users to the supervision space if any are specified
	if len(space.UserIds) > 0 {
		err = s.doRequest(http.MethodPost, fmt.Sprintf("/supervision_spaces/%d/users", responseSpace.Id), space.UserIds, "supervision_space", &responseSpace)
		if err != nil {
			return responseSpace, err
		}
	}

	// Add user groups to the supervision space if any are specified
	if len(space.UserGroupIds) > 0 {
		err = s.doRequest(http.MethodPost, fmt.Sprintf("/supervision_spaces/%d/user_groups", responseSpace.Id), space.UserGroupIds, "supervision_space", &responseSpace)
		if err != nil {
			return responseSpace, err
		}
	}

	// Ensure UserIds in the response reflect the requested user IDs so
	// Terraform state stays consistent with the plan even if the API
	// response omits or does not populate user_ids.
	responseSpace.UserIds = space.UserIds

	return responseSpace, nil
}

func (s *Client) GetSupervisionSpaceById(ctx context.Context, spaceId int64) (SupervisionSpace, error) {
	var responseSupervisionSpace SupervisionSpace
	endpoint := fmt.Sprintf("/supervision_spaces/%v", spaceId)

	err := s.doRequest(http.MethodGet, endpoint, nil, "supervision_space", &responseSupervisionSpace)
	if err != nil {
		return SupervisionSpace{}, err
	}

	for _, user := range responseSupervisionSpace.Users {
		responseSupervisionSpace.UserIds = append(responseSupervisionSpace.UserIds, user.UserId)
	}

	for _, group := range responseSupervisionSpace.UserGroups {
		responseSupervisionSpace.UserGroupIds = append(responseSupervisionSpace.UserGroupIds, group.Id)
	}

	return responseSupervisionSpace, nil
}

func (s *Client) UpdateSupervisionSpace(ctx context.Context, space SupervisionSpace) (SupervisionSpace, error) {
	var responseSpace SupervisionSpace
	endpoint := fmt.Sprintf("/supervision_spaces/%v", space.Id)

	err := s.doRequest(http.MethodPut, endpoint, space, "supervision_space", &responseSpace)
	if err != nil {
		return SupervisionSpace{}, err
	}

	// Update user groups, adding and removing as necessary
	userGroupsUrl := fmt.Sprintf("/supervision_spaces/%d/user_groups", responseSpace.Id)
	if len(space.UserGroupIds) > 0 {
		err = s.doRequest(http.MethodPost, userGroupsUrl, space.UserGroupIds, "supervision_space", &responseSpace)
		if err != nil {
			return responseSpace, err
		}
	}

	for _, group := range responseSpace.UserGroups {
		responseSpace.UserGroupIds = append(responseSpace.UserGroupIds, group.Id)
	}

	for _, user := range responseSpace.Users {
		responseSpace.UserIds = append(responseSpace.UserIds, user.UserId)
	}

	ids := diffIdSets(responseSpace.UserGroupIds, space.UserGroupIds)
	if len(ids) > 0 {
		err = s.doRequest(http.MethodDelete, userGroupsUrl, ids, "supervision_space", &responseSpace)
		if err != nil {
			return responseSpace, err
		}
	}

	for _, group := range responseSpace.UserGroups {
		responseSpace.UserGroupIds = append(responseSpace.UserGroupIds, group.Id)
	}

	for _, user := range responseSpace.Users {
		responseSpace.UserIds = append(responseSpace.UserIds, user.UserId)
	}

	// Update users, adding and removing as necessary
	usersUrls := fmt.Sprintf("/supervision_spaces/%d/users", responseSpace.Id)

	// Add users
	if len(space.UserIds) > 0 {
		err = s.doRequest(http.MethodPost, usersUrls, space.UserIds, "supervision_space", &responseSpace)
		if err != nil {
			return responseSpace, err
		}
	}

	// Remove users
	ids = diffIdSets(responseSpace.UserIds, space.UserIds)
	if len(ids) > 0 {
		err = s.doRequest(http.MethodDelete, fmt.Sprintf("/supervision_spaces/%d/users", responseSpace.Id), ids, "supervision_space", &responseSpace)
		if err != nil {
			return responseSpace, err
		}
	}

	for _, user := range responseSpace.Users {
		responseSpace.UserIds = append(responseSpace.UserIds, user.UserId)
	}

	for _, group := range responseSpace.UserGroups {
		responseSpace.UserGroupIds = append(responseSpace.UserGroupIds, group.Id)
	}

	// As with Create, ensure UserIds in the response reflect the
	// requested user IDs so downstream state mapping remains stable.
	responseSpace.UserIds = space.UserIds

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

// diffIdSets returns the IDs that are present in setOne but not in setTwo
// Used to determine which associations to remove on update
func diffIdSets(setOne []int64, setTwo []int64) []int64 {
	var idsToRemove []int64

	for _, id := range setOne {
		found := false
		for _, newId := range setTwo {
			if id == newId {
				found = true
				break
			}
		}

		if !found {
			idsToRemove = append(idsToRemove, id)
		}
	}

	return idsToRemove
}
