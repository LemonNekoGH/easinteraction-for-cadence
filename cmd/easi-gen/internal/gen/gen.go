package gen

import (
	"bytes"
	_ "embed"
	"errors"
	"github.com/LemonNekoGH/easiteraction-for-cadence/pkg/string_utils"
	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/common"
	"go/format"
	"text/template"
)

var (
	ErrNoTopLevelContract = errors.New("no top level contract found")
)

var (
	//go:embed contract_interaction.gohtml
	templateInteraction string
	//go:embed query_script.gohtml
	templateQueryScript string
	//go:embed tx_script.gohtml
	templateTxScript string
)

type Generator struct {
	pkgName  string
	contract *ast.CompositeDeclaration
	output   *bytes.Buffer
}

type functionParam struct {
	Label  string
	Name   string
	Type   string
	GoType string
}

type contractFunction struct {
	ContractName    string
	Name            string
	GoName          string // first letter uppercase
	Params          []functionParam
	ReturnType      string
	ReturnGoType    string
	GeneratedScript []byte
}

type contractType struct {
	PkgName   string
	Name      string
	Functions []contractFunction
}

func (fn *contractFunction) AuthorizerCount() int {
	var ret int
	for _, p := range fn.Params {
		if p.Type == "AuthAccount" {
			ret++
		}
	}
	return ret
}

func (fn *contractFunction) GenCadenceScript() (string, error) {
	var (
		t   *template.Template
		err error
	)
	if fn.AuthorizerCount() > 0 {
		// should generate transaction script
		t, err = template.New("TxScript").Parse(templateTxScript)
		if err != nil {
			return "", err
		}
	} else {
		// should generate query script
		t, err = template.New("QueryScript").Parse(templateQueryScript)
		if err != nil {
			return "", err
		}
	}
	result := bytes.NewBuffer([]byte{})
	err = t.Execute(result, fn)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

// CommaCountAll returns length of all params - 1, used to check necessary for comma adding.
func (fn *contractFunction) CommaCountAll() int {
	return len(fn.Params) - 1
}

// CommaCountCommon returns length of all common params - 1, used to check necessary for comma adding.
func (fn *contractFunction) CommaCountCommon() int {
	return len(fn.Params) - fn.AuthorizerCount() - 1
}

// CommaCountAuth returns length of all auth params - 1, used to check necessary for comma adding.
func (fn *contractFunction) CommaCountAuth() int {
	return fn.AuthorizerCount() - 1
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

// get all functions of contract
func (g *Generator) collectContractInfos() contractType {
	contract := contractType{
		PkgName: g.pkgName,
		Name:    g.contract.Identifier.String(),
	}
	// travel all public functions
	var fns []contractFunction
	for _, m := range g.contract.DeclarationMembers().Declarations() {
		// skip not function or not public
		if m.DeclarationKind() != common.DeclarationKindFunction || m.DeclarationAccess() != ast.AccessPublic {
			continue
		}
		f := m.(*ast.FunctionDeclaration)
		fnName := f.Identifier.String()
		contractFn := contractFunction{
			ContractName: contract.Name,
			Name:         fnName,
			GoName:       string_utils.FirstLetterUppercase(fnName),
		}
		// get return type
		retType := f.ReturnTypeAnnotation
		if retType != nil {
			contractFn.ReturnType = retType.Type.String()
			contractFn.ReturnGoType = typeMap[retType.Type.String()]
		}
		// travel all params
		var params []functionParam
		for _, p := range f.ParameterList.Parameters {
			params = append(params, functionParam{
				Name:   p.Identifier.String(),
				Label:  p.Label,
				Type:   p.TypeAnnotation.Type.String(),
				GoType: typeMap[p.TypeAnnotation.Type.String()],
			})
		}
		contractFn.Params = params
		fns = append(fns, contractFn)
	}
	contract.Functions = fns
	return contract
}

func (g *Generator) gen(contract contractType) error {
	t, err := template.New("Interaction").Parse(templateInteraction)
	if err != nil {
		return err
	}
	// generate go code
	return t.Execute(g.output, &contract)
}

func (g *Generator) findTopLevelContract(cdc *ast.Program) {
	ast.Walk(g, cdc)
}

// Gen generates go functions
func (g *Generator) Gen(cdc *ast.Program) error {
	g.findTopLevelContract(cdc)
	// no contracts found, return error
	if g.contract == nil {
		return ErrNoTopLevelContract
	}
	contract := g.collectContractInfos()
	// do generate
	err := g.gen(contract)
	if err != nil {
		return err
	}
	// format code
	formatted, err := format.Source(g.output.Bytes())
	if err != nil {
		return err
	}
	g.output = bytes.NewBuffer(formatted)
	return nil
}

// GetOutput returns generated code
func (g *Generator) GetOutput() *bytes.Buffer {
	return g.output
}
