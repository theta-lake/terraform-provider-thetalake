package thetalake

import (
	"context"
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

type User struct {
	CreatedAt             time.Time        `json:"created_at"`
	Id                    int64            `json:"id"`
	Email                 string           `json:"email"`
	Name                  string           `json:"name"`
	Password              string           `json:"password"`
	PasswordConfirmation  string           `json:"password_confirmation"`
	RoleId                int64            `json:"role_id"`
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
	Role                  Role             `json:"role"`
	SecurityFilter        *SecurityFilter  `json:"security_filter"`
}

func (s *Client) CreateUser(ctx context.Context, user User) (User, error) {
	// Placeholder for actual API call to create user
	user.Id = 1 // Simulate assigned ID
	user.CreatedAt = time.Now()
	return user, nil
}

func (s *Client) GetUserById(ctx context.Context, userId int64) (User, error) {
	// Placeholder for actual API call to get user by ID
	return User{
		Id:    userId,
		Email: "",
		Name:  "",
		Role:  Role{Id: 1, Name: "User"},
	}, nil
}

func (s *Client) UpdateUser(ctx context.Context, user User) (User, error) {
	// Placeholder for actual API call to update user.
	// Simulate persistence by ensuring an ID is set so Terraform state
	// retains a stable identifier across updates.
	if user.Id == 0 {
		user.Id = 1
	}

	// user.UpdatedAt = ptrToTime(time.Now())
	return user, nil
}

func (s *Client) DeleteUser(ctx context.Context, userId int64) error {
	// Placeholder for actual API call to delete user
	return nil
}
