package gen

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"github.com/LemonNekoGH/easinteraction-for-cadence/cmd/easi-gen/internal/gen/templates"
	"github.com/LemonNekoGH/easinteraction-for-cadence/cmd/easi-gen/internal/types"
	"github.com/LemonNekoGH/easinteraction-for-cadence/cmd/easi-gen/pkg/string_utils"
	"github.com/LemonNekoGH/easinteraction-for-cadence/cmd/easi-gen/pkg/typeconv"
	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/common"
	"go/format"
)

var (
	ErrNoTopLevelContract = errors.New("no top level contract found")
)

type Generator struct {
	pkgName  string
	contract *ast.CompositeDeclaration
	output   *bytes.Buffer
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

func collectCompositeType(cd *ast.CompositeDeclaration, ownerType string) types.CompositeType {
	var (
		fns                 []types.Function
		subTypes            []types.CompositeType
		fields              []types.Field
		name                = cd.Identifier.String()
		ownerTypeForSubType string
	)

	if ownerType != "" {
		ownerTypeForSubType = ownerType + "." + name
	} else {
		ownerTypeForSubType = name
	}

	cdKind := cd.DeclarationKind()
	for _, m := range cd.DeclarationMembers().Declarations() {
		// skip not public
		fmt.Printf("name: %s, type: %s, access: %s\n", m.DeclarationIdentifier(), m.DeclarationKind().Name(), m.DeclarationAccess().Keyword())
		if m.DeclarationAccess() != ast.AccessPublic &&
			m.DeclarationAccess() != ast.AccessPublicSettable &&
			// event initializer has no name and no access keyword
			m.DeclarationKind() != common.DeclarationKindInitializer {
			continue
		}
		mKind := m.DeclarationKind()
		switch {
		// functions
		case mKind == common.DeclarationKindFunction && cdKind == common.DeclarationKindContract:
			f := m.(*ast.FunctionDeclaration)
			contractFn := collectFunctions(f, name)
			fns = append(fns, contractFn)
		// struct, resource, event
		case mKind == common.DeclarationKindStructure, mKind == common.DeclarationKindResource, mKind == common.DeclarationKindEvent:
			d := m.(*ast.CompositeDeclaration)
			subType := collectCompositeType(d, ownerTypeForSubType)
			subTypes = append(subTypes, subType)
		// field
		case mKind == common.DeclarationKindField:
			f := m.(*ast.FieldDeclaration)
			field := collectField(f)
			fields = append(fields, field)
		// event initializer
		case mKind == common.DeclarationKindInitializer && cdKind == common.DeclarationKindEvent:
			f := m.(*ast.SpecialFunctionDeclaration)
			fields = collectEventInitializer(f)
		}
	}

	var c types.CompositeType
	switch cd.DeclarationKind() {
	case common.DeclarationKindStructure:
		c = &types.Struct{}
	case common.DeclarationKindResource:
		c = &types.Resource{}
	case common.DeclarationKindContract:
		c = &types.Contract{}
	case common.DeclarationKindEvent:
		c = &types.Event{}
	}
	c.SetFields(fields)
	c.SetFunctions(fns)
	c.SetSubTypes(subTypes)
	c.SetName(name)
	c.SetOwnerType(ownerType)
	return c
}

// event initializer is function, but it's field in TransactionResult.Event
func collectEventInitializer(s *ast.SpecialFunctionDeclaration) []types.Field {
	f := s.FunctionDeclaration
	// travel all params
	var fields []types.Field
	for _, p := range f.ParameterList.Parameters {
		fmt.Printf("name: %s, type: %s\n", p.Identifier.String(), p.TypeAnnotation.Type.String())
		fields = append(fields, types.Field{
			Name:   p.Identifier.String(),
			GoName: string_utils.FirstLetterUppercase(p.Identifier.String()),
			Type:   p.TypeAnnotation.Type.String(),
		})
	}
	return fields
}

func collectField(f *ast.FieldDeclaration) types.Field {
	return types.Field{
		Name:   f.Identifier.String(),
		GoName: string_utils.FirstLetterUppercase(f.Identifier.String()),
		Type:   f.TypeAnnotation.Type.String(),
	}
}

func collectFunctions(f *ast.FunctionDeclaration, ownerTypeName string) types.Function {
	fnName := f.Identifier.String()
	fn := types.Function{
		OwnerTypeName: ownerTypeName,
		Name:          fnName,
		GoName:        string_utils.FirstLetterUppercase(fnName),
	}
	// get return type
	retType := f.ReturnTypeAnnotation
	if retType != nil {
		fn.ReturnSimpleType = retType.Type.String()
	}
	// travel all params
	var params []types.FunctionParam
	for _, p := range f.ParameterList.Parameters {
		params = append(params, types.FunctionParam{
			Name:  p.Identifier.String(),
			Label: p.Label,
			Type:  p.TypeAnnotation.Type.String(),
		})
	}
	fn.Params = params
	return fn
}

func (g *Generator) gen(contract *types.Contract) error {
	// generate go code
	return templates.TemplateInteraction.Execute(g.output, contract)
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
	top := collectCompositeType(g.contract, "")
	contract := top.(*types.Contract)
	contract.PkgName = g.pkgName
	contract.FlattenSubTypes() // flatten all nested subtypes
	assignGoTypes(contract)

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

// assignGoTypes assigns correct go type for generation
func assignGoTypes(c *types.Contract) {
	// fields
	var newFields []types.Field
	for _, f := range c.GetFields() {
		f.GoType, f.Type = typeconv.ByName(f.Type, c.GetSubTypes())
		newFields = append(newFields, f)
	}
	c.SetFields(newFields)
	// subtypes field
	var newSubTypes []types.CompositeType
	for _, s := range c.GetSubTypes() {
		var newFields2 []types.Field
		for _, f := range s.GetFields() {
			f.GoType, f.Type = typeconv.ByName(f.Type, s.GetSubTypes())
			newFields2 = append(newFields2, f)
		}
		s.SetFields(newFields2)
		newSubTypes = append(newSubTypes, s)
	}
	c.SetSubTypes(newSubTypes)
	// functions
	var newFns []types.Function
	for _, f := range c.GetFunctions() {
		// process function param go type
		var newPs []types.FunctionParam
		for _, p := range f.Params {
			p.GoType, p.Type = typeconv.ByName(p.Type, c.GetSubTypes())
			newPs = append(newPs, p)
		}
		f.Params = newPs

		// process function return go type
		f.ReturnGoType, f.ReturnType = typeconv.ByName(f.ReturnSimpleType, c.GetSubTypes())
		newFns = append(newFns, f)
	}
	c.SetFunctions(newFns)
}
