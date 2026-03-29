package tagaudit

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/fatih/structtag"
)

// FieldChecker is implemented by rules that check individual struct fields.
type FieldChecker interface {
	ID() string
	Description() string
	CheckField(info FieldInfo, cfg *Config) []Finding
}

// StructChecker is implemented by rules that need cross-field analysis
// (e.g., duplicate detection, completeness checks).
type StructChecker interface {
	ID() string
	Description() string
	CheckStruct(info StructInfo, cfg *Config) []Finding
}

// FieldInfo provides context about a single struct field to rules.
type FieldInfo struct {
	Fset       *token.FileSet
	StructName string          // name of the containing struct, "" for anonymous
	Field      *types.Var      // type-checked field
	ASTField   *ast.Field      // syntax node
	Tags       *structtag.Tags // parsed tags, nil if parse failed or no tag
	RawTag     string          // raw tag string from source
	Index      int             // field index within the struct
}

// StructInfo provides context about a struct to cross-field rules.
type StructInfo struct {
	Fset       *token.FileSet
	StructName string
	StructType *types.Struct
	ASTNode    *ast.StructType
	Fields     []FieldInfo
}
