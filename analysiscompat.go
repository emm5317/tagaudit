package tagaudit

import (
	"go/ast"
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
							if field.ASTField.Tag != nil {
								pos = field.ASTField.Tag.Pos()
							}
							pass.Reportf(pos, "[%s] %s", finding.RuleID, finding.Message)
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
						// Use the struct's position as fallback
						pos := st.Pos()
						// Try to find the specific field's AST position
						for _, f := range fields {
							if f.Field != nil && f.Field.Name() == finding.FieldName {
								pos = f.ASTField.Pos()
								if f.ASTField.Tag != nil {
									pos = f.ASTField.Tag.Pos()
								}
								break
							}
						}
						pass.Reportf(pos, "[%s] %s", finding.RuleID, finding.Message)
					}
				}
			})

			return nil, nil
		},
	}
}
