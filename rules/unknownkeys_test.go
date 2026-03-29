package rules

import (
	"testing"

	"github.com/emm5317/tagaudit"
)

func TestUnknownKeysRule_KnownKey(t *testing.T) {
	rule := &UnknownKeysRule{}
	tags, err := parseTag(`json:"name"`)
	if err != nil {
		t.Fatal(err)
	}

	cfg := &tagaudit.Config{
		KnownTagKeys: []string{"json", "db"},
	}
	info := tagaudit.FieldInfo{Tags: tags, RawTag: `json:"name"`, Field: fakeVar("Name")}

	findings := rule.CheckField(info, cfg)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for known key, got %d", len(findings))
	}
}

func TestUnknownKeysRule_UnknownKey(t *testing.T) {
	rule := &UnknownKeysRule{}
	tags, err := parseTag(`josn:"name"`)
	if err != nil {
		t.Fatal(err)
	}

	cfg := &tagaudit.Config{
		KnownTagKeys: []string{"json", "db"},
	}
	info := tagaudit.FieldInfo{Tags: tags, RawTag: `josn:"name"`, Field: fakeVar("Name")}

	findings := rule.CheckField(info, cfg)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for unknown key, got %d", len(findings))
	}
	if findings[0].TagKey != "josn" {
		t.Errorf("expected tag key 'josn', got %q", findings[0].TagKey)
	}
}

func TestUnknownKeysRule_Disabled(t *testing.T) {
	rule := &UnknownKeysRule{}
	tags, _ := parseTag(`josn:"name"`)

	cfg := &tagaudit.Config{} // KnownTagKeys is nil -> disabled
	info := tagaudit.FieldInfo{Tags: tags, RawTag: `josn:"name"`}

	findings := rule.CheckField(info, cfg)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings when disabled, got %d", len(findings))
	}
}
