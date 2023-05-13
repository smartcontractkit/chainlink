package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/store/models/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func initCosmosTxSubCmd(client *Client) cli.Command {
	return cli.Command{
		Name:  "cosmos",
		Usage: "Commands for handling Cosmos transactions",
		Subcommands: []cli.Command{
			{
				Name:   "create",
				Usage:  "Send <amount> Atom from node Cosmos account <fromAddress> to destination <toAddress>.",
				Action: client.CosmosSendAtom,
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "force",
						Usage: "allows to send a higher amount than the account's balance",
					},
					cli.StringFlag{
						Name:  "id",
						Usage: "chain ID",
					},
				},
			},
		},
	}
}

type CosmosMsgPresenter struct {
	JAID
	presenters.CosmosMsgResource
}

// RenderTable implements TableRenderer
func (p *CosmosMsgPresenter) RenderTable(rt RendererTable) error {
	table := rt.newTable([]string{"Chain ID", "Contract ID", "State", "Tx Hash"})
	var hash string
	if p.TxHash != nil {
		hash = *p.TxHash
	}
	table.Append([]string{
		p.ChainID,
		p.ContractID,
		p.State,
		hash,
	})

	render(fmt.Sprintf("Cosmos Message %v", p.ID), table)
	return nil
}

// CosmosSendAtom transfers coins from the node's account to a specified address.
func (cli *Client) CosmosSendAtom(c *cli.Context) (err error) {
	if c.NArg() < 3 {
		return cli.errorOut(errors.New("three arguments expected: amount, fromAddress and toAddress"))
	}

	amount, err := sdk.NewDecFromStr(c.Args().Get(0))
	if err != nil {
		return cli.errorOut(fmt.Errorf("invalid coin: %w", err))
	}

	unparsedFromAddress := c.Args().Get(1)
	fromAddress, err := sdk.AccAddressFromBech32(unparsedFromAddress)
	if err != nil {
		return cli.errorOut(multierr.Combine(
			fmt.Errorf("while parsing withdrawal source address %v",
				unparsedFromAddress), err))
	}

	unparsedDestinationAddress := c.Args().Get(2)
	destinationAddress, err := sdk.AccAddressFromBech32(unparsedDestinationAddress)
	if err != nil {
		return cli.errorOut(multierr.Combine(
			fmt.Errorf("while parsing withdrawal destination address %v",
				unparsedDestinationAddress), err))
	}

	chainID := c.String("id")
	if chainID == "" {
		return cli.errorOut(errors.New("missing id"))
	}

	request := cosmos.SendRequest{
		DestinationAddress: destinationAddress,
		FromAddress:        fromAddress,
		Amount:             amount,
		CosmosChainID:      chainID,
		AllowHigherAmounts: c.IsSet("force"),
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		return cli.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)

	resp, err := cli.HTTP.Post("/v2/transfers/cosmos", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	err = cli.renderAPIResponse(resp, &CosmosMsgPresenter{})
	return err
}
