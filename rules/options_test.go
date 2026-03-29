package rules

import (
	"testing"

	"github.com/emm5317/tagaudit"
)

func TestOptionsRule_ValidOptions(t *testing.T) {
	rule := &OptionsRule{}
	tags, err := parseTag(`json:"name,omitempty"`)
	if err != nil {
		t.Fatal(err)
	}

	info := tagaudit.FieldInfo{
		RawTag: `json:"name,omitempty"`,
		Tags:   tags,
		Field:  fakeVar("Name"),
	}

	findings := rule.CheckField(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for valid options, got %d", len(findings))
	}
}

func TestOptionsRule_InvalidOption(t *testing.T) {
	rule := &OptionsRule{}
	tags, err := parseTag(`json:"name,omitemtpy"`)
	if err != nil {
		t.Fatal(err)
	}

	info := tagaudit.FieldInfo{
		RawTag: `json:"name,omitemtpy"`,
		Tags:   tags,
		Field:  fakeVar("Name"),
	}

	findings := rule.CheckField(info, nil)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for invalid option, got %d", len(findings))
	}
	if findings[0].Severity != tagaudit.SeverityWarning {
		t.Errorf("expected SeverityWarning, got %v", findings[0].Severity)
	}
}

func TestOptionsRule_NonWellKnownTag(t *testing.T) {
	rule := &OptionsRule{}
	tags, err := parseTag(`db:"name,readonly"`)
	if err != nil {
		t.Fatal(err)
	}

	info := tagaudit.FieldInfo{
		RawTag: `db:"name,readonly"`,
		Tags:   tags,
	}

	findings := rule.CheckField(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for non-well-known tag, got %d", len(findings))
	}
}
