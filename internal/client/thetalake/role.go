package thetalake

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type Role struct {
	CreatedAt     time.Time  `json:"created_at"`
	Default       bool       `json:"default"`
	Description   string     `json:"description"`
	Id            int64      `json:"id"`
	IsBuiltIn     bool       `json:"is_built_in"`
	Name          string     `json:"name"`
	NumberOfUsers int64      `json:"number_of_users"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

func (s *Client) GetRoleByName(ctx context.Context, name string) (Role, error) {
	var roles []Role

	err := s.doRequest(http.MethodGet, "/roles", nil, "roles", &roles)
	if err != nil {
		return Role{}, err
	}

	for _, role := range roles {
		if role.Name == name {
			return role, nil
		}
	}

	return Role{}, errors.New("role not found")
}
