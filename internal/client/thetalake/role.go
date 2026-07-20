package thetalake

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Role struct {
	CreatedAt     *time.Time `json:"created_at"`
	Default       bool       `json:"default"`
	Description   string     `json:"description"`
	Id            int64      `json:"id"`
	IsBuiltIn     bool       `json:"is_built_in"`
	Name          string     `json:"name"`
	NumberOfUsers int64      `json:"number_of_users"`
	Permissions   []string   `json:"permissions"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

func (s *Client) GetRoleByName(ctx context.Context, name string) (Role, error) {
	var roles []Role

	err := s.doRequest(ctx, http.MethodGet, "/roles", nil, "roles", &roles)
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

func (s *Client) GetRolePermissions(ctx context.Context) ([]string, error) {
	var permissions []string

	err := s.doRequest(ctx, http.MethodGet, "/roles/permissions", nil, "permissions", &permissions)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

type createRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

type updateRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

func (s *Client) CreateRole(ctx context.Context, role Role) (Role, error) {
	var responseRole Role
	req := createRoleRequest{
		Name:        role.Name,
		Description: role.Description,
		Permissions: role.Permissions,
	}
	err := s.doRequest(ctx, http.MethodPost, "/roles", req, "role", &responseRole)
	if err != nil {
		return Role{}, err
	}
	return responseRole, nil
}

func (s *Client) GetRoleById(ctx context.Context, roleId int64) (Role, error) {
	var responseRole Role
	endpoint := fmt.Sprintf("/roles/%d", roleId)
	err := s.doRequest(ctx, http.MethodGet, endpoint, nil, "role", &responseRole)
	if err != nil {
		return Role{}, err
	}
	return responseRole, nil
}

func (s *Client) UpdateRole(ctx context.Context, role Role) (Role, error) {
	var responseRole Role
	endpoint := fmt.Sprintf("/roles/%d", role.Id)
	req := updateRoleRequest{
		Name:        role.Name,
		Description: role.Description,
		Permissions: role.Permissions,
	}
	err := s.doRequest(ctx, http.MethodPut, endpoint, req, "role", &responseRole)
	if err != nil {
		return Role{}, err
	}
	return responseRole, nil
}

func (s *Client) DeleteRole(ctx context.Context, roleId int64) error {
	endpoint := fmt.Sprintf("/roles/%d", roleId)
	return s.doRequest(ctx, http.MethodDelete, endpoint, nil, "", nil)
}
