// +build !test

package web

import (
	"github.com/gobuffalo/packr"
)

// NewBox returns production build artifacts when the test tag is not present
func NewBox() packr.Box {
	return packr.NewBox("../gui/dist/")
}
