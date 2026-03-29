package rules

import (
	"fmt"
	"go/types"

	"github.com/emm5317/tagaudit"
	"github.com/fatih/structtag"
)

// ShadowRule detects when an outer struct field's tag value shadows
// an embedded struct field's tag value.
type ShadowRule struct{}

func (r *ShadowRule) ID() string          { return "shadow" }
func (r *ShadowRule) Description() string { return "detects shadowed tags in embedded structs" }

func (r *ShadowRule) CheckStruct(info tagaudit.StructInfo, _ *tagaudit.Config) []tagaudit.Finding {
	// Collect tag values from embedded structs
	// Map of tagKey -> tagValue -> embedded field name
	embeddedTags := make(map[string]map[string]string)

	for _, f := range info.Fields {
		if f.Field == nil || !f.Field.Anonymous() {
			continue
		}
		collectEmbeddedTagNames(f.Field.Type(), embeddedTags)
	}

	if len(embeddedTags) == 0 {
		return nil
	}

	// Check direct fields for shadows
	var findings []tagaudit.Finding

	for _, f := range info.Fields {
		if f.Tags == nil || f.Field == nil || f.Field.Anonymous() {
			continue
		}

		for _, tag := range f.Tags.Tags() {
			if tag.Name == "" || tag.Name == "-" {
				continue
			}
			if embFieldName, ok := embeddedTags[tag.Key][tag.Name]; ok {
				pos := posFromInfo(f)
				findings = append(findings, tagaudit.Finding{
					Pos:       pos,
					RuleID:    r.ID(),
					Severity:  tagaudit.SeverityWarning,
					Message:   fmt.Sprintf("field %s: %s tag %q shadows embedded field %s with the same tag value", f.Field.Name(), tag.Key, tag.Name, embFieldName),
					FieldName: f.Field.Name(),
					TagKey:    tag.Key,
				})
			}
		}
	}

	return findings
}

// collectEmbeddedTagNames collects tag key/value -> field name from embedded structs.
func collectEmbeddedTagNames(t types.Type, out map[string]map[string]string) {
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}
	if named, ok := t.(*types.Named); ok {
		t = named.Underlying()
	}
	st, ok := t.(*types.Struct)
	if !ok {
		return
	}

	for i := range st.NumFields() {
		field := st.Field(i)
		rawTag := st.Tag(i)

		if field.Anonymous() {
			collectEmbeddedTagNames(field.Type(), out)
			continue
		}

		if rawTag == "" {
			continue
		}

		tags, err := structtag.Parse(rawTag)
		if err != nil {
			continue
		}

		for _, tag := range tags.Tags() {
			if tag.Name == "" || tag.Name == "-" {
				continue
			}
			if out[tag.Key] == nil {
				out[tag.Key] = make(map[string]string)
			}
			out[tag.Key][tag.Name] = field.Name()
		}
	}
}
