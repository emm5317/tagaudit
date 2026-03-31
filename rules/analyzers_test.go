package rules

import (
	"testing"

	"github.com/emm5317/tagaudit"
)

func TestAnalyzers(t *testing.T) {
	cfg := &tagaudit.Config{
		Rules: All(),
	}
	analyzers := Analyzers(cfg)
	if len(analyzers) != len(All()) {
		t.Errorf("expected %d analyzers, got %d", len(All()), len(analyzers))
	}

	for _, a := range analyzers {
		if a.Name == "" {
			t.Error("analyzer has empty name")
		}
		if a.Doc == "" {
			t.Error("analyzer has empty doc")
		}
		if len(a.Requires) == 0 {
			t.Errorf("analyzer %s should require inspect.Analyzer", a.Name)
		}
	}
}
