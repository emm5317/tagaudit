package tagaudit_test

import (
	"testing"

	"github.com/emm5317/tagaudit"
	"github.com/emm5317/tagaudit/rules"
)

func TestAnalyzePackages_BasicSyntax(t *testing.T) {
	a := tagaudit.New(&tagaudit.Config{Rules: rules.All()})

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
		Rules:             rules.All(),
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
	a := tagaudit.New(&tagaudit.Config{Rules: rules.All()})

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
		Rules:        rules.All(),
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
		Rules:           rules.All(),
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
	a := tagaudit.New(&tagaudit.Config{Rules: rules.All()})

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
	a := tagaudit.New(&tagaudit.Config{Rules: rules.All()})

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
	a := tagaudit.New(&tagaudit.Config{Rules: rules.All()})

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

func TestAnalyzePackages_TypeAlias(t *testing.T) {
	a := tagaudit.New(&tagaudit.Config{
		Rules:             rules.All(),
		NamingConventions: map[string]string{"json": "snake_case"},
	})

	findings, err := a.AnalyzePackages("./testdata/alias")
	if err != nil {
		t.Fatalf("AnalyzePackages failed: %v", err)
	}

	// AliasedStruct = Base should NOT be analyzed separately (it's an alias,
	// not a new named type). But UsesAlias should be analyzed and userName
	// should trigger the naming rule.
	var namingFindings []tagaudit.Finding
	for _, f := range findings {
		if f.RuleID == "naming" {
			namingFindings = append(namingFindings, f)
		}
	}

	if len(namingFindings) == 0 {
		t.Error("expected at least one naming finding for userName in alias testdata, got 0")
	}

	for _, f := range namingFindings {
		t.Logf("alias naming finding: %s", f)
	}
}

func TestAnalyzePackages_AnonymousStructs(t *testing.T) {
	a := tagaudit.New(&tagaudit.Config{
		Rules:             rules.All(),
		NamingConventions: map[string]string{"json": "snake_case"},
	})

	findings, err := a.AnalyzePackages("./testdata/anon")
	if err != nil {
		t.Fatalf("AnalyzePackages failed: %v", err)
	}

	var namingFindings, syntaxFindings []tagaudit.Finding
	for _, f := range findings {
		switch f.RuleID {
		case "naming":
			namingFindings = append(namingFindings, f)
		case "syntax":
			syntaxFindings = append(syntaxFindings, f)
		}
	}

	if len(namingFindings) == 0 {
		t.Error("expected naming finding for userName in anonymous composite literal, got 0")
	}
	if len(syntaxFindings) == 0 {
		t.Error("expected syntax finding for bad tag in anonymous var declaration, got 0")
	}

	for _, f := range findings {
		t.Logf("anon finding: %s", f)
	}
}

func TestAnalyzePackages_InvalidPattern(t *testing.T) {
	a := tagaudit.New(&tagaudit.Config{Rules: rules.All()})

	_, err := a.AnalyzePackages("./nonexistent_package_xyz_123")
	if err == nil {
		t.Error("expected error for invalid package pattern, got nil")
	}
}

func TestAnalyzePackages_NonStructTypes(t *testing.T) {
	// testdata/basic now contains MyString, MyInterface, MyAlias — these
	// should be silently skipped without errors
	a := tagaudit.New(&tagaudit.Config{Rules: rules.All()})

	_, err := a.AnalyzePackages("./testdata/basic")
	if err != nil {
		t.Fatalf("AnalyzePackages should handle non-struct types gracefully: %v", err)
	}
}

func TestSeverityAllowed_ZeroValue(t *testing.T) {
	// A default Config{} (MinSeverity zero value = SeverityInfo) should include all findings.
	a := tagaudit.New(&tagaudit.Config{Rules: rules.All()})

	findings, err := a.AnalyzePackages("./testdata/basic")
	if err != nil {
		t.Fatalf("AnalyzePackages failed: %v", err)
	}

	hasError, hasInfo := false, false
	for _, f := range findings {
		if f.Severity == tagaudit.SeverityError {
			hasError = true
		}
		if f.Severity == tagaudit.SeverityInfo {
			hasInfo = true
		}
	}

	if !hasError {
		t.Error("default config should include error-level findings")
	}
	// Info-level findings require naming config; test with naming testdata
	a2 := tagaudit.New(&tagaudit.Config{
		Rules:             rules.All(),
		NamingConventions: map[string]string{"json": "snake_case"},
	})
	findings2, err := a2.AnalyzePackages("./testdata/naming")
	if err != nil {
		t.Fatalf("AnalyzePackages failed: %v", err)
	}
	for _, f := range findings2 {
		if f.Severity == tagaudit.SeverityInfo {
			hasInfo = true
		}
	}
	if !hasInfo {
		t.Error("default config should include info-level findings")
	}
}

func TestSeverityAllowed_ErrorsOnly(t *testing.T) {
	a := tagaudit.New(&tagaudit.Config{
		Rules:             rules.All(),
		NamingConventions: map[string]string{"json": "snake_case"},
		MinSeverity:       tagaudit.SeverityError,
	})

	findings, err := a.AnalyzePackages("./testdata/basic")
	if err != nil {
		t.Fatalf("AnalyzePackages failed: %v", err)
	}

	for _, f := range findings {
		if f.Severity != tagaudit.SeverityError {
			t.Errorf("with MinSeverity=SeverityError, got finding with severity %v: %s", f.Severity, f)
		}
	}
}

func TestSeverity_String(t *testing.T) {
	tests := []struct {
		s    tagaudit.Severity
		want string
	}{
		{tagaudit.SeverityError, "error"},
		{tagaudit.SeverityWarning, "warning"},
		{tagaudit.SeverityInfo, "info"},
		{tagaudit.Severity(99), "unknown"},
	}
	for _, tt := range tests {
		got := tt.s.String()
		if got != tt.want {
			t.Errorf("Severity(%d).String() = %q, want %q", tt.s, got, tt.want)
		}
	}
}

func TestFinding_String(t *testing.T) {
	f := tagaudit.Finding{
		RuleID:   "test",
		Severity: tagaudit.SeverityError,
		Message:  "test message",
	}
	s := f.String()
	if s == "" {
		t.Error("Finding.String() returned empty string")
	}
}

func TestNew_NilConfig(t *testing.T) {
	a := tagaudit.New(nil)
	if a == nil {
		t.Fatal("New(nil) returned nil")
	}
}

func TestNew_EmptyRulesNoFindings(t *testing.T) {
	// Empty rules list means no rules run — no findings.
	a := tagaudit.New(&tagaudit.Config{})
	if a == nil {
		t.Fatal("New returned nil")
	}
	findings, err := a.AnalyzePackages("./testdata/basic")
	if err != nil {
		t.Fatalf("AnalyzePackages failed: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings with no rules, got %d", len(findings))
	}
}

func TestRulesDefaultConfig(t *testing.T) {
	// rules.DefaultConfig provides a fully populated config with all rules.
	cfg := rules.DefaultConfig()
	a := tagaudit.New(cfg)
	findings, err := a.AnalyzePackages("./testdata/basic")
	if err != nil {
		t.Fatalf("AnalyzePackages failed: %v", err)
	}
	hasSyntax := false
	for _, f := range findings {
		if f.RuleID == "syntax" {
			hasSyntax = true
		}
	}
	if !hasSyntax {
		t.Error("expected rules.DefaultConfig to catch syntax errors")
	}
}

func TestDefaultConfig(t *testing.T) {
	// DefaultConfig returns sensible defaults but no rules.
	// Use rules.DefaultConfig() to get a config with rules.
	cfg := tagaudit.DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}
	if cfg.NamingConventions["json"] != "snake_case" {
		t.Errorf("expected json naming convention 'snake_case', got %q", cfg.NamingConventions["json"])
	}
	if len(cfg.Rules) != 0 {
		t.Error("DefaultConfig should not include rules; use rules.DefaultConfig() instead")
	}

	// rules.DefaultConfig() should include built-in rules.
	rcfg := rules.DefaultConfig()
	if len(rcfg.Rules) == 0 {
		t.Error("rules.DefaultConfig should include built-in rules")
	}
}

func TestBaseConfig(t *testing.T) {
	cfg := tagaudit.BaseConfig()
	if cfg == nil {
		t.Fatal("BaseConfig returned nil")
	}
	if cfg.NamingConventions["json"] != "snake_case" {
		t.Errorf("expected json naming convention 'snake_case', got %q", cfg.NamingConventions["json"])
	}
	if len(cfg.Rules) != 0 {
		t.Error("BaseConfig should return a skeleton config with no rules")
	}
}
