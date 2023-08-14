package legacygasstation

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/legacygasstation/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/jsonrpc"
)

var (
	forwarderABI = evmtypes.MustGetABI(forwarder.ForwarderABI)
	typeHashRaw  = crypto.Keccak256([]byte("ForwardRequest(address from,address target,uint256 nonce,bytes data,uint256 validUntilTime)"))
)

const calldataDefinition = `
[
	{
		"inputs": [{
			"internalType": "address",
			"name": "receiver",
			"type": "address"
		}, {
			"internalType": "uint256",
			"name": "amount",
			"type": "uint256"
		}, {
			"internalType": "uint64",
			"name": "destinationChainId",
			"type": "uint64"
		}],
		"name": "metaTransfer",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]
`

type RequestHandler struct {
	lggr              logger.Logger
	forwarder         forwarder.ForwarderInterface
	chainID           uint64
	ccipChainSelector uint64
	txm               txmgr.TxManager
	domainTypeHash    common.Hash
	typeHash          common.Hash
	gethks            keystore.Eth
	q                 pg.Q
	job               job.Job
	cfg               Config
	orm               ORM
	fromAddresses     []ethkey.EIP55Address
	pr                pipeline.Runner
	jb                job.Job
	metaTransferAbi   abi.ABI
}

func NewRequestHandler(
	lggr logger.Logger,
	forwarder forwarder.ForwarderInterface,
	chainID uint64,
	ccipChainSelector uint64,
	txm txmgr.TxManager,
	gethks keystore.Eth,
	q pg.Q,
	job job.Job,
	cfg Config,
	orm ORM,
	fromAddresses []ethkey.EIP55Address,
	pr pipeline.Runner,
) (*RequestHandler, error) {
	domainType, err := forwarder.EIP712DOMAINTYPE(nil)
	if err != nil {
		return nil, errors.Wrap(err, "eip712 domain type")
	}
	domainTypeHashRaw := crypto.Keccak256([]byte(domainType))
	var (
		domainTypeHash [32]byte
		typeHash       [32]byte
	)
	copy(domainTypeHash[:], domainTypeHashRaw[:])
	copy(typeHash[:], typeHashRaw[:])
	metaTransferAbi, err := abi.JSON(strings.NewReader(calldataDefinition))
	if err != nil {
		return nil, errors.Wrap(err, "Error while reading metaTransfer ABI definition")
	}
	return &RequestHandler{
		lggr:              lggr,
		forwarder:         forwarder,
		chainID:           chainID,
		ccipChainSelector: ccipChainSelector,
		txm:               txm,
		domainTypeHash:    domainTypeHash,
		typeHash:          typeHash,
		gethks:            gethks,
		q:                 q,
		job:               job,
		cfg:               cfg,
		orm:               orm,
		fromAddresses:     fromAddresses,
		pr:                pr,
		metaTransferAbi:   metaTransferAbi,
	}, nil
}

func (rh *RequestHandler) CCIPChainSelector() *utils.Big {
	return utils.NewBig(new(big.Int).SetUint64(rh.ccipChainSelector))
}

