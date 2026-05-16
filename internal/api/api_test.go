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
	"sync"
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
	ts, conn := newTestServer(t)

	cases := []string{
		"",            // empty
		"Foo",         // uppercase
		"foo_bar",     // underscore
		"foo bar",     // space
		"-leading",    // leading hyphen
		"trailing-",   // trailing hyphen
		"foo--double", // double hyphen
		"m1",          // bare position segment is never a root
		"m1-foo",      // position-segment token within the root
		"foo-t1",      // position-segment token within the root
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

	var total int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM work_instances`).Scan(&total); err != nil {
		t.Fatalf("count work_instances: %v", err)
	}
	if total != 0 {
		t.Errorf("invalid root_slug requests created %d rows, want 0", total)
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

func TestRegisterDescendantsConcurrent(t *testing.T) {
	ts, _ := newTestServer(t)

	rootStatus, _ := post(t, ts, "/work-instances", map[string]any{
		"root_slug": "epic-concurrent",
		"node_type": "epic",
		"actor":     "root-agent",
	})
	if rootStatus != http.StatusCreated {
		t.Fatalf("root: status %d", rootStatus)
	}

	const registrations = 10
	emptyParent := ""
	slugs := make(chan string, registrations)
	errs := make(chan string, registrations)

	var wg sync.WaitGroup
	for i := 0; i < registrations; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			status, body := post(t, ts, "/work-instances", map[string]any{
				"root_slug":   "epic-concurrent",
				"parent_path": emptyParent,
				"node_type":   "milestone",
				"actor":       fmt.Sprintf("agent-%d", i),
			})
			if status != http.StatusCreated {
				errs <- fmt.Sprintf("status = %d, body = %+v", status, body)
				return
			}
			slug, ok := body["slug"].(string)
			if !ok {
				errs <- fmt.Sprintf("slug missing or not string: %+v", body)
				return
			}
			slugs <- slug
		}(i)
	}
	wg.Wait()
	close(slugs)
	close(errs)

	for err := range errs {
		t.Error(err)
	}
	if t.Failed() {
		t.FailNow()
	}

	got := make([]string, 0, registrations)
	for slug := range slugs {
		got = append(got, slug)
	}

	seen := map[string]bool{}
	for _, slug := range got {
		if seen[slug] {
			t.Fatalf("duplicate slug %q in %v", slug, got)
		}
		seen[slug] = true
	}
	for i := 1; i <= registrations; i++ {
		want := fmt.Sprintf("epic-concurrent-m%d", i)
		if !seen[want] {
			t.Fatalf("slugs = %v, missing %q", got, want)
		}
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

// eventCount returns the number of events for a work-instance id.
func eventCount(t *testing.T, conn *sql.DB, wid string) int {
	t.Helper()
	var n int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM events WHERE work_instance_id = ?`, wid).Scan(&n); err != nil {
		t.Fatalf("count events: %v", err)
	}
	return n
}

// TestExactSlugCreatesFirstWorkInstance is t2's first-registration
// case: an exact-slug register against a slug with no prior
// work-instance creates the slug's first one. No prior root-create
// is required.
func TestExactSlugCreatesFirstWorkInstance(t *testing.T) {
	ts, conn := newTestServer(t)

	status, body := post(t, ts, "/work-instances", map[string]any{
		"exact_slug": "workstream-tracker-1-0-m1-t1",
		"actor":      "impl-agent",
	})
	if status != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %+v", status, body)
	}
	if body["slug"] != "workstream-tracker-1-0-m1-t1" {
		t.Errorf("slug = %v, want workstream-tracker-1-0-m1-t1", body["slug"])
	}
	wid, ok := body["id"].(string)
	if !ok || wid == "" {
		t.Fatalf("id missing or not a string: %+v", body)
	}

	var count int
	if err := conn.QueryRow(
		`SELECT COUNT(*) FROM work_instances WHERE slug = ?`,
		"workstream-tracker-1-0-m1-t1",
	).Scan(&count); err != nil {
		t.Fatalf("count work_instances: %v", err)
	}
	if count != 1 {
		t.Errorf("work_instances at slug = %d, want 1", count)
	}
}

