package utils

import (
	"fmt"
	"strings"

	"go.uber.org/multierr"
)

type multiErrorList []error

// MultiErrorList returns an error which formats underlying errors as a list, or nil if err is nil.
func MultiErrorList(err error) error {
	if err == nil {
		return nil
	}

	return multiErrorList(multierr.Errors(err))
}

func (m multiErrorList) Error() string {
	l := len(m)
	if l == 1 {
		return m[0].Error()
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d errors:", l)
	for i, e := range m {
		fmt.Fprintf(&sb, "\n\t%d) %v", i+1, e)
	}
	return sb.String()
}
