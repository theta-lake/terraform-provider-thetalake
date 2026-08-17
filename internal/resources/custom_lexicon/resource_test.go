package customlexicon

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

// newCustomLexiconTestClient creates a thetalake.Client backed by an
// httptest.Server that dispatches all requests to handler.
func newCustomLexiconTestClient(t *testing.T, path string, handler http.HandlerFunc) *thetalake.Client {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"test-access-token","token_type":"Bearer","expires_in":3600}`))
	})
	mux.HandleFunc("/api/v1"+path, handler)

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return thetalake.NewClient(server.URL, "test-client-id", "test-client-secret")
}

// TestUpdateCustomLexiconWithRetry_RetriesOnServiceUnavailable verifies that
// updateCustomLexiconWithRetry retries a 503/Retry-After response (as can
// happen when updating a lexicon shortly after creation, before the create
// has fully settled) rather than surfacing it as a failure.
func TestUpdateCustomLexiconWithRetry_RetriesOnServiceUnavailable(t *testing.T) {
	var attempts atomic.Int32

	client := newCustomLexiconTestClient(t, "/analysis/lexicons/1235", func(w http.ResponseWriter, r *http.Request) {
		if attempts.Add(1) <= 2 {
			// Retry-After: 0 so the test doesn't actually wait.
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status_code": 200,
			"status_string": "OK",
			"request_id": "test-request-id",
			"lexicon": {"id": 1235, "name": "Updated Lexicon"}
		}`))
	})

	r := &customLexiconResource{client: client}

	name := "Updated Lexicon"
	lexicon, err := r.updateCustomLexiconWithRetry(context.Background(), 1235, thetalake.UpdateCustomLexiconRequest{
		Name: &name,
	})
	if err != nil {
		t.Fatalf("expected retries to eventually succeed, got error: %v", err)
	}
	if lexicon.Name != "Updated Lexicon" {
		t.Fatalf("expected updated lexicon name, got %q", lexicon.Name)
	}
	if got := attempts.Load(); got != 3 {
		t.Fatalf("expected 3 attempts (2 failures + 1 success), got %d", got)
	}
}

// TestUpdateCustomLexiconWithRetry_NonRetryableError verifies that a non-503
// error is surfaced immediately without retrying.
func TestUpdateCustomLexiconWithRetry_NonRetryableError(t *testing.T) {
	var attempts atomic.Int32

	client := newCustomLexiconTestClient(t, "/analysis/lexicons/1235", func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"status_code": 400, "status_string": "Bad Request", "message": "invalid request"}`))
	})

	r := &customLexiconResource{client: client}

	name := "Updated Lexicon"
	_, err := r.updateCustomLexiconWithRetry(context.Background(), 1235, thetalake.UpdateCustomLexiconRequest{
		Name: &name,
	})
	if err == nil {
		t.Fatal("expected error to be returned")
	}
	if got := attempts.Load(); got != 1 {
		t.Fatalf("expected exactly 1 attempt for a non-retryable error, got %d", got)
	}
}
