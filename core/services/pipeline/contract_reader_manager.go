package pipeline

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type contractReaderManager struct {
	services.Service
	eng                    *services.Engine
	stopCh                 services.StopChan
	relayers               map[types.RelayID]loop.Relayer
	crs                    map[string]*contractReaderWithIdentifier
	lggr                   logger.Logger
	heartBeatCh            chan string
	lastHeartBeatTime      map[string]time.Time
	hearthBeatTimeout      time.Duration
	heartBeatCheckInterval time.Duration

	mu sync.RWMutex
}

type contractReaderWithIdentifier struct {
	cr  types.ContractReader
	rID string
}

func newContractReaderManager(relayers map[types.RelayID]loop.Relayer, stopCh services.StopChan, lggr logger.Logger) (*contractReaderManager, error) {
	c := contractReaderManager{
		stopCh:                 stopCh,
		relayers:               relayers,
		crs:                    make(map[string]*contractReaderWithIdentifier),
		lggr:                   lggr,
		lastHeartBeatTime:      make(map[string]time.Time),
		heartBeatCh:            make(chan string),
		heartBeatCheckInterval: time.Minute,
		hearthBeatTimeout:      time.Minute * 5,
	}
	c.Service, c.eng = services.Config{
		Name: "ContractReaderManager",
	}.NewServiceEngine(lggr)
	ctx, _ := stopCh.NewCtx()
	if err := c.Start(ctx); err != nil {
		lggr.Errorw("Failed to start contractReaderManager", "err", err)
		return nil, err
	}

	go c.checkForUnusedClients()
	return &c, nil
}

func (c *contractReaderManager) GetOrCreate(relayID types.RelayID, contractName string, contractAddress string, methodName string, config []byte) (reader types.ContractReader, identifier string, err error) {
	id, err := createID(relayID, contractAddress, methodName)
	if err != nil {
		return nil, "", err
	}
	c.mu.RLock()
	csr, found := c.crs[id]
	c.mu.RUnlock()
	if !found {
		csr, err = c.create(relayID, contractName, contractAddress, methodName, config)
		if err != nil {
			return nil, "", err
		}
	}
	c.mu.Lock()
	c.lastHeartBeatTime[id] = time.Now()
	c.mu.Unlock()
	return csr.cr, csr.rID, nil
}

func (c *contractReaderManager) checkForUnusedClients() {
	for {
		select {
		case <-c.eng.StopChan:
			c.lggr.Debug("closing contractReaderManager checkForUnusedClients loop")
			return
		case <-time.After(c.heartBeatCheckInterval):
			for id, lastSeen := range c.lastHeartBeatTime {
				diff := time.Since(lastSeen)
				if diff > c.hearthBeatTimeout {
					c.lggr.Infof("closing contractReader with ID %q", id)

					c.mu.Lock()
					_, found := c.crs[id]
					c.mu.Unlock()

					if !found {
						c.lggr.Errorf("contractReader with ID %q has timed out but cant be found in manager", id)
						continue
					}
					c.mu.Lock()
					delete(c.crs, id)
					delete(c.lastHeartBeatTime, id)
					c.mu.Unlock()
				}
			}
		}
	}
}

func (c *contractReaderManager) create(relayID types.RelayID, contractName string, contractAddress string, methodName string, config []byte) (*contractReaderWithIdentifier, error) {
	r, found := c.relayers[relayID]
	if !found {
		return nil, fmt.Errorf("no relayer found for network %q and chainID %q", relayID.Network, relayID.ChainID)
	}

	id, err := createID(relayID, contractAddress, methodName)
	if err != nil {
		return nil, err
	}

	ctx, _ := c.stopCh.NewCtx()
	csr, err := r.NewContractReader(ctx, config)
	if err != nil {
		return nil, err
	}

	cb := types.BoundContract{
		Address: contractAddress,
		Name:    contractName,
	}
	err = csr.Bind(ctx, []types.BoundContract{cb})
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	cri := contractReaderWithIdentifier{
		cr:  csr,
		rID: cb.ReadIdentifier(methodName),
	}
	c.crs[id] = &cri
	c.lastHeartBeatTime[id] = time.Now()
	c.mu.Unlock()
	return &cri, nil
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
