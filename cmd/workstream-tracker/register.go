package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
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
	nameFlag := fs.String("name", "", "human session name to report in the roster (or WST_NAME; optional)")
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

	actor, err := resolveActor(*actorFlag, getenv)
	if err != nil {
		fmt.Fprintf(stderr, "register: could not derive a per-session actor id (%v); skipping registration, session proceeds\n", err)
		return 0
	}

	// The sole conventionally-read metadata key is the optional
	// human `name` (the roster's display label). The CLI owns this
	// convention; the client stays schema-agnostic. When no name
	// is supplied no metadata is sent at all (not an empty name),
	// so a v0.2-minimum session still lists and the roster falls
	// back to the slug. json.Marshal of a one-string-key map
	// cannot fail, so a name never blocks the best-effort attempt.
	var metadata json.RawMessage
	if name := firstNonEmpty(*nameFlag, getenv("WST_NAME")); name != "" {
		metadata, _ = json.Marshal(map[string]string{"name": name})
	}

	ctx, cancel := context.WithTimeout(context.Background(), registerTimeout)
	defer cancel()

	res, err := registerclient.Register(ctx, server, slug, actor, metadata)
	if err != nil {
		fmt.Fprintf(stderr,
			"register: attempt failed (%v); session proceeds — this session will not appear in the tree. "+
				"Run `go run github.com/kcrobinson-1/workstream-tracker/cmd/workstream-tracker register --slug %s` by hand against a running server to register it.\n",
			err, slug)
		return 0
	}

	fmt.Fprintf(stdout,
		"register: ok work_instance_id=%s slug=%s actor=%s http_status=%d server=%s\n",
		res.WorkInstanceID, res.Slug, actor, res.HTTPStatus, server)
	return 0
}

// cacheBaseDir holds the per-worktree actor-cache file. It is
// os.TempDir() in production; tests repoint it at a temp dir so
// `go test` (which runs inside this repo, i.e. a real worktree root)
// never reads, writes, or deletes a live session's actual cache.
var cacheBaseDir = os.TempDir()

// wstCachePath returns the per-worktree temp-file path for the
// generated per-session actor (the only cached value). Keying by the
// resolved worktree root (not cwd) namespaces parallel worktrees so
// their sessions never collide on one actor.
func wstCachePath(kind string) (string, error) {
	root, err := worktreeRoot()
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(root))
	return filepath.Join(cacheBaseDir, "wst-"+kind+"-"+hex.EncodeToString(sum[:8])+".id"), nil
}

// worktreeRoot resolves the git worktree root so a session's cached
// values are shared across every directory within the same checkout.
// It falls back to the working directory when git is unavailable or
// the path is not a work tree — the best-effort contract must never
// block on cache keying.
func worktreeRoot() (string, error) {
	if out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output(); err == nil {
		if root := strings.TrimSpace(string(out)); root != "" {
			return root, nil
		}
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve worktree root: %w", err)
	}
	return wd, nil
}

// resolveActor returns the actor identity for a register: an
// explicit --actor, else WST_ACTOR, else a generated per-session id
// (server-side idempotency dedupes active registrations by
// (slug, actor)).
func resolveActor(actorFlag string, getenv func(string) string) (string, error) {
	if a := firstNonEmpty(actorFlag, getenv("WST_ACTOR")); a != "" {
		return a, nil
	}
	return sessionActor()
}

// sessionActorIdleWindow bounds how long a cached actor id is
// reused. Repeat register calls within one working session are
// temporally clustered (the handshake fires at session start;
// resume/restart follows within the same sitting), so a sliding
// idle window collapses them onto one actor while a genuinely
// later session — separated by a gap longer than this — gets a
// fresh actor. A purely worktree-keyed cache with no expiry would
// weld every future session in the worktree to one actor forever:
// because the server dedupes active registrations by
// (slug, actor), a later session on a slug that still has a
// lingering active row would silently collapse onto the stale
// work-instance and emit no marker, defeating cross-agent
// visibility. A harness that can supply a real per-session id
// should set WST_ACTOR; this default is the best harness-agnostic
// approximation, with the residual (two sessions within the same
// idle window, same worktree, same slug collapsing) tracked by the
// deterministic-interactive-registration backlog tripwire.
const sessionActorIdleWindow = 6 * time.Hour

// sessionActor returns a generated per-session actor id: stable
// across a session's clustered repeat invocations (so a
// restart/resume collapses via the server's (slug, actor, active)
// idempotency), distinct for a later session (idle window expiry),
// and distinct across parallel worktrees (so parallel agents never
// collapse onto one marker). It is never the git user. The id is
// cached in a temp file keyed by the resolved worktree root; each
// reuse slides the window forward so an active session keeps its
// actor, and a stale cache regenerates.
func sessionActor() (string, error) {
	cachePath, err := wstCachePath("actor")
	if err != nil {
		return "", err
	}

	if info, err := os.Stat(cachePath); err == nil && time.Since(info.ModTime()) <= sessionActorIdleWindow {
		if existing, err := os.ReadFile(cachePath); err == nil {
			if id := strings.TrimSpace(string(existing)); id != "" {
				// Slide the idle window so an active session that
				// keeps calling register keeps its actor.
				now := time.Now()
				_ = os.Chtimes(cachePath, now, now)
				return id, nil
			}
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
