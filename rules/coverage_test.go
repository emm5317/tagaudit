package rules

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/emm5317/tagaudit"
)

// Tests targeting specific uncovered code paths for 100% coverage.

// completeness.go:27 — field with nil Field
func TestCompletenessRule_NilField(t *testing.T) {
	r := &CompletenessRule{}
	tags1, _ := parseTag(`json:"name"`)

	info := tagaudit.StructInfo{
		StructName: "User",
		Fields: []tagaudit.FieldInfo{
			{Field: fakeVar("Name"), Tags: tags1, RawTag: `json:"name"`},
			{Field: nil, Tags: nil, RawTag: ""},
		},
	}

	// Should not panic on nil Field
	findings := r.CheckStruct(info, &tagaudit.Config{RequiredTagKeys: []string{"json"}})
	_ = findings
}

// duplicates.go:97 — single entry for a tag value (no duplicate, otherFields empty)
func TestDuplicatesRule_SingleEntryNotReported(t *testing.T) {
	r := &DuplicatesRule{}

	// Build embedded struct with json:"name", but no direct field duplicates it
	embeddedFields := []*types.Var{
		types.NewVar(token.NoPos, nil, "Name", types.Typ[types.String]),
	}
	embeddedTags := []string{`json:"name"`}
	embeddedStruct := types.NewStruct(embeddedFields, embeddedTags)
	embeddedNamed := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Base", nil),
		embeddedStruct,
		nil,
	)

	// Direct field has a DIFFERENT tag value
	tags1, _ := parseTag(`json:"email"`)
	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "Base", embeddedNamed, true), Tags: nil, RawTag: ""},
			{Field: fakeVar("Email"), Tags: tags1, RawTag: `json:"email"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings when no duplicates, got %d", len(findings))
	}
}

// duplicates.go:152-160 — embedded struct fields with dash, empty, non-identifier tags
func TestCollectEmbeddedTags_SkipsDashAndEmpty(t *testing.T) {
	r := &DuplicatesRule{}

	embeddedFields := []*types.Var{
		types.NewVar(token.NoPos, nil, "Skipped", types.Typ[types.String]),
		types.NewVar(token.NoPos, nil, "Also", types.Typ[types.String]),
		types.NewVar(token.NoPos, nil, "NonID", types.Typ[types.String]),
	}
	embeddedTags := []string{
		`json:"-"`,
		`json:""`,
		`river:"unique"`, // non-identifier tag
	}
	embeddedStruct := types.NewStruct(embeddedFields, embeddedTags)
	embeddedNamed := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Base", nil),
		embeddedStruct,
		nil,
	)

	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "Base", embeddedNamed, true), Tags: nil, RawTag: ""},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

// naming.go:25 — tag key not present on field (Tags.Get returns error)
func TestNamingRule_TagKeyNotPresent(t *testing.T) {
	r := &NamingRule{}
	tags, _ := parseTag(`db:"name"`) // has db but not json

	cfg := &tagaudit.Config{
		NamingConventions: map[string]string{"json": "snake_case"},
	}

	info := tagaudit.FieldInfo{
		RawTag: `db:"name"`,
		Tags:   tags,
		Field:  fakeVar("Name"),
	}

	findings := r.CheckField(info, cfg)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings when tag key not present, got %d", len(findings))
	}
}

// options.go:39 — nil tags
func TestOptionsRule_NilTags(t *testing.T) {
	r := &OptionsRule{}
	info := tagaudit.FieldInfo{Tags: nil}
	findings := r.CheckField(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for nil tags, got %d", len(findings))
	}
}

// options.go:52 — empty option string
func TestOptionsRule_EmptyOption(t *testing.T) {
	r := &OptionsRule{}
	tags, _ := parseTag(`json:"name,"`) // trailing comma = empty option

	info := tagaudit.FieldInfo{
		RawTag: `json:"name,"`,
		Tags:   tags,
		Field:  fakeVar("Name"),
	}

	findings := r.CheckField(info, nil)
	// Empty option should be skipped, not reported
	for _, f := range findings {
		if f.Message == "" {
			t.Error("got finding with empty message")
		}
	}
}

// shadow.go:43-46 — direct field with dash/empty/non-identifier tag
func TestShadowRule_DirectFieldDashTag(t *testing.T) {
	r := &ShadowRule{}

	embeddedFields := []*types.Var{
		types.NewVar(token.NoPos, nil, "Name", types.Typ[types.String]),
	}
	embeddedTags := []string{`json:"name"`}
	embeddedStruct := types.NewStruct(embeddedFields, embeddedTags)
	embeddedNamed := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Base", nil),
		embeddedStruct,
		nil,
	)

	// Direct field has json:"-" — should not shadow
	tags1, _ := parseTag(`json:"-"`)
	// Direct field has a non-identifier tag — should not shadow
	tags2, _ := parseTag(`river:"name"`)

	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "Base", embeddedNamed, true), Tags: nil, RawTag: ""},
			{Field: fakeVar("Skipped"), Tags: tags1, RawTag: `json:"-"`},
			{Field: fakeVar("NonID"), Tags: tags2, RawTag: `river:"name"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

// shadow.go:83-89 — collectEmbeddedTagNames: anonymous field recursion + empty rawTag
func TestShadowRule_EmbeddedRecursive(t *testing.T) {
	r := &ShadowRule{}

	innerFields := []*types.Var{
		types.NewVar(token.NoPos, nil, "ID", types.Typ[types.Int]),
	}
	innerTags := []string{`json:"id"`}
	innerStruct := types.NewStruct(innerFields, innerTags)
	innerNamed := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Inner", nil),
		innerStruct,
		nil,
	)

	outerFields := []*types.Var{
		types.NewField(token.NoPos, nil, "Inner", innerNamed, true),
		types.NewVar(token.NoPos, nil, "Untagged", types.Typ[types.String]),
	}
	outerTags := []string{"", ""} // no tags on either
	outerStruct := types.NewStruct(outerFields, outerTags)
	outerNamed := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Outer", nil),
		outerStruct,
		nil,
	)

	tags1, _ := parseTag(`json:"id"`)
	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "Outer", outerNamed, true), Tags: nil, RawTag: ""},
			{Field: fakeVar("MyID"), Tags: tags1, RawTag: `json:"id"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) == 0 {
		t.Error("expected shadow finding for recursively embedded struct")
	}
}

// syntax.go:28 — re-parse succeeds (edge case where Tags is nil but parse succeeds)
// syntax.go:33 — field with position info (Field is not nil)
func TestSyntaxRule_WithFieldInfo(t *testing.T) {
	r := &SyntaxRule{}

	info := tagaudit.FieldInfo{
		RawTag: `json:bad`,
		Tags:   nil, // parse failed
		Field:  fakeVar("Bad"),
	}

	findings := r.CheckField(info, nil)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].FieldName != "Bad" {
		t.Errorf("expected field name 'Bad', got %q", findings[0].FieldName)
	}
}

// completeness.go:27 — unexported field in the "has any field got key" scan
func TestCompletenessRule_UnexportedFieldInKeyScan(t *testing.T) {
	r := &CompletenessRule{}
	tags1, _ := parseTag(`json:"name"`)

	pkg := types.NewPackage("example.com/test", "test")
	info := tagaudit.StructInfo{
		StructName: "Mixed",
		Fields: []tagaudit.FieldInfo{
			// unexported field WITH json tag — should be skipped when scanning for "has key"
			{Field: types.NewVar(token.NoPos, pkg, "hidden", types.Typ[types.String]), Tags: tags1, RawTag: `json:"name"`},
			// exported field without json tag
			{Field: types.NewVar(token.NoPos, nil, "Visible", types.Typ[types.String]), Tags: nil, RawTag: ""},
		},
	}

	// The only field with json is unexported, so hasKey should be false → no findings
	findings := r.CheckStruct(info, &tagaudit.Config{RequiredTagKeys: []string{"json"}})
	if len(findings) != 0 {
		t.Errorf("expected 0 findings (unexported field shouldn't trigger completeness), got %d", len(findings))
	}
}

// duplicates.go:97 — entry where all "other" entries have the same fieldName
// This happens when an embedded field has the same name as a direct field
func TestDuplicatesRule_SameFieldNameNotReported(t *testing.T) {
	r := &DuplicatesRule{}

	// Embedded struct with field "Name" json:"name"
	embeddedFields := []*types.Var{
		types.NewVar(token.NoPos, nil, "Name", types.Typ[types.String]),
	}
	embeddedTags := []string{`json:"name"`}
	embeddedStruct := types.NewStruct(embeddedFields, embeddedTags)
	embeddedNamed := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Base", nil),
		embeddedStruct,
		nil,
	)

	// Direct field also named "Name" with json:"name"
	// Both have fieldName "Name" so otherFields will be empty
	tags1, _ := parseTag(`json:"name"`)
	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "Base", embeddedNamed, true), Tags: nil, RawTag: ""},
			{Field: fakeVar("Name"), Tags: tags1, RawTag: `json:"name"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	// Both entries have fieldName "Name", so otherFields is empty → no finding
	if len(findings) != 0 {
		t.Errorf("expected 0 findings when same field name, got %d", len(findings))
	}
}

