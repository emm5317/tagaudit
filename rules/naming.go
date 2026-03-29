package rules

import (
	"fmt"

	"github.com/emm5317/tagaudit"
	"github.com/emm5317/tagaudit/internal/naming"
	"github.com/fatih/structtag"
)

// NamingRule enforces naming conventions on tag values.
type NamingRule struct{}

func (r *NamingRule) ID() string          { return "naming" }
func (r *NamingRule) Description() string { return "enforces naming conventions on tag values" }

func (r *NamingRule) CheckField(info tagaudit.FieldInfo, cfg *tagaudit.Config) []tagaudit.Finding {
	if info.Tags == nil || cfg == nil || len(cfg.NamingConventions) == 0 {
		return nil
	}

	var findings []tagaudit.Finding

	for key, convention := range cfg.NamingConventions {
		tag, err := info.Tags.Get(key)
		if err != nil {
			continue // tag key not present on this field
		}

		name := tag.Name
		if name == "" || name == "-" {
			continue // field is skipped
		}

		if !naming.MatchesConvention(name, convention) {
			expected := naming.Convert(name, convention)
			pos := posFromInfo(info)

			var fix *tagaudit.SuggestedFix
			if tags, err := structtag.Parse(info.RawTag); err == nil {
				if t, err := tags.Get(key); err == nil {
					t.Name = expected
					tags.Set(t)
					fix = &tagaudit.SuggestedFix{
						Description: fmt.Sprintf("rename %s tag to %s", key, expected),
						NewTagValue: tags.String(),
					}
				}
			}

			findings = append(findings, tagaudit.Finding{
				Pos:       pos,
				RuleID:    r.ID(),
				Severity:  tagaudit.SeverityInfo,
				Message:   fmt.Sprintf("field %s: %s tag %q does not follow %s convention, expected %q", info.Field.Name(), key, name, convention, expected),
				FieldName: info.Field.Name(),
				TagKey:    key,
				Fix:       fix,
			})
		}
	}

	return findings
}
