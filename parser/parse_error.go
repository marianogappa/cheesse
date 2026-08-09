package parser

import "github.com/marianogappa/cheesse/core"

// ParseError is a structured parse-failure error that carries machine-readable
// detail alongside the human-readable message. The api layer inspects these
// fields to build OutputParseError without reparsing the error string.
type ParseError struct {
	Reason       string        // "unknown-token" | "illegal-move" | "ambiguous" | "no-moves-found"
	Token        string        // the failing token text (empty for unknown-token / no-moves-found)
	MoveNumber   int           // 1-based full-move number at the failure point (0 = before first move)
	LegalActions []core.Action // legal actions at the position (may be nil)
	Msg          string        // user-facing prose
}

func (e *ParseError) Error() string { return e.Msg }
