package thetalake

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Identity struct {
	Email      *string `json:"email"`
	ExternalId *string `json:"external_id"`
	Id         int64   `json:"id"`
	Name       string  `json:"name"`
}

func (s *Client) GetIdentityByEmail(ctx context.Context, email string) (Identity, error) {
	var identities []Identity

	endpoint := fmt.Sprintf("/identities?query=%s&field_name=email&max=25", url.QueryEscape(email))
	err := s.doRequest(ctx, http.MethodGet, endpoint, nil, "identities", &identities)
	if err != nil {
		return Identity{}, err
	}

	for _, identity := range identities {
		if identity.Email != nil && *identity.Email == email {
			return identity, nil
		}
	}

	return Identity{}, errors.New("identity not found")
}
