// Package cronitor provides a minimal interface for sending telemetry events to Cronitor.io.
package cronitor

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

// Cronitor is a monitoring service.
// https://cronitor.io/
// https://cronitor.io/docs/telemetry-api
type Cronitor struct {
	HTTPClient  *http.Client
	cronitorURL string
}

// state represents a cronitor event type.
type state string

const (
	run      state = "run"
	complete state = "complete"
	fail     state = "fail"

	defaultClientTimeout = 5 * time.Second
)

var (
	// ErrFailedToSendCronitorEvent is returned when sending a cronitor event fails.
	ErrFailedToSendCronitorEvent = errors.New("failed to send cronitor event")
)

// New creates a new cronitor instance with the given URL.
//
// The URL should be the cronitor URL with the API key.
// Example: https://cronitor.io/p/your-api-key
func New(url string) Cronitor {
	return Cronitor{
		cronitorURL: url,
		HTTPClient:  &http.Client{Timeout: defaultClientTimeout},
	}
}

// Run indicates that the job has started running.
func (c Cronitor) Run(ctx context.Context) error {
	return c.sendEvent(ctx, run)
}

// Complete indicates that the job has completed successfully.
func (c Cronitor) Complete(ctx context.Context) error {
	return c.sendEvent(ctx, complete)
}

// Fail indicates that the job has failed.
func (c Cronitor) Fail(ctx context.Context) error {
	return c.sendEvent(ctx, fail)
}

func (c Cronitor) sendEvent(ctx context.Context, state state) error {
	query := url.Values{}
	query.Add("state", string(state))
	finalURL := c.cronitorURL + "?" + query.Encode()

	resp, err := ctxhttp.Get(ctx, c.HTTPClient, finalURL)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSendCronitorEvent, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %s", ErrFailedToSendCronitorEvent, resp.Status)
	}

	return nil
}
