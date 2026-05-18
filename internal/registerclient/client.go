// Package registerclient is a thin, single-request HTTP client for
// the workstream-tracker registration endpoint. It posts one
// exact-slug create-or-attach request (design/v0.1-design.md §4)
// and surfaces the observed HTTP status and response body so the
// caller can echo the real receipt. It performs no retry, no
// derivation, and no plan-tree reading.
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
//
// metadata is forwarded verbatim into the request `metadata`
// field when non-empty and omitted when empty. The client stays
// schema-agnostic: it does not know or type the conventionally-
// read `name` key — the caller (the CLI) owns that convention
// (scoping SD4; the milestone schema-loose posture). A metadata
// send failure is not a distinct failure mode: metadata rides
// this single register request, so it surfaces through the same
// returned error the caller already narrates.
func Register(ctx context.Context, serverBase, slug, actor string, metadata json.RawMessage) (Result, error) {
	body, err := json.Marshal(struct {
		ExactSlug string          `json:"exact_slug"`
		Actor     string          `json:"actor"`
		Metadata  json.RawMessage `json:"metadata,omitempty"`
	}{
		ExactSlug: slug,
		Actor:     actor,
		Metadata:  metadata,
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
