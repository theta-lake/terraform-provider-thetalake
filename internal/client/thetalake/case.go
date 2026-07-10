package thetalake

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type CaseManager struct {
	CaseId    int64     `json:"case_id"`
	CreatedAt time.Time `json:"created_at"`
	Id        int64     `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
	UserEmail string    `json:"user_email"`
	UserId    int64     `json:"user_id"`
	UserName  string    `json:"user_name"`
}

type Case struct {
	CloseDate    *time.Time    `json:"close_date"`
	CreatedAt    time.Time     `json:"created_at"`
	Description  string        `json:"description"`
	Id           int64         `json:"id,omitempty"`
	ManagerIds   []int64       `json:"-"` // not serialized; derived from Managers
	Managers     []CaseManager `json:"managers,omitempty"`
	Name         string        `json:"name"`
	Number       string        `json:"number"`
	OpenDate     time.Time     `json:"open_date"`
	RecordsCount int64         `json:"records_count"`
	Status       string        `json:"status,omitempty"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Visibility   string        `json:"visibility"`
}

// caseUpdateRequest is the body sent to the update endpoint. It deliberately
// excludes close_date: closing/reopening a case is handled exclusively via
// CloseCase/ReopenCase, so this avoids sending "close_date": null on every
// update and inadvertently reopening a closed case.
type caseUpdateRequest struct {
	Description string    `json:"description"`
	Name        string    `json:"name"`
	Number      string    `json:"number"`
	OpenDate    time.Time `json:"open_date"`
	Visibility  string    `json:"visibility"`
}

func addCaseManager(c *Client, caseId int64, userId int64) error {
	endpoint := fmt.Sprintf("/cases/%d/managers", caseId)
	body := struct {
		UserId int64 `json:"user_id"`
	}{UserId: userId}
	return c.doRequest(http.MethodPut, endpoint, body, "", nil)
}

func removeCaseManager(c *Client, caseId int64, userId int64) error {
	endpoint := fmt.Sprintf("/cases/%d/managers/%d", caseId, userId)
	return c.doRequest(http.MethodDelete, endpoint, nil, "", nil)
}

func (c *Client) CreateCase(ctx context.Context, newCase Case) (Case, error) {
	var responseCase Case
	err := c.doRequest(http.MethodPost, "/cases", newCase, "case", &responseCase)
	if err != nil {
		return Case{}, err
	}

	for _, userId := range newCase.ManagerIds {
		if err := addCaseManager(c, responseCase.Id, userId); err != nil {
			return responseCase, err
		}
	}

	responseCase.ManagerIds = newCase.ManagerIds
	return responseCase, nil
}

func (c *Client) GetCaseById(ctx context.Context, caseId int64) (Case, error) {
	var responseCase Case
	endpoint := fmt.Sprintf("/cases/%d", caseId)
	err := c.doRequest(http.MethodGet, endpoint, nil, "case", &responseCase)
	if err != nil {
		return Case{}, err
	}

	for _, manager := range responseCase.Managers {
		responseCase.ManagerIds = append(responseCase.ManagerIds, manager.UserId)
	}

	return responseCase, nil
}

func (c *Client) UpdateCase(ctx context.Context, updatedCase Case) (Case, error) {
	var responseCase Case
	endpoint := fmt.Sprintf("/cases/%d", updatedCase.Id)
	requestBody := caseUpdateRequest{
		Description: updatedCase.Description,
		Name:        updatedCase.Name,
		Number:      updatedCase.Number,
		OpenDate:    updatedCase.OpenDate,
		Visibility:  updatedCase.Visibility,
	}
	err := c.doRequest(http.MethodPut, endpoint, requestBody, "case", &responseCase)
	if err != nil {
		return Case{}, err
	}

	// Get current state to retrieve the existing manager IDs
	currentCase, err := c.GetCaseById(ctx, responseCase.Id)
	if err != nil {
		return responseCase, err
	}

	// Add managers that are in desired state but not in current state
	idsToAdd := diffIdSets(updatedCase.ManagerIds, currentCase.ManagerIds)
	for _, userId := range idsToAdd {
		if err := addCaseManager(c, responseCase.Id, userId); err != nil {
			return responseCase, err
		}
	}

	// Remove managers that are in current state but not in desired state
	idsToRemove := diffIdSets(currentCase.ManagerIds, updatedCase.ManagerIds)
	for _, userId := range idsToRemove {
		if err := removeCaseManager(c, responseCase.Id, userId); err != nil {
			return responseCase, err
		}
	}

	responseCase.CreatedAt = currentCase.CreatedAt
	responseCase.ManagerIds = updatedCase.ManagerIds
	return responseCase, nil
}

func (c *Client) DeleteCase(ctx context.Context, caseId int64) error {
	endpoint := fmt.Sprintf("/cases/%d", caseId)
	return c.doRequest(http.MethodDelete, endpoint, nil, "", nil)
}

func (c *Client) ReopenCase(ctx context.Context, caseId int64) (Case, error) {
	endpoint := fmt.Sprintf("/cases/%d/open", caseId)
	var responseCase Case
	err := c.doRequest(http.MethodPut, endpoint, nil, "case", &responseCase)
	return responseCase, err
}

func (c *Client) CloseCase(ctx context.Context, caseId int64, closeDate time.Time) (Case, error) {
	endpoint := fmt.Sprintf("/cases/%d/close", caseId)
	body := struct {
		CloseDate time.Time `json:"close_date"`
	}{CloseDate: closeDate}

	var responseCase Case
	err := c.doRequest(http.MethodPut, endpoint, body, "case", &responseCase)
	return responseCase, err
}
