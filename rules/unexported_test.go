package rules

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/emm5317/tagaudit"
)

func TestUnexportedRule_ID(t *testing.T) {
	r := &UnexportedRule{}
	if r.ID() != "unexported" {
		t.Errorf("expected 'unexported', got %q", r.ID())
	}
	if r.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestUnexportedRule_ExportedField(t *testing.T) {
	r := &UnexportedRule{}
	tags, _ := parseTag(`json:"name"`)

	info := tagaudit.FieldInfo{
		Field:  fakeVar("Name"), // uppercase = exported
		Tags:   tags,
		RawTag: `json:"name"`,
	}

	findings := r.CheckField(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for exported field, got %d", len(findings))
	}
}

func TestUnexportedRule_UnexportedWithJsonTag(t *testing.T) {
	r := &UnexportedRule{}
	tags, _ := parseTag(`json:"name"`)

	pkg := types.NewPackage("example.com/test", "test")
	info := tagaudit.FieldInfo{
		Field:  types.NewVar(token.NoPos, pkg, "name", types.Typ[types.String]),
		Tags:   tags,
		RawTag: `json:"name"`,
	}

	findings := r.CheckField(info, nil)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for unexported field with json tag, got %d", len(findings))
	}
	if findings[0].Severity != tagaudit.SeverityWarning {
		t.Errorf("expected SeverityWarning, got %v", findings[0].Severity)
	}
	if findings[0].TagKey != "json" {
		t.Errorf("expected tag key 'json', got %q", findings[0].TagKey)
	}
}

func TestUnexportedRule_UnexportedWithXmlTag(t *testing.T) {
	r := &UnexportedRule{}
	tags, _ := parseTag(`xml:"name"`)

	pkg := types.NewPackage("example.com/test", "test")
	info := tagaudit.FieldInfo{
		Field:  types.NewVar(token.NoPos, pkg, "name", types.Typ[types.String]),
		Tags:   tags,
		RawTag: `xml:"name"`,
	}

	findings := r.CheckField(info, nil)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for unexported field with xml tag, got %d", len(findings))
	}
}

func TestUnexportedRule_UnexportedWithNonEncodingTag(t *testing.T) {
	r := &UnexportedRule{}
	tags, _ := parseTag(`db:"name"`)

	pkg := types.NewPackage("example.com/test", "test")
	info := tagaudit.FieldInfo{
		Field:  types.NewVar(token.NoPos, pkg, "name", types.Typ[types.String]),
		Tags:   tags,
		RawTag: `db:"name"`,
	}

	findings := r.CheckField(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for non-encoding tag on unexported field, got %d", len(findings))
	}
}

func TestUnexportedRule_NilTags(t *testing.T) {
	r := &UnexportedRule{}
	info := tagaudit.FieldInfo{
		Field: fakeVar("name"),
		Tags:  nil,
	}

	findings := r.CheckField(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for nil tags, got %d", len(findings))
	}
}

func TestUnexportedRule_NilField(t *testing.T) {
	r := &UnexportedRule{}
	tags, _ := parseTag(`json:"name"`)
	info := tagaudit.FieldInfo{
		Field:  nil,
		Tags:   tags,
		RawTag: `json:"name"`,
	}

	findings := r.CheckField(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for nil field, got %d", len(findings))
	}
}

func TestUnexportedRule_MultipleEncodingTags(t *testing.T) {
	r := &UnexportedRule{}
	tags, _ := parseTag(`json:"name" xml:"name" yaml:"name"`)

	pkg := types.NewPackage("example.com/test", "test")
	info := tagaudit.FieldInfo{
		Field:  types.NewVar(token.NoPos, pkg, "secret", types.Typ[types.String]),
		Tags:   tags,
		RawTag: `json:"name" xml:"name" yaml:"name"`,
	}

	findings := r.CheckField(info, nil)
	if len(findings) != 3 {
		t.Errorf("expected 3 findings (one per encoding tag), got %d", len(findings))
	}
}
