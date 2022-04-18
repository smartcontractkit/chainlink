package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/store/models/terra"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

type TerraMsgPresenter struct {
	JAID
	presenters.TerraMsgResource
}

// RenderTable implements TableRenderer
func (p *TerraMsgPresenter) RenderTable(rt RendererTable) error {
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

	render(fmt.Sprintf("Terra Message %v", p.ID), table)
	return nil
}

// TerraSendLuna transfers coins from the node's account to a specified address.
func (cli *Client) TerraSendLuna(c *cli.Context) (err error) {
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

	request := terra.SendRequest{
		DestinationAddress: destinationAddress,
		FromAddress:        fromAddress,
		Amount:             amount,
		TerraChainID:       chainID,
		AllowHigherAmounts: c.IsSet("force"),
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		return cli.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)

	resp, err := cli.HTTP.Post("/v2/transfers/terra", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	err = cli.renderAPIResponse(resp, &TerraMsgPresenter{})
	return err
}
