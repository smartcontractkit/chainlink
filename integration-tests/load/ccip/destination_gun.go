package ccip

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solstate "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	soltokens "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"go.uber.org/atomic"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipchangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/testconfig/ccip"
)

type SeqNumRange struct {
	Start *atomic.Uint64
	End   *atomic.Uint64
}

type DestinationGun struct {
	l                logger.Logger
	env              deployment.Environment
	state            *ccipchangeset.CCIPOnChainState
	roundNum         *atomic.Int32
	chainSelector    uint64
	receiver         common.Address
	testConfig       *ccip.LoadConfig
	evmSourceKeys    map[uint64]*bind.TransactOpts
	solanaSourceKeys map[uint64]*solana.PrivateKey
	chainOffset      int
	metricPipe       chan messageData
}

func NewDestinationGun(
	l logger.Logger,
	chainSelector uint64,
	env deployment.Environment,
	state *ccipchangeset.CCIPOnChainState,
	receiver common.Address,
	overrides *ccip.LoadConfig,
	evmSourceKeys map[uint64]*bind.TransactOpts,
	solanaSourceKeys map[uint64]*solana.PrivateKey,
	chainOffset int,
	metricPipe chan messageData,
) (*DestinationGun, error) {
	dg := DestinationGun{
		l:                l,
		env:              env,
		state:            state,
		roundNum:         &atomic.Int32{},
		chainSelector:    chainSelector,
		receiver:         receiver,
		testConfig:       overrides,
		evmSourceKeys:    evmSourceKeys,
		solanaSourceKeys: solanaSourceKeys,
		chainOffset:      chainOffset,
		metricPipe:       metricPipe,
	}

	return &dg, nil
}

func (m *DestinationGun) Call(_ *wasp.Generator) *wasp.Response {
	m.roundNum.Add(1)
	src, err := m.MustSourceChain()
	if err != nil {
		return &wasp.Response{Error: err.Error(), Group: "", Failed: true}
	}
	waspGroup := fmt.Sprintf("%d->%d", src, m.chainSelector)

	selectorFamily, err := selectors.GetSelectorFamily(src)
	if err != nil {
		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
	}

	switch selectorFamily {
	case selectors.FamilyEVM:
		err = m.sendEVMMessage(src)
	case selectors.FamilySolana:
		err = m.sendSolanaMessage(src)
	}
	if err != nil {
		m.l.Errorw("Failed to transmit message",
			"gun", waspGroup,
			"sourceChainFamily", selectorFamily,
			err, deployment.MaybeDataErr(err))
		if m.metricPipe != nil {
			// in the event of an error, still push a metric
			// sequence numbers start at 1 so using 0 as a sentinel value
			data := messageData{
				eventType: transmitted,
				srcDstSeqNum: srcDstSeqNum{
					src:    src,
					dst:    m.chainSelector,
					seqNum: 0,
				},
				timestamp: uint64(time.Now().Unix()), //nolint:gosec // G115
			}
			m.metricPipe <- data
		}

		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
	}

	return &wasp.Response{Failed: false, Group: waspGroup}
}

