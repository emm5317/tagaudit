package rules

import (
	"fmt"

	"github.com/emm5317/tagaudit"
	"github.com/fatih/structtag"
)

// SyntaxRule validates that struct tags are well-formed.
type SyntaxRule struct{}

func (r *SyntaxRule) ID() string          { return "syntax" }
func (r *SyntaxRule) Description() string { return "validates struct tag syntax" }

func (r *SyntaxRule) CheckField(info tagaudit.FieldInfo, _ *tagaudit.Config) []tagaudit.Finding {
	if info.RawTag == "" {
		return nil
	}

	// If tags parsed successfully, no syntax error
	if info.Tags != nil {
		return nil
	}

	// Parse again to get the error message
	_, err := structtag.Parse(info.RawTag)
	if err == nil {
		return nil
	}

	var fieldName string
	if info.Field != nil {
		fieldName = info.Field.Name()
	}

	pos := posFromInfo(info)

	return []tagaudit.Finding{{
		Pos:       pos,
		RuleID:    r.ID(),
		Severity:  tagaudit.SeverityError,
		Message:   fmt.Sprintf("malformed struct tag on field %s: %v", fieldName, err),
		FieldName: fieldName,
	}}
}
