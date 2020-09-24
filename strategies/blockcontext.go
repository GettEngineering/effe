package strategies

import (
	"go/ast"

	"github.com/GettEngineering/effe/fields"
)

type VarBuilder func(t ast.Expr) *ast.Ident

func NewBlockContext() *BlockContext {
	return &BlockContext{
		Input:  new(ast.FieldList),
		Output: new(ast.FieldList),
		Vars:   make(map[string]*ast.Ident),
	}
}

type BlockContext struct {
	Input   *ast.FieldList
	Output  *ast.FieldList
	Vars    map[string]*ast.Ident
	Builder VarBuilder
}

//nolint:gocognit
func (b *BlockContext) CalculateInput(calls []ComponentCall) {
	for index, c := range calls {
		if c.Input() == nil {
			continue
		}
		for _, inputField := range c.Input().List {
			var foundSourceOfArg bool
			for _, previous := range calls[:index] {
				if previous.Output() != nil && fields.FindFieldWithType(previous.Output().List, inputField.Type) != nil {
					foundSourceOfArg = true
					break
				}
			}

			if !foundSourceOfArg {
				existsInInput := fields.FindFieldWithType(b.Input.List, inputField.Type)
				ident, _ := b.genVariable(inputField.Type)
				if existsInInput == nil {
					b.addInput(ident, inputField.Type)
				} else {
					existsInInput.Names = []*ast.Ident{ident}
				}
			}
		}
	}
}

func (b *BlockContext) OutputList() []*ast.Field {
	return b.Output.List
}

//nolint:gocognit
func (b *BlockContext) CalculateOutput(calls []ComponentCall) {
	for index, c := range calls {
		if c.Output() == nil {
			continue
		}
		for _, outputField := range c.Output().List {
			var foundUsageOfOutput bool

			for _, next := range calls[index+1:] {
				if next.Input() == nil {
					continue
				}
				if next.Input() != nil && fields.FindFieldWithType(next.Input().List, outputField.Type) != nil {
					foundUsageOfOutput = true
					break
				}
			}

			if !foundUsageOfOutput && fields.FindFieldWithType(b.Output.List, outputField.Type) == nil {
				b.addOutput(outputField.Type)
			}
		}
	}
}
func (b *BlockContext) addInput(name *ast.Ident, t ast.Expr) {
	b.Input.List = append(b.Input.List, &ast.Field{
		Names: []*ast.Ident{name},
		Type:  t,
	})
}

func (b *BlockContext) addOutput(t ast.Expr) {
	b.Output.List = append(b.Output.List, &ast.Field{
		Type: t,
	})
}

func (b BlockContext) FindInputByType(t ast.Expr) *ast.Field {
	return fields.FindFieldWithType(b.Input.List, t)
}

func (b *BlockContext) AddInput(t ast.Expr) *ast.Field {
	f := b.FindInputByType(t)
	if f != nil {
		return f
	}

	v := b.Builder(t)
	field := &ast.Field{
		Names: []*ast.Ident{v},
		Type:  t,
	}
	b.Input.List = append(b.Input.List, field)
	b.Vars[fields.GetTypeStrName(t)] = v
	return field
}

func (b *BlockContext) genVariable(t ast.Expr) (*ast.Ident, bool) {
	typeStr := fields.GetTypeStrName(t)
	v, ok := b.Vars[typeStr]
	if !ok {
		v = b.Builder(t)
		b.Vars[typeStr] = v
	}
	return v, ok
}

func (b *BlockContext) BuildInputVars(input *ast.FieldList) []*ast.Field {
	args := make([]*ast.Field, len(input.List))
	for index, inputField := range input.List {
		v, ok := b.genVariable(inputField.Type)
		if !ok && fields.FindFieldWithType(b.Input.List, inputField.Type) == nil {
			b.addInput(v, inputField.Type)
		}
		args[index] = &ast.Field{
			Type:  inputField.Type,
			Names: []*ast.Ident{v},
		}
	}
	return args
}

func (b *BlockContext) BuildOutputVars(output *ast.FieldList) ([]*ast.Field, bool) {
	vars := make([]*ast.Field, len(output.List))

	allOutputVarsExist := true
	for index, outputField := range output.List {
		v, ok := b.genVariable(outputField.Type)
		vars[index] = &ast.Field{
			Type:  outputField.Type,
			Names: []*ast.Ident{v},
		}
		if !ok {
			allOutputVarsExist = false
		}
	}
	return vars, allOutputVarsExist
}
