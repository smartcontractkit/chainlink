package pipeline

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type Expression interface {
	Evaluate(inputs []Result) bool
}

type alwaysTrueExpression struct{}

func (alwaysTrueExpression) Evaluate([]Result) bool {
	return true
}

type ConditionalExpression struct {
	comparator string
}

func (e ConditionalExpression) Evaluate(inputs []Result) bool {
	var (
		lhs ObjectParam
		rhs ObjectParam
	)

	err := multierr.Combine(
		errors.Wrap(ResolveParam(&lhs, From(Input(inputs, 0))), "left hand side of expression"),
		errors.Wrap(ResolveParam(&rhs, From(Input(inputs, 0))), "right hand side of expression"),
	)
	if err != nil {
		// TODO: need a way to error here, or... can we error at parse time?
	}
	fmt.Println("lhs ==>", lhs)
	fmt.Println("rhs ==>", rhs)

	if e.comparator == "=" {
		return lhs.Equals(rhs)
	}

	return false
}
