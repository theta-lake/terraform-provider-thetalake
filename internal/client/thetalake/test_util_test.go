package thetalake

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

// testRoute describes a single API route to be handled by the test server.
type testRoute struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

// newTestClient is a convenience wrapper for the common case where a test
// only needs to stub a single API route.
func newTestClient(t *testing.T, method, path string, handler http.HandlerFunc) *Client {
	return newTestClientWithRoutes(t, testRoute{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
}

// newTestClientWithRoutes creates a Client whose underlying HTTP server
// can handle multiple API routes. This is useful for testing client
// functions that invoke more than one Theta Lake endpoint.
//
// Multiple methods on the same path are supported; a single mux handler
// is registered per path and dispatches based on the HTTP method.
func newTestClientWithRoutes(t *testing.T, routes ...testRoute) *Client {
	t.Helper() // Mark this function as a test helper, so errors point to the caller

	mux := http.NewServeMux()

	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		// Assert parameters
		assert.Equal(t, http.MethodPost, r.Method)

		assert.Equal(t, "client_credentials", r.FormValue("grant_type"))
		assert.Equal(t, "test-client-id", r.FormValue("client_id"))
		assert.Equal(t, "test-client-secret", r.FormValue("client_secret"))

		// Return test token
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "test-access-token",
			"token_type": "Bearer",
			"expires_in": 3600
			}`))
	})

	// Group routes by path so we can support multiple methods per path.
	byPath := map[string]map[string]http.HandlerFunc{}
	for _, rt := range routes {
		if byPath[rt.Path] == nil {
			byPath[rt.Path] = map[string]http.HandlerFunc{}
		}
		if _, exists := byPath[rt.Path][rt.Method]; exists {
			t.Fatalf("duplicate route registered for %s %s", rt.Method, rt.Path)
		}
		byPath[rt.Path][rt.Method] = rt.Handler
	}

	for path, methods := range byPath {
		p := path
		m := methods
		mux.HandleFunc("/api/v1"+p, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "Bearer test-access-token", r.Header.Get("Authorization"))
			if handler, ok := m[r.Method]; ok {
				handler(w, r)
				return
			}
			t.Fatalf("unexpected request method for path: %s %s", r.Method, r.URL.Path)
		})
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Helper()
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
	})

	ts := httptest.NewServer(mux)

	client := NewClient(ts.URL, "test-client-id", "test-client-secret")

	t.Cleanup(ts.Close) // Ensure server is closed when test ends

	return client
}

func readTestData(filename string) []byte {
	data, err := os.ReadFile("testdata/" + filename)
	if err != nil {
		// This is a test utility; it's acceptable to panic here
		// The test data file should always be present
		panic(fmt.Sprintf("Failed to read test data file %s: %v", filename, err))
	}

	return data
}
