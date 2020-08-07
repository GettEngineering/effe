package fields

import (
	"go/ast"
	"go/types"
	"strings"

	"github.com/iancoleman/strcase"
)

// Generates an identifier for a variable by type
func NewIdentWithType(t ast.Expr) *ast.Ident {
	if GetTypeStrName(t) == "error" {
		return ast.NewIdent("err")
	}
	if GetTypeStrName(t) == "context.Context" {
		return ast.NewIdent("ctx")
	}

	var camelTypeName string
	switch t := t.(type) {
	case *ast.StarExpr:
		switch x := t.X.(type) {
		case *ast.SelectorExpr:
			camelTypeName = strcase.ToLowerCamel(GetTypeStrName(x.Sel) + "Ptr")
		default:
			camelTypeName = GetTypeStrName(t.X) + "Ptr"
		}
	case *ast.Ident:
		n := t.Name

		dotIndex := strings.Index(n, ".")
		if dotIndex != -1 {
			camelTypeName = strcase.ToLowerCamel(string([]byte(n)[dotIndex+1:]))
		} else {
			camelTypeName = n
		}
	case *ast.ArrayType:
		ident := NewIdentWithType(t.Elt)
		ident.Name += "Ar"
		return ident
	case *ast.SelectorExpr:
		camelTypeName = strcase.ToLowerCamel(GetTypeStrName(t.Sel))
		return ast.NewIdent(camelTypeName + "Val")
	default:
		camelTypeName = GetTypeStrName(t)
	}
	return ast.NewIdent(camelTypeName + "Val")
}

// Searches a field with a specific type in an array. Returns a field or nil.
func FindFieldWithType(fields []*ast.Field, t ast.Expr) *ast.Field {
	for _, f := range fields {
		if GetTypeStrName(f.Type) == GetTypeStrName(t) {
			return f
		}
	}
	return nil
}

// Represents a string by expression type
func GetTypeStrName(t ast.Expr) string {
	return types.ExprString(t)
}
