package busy

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	// DefaultRetryAttempts is the number of retries [WithRetry] performs
	// unless told otherwise.
	DefaultRetryAttempts = 3
	// DefaultRetryMaxWait is the longest [WithRetry] waits for a single retry
	// unless told otherwise. A throttling window longer than this is handed
	// back to the caller as a 429 rather than slept through.
	DefaultRetryMaxWait = 60 * time.Second
)

// RetryOption configures the retrying installed by [WithRetry].
type RetryOption func(*retryTransport)

// WithRetryAttempts sets how many times a throttled request is retried.
// Defaults to [DefaultRetryAttempts].
func WithRetryAttempts(attempts int) RetryOption {
	return func(t *retryTransport) {
		t.attempts = attempts
	}
}

// WithRetryMaxWait caps how long a single retry waits. When Busy asks for
// longer, the 429 is returned to the caller instead.
// Defaults to [DefaultRetryMaxWait].
func WithRetryMaxWait(maxWait time.Duration) RetryOption {
	return func(t *retryTransport) {
		t.maxWait = maxWait
	}
}

// WithRetry retries requests that Busy's rate limiter rejects with a
// 429, waiting as long as the response asks for.
//
// The wait comes from the `Retry-After` header, falling back to
// `RateLimit-Reset` and then to exponential backoff. Only 429 is retried:
// other failures, including 5xx, are returned as they are, and the wait is
// abandoned if the request's context is cancelled.
//
// Retrying is off by default, since it makes a call block for as long as the
// limiter says. It wraps the transport of the client set by [WithHttpClient],
// so the two compose in any order.
func WithRetry(options ...RetryOption) Option {
	return func(c *BusyClient) {
		t := &retryTransport{
			attempts: DefaultRetryAttempts,
			maxWait:  DefaultRetryMaxWait,
		}
		for _, option := range options {
			option(t)
		}
		c.retry = t
	}
}

type retryTransport struct {
	base     http.RoundTripper
	attempts int
	maxWait  time.Duration
}

func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}

	for attempt := 0; ; attempt++ {
		// RoundTrip must not mutate the request it is handed, and a replayed
		// body has to be a fresh reader.
		attemptReq := req.Clone(req.Context())
		if req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return nil, err
			}
			attemptReq.Body = body
		}

		res, err := base.RoundTrip(attemptReq)
		if err != nil || res.StatusCode != http.StatusTooManyRequests {
			return res, err
		}

		// A body we cannot rewind can only be sent once.
		if attempt >= t.attempts || (req.Body != nil && req.GetBody == nil) {
			return res, nil
		}

		wait := retryWait(res, attempt)
		if wait > t.maxWait {
			return res, nil
		}

		// Drain before closing so the connection can be reused.
		_, _ = io.Copy(io.Discard, io.LimitReader(res.Body, 4<<10))
		_ = res.Body.Close()

		if err := sleep(req.Context(), wait); err != nil {
			return nil, err
		}
	}
}

// retryWait reads how long to wait from the response, preferring Busy's
// explicit instruction over a guess.
func retryWait(res *http.Response, attempt int) time.Duration {
	for _, header := range []string{"Retry-After", "RateLimit-Reset"} {
		value := res.Header.Get(header)
		if value == "" {
			continue
		}
		if seconds, err := strconv.Atoi(value); err == nil && seconds >= 0 {
			return time.Duration(seconds) * time.Second
		}
		// Retry-After may also be an HTTP date.
		if date, err := http.ParseTime(value); err == nil {
			if wait := time.Until(date); wait > 0 {
				return wait
			}
			return 0
		}
	}

	return time.Duration(1<<attempt) * time.Second
}

// sleep waits for d, or until the context is done.
func sleep(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}

	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
