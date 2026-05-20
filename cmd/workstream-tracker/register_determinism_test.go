package main

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

// proofSlug is the construction-known, position-segmented slug
// driven through the bare --slug producer (the stand-in for m2's
// future construction-time producer). The mN/tN suffix is
// load-bearing: a slug carrying its own position segments
// surviving verbatim distinguishes the exact-slug branch (honors
// the caller's slug) from any derivation or server-side
// descendant-generation path (which would mint its own position
// number).
const proofSlug = "proof-root-m1-t2"

type receipt struct {
	workInstanceID string
	slug           string
	httpStatus     string
}

var (
	receiptSlugPattern   = regexp.MustCompile(`slug=(\S+)`)
	receiptStatusPattern = regexp.MustCompile(`http_status=(\S+)`)
)

func parseReceipt(t *testing.T, stdout string) receipt {
	t.Helper()
	wid := widPattern.FindStringSubmatch(stdout)
	if wid == nil {
		t.Fatalf("receipt missing work_instance_id (real-receipt echo is the sole observability surface): %q", stdout)
	}
	slug := receiptSlugPattern.FindStringSubmatch(stdout)
	if slug == nil {
		t.Fatalf("receipt missing slug: %q", stdout)
	}
	status := receiptStatusPattern.FindStringSubmatch(stdout)
	if status == nil {
		t.Fatalf("receipt missing http_status: %q", stdout)
	}
	return receipt{workInstanceID: wid[1], slug: slug[1], httpStatus: status[1]}
}

// registerWithSlug drives the deterministic, handshake-free
// registration path the way m2's construction-time slug producer
// will: the bare --slug argument, an explicit --actor (no
// generated id), an empty environment (noEnv: no WST_* leak), and
// nothing else. runRegister takes no prompt argument and has no
// resolve/confirm step, so this invocation shape itself
// demonstrates the structural absence of a handshake (C5).
func registerWithSlug(t *testing.T, serverURL, slug, actor string) (string, string, int) {
	t.Helper()
	var out, errb bytes.Buffer
	code := runRegister(
		[]string{"--slug", slug, "--actor", actor, "--server", serverURL},
		noEnv, &out, &errb,
	)
	return out.String(), errb.String(), code
}

