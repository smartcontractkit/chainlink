package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils/stringutils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func initEVMTxSubCmd(s *Shell) cli.Command {
	return cli.Command{
		Name:  "evm",
		Usage: "Commands for handling EVM transactions",
		Subcommands: []cli.Command{
			{
				Name:   "create",
				Usage:  "Send <amount> ETH (or wei) from node ETH account <fromAddress> to destination <toAddress>.",
				Action: s.SendEther,
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "force",
						Usage: "allows to send a higher amount than the account's balance",
					},
					cli.BoolFlag{
						Name:  "eth",
						Usage: "allows to send ETH amounts (Default behavior)",
					},
					cli.BoolFlag{
						Name:  "wei",
						Usage: "allows to send WEI amounts",
					},
					cli.Int64Flag{
						Name:  "id",
						Usage: "chain ID",
					},
				},
			},
			{
				Name:   "list",
				Usage:  "List the Ethereum Transactions in descending order",
				Action: s.IndexTransactions,
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "page",
						Usage: "page of results to display",
					},
				},
			},
			{
				Name:   "show",
				Usage:  "get information on a specific Ethereum Transaction",
				Action: s.ShowTransaction,
			},
		},
	}
}

type EthTxPresenter struct {
	JAID
	presenters.EthTxResource
}

// RenderTable implements TableRenderer
func (p *EthTxPresenter) RenderTable(rt RendererTable) error {
	table := rt.newTable([]string{"From", "Nonce", "To", "State"})
	table.Append([]string{
		p.From.Hex(),
		p.Nonce,
		p.To.Hex(),
		fmt.Sprint(p.State),
	})

	render(fmt.Sprintf("Ethereum Transaction %v", p.Hash.Hex()), table)
	return nil
}

type EthTxPresenters []EthTxPresenter

// RenderTable implements TableRenderer
func (ps EthTxPresenters) RenderTable(rt RendererTable) error {
	table := rt.newTable([]string{"Hash", "Nonce", "From", "GasPrice", "SentAt", "State"})
	for _, p := range ps {
		table.Append([]string{
			p.Hash.Hex(),
			p.Nonce,
			p.From.Hex(),
			p.GasPrice,
			p.SentAt,
			fmt.Sprint(p.State),
		})
	}

	render("Ethereum Transactions", table)
	return nil
}

// IndexTransactions returns the list of transactions in descending order,
// taking an optional page parameter
func (s *Shell) IndexTransactions(c *cli.Context) error {
	return s.getPage("/v2/transactions/evm", c.Int("page"), &EthTxPresenters{})
}

// ShowTransaction returns the info for the given transaction hash
func (s *Shell) ShowTransaction(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return s.errorOut(errors.New("must pass the hash of the transaction"))
	}
	hash := c.Args().First()
	resp, err := s.HTTP.Get(s.ctx(), "/v2/transactions/evm/"+hash)
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	err = s.renderAPIResponse(resp, &EthTxPresenter{})
	return err
}

// SendEther transfers ETH from the node's account to a specified address.
func (s *Shell) SendEther(c *cli.Context) (err error) {
	if c.NArg() < 3 {
		return s.errorOut(errors.New("the following arguments expected: (chain) id (in multi-chain setup), amount, fromAddress and toAddress"))
	}

	var amount assets.Eth

	if c.IsSet("wei") {
		var value int64

		value, err = stringutils.ToInt64(c.Args().Get(0))
		if err != nil {
			return s.errorOut(multierr.Combine(
				errors.New("while parsing WEI transfer amount"), err))
		}

		amount = assets.NewEthValue(value)
	} else {
		amount, err = assets.NewEthValueS(c.Args().Get(0))
		if err != nil {
			return s.errorOut(multierr.Combine(
				errors.New("while parsing ETH transfer amount"), err))
		}
	}

	unparsedFromAddress := c.Args().Get(1)
	fromAddress, err := utils.ParseEthereumAddress(unparsedFromAddress)
	if err != nil {
		return s.errorOut(multierr.Combine(
			fmt.Errorf("while parsing withdrawal source address %v",
				unparsedFromAddress), err))
	}

	unparsedDestinationAddress := c.Args().Get(2)
	destinationAddress, err := utils.ParseEthereumAddress(unparsedDestinationAddress)
	if err != nil {
		return s.errorOut(multierr.Combine(
			fmt.Errorf("while parsing withdrawal destination address %v",
				unparsedDestinationAddress), err))
	}

	var evmChainID *big.Int
	if c.IsSet("id") {
		str := c.String("id")
		var ok bool
		evmChainID, ok = new(big.Int).SetString(str, 10)
		if !ok {
			return s.errorOut(errors.New(""))
		}
	}

	request := models.SendEtherRequest{
		DestinationAddress: destinationAddress,
		FromAddress:        fromAddress,
		Amount:             amount,
		EVMChainID:         (*ubig.Big)(evmChainID),
		AllowHigherAmounts: c.IsSet("force"),
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		return s.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)

	resp, err := s.HTTP.Post(s.ctx(), "/v2/transfers/evm", buf)
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	err = s.renderAPIResponse(resp, &EthTxPresenter{})
	return err
}
