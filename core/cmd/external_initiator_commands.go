package cmd

import (
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	clipkg "github.com/urfave/cli"
)

type ExternalInitiatorPresenter struct {
	JAID
	presenters.ExternalInitiatorResource
}

func (eip *ExternalInitiatorPresenter) RenderTable(rt RendererTable) error {
	table := rt.newTable([]string{"ID", "Name", "URL", "AccessKey", "OutgoingToken", "CreatedAt", "UpdatedAt"})
	table.Append(eip.ToRow())
	render("External Initiator:", table)
	return nil
}

func (eip *ExternalInitiatorPresenter) ToRow() []string {
	var urlS string
	if eip.URL != nil {
		urlS = eip.URL.String()
	}
	return []string{
		eip.ID,
		eip.Name,
		urlS,
		eip.AccessKey,
		eip.OutgoingToken,
		eip.CreatedAt.String(),
		eip.UpdatedAt.String(),
	}
}

type ExternalInitiatorPresenters []ExternalInitiatorPresenter

func (eips *ExternalInitiatorPresenters) RenderTable(rt RendererTable) error {
	table := rt.newTable([]string{"ID", "Name", "URL", "AccessKey", "OutgoingToken", "CreatedAt", "UpdatedAt"})
	for _, eip := range *eips {
		table.Append(eip.ToRow())
	}
	render("External Initiators:", table)
	return nil
}

// IndexExternalInitiators lists external initiators
func (cli *Client) IndexExternalInitiators(c *clipkg.Context) (err error) {
	return cli.getPage("/v2/external_initiators", c.Int("page"), &ExternalInitiatorPresenters{})
}
