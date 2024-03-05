This package contains test implementations for all the components that compose OCR2  LOOPPs. The existing of this directory is partially an intermediate term workaround for lack of packages within `loop/internal` the https://smartcontract-it.atlassian.net/browse/BCF-3039. When that issue is addressed, we can move these test implementations into specific `test` packages for each packaged component.

With all that said, the directory layout is intended to mirror the natural interfaces composition of libocr:
- ocr2: contains static test implementations for reuseable reporting plugin, factory, and plugin provider
- core: represents core-node provided resources used to in calls to core node api(s) eg pipeline_runner, keystore, etc
- types: type definitions used in through tests

Every test implementation follows the pattern wrapping an interface and provider one or more funcs to compare to another
instance of the interface. Every package attempts to share exposed the bare minimum surface area to avoid entanglement with logically separated tests and domains.

In practice this is accomplished by exporting a static implementation and an interface that the implementation satisfies. The interface is used by other packages to declaration dependencies and the implementation is used in the test implementation of those other packaged

Example

```
go
package types

type Evaluator[T any] interface {
     // run methods of other, returns first error or error if
    // result of method invocation does not match expected value
    Evaluate(ctx context.Context, other T) error
}

```
package x

import (
    testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var FooEvaluator = staticFoo{
    expectedStr = "a"
    expectedInt = 1
}


type staticFoo struct {
    types.Foo // the interface to be tested, used to compose a loop
    expectedStr string
    expectedInt int
    ...
}

var _ testtypes.Evaulator[Foo] = staticFoo
var _ types.Foo = staticFoo

func (f Foo) Evaluate(ctx context.Context, other Foo) {
    // test implementation of types.Foo interface
    s, err := other.GetStr()
    if err ! = nil {
        return fmt.Errorf("other failed to get str: %w", err)
    } 
    if s != f.expectedStr {
        return fmt.Errorf(" expected str %s got %s", s.expectedStr, s)
    }
    ...
}

// implements types.Foo
func (f Foo) GetStr() (string, error) {
    return f.expectedStr, nil
}
...

```

```
package y

import (
    testx "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/x"
    testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var BarEvaluator = staticBar{
    expectedFoo = x_test.FooImpl
    expectedBytes = []bytes{1:1}
}


type staticBar struct {
    types.Bar // test implementation of Bar interface
    expectedFoo types_test[types.Foo]
    expectedInt int
    ...
}

// BarEvaluation implement [types.Bar] and [types_test.Evaluator[types.Bar]] to be used in tests
var _ BarEvaluator = staticBar {
    expectedFoo x_test.FooEvaluator
    expectedInt = 7
    ...
}

var _ testtypes[types.Bar] = staticBar

// implement types.Bar
...
// implement types_test.Evaluator[types.Bar]
...

```