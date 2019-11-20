package bulletin_board

import (
	"fmt"
	"regexp"
	"testing"
)

var r regexp.Regexp

func TestEmptyRegexp(t *testing.T) {
	fmt.Println("r", r)
}
