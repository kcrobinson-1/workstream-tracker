// Package registerclient is a thin, single-request HTTP client for
// the workstream-tracker agent endpoints (design/v0.1-design.md §4):
// the exact-slug create-or-attach registration POST and the
// terminal-state event POST. Each call surfaces the observed HTTP
// status and response body so the caller can echo the real receipt.
// It performs no retry, no derivation, and no plan-tree reading.
package registerclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Result is the observed outcome of a registration attempt. On a
// non-success status HTTPStatus is still populated so the caller
// can echo the real status alongside the error.
type Result struct {
	WorkInstanceID string
	Slug           string
	HTTPStatus     int
}

// Register posts a single exact-slug create-or-attach request to
// serverBase's /work-instances endpoint. The context carries the
// caller's timeout; this client does not retry. A non-2xx status,
// an unreachable server, or a malformed body is returned as an
// error — the caller decides it is non-fatal, not this client.
func Register(ctx context.Context, serverBase, slug, actor string) (Result, error) {
	body, err := json.Marshal(map[string]string{
		"exact_slug": slug,
		"actor":      actor,
	})
	if err != nil {
		return Result{}, fmt.Errorf("marshal request: %w", err)
	}

	url := strings.TrimRight(serverBase, "/") + "/work-instances"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Result{}, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{}, fmt.Errorf("post to %s: %w", url, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return Result{HTTPStatus: resp.StatusCode}, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return Result{HTTPStatus: resp.StatusCode},
			fmt.Errorf("server returned HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var rr struct {
		ID   string `json:"id"`
		Slug string `json:"slug"`
	}
	if err := json.Unmarshal(raw, &rr); err != nil {
		return Result{HTTPStatus: resp.StatusCode},
			fmt.Errorf("malformed response body: %w", err)
	}
	if rr.ID == "" {
		return Result{HTTPStatus: resp.StatusCode},
			fmt.Errorf("response missing work-instance id: %s", strings.TrimSpace(string(raw)))
	}

	return Result{WorkInstanceID: rr.ID, Slug: rr.Slug, HTTPStatus: resp.StatusCode}, nil
}

// EventResult is the observed outcome of a terminal-state event
// POST. On a non-success status HTTPStatus is still populated so the
// caller can echo the real status alongside the error.
type EventResult struct {
	EventID    string
	HTTPStatus int
}

// RecordState posts a single terminal-state transition
// ("completed" or "abandoned") to serverBase's
// /work-instances/{id}/events endpoint. Symmetric to Register: the
// context carries the caller's timeout, there is no retry, and a
// non-2xx status, an unreachable server, or a malformed body is
// returned as an error for the caller to decide is non-fatal.
func RecordState(ctx context.Context, serverBase, id, state string) (EventResult, error) {
	body, err := json.Marshal(map[string]string{"state": state})
	if err != nil {
		return EventResult{}, fmt.Errorf("marshal request: %w", err)
	}

	url := strings.TrimRight(serverBase, "/") + "/work-instances/" + id + "/events"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return EventResult{}, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return EventResult{}, fmt.Errorf("post to %s: %w", url, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return EventResult{HTTPStatus: resp.StatusCode}, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return EventResult{HTTPStatus: resp.StatusCode},
			fmt.Errorf("server returned HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var er struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(raw, &er); err != nil {
		return EventResult{HTTPStatus: resp.StatusCode},
			fmt.Errorf("malformed response body: %w", err)
	}
	if er.ID == "" {
		return EventResult{HTTPStatus: resp.StatusCode},
			fmt.Errorf("response missing event id: %s", strings.TrimSpace(string(raw)))
	}

	return EventResult{EventID: er.ID, HTTPStatus: resp.StatusCode}, nil
}
