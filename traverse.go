package tagaudit

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/fatih/structtag"
)

// structDecl holds the AST and type information for a discovered struct.
type structDecl struct {
	name    string
	astNode *ast.StructType
	styp    *types.Struct
}

// resolveStruct checks whether a TypeSpec declares a named struct type
// and returns its name, AST node, and resolved types.Struct.
func resolveStruct(ts *ast.TypeSpec, info *types.Info) (name string, astNode *ast.StructType, styp *types.Struct, ok bool) {
	st, isStruct := ts.Type.(*ast.StructType)
	if !isStruct {
		return "", nil, nil, false
	}
	obj := info.Defs[ts.Name]
	if obj == nil {
		return "", nil, nil, false
	}
	// Unalias handles type aliases (e.g., type Foo = Bar) transparently;
	// it is a no-op for non-alias types.
	named, isNamed := types.Unalias(obj.Type()).(*types.Named)
	if !isNamed {
		return "", nil, nil, false
	}
	underlying, isStruct := named.Underlying().(*types.Struct)
	if !isStruct {
		return "", nil, nil, false
	}
	return ts.Name.Name, st, underlying, true
}

// extractStructs walks AST files and finds all struct type declarations,
// including named types, anonymous structs in composite literals, and
// anonymous structs in variable declarations.
func extractStructs(files []*ast.File, info *types.Info) []structDecl {
	var decls []structDecl
	seen := make(map[*ast.StructType]bool)

	for _, file := range files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.TypeSpec:
				name, st, styp, ok := resolveStruct(node, info)
				if ok && !seen[st] {
					seen[st] = true
					decls = append(decls, structDecl{name: name, astNode: st, styp: styp})
				}

			case *ast.CompositeLit:
				st, ok := node.Type.(*ast.StructType)
				if !ok || seen[st] {
					break
				}
				tv, exists := info.Types[node.Type]
				if !exists {
					break
				}
				if styp, ok := tv.Type.(*types.Struct); ok {
					seen[st] = true
					decls = append(decls, structDecl{name: "", astNode: st, styp: styp})
				}

			case *ast.ValueSpec:
				st, ok := node.Type.(*ast.StructType)
				if !ok || seen[st] {
					break
				}
				tv, exists := info.Types[node.Type]
				if !exists {
					break
				}
				if styp, ok := tv.Type.(*types.Struct); ok {
					seen[st] = true
					decls = append(decls, structDecl{name: "", astNode: st, styp: styp})
				}
			}
			return true
		})
	}

	return decls
}

// buildFieldInfos constructs FieldInfo for each field in a struct,
// handling the AST/types index divergence from multi-name declarations.
func buildFieldInfos(fset *token.FileSet, structName string, styp *types.Struct, astNode *ast.StructType) []FieldInfo {
	var fields []FieldInfo
	typesIdx := 0

	for _, astField := range astNode.Fields.List {
		// Number of names this AST field declares (0 means embedded/anonymous)
		nameCount := len(astField.Names)
		if nameCount == 0 {
			nameCount = 1
		}

		for i := 0; i < nameCount; i++ {
			if typesIdx >= styp.NumFields() {
				break
			}

			field := styp.Field(typesIdx)
			rawTag := ""
			if astField.Tag != nil {
				// Strip the backtick quotes
				rawTag = astField.Tag.Value
				rawTag = strings.Trim(rawTag, "`")
			}

			var tags *structtag.Tags
			if rawTag != "" {
				parsed, err := structtag.Parse(rawTag)
				if err == nil {
					tags = parsed
				}
				// If parse fails, tags stays nil — the syntax rule will report it
			}

			fields = append(fields, FieldInfo{
				Fset:       fset,
				StructName: structName,
				Field:      field,
				ASTField:   astField,
				Tags:       tags,
				RawTag:     rawTag,
				Index:      typesIdx,
			})
			typesIdx++
		}
	}

	return fields
}
