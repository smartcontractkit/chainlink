package pipeline

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var ContractReaderNotFound = errors.New("contractReader not found")

type contractReaderManager struct {
	services.Service
	eng                    *services.Engine
	ctx                    context.Context
	relayers               map[types.RelayID]loop.Relayer
	crs                    map[string]types.ContractReader
	lggr                   logger.Logger
	heartBeatCh            chan string
	lastHeartBeatTime      map[string]time.Time
	hearthBeatTimeout      time.Duration
	heartBeatCheckInterval time.Duration

	mu sync.RWMutex
}

func newContractReaderManager(ctx context.Context, relayers map[types.RelayID]loop.Relayer, lggr logger.Logger) (*contractReaderManager, error) {
	c := contractReaderManager{
		ctx:                    ctx,
		relayers:               relayers,
		crs:                    make(map[string]types.ContractReader),
		lggr:                   lggr,
		lastHeartBeatTime:      make(map[string]time.Time),
		heartBeatCh:            make(chan string),
		heartBeatCheckInterval: time.Minute,
		hearthBeatTimeout:      time.Minute * 5,
	}
	c.Service, c.eng = services.Config{
		Name: "ContractReaderManager",
	}.NewServiceEngine(lggr)
	if err := c.Start(ctx); err != nil {
		lggr.Errorw("Failed to start contractReaderManager", "err", err)
		return nil, err
	}

	go c.checkForUnusedClients()
	return &c, nil
}

func (c *contractReaderManager) Get(relayID types.RelayID, contractAddress string, methodName string) (types.ContractReader, error) {
	id, err := createID(relayID, contractAddress, methodName)
	if err != nil {
		return nil, err
	}
	c.mu.RLock()
	csr, found := c.crs[id]
	c.mu.RUnlock()
	if !found {
		return nil, ContractReaderNotFound
	}
	c.mu.Lock()
	c.lastHeartBeatTime[id] = time.Now()
	c.mu.Unlock()
	return csr, nil
}
func (c *contractReaderManager) Create(relayID types.RelayID, contractAddress string, methodName string, config []byte) (types.ContractReader, error) {
	id, err := createID(relayID, contractAddress, methodName)
	if err != nil {
		return nil, err
	}
	c.mu.RLock()
	csr, found := c.crs[id]
	c.mu.RUnlock()
	if found {
		return nil, fmt.Errorf("contractReader already exists for network %q, chainID %q, contractAddress %q, methodName %q", relayID.Network, relayID.ChainID, contractAddress, methodName)
	}

	csr, err = c.create(relayID, contractAddress, methodName, config)
	if err != nil {
		return nil, err
	}
	//crs.Start(c.ctx)
	return csr, nil
}

func (c *contractReaderManager) checkForUnusedClients() {
	for {
		select {
		case <-c.eng.StopChan:
			c.lggr.Debug("closing contractReaderManager checkForUnusedClients loop")
			return
		case <-time.After(c.heartBeatCheckInterval):
			for id, lastSeen := range c.lastHeartBeatTime {
				diff := time.Now().Sub(lastSeen)
				if diff > c.hearthBeatTimeout {
					c.lggr.Infof("closing contractReader with ID %q", id)

					c.mu.Lock()
					_, found := c.crs[id]
					c.mu.Unlock()

					if !found {
						c.lggr.Errorf("contractReader with ID %q has timed out but cant be found in manager", id)
						continue
					}
					//c.crs[id].Close()
					c.mu.Lock()
					delete(c.crs, id)
					delete(c.lastHeartBeatTime, id)
					c.mu.Unlock()
				}
			}
		}
	}
}

func (c *contractReaderManager) create(relayID types.RelayID, contractAddress string, methodName string, config []byte) (types.ContractReader, error) {
	r, found := c.relayers[relayID]
	if !found {
		return nil, fmt.Errorf("no relayer found for network %q and chainID %q", relayID.Network, relayID.ChainID)
	}

	id, err := createID(relayID, contractAddress, methodName)
	if err != nil {
		return nil, err
	}

	csr, err := r.NewContractReader(c.ctx, config)
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	c.crs[id] = csr
	c.lastHeartBeatTime[id] = time.Now()
	c.mu.Unlock()
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

func (c *contractReaderManager) Name() string {
	return c.lggr.Name()
}

func (c *contractReaderManager) Healthy() error {
	return nil
}

func (c *contractReaderManager) HealthReport() map[string]error {
	return map[string]error{c.Name(): c.Healthy()}
}
