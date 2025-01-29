package load

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/common/model"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testreporters"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
)

// CCIPMultiCallLoadGenerator represents a load generator for the CCIP lanes originating from same network
// The purpose of this load generator is to group ccip-send calls for the CCIP lanes originating from same network
// This is to avoid the scenario of hitting rpc rate limit for the same network if the load generator is sending
// too many ccip-send calls to the same network hitting the rpc rate limit
type CCIPMultiCallLoadGenerator struct {
	t                       *testing.T
	logger                  zerolog.Logger
	client                  blockchain.EVMClient
	E2ELoads                map[string]*CCIPE2ELoad
	MultiCall               string
	NoOfRequestsPerUnitTime int64
	labels                  model.LabelSet
	loki                    *wasp.LokiClient
	responses               chan map[string]MultiCallReturnValues
	Done                    chan struct{}
}

type MultiCallReturnValues struct {
	Msgs  []contracts.CCIPMsgData
	Stats []*testreporters.RequestStat
}

func NewMultiCallLoadGenerator(testCfg *testsetups.CCIPTestConfig, lanes []*actions.CCIPLane, noOfRequestsPerUnitTime int64, labels map[string]string) (*CCIPMultiCallLoadGenerator, error) {
	// check if all lanes are from same network
	source := lanes[0].Source.Common.ChainClient.GetChainID()
	multiCall := lanes[0].Source.Common.MulticallContract.Hex()
	if multiCall == "" {
		return nil, fmt.Errorf("multicall address cannot be empty")
	}
	for i := 1; i < len(lanes); i++ {
		if source.String() != lanes[i].Source.Common.ChainClient.GetChainID().String() {
			return nil, fmt.Errorf("all lanes should be from same network; expected %s, got %s", source, lanes[i].Source.Common.ChainClient.GetChainID())
		}
		if lanes[i].Source.Common.MulticallContract.Hex() != multiCall {
			return nil, fmt.Errorf("multicall address should be same for all lanes")
		}
	}
	client := lanes[0].Source.Common.ChainClient
	lggr := logging.GetTestLogger(testCfg.Test).With().Str("Source Network", client.GetNetworkName()).Logger()
	ls := wasp.LabelsMapToModel(labels)
	if err := ls.Validate(); err != nil {
		return nil, err
	}
	lokiConfig := testCfg.EnvInput.Logging.Loki
	loki, err := wasp.NewLokiClient(wasp.NewLokiConfig(lokiConfig.Endpoint, lokiConfig.TenantId, nil, nil))
	if err != nil {
		return nil, err
	}
	m := &CCIPMultiCallLoadGenerator{
		t:                       testCfg.Test,
		client:                  client,
		MultiCall:               multiCall,
		logger:                  lggr,
		NoOfRequestsPerUnitTime: noOfRequestsPerUnitTime,
		E2ELoads:                make(map[string]*CCIPE2ELoad),
		labels:                  ls,
		loki:                    loki,
		responses:               make(chan map[string]MultiCallReturnValues),
		Done:                    make(chan struct{}),
	}
	for _, lane := range lanes {
		// for multicall load generator, we don't want to send max data intermittently, it might
		// cause oversized data for multicall
		ccipLoad := NewCCIPLoad(
			testCfg.Test, lane, testCfg.TestGroupInput.PhaseTimeout.Duration(),
			100000,
			testCfg.TestGroupInput.LoadProfile.MsgProfile, 0,
			testCfg.TestGroupInput.SkipRequestIfAnotherRequestTriggeredWithin,
		)
		ccipLoad.BeforeAllCall()
		m.E2ELoads[fmt.Sprintf("%s-%s", lane.SourceNetworkName, lane.DestNetworkName)] = ccipLoad
	}

	m.StartLokiStream()
	return m, nil
}

func (m *CCIPMultiCallLoadGenerator) Stop() error {
	m.Done <- struct{}{}
	tokenMap := make(map[string]struct{})
	var tokens []*contracts.ERC20Token
	for _, e2eLoad := range m.E2ELoads {
		for i := range e2eLoad.Lane.Source.TransferAmount {
			// if length of sourceCCIP.TransferAmount is more than available bridge token use first bridge token
			token := e2eLoad.Lane.Source.Common.BridgeTokens[0]
			if i < len(e2eLoad.Lane.Source.Common.BridgeTokens) {
				token = e2eLoad.Lane.Source.Common.BridgeTokens[i]
			}
			if _, ok := tokenMap[token.Address()]; !ok {
				tokens = append(tokens, e2eLoad.Lane.Source.Common.BridgeTokens[i])
			}
		}
	}
	if len(tokens) > 0 {
		return contracts.TransferTokens(m.client, common.HexToAddress(m.MultiCall), tokens)
	}
	return nil
}

func (m *CCIPMultiCallLoadGenerator) StartLokiStream() {
	go func() {
		for {
			select {
			case <-m.Done:
				m.logger.Info().Msg("stopping loki client from multi call load generator")
				m.loki.Stop()
				return
			case rValues := <-m.responses:
				m.HandleLokiLogs(rValues)
			}
		}
	}()
}

