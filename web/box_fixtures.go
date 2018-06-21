// +build test

package web

import (
	"github.com/gobuffalo/packr"
)

// NewBox returns fixtures when a test tag is present
func NewBox() packr.Box {
	return packr.NewBox("../internal/fixtures/gui/dist/")
}
