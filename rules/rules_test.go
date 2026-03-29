package rules

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/emm5317/tagaudit"
)

func TestAll_ReturnsAllRules(t *testing.T) {
	all := All()
	if len(all) == 0 {
		t.Fatal("All() returned empty slice")
	}

	// Verify each rule implements at least one interface
	for _, r := range all {
		_, isField := r.(tagaudit.FieldChecker)
		_, isStruct := r.(tagaudit.StructChecker)
		if !isField && !isStruct {
			t.Errorf("rule %T implements neither FieldChecker nor StructChecker", r)
		}
	}
}

func TestAllDescriptions(t *testing.T) {
	// Every rule must have a non-empty ID and Description
	all := All()
	for _, r := range all {
		if fc, ok := r.(tagaudit.FieldChecker); ok {
			if fc.ID() == "" {
				t.Errorf("%T has empty ID", r)
			}
			if fc.Description() == "" {
				t.Errorf("%T has empty Description", r)
			}
		}
		if sc, ok := r.(tagaudit.StructChecker); ok {
			if sc.ID() == "" {
				t.Errorf("%T has empty ID", r)
			}
			if sc.Description() == "" {
				t.Errorf("%T has empty Description", r)
			}
		}
	}
}

func TestPosFromInfo_WithTagPos(t *testing.T) {
	fset := token.NewFileSet()
	f := fset.AddFile("test.go", -1, 100)
	f.AddLine(0)
	f.AddLine(20)
	f.AddLine(40)

	astField := &ast.Field{
		Tag: &ast.BasicLit{
			ValuePos: f.Pos(25),
			Value:    "`json:\"name\"`",
		},
	}

	info := tagaudit.FieldInfo{
		Fset:     fset,
		ASTField: astField,
	}

	pos := posFromInfo(info)
	if pos.Line != 2 {
		t.Errorf("expected line 2, got %d", pos.Line)
	}
}

func TestPosFromInfo_WithoutTag(t *testing.T) {
	fset := token.NewFileSet()
	f := fset.AddFile("test.go", -1, 100)
	f.AddLine(0)
	f.AddLine(20)

	astField := &ast.Field{}
	// Field.Pos() returns the position of the first token — for a field without names it's 0
	// We just verify it doesn't panic
	info := tagaudit.FieldInfo{
		Fset:     fset,
		ASTField: astField,
	}

	pos := posFromInfo(info)
	// Should return a valid position (from field pos, not tag pos)
	_ = pos
}

func TestPosFromInfo_NilFset(t *testing.T) {
	info := tagaudit.FieldInfo{
		Fset:     nil,
		ASTField: &ast.Field{},
	}

	pos := posFromInfo(info)
	if pos.IsValid() {
		t.Error("expected invalid position for nil Fset")
	}
}

func TestPosFromInfo_NilASTField(t *testing.T) {
	fset := token.NewFileSet()
	info := tagaudit.FieldInfo{
		Fset:     fset,
		ASTField: nil,
	}

	pos := posFromInfo(info)
	if pos.IsValid() {
		t.Error("expected invalid position for nil ASTField")
	}
}
