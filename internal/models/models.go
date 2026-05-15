// Package models defines the shared types used by the
// workstream-tracker server, API, and DB layers. See
// design/v0.1-design.md Section 5 for the data model.
package models

import "time"

// EventType enumerates the recognized event types in the
// append-only event log.
type EventType string

const (
	EventRegister        EventType = "register"
	EventHeartbeat       EventType = "heartbeat"
	EventStateTransition EventType = "state_transition"
)

// Event is one row in the event log.
type Event struct {
	ID             string
	WorkInstanceID string
	Type           EventType
	Payload        []byte // JSON, type-specific
	Metadata       []byte // JSON, optional free-form
	ReceivedAt     time.Time
}

// WorkInstanceState enumerates the runtime states of a
// work-instance.
//
// Distinct from plan-doc Status, which lives in plan-file
// frontmatter and is owned by the spec (see
// spec/planning/shared.md "Plan-doc Status"). Work-instance
// state describes "is an agent currently doing work against this
// node"; plan-doc Status describes the durable lifecycle of the
// planning artifact. The two evolve independently.
type WorkInstanceState string

const (
	StateActive    WorkInstanceState = "active"
	StateCompleted WorkInstanceState = "completed"
	StateAbandoned WorkInstanceState = "abandoned"
)

// WorkInstance is the derived current-state record for one
// work-instance. Built by folding the event log.
type WorkInstance struct {
	ID            string
	Slug          string
	Actor         string
	State         WorkInstanceState
	CreatedAt     time.Time
	LastUpdatedAt time.Time
	TerminalAt    *time.Time // nil until terminal state reached
}
