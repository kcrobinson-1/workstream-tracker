package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func wiCachePath(t *testing.T) string {
	t.Helper()
	root, err := worktreeRoot()
	if err != nil {
		t.Fatalf("worktreeRoot: %v", err)
	}
	sum := sha256.Sum256([]byte(root))
	return filepath.Join(os.TempDir(), "wst-wi-"+hex.EncodeToString(sum[:8])+".id")
}

// pinnedEnv models the documented contract: a session that wants the
// owner-scoped cache to resolve pins WST_ACTOR (and here WST_SERVER)
// consistently across register and complete.
func pinnedEnv(actor, server string) func(string) string {
	return func(k string) string {
		switch k {
		case "WST_ACTOR":
			return actor
		case "WST_SERVER":
			return server
		}
		return ""
	}
}

func cleanupCaches(t *testing.T) {
	t.Cleanup(func() {
		_ = os.Remove(sessionActorCachePath(t))
		_ = os.Remove(wiCachePath(t))
	})
}

// TestRegisterThenCompleteResolvesCachedID exercises the symmetric
// handshake end to end: a successful register caches the real
// work-instance id scoped to the session actor, and complete resolves
// it from that cache with no --id threaded through.
func TestRegisterThenCompleteResolvesCachedID(t *testing.T) {
	ts := startAPIServer(t)
	cleanupCaches(t)
	env := pinnedEnv("wst-fixed", ts.URL)

	var rout, rerr bytes.Buffer
	if code := runRegister([]string{"--slug", "demo-root-m1-t2"}, env, &rout, &rerr); code != 0 {
		t.Fatalf("register exit = %d; stderr=%q", code, rerr.String())
	}
	want := widPattern.FindStringSubmatch(rout.String())
	if want == nil {
		t.Fatalf("register echoed no work_instance_id: %q", rout.String())
	}

	var cout, cerr bytes.Buffer
	if code := runTerminal("complete", "completed", nil, env, &cout, &cerr); code != 0 {
		t.Fatalf("complete exit = %d; stderr=%q", code, cerr.String())
	}
	got := cout.String()
	if !strings.Contains(got, "complete: ok") || !strings.Contains(got, "state=completed") || !strings.Contains(got, "http_status=201") {
		t.Fatalf("complete receipt missing real outcome: %q", got)
	}
	if !strings.Contains(got, "work_instance_id="+want[1]) {
		t.Fatalf("complete should resolve the cached register id %q: %q", want[1], got)
	}
}

func TestAbandonCommandRecordsAbandonedState(t *testing.T) {
	ts := startAPIServer(t)
	cleanupCaches(t)
	env := pinnedEnv("wst-fixed", ts.URL)

	var rout, rerr bytes.Buffer
	if code := runRegister([]string{"--slug", "demo-root-m1-t2"}, env, &rout, &rerr); code != 0 {
		t.Fatalf("register exit = %d; stderr=%q", code, rerr.String())
	}

	var cout, cerr bytes.Buffer
	if code := runTerminal("abandon", "abandoned", nil, env, &cout, &cerr); code != 0 {
		t.Fatalf("abandon exit = %d; stderr=%q", code, cerr.String())
	}
	if !strings.Contains(cout.String(), "abandon: ok") || !strings.Contains(cout.String(), "state=abandoned") {
		t.Fatalf("abandon receipt missing real outcome: %q", cout.String())
	}
}

// TestCompleteConsumesCachedIDOnSuccess guards the symmetric stale-id
// hazard at session end: once a terminal transition succeeds, the
// owned cached id must be consumed so a later no---id complete in the
// same worktree narrates an explicit skip instead of re-transitioning
// the already-terminal work-instance.
func TestCompleteConsumesCachedIDOnSuccess(t *testing.T) {
	ts := startAPIServer(t)
	cleanupCaches(t)
	env := pinnedEnv("wst-fixed", ts.URL)

	var rout, rerr bytes.Buffer
	if code := runRegister([]string{"--slug", "demo-root-m1-t2"}, env, &rout, &rerr); code != 0 {
		t.Fatalf("register exit = %d; stderr=%q", code, rerr.String())
	}

	var c1out, c1err bytes.Buffer
	if code := runTerminal("complete", "completed", nil, env, &c1out, &c1err); code != 0 {
		t.Fatalf("first complete exit = %d; stderr=%q", code, c1err.String())
	}
	if !strings.Contains(c1out.String(), "complete: ok") {
		t.Fatalf("first complete should succeed: %q", c1out.String())
	}
	if _, err := os.Stat(wiCachePath(t)); !os.IsNotExist(err) {
		t.Fatalf("successful terminal transition must consume the owned cached id, stat err=%v", err)
	}

	var c2out, c2err bytes.Buffer
	if code := runTerminal("complete", "completed", nil, env, &c2out, &c2err); code != 0 {
		t.Fatalf("second complete exit = %d; stderr=%q", code, c2err.String())
	}
	if !strings.Contains(c2err.String(), "no work-instance id") {
		t.Fatalf("second complete must skip, not reuse the terminal id; stderr=%q stdout=%q", c2err.String(), c2out.String())
	}
}

