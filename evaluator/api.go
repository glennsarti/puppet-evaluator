package evaluator

import (
	. "github.com/puppetlabs/go-parser/issue"
	. "github.com/puppetlabs/go-parser/parser"
)

type (
	Evaluator interface {
		AddDefinitions(expression Expression)

		ResolveDefinitions()

		Evaluate(expression Expression, scope Scope, loader Loader) (PValue, *ReportedIssue)

		Eval(expression Expression, c EvalContext) PValue

		Logger() Logger
	}

	EvalContext interface {
		Call(name string, args []PValue, block Lambda) PValue

		Evaluate(expr Expression) PValue

		EvaluateIn(expr Expression, scope Scope) PValue

		Evaluator() Evaluator

		Fail(message string)

		WithScope(scope Scope) EvalContext

		ParseAndValidate(str, filename string, singleExpression bool) Expression

		ParseType(str PValue) PType

		Loader() Loader

		Logger() Logger

		StackPop()

		StackPush(location Location)

		StackTop() Location

		Scope() Scope

		Stack() []Location
	}
)

var GeneralFailure func(error string)

var CurrentContext func() EvalContext
