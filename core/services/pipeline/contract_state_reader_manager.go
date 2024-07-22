package pipeline

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var CSRNotFoundErr = errors.New("contractStateReader not found")

type contractStateReaderManager struct {
	ctx               context.Context
	relayers          map[types.RelayID]loop.Relayer
	csr               map[string]types.ContractStateReader
	lggr              logger.Logger
	heartBeatCh       chan string
	lastHeartBeatTime map[string]time.Time
}

func newContractStateReaderManager(ctx context.Context, relayers map[types.RelayID]loop.Relayer, lggr logger.Logger) contractStateReaderManager {
	c := contractStateReaderManager{
		ctx:               ctx,
		relayers:          relayers,
		csr:               make(map[string]types.ContractStateReader),
		lggr:              lggr,
		lastHeartBeatTime: make(map[string]time.Time),
		heartBeatCh:       make(chan string),
	}
	c.checkForUnusedClients(ctx)
	return c
}

func (c *contractStateReaderManager) Get(relayID types.RelayID, contractAddress string, methodName string) (types.ContractStateReader, error) {
	id, err := createID(relayID, contractAddress, methodName)
	if err != nil {
		return nil, err
	}
	csr, found := c.csr[id]
	if !found {
		return nil, CSRNotFoundErr
	}
	c.lastHeartBeatTime[id] = time.Now()
	return csr, nil
}
func (c *contractStateReaderManager) Create(relayID types.RelayID, contractAddress string, methodName string, config []byte) (types.ContractStateReader, error) {
	id, err := createID(relayID, contractAddress, methodName)
	if err != nil {
		return nil, err
	}
	csr, found := c.csr[id]
	if found {
		return nil, fmt.Errorf("contractStateReader already exists for network %q, chainID %q, contractAddress %q, methodName %q", relayID.Network, relayID.ChainID, contractAddress, methodName)
	}

	csr, err = c.create(relayID, contractAddress, methodName, config)
	if err != nil {
		return nil, err
	}
	csr.Start(c.ctx)
	return csr, nil
}

func (c *contractStateReaderManager) checkForUnusedClients(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Minute):
			for id, lastSeen := range c.lastHeartBeatTime {
				diff := time.Now().Sub(lastSeen)
				if diff > (time.Minute * 5) {
					c.lggr.Infof("closing contractStateReader with ID %q", id)
					_, found := c.csr[id]
					if !found {
						c.lggr.Errorf("contractStateReader with ID %q has timed out but cant be found in manager", id)
						continue
					}
					c.csr[id].Close()
					delete(c.csr, id)
					delete(c.lastHeartBeatTime, id)
				}
			}
		}
	}
}

func (c *contractStateReaderManager) create(relayID types.RelayID, contractAddress string, methodName string, config []byte) (types.ContractStateReader, error) {
	r, found := c.relayers[relayID]
	if !found {
		return nil, fmt.Errorf("no relayer found for network %q and chainID %q", relayID.Network, relayID.ChainID)
	}

	id, err := createID(relayID, contractAddress, methodName)
	if err != nil {
		return nil, err
	}

	csr, err := r.NewContractStateReader(c.ctx, config)
	if err != nil {
		return nil, err
	}

	c.csr[id] = csr
	return csr, nil
}

func createID(relayID types.RelayID, contractAddress string, methodName string) (string, error) {
	if relayID.ChainID == "" {
		return "", errors.New("cannot create ID, chainID is empty")
	}

	if relayID.Network == "" {
		return "", errors.New("cannot create ID, network is empty")
	}

	if contractAddress == "" {
		return "", errors.New("cannot create ID, contractAddress is empty")
	}

	if methodName == "" {
		return "", errors.New("cannot create ID, methodName is empty")
	}

	return fmt.Sprintf("%s_%s_%s_%s", relayID.Network, relayID.ChainID, contractAddress, methodName), nil
}

/*TODO @george-dorin:
- cleanup
- add context
*/
