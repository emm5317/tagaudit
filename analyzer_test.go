package tagaudit_test

import (
	"testing"

	"github.com/emm5317/tagaudit"
	_ "github.com/emm5317/tagaudit/rules" // register built-in rules
)

func TestAnalyzePackages_BasicSyntax(t *testing.T) {
	a := tagaudit.New(&tagaudit.Config{})

	findings, err := a.AnalyzePackages("./testdata/basic")
	if err != nil {
		t.Fatalf("AnalyzePackages failed: %v", err)
	}

	var syntaxFindings []tagaudit.Finding
	for _, f := range findings {
		if f.RuleID == "syntax" {
			syntaxFindings = append(syntaxFindings, f)
		}
	}

	if len(syntaxFindings) == 0 {
		t.Error("expected at least one syntax finding from testdata/basic, got 0")
	}

	for _, f := range syntaxFindings {
		if f.Severity != tagaudit.SeverityError {
			t.Errorf("syntax findings should be SeverityError, got %v", f.Severity)
		}
		t.Logf("finding: %s", f)
	}
}

func TestAnalyzePackages_Naming(t *testing.T) {
	a := tagaudit.New(&tagaudit.Config{
		NamingConventions: map[string]string{"json": "snake_case"},
	})

	findings, err := a.AnalyzePackages("./testdata/naming")
	if err != nil {
		t.Fatalf("AnalyzePackages failed: %v", err)
	}

	var namingFindings []tagaudit.Finding
	for _, f := range findings {
		if f.RuleID == "naming" {
			namingFindings = append(namingFindings, f)
		}
	}

	// "userName" and "createdAt" should be flagged
	if len(namingFindings) < 2 {
		t.Errorf("expected at least 2 naming findings, got %d", len(namingFindings))
	}

	for _, f := range namingFindings {
		t.Logf("naming finding: %s", f)
		if f.Fix == nil {
			t.Errorf("naming findings should have a fix suggestion")
		}
	}
}

func TestAnalyzePackages_Unexported(t *testing.T) {
	a := tagaudit.New(&tagaudit.Config{})

	findings, err := a.AnalyzePackages("./testdata/basic")
	if err != nil {
		t.Fatalf("AnalyzePackages failed: %v", err)
	}

	var unexportedFindings []tagaudit.Finding
	for _, f := range findings {
		if f.RuleID == "unexported" {
			unexportedFindings = append(unexportedFindings, f)
		}
	}

	// "email" and "secret" in WithUnexported should be flagged
	if len(unexportedFindings) < 2 {
		t.Errorf("expected at least 2 unexported findings, got %d", len(unexportedFindings))
	}

	for _, f := range unexportedFindings {
		t.Logf("unexported finding: %s", f)
	}
}

func TestAnalyzePackages_UnknownKeys(t *testing.T) {
	a := tagaudit.New(&tagaudit.Config{
		KnownTagKeys: []string{"json"},
	})

	findings, err := a.AnalyzePackages("./testdata/basic")
	if err != nil {
		t.Fatalf("AnalyzePackages failed: %v", err)
	}

	var unknownFindings []tagaudit.Finding
	for _, f := range findings {
		if f.RuleID == "unknownkeys" {
			unknownFindings = append(unknownFindings, f)
		}
	}

	// GoodStruct has a "db" tag which is not in known keys
	if len(unknownFindings) == 0 {
		t.Error("expected at least one unknownkeys finding for 'db' tag, got 0")
	}

	for _, f := range unknownFindings {
		t.Logf("unknownkeys finding: %s", f)
	}
}

func TestAnalyzePackages_Completeness(t *testing.T) {
	a := tagaudit.New(&tagaudit.Config{
		RequiredTagKeys: []string{"json"},
	})

	findings, err := a.AnalyzePackages("./testdata/embedded")
	if err != nil {
		t.Fatalf("AnalyzePackages failed: %v", err)
	}

	var completenessFindings []tagaudit.Finding
	for _, f := range findings {
		if f.RuleID == "completeness" {
			completenessFindings = append(completenessFindings, f)
		}
	}

	// Incomplete struct has Age and Phone without json tags
	if len(completenessFindings) < 2 {
		t.Errorf("expected at least 2 completeness findings, got %d", len(completenessFindings))
	}

	for _, f := range completenessFindings {
		t.Logf("completeness finding: %s", f)
	}
}

func TestAnalyzePackages_Duplicates(t *testing.T) {
	a := tagaudit.New(&tagaudit.Config{})

	findings, err := a.AnalyzePackages("./testdata/embedded")
	if err != nil {
		t.Fatalf("AnalyzePackages failed: %v", err)
	}

	var dupFindings []tagaudit.Finding
	for _, f := range findings {
		if f.RuleID == "duplicates" {
			dupFindings = append(dupFindings, f)
		}
	}

	// DuplicateTag.Name duplicates Base.Name's json:"name"
	if len(dupFindings) == 0 {
		t.Error("expected at least one duplicate finding, got 0")
	}

	for _, f := range dupFindings {
		t.Logf("duplicate finding: %s", f)
	}
}

func TestAnalyzePackages_Shadow(t *testing.T) {
	a := tagaudit.New(&tagaudit.Config{})

	findings, err := a.AnalyzePackages("./testdata/embedded")
	if err != nil {
		t.Fatalf("AnalyzePackages failed: %v", err)
	}

	var shadowFindings []tagaudit.Finding
	for _, f := range findings {
		if f.RuleID == "shadow" {
			shadowFindings = append(shadowFindings, f)
		}
	}

	// ShadowedTag.DisplayName shadows Base.Name's json:"name"
	if len(shadowFindings) == 0 {
		t.Error("expected at least one shadow finding, got 0")
	}

	for _, f := range shadowFindings {
		t.Logf("shadow finding: %s", f)
	}
}

func TestAnalyzePackages_NoDuplicatesOnClean(t *testing.T) {
	a := tagaudit.New(&tagaudit.Config{})

	findings, err := a.AnalyzePackages("./testdata/embedded")
	if err != nil {
		t.Fatalf("AnalyzePackages failed: %v", err)
	}

	// NoDuplicates struct should not produce duplicate findings
	for _, f := range findings {
		if f.RuleID == "duplicates" && f.FieldName == "Email" {
			t.Errorf("unexpected duplicate finding on NoDuplicates.Email: %s", f)
		}
	}
}

func TestNew_NilConfig(t *testing.T) {
	a := tagaudit.New(nil)
	if a == nil {
		t.Fatal("New(nil) returned nil")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := tagaudit.DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}
	if cfg.NamingConventions["json"] != "snake_case" {
		t.Errorf("expected json naming convention 'snake_case', got %q", cfg.NamingConventions["json"])
	}
}
