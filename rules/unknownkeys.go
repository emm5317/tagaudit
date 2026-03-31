package rules

import (
	"fmt"

	"github.com/emm5317/tagaudit"
	"github.com/emm5317/tagaudit/internal/distance"
	"github.com/fatih/structtag"
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

			msg := fmt.Sprintf("field %s: unknown tag key %q (not in known keys list)", fieldName, tag.Key)

			var fix *tagaudit.SuggestedFix
			if suggestion, ok := distance.ClosestMatch(tag.Key, cfg.KnownTagKeys, 2); ok {
				msg += fmt.Sprintf(", did you mean %q?", suggestion)
				// Build a fix that replaces the misspelled key
				if tags, err := structtag.Parse(info.Tags.String()); err == nil {
					if oldTag, err := tags.Get(tag.Key); err == nil {
						tags.Delete(tag.Key)
						newTag := &structtag.Tag{
							Key:     suggestion,
							Name:    oldTag.Name,
							Options: oldTag.Options,
						}
						tags.Set(newTag)
						tagStart, tagEnd := tagSpanFromInfo(info)
						fix = &tagaudit.SuggestedFix{
							Description: fmt.Sprintf("rename tag key %s to %s", tag.Key, suggestion),
							NewTagValue: tags.String(),
							TagStart:    tagStart,
							TagEnd:      tagEnd,
						}
					}
				}
			}

			findings = append(findings, tagaudit.Finding{
				Pos:       pos,
				RuleID:    r.ID(),
				Severity:  tagaudit.SeverityWarning,
				Message:   msg,
				FieldName: fieldName,
				TagKey:    tag.Key,
				Fix:       fix,
			})
		}
	}

	return findings
}
