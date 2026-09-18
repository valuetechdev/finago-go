package busy

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// requireOnline skips tests that talk to the live Busy API. `go test -short`
// must pass offline and without credentials, and a plain `go test` without
// TFSO_BUSY_TOKEN in the environment skips too rather than failing on a 401 —
// the secret only resolves under `fnox exec` (see `mise run test:full`).
func requireOnline(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping test that requires the Busy API")
	}
	if os.Getenv("TFSO_BUSY_TOKEN") == "" {
		t.Skip("skipping test that requires the Busy API: TFSO_BUSY_TOKEN is not set, run via `mise run test:full`")
	}
}

// getClient builds a client from TFSO_BUSY_TOKEN. Busy tokens are bound to one
// environment, so a demo token needs TFSO_BUSY_BASE_URL set to [DemoBaseUrl];
// without it the tests talk to production. Keep the live tests read-only.
func getClient() *BusyClient {
	options := []Option{}
	if baseUrl := os.Getenv("TFSO_BUSY_BASE_URL"); baseUrl != "" {
		options = append(options, WithBaseUrl(baseUrl))
	}
	return New(os.Getenv("TFSO_BUSY_TOKEN"), options...)
}

func TestClientDefaults(t *testing.T) {
	require := require.New(t)

	c := New("test-token")
	require.NotNil(c)
	require.Equal(ProdBaseUrl, c.baseUrl, "should default to production")

	require.Equal(DemoBaseUrl, New("test-token", WithDemo()).baseUrl)
	require.Equal("https://example.org", New("test-token", WithBaseUrl("https://example.org")).baseUrl)
}

func TestClientSendsBearerToken(t *testing.T) {
	require := require.New(t)

	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer srv.Close()

	c := New("test-token", WithBaseUrl(srv.URL))
	_, err := c.GetAllUsersWithResponse(t.Context(), &GetAllUsersParams{})
	require.NoError(err)
	require.Equal("Bearer test-token", gotAuth, "token should be sent as a bearer token")
}

func TestClient(t *testing.T) {
	requireOnline(t)
	require := require.New(t)

	c := getClient()
	require.NotNil(c)

	res, err := c.GetAllUsersWithResponse(t.Context(), &GetAllUsersParams{})
	require.NoError(err, "GetAllUsers should not error")
	require.Equal(http.StatusOK, res.StatusCode(), "GetAllUsers status should be OK", string(res.Body))
	require.NotNil(res.JSON200, "GetAllUsers JSON200 should not be nil")
}
