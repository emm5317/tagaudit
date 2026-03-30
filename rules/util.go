package rules

import (
	"go/token"

	"github.com/emm5317/tagaudit"
)

// posFromInfo extracts the best available position from a FieldInfo.
// Returns a zero Position if no position info is available.
func posFromInfo(info tagaudit.FieldInfo) token.Position {
	if info.Fset == nil || info.ASTField == nil {
		return token.Position{}
	}
	if info.ASTField.Tag != nil {
		return info.Fset.Position(info.ASTField.Tag.Pos())
	}
	return info.Fset.Position(info.ASTField.Pos())
}
