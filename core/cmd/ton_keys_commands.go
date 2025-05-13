package cmd

import (
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tonkey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

type TonKeyPresenter struct {
	JAID
	presenters.TonKeyResource
}

// RenderTable implements TableRenderer
func (p TonKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Ton Public Key", "User-Friendly Address", "Raw Address"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 Ton Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func (p *TonKeyPresenter) ToRow() []string {
	row := []string{
		p.ID,
		p.PubKey,
		p.UserFriendlyAddress,
		p.RawAddress,
	}

	return row
}

type TonKeyPresenters []TonKeyPresenter

// RenderTable implements TableRenderer
func (ps TonKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Ton Public Key", "User-Friendly Address", "Raw Address"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 Ton Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func NewTonKeysClient(s *Shell) KeysClient {
	return newKeysClient[tonkey.Key, TonKeyPresenter, TonKeyPresenters]("Ton", s)
}
