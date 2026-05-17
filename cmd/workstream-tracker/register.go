package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/kcrobinson-1/workstream-tracker/internal/registerclient"
)

// registerTimeout bounds the single best-effort registration
// attempt. Registration must never slow a session, so this is
// short and there is no retry.
const registerTimeout = 3 * time.Second

// defaultServer is the registration target when neither --server
// nor WST_SERVER is set. It mirrors the server's own PORT=8080
// default in runServer.
const defaultServer = "http://localhost:8080"

// runRegister performs one best-effort work-instance registration
// and always returns 0: registration is a consequence of starting
// work, never a gate on it. Every failure path logs an explicit,
// actionable line to stderr and still exits success so the calling
// session proceeds. On success it prints the real observed receipt
// (work-instance id, slug, HTTP status) to stdout so the agent can
// echo fact, not a prose success claim.
func runRegister(args []string, getenv func(string) string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("register", flag.ContinueOnError)
	fs.SetOutput(stderr)
	slugFlag := fs.String("slug", "", "canonical plan-doc slug to register a work-instance for (or WST_SLUG)")
	actorFlag := fs.String("actor", "", "actor label (or WST_ACTOR; defaults to a generated per-session id)")
	serverFlag := fs.String("server", "", "server base URL (or WST_SERVER; default "+defaultServer+")")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "register: %v; skipping registration, session proceeds\n", err)
		return 0
	}

	slug := firstNonEmpty(*slugFlag, getenv("WST_SLUG"))
	if slug == "" {
		fmt.Fprintln(stderr, "register: no slug supplied (--slug or WST_SLUG); skipping registration, session proceeds")
		return 0
	}

	server := firstNonEmpty(*serverFlag, getenv("WST_SERVER"), defaultServer)

	actor := firstNonEmpty(*actorFlag, getenv("WST_ACTOR"))
	if actor == "" {
		a, err := sessionActor()
		if err != nil {
			fmt.Fprintf(stderr, "register: could not derive a per-session actor id (%v); skipping registration, session proceeds\n", err)
			return 0
		}
		actor = a
	}

	ctx, cancel := context.WithTimeout(context.Background(), registerTimeout)
	defer cancel()

	res, err := registerclient.Register(ctx, server, slug, actor)
	if err != nil {
		fmt.Fprintf(stderr,
			"register: attempt failed (%v); session proceeds — this session will not appear in the tree. "+
				"Run `workstream-tracker register --slug %s` by hand against a running server to register it.\n",
			err, slug)
		return 0
	}

	fmt.Fprintf(stdout,
		"register: ok work_instance_id=%s slug=%s actor=%s http_status=%d server=%s\n",
		res.WorkInstanceID, res.Slug, actor, res.HTTPStatus, server)
	return 0
}

// sessionActor returns a generated per-session actor id, stable
// across repeat invocations within the same worktree (so a
// restart/resume collapses via the server's (slug, actor, active)
// idempotency) and distinct across parallel worktrees (so parallel
// agents never collapse onto one marker). It is never the git
// user. The id is cached in a temp file keyed by the absolute
// working directory; the first invocation generates it, later ones
// read it back.
func sessionActor() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve working directory: %w", err)
	}
	sum := sha256.Sum256([]byte(wd))
	cachePath := filepath.Join(os.TempDir(), "wst-actor-"+hex.EncodeToString(sum[:8])+".id")

	if existing, err := os.ReadFile(cachePath); err == nil {
		if id := strings.TrimSpace(string(existing)); id != "" {
			return id, nil
		}
	}

	id := "wst-" + uuid.NewString()
	if err := os.WriteFile(cachePath, []byte(id), 0o600); err != nil {
		// A non-writable temp dir must not block the session; fall
		// back to the freshly generated id without caching it.
		return id, nil
	}
	return id, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
