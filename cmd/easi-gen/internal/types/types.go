package types

import (
	"bytes"
	"fmt"
	"github.com/LemonNekoGH/easinteraction-for-cadence/cmd/easi-gen/internal/gen/templates"
	"github.com/onflow/cadence/runtime/common"
	"path/filepath"
	"strings"
	"text/template"
)

type CompositeType interface {
	Kind() common.CompositeKind
	GetFields() []Field
	SetFields(fields []Field)
	GetFunctions() []Function
	SetFunctions(fns []Function)
	GetSubTypes() []CompositeType
	SetSubTypes(sb []CompositeType)
	SetName(n string)
	GetName() string
	GetSimpleName() string
	GetGoName() string // GetGoName returns name in go, if it's kind is struct, will add "Struct" prefix, if it's kind is resource, will add "Resource" prefix.
	SetOwnerType(own string)
}

type compositeTypeImpl struct {
	ownerType string
	fields    []Field
	functions []Function
	subTypes  []CompositeType // all nested types will flatten into this field
	name      string
}

func (c *compositeTypeImpl) GetFunctions() []Function {
	return c.functions
}

func (c *compositeTypeImpl) SetFunctions(fns []Function) {
	c.functions = fns
}
func (c *compositeTypeImpl) GetFields() []Field {
	return c.fields
}

func (c *compositeTypeImpl) SetFields(fields []Field) {
	c.fields = fields
}

func (c *compositeTypeImpl) Kind() common.CompositeKind {
	return common.CompositeKindUnknown
}

func (c *compositeTypeImpl) GetSubTypes() []CompositeType {
	return c.subTypes
}

func (c *compositeTypeImpl) SetSubTypes(sb []CompositeType) {
	c.subTypes = sb
}

func (c *compositeTypeImpl) GetName() string {
	if c.ownerType != "" {
		return c.ownerType + "." + c.name
	}
	return c.name
}

func (c *compositeTypeImpl) GetSimpleName() string {
	return c.name
}

func (c *compositeTypeImpl) SetName(n string) {
	c.name = n
}
func (c *compositeTypeImpl) GetGoName() string {
	return c.name
}

func (c *compositeTypeImpl) SetOwnerType(owner string) {
	c.ownerType = owner
}

type Contract struct {
	compositeTypeImpl
	PkgName string
}

func (c *Contract) GetGoName() string {
	return "Contract" + c.name
}

func flattenSubTypes(com CompositeType) []CompositeType {
	var subTypes []CompositeType

	for _, c := range com.GetSubTypes() {
		subTypes = append(subTypes, c)
		subTypes = append(subTypes, flattenSubTypes(c)...)
	}
	com.SetSubTypes(nil)

	return subTypes
}

func (c *Contract) Kind() common.CompositeKind {
	return common.CompositeKindContract
}

// FlattenSubTypes move all subtypes to contracts subtypes
func (c *Contract) FlattenSubTypes() {
	c.SetSubTypes(flattenSubTypes(c))
}

type Struct struct {
	compositeTypeImpl
}

func (s *Struct) Kind() common.CompositeKind {
	return common.CompositeKindStructure
}

func (s *Struct) GetGoName() string {
	return "Struct" + s.name
}

type Resource struct {
	compositeTypeImpl
}

func (r *Resource) Kind() common.CompositeKind {
	return common.CompositeKindResource
}

func (r *Resource) GetGoName() string {
	return "Resource" + r.name
}

type Event struct {
	compositeTypeImpl
}

func (*Event) Kind() common.CompositeKind {
	return common.CompositeKindEvent
}

func (e *Event) GetGoName() string {
	return "Event" + e.name
}

type Field struct {
	Name   string
	GoName string
	Type   string
	GoType string
}

type FunctionParam struct {
	Label  string
	Name   string
	Type   string
	GoType string
}

type Function struct {
	OwnerTypeName    string
	Name             string
	GoName           string // first letter uppercase
	Params           []FunctionParam
	ReturnType       string // for return type of script
	ReturnSimpleType string // for convert to go type
	ReturnGoType     string
	usedCommaCommon  int // use index to check is need to add comma is bad, because AuthAccount generate will skip but index will plus one
	usedCommaAuth    int // use index to check is need to add comma is bad, because not AuthAccount generate will skip but index will plus one
}

func (fn *Function) IsReturnMap() bool {
	return strings.HasPrefix(fn.ReturnGoType, "map")
}

func (fn *Function) AuthorizerCount() int {
	var ret int
	for _, p := range fn.Params {
		if p.Type == "AuthAccount" {
			ret++
		}
	}
	return ret
}

func (fn *Function) GenCadenceScript() (string, error) {
	var (
		t   *template.Template
		err error
	)
	if fn.AuthorizerCount() > 0 {
		// should generate transaction script
		t = templates.TemplateTxScript
	} else {
		// should generate query script
		t = templates.TemplateQueryScript
	}
	result := bytes.NewBuffer([]byte{})
	err = t.Execute(result, fn)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

// CommaCountAll returns length of all params - 1, used to check necessary for comma adding.
func (fn *Function) CommaCountAll() int {
	return len(fn.Params) - 1
}

// CommaCountCommon returns length of all common params - 1, used to check necessary for comma adding.
func (fn *Function) CommaCountCommon() int {
	return len(fn.Params) - fn.AuthorizerCount() - 1
}

// CommaCountAuth returns length of all auth params - 1, used to check necessary for comma adding.
func (fn *Function) CommaCountAuth() int {
	return fn.AuthorizerCount() - 1
}

// AddUsedCommaCommon manually add index
func (fn *Function) AddUsedCommaCommon() int {
	fn.usedCommaCommon += 1
	return fn.usedCommaCommon
}

// AddUsedCommaAuth manually add index
func (fn *Function) AddUsedCommaAuth() int {
	fn.usedCommaAuth += 1
	return fn.usedCommaAuth
}

// getContractPath there are two variants of contracts object
func getContractPath(d any) string {
	switch v := d.(type) {
	case string:
		return v
	case map[string]any:
		return v["source"].(string)
	}
	return ""
}

type FlowJson struct {
	Contracts map[string]any `json:"contracts"`
}

// ResolvePath resolve contracts file path, and concat contract output path
func (f *FlowJson) ResolvePath(flowJsonPath, pkgName, outputDir string) ([]string, []string) {
	flowJsonDir := filepath.Dir(flowJsonPath)
	if outputDir == "" {
		outputDir = filepath.Join(flowJsonDir, pkgName)
		fmt.Println("auto set output path: " + outputDir)
	}
	var (
		sourcePaths []string
		outputPaths []string
	)
	for name, path := range f.Contracts {
		p := getContractPath(path)
		sourcePaths = append(sourcePaths, filepath.Join(flowJsonDir, p))
		outputPaths = append(outputPaths, filepath.Join(outputDir, strings.ToLower(name)+".go")) // use lower case contract name for file name
	}
	return sourcePaths, outputPaths
}
