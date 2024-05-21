package cmd

import (
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

type StarkNetKeyPresenter struct {
	JAID
	presenters.StarkNetKeyResource
}

// RenderTable implements TableRenderer
func (p StarkNetKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Stark Public Key"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 StarkNet Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func (p *StarkNetKeyPresenter) ToRow() []string {
	row := []string{
		p.ID,
		p.StarkKey,
	}

	return row
}

type StarkNetKeyPresenters []StarkNetKeyPresenter

// RenderTable implements TableRenderer
func (ps StarkNetKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Stark Public Key"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 StarkNet Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func NewStarkNetKeysClient(s *Shell) KeysClient {
	return newKeysClient[starkkey.Key, StarkNetKeyPresenter, StarkNetKeyPresenters]("StarkNet", s)
}
