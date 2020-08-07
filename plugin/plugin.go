//Package plugin implements a plugin layer for component statements.
// Exports the required interface that must be implemented by plugins.
package plugin

import (
	"go/ast"
)

// Context interface uses for searching variable in stack of variables and adding
// a new input to a component function
type Context interface {
	FindInputByType(ast.Expr) *ast.Field
	AddInput(ast.Expr) *ast.Field
	OutputList() []*ast.Field
}

// ComponentStmt interface is used by a plugin
type ComponentStmt interface {
	// Generated code for calling. It can be ast.AssignStmt or  ast.ExprStmt
	Stmt() ast.Stmt

	// Varaiables are used for calling with type.
	InputFields() []*ast.Field

	// Result variables for a component.
	OutputFields() []*ast.Field

	// Statements are called after a component statement.
	ErrStmt() *ast.BlockStmt

	// Component name
	Name() string
}

// Plugin interface that must be implemented by plugins
type Plugin interface {

	// Imports which are added if plugins are used
	Imports() []string

	// Hook for adding new statements before calling a component statement.
	//
	// Second argument - component name
	//
	// Third argument - result of calling InputFields()
	Before(Context, string, []*ast.Field) []ast.Stmt

	// Hook for adding new statements after calling a component statement.
	//
	// Second argument - component name
	//
	// Third argument - result of calling OutputFields()
	Success(Context, string, []*ast.Field) []ast.Stmt

	// Hook for changing a component statement and failure component statements.
	// Return value is a mark - code changed or not. If code has changes then imports are added.
	Change(Context, ComponentStmt) bool
}
