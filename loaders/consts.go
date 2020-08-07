package loaders

import "github.com/pkg/errors"

const (
	// wrap helpers
	WrapExprType    = "Wrap"
	BeforeExprType  = "Before"
	SuccessExprType = "Success"

	// decision
	DecisionExprType = "Decision"
	CaseExprType     = "Case"

	//others
	StepExprType    = "Step"
	FailureExprType = "Failure"
)

var (
	ErrNoExpr = errors.New("no expr in args")
)
