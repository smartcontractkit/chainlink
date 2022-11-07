package services_test

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

type example struct {
	services.ServiceCtx
	g      *services.Group
	workCh chan func() (name string, err error)
}

func (e *example) start(context.Context) error {
	e.g.Go(func() {
		for {
			select {
			case <-e.g.StopChan:
			case <-time.After(time.Minute):
				e.do(e.workCh)
			}
		}
	})
	return nil
}

func (e *example) do(workCh <-chan func() (name string, err error)) {
	// do work until none is left
	for {
		select {
		case <-e.g.StopChan:
			return
		case work, ok := <-workCh:
			if !ok {
				return
			}
			name, err := work()
			if err != nil {
				e.g.SetUnwell(name, err)
			} else {
				e.g.SetWell(name)
			}
		default:
			return
		}
	}
}

func NewExample(lggr logger.Logger) services.ServiceCtx {
	e := &example{
		workCh: make(chan func()),
	}
	e.ServiceCtx, e.g = services.New(services.Spec{
		Name:        "Example",
		Start:       e.start,
		SubServices: nil, // optional
	}, lggr)
	return e
}

func Example() {
	NewExample(logger.NullLogger)
}
