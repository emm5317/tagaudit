package rules

import (
	"github.com/emm5317/tagaudit"
	"golang.org/x/tools/go/analysis"
)

// Analyzers returns one analysis.Analyzer per built-in rule.
// Each analyzer runs independently and can be used with multichecker
// or selectively enabled in golangci-lint.
func Analyzers(cfg *tagaudit.Config) []*analysis.Analyzer {
	var analyzers []*analysis.Analyzer
	for _, r := range All() {
		analyzers = append(analyzers, tagaudit.NewSingleRuleAnalyzer(r, cfg))
	}
	return analyzers
}
