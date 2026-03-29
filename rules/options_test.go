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

func TestOptionsRule_CustomKnownOptions(t *testing.T) {
	rule := &OptionsRule{}
	tags, err := parseTag(`gorm:"name,primaryKey"`)
	if err != nil {
		t.Fatal(err)
	}

	info := tagaudit.FieldInfo{
		RawTag: `gorm:"name,primaryKey"`,
		Tags:   tags,
		Field:  fakeVar("ID"),
	}

	// Without custom options, gorm is skipped entirely
	findings := rule.CheckField(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings without custom options, got %d", len(findings))
	}

	// With custom options, valid options pass
	cfg := &tagaudit.Config{
		KnownOptions: map[string][]string{"gorm": {"primaryKey", "autoIncrement"}},
	}
	findings = rule.CheckField(info, cfg)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for valid custom option, got %d", len(findings))
	}

	// With custom options, invalid options are flagged
	tags2, _ := parseTag(`gorm:"name,badOption"`)
	info2 := tagaudit.FieldInfo{
		RawTag: `gorm:"name,badOption"`,
		Tags:   tags2,
		Field:  fakeVar("ID"),
	}
	findings = rule.CheckField(info2, cfg)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for invalid custom option, got %d", len(findings))
	}
}

func TestOptionsRule_CustomOptionsMergeWithBuiltin(t *testing.T) {
	rule := &OptionsRule{}
	tags, err := parseTag(`json:"name,omitempty,inline"`)
	if err != nil {
		t.Fatal(err)
	}

	// "inline" is not a built-in json option
	info := tagaudit.FieldInfo{
		RawTag: `json:"name,omitempty,inline"`,
		Tags:   tags,
		Field:  fakeVar("Name"),
	}

	// Without custom options, "inline" is flagged
	findings := rule.CheckField(info, nil)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for unknown json option, got %d", len(findings))
	}

	// With custom options adding "inline", it passes
	cfg := &tagaudit.Config{
		KnownOptions: map[string][]string{"json": {"inline"}},
	}
	findings = rule.CheckField(info, cfg)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings with custom json option, got %d", len(findings))
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
