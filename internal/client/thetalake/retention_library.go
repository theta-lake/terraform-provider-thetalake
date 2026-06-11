package thetalake

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type RetentionLibrary struct {
	CreatedAt                    time.Time `json:"created_at"`
	DatumCount                   int64     `json:"datum_count"`
	DatumSize                    int64     `json:"datum_size"`
	DeleteOnExpiration           bool      `json:"delete_on_expiration"`
	Description                  string    `json:"description"`
	DisplayName                  string    `json:"display_name"`
	ExternalId                   *string   `json:"external_id"`
	Id                           int64     `json:"id"`
	LegalHoldCount               int64     `json:"legal_hold_count"`
	Name                         string    `json:"name"`
	RetainInReview               bool      `json:"retain_in_review"`
	RetentionPeriodDays          int64     `json:"retention_period_days"`
	RetentionPeriodEnabled       bool      `json:"retention_period_enabled"`
	RetentionSummaryText         string    `json:"retention_summary_text"`
	SecCompliantStorageConfirmed bool      `json:"sec_compliant_storage_confirmed"`
	SecCompliantStorageEnabled   bool      `json:"sec_compliant_storage_enabled"`
	StorageAccountId             int64     `json:"storage_account_id"`
	SwrvRuleCount                int64     `json:"swrv_rule_count"`
	UpdatedAt                    time.Time `json:"updated_at"`
}

func (s *Client) GetRetentionLibraryByName(ctx context.Context, name string) (RetentionLibrary, error) {
	var retentionLibraries []RetentionLibrary

	err := s.doRequest(http.MethodGet, "/retention_libraries", nil, "retention_libraries", &retentionLibraries)
	if err != nil {
		return RetentionLibrary{}, err
	}

	for _, retentionLibrary := range retentionLibraries {
		if retentionLibrary.Name == name {
			return retentionLibrary, nil
		}
	}

	return RetentionLibrary{}, errors.New("retention library not found")
}

func (c *Client) CreateRetentionLibrary(ctx context.Context, library RetentionLibrary) (RetentionLibrary, error) {
	var responseLibrary RetentionLibrary
	err := c.doRequest(http.MethodPost, "/retention_libraries", library, "retention_library", &responseLibrary)
	if err != nil {
		return RetentionLibrary{}, err
	}
	return responseLibrary, nil
}

func (c *Client) GetRetentionLibraryById(ctx context.Context, libraryId int64) (RetentionLibrary, error) {
	var responseLibrary RetentionLibrary
	endpoint := fmt.Sprintf("/retention_libraries/%d", libraryId)
	err := c.doRequest(http.MethodGet, endpoint, nil, "retention_library", &responseLibrary)
	if err != nil {
		return RetentionLibrary{}, err
	}
	return responseLibrary, nil
}

func (c *Client) UpdateRetentionLibrary(ctx context.Context, library RetentionLibrary) (RetentionLibrary, error) {
	var responseLibrary RetentionLibrary
	endpoint := fmt.Sprintf("/retention_libraries/%d", library.Id)
	err := c.doRequest(http.MethodPut, endpoint, library, "retention_library", &responseLibrary)
	if err != nil {
		return RetentionLibrary{}, err
	}
	return responseLibrary, nil
}

func (c *Client) DeleteRetentionLibrary(ctx context.Context, libraryId int64) error {
	endpoint := fmt.Sprintf("/retention_libraries/%d", libraryId)
	return c.doRequest(http.MethodDelete, endpoint, nil, "", nil)
}
