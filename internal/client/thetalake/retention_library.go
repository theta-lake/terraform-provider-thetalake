package thetalake

import (
	"context"
	"errors"
	"net/http"
)

type RetentionLibrary struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (s *Client) GetRetentionLibraryByName(ctx context.Context, name string) (RetentionLibrary, error) {
	var retentionLibraries []RetentionLibrary

	err := s.doRequestWithPagination(http.MethodGet, "/retention_libraries", nil, "retention_libraries", &retentionLibraries, 500)
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