func (m *CCIPMultiCallLoadGenerator) HandleLokiLogs(rValues map[string]MultiCallReturnValues) {
	for dest, rValue := range rValues {
		labels := m.labels.Merge(model.LabelSet{
			"dest_chain":     model.LabelValue(dest),
			"test_data_type": "responses",
			"go_test_name":   model.LabelValue(m.t.Name()),
		})
		for _, stat := range rValue.Stats {
			err := m.loki.HandleStruct(labels, time.Now().UTC(), stat.StatusByPhase)
			if err != nil {
				m.logger.Error().Err(err).Msg("error while handling loki logs")
			}
		}
	}
}

func (m *CCIPMultiCallLoadGenerator) Call(_ *wasp.Generator) *wasp.Response {
	res := &wasp.Response{}
	msgs, returnValuesByDest, err := m.MergeCalls()
	if err != nil {
		res.Error = err.Error()
		res.Failed = true
		return res
	}
	defer func() {
		m.responses <- returnValuesByDest
	}()
	m.logger.Info().Interface("msgs", msgs).Msgf("Sending %d ccip-send calls", len(msgs))
	startTime := time.Now().UTC()
	// for now we are using all ccip-sends with native
	sendTx, err := contracts.MultiCallCCIP(m.client, m.MultiCall, msgs, true)
	if err != nil {
		res.Error = err.Error()
		res.Failed = true
		return res
	}

	lggr := m.logger.With().Str("Msg Tx", sendTx.Hash().String()).Logger()
	txConfirmationTime := time.Now().UTC()
	for _, rValues := range returnValuesByDest {
		if len(rValues.Stats) != len(rValues.Msgs) {
			res.Error = fmt.Sprintf("number of stats %d and msgs %d should be same", len(rValues.Stats), len(rValues.Msgs))
			res.Failed = true
			return res
		}
		for _, stat := range rValues.Stats {
			stat.UpdateState(&lggr, 0, testreporters.TX, startTime.Sub(txConfirmationTime), testreporters.Success, nil)
		}
	}

	validateGrp := errgroup.Group{}
	// wait for
	// - CCIPSendRequested Event log to be generated,
	for _, rValues := range returnValuesByDest {
		key := fmt.Sprintf("%s-%s", rValues.Stats[0].SourceNetwork, rValues.Stats[0].DestNetwork)
		c, ok := m.E2ELoads[key]
		if !ok {
			res.Error = fmt.Sprintf("load for %s not found", key)
			res.Failed = true
			return res
		}

		lggr = lggr.With().Str("Source Network", c.Lane.Source.Common.ChainClient.GetNetworkName()).Str("Dest Network", c.Lane.Dest.Common.ChainClient.GetNetworkName()).Logger()
		stats := rValues.Stats
		txConfirmationTime := txConfirmationTime
		sendTx := sendTx
		lggr := lggr
		validateGrp.Go(func() error {
			return c.Validate(lggr, sendTx, txConfirmationTime, stats)
		})
	}
	err = validateGrp.Wait()
	if err != nil {
		res.Error = err.Error()
		res.Failed = true
		return res
	}

	return res
}

func (m *CCIPMultiCallLoadGenerator) MergeCalls() ([]contracts.CCIPMsgData, map[string]MultiCallReturnValues, error) {
	var ccipMsgs []contracts.CCIPMsgData
	statDetails := make(map[string]MultiCallReturnValues)

	for _, e2eLoad := range m.E2ELoads {
		destChainSelector, err := chain_selectors.SelectorFromChainId(e2eLoad.Lane.Source.DestinationChainId)
		if err != nil {
			return ccipMsgs, statDetails, err
		}

		allFee := big.NewInt(0)
		var allStatsForDest []*testreporters.RequestStat
		var allMsgsForDest []contracts.CCIPMsgData
		for i := int64(0); i < m.NoOfRequestsPerUnitTime; i++ {
			msg, stats, err := e2eLoad.CCIPMsg()
			if err != nil {
				return ccipMsgs, statDetails, err
			}
			msg.FeeToken = common.Address{}
			fee, err := e2eLoad.Lane.Source.Common.Router.GetFee(destChainSelector, msg)
			if err != nil {
				return ccipMsgs, statDetails, err
			}
			// transfer fee to the multicall address
			if msg.FeeToken != (common.Address{}) {
				allFee = new(big.Int).Add(allFee, fee)
			}
			msgData := contracts.CCIPMsgData{
				RouterAddr:    e2eLoad.Lane.Source.Common.Router.EthAddress,
				ChainSelector: destChainSelector,
				Msg:           msg,
				Fee:           fee,
			}
			ccipMsgs = append(ccipMsgs, msgData)

			allStatsForDest = append(allStatsForDest, stats)
			allMsgsForDest = append(allMsgsForDest, msgData)
		}
		statDetails[e2eLoad.Lane.DestNetworkName] = MultiCallReturnValues{
			Stats: allStatsForDest,
			Msgs:  allMsgsForDest,
		}
		// transfer fee to the multicall address
		if allFee.Cmp(big.NewInt(0)) > 0 {
			if err := e2eLoad.Lane.Source.Common.FeeToken.Transfer(e2eLoad.Lane.Source.Common.MulticallContract.Hex(), allFee); err != nil {
				return ccipMsgs, statDetails, err
			}
		}
	}
	return ccipMsgs, statDetails, nil
}
