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
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		NamingConventions: map[string]string{
			"json": "snake_case",
		},
		RequiredTagKeys: []string{"json"},
	}
}
