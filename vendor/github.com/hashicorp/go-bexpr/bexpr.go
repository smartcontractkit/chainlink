// bexpr is an implementation of a generic boolean expression evaluator.
// The general goal is to be able to evaluate some expression against some
// arbitrary data and get back a boolean of whether or not the data
// was matched by the expression
package bexpr

//go:generate pigeon -o grammar/grammar.go -optimize-parser grammar/grammar.peg
//go:generate goimports -w grammar/grammar.go

import (
	"github.com/hashicorp/go-bexpr/grammar"
	"github.com/mitchellh/pointerstructure"
)

// HookFn provides a way to translate one reflect.Value to another during
// evaluation by bexpr. This facilitates making Go structures appear in a way
// that matches the expected JSON Pointers used for evaluation. This is
// helpful, for example, when working with protocol buffers' well-known types.
type ValueTransformationHookFn = pointerstructure.ValueTransformationHookFn

type Evaluator struct {
	// The syntax tree
	ast                     grammar.Expression
	tagName                 string
	valueTransformationHook ValueTransformationHookFn
}

func CreateEvaluator(expression string, opts ...Option) (*Evaluator, error) {
	parsedOpts := getOpts(opts...)
	var parserOpts []grammar.Option
	if parsedOpts.withMaxExpressions != 0 {
		parserOpts = append(parserOpts, grammar.MaxExpressions(parsedOpts.withMaxExpressions))
	}

	ast, err := grammar.Parse("", []byte(expression), parserOpts...)
	if err != nil {
		return nil, err
	}

	eval := &Evaluator{
		ast:                     ast.(grammar.Expression),
		tagName:                 parsedOpts.withTagName,
		valueTransformationHook: parsedOpts.withHookFn,
	}

	return eval, nil
}

func (eval *Evaluator) Evaluate(datum interface{}) (bool, error) {
	return evaluate(eval.ast, datum, WithTagName(eval.tagName), WithHookFn(eval.valueTransformationHook))
}
