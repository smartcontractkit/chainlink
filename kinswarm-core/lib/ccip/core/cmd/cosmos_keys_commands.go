package cmd

import (
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/cosmoskey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

type CosmosKeyPresenter struct {
	JAID
	presenters.CosmosKeyResource
}

// RenderTable implements TableRenderer
func (p CosmosKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 Cosmos Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func (p *CosmosKeyPresenter) ToRow() []string {
	row := []string{
		p.ID,
		p.PubKey,
	}

	return row
}

type CosmosKeyPresenters []CosmosKeyPresenter

// RenderTable implements TableRenderer
func (ps CosmosKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 Cosmos Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func NewCosmosKeysClient(s *Shell) KeysClient {
	return newKeysClient[cosmoskey.Key, CosmosKeyPresenter, CosmosKeyPresenters]("Cosmos", s)
}
