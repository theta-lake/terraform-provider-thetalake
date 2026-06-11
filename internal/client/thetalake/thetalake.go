package thetalake

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// ErrNotFound is returned when the API responds with a 404 status code.
var ErrNotFound = errors.New("not found")

type Client struct {
	apiServerUrl string
	clientId     string
	clientSecret string
	bearerToken  string
	httpClient   *http.Client
	version      string
}

type apiResponse struct {
	StatusCode   int             `json:"status_code"`
	StatusString string          `json:"status_string"`
	RequestId    string          `json:"request_id"`
	PagingInfo   PagingInfo      `json:"paging_info"`
	Data         json.RawMessage `json:"data"`
}

type apiErrorResponse struct {
	apiResponse
	ErrorMessage string `json:"message"`
}

type PagingInfo struct {
	NextPageToken string `json:"next_page_token"`
}

func NewClient(endpoint, clientId, clientSecret string) *Client {
	c := &Client{
		apiServerUrl: strings.TrimRight(endpoint, "/"),
		clientId:     clientId,
		clientSecret: clientSecret,
		httpClient:   &http.Client{},
		version:      "localdev",
	}

	if token, err := c.getToken(context.Background()); err == nil {
		c.bearerToken = token
	}

	return c
}

func (c *Client) SetVersion(version string) {
	c.version = version
}

func (c *Client) doRequestWithPagination(method, endpoint string, body any, responseObjectName string, responseObject any, max int) error {
	// Generate a slice to hold individual page results
	v := reflect.ValueOf(responseObject)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("responseObject must be pointer to slice")
	}
	sliceValue := v.Elem()         // []T
	sliceType := sliceValue.Type() // type of []T

	nextPageToken := ""
	for {
		var err error
		paginatedEndpoint := fmt.Sprintf("%s?max=%d", endpoint, max) // Only need max on first page
		if nextPageToken != "" {
			paginatedEndpoint = fmt.Sprintf("%s?page_token=%s", endpoint, url.QueryEscape(nextPageToken))
		}

		// Create a fresh *[]T for this page
		pagePtr := reflect.New(sliceType) // *([]T)
		pageSlice := pagePtr.Interface()

		if nextPageToken, err = c.doRequestInner(method, paginatedEndpoint, body, responseObjectName, pageSlice); err != nil {
			return err
		}

		// Append page results into the final slice
		sliceValue.Set(reflect.AppendSlice(sliceValue, pagePtr.Elem()))

		if nextPageToken == "" {
			break
		}
	}

	return nil
}

func (c *Client) doRequest(method, endpoint string, body any, responseObjectName string, responseObject any) error {
	_, err := c.doRequestInner(method, endpoint, body, responseObjectName, responseObject)
	return err
}

// Returns the page token if present, nil otherwise.
func (c *Client) doRequestInner(method, endpoint string, body any, responseObjectName string, responseObject any) (string, error) {
	respBody, err := c.doRequestBytes(method, endpoint, body)
	if err != nil {
		return "", err
	}

	var responseMap map[string]json.RawMessage
	err = json.Unmarshal(respBody, &responseMap)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if responseData, ok := responseMap[responseObjectName]; ok && responseObject != nil {
		if err := json.Unmarshal(responseData, responseObject); err != nil {
			return "", fmt.Errorf("failed to unmarshal %s from response: %w", responseObjectName, err)
		}
	}

	if pagingInfoData, ok := responseMap["paging"]; ok {
		var pagingInfo PagingInfo
		if err := json.Unmarshal(pagingInfoData, &pagingInfo); err != nil {
			return "", fmt.Errorf("failed to unmarshal paging from response: %w", err)
		}

		if pagingInfo.NextPageToken != "" {
			return pagingInfo.NextPageToken, nil
		}
	}

	return "", nil
}

func (c *Client) doRequestBytes(method, endpoint string, body any) ([]byte, error) {
	var bodyBytes []byte

	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// Build the full request URL without escaping query delimiters in endpoint.
	path := c.apiServerUrl + "/api/v1" + endpoint

	var bodyReader io.Reader
	if bodyBytes != nil {
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	req.Header.Set("User-Agent", fmt.Sprintf("ThetaLake-Terraform-Provider/%s", c.version))

	if bodyBytes != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode >= 400 {
		var apiErr apiErrorResponse
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
		}
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, apiErr.ErrorMessage)
	}

	return respBody, nil
}

// getToken retrieves an OAuth2 access token using the client_credentials grant
// type against the Theta Lake token endpoint.
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
	req.Header.Set("User-Agent", fmt.Sprintf("ThetaLake-Terraform-Provider/%s", c.version))

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