// MustSourceChain will return a chain selector to send a message from
func (m *DestinationGun) MustSourceChain() (uint64, error) {
	otherCS := m.env.AllChainSelectorsExcluding([]uint64{m.chainSelector})
	// todo: uncomment to enable solana as a source chain
	// otherCS := m.env.AllChainSelectorsAllFamilliesExcluding([]uint64{m.chainSelector})
	if len(otherCS) == 0 {
		return 0, errors.New("no other chains to send from")
	}
	index := (int(m.roundNum.Load()) + m.chainOffset) % len(otherCS)
	return otherCS[index], nil
}
func (m *DestinationGun) sendEVMMessage(src uint64) error {
	acc := m.evmSourceKeys[src]
	r := m.state.Chains[src].Router

	msg, gasLimit, err := m.GetEVMMessage(src)
	if err != nil {
		return err
	}
	// Set the gas limit for this tx
	if gasLimit != 0 {
		//nolint:gosec // it's okay here
		acc.GasLimit = uint64(gasLimit)
	}

	fee, err := r.GetFee(
		&bind.CallOpts{Context: context.Background()}, m.chainSelector, msg)
	if err != nil {
		m.l.Errorw("could not get fee ",
			"dstChainSelector", m.chainSelector,
			"fee", fee,
			"err", deployment.MaybeDataErr(err))
		return err
	}
	if msg.FeeToken == common.HexToAddress("0x0") {
		acc.Value = fee
	}
	msgWithoutData := msg
	msgWithoutData.Data = nil
	m.l.Debugw("sending message ",
		"srcChain", src,
		"dstChain", m.chainSelector,
		"fee", fee,
		"msg size", len(msg.Data),
		"msgWithoutData", msgWithoutData)
	tx, err := r.CcipSend(
		acc,
		m.chainSelector,
		msg)
	if err != nil {
		m.l.Errorw("execution reverted from ",
			"sourceChain", src,
			"destchain", m.chainSelector,
			"err", deployment.MaybeDataErr(err))
		return err
	}

	_, err = m.env.Chains[src].Confirm(tx)
	if err != nil {
		m.l.Errorw("could not confirm tx on source", "tx", tx, "err", deployment.MaybeDataErr(err))
		return err
	}

	return nil
}

// GetEVMMessage will return the message to be sent while considering expected load of different messages
// returns the message, gas limit
func (m *DestinationGun) GetEVMMessage(src uint64) (router.ClientEVM2AnyMessage, int64, error) {
	rcv, err := utils.ABIEncode(`[{"type":"address"}]`, m.receiver)
	if err != nil {
		m.l.Error("Error encoding receiver address")
		return router.ClientEVM2AnyMessage{}, 0, err
	}
	extraArgs, err := GetEVMExtraArgsV2(big.NewInt(0), *m.testConfig.OOOExecution)
	if err != nil {
		m.l.Error("Error encoding extra args")
		return router.ClientEVM2AnyMessage{}, 0, err
	}

	// Select a message type based on ratio
	randomValue := mathrand.Intn(100)
	accumulatedRatio := 0
	var selectedMsgDetails *ccip.MsgDetails

	for _, msg := range *m.testConfig.MessageDetails {
		accumulatedRatio += *msg.Ratio
		if randomValue < accumulatedRatio {
			selectedMsgDetails = &msg
			break
		}
	}

	if selectedMsgDetails == nil {
		return router.ClientEVM2AnyMessage{}, 0, errors.New("failed to select message type")
	}

	m.l.Infow("Selected message type", "msgType", *selectedMsgDetails.MsgType)

	message := router.ClientEVM2AnyMessage{
		Receiver:  rcv,
		FeeToken:  common.HexToAddress("0x0"),
		ExtraArgs: extraArgs,
	}

	// Set data length if it's a data transfer
	if selectedMsgDetails.IsDataTransfer() {
		dataLength := *selectedMsgDetails.DataLengthBytes
		data := make([]byte, dataLength)
		_, err2 := rand.Read(data)
		if err2 != nil {
			return router.ClientEVM2AnyMessage{}, 0, err2
		}
		message.Data = data
	}

	// When it's not a programmable token transfer the receiver can be an EOA, we use a random address to denote that
	if selectedMsgDetails.IsTokenOnlyTransfer() {
		receiver, err := utils.ABIEncode(`[{"type":"address"}]`, common.HexToAddress(utils.RandomAddress().Hex()))
		if err != nil {
			m.l.Error("Error encoding receiver address")
			return router.ClientEVM2AnyMessage{}, 0, err
		}
		message.Receiver = receiver
	}

	// Set token amounts if it's a token transfer
	if selectedMsgDetails.IsTokenTransfer() {
		message.TokenAmounts = []router.ClientEVMTokenAmount{
			{
				Token:  m.state.Chains[src].LinkToken.Address(),
				Amount: big.NewInt(1),
			},
		}
	}

	gasLimit := int64(0)
	if selectedMsgDetails.DestGasLimit != nil {
		gasLimit = *selectedMsgDetails.DestGasLimit
	}

	return message, gasLimit, nil
}

