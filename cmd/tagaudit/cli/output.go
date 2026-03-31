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
	fmt.Fprintln(w, summary(findings))
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
	if err := enc.Encode(out); err != nil {
		return err
	}
	return nil
}

func summary(findings []tagaudit.Finding) string {
	if len(findings) == 0 {
		return "No findings."
	}

	var errors, warnings, infos, fixable int
	for _, f := range findings {
		switch f.Severity {
		case tagaudit.SeverityError:
			errors++
		case tagaudit.SeverityWarning:
			warnings++
		case tagaudit.SeverityInfo:
			infos++
		}
		if f.Fix != nil {
			fixable++
		}
	}

	s := fmt.Sprintf("%d finding(s):", len(findings))
	if errors > 0 {
		s += fmt.Sprintf(" %d error", errors)
		if errors > 1 {
			s += "s"
		}
	}
	if warnings > 0 {
		s += fmt.Sprintf(" %d warning", warnings)
		if warnings > 1 {
			s += "s"
		}
	}
	if infos > 0 {
		s += fmt.Sprintf(" %d info", infos)
	}
	if fixable > 0 {
		s += fmt.Sprintf(" (%d fixable)", fixable)
	}
	return s
}
