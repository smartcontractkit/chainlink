package cmd

import (
	starkkey "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/keys"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

type StarknetKeyPresenter struct {
	JAID
	presenters.StarknetKeyResource
}

// RenderTable implements TableRenderer
func (p StarknetKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Account Address", "Stark Public Key"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 Starknet Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func (p *StarknetKeyPresenter) ToRow() []string {
	row := []string{
		p.ID,
		p.AccountAddr,
		p.StarkKey,
	}

	return row
}

type StarknetKeyPresenters []StarknetKeyPresenter

// RenderTable implements TableRenderer
func (ps StarknetKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Account Address", "Stark Public Key"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 Starknet Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func NewStarknetKeysClient(c *Client) KeysClient {
	return newKeysClient[starkkey.Key, StarknetKeyPresenter, StarknetKeyPresenters]("Starknet", c)
}
