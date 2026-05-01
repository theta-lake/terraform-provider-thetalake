package thetalake

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type UserGroupCategory struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type UserGroupUser struct {
	Id int64 `json:"id"`
}

type UserGroup struct {
	CategoryIds []int64             `json:"category_ids,omitempty"`
	Categories  []UserGroupCategory `json:"categories,omitempty"`
	CreatedAt   *time.Time          `json:"created_at,omitempty"`
	Description string              `json:"description"`
	ExternalId  *string             `json:"external_id"`
	Id          int64               `json:"id"`
	Name        string              `json:"name"`
	UpdatedAt   *time.Time          `json:"updated_at,omitempty"`
	UserIds     []int64             `json:"user_ids,omitempty"`
	Users       []UserGroupUser     `json:"users,omitempty"`
}

func (s *Client) GetUserGroupByName(ctx context.Context, name string) (UserGroup, error) {
	var userGroups []UserGroup

	err := s.doRequest(http.MethodGet, "/user_groups", nil, "user_groups", &userGroups)
	if err != nil {
		return UserGroup{}, err
	}

	for _, userGroup := range userGroups {
		if userGroup.Name == name {
			return userGroup, nil
		}
	}

	return UserGroup{}, errors.New("user group not found")
}

func (s *Client) CreateUserGroup(ctx context.Context, userGroup UserGroup) (UserGroup, error) {
	var responseUserGroup UserGroup
	err := s.doRequest(http.MethodPost, "/user_groups", userGroup, "user_group", &responseUserGroup)
	if err != nil {
		return UserGroup{}, err
	}

	if len(userGroup.UserIds) > 0 {
		addUrl := fmt.Sprintf("/user_groups/%d/add_users", responseUserGroup.Id)
		err = s.doRequest(http.MethodPut, addUrl, userGroup.UserIds, "user_group", &responseUserGroup)
		if err != nil {
			return responseUserGroup, err
		}
	}

	responseUserGroup.UserIds = userGroup.UserIds

	return responseUserGroup, nil
}

func (s *Client) GetUserGroupById(ctx context.Context, userGroupId int64) (UserGroup, error) {
	var responseUserGroup UserGroup
	endpoint := fmt.Sprintf("/user_groups/%v", userGroupId)

	err := s.doRequest(http.MethodGet, endpoint, nil, "user_group", &responseUserGroup)
	if err != nil {
		return UserGroup{}, err
	}

	for _, user := range responseUserGroup.Users {
		responseUserGroup.UserIds = append(responseUserGroup.UserIds, user.Id)
	}

	return responseUserGroup, nil
}

func (s *Client) UpdateUserGroup(ctx context.Context, userGroup UserGroup) (UserGroup, error) {
	var responseUserGroup UserGroup
	endpoint := fmt.Sprintf("/user_groups/%v", userGroup.Id)

	err := s.doRequest(http.MethodPut, endpoint, userGroup, "user_group", &responseUserGroup)
	if err != nil {
		return UserGroup{}, err
	}

	// Collect current user IDs from the PUT response
	var currentUserIds []int64
	for _, user := range responseUserGroup.Users {
		currentUserIds = append(currentUserIds, user.Id)
	}

	addUrl := fmt.Sprintf("/user_groups/%d/add_users", responseUserGroup.Id)
	removeUrl := fmt.Sprintf("/user_groups/%d/remove_users", responseUserGroup.Id)

	// Only add users that are not already in the group
	idsToAdd := findIdsToRemove(userGroup.UserIds, currentUserIds) // Switch the order of arguments to findIdsToRemove to get the correct IDs to add
	if len(idsToAdd) > 0 {
		err = s.doRequest(http.MethodPut, addUrl, idsToAdd, "user_group", &responseUserGroup)
		if err != nil {
			return responseUserGroup, err
		}
	}

	// Remove users that are no longer desired
	idsToRemove := findIdsToRemove(currentUserIds, userGroup.UserIds)
	if len(idsToRemove) > 0 {
		err = s.doRequest(http.MethodPut, removeUrl, idsToRemove, "user_group", &responseUserGroup)
		if err != nil {
			return responseUserGroup, err
		}
	}

	responseUserGroup.UserIds = userGroup.UserIds

	return responseUserGroup, nil
}

func (s *Client) DeleteUserGroup(ctx context.Context, userGroupId int64) error {
	endpoint := fmt.Sprintf("/user_groups/%v", userGroupId)
	err := s.doRequest(http.MethodDelete, endpoint, nil, "", nil)
	if err != nil {
		return err
	}

	return nil
}
