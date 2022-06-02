package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/urfave/cli"
	"go.uber.org/multierr"
)

type OperatorPresenter struct {
	JAID // This is needed to render the id for a JSONAPI Resource as normal JSON
	presenters.OperatorResource
}

var operatorHeaders = []string{"ID", "Address", "Chain ID", "Created At", "Updated At"}

// ToRow presents the EVMForwarderResource as a slice of strings.
func (p *OperatorPresenter) ToRow() []string {
	row := []string{
		p.GetID(),
		p.Address.String(),
		p.ChainID.String(),
		p.CreatedAt.Format(time.RFC3339),
		p.UpdatedAt.Format(time.RFC3339),
	}
	return row
}

// RenderTable implements TableRenderer
func (p *OperatorPresenter) RenderTable(rt RendererTable) error {
	var rows [][]string
	rows = append(rows, p.ToRow())
	renderList(operatorHeaders, rows, rt.Writer)

	return nil
}

// OperatorPresenters implements TableRenderer for a slice of OperatorPresenter.
type OperatorPresenters []OperatorPresenter

// RenderTable implements TableRenderer
func (ps OperatorPresenters) RenderTable(rt RendererTable) error {
	var rows [][]string

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(operatorHeaders, rows, rt.Writer)

	return nil
}

// // ListForwarders list all forwarder addresses tracked by node
// func (cli *Client) ListForwarders(c *cli.Context) (err error) {
// 	return cli.getPage("/v2/nodes/evm/forwarders", c.Int("page"), &EVMForwarderPresenters{})
// }

// // DeleteForwarder deletes forwarder address from node db by id.
// func (cli *Client) DeleteForwarder(c *cli.Context) (err error) {
// 	if !c.Args().Present() {
// 		return cli.errorOut(errors.New("must pass the forwarder id to be archived"))
// 	}
// 	resp, err := cli.HTTP.Delete("/v2/nodes/evm/forwarders/" + c.Args().First())
// 	if err != nil {
// 		return cli.errorOut(err)
// 	}
// 	_, err = cli.parseResponse(resp)
// 	if err != nil {
// 		return cli.errorOut(err)
// 	}

// 	fmt.Printf("Forwarder %v Deleted\n", c.Args().First())
// 	return nil
// }

func (cli *Client) DeployOperator(c *cli.Context) (err error) {
	addressHex := c.String("owner")
	chainIDStr := c.String("chainId")

	addressBytes, err := hexutil.Decode(addressHex)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "could not decode address"))
	}
	address := gethCommon.BytesToAddress(addressBytes)

	var chainID *big.Int
	if chainIDStr != "" {
		var ok bool
		chainID, ok = big.NewInt(0).SetString(chainIDStr, 10)
		if !ok {
			return cli.errorOut(errors.Wrap(err, "invalid evmChainID"))
		}
	}

	request, err := json.Marshal(web.DeployOperatorRequest{
		ChainID: (*utils.Big)(chainID),
		Owner:   address,
	})
	if err != nil {
		return cli.errorOut(err)
	}

	dresp, err := cli.HTTP.Post("/v2/operators/deploy", bytes.NewReader(request))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := dresp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if dresp.StatusCode >= 400 {
		body, rerr := ioutil.ReadAll(dresp.Body)
		if err != nil {
			err = multierr.Append(err, rerr)
			return cli.errorOut(err)
		}
		fmt.Printf("Response: '%v', Status: %d\n", string(body), dresp.StatusCode)
		return cli.errorOut(err)
	}

	// err = cli.renderAPIResponse(dresp, &EthTxPresenter{}, "Tx created")
	// if err != nil {
	// 	return err
	// }

	var links jsonapi.Links
	ethtxId := JAID{}
	if err := cli.deserializeAPIResponse(dresp, &ethtxId, &links); err != nil {
		return cli.errorOut(err)
	}
	// err = cli.Render(&ethtxp, "Tx created")
	// if err != nil {
	// 	return cli.errorOut(err)
	// }
	cli.Logger.Infof("Tx ID: %s", ethtxId.GetID())

	pollInterval := 6500 * time.Millisecond
	ticker := time.NewTicker(utils.WithJitter(pollInterval))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cli.Logger.Infof("Polling for operator deployment status...")
			resp, err := cli.HTTP.Get(fmt.Sprintf("/v2/operators/%s/%s/status", chainIDStr, ethtxId.GetID()))
			if err != nil {
				continue
			}
			defer func() {
				if cerr := resp.Body.Close(); cerr != nil {
					err = multierr.Append(err, cerr)
				}
			}()
			if resp.StatusCode == 201 {
				err := cli.renderAPIResponse(resp, &OperatorPresenter{}, "Operator Created")
				return err
			}
			if resp.StatusCode == 417 {
				cli.Logger.Infof("Tx %s attempt failed", ethtxId.GetID())
				return err
			}
			if resp.StatusCode == 404 {
				cli.Logger.Infof("Tx %s not on chain yet", ethtxId.GetID())
				continue
			} else if resp.StatusCode >= 400 {
				cli.Logger.Infof("Request status: %d %s", resp.StatusCode)
				continue
			}
		case <-time.After(2 * time.Minute):
			return cli.errorOut(errors.Errorf("timeout polling for tx %s", ethtxId.GetID()))
		}
	}
}