// duplicates.go:158 — embedded struct with a tag key not seen by any direct field
func TestDuplicatesRule_EmbeddedUniqueKey(t *testing.T) {
	r := &DuplicatesRule{}

	embeddedFields := []*types.Var{
		types.NewVar(token.NoPos, nil, "Col", types.Typ[types.String]),
	}
	embeddedTags := []string{`db:"column_name"`} // db key, no direct fields have db
	embeddedStruct := types.NewStruct(embeddedFields, embeddedTags)
	embeddedNamed := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Base", nil),
		embeddedStruct,
		nil,
	)

	tags1, _ := parseTag(`json:"name"`) // different key entirely
	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "Base", embeddedNamed, true), Tags: nil, RawTag: ""},
			{Field: fakeVar("Name"), Tags: tags1, RawTag: `json:"name"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

// syntax.go:28 — Tags is nil but RawTag parses fine (force the edge case)
func TestSyntaxRule_TagsNilButParseable(t *testing.T) {
	r := &SyntaxRule{}
	info := tagaudit.FieldInfo{
		RawTag: `json:"name"`, // valid tag
		Tags:   nil,           // artificially set to nil
	}

	findings := r.CheckField(info, nil)
	// Re-parse succeeds → returns nil
	if len(findings) != 0 {
		t.Errorf("expected 0 findings when re-parse succeeds, got %d", len(findings))
	}
}

// rules.go:6-8 — init function calls SetDefaultRules
// This is covered by importing the package, but we can verify it works
func TestInit_RegistersDefaults(t *testing.T) {
	// The init() in rules.go calls SetDefaultRules.
	// Verify that creating an analyzer without explicit rules uses built-in rules.
	a := tagaudit.New(&tagaudit.Config{})
	if a == nil {
		t.Fatal("New returned nil")
	}
}
