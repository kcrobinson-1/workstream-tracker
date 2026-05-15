package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/kcrobinson-1/workstream-tracker/internal/db"
)

// newTestServer spins up an httptest.Server backed by a fresh
// SQLite DB in t.TempDir().
func newTestServer(t *testing.T) (*httptest.Server, *sql.DB) {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	conn, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	if err := db.Init(context.Background(), conn); err != nil {
		t.Fatalf("db.Init: %v", err)
	}

	apiServer := New(conn)
	r := chi.NewRouter()
	r.Route("/work-instances", apiServer.MountRoutes)

	ts := httptest.NewServer(r)
	t.Cleanup(ts.Close)

	return ts, conn
}

// post sends a JSON POST and returns the status code and decoded body.
func post(t *testing.T, ts *httptest.Server, path string, body any) (int, map[string]any) {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req, err := http.NewRequest(http.MethodPost, ts.URL+path, &buf)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	out := map[string]any{}
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return resp.StatusCode, out
}

func TestRegisterRoot(t *testing.T) {
	ts, _ := newTestServer(t)

	status, body := post(t, ts, "/work-instances", map[string]any{
		"root_slug": "madrona-feedback",
		"node_type": "epic",
		"actor":     "test-agent",
	})
	if status != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %+v", status, body)
	}
	if body["slug"] != "madrona-feedback" {
		t.Errorf("slug = %v, want madrona-feedback", body["slug"])
	}
	if _, ok := body["id"].(string); !ok {
		t.Errorf("id missing or not a string: %+v", body)
	}
}

func TestRegisterRootCollision(t *testing.T) {
	ts, _ := newTestServer(t)

	first, _ := post(t, ts, "/work-instances", map[string]any{
		"root_slug": "foo",
		"node_type": "epic",
		"actor":     "a1",
	})
	if first != http.StatusCreated {
		t.Fatalf("first registration: status %d", first)
	}

	second, body := post(t, ts, "/work-instances", map[string]any{
		"root_slug": "foo",
		"node_type": "epic",
		"actor":     "a2",
	})
	if second != http.StatusConflict {
		t.Errorf("second registration: status = %d, want 409; body = %+v", second, body)
	}
}

func TestRegisterRootInvalidSlug(t *testing.T) {
	ts, _ := newTestServer(t)

	cases := []string{
		"",            // empty
		"Foo",         // uppercase
		"foo_bar",     // underscore
		"foo bar",     // space
		"-leading",    // leading hyphen
		"trailing-",   // trailing hyphen
		"foo--double", // double hyphen
	}
	for _, slug := range cases {
		t.Run(fmt.Sprintf("slug=%q", slug), func(t *testing.T) {
			status, _ := post(t, ts, "/work-instances", map[string]any{
				"root_slug": slug,
				"node_type": "epic",
				"actor":     "a",
			})
			if status != http.StatusBadRequest {
				t.Errorf("status = %d, want 400 for slug %q", status, slug)
			}
		})
	}
}

func TestRegisterDescendants(t *testing.T) {
	ts, _ := newTestServer(t)

	// Register the root.
	rootStatus, _ := post(t, ts, "/work-instances", map[string]any{
		"root_slug": "epic-x",
		"node_type": "epic",
		"actor":     "a",
	})
	if rootStatus != http.StatusCreated {
		t.Fatalf("root: status %d", rootStatus)
	}

	emptyParent := ""

	// First milestone of epic-x → epic-x-m1.
	_, m1Body := post(t, ts, "/work-instances", map[string]any{
		"root_slug":   "epic-x",
		"parent_path": emptyParent,
		"node_type":   "milestone",
		"actor":       "a",
	})
	if m1Body["slug"] != "epic-x-m1" {
		t.Errorf("m1 slug = %v, want epic-x-m1", m1Body["slug"])
	}

	// Second milestone → epic-x-m2.
	_, m2Body := post(t, ts, "/work-instances", map[string]any{
		"root_slug":   "epic-x",
		"parent_path": emptyParent,
		"node_type":   "milestone",
		"actor":       "a",
	})
	if m2Body["slug"] != "epic-x-m2" {
		t.Errorf("m2 slug = %v, want epic-x-m2", m2Body["slug"])
	}

	// First task under m1 → epic-x-m1-t1.
	m1Path := "m1"
	_, t1Body := post(t, ts, "/work-instances", map[string]any{
		"root_slug":   "epic-x",
		"parent_path": m1Path,
		"node_type":   "task",
		"actor":       "a",
	})
	if t1Body["slug"] != "epic-x-m1-t1" {
		t.Errorf("t1 slug = %v, want epic-x-m1-t1", t1Body["slug"])
	}

	// First phase under m1-t1 → epic-x-m1-t1-p1.
	m1t1Path := "m1-t1"
	_, p1Body := post(t, ts, "/work-instances", map[string]any{
		"root_slug":   "epic-x",
		"parent_path": m1t1Path,
		"node_type":   "phase",
		"actor":       "a",
	})
	if p1Body["slug"] != "epic-x-m1-t1-p1" {
		t.Errorf("p1 slug = %v, want epic-x-m1-t1-p1", p1Body["slug"])
	}

	// Second milestone's first task — m2 numbering is independent.
	m2Path := "m2"
	_, t1OfM2Body := post(t, ts, "/work-instances", map[string]any{
		"root_slug":   "epic-x",
		"parent_path": m2Path,
		"node_type":   "task",
		"actor":       "a",
	})
	if t1OfM2Body["slug"] != "epic-x-m2-t1" {
		t.Errorf("m2-t1 slug = %v, want epic-x-m2-t1", t1OfM2Body["slug"])
	}
}

