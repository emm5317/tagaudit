# tagaudit

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A composable struct tag validation library for Go.

`tagaudit` fills the gap between `fatih/structtag` (parsing only) and `go vet` (not usable as a library). It provides a pluggable rule system with built-in rules for common issues and supports user-defined custom rules.

## Install

```bash
go get github.com/emm5317/tagaudit
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/emm5317/tagaudit"
	_ "github.com/emm5317/tagaudit/rules" // register built-in rules
)

func main() {
	a := tagaudit.New(&tagaudit.Config{
		NamingConventions: map[string]string{
			"json": "snake_case",
			"yaml": "snake_case",
		},
		RequiredTagKeys: []string{"json"},
	})

	findings, err := a.AnalyzePackages("./internal/models/...")
	if err != nil {
		panic(err)
	}

	for _, f := range findings {
		fmt.Println(f)
	}
}
```

## Built-in Rules

| Rule | Type | What it catches |
|------|------|-----------------|
| `syntax` | per-field | Malformed struct tag strings |
| `naming` | per-field | Naming convention violations (e.g., camelCase in json tags when snake_case is expected). Provides auto-fix suggestions. |
| `options` | per-field | Invalid options for well-known tags (e.g., `json:"foo,omitemtpy"`) |
| `unexported` | per-field | Encoding tags on unexported fields (silently ignored at runtime) |
| `unknownkeys` | per-field | Tag keys not in a configured known set (catches typos like `josn`) |
| `completeness` | per-struct | Missing tags when other fields in the struct have them |
| `duplicates` | per-struct | Duplicate tag values within a struct, including via embedding |
| `shadow` | per-struct | Outer field tags that silently override embedded field tags |

## Custom Rules

Implement `FieldChecker` for per-field checks or `StructChecker` for cross-field checks:

```go
type MyRule struct{}

func (r *MyRule) ID() string          { return "my-rule" }
func (r *MyRule) Description() string { return "my custom rule" }

func (r *MyRule) CheckField(info tagaudit.FieldInfo, cfg *tagaudit.Config) []tagaudit.Finding {
	// your logic here
	return nil
}
```

Pass custom rules via config:

```go
a := tagaudit.New(&tagaudit.Config{
	Rules: append(rules.All(), &MyRule{}),
})
```

## Configuration

```go
&tagaudit.Config{
	// Naming conventions per tag key
	NamingConventions: map[string]string{
		"json": "snake_case",  // "snake_case", "camelCase", "PascalCase", "kebab-case"
	},

	// Require these tag keys on all exported fields if any field has them
	RequiredTagKeys: []string{"json"},

	// Only allow these tag keys (catches typos). nil = disabled.
	KnownTagKeys: []string{"json", "db", "yaml"},
}
```

## License

[MIT](LICENSE)
