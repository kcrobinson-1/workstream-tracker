package main

import (
	"bytes"
	"strings"
	"testing"
)

// TestRegisterThenCompleteWithExplicitID is the end-to-end handshake:
// register prints a work_instance_id; the session passes that id to
// complete via --id. No hidden cross-invocation state.
func TestRegisterThenCompleteWithExplicitID(t *testing.T) {
	isolateCaches(t)
	ts := startAPIServer(t)

	var rout, rerr bytes.Buffer
	if code := runRegister(
		[]string{"--slug", "demo-root-m1-t2", "--actor", "wst-fixed", "--server", ts.URL},
		noEnv, &rout, &rerr,
	); code != 0 {
		t.Fatalf("register exit = %d; stderr=%q", code, rerr.String())
	}
	m := widPattern.FindStringSubmatch(rout.String())
	if m == nil {
		t.Fatalf("register echoed no work_instance_id: %q", rout.String())
	}
	wid := m[1]

	var cout, cerr bytes.Buffer
	if code := runTerminal("complete", "completed", []string{"--id", wid, "--server", ts.URL}, noEnv, &cout, &cerr); code != 0 {
		t.Fatalf("complete exit = %d; stderr=%q", code, cerr.String())
	}
	got := cout.String()
	if !strings.Contains(got, "complete: ok") || !strings.Contains(got, "state=completed") ||
		!strings.Contains(got, "http_status=201") || !strings.Contains(got, "work_instance_id="+wid) {
		t.Fatalf("complete receipt missing real outcome: %q", got)
	}
}

func TestAbandonCommandRecordsAbandonedState(t *testing.T) {
	isolateCaches(t)
	ts := startAPIServer(t)

	var rout, rerr bytes.Buffer
	if code := runRegister(
		[]string{"--slug", "demo-root-m1-t2", "--actor", "wst-fixed", "--server", ts.URL},
		noEnv, &rout, &rerr,
	); code != 0 {
		t.Fatalf("register exit = %d; stderr=%q", code, rerr.String())
	}
	wid := widPattern.FindStringSubmatch(rout.String())[1]

	var cout, cerr bytes.Buffer
	if code := runTerminal("abandon", "abandoned", []string{"--id", wid, "--server", ts.URL}, noEnv, &cout, &cerr); code != 0 {
		t.Fatalf("abandon exit = %d; stderr=%q", code, cerr.String())
	}
	if !strings.Contains(cout.String(), "abandon: ok") || !strings.Contains(cout.String(), "state=abandoned") {
		t.Fatalf("abandon receipt missing real outcome: %q", cout.String())
	}
}

func TestCompleteCommandNoIDSkips(t *testing.T) {
	// No id supplied: narrate an explicit skip and still exit success
	// — never blocks the session.
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

func TestCompleteCommandServerDownProceeds(t *testing.T) {
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
	isolateCaches(t)
	ts := startAPIServer(t)

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
		t.Fatalf("WST_WI_ID not honored: %q", out.String())
	}
}
