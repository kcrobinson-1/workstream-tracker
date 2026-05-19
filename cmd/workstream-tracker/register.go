package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	// The early skip paths below deliberately do NOT clear the wi
	// cache: owner-scoping already makes a skipped new session safe
	// (a different actor's complete will not match), and clearing on
	// every skip would destroy a still-valid id when a tool
	// re-invokes register without a slug mid-session. The same-actor
	// skip case is an accepted residual (see the wi-cache contract).
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

	// Invalidate any cached work-instance id before attempting: if
	// this registration fails, a stale id must not survive for
	// `complete` to act on. The owner-scoped re-write on success
	// below is the load-bearing guard (see the wi-cache contract);
	// this clear is defense-in-depth for the same-actor case.
	clearWICache()

	ctx, cancel := context.WithTimeout(context.Background(), registerTimeout)
	defer cancel()

	res, err := registerclient.Register(ctx, server, slug, actor)
	if err != nil {
		fmt.Fprintf(stderr,
			"register: attempt failed (%v); session proceeds — this session will not appear in the tree. "+
				"Run `go run github.com/kcrobinson-1/workstream-tracker/cmd/workstream-tracker register --slug %s` by hand against a running server to register it.\n",
			err, slug)
		return 0
	}

	// Cache the real work-instance id, scoped to the actor that owns
	// it, so the symmetric `complete` resolves it without the session
	// threading the id through the prompt — and only when the same
	// session owns it. A cache-write failure is non-fatal: the
	// session still registered, and complete falls back to --id.
	writeWICache(actor, res.WorkInstanceID)

	fmt.Fprintf(stdout,
		"register: ok work_instance_id=%s slug=%s actor=%s http_status=%d server=%s\n",
		res.WorkInstanceID, res.Slug, actor, res.HTTPStatus, server)
	return 0
}

// wstCachePath returns the per-worktree temp-file path for a cached
// session value of the given kind ("actor", "wi"). Keying by the
// resolved worktree root (not cwd) namespaces parallel worktrees so
// their sessions never collide, while letting a register from the
// repo root and a complete from a subdirectory of the same checkout
// resolve the same cached receipt.
func wstCachePath(kind string) (string, error) {
	root, err := worktreeRoot()
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(root))
	return filepath.Join(os.TempDir(), "wst-"+kind+"-"+hex.EncodeToString(sum[:8])+".id"), nil
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

// The wi-cache contract — the single invariant every register and
// complete path must preserve. Each of the four review rounds on
// this change was one face of this contract being implicit:
//
//	The wi cache holds "<actor>\n<work-instance-id>" for a worktree
//	iff a live, non-terminal work-instance, owned by that actor, is
//	the current session's instance in that worktree. It is keyed by
//	the resolved worktree root (not cwd) so every directory in one
//	checkout shares it and parallel worktrees stay distinct; it is
//	scoped by actor so a stale or foreign id is never acted on.
//
// Enforced by exactly three operations, used everywhere instead of
// touching the cache file directly:
//   - register clears it before attempting and writes (actor, id)
//     only on success.
//   - complete/abandon reads it, USES the id only when the cached
//     actor equals the actor this invocation resolves, and removes
//     it on a successful terminal transition of an owned entry.
//
// Accepted residuals (the irreducible part, deliberately not
// engineered away; tracked by the deterministic-interactive-
// registration backlog entry):
//   - Two sessions sharing one resolved actor in one worktree
//     (generated actors within sessionActorIdleWindow, or an
//     identically pinned WST_ACTOR across true-parallel sessions)
//     still collapse onto one slot — inherited from the actor
//     cache, no worse.
//   - WST_ACTOR set inconsistently between register and complete
//     makes complete safely SKIP (it narrates and suggests --id),
//     never misattribute.
//   - An explicit --id / WST_WI_ID is honored without the owner
//     check: naming the instance is a deliberate operator override.

// resolveActor returns the actor identity for this invocation,
// resolved identically by register and complete so the wi-cache
// owner check compares like with like: an explicit --actor, else
// WST_ACTOR, else the generated per-session id. A session that pins
// an explicit actor must pass it to both halves (or pin WST_ACTOR
// for the whole session); otherwise complete cannot re-derive it
// and safely skips.
func resolveActor(actorFlag string, getenv func(string) string) (string, error) {
	if a := firstNonEmpty(actorFlag, getenv("WST_ACTOR")); a != "" {
		return a, nil
	}
	return sessionActor()
}

// writeWICache records the owning actor and work-instance id. Best
// effort by contract: a non-writable temp dir must not fail the
// session — complete then falls back to --id.
func writeWICache(actor, id string) {
	path, err := wstCachePath("wi")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, []byte(actor+"\n"+id), 0o600)
}

// readWICache returns the cached owning actor and work-instance id.
// ok is false when the cache is absent, unreadable, or not in the
// owner-scoped two-line form — all of which mean "no usable cached
// id," never "use it unverified."
func readWICache() (actor, id string, ok bool) {
	path, err := wstCachePath("wi")
	if err != nil {
		return "", "", false
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", "", false
	}
	parts := strings.SplitN(strings.TrimRight(string(raw), "\n"), "\n", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	actor = strings.TrimSpace(parts[0])
	id = strings.TrimSpace(parts[1])
	if actor == "" || id == "" {
		return "", "", false
	}
	return actor, id, true
}

// clearWICache removes the cached entry. Best effort by contract.
func clearWICache() {
	if path, err := wstCachePath("wi"); err == nil {
		_ = os.Remove(path)
	}
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
