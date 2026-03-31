package tagaudit_test

import (
	"testing"

	"github.com/emm5317/tagaudit"
	"github.com/emm5317/tagaudit/rules"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestNewAnalyzer(t *testing.T) {
	cfg := &tagaudit.Config{
		Rules:             rules.All(),
		NamingConventions: map[string]string{"json": "snake_case"},
		RequiredTagKeys:   []string{"json"},
	}
	analyzer := tagaudit.NewAnalyzer(cfg)
	if analyzer == nil {
		t.Fatal("NewAnalyzer returned nil")
	}
	if analyzer.Name != "tagaudit" {
		t.Errorf("expected analyzer name 'tagaudit', got %q", analyzer.Name)
	}

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer, "analysistest/basic")
}

func TestNewSingleRuleAnalyzer(t *testing.T) {
	cfg := &tagaudit.Config{
		NamingConventions: map[string]string{"json": "snake_case"},
		RequiredTagKeys:   []string{"json"},
	}

	// Create a syntax-only analyzer
	syntaxAnalyzer := tagaudit.NewSingleRuleAnalyzer(&rules.SyntaxRule{}, cfg)
	if syntaxAnalyzer == nil {
		t.Fatal("NewSingleRuleAnalyzer returned nil")
	}
	if syntaxAnalyzer.Name != "tagaudit_syntax" {
		t.Errorf("expected analyzer name 'tagaudit_syntax', got %q", syntaxAnalyzer.Name)
	}
}

func TestAnalyzers_Count(t *testing.T) {
	cfg := &tagaudit.Config{}
	analyzers := rules.Analyzers(cfg)
	if len(analyzers) != 8 {
		t.Errorf("expected 8 analyzers, got %d", len(analyzers))
	}

	names := make(map[string]bool)
	for _, a := range analyzers {
		if names[a.Name] {
			t.Errorf("duplicate analyzer name: %s", a.Name)
		}
		names[a.Name] = true
	}
}

func TestNewAnalyzer_AnonymousStructs(t *testing.T) {
	cfg := &tagaudit.Config{
		Rules:             rules.All(),
		NamingConventions: map[string]string{"json": "snake_case"},
	}
	analyzer := tagaudit.NewAnalyzer(cfg)
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer, "analysistest/anon")
}

func TestNewAnalyzer_Fixes(t *testing.T) {
	cfg := &tagaudit.Config{
		Rules:             rules.All(),
		NamingConventions: map[string]string{"json": "snake_case"},
		RequiredTagKeys:   []string{"json"},
	}
	analyzer := tagaudit.NewAnalyzer(cfg)
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer, "analysistest/fixes")
}

func TestNewAnalyzer_SeverityFiltering(t *testing.T) {
	cfg := &tagaudit.Config{
		Rules:             rules.All(),
		NamingConventions: map[string]string{"json": "snake_case"},
		RequiredTagKeys:   []string{"json"},
		MinSeverity:       tagaudit.SeverityError,
	}
	analyzer := tagaudit.NewAnalyzer(cfg)

	// With MinSeverity=Error, naming (warning) and completeness (warning)
	// should be filtered out; only syntax (error) should remain.
	testdata := analysistest.TestData()
	// Run against basic which has syntax errors (error level)
	analysistest.Run(t, testdata, analyzer, "analysistest/severityfilter")
}

func TestNewSingleRuleAnalyzer_Run(t *testing.T) {
	cfg := &tagaudit.Config{
		NamingConventions: map[string]string{"json": "snake_case"},
	}
	syntaxAnalyzer := tagaudit.NewSingleRuleAnalyzer(&rules.SyntaxRule{}, cfg)
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, syntaxAnalyzer, "analysistest/severityfilter")
}
