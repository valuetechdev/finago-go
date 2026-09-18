package busy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// throttleServer returns 429 for the first n requests, then 200.
func throttleServer(t *testing.T, n int32, header, value string) (*httptest.Server, *atomic.Int32) {
	t.Helper()

	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) <= n {
			if header != "" {
				w.Header().Set(header, value)
			}
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":{"name":"TooManyRequests"}}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	t.Cleanup(srv.Close)

	return srv, &calls
}

func TestRetryOffByDefault(t *testing.T) {
	require := require.New(t)

	srv, calls := throttleServer(t, 1, "Retry-After", "0")

	c := New("test-token", WithURL(srv.URL))
	res, err := c.GetAllUsersWithResponse(t.Context(), &GetAllUsersParams{})
	require.NoError(err)
	require.Equal(http.StatusTooManyRequests, res.StatusCode(), "should not retry unless asked to")
	require.Equal(int32(1), calls.Load(), "should send exactly one request")
}

func TestRetryOnThrottle(t *testing.T) {
	require := require.New(t)

	srv, calls := throttleServer(t, 2, "Retry-After", "0")

	c := New("test-token", WithURL(srv.URL), WithRetry())
	res, err := c.GetAllUsersWithResponse(t.Context(), &GetAllUsersParams{})
	require.NoError(err)
	require.Equal(http.StatusOK, res.StatusCode(), "should retry past the throttling")
	require.Equal(int32(3), calls.Load(), "should send two retries and succeed on the third")
}

func TestRetryGivesUpAfterAttempts(t *testing.T) {
	require := require.New(t)

	srv, calls := throttleServer(t, 100, "Retry-After", "0")

	c := New("test-token", WithURL(srv.URL), WithRetry(WithRetryAttempts(2)))
	res, err := c.GetAllUsersWithResponse(t.Context(), &GetAllUsersParams{})
	require.NoError(err)
	require.Equal(http.StatusTooManyRequests, res.StatusCode(), "should hand back the 429")
	require.Equal(int32(3), calls.Load(), "should send the original request and two retries")
}

func TestRetryRespectsMaxWait(t *testing.T) {
	require := require.New(t)

	srv, calls := throttleServer(t, 100, "Retry-After", "600")

	start := time.Now()
	c := New("test-token", WithURL(srv.URL), WithRetry(WithRetryMaxWait(time.Second)))
	res, err := c.GetAllUsersWithResponse(t.Context(), &GetAllUsersParams{})
	require.NoError(err)
	require.Equal(http.StatusTooManyRequests, res.StatusCode(), "a too-long wait is the caller's call")
	require.Equal(int32(1), calls.Load(), "should not retry")
	require.Less(time.Since(start), 5*time.Second, "should not have slept")
}

func TestRetryReplaysBody(t *testing.T) {
	require := require.New(t)

	var bodies []string
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(body)
		bodies = append(bodies, string(body))

		if calls.Add(1) <= 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"data":{"id":"1"}}`))
	}))
	defer srv.Close()

	c := New("test-token", WithURL(srv.URL), WithRetry())
	_, err := c.CreateClientWithResponse(t.Context(), CreateClientJSONRequestBody{Name: "ACME INC."})
	require.NoError(err)
	require.Len(bodies, 2, "should have sent the request twice")
	require.Equal(bodies[0], bodies[1], "the retried request should carry the same body")
	require.Contains(bodies[0], "ACME INC.")
}

func TestRetryStopsOnContextCancel(t *testing.T) {
	require := require.New(t)

	srv, calls := throttleServer(t, 100, "Retry-After", "30")

	ctx, cancel := context.WithCancel(t.Context())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	c := New("test-token", WithURL(srv.URL), WithRetry())
	start := time.Now()
	_, err := c.GetAllUsersWithResponse(ctx, &GetAllUsersParams{})
	require.Error(err, "a cancelled wait should surface as an error")
	require.Less(time.Since(start), 5*time.Second, "should abandon the wait immediately")
	require.Equal(int32(1), calls.Load())
}

func TestRetryUsesRateLimitReset(t *testing.T) {
	require := require.New(t)

	srv, calls := throttleServer(t, 1, "RateLimit-Reset", "0")

	c := New("test-token", WithURL(srv.URL), WithRetry())
	res, err := c.GetAllUsersWithResponse(t.Context(), &GetAllUsersParams{})
	require.NoError(err)
	require.Equal(http.StatusOK, res.StatusCode(), "should fall back to RateLimit-Reset")
	require.Equal(int32(2), calls.Load())
}

// TestRetryLeavesHttpClientAlone guards against the retrying transport being
// written into a client the caller owns, or worse, http.DefaultClient.
func TestRetryLeavesHttpClientAlone(t *testing.T) {
	require := require.New(t)

	caller := &http.Client{}
	New("test-token", WithHttpClient(caller), WithRetry())
	require.Nil(caller.Transport, "the caller's client should not be modified")
	require.Nil(http.DefaultClient.Transport, "http.DefaultClient should not be modified")
}
