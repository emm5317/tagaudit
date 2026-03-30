package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/emm5317/tagaudit"
)

// JSONFinding is the JSON-serializable representation of a finding.
type JSONFinding struct {
	File      string `json:"file"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	RuleID    string `json:"rule_id"`
	Severity  string `json:"severity"`
	Message   string `json:"message"`
	FieldName string `json:"field_name,omitempty"`
	TagKey    string `json:"tag_key,omitempty"`
	HasFix    bool   `json:"has_fix"`
}

func outputText(w io.Writer, findings []tagaudit.Finding) error {
	for _, f := range findings {
		fmt.Fprintln(w, f.String())
	}
	return nil
}

func outputJSON(w io.Writer, findings []tagaudit.Finding) error {
	out := make([]JSONFinding, len(findings))
	for i, f := range findings {
		out[i] = JSONFinding{
			File:      f.Pos.Filename,
			Line:      f.Pos.Line,
			Column:    f.Pos.Column,
			RuleID:    f.RuleID,
			Severity:  f.Severity.String(),
			Message:   f.Message,
			FieldName: f.FieldName,
			TagKey:    f.TagKey,
			HasFix:    f.Fix != nil,
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
