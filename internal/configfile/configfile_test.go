package configfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/emm5317/tagaudit"
)

func TestLoad(t *testing.T) {
	content := `
rules:
  enable: [syntax, naming]
  disable: [duplicates]
naming_conventions:
  json: snake_case
required_tag_keys: [json]
known_tag_keys: [json, yaml]
known_options:
  gorm: [primaryKey, column]
min_severity: warning
`
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fc, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(fc.Rules.Enable) != 2 {
		t.Errorf("expected 2 enabled rules, got %d", len(fc.Rules.Enable))
	}
	if fc.NamingConventions["json"] != "snake_case" {
		t.Errorf("expected json naming snake_case, got %q", fc.NamingConventions["json"])
	}
	if fc.MinSeverity != "warning" {
		t.Errorf("expected min_severity warning, got %q", fc.MinSeverity)
	}
}

func TestLoad_BadPath(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestLoad_BadYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(path, []byte(":::bad yaml"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for bad YAML")
	}
}

type mockRule struct {
	id string
}

func (r *mockRule) ID() string          { return r.id }
func (r *mockRule) Description() string { return "" }

func TestToConfig_FilterRules(t *testing.T) {
	allRules := []tagaudit.Rule{
		&mockRule{"syntax"},
		&mockRule{"naming"},
		&mockRule{"duplicates"},
	}

	fc := &FileConfig{
		Rules: RulesConfig{
			Enable: []string{"syntax", "naming"},
		},
		MinSeverity: "warning",
	}

	cfg, err := fc.ToConfig(allRules)
	if err != nil {
		t.Fatalf("ToConfig failed: %v", err)
	}

	if len(cfg.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(cfg.Rules))
	}
	if cfg.MinSeverity != tagaudit.SeverityWarning {
		t.Errorf("expected SeverityWarning, got %v", cfg.MinSeverity)
	}
}

func TestToConfig_DisableOnly(t *testing.T) {
	allRules := []tagaudit.Rule{
		&mockRule{"syntax"},
		&mockRule{"naming"},
		&mockRule{"duplicates"},
	}

	fc := &FileConfig{
		Rules: RulesConfig{
			Disable: []string{"duplicates"},
		},
	}

	cfg, err := fc.ToConfig(allRules)
	if err != nil {
		t.Fatalf("ToConfig failed: %v", err)
	}

	if len(cfg.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(cfg.Rules))
	}
}

func TestToConfig_BadSeverity(t *testing.T) {
	fc := &FileConfig{MinSeverity: "critical"}
	_, err := fc.ToConfig(nil)
	if err == nil {
		t.Error("expected error for bad severity")
	}
}

func TestParseSeverity(t *testing.T) {
	tests := []struct {
		input string
		want  tagaudit.Severity
		err   bool
	}{
		{"error", tagaudit.SeverityError, false},
		{"warning", tagaudit.SeverityWarning, false},
		{"warn", tagaudit.SeverityWarning, false},
		{"info", tagaudit.SeverityInfo, false},
		{"INFO", tagaudit.SeverityInfo, false},
		{"bad", 0, true},
	}
	for _, tt := range tests {
		got, err := ParseSeverity(tt.input)
		if (err != nil) != tt.err {
			t.Errorf("ParseSeverity(%q) error = %v, want error = %v", tt.input, err, tt.err)
		}
		if err == nil && got != tt.want {
			t.Errorf("ParseSeverity(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