func GetEVMExtraArgsV2(gasLimit *big.Int, allowOutOfOrder bool) ([]byte, error) {
	EVMV2Tag := hexutil.MustDecode("0x181dcf10")

	encodedArgs, err := utils.ABIEncode(`[{"type":"uint256"},{"type":"bool"}]`, gasLimit, allowOutOfOrder)
	if err != nil {
		return nil, err
	}

	return append(EVMV2Tag, encodedArgs...), nil
}

func (m *DestinationGun) sendSolanaMessage(src uint64) error {
	acc := m.solanaSourceKeys[src]
	s := m.state.SolChains[src]
	sourceKey := m.solanaSourceKeys[src]

	msg, err := m.getSolanaMessage(src, acc)
	if err != nil {
		return err
	}

	// if fee token is 0, fallback to WSOL
	if msg.FeeToken.IsZero() {
		msg.FeeToken = m.state.SolChains[src].WSOL
	}

	destinationChainStatePDA, err := solstate.FindDestChainStatePDA(m.chainSelector, s.Router)
	if err != nil {
		return err
	}

	noncePDA, err := solstate.FindNoncePDA(m.chainSelector, sourceKey.PublicKey(), s.Router)
	if err != nil {
		return err
	}
	feeToken := msg.FeeToken

	linkFqBillingConfigPDA, _, err := solstate.FindFqBillingTokenConfigPDA(s.LinkToken, s.FeeQuoter)
	if err != nil {
		return err
	}
	feeTokenFqBillingConfigPDA, _, err := solstate.FindFqBillingTokenConfigPDA(feeToken, s.FeeQuoter)
	if err != nil {
		return err
	}

	billingSignerPDA, _, err := solstate.FindFeeBillingSignerPDA(s.Router)
	if err != nil {
		return err
	}

	feeTokenUserATA, _, err := soltokens.FindAssociatedTokenAddress(solana.TokenProgramID, feeToken, sourceKey.PublicKey())
	if err != nil {
		return err
	}
	feeTokenReceiverATA, _, err := soltokens.FindAssociatedTokenAddress(solana.TokenProgramID, feeToken, billingSignerPDA)
	if err != nil {
		return err
	}
	fqDestChainPDA, _, err := solstate.FindFqDestChainPDA(m.chainSelector, s.FeeQuoter)
	if err != nil {
		return err
	}

	base := ccip_router.NewCcipSendInstruction(
		m.chainSelector,
		msg,
		[]byte{}, // starting indices for accounts, calculated later
		s.RouterConfigPDA,
		destinationChainStatePDA,
		noncePDA,
		sourceKey.PublicKey(),
		solana.SystemProgramID,
		solana.TokenProgramID,
		feeToken,
		feeTokenUserATA,
		feeTokenReceiverATA,
		billingSignerPDA,
		s.FeeQuoter,
		s.FeeQuoterConfigPDA,
		fqDestChainPDA,
		feeTokenFqBillingConfigPDA,
		linkFqBillingConfigPDA,
		solana.PublicKey{},
		solana.PublicKey{},
		solana.PublicKey{},
	)
	base.GetFeeTokenUserAssociatedAccountAccount().WRITE()
	instruction, err := base.ValidateAndBuild()
	if err != nil {
		m.l.Errorw("failed to build instruction",
			"src", src,
			"dest", m.chainSelector,
			"err", deployment.MaybeDataErr(err))
		return err
	}

	result, err := solcommon.SendAndConfirm(
		m.env.GetContext(),
		m.env.SolChains[src].Client,
		[]solana.Instruction{instruction},
		*sourceKey,
		rpc.CommitmentConfirmed)
	if err != nil || result == nil {
		m.l.Errorw("could not confirm solana tx on source",
			"src", src,
			"dest", m.chainSelector)
		return err
	}

	return nil
}

func (m *DestinationGun) getSolanaMessage(src uint64, account *solana.PrivateKey) (ccip_router.SVM2AnyMessage, error) {
	return ccip_router.SVM2AnyMessage{
		Receiver:     common.LeftPadBytes(m.receiver.Bytes(), 32),
		TokenAmounts: nil,
		Data:         []byte("hello world"),
		ExtraArgs:    []byte{},
	}, nil
}
