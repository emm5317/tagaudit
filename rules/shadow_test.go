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

func TestShadowRule_ThreeLevelEmbedding(t *testing.T) {
	r := &ShadowRule{}

	// Level 3
	l3Fields := []*types.Var{types.NewVar(token.NoPos, nil, "Name", types.Typ[types.String])}
	l3Tags := []string{`json:"name"`}
	l3Struct := types.NewStruct(l3Fields, l3Tags)
	l3Named := types.NewNamed(types.NewTypeName(token.NoPos, nil, "L3", nil), l3Struct, nil)

	// Level 2 embeds Level 3
	l2Fields := []*types.Var{types.NewField(token.NoPos, nil, "L3", l3Named, true)}
	l2Tags := []string{""}
	l2Struct := types.NewStruct(l2Fields, l2Tags)
	l2Named := types.NewNamed(types.NewTypeName(token.NoPos, nil, "L2", nil), l2Struct, nil)

	tags1, _ := parseTag(`json:"name"`)
	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "L2", l2Named, true), Tags: nil, RawTag: ""},
			{Field: fakeVar("Title"), Tags: tags1, RawTag: `json:"name"`},
		},
	}

	findings := r.CheckStruct(info, nil)
	if len(findings) == 0 {
		t.Error("expected shadow finding with 3-level embedded struct")
	}
}

func TestShadowRule_CyclicEmbedding(t *testing.T) {
	r := &ShadowRule{}

	aName := types.NewTypeName(token.NoPos, nil, "A", nil)
	aNamed := types.NewNamed(aName, nil, nil)
	bName := types.NewTypeName(token.NoPos, nil, "B", nil)
	bNamed := types.NewNamed(bName, nil, nil)

	bStruct := types.NewStruct(
		[]*types.Var{types.NewField(token.NoPos, nil, "A", aNamed, true)},
		[]string{""},
	)
	bNamed.SetUnderlying(bStruct)

	aStruct := types.NewStruct(
		[]*types.Var{types.NewField(token.NoPos, nil, "B", bNamed, true)},
		[]string{""},
	)
	aNamed.SetUnderlying(aStruct)

	info := tagaudit.StructInfo{
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewField(token.NoPos, nil, "A", aNamed, true), Tags: nil, RawTag: ""},
		},
	}

	// Should not panic or infinite loop
	findings := r.CheckStruct(info, nil)
	_ = findings
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
