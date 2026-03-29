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