// TestRegisterDeterminismProof exercises the merged exact-slug
// create-or-attach path through the CLI with a construction-known
// slug and no narration handshake, asserting each named
// determinism property of the deterministic-path contract
// (docs/agents/local/session-registration.md "The deterministic
// path", spec/planning/shared.md "Deterministic assertion when
// the slug is construction-known") with a discriminator strong
// enough that the assertion fails if the property does not hold.
//
// Falsifiability discipline: a bare http_status=201 is the trap
// for C3 (multiple causes satisfy it) and for C4 (a real
// plan-tree slug would also 201). Each property carries a
// positive/negative discriminator a single 201 cannot satisfy.
func TestRegisterDeterminismProof(t *testing.T) {
	isolateCaches(t)
	ts := startAPIServer(t)

	t.Run("C1 identity-by-construction: verbatim round-trip", func(t *testing.T) {
		out, errb, code := registerWithSlug(t, ts.URL, proofSlug, "wst-c1-actor")
		if code != 0 {
			t.Fatalf("runRegister exit = %d, want 0; stderr=%q", code, errb)
		}
		r := parseReceipt(t, out)
		if r.httpStatus != "201" {
			t.Fatalf("C1: want http_status=201, got %q (stdout=%q)", r.httpStatus, out)
		}
		// Byte-identity is the discriminator: a derivation would
		// mint its own mN/tN; the exact-slug branch honors the
		// caller's slug verbatim.
		if r.slug != proofSlug {
			t.Fatalf("C1: receipt slug must be byte-identical to the supplied slug; supplied=%q got=%q", proofSlug, r.slug)
		}
	})

	t.Run("C2a idempotent attach on same (slug, actor) while active", func(t *testing.T) {
		const actor = "wst-c2a-actor"
		out1, _, code1 := registerWithSlug(t, ts.URL, proofSlug, actor)
		if code1 != 0 {
			t.Fatalf("first call exit = %d, want 0", code1)
		}
		r1 := parseReceipt(t, out1)
		out2, _, code2 := registerWithSlug(t, ts.URL, proofSlug, actor)
		if code2 != 0 {
			t.Fatalf("repeat call exit = %d, want 0", code2)
		}
		r2 := parseReceipt(t, out2)
		if r1.workInstanceID != r2.workInstanceID {
			t.Fatalf("C2(a): same-(slug, actor) repeat must collapse to same work_instance_id; first=%q repeat=%q",
				r1.workInstanceID, r2.workInstanceID)
		}
		// Both receipts must report 201 — the idempotent attach
		// shares the created status with a fresh row; this is
		// what distinguishes the exact-slug branch (C3) from the
		// bare-root-create branch, which would 409 on a repeat.
		if r1.httpStatus != "201" || r2.httpStatus != "201" {
			t.Fatalf("C2(a): both calls must echo http_status=201; got first=%q repeat=%q", r1.httpStatus, r2.httpStatus)
		}
	})

	t.Run("C2b distinct work-instance for a different actor on the same slug", func(t *testing.T) {
		const actorA = "wst-c2b-actor-a"
		const actorB = "wst-c2b-actor-b"
		outA, _, codeA := registerWithSlug(t, ts.URL, proofSlug, actorA)
		if codeA != 0 {
			t.Fatalf("actor-A call exit = %d, want 0", codeA)
		}
		outB, _, codeB := registerWithSlug(t, ts.URL, proofSlug, actorB)
		if codeB != 0 {
			t.Fatalf("actor-B call exit = %d, want 0", codeB)
		}
		rA := parseReceipt(t, outA)
		rB := parseReceipt(t, outB)
		// (slug, actor) — not slug alone — is the attach
		// identity. If attach keyed on slug, B would collapse
		// onto A's row; demanding a distinct id is the
		// discriminator that proves the keys are paired.
		if rA.workInstanceID == rB.workInstanceID {
			t.Fatalf("C2(b): distinct actors on same slug must yield distinct work_instance_ids; got %q for both", rA.workInstanceID)
		}
		if rA.slug != proofSlug || rB.slug != proofSlug {
			t.Fatalf("C2(b): both receipts must still round-trip the supplied slug verbatim; got A=%q B=%q", rA.slug, rB.slug)
		}
	})

	t.Run("C3 no NL-resolution path engaged (C1 conjunction C2(a))", func(t *testing.T) {
		// C3 is fingerprinted by C1 ∧ C2(a) together: a
		// position-segmented slug surviving verbatim rules out
		// interpretation/derivation, AND a 201 + same id on a
		// repeat positively identifies the exact-slug attach
		// branch (the bare-root-create branch 409s on a repeat;
		// the descendant-generation branch never carries a
		// caller exact slug). A status-code-only assertion is
		// explicitly insufficient (the falsifiability trap).
		const actor = "wst-c3-actor"
		out1, _, code1 := registerWithSlug(t, ts.URL, proofSlug, actor)
		if code1 != 0 {
			t.Fatalf("first call exit = %d, want 0", code1)
		}
		r1 := parseReceipt(t, out1)
		out2, _, code2 := registerWithSlug(t, ts.URL, proofSlug, actor)
		if code2 != 0 {
			t.Fatalf("repeat call exit = %d, want 0", code2)
		}
		r2 := parseReceipt(t, out2)
		if r1.slug != proofSlug || r2.slug != proofSlug {
			t.Fatalf("C3: verbatim leg fails — both receipts must echo the supplied slug byte-identically; got first=%q repeat=%q", r1.slug, r2.slug)
		}
		if r1.httpStatus != "201" || r2.httpStatus != "201" {
			t.Fatalf("C3: attach leg fails — both calls must echo http_status=201 (the bare-root-create branch would 409 on the repeat); got first=%q repeat=%q", r1.httpStatus, r2.httpStatus)
		}
		if r1.workInstanceID != r2.workInstanceID {
			t.Fatalf("C3: attach leg fails — same-(slug, actor) repeat must collapse to same work_instance_id; got first=%q repeat=%q", r1.workInstanceID, r2.workInstanceID)
		}
	})

	t.Run("C4 server repo-blind: paired grammar-gate fingerprint", func(t *testing.T) {
		// (a) A well-formed slug that names no real plan-tree
		// doc still registers 201 + verbatim — existence is
		// never checked, so the slug's repo-resolvability is
		// irrelevant. The synthetic root token plus the m9/t9
		// positions could not be produced by any real plan tree
		// in this repo.
		const syntheticSlug = "synthetic-repo-blind-proof-m9-t9"
		outA, errbA, codeA := registerWithSlug(t, ts.URL, syntheticSlug, "wst-c4-synthetic")
		if codeA != 0 {
			t.Fatalf("C4(a) exit = %d, want 0; stderr=%q", codeA, errbA)
		}
		rA := parseReceipt(t, outA)
		if rA.httpStatus != "201" {
			t.Fatalf("C4(a): synthetic well-formed slug must register 201; got %q", rA.httpStatus)
		}
		if rA.slug != syntheticSlug {
			t.Fatalf("C4(a): synthetic slug must round-trip verbatim; supplied=%q got=%q", syntheticSlug, rA.slug)
		}

		// (b) A grammatically malformed slug (underscore in the
		// root token) is rejected 400 by the grammar gate —
		// IsWellFormed is the only acceptance check, and it is
		// repo-blind. The CLI is best-effort, so a server
		// rejection still exits 0 and surfaces the real HTTP 400
		// in the failure narration on stderr; no success receipt
		// is printed.
		const malformedSlug = "bad_underscored-root-m1"
		outB, errbB, codeB := registerWithSlug(t, ts.URL, malformedSlug, "wst-c4-malformed")
		if codeB != 0 {
			t.Fatalf("C4(b) exit = %d, want 0 (best-effort: failure narrated, session proceeds); stderr=%q", codeB, errbB)
		}
		if outB != "" {
			t.Fatalf("C4(b): malformed slug must produce no success receipt; got stdout=%q", outB)
		}
		if !strings.Contains(errbB, "HTTP 400") {
			t.Fatalf("C4(b): malformed slug must be rejected with HTTP 400 on stderr (the grammar gate is the only acceptance check); got %q", errbB)
		}
	})

	t.Run("C5 handshake-free (structural)", func(t *testing.T) {
		// "Handshake-free" is shown structurally: runRegister
		// takes no prompt argument and has no resolve/confirm
		// step, so the invocation shape itself — bare
		// --slug/--actor/--server with noEnv (no WST_* leak) —
		// cannot engage a handshake. A real 201 receipt from
		// exactly that invocation closes the proof: the
		// deterministic path lands without a handshake. If a
		// future change wired a resolve/confirm interaction into
		// runRegister, this assertion would either hang reading
		// from stdin or fail at the receipt parse — a loud
		// regression signal, not a silent drift.
		out, errb, code := registerWithSlug(t, ts.URL, proofSlug, "wst-c5-actor")
		if code != 0 {
			t.Fatalf("C5 exit = %d, want 0; stderr=%q", code, errb)
		}
		r := parseReceipt(t, out)
		if r.httpStatus != "201" {
			t.Fatalf("C5: handshake-free invocation must still register 201; got %q", r.httpStatus)
		}
		if r.workInstanceID == "" {
			t.Fatalf("C5: receipt must echo a real work_instance_id from a real registration; got empty (stdout=%q)", out)
		}
	})
}
