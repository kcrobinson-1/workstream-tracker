package main

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/kcrobinson-1/workstream-tracker/internal/registerclient"
)

// runTerminal performs one best-effort terminal-state transition for
// the current session's work-instance and always returns 0: marking
// a session complete is a consequence of finishing work, never a
// gate on it. cmd is the subcommand label ("complete" / "abandon");
// state is the API state it records ("completed" / "abandoned").
//
// Symmetric to runRegister: every failure path logs an explicit,
// actionable line to stderr and still exits success so the calling
// session is never blocked. On success it prints the real observed
// receipt (event id, work-instance id, HTTP status) to stdout so the
// agent echoes fact, not a prose success claim.
//
// Work-instance id resolution honors the wi-cache contract (see
// register.go): an explicit --id / WST_WI_ID is an operator override
// used verbatim; otherwise the cached receipt is used only when its
// owning actor equals the actor this invocation resolves, so a stale
// or foreign id is narrated as a skip, never acted on.
func runTerminal(cmd, state string, args []string, getenv func(string) string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
	fs.SetOutput(stderr)
	idFlag := fs.String("id", "", "work-instance id to transition (or WST_WI_ID; defaults to the owned id cached by a prior register)")
	actorFlag := fs.String("actor", "", "actor label (or WST_ACTOR; must match the register that cached the receipt)")
	serverFlag := fs.String("server", "", "server base URL (or WST_SERVER; default "+defaultServer+")")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "%s: %v; skipping, session proceeds\n", cmd, err)
		return 0
	}

	manualCmd := "go run github.com/kcrobinson-1/workstream-tracker/cmd/workstream-tracker " + cmd

	explicitID := firstNonEmpty(*idFlag, getenv("WST_WI_ID"))
	// Explicit actor only — never sessionActor(). Re-deriving the
	// generated per-session id here would rotate after the idle
	// window and make an ordinary long session's complete skip its
	// own still-cached receipt; and within the window two sessions
	// derive the same generated id anyway, so it is no reliable owner
	// signal. The reliable owner signal is a pinned WST_ACTOR/--actor;
	// with none, the register-side lifecycle (clear-before-attempt +
	// consume-on-success) is the guard (see the wi-cache contract).
	explicitActor := firstNonEmpty(*actorFlag, getenv("WST_ACTOR"))

	var id string
	if explicitID != "" {
		id = explicitID
	} else {
		cachedActor, cachedID, ok := readWICache()
		if !ok {
			fmt.Fprintf(stderr,
				"%s: no work-instance id (--id or WST_WI_ID, and no cached receipt from a register in this worktree); "+
					"skipping, session proceeds — the tree will keep showing this session active. "+
					"Run `%s --id <work_instance_id>` by hand against a running server to record it.\n",
				cmd, manualCmd)
			return 0
		}
		if explicitActor != "" && cachedActor != explicitActor {
			fmt.Fprintf(stderr,
				"%s: a cached receipt exists but is owned by a different session (cached actor does not match the pinned WST_ACTOR/--actor); "+
					"skipping, session proceeds. "+
					"Run `%s --id <work_instance_id>` by hand to act on a specific work-instance.\n",
				cmd, manualCmd)
			return 0
		}
		id = cachedID
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

	// Consume the cache only when the entry IS the id just made
	// terminal: a later no---id complete must then narrate "no
	// current session" rather than re-transition it. Keyed on the id,
	// not the actor — an explicit --id transitioning a different
	// work-instance must not wipe this session's still-live receipt.
	if _, cachedID, ok := readWICache(); ok && cachedID == id {
		clearWICache()
	}

	fmt.Fprintf(stdout,
		"%s: ok event_id=%s work_instance_id=%s state=%s http_status=%d server=%s\n",
		cmd, res.EventID, id, state, res.HTTPStatus, server)
	return 0
}
