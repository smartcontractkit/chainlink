package cmd

import (
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func initInitiatorsSubCmds(s *Shell) []cli.Command {
	return []cli.Command{
		{
			Name:   "create",
			Usage:  "Create an authentication key for a user of External Initiators",
			Action: s.CreateExternalInitiator,
		},
		{
			Name:   "destroy",
			Usage:  "Remove an external initiator by name",
			Action: s.DeleteExternalInitiator,
		},
		{
			Name:   "list",
			Usage:  "List all external initiators",
			Action: s.IndexExternalInitiators,
		},
	}
}

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
func (s *Shell) IndexExternalInitiators(c *cli.Context) (err error) {
	return s.getPage("/v2/external_initiators", c.Int("page"), &ExternalInitiatorPresenters{})
}