// TestCompleteSkipsForeignOwnedCache is the load-bearing guard the
// four review rounds converged on: a cached receipt owned by a
// different session's actor must never be acted on — it is narrated
// as a skip and left intact for its real owner.
func TestCompleteSkipsForeignOwnedCache(t *testing.T) {
	ts := startAPIServer(t)
	cleanupCaches(t)

	var rout, rerr bytes.Buffer
	if code := runRegister([]string{"--slug", "demo-root-m1-t2"}, pinnedEnv("wst-A", ts.URL), &rout, &rerr); code != 0 {
		t.Fatalf("register exit = %d; stderr=%q", code, rerr.String())
	}

	var cout, cerr bytes.Buffer
	code := runTerminal("complete", "completed", nil, pinnedEnv("wst-B", ts.URL), &cout, &cerr)
	if code != 0 {
		t.Fatalf("foreign-owned cache must still exit success; got %d", code)
	}
	if !strings.Contains(cerr.String(), "owned by a different session") {
		t.Fatalf("foreign-owned cache should narrate an owner-mismatch skip: %q", cerr.String())
	}
	if cout.String() != "" {
		t.Fatalf("foreign-owned cache must not record a terminal event: %q", cout.String())
	}
	if _, err := os.Stat(wiCachePath(t)); err != nil {
		t.Fatalf("a foreign-owned cache entry must be left intact for its owner: %v", err)
	}
}

func TestCompleteCommandNoIDSkips(t *testing.T) {
	// No cached id, no flag, no env: must narrate an explicit skip
	// and still exit success — never blocks the session.
	_ = os.Remove(wiCachePath(t))

	var out, errb bytes.Buffer
	code := runTerminal("complete", "completed", nil, noEnv, &out, &errb)
	if code != 0 {
		t.Fatalf("no-id must exit success; got %d", code)
	}
	if !strings.Contains(errb.String(), "no work-instance id") || !strings.Contains(errb.String(), "session proceeds") {
		t.Fatalf("no-id should narrate an explicit, actionable skip: %q", errb.String())
	}
	if out.String() != "" {
		t.Fatalf("no-id should emit no success receipt: %q", out.String())
	}
}

// TestFailedRegisterInvalidatesPriorCachedID guards the false-signal
// hazard: a prior session cached an id, then a later session's
// register fails. The stale id must not survive for complete to mark
// an unrelated earlier work-instance terminal.
func TestFailedRegisterInvalidatesPriorCachedID(t *testing.T) {
	ts := startAPIServer(t)
	cleanupCaches(t)
	env := pinnedEnv("wst-fixed", ts.URL)

	var rout, rerr bytes.Buffer
	if code := runRegister([]string{"--slug", "demo-root-m1-t2"}, env, &rout, &rerr); code != 0 {
		t.Fatalf("first register exit = %d; stderr=%q", code, rerr.String())
	}
	if _, err := os.Stat(wiCachePath(t)); err != nil {
		t.Fatalf("successful register should have cached an id: %v", err)
	}

	// A later session's register fails (server unreachable, via the
	// flag which beats the pinned WST_SERVER).
	var fout, ferr bytes.Buffer
	runRegister([]string{"--slug", "demo-root-m1-t2", "--server", "http://127.0.0.1:0"}, env, &fout, &ferr)
	if _, err := os.Stat(wiCachePath(t)); !os.IsNotExist(err) {
		t.Fatalf("failed register must invalidate the prior cached id, stat err=%v", err)
	}

	// complete must now narrate an explicit skip, not act on the
	// stale id from the first session.
	var cout, cerr bytes.Buffer
	if code := runTerminal("complete", "completed", nil, env, &cout, &cerr); code != 0 {
		t.Fatalf("complete exit = %d; stderr=%q", code, cerr.String())
	}
	if !strings.Contains(cerr.String(), "no work-instance id") {
		t.Fatalf("complete should skip after invalidated cache, got stderr=%q stdout=%q", cerr.String(), cout.String())
	}
}

func TestCompleteCommandServerDownProceeds(t *testing.T) {
	cleanupCaches(t)
	var out, errb bytes.Buffer
	code := runTerminal(
		"complete", "completed",
		[]string{"--id", "wi-7", "--server", "http://127.0.0.1:0"},
		noEnv, &out, &errb,
	)
	if code != 0 {
		t.Fatalf("server-down must exit success; got %d", code)
	}
	if !strings.Contains(errb.String(), "attempt failed") || !strings.Contains(errb.String(), "session proceeds") {
		t.Fatalf("server-down should narrate an explicit, actionable failure: %q", errb.String())
	}
}

func TestCompleteCommandIDFromEnv(t *testing.T) {
	ts := startAPIServer(t)
	cleanupCaches(t)

	var rout, rerr bytes.Buffer
	if code := runRegister(
		[]string{"--slug", "demo-root-m1-t2", "--actor", "wst-fixed", "--server", ts.URL},
		noEnv, &rout, &rerr,
	); code != 0 {
		t.Fatalf("register exit = %d; stderr=%q", code, rerr.String())
	}
	wid := widPattern.FindStringSubmatch(rout.String())[1]

	env := func(k string) string {
		switch k {
		case "WST_WI_ID":
			return wid
		case "WST_SERVER":
			return ts.URL
		}
		return ""
	}
	var out, errb bytes.Buffer
	if code := runTerminal("complete", "completed", nil, env, &out, &errb); code != 0 {
		t.Fatalf("env-driven complete exit = %d; stderr=%q", code, errb.String())
	}
	if !strings.Contains(out.String(), "work_instance_id="+wid) {
		t.Fatalf("WST_WI_ID not honored (explicit id is an operator override): %q", out.String())
	}
}
