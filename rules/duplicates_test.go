package rules

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/emm5317/tagaudit"
)

func TestDuplicatesRule_ID(t *testing.T) {
	r := &DuplicatesRule{}
	if r.ID() != "duplicates" {
		t.Errorf("expected 'duplicates', got %q", r.ID())
	}
	if r.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestDuplicatesRule_NoDuplicates(t *testing.T) {
	r := &DuplicatesRule{}
	tags1, _ := parseTag(`json:"name"`)
	tags2, _ := parseTag(`json:"email"`)

	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: fakeVar("Name"), Tags: tags1, RawTag: `json:"name"`},
			{Field: fakeVar("Email"), Tags: tags2, RawTag: `json:"email"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestDuplicatesRule_DuplicateJsonTag(t *testing.T) {
	r := &DuplicatesRule{}
	tags1, _ := parseTag(`json:"name"`)
	tags2, _ := parseTag(`json:"name"`)

	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: fakeVar("Name"), Tags: tags1, RawTag: `json:"name"`},
			{Field: fakeVar("DisplayName"), Tags: tags2, RawTag: `json:"name"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) < 1 {
		t.Fatalf("expected at least 1 finding, got %d", len(findings))
	}
}

func TestDuplicatesRule_SkipsNonIdentifierTags(t *testing.T) {
	r := &DuplicatesRule{}
	tags1, _ := parseTag(`river:"unique"`)
	tags2, _ := parseTag(`river:"unique"`)

	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: fakeVar("Season"), Tags: tags1, RawTag: `river:"unique"`},
			{Field: fakeVar("Year"), Tags: tags2, RawTag: `river:"unique"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for non-identifier tag, got %d", len(findings))
	}
}

func TestDuplicatesRule_SkipsDash(t *testing.T) {
	r := &DuplicatesRule{}
	tags1, _ := parseTag(`json:"-"`)
	tags2, _ := parseTag(`json:"-"`)

	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: fakeVar("A"), Tags: tags1, RawTag: `json:"-"`},
			{Field: fakeVar("B"), Tags: tags2, RawTag: `json:"-"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for skipped fields, got %d", len(findings))
	}
}

func TestDuplicatesRule_NilTags(t *testing.T) {
	r := &DuplicatesRule{}
	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: fakeVar("A"), Tags: nil},
			{Field: fakeVar("B"), Tags: nil},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for nil tags, got %d", len(findings))
	}
}

func TestDuplicatesRule_WithEmbeddedStruct(t *testing.T) {
	r := &DuplicatesRule{}

	// Build an embedded struct type with a json:"name" field
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

	tags1, _ := parseTag(`json:"name"`)

	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "Base", embeddedNamed, true), Tags: nil, RawTag: ""},
			{Field: fakeVar("DisplayName"), Tags: tags1, RawTag: `json:"name"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) == 0 {
		t.Error("expected at least 1 finding for duplicate with embedded struct")
	}
}

func TestDuplicatesRule_EmbeddedViaPointer(t *testing.T) {
	r := &DuplicatesRule{}

	embeddedFields := []*types.Var{
		types.NewVar(token.NoPos, nil, "ID", types.Typ[types.Int]),
	}
	embeddedTags := []string{`json:"id"`}
	embeddedStruct := types.NewStruct(embeddedFields, embeddedTags)
	embeddedNamed := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Base", nil),
		embeddedStruct,
		nil,
	)
	ptrType := types.NewPointer(embeddedNamed)

	tags1, _ := parseTag(`json:"id"`)

	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "Base", ptrType, true), Tags: nil, RawTag: ""},
			{Field: fakeVar("MyID"), Tags: tags1, RawTag: `json:"id"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) == 0 {
		t.Error("expected at least 1 finding for duplicate with pointer-embedded struct")
	}
}

func TestDuplicatesRule_EmbeddedNoTag(t *testing.T) {
	r := &DuplicatesRule{}

	embeddedFields := []*types.Var{
		types.NewVar(token.NoPos, nil, "Name", types.Typ[types.String]),
	}
	embeddedTags := []string{""} // no tag
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

func TestDuplicatesRule_EmbeddedNonStruct(t *testing.T) {
	r := &DuplicatesRule{}

	// Embed a non-struct type (e.g., an interface)
	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "Stringer", types.Typ[types.String], true), Tags: nil, RawTag: ""},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for non-struct embed, got %d", len(findings))
	}
}

func TestDuplicatesRule_EmbeddedRecursive(t *testing.T) {
	r := &DuplicatesRule{}

	// Inner embedded struct
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

	// Outer embedded struct that embeds Inner
	outerFields := []*types.Var{
		types.NewField(token.NoPos, nil, "Inner", innerNamed, true),
	}
	outerTags := []string{""}
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
		t.Error("expected finding for duplicate with recursively embedded struct")
	}
}

func TestDuplicatesRule_EmbeddedBadTag(t *testing.T) {
	r := &DuplicatesRule{}

	embeddedFields := []*types.Var{
		types.NewVar(token.NoPos, nil, "Bad", types.Typ[types.String]),
	}
	embeddedTags := []string{`json:bad_syntax`}
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
		t.Errorf("expected 0 findings for bad embedded tag, got %d", len(findings))
	}
}
