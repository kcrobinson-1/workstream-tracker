package main

import (
	"bytes"
	"context"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/kcrobinson-1/workstream-tracker/internal/api"
	"github.com/kcrobinson-1/workstream-tracker/internal/db"
)

// startAPIServer spins an httptest server backed by a fresh SQLite
// DB so the register subcommand can be exercised end-to-end
// against the real exact-slug create-or-attach flow.
func startAPIServer(t *testing.T) *httptest.Server {
	t.Helper()
	conn, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	if err := db.Init(context.Background(), conn); err != nil {
		t.Fatalf("db.Init: %v", err)
	}
	r := chi.NewRouter()
	r.Route("/work-instances", api.New(conn).MountRoutes)
	ts := httptest.NewServer(r)
	t.Cleanup(ts.Close)
	return ts
}

func noEnv(string) string { return "" }

var widPattern = regexp.MustCompile(`work_instance_id=(\S+)`)

func TestRegisterCommandSuccessAndIdempotentRepeat(t *testing.T) {
	isolateCaches(t)
	ts := startAPIServer(t)

	run := func() (string, string, int) {
		var out, errb bytes.Buffer
		code := runRegister(
			[]string{"--slug", "demo-root-m1-t2", "--actor", "wst-fixed", "--server", ts.URL},
			noEnv, &out, &errb,
		)
		return out.String(), errb.String(), code
	}

	out1, err1, code1 := run()
	if code1 != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%q", code1, err1)
	}
	if !strings.Contains(out1, "register: ok") || !strings.Contains(out1, "http_status=201") {
		t.Fatalf("success receipt missing real status: %q", out1)
	}
	m1 := widPattern.FindStringSubmatch(out1)
	if m1 == nil {
		t.Fatalf("no work_instance_id echoed: %q", out1)
	}

	out2, _, code2 := run()
	if code2 != 0 {
		t.Fatalf("repeat exit code = %d, want 0", code2)
	}
	m2 := widPattern.FindStringSubmatch(out2)
	if m2 == nil {
		t.Fatalf("repeat echoed no work_instance_id: %q", out2)
	}
	if m1[1] != m2[1] {
		t.Fatalf("idempotent repeat should collapse to same work-instance: %q vs %q", m1[1], m2[1])
	}
}

func TestRegisterCommandServerDownProceeds(t *testing.T) {
	isolateCaches(t)
	var out, errb bytes.Buffer
	code := runRegister(
		[]string{"--slug", "demo-root-m1-t2", "--actor", "a", "--server", "http://127.0.0.1:0"},
		noEnv, &out, &errb,
	)
	if code != 0 {
		t.Fatalf("server-down must exit success; got %d", code)
	}
	if !strings.Contains(errb.String(), "attempt failed") || !strings.Contains(errb.String(), "session proceeds") {
		t.Fatalf("server-down should narrate an explicit, actionable failure: %q", errb.String())
	}
}

func TestRegisterCommandNoSlugSkips(t *testing.T) {
	isolateCaches(t)
	var out, errb bytes.Buffer
	code := runRegister(nil, noEnv, &out, &errb)
	if code != 0 {
		t.Fatalf("no-slug must exit success; got %d", code)
	}
	if !strings.Contains(errb.String(), "no slug supplied") {
		t.Fatalf("no-slug should narrate an explicit skip: %q", errb.String())
	}
	if out.String() != "" {
		t.Fatalf("no-slug should emit no success receipt: %q", out.String())
	}
}

func TestRegisterCommandSlugFromEnv(t *testing.T) {
	isolateCaches(t)
	ts := startAPIServer(t)
	env := func(k string) string {
		switch k {
		case "WST_SLUG":
			return "demo-root-m1-t2"
		case "WST_SERVER":
			return ts.URL
		case "WST_ACTOR":
			return "wst-env-actor"
		}
		return ""
	}
	var out, errb bytes.Buffer
	if code := runRegister(nil, env, &out, &errb); code != 0 {
		t.Fatalf("env-driven register exit = %d, want 0; stderr=%q", code, errb.String())
	}
	if !strings.Contains(out.String(), "actor=wst-env-actor") {
		t.Fatalf("WST_ACTOR not honored: %q", out.String())
	}
}

// isolateCaches points the session cache at a per-test temp dir for
// the duration of the test, so the suite never reads, writes, or
// deletes a live session's real /tmp/wst-* files (the cache is keyed
// by the real worktree root, which during tests is this repo). Every
// cache-touching test must call this first.
func isolateCaches(t *testing.T) {
	t.Helper()
	prev := cacheBaseDir
	cacheBaseDir = t.TempDir()
	t.Cleanup(func() { cacheBaseDir = prev })
}

func sessionActorCachePath(t *testing.T) string {
	t.Helper()
	path, err := wstCachePath("actor")
	if err != nil {
		t.Fatalf("wstCachePath: %v", err)
	}
	return path
}

func TestSessionActorStableAndNamespaced(t *testing.T) {
	isolateCaches(t)

	a, err := sessionActor()
	if err != nil {
		t.Fatalf("sessionActor: %v", err)
	}
	b, err := sessionActor()
	if err != nil {
		t.Fatalf("sessionActor (repeat): %v", err)
	}
	if a != b {
		t.Fatalf("per-session actor must be stable across repeat invocations: %q vs %q", a, b)
	}
	if !strings.HasPrefix(a, "wst-") {
		t.Fatalf("generated actor must be namespaced (never the git user): %q", a)
	}
}

func TestSessionActorRotatesAfterIdleWindow(t *testing.T) {
	isolateCaches(t)
	cachePath := sessionActorCachePath(t)

	first, err := sessionActor()
	if err != nil {
		t.Fatalf("sessionActor: %v", err)
	}

	// Age the cache past the idle window: a genuinely later session
	// must NOT collapse onto the prior session's actor.
	stale := time.Now().Add(-sessionActorIdleWindow - time.Minute)
	if err := os.Chtimes(cachePath, stale, stale); err != nil {
		t.Fatalf("chtimes: %v", err)
	}

	second, err := sessionActor()
	if err != nil {
		t.Fatalf("sessionActor (post-expiry): %v", err)
	}
	if second == first {
		t.Fatalf("actor must rotate after the idle window, got same id %q", second)
	}
	if !strings.HasPrefix(second, "wst-") {
		t.Fatalf("rotated actor must still be namespaced: %q", second)
	}

	// A fresh cache (just written by the call above) is reused, and
	// the reuse slides the window forward.
	third, err := sessionActor()
	if err != nil {
		t.Fatalf("sessionActor (within window): %v", err)
	}
	if third != second {
		t.Fatalf("actor must be stable within the idle window: %q vs %q", second, third)
	}
}
