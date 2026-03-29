package rules

import (
	"fmt"

	"github.com/emm5317/tagaudit"
)

// CompletenessRule checks that if any exported field has a required tag key,
// all exported fields in the struct have it too.
type CompletenessRule struct{}

func (r *CompletenessRule) ID() string          { return "completeness" }
func (r *CompletenessRule) Description() string { return "checks for missing required tag keys" }

func (r *CompletenessRule) CheckStruct(info tagaudit.StructInfo, cfg *tagaudit.Config) []tagaudit.Finding {
	if cfg == nil || len(cfg.RequiredTagKeys) == 0 {
		return nil
	}

	var findings []tagaudit.Finding

	for _, reqKey := range cfg.RequiredTagKeys {
		// Check if any exported field has this tag key
		hasKey := false
		for _, f := range info.Fields {
			if f.Field == nil || !f.Field.Exported() {
				continue
			}
			if f.Tags != nil {
				if _, err := f.Tags.Get(reqKey); err == nil {
					hasKey = true
					break
				}
			}
		}

		if !hasKey {
			continue // no field has this key, nothing to enforce
		}

		// Now flag exported fields missing the key
		for _, f := range info.Fields {
			if f.Field == nil || !f.Field.Exported() {
				continue
			}
			if f.Field.Anonymous() {
				continue // skip embedded fields
			}

			hasTag := false
			if f.Tags != nil {
				if _, err := f.Tags.Get(reqKey); err == nil {
					hasTag = true
				}
			}

			if !hasTag {
				pos := posFromInfo(f)
				findings = append(findings, tagaudit.Finding{
					Pos:       pos,
					RuleID:    r.ID(),
					Severity:  tagaudit.SeverityInfo,
					Message:   fmt.Sprintf("field %s in %s: missing %s tag (other fields in this struct have it)", f.Field.Name(), info.StructName, reqKey),
					FieldName: f.Field.Name(),
					TagKey:    reqKey,
				})
			}
		}
	}

	return findings
}
