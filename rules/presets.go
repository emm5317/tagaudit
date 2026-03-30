package rules

import "github.com/emm5317/tagaudit"

// DefaultKnownTagKeys is a reasonable set of well-known struct tag keys.
var DefaultKnownTagKeys = []string{
	"json", "yaml", "xml", "toml", "csv",
	"db", "gorm", "bson", "redis",
	"mapstructure", "env", "validate",
	"form", "query", "param", "header", "binding",
	"protobuf", "msgpack",
}

// JSONPreset returns a Config tuned for JSON-heavy codebases.
func JSONPreset() *tagaudit.Config {
	return &tagaudit.Config{
		Rules: All(),
		NamingConventions: map[string]string{
			"json": "snake_case",
		},
		RequiredTagKeys: []string{"json"},
		KnownTagKeys:    DefaultKnownTagKeys,
	}
}

// APIModelPreset returns a Config for API model structs using json + yaml + validate.
func APIModelPreset() *tagaudit.Config {
	return &tagaudit.Config{
		Rules: All(),
		NamingConventions: map[string]string{
			"json": "snake_case",
			"yaml": "snake_case",
		},
		RequiredTagKeys: []string{"json"},
		KnownTagKeys:    DefaultKnownTagKeys,
		KnownOptions: map[string][]string{
			"validate": {"required", "min", "max", "len", "email", "url", "uuid", "oneof", "gte", "lte", "gt", "lt", "dive", "omitempty"},
		},
	}
}

// GORMPreset returns a Config for GORM model structs.
func GORMPreset() *tagaudit.Config {
	return &tagaudit.Config{
		Rules: All(),
		NamingConventions: map[string]string{
			"json": "snake_case",
		},
		RequiredTagKeys: []string{"json"},
		KnownTagKeys:    DefaultKnownTagKeys,
		KnownOptions: map[string][]string{
			"gorm": {
				"primaryKey", "autoIncrement", "column", "type", "size",
				"index", "unique", "not null", "default", "embedded",
				"embeddedPrefix", "foreignKey", "references", "constraint",
				"many2many", "joinForeignKey", "joinReferences",
			},
		},
	}
}
