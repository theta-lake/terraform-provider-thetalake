package thetalake

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type DirectoryGroupIdentity struct {
	CreatedAt   time.Time `json:"created_at"`
	Email       *string   `json:"email"`
	ExternalId  *string   `json:"external_id"`
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	PhoneNumber *string   `json:"phone_number"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DirectoryGroup struct {
	CreatedAt     *time.Time               `json:"created_at,omitempty"`
	Description   *string                  `json:"description,omitempty"`
	ExternalId    *string                  `json:"external_id,omitempty"`
	Id            int64                    `json:"id,omitempty"`
	IdentityCount int64                    `json:"identity_count,omitempty"`
	Identities    []DirectoryGroupIdentity `json:"identities,omitempty"`
	IdentityIds   []int64                  `json:"-"` // not serialized; derived from Identities
	Name          string                   `json:"name"`
	UpdatedAt     *time.Time               `json:"updated_at,omitempty"`
}

func (s *Client) GetDirectoryGroupByName(ctx context.Context, name string) (DirectoryGroup, error) {
	var directoryGroups []DirectoryGroup

	err := s.doRequestWithPagination(ctx, http.MethodGet, "/directory_groups", nil, "directory_groups", &directoryGroups, 500)
	if err != nil {
		return DirectoryGroup{}, err
	}

	for _, directoryGroup := range directoryGroups {
		if directoryGroup.Name == name {
			return directoryGroup, nil
		}
	}

	return DirectoryGroup{}, errors.New("directory group not found")
}

func (s *Client) CreateDirectoryGroup(ctx context.Context, dg DirectoryGroup) (DirectoryGroup, error) {
	var responseDg DirectoryGroup

	err := s.doRequest(ctx, http.MethodPost, "/directory_groups", dg, "directory_group", &responseDg)
	if err != nil {
		return DirectoryGroup{}, err
	}

	if len(dg.IdentityIds) > 0 {
		addEndpoint := fmt.Sprintf("/directory_groups/%d/identities", responseDg.Id)
		err = s.doRequest(ctx, http.MethodPost, addEndpoint, dg.IdentityIds, "identities", nil)
		if err != nil {
			return responseDg, err
		}
	}

	responseDg.IdentityIds = dg.IdentityIds
	return responseDg, nil
}

func (s *Client) GetDirectoryGroupById(ctx context.Context, id int64) (DirectoryGroup, error) {
	var responseDg DirectoryGroup

	endpoint := fmt.Sprintf("/directory_groups/%d", id)
	err := s.doRequest(ctx, http.MethodGet, endpoint, nil, "directory_group", &responseDg)
	if err != nil {
		return DirectoryGroup{}, err
	}

	for _, identity := range responseDg.Identities {
		responseDg.IdentityIds = append(responseDg.IdentityIds, identity.Id)
	}

	return responseDg, nil
}

func (s *Client) UpdateDirectoryGroup(ctx context.Context, dg DirectoryGroup) (DirectoryGroup, error) {
	var responseDg DirectoryGroup

	endpoint := fmt.Sprintf("/directory_groups/%d", dg.Id)
	err := s.doRequest(ctx, http.MethodPut, endpoint, dg, "directory_group", &responseDg)
	if err != nil {
		return DirectoryGroup{}, err
	}

	// Get current state to retrieve timestamps and existing identity IDs
	currentDg, err := s.GetDirectoryGroupById(ctx, responseDg.Id)
	if err != nil {
		return responseDg, err
	}

	// Add identities that are in desired state but not in current state
	idsToAdd := diffIdSets(dg.IdentityIds, currentDg.IdentityIds)
	if len(idsToAdd) > 0 {
		addEndpoint := fmt.Sprintf("/directory_groups/%d/identities", responseDg.Id)
		err = s.doRequest(ctx, http.MethodPost, addEndpoint, idsToAdd, "identities", nil)
		if err != nil {
			return responseDg, err
		}
	}

	// Remove identities that are in current state but not in desired state
	idsToRemove := diffIdSets(currentDg.IdentityIds, dg.IdentityIds)
	for _, identityId := range idsToRemove {
		removeEndpoint := fmt.Sprintf("/directory_groups/%d/identity/%d", responseDg.Id, identityId)
		err = s.doRequest(ctx, http.MethodDelete, removeEndpoint, nil, "", nil)
		if err != nil {
			return responseDg, err
		}
	}

	// Populate timestamps from the GET response and stamp final identity IDs
	responseDg.CreatedAt = currentDg.CreatedAt
	responseDg.UpdatedAt = currentDg.UpdatedAt
	responseDg.IdentityIds = dg.IdentityIds
	return responseDg, nil
}

func (s *Client) DeleteDirectoryGroup(ctx context.Context, id int64) error {
	endpoint := fmt.Sprintf("/directory_groups/%d", id)
	return s.doRequest(ctx, http.MethodDelete, endpoint, nil, "", nil)
}
