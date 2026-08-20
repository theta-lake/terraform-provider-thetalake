package thetalake

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	IntegrationTypeGenericJournaling    = "generic_journaling"
	IntegrationTypeGoogleWorkspaceEmail = "google_workspace_email"
	IntegrationTypeThetaLakeApi         = "theta_lake_api"
)

// integrationTypeIdToSlug maps the API's internal integration_type_id to the
// snake_case type token accepted by CreateIntegrationRequest/UpdateIntegrationRequest.
// These IDs are not documented as an enum in the spec; they only appear in examples,
// so they are treated as a secondary, best-effort signal: IntegrationTypeSlug only
// consults this map when the human-readable integration_type name (see
// integrationTypeNameToSlug) does not resolve to a known slug.
var integrationTypeIdToSlug = map[int64]string{
	41: IntegrationTypeGenericJournaling,
	57: IntegrationTypeGoogleWorkspaceEmail,
	80: IntegrationTypeThetaLakeApi,
}

// integrationTypeNameToSlug maps the API's human-readable integration_type display
// name to the snake_case type token. This is preferred over integrationTypeIdToSlug
// by IntegrationTypeSlug because the display name is documented behavior, whereas the
// numeric integration_type_id values are not documented as a stable enum and could
// differ on a server whose IDs don't match the spec's examples.
var integrationTypeNameToSlug = map[string]string{
	"Generic Journaling":     IntegrationTypeGenericJournaling,
	"Google Workspace Email": IntegrationTypeGoogleWorkspaceEmail,
	"Theta Lake API":         IntegrationTypeThetaLakeApi,
}

// IntegrationTypeSlug resolves the snake_case type token for an integration from its
// integration_type display name and integration_type_id. The display name is
// preferred, since the numeric ID values are not documented as a stable enum (see
// integrationTypeIdToSlug); the ID is only consulted as a fallback when the name is
// not recognized. Returns "" if neither is recognized.
func IntegrationTypeSlug(typeId int64, typeName string) string {
	if slug, ok := integrationTypeNameToSlug[typeName]; ok {
		return slug
	}

	return integrationTypeIdToSlug[typeId]
}

// IntegrationOptions is the union of every field across the three anyOf option variants
// (genericJournalingIntegrationOptions, googleWorkspaceEmailIntegrationOptions,
// thetaLakeAPIIntegrationOptions) documented for the integration types this provider
// supports. Which subset is legal for a given integration type is enforced by the
// resource layer, not here.
type IntegrationOptions struct {
	DownloadO365OnedriveLinks  *bool   `json:"download_o365_onedrive_links,omitempty"`
	DownloadSalesforceDoclinks *bool   `json:"download_salesforce_doclinks,omitempty"`
	IndexHeaders               *string `json:"index_headers,omitempty"`
	SenderSpfOverride          *string `json:"sender_spf_override,omitempty"`
	UndeliverableDisabled      *bool   `json:"undeliverable_disabled,omitempty"`
	UndeliverableEmailAddress  *string `json:"undeliverable_email_address,omitempty"`
	UndeliverableEmailPassword *string `json:"undeliverable_email_password,omitempty"`
	UndeliverableEmailPort     *int64  `json:"undeliverable_email_port,omitempty"`
	UndeliverableEmailServer   *string `json:"undeliverable_email_server,omitempty"`
	UndeliverableEmailUser     *string `json:"undeliverable_email_user,omitempty"`
}

// Integration represents a Theta Lake integration. Options/Paused/Type are
// request-only inputs, serialized separately by integrationRequest. The response
// (including the plain GetIntegrationById response) instead carries ServicePaused
// (service_paused) and may carry ServiceParams (service_params); GetIntegrationConfiguration
// is the endpoint that reliably supplies the type-specific options on Read.
type Integration struct {
	CreatedAt         *time.Time          `json:"created_at,omitempty"`
	Id                int64               `json:"id,omitempty"`
	IntegrationGroup  string              `json:"integration_group,omitempty"`
	IntegrationType   string              `json:"integration_type,omitempty"`
	IntegrationTypeId int64               `json:"integration_type_id,omitempty"`
	Name              string              `json:"name"`
	Options           *IntegrationOptions `json:"-"`
	Paused            *bool               `json:"-"`
	ServiceParams     *IntegrationOptions `json:"service_params,omitempty"`
	ServicePaused     bool                `json:"service_paused,omitempty"`
	Status            *string             `json:"status,omitempty"`
	Type              string              `json:"-"`
	UpdatedAt         *time.Time          `json:"updated_at,omitempty"`
}

