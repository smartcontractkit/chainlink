package cmd

import (
	"github.com/olekukonko/tablewriter"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

type ExternalInitiatorPresenter struct {
	presenters.ExternalInitiatorResource
}

func (eip ExternalInitiatorPresenter) RenderTable(rt RendererTable) error {
	table := rt.newTable([]string{"ID", "Name", "URL", "AccessKey", "OutgoingToken", "CreatedAt", "UpdatedAt"})
	eip.render(table)
	render("External Initiator:", table)
	return nil
}

func (eip ExternalInitiatorPresenter) render(table *tablewriter.Table) {
	var urlS string
	if eip.URL != nil {
		urlS = eip.URL.String()
	}
	table.Append([]string{
		eip.ID,
		eip.Name,
		urlS,
		eip.AccessKey,
		eip.OutgoingToken,
		eip.CreatedAt.String(),
		eip.UpdatedAt.String(),
	})
}

type ExternalInitiatorPresenters []ExternalInitiatorPresenter

func (eips ExternalInitiatorPresenters) RenderTable(rt RendererTable) error {
	table := rt.newTable([]string{"ID", "Name", "URL", "AccessKey", "OutgoingToken", "CreatedAt", "UpdatedAt"})
	for _, eip := range eips {
		eip.render(table)
	}
	render("External Initiators:", table)
	return nil
}
