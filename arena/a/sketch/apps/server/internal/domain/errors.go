// Package domain is the pure rules of the app: the Order state machine, Balance, effective
// statuses, checkout and profile validation. No database, no HTTP, no clock reads: time
// arrives as an argument. Every rule the spec states lives here exactly once.
package domain

// Kind maps a rule failure to an HTTP status in exactly one place (server.fail).
type Kind uint8

const (
	Invalid   Kind = iota + 1 // 400: input breaks a rule
	Forbidden                 // 403: actor may never do this
	NotFound                  // 404: missing, or hidden from this actor (other Member's Order)
	Conflict                  // 409: allowed in general, not in the current state
)

// Error is the only error type handlers translate to a response. Msg is Thai, shown verbatim.
type Error struct {
	Kind Kind
	Msg  string
}

func (e *Error) Error() string { return e.Msg }

// Actor is who is acting. Built from the session-loaded Member on every request,
// so a role change applies on the next request.
type Actor struct {
	MemberID int64
	Admin    bool
}
