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

func initCosmosTxSubCmd(s *Shell) cli.Command {
	return cli.Command{
		Name:  "cosmos",
		Usage: "Commands for handling Cosmos transactions",
		Subcommands: []cli.Command{
			{
				Name:   "create",
				Usage:  "Send <amount> of <token> from node Cosmos account <fromAddress> to destination <toAddress>.",
				Action: s.CosmosSendNativeToken,
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

// CosmosSendNativeToken transfers coins from the node's account to a specified address.
func (s *Shell) CosmosSendNativeToken(c *cli.Context) (err error) {
	if c.NArg() < 3 {
		return s.errorOut(errors.New("four arguments expected: token, amount, fromAddress and toAddress"))
	}

	err = sdk.ValidateDenom(c.Args().Get(0))
	if err != nil {
		return s.errorOut(fmt.Errorf("invalid native token: %w", err))
	}

	amount, err := sdk.NewDecFromStr(c.Args().Get(1))
	if err != nil {
		return s.errorOut(multierr.Combine(
			fmt.Errorf("invalid coin: %w", err)))
	}

	unparsedFromAddress := c.Args().Get(2)
	fromAddress, err := sdk.AccAddressFromBech32(unparsedFromAddress)
	if err != nil {
		return s.errorOut(multierr.Combine(
			fmt.Errorf("while parsing withdrawal source address %v",
				unparsedFromAddress), err))
	}

	unparsedDestinationAddress := c.Args().Get(3)
	destinationAddress, err := sdk.AccAddressFromBech32(unparsedDestinationAddress)
	if err != nil {
		return s.errorOut(multierr.Combine(
			fmt.Errorf("while parsing withdrawal destination address %v",
				unparsedDestinationAddress), err))
	}

	chainID := c.String("id")
	if chainID == "" {
		return s.errorOut(errors.New("missing id"))
	}

	request := cosmos.SendRequest{
		DestinationAddress: destinationAddress,
		FromAddress:        fromAddress,
		Amount:             amount,
		CosmosChainID:      chainID,
		Token:              c.Args().Get(0),
		AllowHigherAmounts: c.IsSet("force"),
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		return s.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)

	resp, err := s.HTTP.Post(s.ctx(), "/v2/transfers/cosmos", buf)
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	err = s.renderAPIResponse(resp, &CosmosMsgPresenter{})
	return err
}
