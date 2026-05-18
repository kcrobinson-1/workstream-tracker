package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// captureServer is a minimal /work-instances stand-in that records
// the last request body and returns a success receipt, so the CLI
// metadata plumbing can be asserted without a real DB.
func captureServer(t *testing.T, body *[]byte) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*body, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"wi-1","slug":"s"}`))
	}))
	t.Cleanup(ts.Close)
	return ts
}

func metadataName(t *testing.T, raw []byte) (string, bool) {
	t.Helper()
	var req struct {
		Metadata json.RawMessage `json:"metadata"`
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		t.Fatalf("unmarshal request: %v", err)
	}
	if len(req.Metadata) == 0 {
		return "", false
	}
	var m struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(req.Metadata, &m); err != nil {
		t.Fatalf("unmarshal metadata: %v", err)
	}
	return m.Name, true
}

// TestRegisterCommandNameFlagReportsMetadata: --name builds the
// conventional one-key metadata object the roster reads as the
// label.
func TestRegisterCommandNameFlagReportsMetadata(t *testing.T) {
	var body []byte
	ts := captureServer(t, &body)

	var out, errb bytes.Buffer
	code := runRegister(
		[]string{"--slug", "s", "--actor", "a", "--server", ts.URL, "--name", "Roster work"},
		noEnv, &out, &errb,
	)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%q", code, errb.String())
	}
	if name, ok := metadataName(t, body); !ok || name != "Roster work" {
		t.Fatalf("--name not reported as metadata.name: body=%s", body)
	}
}

// TestRegisterCommandNameFromEnv: WST_NAME is the env fallback.
func TestRegisterCommandNameFromEnv(t *testing.T) {
	var body []byte
	ts := captureServer(t, &body)
	env := func(k string) string {
		if k == "WST_NAME" {
			return "Env session"
		}
		return ""
	}
	var out, errb bytes.Buffer
	code := runRegister(
		[]string{"--slug", "s", "--actor", "a", "--server", ts.URL},
		env, &out, &errb,
	)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%q", code, errb.String())
	}
	if name, ok := metadataName(t, body); !ok || name != "Env session" {
		t.Fatalf("WST_NAME not reported as metadata.name: body=%s", body)
	}
}

// TestRegisterCommandNoNameOmitsMetadata: with no --name / WST_NAME
// the request carries no metadata at all (not an empty name), so a
// v0.2-minimum session still lists and the roster falls back to the
// slug.
func TestRegisterCommandNoNameOmitsMetadata(t *testing.T) {
	var body []byte
	ts := captureServer(t, &body)

	var out, errb bytes.Buffer
	code := runRegister(
		[]string{"--slug", "s", "--actor", "a", "--server", ts.URL},
		noEnv, &out, &errb,
	)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%q", code, errb.String())
	}
	if strings.Contains(string(body), "metadata") {
		t.Fatalf("no-name request must omit metadata, got: %s", body)
	}
}
