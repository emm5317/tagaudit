package tagaudit

import (
	"fmt"
	"go/token"
)

// Severity classifies the importance of a finding.
// Higher values are more severe. The zero value (SeverityInfo) includes all
// findings, making a default Config{} permissive by default.
type Severity int

const (
	SeverityInfo    Severity = iota // Nice to fix: naming conventions, completeness
	SeverityWarning                 // Should fix: duplicates, shadowing
	SeverityError                   // Must fix: broken tags, syntax errors
)

func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "info"
	case SeverityWarning:
		return "warning"
	case SeverityError:
		return "error"
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

	// TagStart and TagEnd are byte offsets in the source file for the
	// tag literal (including backticks). When set (both > 0), the CLI
	// fixer uses these for precise replacement instead of line-based
	// heuristics. The go/analysis path uses AST positions directly.
	TagStart int
	TagEnd   int
}
