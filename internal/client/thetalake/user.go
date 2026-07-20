package thetalake

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type CurrentWorkspace struct {
	ArchiveOnly bool   `json:"archive_only"`
	Id          int64  `json:"id"`
	Name        string `json:"name"`
}

type SecurityFilter struct {
	SearchId int64  `json:"search_id"`
	Name     string `json:"name"`
}

type user struct {
	CreatedAt             time.Time        `json:"created_at"`
	Id                    int64            `json:"id"`
	Email                 string           `json:"email"`
	Name                  string           `json:"name"`
	Password              string           `json:"password"`
	PasswordConfirmation  string           `json:"password_confirmation"`
	SearchId              int64            `json:"search_id,omitempty"`
	CurrentWorkspace      CurrentWorkspace `json:"current_org_unit"`
	Disabled              bool             `json:"disabled"`
	DisabledAt            *time.Time       `json:"disabled_at"`
	DefaultUserTimezone   *string          `json:"default_user_timezone"`
	ForceSso              bool             `json:"force_sso"`
	HasMultipleWorkspaces bool             `json:"has_multiple_workspaces"`
	LastLogin             *time.Time       `json:"last_login"`
	OtpEnabled            bool             `json:"otp_enabled"`
	OtpEnabledAt          *time.Time       `json:"otp_enabled_at"`
	PasswordChangedAt     *time.Time       `json:"password_changed_at"`
	QueuePaused           bool             `json:"queue_paused"`
	UpdatedAt             *time.Time       `json:"updated_at"`
	SecurityFilter        *SecurityFilter  `json:"security_filter"`
	RoleId                int64            `json:"role_id,omitempty"`
}

// EmbeddedRole is a lightweight role representation used in user API responses.
// It matches the role-2 schema which only contains id and name.
type EmbeddedRole struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type User struct {
	user
	Role EmbeddedRole `json:"role"`
}

type SupervisionSpaceUser struct {
	user
	Role   string `json:"role"`
	UserId int64  `json:"user_id"`
}

func (s *Client) CreateUser(ctx context.Context, user User) (User, error) {
	var responseUser User
	err := s.doRequest(ctx, http.MethodPost, "/users", user, "user", &responseUser)
	if err != nil {
		return User{}, err
	}

	return responseUser, nil
}

func (s *Client) GetUserById(ctx context.Context, userId int64) (User, error) {
	var responseUser User
	endpoint := fmt.Sprintf("/users/%v", userId)

	err := s.doRequest(ctx, http.MethodGet, endpoint, nil, "user", &responseUser)
	if err != nil {
		return User{}, err
	}

	return responseUser, nil
}

func (s *Client) UpdateUser(ctx context.Context, user User) (User, error) {
	var responseUser User
	endpoint := fmt.Sprintf("/users/%v", user.Id)

	err := s.doRequest(ctx, http.MethodPut, endpoint, user, "user", &responseUser)
	if err != nil {
		return User{}, err
	}

	return responseUser, nil
}

func (s *Client) DeleteUser(ctx context.Context, userId int64) error {
	endpoint := fmt.Sprintf("/users/%v", userId)
	err := s.doRequest(ctx, http.MethodDelete, endpoint, nil, "", nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *Client) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var users []User

	err := s.doRequest(ctx, http.MethodGet, "/users", nil, "users", &users)
	if err != nil {
		return User{}, err
	}

	for _, u := range users {
		if u.Email == email {
			return u, nil
		}
	}

	return User{}, fmt.Errorf("user with email %q not found", email)
}