func TestRegisterDescendantWithoutRoot(t *testing.T) {
	ts, _ := newTestServer(t)

	emptyParent := ""
	status, body := post(t, ts, "/work-instances", map[string]any{
		"root_slug":   "no-such-root",
		"parent_path": emptyParent,
		"node_type":   "milestone",
		"actor":       "a",
	})
	if status != http.StatusBadRequest {
		t.Errorf("status = %d, want 400; body = %+v", status, body)
	}
}

func TestHeartbeat(t *testing.T) {
	ts, conn := newTestServer(t)

	_, root := post(t, ts, "/work-instances", map[string]any{
		"root_slug": "hb-test",
		"node_type": "epic",
		"actor":     "a",
	})
	wid := root["id"].(string)

	// Read the initial last_updated_at.
	var beforeNanos int64
	if err := conn.QueryRow(`SELECT last_updated_at FROM work_instances WHERE id = ?`, wid).Scan(&beforeNanos); err != nil {
		t.Fatalf("read last_updated_at: %v", err)
	}

	status, body := post(t, ts, "/work-instances/"+wid+"/events", map[string]any{})
	if status != http.StatusCreated {
		t.Fatalf("heartbeat: status %d; body = %+v", status, body)
	}

	var afterNanos int64
	if err := conn.QueryRow(`SELECT last_updated_at FROM work_instances WHERE id = ?`, wid).Scan(&afterNanos); err != nil {
		t.Fatalf("read last_updated_at after heartbeat: %v", err)
	}
	if afterNanos < beforeNanos {
		t.Errorf("last_updated_at went backwards: before=%d after=%d", beforeNanos, afterNanos)
	}

	// State should still be 'active'.
	var state string
	if err := conn.QueryRow(`SELECT state FROM work_instances WHERE id = ?`, wid).Scan(&state); err != nil {
		t.Fatalf("read state: %v", err)
	}
	if state != "active" {
		t.Errorf("state = %q, want active", state)
	}
}

func TestStateTransition(t *testing.T) {
	ts, conn := newTestServer(t)

	_, root := post(t, ts, "/work-instances", map[string]any{
		"root_slug": "st-test",
		"node_type": "epic",
		"actor":     "a",
	})
	wid := root["id"].(string)

	status, _ := post(t, ts, "/work-instances/"+wid+"/events", map[string]any{
		"state": "completed",
	})
	if status != http.StatusCreated {
		t.Fatalf("state transition: status %d", status)
	}

	var state string
	var terminalAt sql.NullInt64
	if err := conn.QueryRow(`SELECT state, terminal_at FROM work_instances WHERE id = ?`, wid).Scan(&state, &terminalAt); err != nil {
		t.Fatalf("read after transition: %v", err)
	}
	if state != "completed" {
		t.Errorf("state = %q, want completed", state)
	}
	if !terminalAt.Valid {
		t.Errorf("terminal_at should be set after state transition")
	}
}

func TestStateTransitionInvalidState(t *testing.T) {
	ts, _ := newTestServer(t)

	_, root := post(t, ts, "/work-instances", map[string]any{
		"root_slug": "bad-state-test",
		"node_type": "epic",
		"actor":     "a",
	})
	wid := root["id"].(string)

	status, _ := post(t, ts, "/work-instances/"+wid+"/events", map[string]any{
		"state": "invented-state",
	})
	if status != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", status)
	}
}

func TestEventOnUnknownWorkInstance(t *testing.T) {
	ts, _ := newTestServer(t)

	status, _ := post(t, ts, "/work-instances/nope/events", map[string]any{})
	if status != http.StatusNotFound {
		t.Errorf("status = %d, want 404", status)
	}
}

func TestEventLogPersists(t *testing.T) {
	ts, conn := newTestServer(t)

	_, root := post(t, ts, "/work-instances", map[string]any{
		"root_slug": "log-test",
		"node_type": "epic",
		"actor":     "a",
	})
	wid := root["id"].(string)

	// One register + two heartbeats + one transition = four events.
	post(t, ts, "/work-instances/"+wid+"/events", map[string]any{})
	post(t, ts, "/work-instances/"+wid+"/events", map[string]any{})
	post(t, ts, "/work-instances/"+wid+"/events", map[string]any{"state": "completed"})

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM events WHERE work_instance_id = ?`, wid).Scan(&count); err != nil {
		t.Fatalf("count events: %v", err)
	}
	if count != 4 {
		t.Errorf("event count = %d, want 4", count)
	}
}
