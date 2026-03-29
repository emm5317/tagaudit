package rules

import "github.com/emm5317/tagaudit"

func init() {
	tagaudit.SetDefaultRules(func() []any {
		return All()
	})
}

// All returns all built-in rules.
func All() []any {
	return []any{
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
