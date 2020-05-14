package bulletprooftxmanager

import (
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type EthConfirmer interface {
	store.HeadTrackable
}

type ethConfirmer struct{}

func (ec *ethConfirmer) Connect(*models.Head) error {
	// TODO
	return nil
}

func (ec *ethConfirmer) Disconnect() {
	// TODO
	return
}

func (ec *ethConfirmer) OnNewHead(*models.Head) {
	// TODO
	return
}
