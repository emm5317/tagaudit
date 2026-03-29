package rules

import (
	"fmt"

	"github.com/emm5317/tagaudit"
)

// Tags on unexported fields that are silently ignored by their encoders.
var encodingTagKeys = map[string]bool{
	"json": true,
	"xml":  true,
	"yaml": true,
	"toml": true,
}

// UnexportedRule reports encoding tags on unexported fields.
type UnexportedRule struct{}

func (r *UnexportedRule) ID() string          { return "unexported" }
func (r *UnexportedRule) Description() string { return "reports encoding tags on unexported fields" }

func (r *UnexportedRule) CheckField(info tagaudit.FieldInfo, _ *tagaudit.Config) []tagaudit.Finding {
	if info.Tags == nil || info.Field == nil {
		return nil
	}

	// Only flag unexported fields
	if info.Field.Exported() {
		return nil
	}

	var findings []tagaudit.Finding

	for _, tag := range info.Tags.Tags() {
		if encodingTagKeys[tag.Key] {
			pos := posFromInfo(info)
			findings = append(findings, tagaudit.Finding{
				Pos:       pos,
				RuleID:    r.ID(),
				Severity:  tagaudit.SeverityWarning,
				Message:   fmt.Sprintf("field %s: %s tag on unexported field is ignored by encoding/%s", info.Field.Name(), tag.Key, tag.Key),
				FieldName: info.Field.Name(),
				TagKey:    tag.Key,
			})
		}
	}

	return findings
}
