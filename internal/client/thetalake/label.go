package thetalake

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Label struct {
	BackgroundColor   string     `json:"background_color"`
	CreatedAt         time.Time  `json:"created_at"`
	Hidden            bool       `json:"hidden"`
	Id                int64      `json:"id"`
	LongName          string     `json:"long_name"`
	OrgUnitId         int64      `json:"org_unit_id"`
	ShortName         string     `json:"short_name"`
	TaggedDatumsCount int64      `json:"tagged_datums_count"`
	UpdatedAt         *time.Time `json:"updated_at"`
	UserId            int64      `json:"user_id"`
}

func (c *Client) CreateLabel(ctx context.Context, label Label) (Label, error) {
	var responseLabel Label
	err := c.doRequest(ctx, http.MethodPost, "/labels", label, "label", &responseLabel)
	if err != nil {
		return Label{}, err
	}

	return responseLabel, nil
}

func (c *Client) GetLabelById(ctx context.Context, labelId int64) (Label, error) {
	var responseLabel Label
	endpoint := fmt.Sprintf("/labels/%v", labelId)

	err := c.doRequest(ctx, http.MethodGet, endpoint, nil, "label", &responseLabel)
	if err != nil {
		return Label{}, err
	}

	return responseLabel, nil
}

func (c *Client) UpdateLabel(ctx context.Context, label Label) (Label, error) {
	var responseLabel Label
	endpoint := fmt.Sprintf("/labels/%v", label.Id)

	err := c.doRequest(ctx, http.MethodPut, endpoint, label, "label", &responseLabel)
	if err != nil {
		return Label{}, err
	}

	return responseLabel, nil
}

func (c *Client) DeleteLabel(ctx context.Context, labelId int64) error {
	endpoint := fmt.Sprintf("/labels/%v", labelId)
	err := c.doRequest(ctx, http.MethodDelete, endpoint, nil, "", nil)
	if err != nil {
		return err
	}

	return nil
}