// integrationRequest is the wire shape for CreateIntegrationRequest/UpdateIntegrationRequest.
// Paused is only meaningful on create; UpdateIntegrationRequest has no paused field in
// the spec, so callers must omit it when building an update request.
type integrationRequest struct {
	Name    string              `json:"name"`
	Options *IntegrationOptions `json:"options,omitempty"`
	Paused  *bool               `json:"paused,omitempty"`
	Type    string              `json:"type"`
}

// IntegrationConfiguration is the response body of GET /integrations/{id}/configuration.
type IntegrationConfiguration struct {
	IntegrationType   string             `json:"integration_type"`
	IntegrationTypeId int64              `json:"integration_type_id"`
	Options           IntegrationOptions `json:"options"`
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

	return Integration{}, fmt.Errorf("integration %q: %w", name, ErrNotFound)
}

func (c *Client) CreateIntegration(ctx context.Context, integration Integration) (Integration, error) {
	var responseIntegration Integration

	request := integrationRequest{
		Name:    integration.Name,
		Options: integration.Options,
		Paused:  integration.Paused,
		Type:    integration.Type,
	}

	err := c.doRequest(ctx, http.MethodPost, "/integrations", request, "integration", &responseIntegration)
	if err != nil {
		return Integration{}, err
	}

	return responseIntegration, nil
}

func (c *Client) GetIntegrationById(ctx context.Context, id int64) (Integration, error) {
	var responseIntegration Integration

	endpoint := fmt.Sprintf("/integrations/%d", id)
	err := c.doRequest(ctx, http.MethodGet, endpoint, nil, "integration", &responseIntegration)
	if err != nil {
		return Integration{}, err
	}

	return responseIntegration, nil
}

func (c *Client) GetIntegrationConfiguration(ctx context.Context, id int64) (IntegrationConfiguration, error) {
	var configuration IntegrationConfiguration

	endpoint := fmt.Sprintf("/integrations/%d/configuration", id)
	err := c.doRequest(ctx, http.MethodGet, endpoint, nil, "configuration", &configuration)
	if err != nil {
		return IntegrationConfiguration{}, err
	}

	return configuration, nil
}

// UpdateIntegration sends a PUT /integrations/{id} request. Note that
// UpdateIntegrationRequest has no paused field in the spec; use PauseIntegration/
// StartIntegration to change the paused state of an existing integration.
func (c *Client) UpdateIntegration(ctx context.Context, integration Integration) (Integration, error) {
	var responseIntegration Integration

	request := integrationRequest{
		Name:    integration.Name,
		Options: integration.Options,
		Type:    integration.Type,
	}

	endpoint := fmt.Sprintf("/integrations/%d", integration.Id)
	err := c.doRequest(ctx, http.MethodPut, endpoint, request, "integration", &responseIntegration)
	if err != nil {
		return Integration{}, err
	}

	return responseIntegration, nil
}

func (c *Client) PauseIntegration(ctx context.Context, id int64) error {
	endpoint := fmt.Sprintf("/integrations/%d/pause", id)
	return c.doRequest(ctx, http.MethodPut, endpoint, nil, "", nil)
}

func (c *Client) StartIntegration(ctx context.Context, id int64) error {
	endpoint := fmt.Sprintf("/integrations/%d/start", id)
	return c.doRequest(ctx, http.MethodPut, endpoint, nil, "", nil)
}

func (c *Client) DeleteIntegration(ctx context.Context, id int64) error {
	endpoint := fmt.Sprintf("/integrations/%d", id)
	return c.doRequest(ctx, http.MethodDelete, endpoint, nil, "", nil)
}
