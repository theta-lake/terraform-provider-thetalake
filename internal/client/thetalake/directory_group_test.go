package thetalake

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestGetDirectoryGroupByName(t *testing.T) {
	pageCounter := 1
	handler := func(w http.ResponseWriter, r *http.Request) {
		if pageCounter > 1 {
			expectedPageToken := fmt.Sprintf("page-%d-token", pageCounter)
			actualPageToken := r.URL.Query().Get("page_token")
			assert.Equal(t, expectedPageToken, actualPageToken)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		responseFile := fmt.Sprintf("get_directory_group_page_%d_response.json", pageCounter)
		w.Write(readTestData(responseFile))

		pageCounter++
	}

	client := newTestClient(t, http.MethodGet, "/directory_groups", handler)

	dg, err := client.GetDirectoryGroupByName(context.Background(), "Test directory group")
	assert.NoError(t, err)

	assert.Equal(t, int64(2065), dg.Id)
	assert.Equal(t, "Test directory group", dg.Name)
}
