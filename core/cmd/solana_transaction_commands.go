package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	solanaGo "github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/store/models/solana"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func initSolanaTxSubCmd(s *Shell) cli.Command {
	return cli.Command{
		Name:  "solana",
		Usage: "Commands for handling Solana transactions",
		Subcommands: []cli.Command{
			{
				Name:   "create",
				Usage:  "Send <amount> lamports from node Solana account <fromAddress> to destination <toAddress>.",
				Action: s.SolanaSendSol,
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "force",
						Usage: "allows to send a higher amount than the account's balance",
					},
					cli.StringFlag{
						Name:  "id",
						Usage: "chain ID, options: [mainnet, testnet, devnet, localnet]",
					},
				},
			},
		},
	}
}

type SolanaMsgPresenter struct {
	JAID
	presenters.SolanaMsgResource
}

// RenderTable implements TableRenderer
func (p *SolanaMsgPresenter) RenderTable(rt RendererTable) error {
	table := rt.newTable([]string{"Chain ID", "From", "To", "Amount"})
	table.Append([]string{
		p.ChainID,
		p.From,
		p.To,
		strconv.FormatUint(p.Amount, 10),
	})

	render(fmt.Sprintf("Solana Message %v", p.ID), table)
	return nil
}

// SolanaSendSol transfers sol from the node's account to a specified address.
func (s *Shell) SolanaSendSol(c *cli.Context) (err error) {
	if c.NArg() < 3 {
		return s.errorOut(errors.New("three arguments expected: amount, fromAddress and toAddress"))
	}

	amount, err := strconv.ParseUint(c.Args().Get(0), 10, 64)
	if err != nil {
		return s.errorOut(fmt.Errorf("invalid amount: %w", err))
	}

	unparsedFromAddress := c.Args().Get(1)
	fromAddress, err := solanaGo.PublicKeyFromBase58(unparsedFromAddress)
	if err != nil {
		return s.errorOut(multierr.Combine(
			errors.Errorf("while parsing withdrawal source address %v",
				unparsedFromAddress), err))
	}

	unparsedDestinationAddress := c.Args().Get(2)
	destinationAddress, err := solanaGo.PublicKeyFromBase58(unparsedDestinationAddress)
	if err != nil {
		return s.errorOut(multierr.Combine(
			errors.Errorf("while parsing withdrawal destination address %v",
				unparsedDestinationAddress), err))
	}

	chainID := c.String("id")
	if chainID == "" {
		return s.errorOut(errors.New("missing id"))
	}

	request := solana.SendRequest{
		To:                 destinationAddress,
		From:               fromAddress,
		Amount:             amount,
		SolanaChainID:      chainID,
		AllowHigherAmounts: c.IsSet("force"),
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		return s.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)

	resp, err := s.HTTP.Post("/v2/transfers/solana", buf)
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	err = s.renderAPIResponse(resp, &SolanaMsgPresenter{})
	return err
}
