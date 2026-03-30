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
			return runAnalysis(a, pass)
		},
	}
}

// NewSingleRuleAnalyzer returns an analysis.Analyzer that runs exactly one rule.
// The analyzer is named "tagaudit_<ruleID>" so each rule can be individually
// enabled/disabled in tools like golangci-lint or multichecker.
func NewSingleRuleAnalyzer(rule Rule, cfg *Config) *analysis.Analyzer {
	singleCfg := *cfg
	singleCfg.Rules = []Rule{rule}
	a := New(&singleCfg)

	return &analysis.Analyzer{
		Name:     "tagaudit_" + rule.ID(),
		Doc:      rule.Description(),
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run: func(pass *analysis.Pass) (any, error) {
			return runAnalysis(a, pass)
		},
	}
}

// runAnalysis is the shared implementation for NewAnalyzer and NewSingleRuleAnalyzer.
func runAnalysis(a *Analyzer, pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	seen := make(map[*ast.StructType]bool)

	nodeFilter := []ast.Node{
		(*ast.TypeSpec)(nil),
		(*ast.CompositeLit)(nil),
		(*ast.ValueSpec)(nil),
	}
	insp.Preorder(nodeFilter, func(n ast.Node) {
		var name string
		var st *ast.StructType
		var styp *types.Struct

		switch node := n.(type) {
		case *ast.TypeSpec:
			var ok bool
			name, st, styp, ok = resolveStruct(node, pass.TypesInfo)
			if !ok {
				return
			}
		case *ast.CompositeLit:
			var ok bool
			st, ok = node.Type.(*ast.StructType)
			if !ok {
				return
			}
			tv, exists := pass.TypesInfo.Types[node.Type]
			if !exists {
				return
			}
			styp, ok = tv.Type.(*types.Struct)
			if !ok {
				return
			}
		case *ast.ValueSpec:
			var ok bool
			st, ok = node.Type.(*ast.StructType)
			if !ok {
				return
			}
			tv, exists := pass.TypesInfo.Types[node.Type]
			if !exists {
				return
			}
			styp, ok = tv.Type.(*types.Struct)
			if !ok {
				return
			}
		default:
			return
		}

		if seen[st] {
			return
		}
		seen[st] = true

		fields := buildFieldInfos(pass.Fset, name, styp, st)

		// Run field-level checks
		for _, field := range fields {
			for _, fc := range a.fieldCheckers {
				for _, finding := range fc.CheckField(field, a.cfg) {
					if !a.severityAllowed(finding.Severity) {
						continue
					}
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
			StructName: name,
			StructType: styp,
			ASTNode:    st,
			Fields:     fields,
		}
		for _, sc := range a.structCheckers {
			for _, finding := range sc.CheckStruct(si, a.cfg) {
				if !a.severityAllowed(finding.Severity) {
					continue
				}
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
