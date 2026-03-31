package rules

import (
	"fmt"
	"go/types"

	"github.com/emm5317/tagaudit"
	"github.com/fatih/structtag"
)

// DuplicatesRule detects duplicate tag values within a struct, including
// values from embedded struct fields. Only checks tags where the value
// acts as a field name/identifier (encoding tags like json, xml, db, etc.),
// not constraint tags where sharing values is expected (e.g., validate, river).
type DuplicatesRule struct{}

func (r *DuplicatesRule) ID() string          { return "duplicates" }
func (r *DuplicatesRule) Description() string { return "detects duplicate tag values" }

// DefaultIdentifierTagKeys are tag keys where the value represents a unique
// field name. Duplicate values in these tags indicate a real conflict.
var DefaultIdentifierTagKeys = []string{
	"json", "xml", "yaml", "toml", "db", "bson", "csv", "avro", "parquet",
}

// identifierTagSet returns the set of identifier tag keys to use,
// preferring cfg.IdentifierTagKeys if set, else DefaultIdentifierTagKeys.
func identifierTagSet(cfg *tagaudit.Config) map[string]bool {
	keys := DefaultIdentifierTagKeys
	if cfg != nil && len(cfg.IdentifierTagKeys) > 0 {
		keys = cfg.IdentifierTagKeys
	}
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}

type tagEntry struct {
	fieldName string
	depth     int // 0 = direct field, 1+ = embedded
	info      tagaudit.FieldInfo
}

func (r *DuplicatesRule) CheckStruct(info tagaudit.StructInfo, cfg *tagaudit.Config) []tagaudit.Finding {
	idTags := identifierTagSet(cfg)

	// Map of tagKey -> tagValue -> []tagEntry
	seen := make(map[string]map[string][]tagEntry)

	// Collect tags from direct fields (only identifier tags)
	for _, f := range info.Fields {
		if f.Tags == nil {
			continue
		}
		for _, tag := range f.Tags.Tags() {
			if tag.Name == "" || tag.Name == "-" {
				continue
			}
			if !idTags[tag.Key] {
				continue
			}
			if seen[tag.Key] == nil {
				seen[tag.Key] = make(map[string][]tagEntry)
			}
			var fieldName string
			if f.Field != nil {
				fieldName = f.Field.Name()
			}
			seen[tag.Key][tag.Name] = append(seen[tag.Key][tag.Name], tagEntry{
				fieldName: fieldName,
				depth:     0,
				info:      f,
			})
		}
	}

	// Collect tags from embedded structs
	for _, f := range info.Fields {
		if f.Field == nil || !f.Field.Anonymous() {
			continue
		}
		collectEmbeddedTags(f.Field.Type(), seen, 1, make(map[*types.Named]bool), idTags)
	}

	// Report duplicates
	var findings []tagaudit.Finding
	for tagKey, names := range seen {
		for tagName, entries := range names {
			if len(entries) < 2 {
				continue
			}
			// Only report for entries that have position info (direct fields)
			for _, e := range entries {
				if e.depth > 0 {
					continue // don't report on embedded fields — report on the struct that embeds them
				}
				otherFields := make([]string, 0, len(entries)-1)
				for _, other := range entries {
					if other.fieldName != e.fieldName {
						otherFields = append(otherFields, other.fieldName)
					}
				}
				if len(otherFields) == 0 {
					continue
				}
				pos := posFromInfo(e.info)
				findings = append(findings, tagaudit.Finding{
					Pos:       pos,
					RuleID:    r.ID(),
					Severity:  tagaudit.SeverityWarning,
					Message:   fmt.Sprintf("field %s: duplicate %s tag value %q (also on: %v)", e.fieldName, tagKey, tagName, otherFields),
					FieldName: e.fieldName,
					TagKey:    tagKey,
				})
			}
		}
	}

	return findings
}

// collectEmbeddedTags recursively collects tag values from embedded struct types.
// visited tracks named types to prevent infinite recursion on cyclic definitions.
func collectEmbeddedTags(t types.Type, seen map[string]map[string][]tagEntry, depth int, visited map[*types.Named]bool, idTags map[string]bool) {
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
			collectEmbeddedTags(field.Type(), seen, depth+1, visited, idTags)
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
			if seen[tag.Key] == nil {
				seen[tag.Key] = make(map[string][]tagEntry)
			}
			seen[tag.Key][tag.Name] = append(seen[tag.Key][tag.Name], tagEntry{
				fieldName: field.Name(),
				depth:     depth,
			})
		}
	}
}
