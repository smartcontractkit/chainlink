package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

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
func (cli *Client) IndexTransactions(c *cli.Context) error {
	return cli.getPage("/v2/transactions/evm", c.Int("page"), &EthTxPresenters{})
}

// ShowTransaction returns the info for the given transaction hash
func (cli *Client) ShowTransaction(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass the hash of the transaction"))
	}
	hash := c.Args().First()
	resp, err := cli.HTTP.Get("/v2/transactions/evm/" + hash)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	err = cli.renderAPIResponse(resp, &EthTxPresenter{})
	return err
}

// IndexTxAttempts returns the list of transactions in descending order,
// taking an optional page parameter
func (cli *Client) IndexTxAttempts(c *cli.Context) error {
	return cli.getPage("/v2/tx_attempts/evm", c.Int("page"), &EthTxPresenters{})
}

// SendEther transfers ETH from the node's account to a specified address.
func (cli *Client) SendEther(c *cli.Context) (err error) {
	if c.NArg() < 3 {
		return cli.errorOut(errors.New("three arguments expected: amount, fromAddress and toAddress"))
	}

	var amount assets.Eth

	if c.IsSet("wei") {
		var value int64

		value, err = stringutils.ToInt64(c.Args().Get(0))
		if err != nil {
			return cli.errorOut(multierr.Combine(
				errors.New("while parsing WEI transfer amount"), err))
		}

		amount = assets.NewEthValue(value)
	} else {
		amount, err = assets.NewEthValueS(c.Args().Get(0))
		if err != nil {
			return cli.errorOut(multierr.Combine(
				errors.New("while parsing ETH transfer amount"), err))
		}
	}

	unparsedFromAddress := c.Args().Get(1)
	fromAddress, err := utils.ParseEthereumAddress(unparsedFromAddress)
	if err != nil {
		return cli.errorOut(multierr.Combine(
			fmt.Errorf("while parsing withdrawal source address %v",
				unparsedFromAddress), err))
	}

	unparsedDestinationAddress := c.Args().Get(2)
	destinationAddress, err := utils.ParseEthereumAddress(unparsedDestinationAddress)
	if err != nil {
		return cli.errorOut(multierr.Combine(
			fmt.Errorf("while parsing withdrawal destination address %v",
				unparsedDestinationAddress), err))
	}

	var evmChainID *big.Int
	if c.IsSet("id") {
		s := c.String("id")
		var ok bool
		evmChainID, ok = new(big.Int).SetString(s, 10)
		if !ok {
			return cli.errorOut(errors.New(""))
		}
	}

	request := models.SendEtherRequest{
		DestinationAddress: destinationAddress,
		FromAddress:        fromAddress,
		Amount:             amount,
		EVMChainID:         (*utils.Big)(evmChainID),
		AllowHigherAmounts: c.IsSet("force"),
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		return cli.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)

	resp, err := cli.HTTP.Post("/v2/transfers/evm", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	err = cli.renderAPIResponse(resp, &EthTxPresenter{})
	return err
}
