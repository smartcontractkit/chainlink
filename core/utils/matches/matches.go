package matches

import (
	"context"

	"github.com/stretchr/testify/mock"
)

func anyContext(_ context.Context) bool {
	return true
}

func anyString(_ string) bool {
	return true
}

// AnyContext is an argument matcher that matches any argument of type context.Context.
var AnyContext = mock.MatchedBy(anyContext)

// AnyString is an argument matcher that matches any argument of type string.
var AnyString = mock.MatchedBy(anyString)
