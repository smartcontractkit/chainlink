package cmd

import (
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/aptoskey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

type AptosKeyPresenter struct {
	JAID
	presenters.AptosKeyResource
}

// RenderTable implements TableRenderer
func (p AptosKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Aptos Public Key"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 Aptos Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func (p *AptosKeyPresenter) ToRow() []string {
	row := []string{
		p.ID,
		p.PubKey,
	}

	return row
}

type AptosKeyPresenters []AptosKeyPresenter

// RenderTable implements TableRenderer
func (ps AptosKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Aptos Public Key"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 Aptos Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func NewAptosKeysClient(s *Shell) KeysClient {
	return newKeysClient[aptoskey.Key, AptosKeyPresenter, AptosKeyPresenters]("Aptos", s)
}
