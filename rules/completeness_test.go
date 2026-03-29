package rules

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/emm5317/tagaudit"
)

func TestCompletenessRule_ID(t *testing.T) {
	r := &CompletenessRule{}
	if r.ID() != "completeness" {
		t.Errorf("expected 'completeness', got %q", r.ID())
	}
	if r.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestCompletenessRule_NilConfig(t *testing.T) {
	r := &CompletenessRule{}
	findings := r.CheckStruct(tagaudit.StructInfo{}, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for nil config, got %d", len(findings))
	}
}

func TestCompletenessRule_NoRequiredKeys(t *testing.T) {
	r := &CompletenessRule{}
	findings := r.CheckStruct(tagaudit.StructInfo{}, &tagaudit.Config{})
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for empty required keys, got %d", len(findings))
	}
}

func TestCompletenessRule_AllTagged(t *testing.T) {
	r := &CompletenessRule{}
	tags1, _ := parseTag(`json:"name"`)
	tags2, _ := parseTag(`json:"email"`)

	info := tagaudit.StructInfo{
		StructName: "User",
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewVar(token.NoPos, nil, "Name", types.Typ[types.String]), Tags: tags1, RawTag: `json:"name"`},
			{Field: types.NewVar(token.NoPos, nil, "Email", types.Typ[types.String]), Tags: tags2, RawTag: `json:"email"`},
		},
	}

	findings := r.CheckStruct(info, &tagaudit.Config{RequiredTagKeys: []string{"json"}})
	if len(findings) != 0 {
		t.Errorf("expected 0 findings when all fields tagged, got %d", len(findings))
	}
}

func TestCompletenessRule_MissingTag(t *testing.T) {
	r := &CompletenessRule{}
	tags1, _ := parseTag(`json:"name"`)

	info := tagaudit.StructInfo{
		StructName: "User",
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewVar(token.NoPos, nil, "Name", types.Typ[types.String]), Tags: tags1, RawTag: `json:"name"`},
			{Field: types.NewVar(token.NoPos, nil, "Email", types.Typ[types.String]), Tags: nil, RawTag: ""},
		},
	}

	findings := r.CheckStruct(info, &tagaudit.Config{RequiredTagKeys: []string{"json"}})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].FieldName != "Email" {
		t.Errorf("expected finding on Email, got %s", findings[0].FieldName)
	}
}

func TestCompletenessRule_NoFieldsHaveKey(t *testing.T) {
	r := &CompletenessRule{}

	info := tagaudit.StructInfo{
		StructName: "User",
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewVar(token.NoPos, nil, "Name", types.Typ[types.String]), Tags: nil, RawTag: ""},
			{Field: types.NewVar(token.NoPos, nil, "Email", types.Typ[types.String]), Tags: nil, RawTag: ""},
		},
	}

	findings := r.CheckStruct(info, &tagaudit.Config{RequiredTagKeys: []string{"json"}})
	if len(findings) != 0 {
		t.Errorf("expected 0 findings when no field has the key, got %d", len(findings))
	}
}

func TestCompletenessRule_SkipsUnexported(t *testing.T) {
	r := &CompletenessRule{}
	tags1, _ := parseTag(`json:"name"`)

	pkg := types.NewPackage("example.com/test", "test")
	info := tagaudit.StructInfo{
		StructName: "User",
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewVar(token.NoPos, pkg, "Name", types.Typ[types.String]), Tags: tags1, RawTag: `json:"name"`},
			{Field: types.NewVar(token.NoPos, pkg, "secret", types.Typ[types.String]), Tags: nil, RawTag: ""},
		},
	}

	findings := r.CheckStruct(info, &tagaudit.Config{RequiredTagKeys: []string{"json"}})
	if len(findings) != 0 {
		t.Errorf("expected 0 findings (unexported field should be skipped), got %d", len(findings))
	}
}

func TestCompletenessRule_SkipsAnonymous(t *testing.T) {
	r := &CompletenessRule{}
	tags1, _ := parseTag(`json:"name"`)

	embeddedType := types.NewNamed(
		types.NewTypeName(token.NoPos, nil, "Base", nil),
		types.NewStruct(nil, nil),
		nil,
	)

	info := tagaudit.StructInfo{
		StructName: "User",
		Fields: []tagaudit.FieldInfo{
			{Field: types.NewVar(token.NoPos, nil, "Name", types.Typ[types.String]), Tags: tags1, RawTag: `json:"name"`},
			{Field: types.NewField(token.NoPos, nil, "Base", embeddedType, true), Tags: nil, RawTag: ""},
		},
	}

	findings := r.CheckStruct(info, &tagaudit.Config{RequiredTagKeys: []string{"json"}})
	if len(findings) != 0 {
		t.Errorf("expected 0 findings (anonymous field should be skipped), got %d", len(findings))
	}
}
