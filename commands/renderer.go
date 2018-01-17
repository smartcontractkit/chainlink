package commands

import (
	"io"

	"github.com/smartcontractkit/chainlink/utils"
)

type Renderer interface {
	Render(interface{}) error
}

type RendererJSON struct {
	io.Writer
}

func (rj RendererJSON) Render(v interface{}) error {
	b, err := utils.FormatJSON(v)
	if err != nil {
		return err
	}
	if _, err = rj.Write(b); err != nil {
		return err
	}
	return nil
}

type RendererNoOp struct{}

func (rj RendererNoOp) Render(v interface{}) error { return nil }
