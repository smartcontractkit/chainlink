package matches

import (
	"context"

	"github.com/stretchr/testify/mock"
)

func anyContext(_ context.Context) bool {
	return true
}

// AnyContext is an argument matcher that matches any argument of type context.Context.
var AnyContext = mock.MatchedBy(anyContext)
