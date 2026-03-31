package rules

import (
	"testing"
)

func TestJSONPreset(t *testing.T) {
	cfg := JSONPreset()
	if cfg == nil {
		t.Fatal("JSONPreset returned nil")
	}
	if len(cfg.Rules) == 0 {
		t.Error("JSONPreset should include all rules")
	}
	if cfg.NamingConventions["json"] != "snake_case" {
		t.Errorf("expected json naming convention 'snake_case', got %q", cfg.NamingConventions["json"])
	}
	if len(cfg.RequiredTagKeys) == 0 || cfg.RequiredTagKeys[0] != "json" {
		t.Error("JSONPreset should require json tag key")
	}
	if len(cfg.KnownTagKeys) == 0 {
		t.Error("JSONPreset should have known tag keys")
	}
}

func TestAPIModelPreset(t *testing.T) {
	cfg := APIModelPreset()
	if cfg == nil {
		t.Fatal("APIModelPreset returned nil")
	}
	if len(cfg.Rules) == 0 {
		t.Error("APIModelPreset should include all rules")
	}
	if cfg.NamingConventions["json"] != "snake_case" {
		t.Errorf("expected json naming convention 'snake_case', got %q", cfg.NamingConventions["json"])
	}
	if cfg.NamingConventions["yaml"] != "snake_case" {
		t.Errorf("expected yaml naming convention 'snake_case', got %q", cfg.NamingConventions["yaml"])
	}
	if len(cfg.KnownOptions) == 0 {
		t.Error("APIModelPreset should have known options for validate")
	}
	if _, ok := cfg.KnownOptions["validate"]; !ok {
		t.Error("APIModelPreset should have validate in known options")
	}
}

func TestGORMPreset(t *testing.T) {
	cfg := GORMPreset()
	if cfg == nil {
		t.Fatal("GORMPreset returned nil")
	}
	if len(cfg.Rules) == 0 {
		t.Error("GORMPreset should include all rules")
	}
	if cfg.NamingConventions["json"] != "snake_case" {
		t.Errorf("expected json naming convention 'snake_case', got %q", cfg.NamingConventions["json"])
	}
	if _, ok := cfg.KnownOptions["gorm"]; !ok {
		t.Error("GORMPreset should have gorm in known options")
	}
}