// TestExactSlugAttachesAdditionalWorkInstance: two different actors
// at one slug both get active work-instances.
func TestExactSlugMultiActorAttach(t *testing.T) {
	ts, conn := newTestServer(t)

	s1, b1 := post(t, ts, "/work-instances", map[string]any{
		"exact_slug": "epic-shared-m1",
		"actor":      "actor-one",
	})
	s2, b2 := post(t, ts, "/work-instances", map[string]any{
		"exact_slug": "epic-shared-m1",
		"actor":      "actor-two",
	})
	if s1 != http.StatusCreated || s2 != http.StatusCreated {
		t.Fatalf("statuses = %d, %d; want 201, 201; bodies = %+v %+v", s1, s2, b1, b2)
	}
	if b1["id"] == b2["id"] {
		t.Errorf("two actors got the same work-instance id %v", b1["id"])
	}

	var active int
	if err := conn.QueryRow(
		`SELECT COUNT(*) FROM work_instances WHERE slug = ? AND state = 'active'`,
		"epic-shared-m1",
	).Scan(&active); err != nil {
		t.Fatalf("count active: %v", err)
	}
	if active != 2 {
		t.Errorf("active work_instances at slug = %d, want 2", active)
	}
}

// TestExactSlugIdempotent: the same actor registering twice against
// one active slug gets one work-instance. The discriminator is the
// returned id equality AND no new event row — a second row created
// but the test only checking 201 would not slip through.
func TestExactSlugIdempotent(t *testing.T) {
	ts, conn := newTestServer(t)

	s1, b1 := post(t, ts, "/work-instances", map[string]any{
		"exact_slug": "epic-idem-m1",
		"actor":      "same-actor",
	})
	if s1 != http.StatusCreated {
		t.Fatalf("first register: status %d, body %+v", s1, b1)
	}
	firstID := b1["id"].(string)
	if got := eventCount(t, conn, firstID); got != 1 {
		t.Fatalf("after first register, event count = %d, want 1", got)
	}

	s2, b2 := post(t, ts, "/work-instances", map[string]any{
		"exact_slug": "epic-idem-m1",
		"actor":      "same-actor",
	})
	if s2 != http.StatusCreated {
		t.Fatalf("second register: status %d, body %+v", s2, b2)
	}
	if b2["id"] != firstID {
		t.Errorf("idempotent register returned id %v, want first id %v", b2["id"], firstID)
	}

	var rowCount int
	if err := conn.QueryRow(
		`SELECT COUNT(*) FROM work_instances WHERE slug = ? AND actor = ?`,
		"epic-idem-m1", "same-actor",
	).Scan(&rowCount); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if rowCount != 1 {
		t.Errorf("work_instances for (slug, actor) = %d, want 1 (idempotent)", rowCount)
	}
	if got := eventCount(t, conn, firstID); got != 1 {
		t.Errorf("event count after idempotent re-register = %d, want 1 (no new register event)", got)
	}
}

// TestExactSlugSerialResume: once the prior work-instance for a
// (slug, actor) pair is terminal, a fresh register for that pair is
// permitted (serial reuse after pause/resume) and mints a new row.
func TestExactSlugSerialResume(t *testing.T) {
	ts, conn := newTestServer(t)

	_, b1 := post(t, ts, "/work-instances", map[string]any{
		"exact_slug": "epic-resume-m1",
		"actor":      "resume-actor",
	})
	firstID := b1["id"].(string)

	// Drive the first work-instance to a terminal state.
	st, _ := post(t, ts, "/work-instances/"+firstID+"/events", map[string]any{"state": "completed"})
	if st != http.StatusCreated {
		t.Fatalf("complete first WI: status %d", st)
	}

	s2, b2 := post(t, ts, "/work-instances", map[string]any{
		"exact_slug": "epic-resume-m1",
		"actor":      "resume-actor",
	})
	if s2 != http.StatusCreated {
		t.Fatalf("resume register: status %d, body %+v", s2, b2)
	}
	if b2["id"] == firstID {
		t.Errorf("serial resume returned the terminal WI id %v; want a new work-instance", firstID)
	}

	var rowCount int
	if err := conn.QueryRow(
		`SELECT COUNT(*) FROM work_instances WHERE slug = ? AND actor = ?`,
		"epic-resume-m1", "resume-actor",
	).Scan(&rowCount); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if rowCount != 2 {
		t.Errorf("work_instances for (slug, actor) = %d, want 2 (serial resume)", rowCount)
	}
}

