package rules

import (
	"testing"

	"github.com/emm5317/tagaudit"
)

func TestNamingRule_SnakeCaseViolation(t *testing.T) {
	rule := &NamingRule{}
	tags, err := parseTag(`json:"userName"`)
	if err != nil {
		t.Fatal(err)
	}

	cfg := &tagaudit.Config{
		NamingConventions: map[string]string{"json": "snake_case"},
	}

	info := tagaudit.FieldInfo{
		RawTag: `json:"userName"`,
		Tags:   tags,
		Field:  fakeVar("UserName"),
	}

	findings := rule.CheckField(info, cfg)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Fix == nil {
		t.Fatal("expected a fix suggestion")
	}
	if findings[0].Fix.NewTagValue != "user_name" {
		t.Errorf("expected fix value 'user_name', got %q", findings[0].Fix.NewTagValue)
	}
}

func TestNamingRule_SnakeCaseOK(t *testing.T) {
	rule := &NamingRule{}
	tags, err := parseTag(`json:"user_name"`)
	if err != nil {
		t.Fatal(err)
	}

	cfg := &tagaudit.Config{
		NamingConventions: map[string]string{"json": "snake_case"},
	}

	info := tagaudit.FieldInfo{
		RawTag: `json:"user_name"`,
		Tags:   tags,
	}

	findings := rule.CheckField(info, cfg)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for valid snake_case, got %d", len(findings))
	}
}

func TestNamingRule_SkipDash(t *testing.T) {
	rule := &NamingRule{}
	tags, err := parseTag(`json:"-"`)
	if err != nil {
		t.Fatal(err)
	}

	cfg := &tagaudit.Config{
		NamingConventions: map[string]string{"json": "snake_case"},
	}

	info := tagaudit.FieldInfo{
		RawTag: `json:"-"`,
		Tags:   tags,
	}

	findings := rule.CheckField(info, cfg)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for skipped field, got %d", len(findings))
	}
}

func TestNamingRule_NilConfig(t *testing.T) {
	rule := &NamingRule{}
	tags, _ := parseTag(`json:"foo"`)
	info := tagaudit.FieldInfo{Tags: tags, RawTag: `json:"foo"`}

	findings := rule.CheckField(info, nil)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for nil config, got %d", len(findings))
	}
}
