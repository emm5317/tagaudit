package tagaudit

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NewAnalyzer returns a go/analysis.Analyzer that can be used with
// golangci-lint, multichecker, or singlechecker.
func NewAnalyzer(cfg *Config) *analysis.Analyzer {
	a := New(cfg)

	return &analysis.Analyzer{
		Name:     "tagaudit",
		Doc:      "comprehensive struct tag validation",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run: func(pass *analysis.Pass) (any, error) {
			insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

			nodeFilter := []ast.Node{(*ast.TypeSpec)(nil)}
			insp.Preorder(nodeFilter, func(n ast.Node) {
				ts := n.(*ast.TypeSpec)
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					return
				}

				obj := pass.TypesInfo.Defs[ts.Name]
				if obj == nil {
					return
				}
				named, ok := obj.Type().(*types.Named)
				if !ok {
					return
				}
				styp, ok := named.Underlying().(*types.Struct)
				if !ok {
					return
				}

				fields := buildFieldInfos(pass.Fset, ts.Name.Name, styp, st)

				// Run field-level checks
				for _, field := range fields {
					for _, fc := range a.fieldCheckers {
						for _, finding := range fc.CheckField(field, a.cfg) {
							pos := field.ASTField.Pos()
							end := token.NoPos
							if field.ASTField.Tag != nil {
								pos = field.ASTField.Tag.Pos()
								end = field.ASTField.Tag.End()
							}
							pass.Report(buildDiagnostic(pos, end, field.ASTField, finding))
						}
					}
				}

				// Run struct-level checks
				si := StructInfo{
					Fset:       pass.Fset,
					StructName: ts.Name.Name,
					StructType: styp,
					ASTNode:    st,
					Fields:     fields,
				}
				for _, sc := range a.structCheckers {
					for _, finding := range sc.CheckStruct(si, a.cfg) {
						pos := st.Pos()
						end := token.NoPos
						var astField *ast.Field
						for _, f := range fields {
							if f.Field != nil && f.Field.Name() == finding.FieldName {
								astField = f.ASTField
								pos = f.ASTField.Pos()
								if f.ASTField.Tag != nil {
									pos = f.ASTField.Tag.Pos()
									end = f.ASTField.Tag.End()
								}
								break
							}
						}
						pass.Report(buildDiagnostic(pos, end, astField, finding))
					}
				}
			})

			return nil, nil
		},
	}
}

// buildDiagnostic converts a Finding into an analysis.Diagnostic with
// category, end position, and suggested fixes when available.
func buildDiagnostic(pos, end token.Pos, astField *ast.Field, f Finding) analysis.Diagnostic {
	d := analysis.Diagnostic{
		Pos:      pos,
		End:      end,
		Category: f.Severity.String(),
		Message:  fmt.Sprintf("[%s] %s", f.RuleID, f.Message),
	}

	if f.Fix != nil && astField != nil && astField.Tag != nil {
		d.SuggestedFixes = []analysis.SuggestedFix{{
			Message: f.Fix.Description,
			TextEdits: []analysis.TextEdit{{
				Pos:     astField.Tag.Pos(),
				End:     astField.Tag.End(),
				NewText: []byte("`" + f.Fix.NewTagValue + "`"),
			}},
		}}
	}

	return d
}
