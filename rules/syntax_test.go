package rules

import (
	"testing"

	"github.com/emm5317/tagaudit"
)

func TestSyntaxRule_CheckField_NoTag(t *testing.T) {
	rule := &SyntaxRule{}
	info := tagaudit.FieldInfo{
		RawTag: "",
	}
	findings := rule.CheckField(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for empty tag, got %d", len(findings))
	}
}

func TestSyntaxRule_CheckField_ValidTag(t *testing.T) {
	rule := &SyntaxRule{}

	// Simulate a valid tag by setting Tags to non-nil
	// In real usage, Tags is populated by structtag.Parse
	tags, err := parseTag(`json:"name,omitempty"`)
	if err != nil {
		t.Fatal(err)
	}

	info := tagaudit.FieldInfo{
		RawTag: `json:"name,omitempty"`,
		Tags:   tags,
	}
	findings := rule.CheckField(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for valid tag, got %d", len(findings))
	}
}

func TestSyntaxRule_CheckField_MalformedTag(t *testing.T) {
	rule := &SyntaxRule{}
	info := tagaudit.FieldInfo{
		RawTag: `json:bad_no_quotes`,
		Tags:   nil, // parse failed
	}
	// Note: without Fset/ASTField set, we can't get position info,
	// but the rule should still produce a finding.
	// In integration tests, these are populated.
	findings := rule.CheckField(info, nil)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for malformed tag, got %d", len(findings))
	}
	if findings[0].RuleID != "syntax" {
		t.Errorf("expected rule ID 'syntax', got %q", findings[0].RuleID)
	}
	if findings[0].Severity != tagaudit.SeverityError {
		t.Errorf("expected SeverityError, got %v", findings[0].Severity)
	}
}

func TestSyntaxRule_ID(t *testing.T) {
	rule := &SyntaxRule{}
	if rule.ID() != "syntax" {
		t.Errorf("expected ID 'syntax', got %q", rule.ID())
	}
}
