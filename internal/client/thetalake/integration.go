package thetalake

import (
	"context"
	"errors"
	"net/http"
)

type Integration struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (s *Client) GetIntegrationByName(ctx context.Context, name string) (Integration, error) {
	var integrations []Integration

	err := s.doRequest(ctx, http.MethodGet, "/integrations", nil, "integrations", &integrations)
	if err != nil {
		return Integration{}, err
	}

	for _, integration := range integrations {
		if integration.Name == name {
			return integration, nil
		}
	}

	return Integration{}, errors.New("integration not found")
}
