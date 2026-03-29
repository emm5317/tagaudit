package rules

import "github.com/emm5317/tagaudit"

// All returns all built-in rules.
func All() []tagaudit.Rule {
	return []tagaudit.Rule{
		&SyntaxRule{},
		&NamingRule{},
		&OptionsRule{},
		&UnexportedRule{},
		&UnknownKeysRule{},
		&CompletenessRule{},
		&DuplicatesRule{},
		&ShadowRule{},
	}
}

// DefaultConfig returns a Config pre-populated with all built-in rules
// and sensible defaults (snake_case for json, json as required tag key).
func DefaultConfig() *tagaudit.Config {
	return &tagaudit.Config{
		Rules: All(),
		NamingConventions: map[string]string{
			"json": "snake_case",
		},
		RequiredTagKeys: []string{"json"},
	}
}