// SendTransaction submits meta transaction to transaction manager and persists meta-transaction data
func (rh *RequestHandler) SendTransaction(ctx *gin.Context, req types.SendTransactionRequest) (*types.SendTransactionResponse, *jsonrpc.Error) {
	err := validateSendTransactionRequest(req)
	if err != nil {
		return nil, &jsonrpc.Error{
			Code:    jsonrpc.InvalidRequestError,
			Message: err.Error(),
		}
	}

	l := rh.lggr.With("from", req.From.Hex(),
		"sourceCCIPChainSelector", req.SourceChainID,
		"destinationCCIPChainSelector", req.DestinationChainID,
		"nonce", req.Nonce,
		"receiver", req.Receiver.Hex(),
		"target", req.Target.Hex(),
		"targetName", req.TargetName,
		"targetVersion", req.Version,
		"validUntilTime", req.ValidUntilTime,
		"amount", req.Amount,
		"signature", hex.EncodeToString(req.Signature),
	)

	calldata, err := rh.metaTransferAbi.Pack("metaTransfer", req.Receiver, req.Amount, req.DestinationChainID)
	if err != nil {
		l.Errorw("Error while packing metaTransfer", "err", err)
		return nil, &jsonrpc.Error{
			Code:    jsonrpc.InternalError,
			Message: jsonrpc.InternalServerErrorMsg,
		}
	}

	fromAddresses := rh.sendingKeys()
	fromAddress, err := rh.gethks.GetRoundRobinAddress(big.NewInt(0).SetUint64(rh.chainID), fromAddresses...)
	if err != nil {
		l.Errorw("Couldn't get next from address", "err", err)
		return nil, &jsonrpc.Error{
			Code:    jsonrpc.InternalError,
			Message: jsonrpc.InternalServerErrorMsg,
		}
	}

	forwardReq := forwarder.IForwarderForwardRequest{
		From:           req.From,
		Target:         req.Target,
		Nonce:          req.Nonce,
		Data:           calldata,
		ValidUntilTime: req.ValidUntilTime,
	}

	requestID := uuid.New().String()

	domainSeparatorHash, err := rh.domainSeparatorHash(req.TargetName, req.Version)
	if err != nil {
		l.Errorw("Error while getting domain separator", "err", err)
		return nil, &jsonrpc.Error{
			Code:    jsonrpc.InternalError,
			Message: jsonrpc.InternalServerErrorMsg,
		}
	}

	payload, err := forwarderABI.Pack("execute", forwardReq, domainSeparatorHash, rh.typeHash, []byte{}, req.Signature)
	if err != nil {
		l.Errorw("Error while packing", "err", err)
		return nil, &jsonrpc.Error{
			Code:    jsonrpc.InternalError,
			Message: jsonrpc.InternalServerErrorMsg,
		}
	}

	pipelineResult, jsonrpcErr := rh.executePipline(ctx, l, payload)
	if jsonrpcErr != nil {
		return nil, jsonrpcErr
	}

	// Creation of eth transaction and persistence of data are done in a transaction
	// to avoid partial failures, which would leave the persistence layer in inconsistent state
	err = rh.q.Transaction(func(tx pg.Queryer) error {
		ethTx, err2 := rh.txm.CreateTransaction(txmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      rh.forwarder.Address(),
			EncodedPayload: payload,
			FeeLimit:       uint32(pipelineResult.gas), // safe down-cast because we cap gas at EvmGasLimitDefault, which is uint32
			// TODO: add new field in eth tx meta
			//Meta: &txmgr.EthTxMeta{
			//	RequestID: requestID,
			//},
			Strategy: txmgrcommon.NewSendEveryStrategy(),
		}, pg.WithQueryer(tx), pg.WithParentCtx(ctx))
		if err2 != nil {
			return err2
		}

		l.Debugw("created Eth tx", "ethTxID", ethTx.GetID())

		gaslessTx := types.LegacyGaslessTx{
			ID:                 requestID,
			From:               req.From,
			Target:             req.Target,
			Forwarder:          rh.forwarder.Address(),
			Nonce:              utils.NewBig(req.Nonce),
			Receiver:           req.Receiver,
			Amount:             utils.NewBig(req.Amount),
			SourceChainID:      req.SourceChainID,
			DestinationChainID: req.DestinationChainID,
			ValidUntilTime:     utils.NewBig(req.ValidUntilTime),
			Signature:          req.Signature,
			Status:             types.Submitted,
			TokenName:          req.TargetName,
			TokenVersion:       req.Version,
			EthTxID:            ethTx.GetID(),
		}

		err2 = rh.orm.InsertLegacyGaslessTx(gaslessTx, pg.WithQueryer(tx), pg.WithParentCtx(ctx))

		return err2
	})

	if err != nil {
		l.Errorw("Error while writing to DB", "err", err)
		return nil, &jsonrpc.Error{
			Code:    jsonrpc.InternalError,
			Message: jsonrpc.InternalServerErrorMsg,
		}
	}

	return &types.SendTransactionResponse{
		RequestID: requestID,
	}, nil
}

func (rh *RequestHandler) sendingKeys() []common.Address {
	var addresses []common.Address
	for _, a := range rh.fromAddresses {
		addresses = append(addresses, a.Address())
	}
	return addresses
}

// domainSeparatorHash is Keccak256 hash of token name, token version, chain ID and forwarder address
func (rh *RequestHandler) domainSeparatorHash(name, version string) (domainSeparatorHash common.Hash, err error) {
	domainSeparator, err := rh.forwarder.GetDomainSeparator(nil, name, version)
	if err != nil {
		return
	}
	domainSeparatorHashRaw := crypto.Keccak256(domainSeparator)
	copy(domainSeparatorHash[:], domainSeparatorHashRaw[:])
	return
}

type pipelineResult struct {
	gas uint32
}

func (rh *RequestHandler) executePipline(ctx *gin.Context, l logger.Logger, payload []byte) (*pipelineResult, *jsonrpc.Error) {
	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"databaseID":    rh.jb.ID,
			"externalJobID": rh.jb.ExternalJobID,
			"name":          rh.jb.Name.ValueOrZero(),
			"maxGasPrice":   rh.cfg.PriceMax().ToInt().String(),
		},
		"jobRun": map[string]interface{}{
			"payload": payload[:],
		},
	})
	runs, trrs, err := rh.pr.ExecuteRun(ctx, *rh.job.PipelineSpec, vars, l)
	if err != nil {
		l.Errorw("Error while exeucting run", "err", err)
		return nil, &jsonrpc.Error{
			Code:    jsonrpc.InternalError,
			Message: jsonrpc.InternalServerErrorMsg,
		}
	}
	if runs.AllErrors.HasError() {
		err = runs.AllErrors.ToError()
		l.Warnw("Pipeline run returned error", "err", err)
		if strings.Contains(err.Error(), "execution reverted") {
			return nil, &jsonrpc.Error{
				Code:    jsonrpc.InvalidRequestError,
				Message: fmt.Sprintf("Error while simulating transaction: %s", err.Error()),
			}
		}
		return nil, &jsonrpc.Error{
			Code:    jsonrpc.InternalError,
			Message: jsonrpc.InternalServerErrorMsg,
		}
	}

	gas := rh.cfg.LimitDefault()
	for _, trr := range trrs {
		if trr.Task.Type() == pipeline.TaskTypeEstimateGasLimit {
			gas = trr.Result.Value.(uint32)
		}
	}

	return &pipelineResult{
		gas: gas,
	}, nil

	//TODO: store pipeline runs to pipelines table
}
