package rules

import (
	"go/token"
	"go/types"

	"github.com/emm5317/tagaudit"
)

// maxEmbedDepth is the maximum embedding depth to traverse.
// This prevents stack overflow on cyclic type definitions.
const maxEmbedDepth = 10

// unwrapToStruct unwraps pointer, alias, and named type wrappers to
// reach the underlying *types.Struct, if any. It returns the
// intermediate *types.Named (for cycle tracking) and the struct.
func unwrapToStruct(t types.Type) (named *types.Named, styp *types.Struct, ok bool) {
	if ptr, isPtr := t.(*types.Pointer); isPtr {
		t = ptr.Elem()
	}
	t = types.Unalias(t)
	named, _ = t.(*types.Named)
	if named != nil {
		t = named.Underlying()
	}
	styp, ok = t.(*types.Struct)
	return named, styp, ok
}

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

// tagSpanFromInfo returns the byte offsets of the tag literal (including
// backticks) in the source file. Returns (0, 0) if unavailable.
func tagSpanFromInfo(info tagaudit.FieldInfo) (start, end int) {
	if info.Fset == nil || info.ASTField == nil || info.ASTField.Tag == nil {
		return 0, 0
	}
	s := info.Fset.Position(info.ASTField.Tag.Pos())
	e := info.Fset.Position(info.ASTField.Tag.End())
	return s.Offset, e.Offset
}
