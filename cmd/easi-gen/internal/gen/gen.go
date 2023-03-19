package gen

import (
	"bytes"
	_ "embed"
	"errors"
	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/common"
	"text/template"
)

var (
	ErrNoTopLevelContract = errors.New("no top level contract found")
)

//go:embed type_definition.template
var templateTypeDefinition string

type Generator struct {
	pkgName  string
	contract *ast.CompositeDeclaration
	output   *bytes.Buffer
}

type contractType struct {
	PkgName      string
	ContractName string
}

func NewGenerator(pkgName string) *Generator {
	return &Generator{
		pkgName: pkgName,
		output:  bytes.NewBuffer([]byte{}),
	}
}

// Walk implements ast.Walker, to find top level contract of program
func (g *Generator) Walk(e ast.Element) ast.Walker {
	if e == nil {
		return nil
	}

	// contract found, skip all
	if g.contract != nil {
		return nil
	}

	// skip not composite declaration
	if e.ElementType() != ast.ElementTypeCompositeDeclaration {
		return g
	}

	// skip not contract declaration
	d, ok := e.(*ast.CompositeDeclaration)
	if !ok || d.DeclarationKind() != common.DeclarationKindContract {
		return g
	}

	g.contract = d
	return g
}

func (g *Generator) genContractTypeDefinition() error {
	cName := g.contract.Identifier.String()
	t, err := template.New("TypeDefinition").Parse(templateTypeDefinition)
	if err != nil {
		return err
	}
	err = t.Execute(g.output, contractType{
		PkgName:      g.pkgName,
		ContractName: cName,
	})
	if err != nil {
		return err
	}
	return nil
}

// Gen generates go functions
func (g *Generator) Gen(cdc *ast.Program) error {
	ast.Walk(g, cdc)
	// no contracts found, return error
	if g.contract == nil {
		return ErrNoTopLevelContract
	}
	// gen contract type definition
	err := g.genContractTypeDefinition()
	if err != nil {
		return err
	}
	// travel all public functions
	for _, m := range g.contract.DeclarationMembers().Declarations() {
		// skip not function or not public
		if m.DeclarationKind() != common.DeclarationKindFunction || m.DeclarationAccess() != ast.AccessPublic {
			continue
		}
		//f := m.(*ast.FunctionDeclaration)
	}
	return nil
}

// GetOutput returns generated code
func (g *Generator) GetOutput() *bytes.Buffer {
	return g.output
}
