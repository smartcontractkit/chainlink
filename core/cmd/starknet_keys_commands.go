package cmd

import (
	starkkey "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/keys"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

type StarkNetKeyPresenter struct {
	JAID
	presenters.StarkNetKeyResource
}

// RenderTable implements TableRenderer
func (p StarkNetKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key"}
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
		p.PubKey,
	}

	return row
}

type StarkNetKeyPresenters []StarkNetKeyPresenter

// RenderTable implements TableRenderer
func (ps StarkNetKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key"}
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

func NewStarkNetKeysClient(c *Client) KeysClient {
	return newKeysClient[starkkey.StarkKey, StarkNetKeyPresenter, StarkNetKeyPresenters]("StarkNet", c)
}
