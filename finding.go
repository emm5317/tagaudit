package tagaudit

import (
	"fmt"
	"go/token"
)

// Severity classifies the importance of a finding.
type Severity int

const (
	SeverityError   Severity = iota // Must fix: broken tags, syntax errors
	SeverityWarning                 // Should fix: duplicates, shadowing
	SeverityInfo                    // Nice to fix: naming conventions, completeness
)

func (s Severity) String() string {
	switch s {
	case SeverityError:
		return "error"
	case SeverityWarning:
		return "warning"
	case SeverityInfo:
		return "info"
	default:
		return "unknown"
	}
}

// Finding is a single diagnostic produced by a rule.
type Finding struct {
	Pos       token.Position // file, line, column
	End       token.Position // optional end position
	RuleID    string         // e.g. "syntax", "naming", "duplicates"
	Severity  Severity
	Message   string
	FieldName string        // the struct field name
	TagKey    string        // the tag key involved, empty if N/A
	Fix       *SuggestedFix // optional fix, nil if none
}

// String returns a human-readable representation of the finding.
func (f Finding) String() string {
	return fmt.Sprintf("%s: [%s] %s: %s", f.Pos, f.Severity, f.RuleID, f.Message)
}

// SuggestedFix describes a replacement to fix the finding.
type SuggestedFix struct {
	Description string
	NewTagValue string // the corrected tag string for the field
}
