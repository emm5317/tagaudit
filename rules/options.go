package rules

import (
	"fmt"

	"github.com/emm5317/tagaudit"
	"github.com/emm5317/tagaudit/internal/distance"
	"github.com/fatih/structtag"
)

// Known valid options for well-known tag keys.
var knownOptions = map[string]map[string]bool{
	"json": {
		"omitempty": true,
		"string":    true,
		"omitzero":  true, // Go 1.24+
	},
	"xml": {
		"attr":      true,
		"chardata":  true,
		"cdata":     true,
		"innerxml":  true,
		"comment":   true,
		"any":       true,
		"omitempty": true,
	},
	"yaml": {
		"omitempty": true,
		"flow":      true,
		"inline":    true,
	},
}

// OptionsRule validates that options for well-known tags are valid.
type OptionsRule struct{}

func (r *OptionsRule) ID() string          { return "options" }
func (r *OptionsRule) Description() string { return "validates options for well-known tag keys" }

func (r *OptionsRule) CheckField(info tagaudit.FieldInfo, cfg *tagaudit.Config) []tagaudit.Finding {
	if info.Tags == nil {
		return nil
	}

	var findings []tagaudit.Finding

	for _, tag := range info.Tags.Tags() {
		validOpts, ok := knownOptions[tag.Key]
		if !ok && (cfg == nil || cfg.KnownOptions == nil || len(cfg.KnownOptions[tag.Key]) == 0) {
			continue // not a well-known tag and no user-provided options, skip
		}
		if !ok {
			validOpts = make(map[string]bool)
		}
		// Merge user-provided options if present
		if cfg != nil && cfg.KnownOptions != nil {
			if extra, has := cfg.KnownOptions[tag.Key]; has {
				// Copy built-in map to avoid mutating the package-level variable
				merged := make(map[string]bool, len(validOpts)+len(extra))
				for k, v := range validOpts {
					merged[k] = v
				}
				for _, o := range extra {
					merged[o] = true
				}
				validOpts = merged
			}
		}

		// Build candidate list for Levenshtein matching
		var validOptsList []string
		for k := range validOpts {
			validOptsList = append(validOptsList, k)
		}

		for _, opt := range tag.Options {
			if opt == "" {
				continue
			}
			if !validOpts[opt] {
				var fieldName string
				if info.Field != nil {
					fieldName = info.Field.Name()
				}
				pos := posFromInfo(info)

				msg := fmt.Sprintf("field %s: unknown option %q for %s tag", fieldName, opt, tag.Key)

				var fix *tagaudit.SuggestedFix
				if suggestion, ok := distance.ClosestMatch(opt, validOptsList, 2); ok {
					msg += fmt.Sprintf(", did you mean %q?", suggestion)
					// Build a fix that replaces the misspelled option
					if tags, err := structtag.Parse(info.Tags.String()); err == nil {
						if t, err := tags.Get(tag.Key); err == nil {
							for i, o := range t.Options {
								if o == opt {
									t.Options[i] = suggestion
									break
								}
							}
							tags.Set(t)
							tagStart, tagEnd := tagSpanFromInfo(info)
							fix = &tagaudit.SuggestedFix{
								Description: fmt.Sprintf("fix option %s to %s", opt, suggestion),
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
	}

	return findings
}
