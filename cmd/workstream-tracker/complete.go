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
	actor, actorErr := resolveActor(*actorFlag, getenv)

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
		if actorErr != nil {
			fmt.Fprintf(stderr,
				"%s: a cached receipt exists but a per-session actor could not be derived to verify it owns this session (%v); "+
					"skipping, session proceeds — pass `--id <work_instance_id>` to act on it explicitly.\n",
				cmd, actorErr)
			return 0
		}
		if cachedActor != actor {
			fmt.Fprintf(stderr,
				"%s: a cached receipt exists but is owned by a different session (cached actor does not match this session's); "+
					"skipping, session proceeds — this session never registered, or WST_ACTOR differs from the register. "+
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

	// Consume the cached id on a successful terminal transition, but
	// only the entry this session owns: the work-instance is now
	// terminal, so a later no---id complete must narrate "no current
	// session" rather than re-transition it. A foreign entry (acted
	// on via an explicit --id) is left untouched.
	if cachedActor, _, ok := readWICache(); ok && actorErr == nil && cachedActor == actor {
		clearWICache()
	}

	fmt.Fprintf(stdout,
		"%s: ok event_id=%s work_instance_id=%s state=%s http_status=%d server=%s\n",
		cmd, res.EventID, id, state, res.HTTPStatus, server)
	return 0
}
