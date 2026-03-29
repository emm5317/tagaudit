package rules

import (
	"go/token"
	"go/types"

	"github.com/fatih/structtag"
)

func parseTag(raw string) (*structtag.Tags, error) {
	return structtag.Parse(raw)
}

// fakeVar creates a *types.Var with the given name for unit testing.
func fakeVar(name string) *types.Var {
	return types.NewVar(token.NoPos, nil, name, types.Typ[types.String])
}
