package rules

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/emm5317/tagaudit"
)

func TestShadowRule_ID(t *testing.T) {
	r := &ShadowRule{}
	if r.ID() != "shadow" {
		t.Errorf("expected 'shadow', got %q", r.ID())
	}
	if r.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestShadowRule_NoEmbedded(t *testing.T) {
	r := &ShadowRule{}
	tags1, _ := parseTag(`json:"name"`)

	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: fakeVar("Name"), Tags: tags1, RawTag: `json:"name"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings without embedded structs, got %d", len(findings))
	}
}

func TestShadowRule_DetectsShadow(t *testing.T) {
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

	tags1, _ := parseTag(`json:"name"`)

	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "Base", embeddedNamed, true), Tags: nil, RawTag: ""},
			{Field: fakeVar("DisplayName"), Tags: tags1, RawTag: `json:"name"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) != 1 {
		t.Fatalf("expected 1 shadow finding, got %d", len(findings))
	}
	if findings[0].FieldName != "DisplayName" {
		t.Errorf("expected finding on DisplayName, got %s", findings[0].FieldName)
	}
	if findings[0].TagKey != "json" {
		t.Errorf("expected tag key 'json', got %q", findings[0].TagKey)
	}
}

func TestShadowRule_NoShadowDifferentValues(t *testing.T) {
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

	tags1, _ := parseTag(`json:"display_name"`)

	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "Base", embeddedNamed, true), Tags: nil, RawTag: ""},
			{Field: fakeVar("DisplayName"), Tags: tags1, RawTag: `json:"display_name"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings when tag values differ, got %d", len(findings))
	}
}

func TestShadowRule_SkipsNonIdentifierTags(t *testing.T) {
	r := &ShadowRule{}

	embeddedFields := []*types.Var{
		types.NewVar(token.NoPos, nil, "A", types.Typ[types.String]),
	}
	embeddedTags := []string{`validate:"required"`}
	embeddedStruct := types.NewStruct(embeddedFields, embeddedTags)
	embeddedNamed := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Base", nil),
		embeddedStruct,
		nil,
	)

	tags1, _ := parseTag(`validate:"required"`)

	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "Base", embeddedNamed, true), Tags: nil, RawTag: ""},
			{Field: fakeVar("B"), Tags: tags1, RawTag: `validate:"required"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for non-identifier tag, got %d", len(findings))
	}
}

func TestShadowRule_SkipsDash(t *testing.T) {
	r := &ShadowRule{}

	embeddedFields := []*types.Var{
		types.NewVar(token.NoPos, nil, "Name", types.Typ[types.String]),
	}
	embeddedTags := []string{`json:"-"`}
	embeddedStruct := types.NewStruct(embeddedFields, embeddedTags)
	embeddedNamed := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Base", nil),
		embeddedStruct,
		nil,
	)

	tags1, _ := parseTag(`json:"-"`)

	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "Base", embeddedNamed, true), Tags: nil, RawTag: ""},
			{Field: fakeVar("Name"), Tags: tags1, RawTag: `json:"-"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for dash tags, got %d", len(findings))
	}
}

func TestShadowRule_EmbeddedViaPointer(t *testing.T) {
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
	ptrType := types.NewPointer(embeddedNamed)

	tags1, _ := parseTag(`json:"name"`)

	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "Base", ptrType, true), Tags: nil, RawTag: ""},
			{Field: fakeVar("Title"), Tags: tags1, RawTag: `json:"name"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) != 1 {
		t.Fatalf("expected 1 shadow finding via pointer embed, got %d", len(findings))
	}
}

func TestShadowRule_SkipsAnonymousFieldTags(t *testing.T) {
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

	// Only an embedded field, no direct fields to shadow
	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "Base", embeddedNamed, true), Tags: nil, RawTag: ""},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings when no direct fields to shadow, got %d", len(findings))
	}
}

func TestShadowRule_EmbeddedBadTag(t *testing.T) {
	r := &ShadowRule{}

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

func TestShadowRule_EmbeddedNonStruct(t *testing.T) {
	r := &ShadowRule{}

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
