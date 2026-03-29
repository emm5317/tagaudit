package rules

import (
	"fmt"

	"github.com/emm5317/tagaudit"
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

func (r *OptionsRule) CheckField(info tagaudit.FieldInfo, _ *tagaudit.Config) []tagaudit.Finding {
	if info.Tags == nil {
		return nil
	}

	var findings []tagaudit.Finding

	for _, tag := range info.Tags.Tags() {
		validOpts, ok := knownOptions[tag.Key]
		if !ok {
			continue // not a well-known tag, skip
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

				findings = append(findings, tagaudit.Finding{
					Pos:       pos,
					RuleID:    r.ID(),
					Severity:  tagaudit.SeverityWarning,
					Message:   fmt.Sprintf("field %s: unknown option %q for %s tag", fieldName, opt, tag.Key),
					FieldName: fieldName,
					TagKey:    tag.Key,
				})
			}
		}
	}

	return findings
}
