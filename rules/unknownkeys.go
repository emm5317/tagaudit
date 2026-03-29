package rules

import (
	"fmt"

	"github.com/emm5317/tagaudit"
)

// UnknownKeysRule reports tag keys not in the configured known set.
type UnknownKeysRule struct{}

func (r *UnknownKeysRule) ID() string          { return "unknownkeys" }
func (r *UnknownKeysRule) Description() string { return "reports unknown tag keys" }

func (r *UnknownKeysRule) CheckField(info tagaudit.FieldInfo, cfg *tagaudit.Config) []tagaudit.Finding {
	if info.Tags == nil || cfg == nil || len(cfg.KnownTagKeys) == 0 {
		return nil
	}

	known := make(map[string]bool, len(cfg.KnownTagKeys))
	for _, k := range cfg.KnownTagKeys {
		known[k] = true
	}

	var findings []tagaudit.Finding

	for _, tag := range info.Tags.Tags() {
		if !known[tag.Key] {
			var fieldName string
			if info.Field != nil {
				fieldName = info.Field.Name()
			}
			pos := posFromInfo(info)

			findings = append(findings, tagaudit.Finding{
				Pos:       pos,
				RuleID:    r.ID(),
				Severity:  tagaudit.SeverityWarning,
				Message:   fmt.Sprintf("field %s: unknown tag key %q (not in known keys list)", fieldName, tag.Key),
				FieldName: fieldName,
				TagKey:    tag.Key,
			})
		}
	}

	return findings
}
