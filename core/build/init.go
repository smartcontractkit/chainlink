//go:build !dev

package build

import "testing"

func init() {
	if testing.Testing() {
		mode = Test
	} else {
		mode = Prod
	}
}
