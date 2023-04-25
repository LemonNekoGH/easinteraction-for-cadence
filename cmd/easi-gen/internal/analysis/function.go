package analysis

import (
	"fmt"
	"github.com/onflow/cadence/runtime/ast"
)

// FunctionBlockAnalysis analyzes a function block
type FunctionBlockAnalysis struct {
	functionName  string
	owner         string
	contract      string
	functionBlock *ast.Block
	result        FunctionBlockAnalysisResult
}

// NewFunctionBlockAnalysis creates a new FunctionBlockAnalysis
func NewFunctionBlockAnalysis(
	functionName string,
	owner string,
	contract string,
	functionBlock *ast.Block,
) *FunctionBlockAnalysis {
	return &FunctionBlockAnalysis{
		functionName:  functionName,
		owner:         owner,
		contract:      contract,
		functionBlock: functionBlock,
	}
}

func (a *FunctionBlockAnalysis) getFinalOwnerOfTarget(memberExp ast.Expression) string {
	switch e := memberExp.(type) {
	case *ast.IdentifierExpression:
		return e.Identifier.String()
	case *ast.MemberExpression:
		return a.getFinalOwnerOfTarget(e.Expression)
	case *ast.IndexExpression:
		return a.getFinalOwnerOfTarget(e.TargetExpression)
	default:
		// unknown
		return ""
	}
}

// analyzeAssignment analyzes an assignment statement
func (a *FunctionBlockAnalysis) analyzeAssignment(s *ast.AssignmentStatement) {
	finalOwner := a.getFinalOwnerOfTarget(s.Target)
	fmt.Println("member final owner: " + finalOwner)
	// if owner a contract, final owner should be "self", if not, it should be the contract name
	if a.isOwnerContract() {
		a.result.WillChangeState = finalOwner == "self"
	} else {
		a.result.WillChangeState = finalOwner == a.contract
	}
}

// Walk implements the ast.Walker interface
func (a *FunctionBlockAnalysis) Walk(e ast.Element) ast.Walker {
	if e != nil {
		//fmt.Printf("%v: %v\n", e, e.ElementType())
	}

	switch s := e.(type) {
	case *ast.AssignmentStatement:
		a.analyzeAssignment(s)
	}

	if a.result.WillChangeState {
		// detected state change, return
		return nil
	}

	// TODO: find invocation statement
	return a
}

// Analyze analyzes the function block
func (a *FunctionBlockAnalysis) Analyze() FunctionBlockAnalysisResult {
	ast.Walk(a, a.functionBlock)
	return a.result
}

// isOwnerContract checks if the owner is a contract
func (a *FunctionBlockAnalysis) isOwnerContract() bool {
	return a.owner == a.contract
}

// FunctionBlockAnalysisResult is the result of the analysis
type FunctionBlockAnalysisResult struct {
	WillChangeState bool
}
