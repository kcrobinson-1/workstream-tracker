package main

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/kcrobinson-1/workstream-tracker/internal/registerclient"
)

// runTerminal performs one best-effort terminal-state transition for
// the session's work-instance and always returns 0: marking a session
// complete is a consequence of finishing work, never a gate on it.
// cmd is the subcommand label ("complete" / "abandon"); state is the
// API state it records ("completed" / "abandoned").
//
// Symmetric to runRegister: every failure path logs an explicit,
// actionable line to stderr and still exits success so the calling
// session is never blocked. On success it prints the real observed
// receipt (event id, work-instance id, HTTP status) to stdout so the
// agent echoes fact, not a prose success claim.
//
// The work-instance id is supplied explicitly via --id / WST_WI_ID —
// the id register printed in its receipt. There is no hidden
// cross-invocation state: with no id the command narrates a skip and
// exits success.
func runTerminal(cmd, state string, args []string, getenv func(string) string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
	fs.SetOutput(stderr)
	idFlag := fs.String("id", "", "work-instance id to transition (or WST_WI_ID; the id from the register receipt)")
	serverFlag := fs.String("server", "", "server base URL (or WST_SERVER; default "+defaultServer+")")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "%s: %v; skipping, session proceeds\n", cmd, err)
		return 0
	}

	manualCmd := "go run github.com/kcrobinson-1/workstream-tracker/cmd/workstream-tracker " + cmd

	id := firstNonEmpty(*idFlag, getenv("WST_WI_ID"))
	if id == "" {
		fmt.Fprintf(stderr,
			"%s: no work-instance id (--id or WST_WI_ID); skipping, session proceeds — "+
				"the tree will keep showing this session active. Pass the work_instance_id "+
				"from the register receipt: `%s --id <work_instance_id>`.\n",
			cmd, manualCmd)
		return 0
	}

	server := firstNonEmpty(*serverFlag, getenv("WST_SERVER"), defaultServer)

	ctx, cancel := context.WithTimeout(context.Background(), registerTimeout)
	defer cancel()

	res, err := registerclient.RecordState(ctx, server, id, state)
	if err != nil {
		fmt.Fprintf(stderr,
			"%s: attempt failed (%v); session proceeds — the tree will keep showing this session active. "+
				"Run `%s --id %s` by hand against a running server to record it.\n",
			cmd, err, manualCmd, id)
		return 0
	}

	fmt.Fprintf(stdout,
		"%s: ok event_id=%s work_instance_id=%s state=%s http_status=%d server=%s\n",
		cmd, res.EventID, id, state, res.HTTPStatus, server)
	return 0
}
