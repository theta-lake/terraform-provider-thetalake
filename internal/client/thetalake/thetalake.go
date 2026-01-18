package thetalake

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	apiServerUrl string
	clientId     string
	clientSecret string
	bearerToken  string
	httpClient   *http.Client
}

type apiResponse struct {
	StatusCode   int             `json:"status_code"`
	StatusString string          `json:"status_string"`
	RequestId    string          `json:"request_id"`
	Data         json.RawMessage `json:"data"`
}

type apiErrorResponse struct {
	apiResponse
	ErrorMessage string `json:"message"`
}

func NewClient(endpoint, clientId, clientSecret string) *Client {
	c := &Client{
		apiServerUrl: strings.TrimRight(endpoint, "/"),
		clientId:     clientId,
		clientSecret: clientSecret,
		httpClient:   &http.Client{},
	}

	// Attempt to fetch an access token immediately. For now, ignore
	// errors and leave bearerToken empty; callers should handle
	// authorization failures from subsequent requests.
	if token, err := c.getToken(context.Background()); err == nil {
		c.bearerToken = token
	}

	return c
}

func (c *Client) doRequest(method, endpoint string, body any, responseObjectName string, responseObject any) error {
	var bodyReader io.Reader

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}

		bodyReader = bytes.NewReader(bodyBytes)
	}

	path, err := url.JoinPath(c.apiServerUrl, "/api/v1", endpoint)
	if err != nil {
		return fmt.Errorf("failed to join URL path: %w", err)
	}

	req, err := http.NewRequest(method, path, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	req.Header.Set("User-Agent", "ThetaLake-Terraform-Provider/0.0.0-dev") // TODO make version dynamic from the build info

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr apiErrorResponse
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
		}
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, apiErr.ErrorMessage)
	}

	var responseMap map[string]json.RawMessage
	err = json.Unmarshal(respBody, &responseMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if responseData, ok := responseMap[responseObjectName]; ok && responseObject != nil {
		if err := json.Unmarshal(responseData, responseObject); err != nil {
			return fmt.Errorf("failed to unmarshal %s from response: %w", responseObjectName, err)
		}
	}

	return nil
}

// getToken retrieves an OAuth2 access token using the
// client_credentials grant type against the Theta Lake token endpoint.
func (c *Client) getToken(ctx context.Context) (string, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", c.clientId)
	form.Set("client_secret", c.clientSecret)

	endpoint := c.apiServerUrl + "/token"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("token request returned non-2xx status: %s", resp.Status)
	}

	var body struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	if body.AccessToken == "" {
		return "", fmt.Errorf("token response did not contain access_token")
	}

	return body.AccessToken, nil
}
