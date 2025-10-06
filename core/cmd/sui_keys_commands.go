package cmd

import (
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/suikey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

type SuiKeyPresenter struct {
	JAID
	presenters.SuiKeyResource
}

// RenderTable implements TableRenderer
func (p SuiKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Sui Public Key"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 Sui Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func (p *SuiKeyPresenter) ToRow() []string {
	row := []string{
		p.ID,
		p.PubKey,
	}

	return row
}

type SuiKeyPresenters []SuiKeyPresenter

// RenderTable implements TableRenderer
func (ps SuiKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Sui Public Key"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 Sui Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func NewSuiKeysClient(s *Shell) KeysClient {
	return newKeysClient[suikey.Key, SuiKeyPresenter, SuiKeyPresenters]("Sui", s)
}
