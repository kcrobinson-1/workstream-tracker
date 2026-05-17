package registerclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRegisterSuccess(t *testing.T) {
	var gotBody map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/work-instances" || r.Method != http.MethodPost {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"wi-123","slug":"epic-m1-t2"}`))
	}))
	defer ts.Close()

	res, err := Register(context.Background(), ts.URL, "epic-m1-t2", "wst-actor")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if res.WorkInstanceID != "wi-123" || res.Slug != "epic-m1-t2" || res.HTTPStatus != http.StatusCreated {
		t.Fatalf("unexpected result: %+v", res)
	}
	if gotBody["exact_slug"] != "epic-m1-t2" || gotBody["actor"] != "wst-actor" {
		t.Fatalf("unexpected request body: %v", gotBody)
	}
}

func TestRegisterTrailingSlashServerBase(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/work-instances" {
			t.Errorf("path not normalized: %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"wi-1","slug":"s"}`))
	}))
	defer ts.Close()

	if _, err := Register(context.Background(), ts.URL+"/", "s", "a"); err != nil {
		t.Fatalf("Register: %v", err)
	}
}

func TestRegisterNonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"root_slug is required"}`))
	}))
	defer ts.Close()

	res, err := Register(context.Background(), ts.URL, "s", "a")
	if err == nil {
		t.Fatal("expected error on non-2xx, got nil")
	}
	if res.HTTPStatus != http.StatusBadRequest {
		t.Fatalf("expected observed status 400, got %d", res.HTTPStatus)
	}
	if !strings.Contains(err.Error(), "400") {
		t.Fatalf("error should carry the real status: %v", err)
	}
}

func TestRegisterMalformedBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`not json`))
	}))
	defer ts.Close()

	if _, err := Register(context.Background(), ts.URL, "s", "a"); err == nil {
		t.Fatal("expected error on malformed body, got nil")
	}
}

func TestRegisterMissingID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"slug":"s"}`))
	}))
	defer ts.Close()

	if _, err := Register(context.Background(), ts.URL, "s", "a"); err == nil {
		t.Fatal("expected error when response omits work-instance id, got nil")
	}
}

func TestRegisterServerUnreachable(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	url := ts.URL
	ts.Close() // nothing is listening now

	if _, err := Register(context.Background(), url, "s", "a"); err == nil {
		t.Fatal("expected error when server is unreachable, got nil")
	}
}

func TestRegisterContextTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	if _, err := Register(ctx, ts.URL, "s", "a"); err == nil {
		t.Fatal("expected error when context deadline exceeded, got nil")
	}
}