// TestExactSlugMalformed: a malformed exact_slug is a 400 and
// creates no row. The rejection surfaces as an explicit API error,
// matching the existing register error shape.
func TestExactSlugMalformed(t *testing.T) {
	ts, conn := newTestServer(t)

	cases := []string{
		"Foo-Bar",       // uppercase
		"foo_bar",       // underscore
		"foo bar",       // space
		"-leading",      // leading hyphen
		"trailing-",     // trailing hyphen
		"foo--double",   // double hyphen
		"m1",            // bare position segment, no root word
		"m1-foo",        // root token is a bare position segment
		"epic-m1-x9",    // trailing non-position segment after a position segment
		"epic-m1-foo",   // non-position word after a position segment
		"root-t1-m2",    // out-of-order: t before m
		"root-m1-m2",    // repeated milestone segment
		"root-m1-p1-t1", // out-of-order: t after p
		"root-p1",       // phase without preceding milestone/task
	}
	for _, slug := range cases {
		t.Run(fmt.Sprintf("slug=%q", slug), func(t *testing.T) {
			status, body := post(t, ts, "/work-instances", map[string]any{
				"exact_slug": slug,
				"actor":      "a",
			})
			if status != http.StatusBadRequest {
				t.Errorf("status = %d, want 400 for exact_slug %q", status, slug)
			}
			if _, ok := body["error"].(string); !ok {
				t.Errorf("malformed exact_slug %q: expected an error field, got %+v", slug, body)
			}
		})
	}

	var total int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM work_instances`).Scan(&total); err != nil {
		t.Fatalf("count work_instances: %v", err)
	}
	if total != 0 {
		t.Errorf("malformed exact_slug requests created %d rows, want 0", total)
	}
}

// TestExactSlugBypassesRootConflict: exact-slug registers at a slug
// that already has a root work-instance attach instead of 409. Bare
// root-create still 409s (asserted by TestRegisterRootCollision).
func TestExactSlugBypassesRootConflict(t *testing.T) {
	ts, _ := newTestServer(t)

	s1, _ := post(t, ts, "/work-instances", map[string]any{
		"root_slug": "conflict-root",
		"node_type": "epic",
		"actor":     "a1",
	})
	if s1 != http.StatusCreated {
		t.Fatalf("root create: status %d", s1)
	}

	s2, b2 := post(t, ts, "/work-instances", map[string]any{
		"exact_slug": "conflict-root",
		"actor":      "a2",
	})
	if s2 != http.StatusCreated {
		t.Errorf("exact-slug at existing root: status = %d, want 201 (bypasses root-conflict); body = %+v", s2, b2)
	}
}

// TestV01CallersUnaffected: a request with no exact_slug behaves
// exactly as v0.1 — root-create returns 201 with the supplied slug.
func TestV01CallersUnaffected(t *testing.T) {
	ts, _ := newTestServer(t)

	status, body := post(t, ts, "/work-instances", map[string]any{
		"root_slug": "v01-root",
		"node_type": "task",
		"actor":     "legacy-agent",
	})
	if status != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %+v", status, body)
	}
	if body["slug"] != "v01-root" {
		t.Errorf("slug = %v, want v01-root", body["slug"])
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
