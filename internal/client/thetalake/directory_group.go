package thetalake

import (
	"context"
	"errors"
	"net/http"
)

type DirectoryGroup struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (s *Client) GetDirectoryGroupByName(ctx context.Context, name string) (DirectoryGroup, error) {
	var directoryGroups []DirectoryGroup

	err := s.doRequestWithPagination(http.MethodGet, "/directory_groups", nil, "directory_groups", &directoryGroups, 500)
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
