package thetalake

import (
	"context"
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
	// Placeholder for actual API call to get role by name
	return Role{
		Id:   1,
		Name: name,
	}, nil
}
