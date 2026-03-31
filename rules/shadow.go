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

func (r *ShadowRule) CheckStruct(info tagaudit.StructInfo, cfg *tagaudit.Config) []tagaudit.Finding {
	idTags := identifierTagSet(cfg)

	// Collect tag values from embedded structs
	// Map of tagKey -> tagValue -> embedded field name
	embeddedTags := make(map[string]map[string]string)

	for _, f := range info.Fields {
		if f.Field == nil || !f.Field.Anonymous() {
			continue
		}
		collectEmbeddedTagNames(f.Field.Type(), embeddedTags, 1, make(map[*types.Named]bool), idTags)
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
			if !idTags[tag.Key] {
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
// visited tracks named types to prevent infinite recursion on cyclic definitions.
func collectEmbeddedTagNames(t types.Type, out map[string]map[string]string, depth int, visited map[*types.Named]bool, idTags map[string]bool) {
	if depth > maxEmbedDepth {
		return
	}

	named, st, ok := unwrapToStruct(t)
	if !ok {
		return
	}
	if named != nil {
		if visited[named] {
			return
		}
		visited[named] = true
	}

	for i := range st.NumFields() {
		field := st.Field(i)
		rawTag := st.Tag(i)

		if field.Anonymous() {
			collectEmbeddedTagNames(field.Type(), out, depth+1, visited, idTags)
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
			if !idTags[tag.Key] {
				continue
			}
			if out[tag.Key] == nil {
				out[tag.Key] = make(map[string]string)
			}
			out[tag.Key][tag.Name] = field.Name()
		}
	}
}
