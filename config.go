package tagaudit

// Config controls which rules run and how they behave.
type Config struct {
	// Rules to run. Each must implement FieldChecker and/or StructChecker.
	// If nil or empty, no rules are applied — use rules.All() or
	// rules.DefaultConfig() to include the built-in rule set.
	Rules []Rule

	// NamingConventions maps a tag key to the expected naming convention.
	// Recognized values: "snake_case", "camelCase", "PascalCase", "kebab-case".
	NamingConventions map[string]string

	// RequiredTagKeys: if any exported field in a struct has one of these
	// tag keys, all exported fields in that struct must also have it.
	RequiredTagKeys []string

	// KnownTagKeys: if non-nil, any tag key not in this set produces a finding.
	// If nil, the unknown-key check is disabled.
	KnownTagKeys []string

	// KnownOptions maps tag keys to additional valid options that supplement
	// the built-in defaults (json, xml, yaml). For example, to allow gorm
	// options: KnownOptions: map[string][]string{"gorm": {"primaryKey", "autoIncrement"}}.
	KnownOptions map[string][]string

	// MinSeverity filters findings to only include those at or above this
	// severity level. The zero value (SeverityInfo) includes all findings.
	// SeverityWarning = warnings and errors, SeverityError = errors only.
	MinSeverity Severity
}

// DefaultRulesFunc, if non-nil, is called by DefaultConfig to populate the
// Rules field. The rules package sets this in its init function so that
// tagaudit.DefaultConfig() automatically includes all built-in rules without
// creating a circular import.
var DefaultRulesFunc func() []Rule

// BaseConfig returns the skeleton Config with sensible defaults but no rules.
// Use DefaultConfig() to also get rules (when the rules package is imported),
// or build the Config manually with rules.All() / rules.DefaultConfig().
func BaseConfig() *Config {
	return &Config{
		NamingConventions: map[string]string{
			"json": "snake_case",
		},
		RequiredTagKeys: []string{"json"},
	}
}

// DefaultConfig returns a Config with sensible defaults. If the rules package
// has been imported (directly or transitively), Rules is populated via
// DefaultRulesFunc; otherwise Rules is nil and no rules are applied.
func DefaultConfig() *Config {
	cfg := BaseConfig()
	if DefaultRulesFunc != nil {
		cfg.Rules = DefaultRulesFunc()
	}
	return cfg
}
