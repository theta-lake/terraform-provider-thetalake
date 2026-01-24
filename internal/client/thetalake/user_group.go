package thetalake

import (
	"context"
	"errors"
	"net/http"
)

type UserGroup struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
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
