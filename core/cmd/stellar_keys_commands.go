package cmd

import (
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/stellarkey"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

type StellarKeyPresenter struct {
	JAID
	presenters.StellarKeyResource
}

// RenderTable implements TableRenderer
func (p StellarKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Stellar Public Key"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 Stellar Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func (p StellarKeyPresenter) ToRow() []string {
	row := []string{
		p.ID,
		p.PubKey,
	}

	return row
}

type StellarKeyPresenters []StellarKeyPresenter

// RenderTable implements TableRenderer
func (ps StellarKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Stellar Public Key"}
	rows := make([][]string, 0, len(ps))

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 Stellar Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func NewStellarKeysClient(s *Shell) KeysClient {
	return newKeysClient[stellarkey.Key, StellarKeyPresenter, StellarKeyPresenters]("Stellar", s)
}
