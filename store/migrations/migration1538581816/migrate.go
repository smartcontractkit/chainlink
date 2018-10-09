package migration1538581816

import (
	"net/url"

	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1538581816/old"
	"github.com/smartcontractkit/chainlink/store/orm"
)

type Migration struct{}

type WebURL url.URL

type TaskType string

type BridgeType struct {
	Name                   TaskType    `json:"name" storm:"id,unique"`
	URL                    WebURL      `json:"url"`
	Confirmations          uint64      `json:"confirmations"`
	IncomingToken          string      `json:"incomingToken"`
	OutgoingToken          string      `json:"outgoingToken"`
	MinimumContractPayment assets.Link `json:"minimumContractPayment"`
}

func (m Migration) Timestamp() string {
	return "1538581816"
}

func (m Migration) Migrate(orm *orm.ORM) error {
	var bridgeTypes []old.BridgeType
	if err := orm.All(&bridgeTypes); err != nil {
		return err
	}

	tx, err := orm.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, obt := range bridgeTypes {
		nbt := BridgeType{
			MinimumContractPayment: *assets.NewLink(0),
			Name:                   TaskType(obt.Name),
			URL:                    WebURL(obt.URL),
			Confirmations:          obt.Confirmations,
			IncomingToken:          obt.IncomingToken,
			OutgoingToken:          obt.OutgoingToken,
		}
		if err := tx.Save(&nbt); err != nil {
			return err
		}
	}

	return tx.Commit()
}
