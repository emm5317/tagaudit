// Package tagaudit provides composable struct tag validation for Go.
//
// It offers a pluggable rule system with built-in rules for common issues
// (syntax errors, naming violations, duplicate tags, etc.) and supports
// user-defined custom rules.
package tagaudit

import (
	"fmt"
	"sort"

	"golang.org/x/tools/go/packages"
)

// Analyzer runs struct tag validation rules against Go packages.
type Analyzer struct {
	cfg            *Config
	fieldCheckers  []FieldChecker
	structCheckers []StructChecker
}

// New creates an Analyzer with the given config.
// If cfg is nil, DefaultConfig() is used.
// If cfg.Rules is empty, no rules are applied; use rules.All() or
// rules.DefaultConfig() to include the built-in rule set.
func New(cfg *Config) *Analyzer {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	a := &Analyzer{cfg: cfg}

	for _, r := range cfg.Rules {
		if fc, ok := r.(FieldChecker); ok {
			a.fieldCheckers = append(a.fieldCheckers, fc)
		}
		if sc, ok := r.(StructChecker); ok {
			a.structCheckers = append(a.structCheckers, sc)
		}
	}

	return a
}

// severityAllowed reports whether a finding with the given severity should be
// included under the current MinSeverity configuration.
func (a *Analyzer) severityAllowed(s Severity) bool {
	return s >= a.cfg.MinSeverity
}

// AnalyzePackages loads the named packages and runs all configured rules
// against every struct found in them. Patterns follow the same conventions
// as `go build` (e.g., "./...", "./internal/models").
func (a *Analyzer) AnalyzePackages(patterns ...string) ([]Finding, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax |
			packages.NeedTypes | packages.NeedTypesInfo,
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("loading packages: %w", err)
	}

	// Check for package load errors
	for _, pkg := range pkgs {
		for _, e := range pkg.Errors {
			return nil, fmt.Errorf("package %s: %s", pkg.PkgPath, e.Msg)
		}
	}

	var findings []Finding

	for _, pkg := range pkgs {
		decls := extractStructs(pkg.Syntax, pkg.TypesInfo)
		for _, decl := range decls {
			fields := buildFieldInfos(pkg.Fset, decl.name, decl.styp, decl.astNode)

			// Run field-level checks
			for _, field := range fields {
				for _, fc := range a.fieldCheckers {
					for _, f := range fc.CheckField(field, a.cfg) {
						if a.severityAllowed(f.Severity) {
							findings = append(findings, f)
						}
					}
				}
			}

			// Run struct-level checks
			si := StructInfo{
				Fset:       pkg.Fset,
				StructName: decl.name,
				StructType: decl.styp,
				ASTNode:    decl.astNode,
				Fields:     fields,
			}
			for _, sc := range a.structCheckers {
				for _, f := range sc.CheckStruct(si, a.cfg) {
					if a.severityAllowed(f.Severity) {
						findings = append(findings, f)
					}
				}
			}
		}
	}

	// Sort findings by file position
	sort.Slice(findings, func(i, j int) bool {
		fi, fj := findings[i], findings[j]
		if fi.Pos.Filename != fj.Pos.Filename {
			return fi.Pos.Filename < fj.Pos.Filename
		}
		if fi.Pos.Line != fj.Pos.Line {
			return fi.Pos.Line < fj.Pos.Line
		}
		return fi.Pos.Column < fj.Pos.Column
	})

	return findings, nil
}
