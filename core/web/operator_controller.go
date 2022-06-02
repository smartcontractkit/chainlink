package web

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/operators"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/operator_factory"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/gin-gonic/gin"
)

var deployOperatorABI = evmtypes.MustGetABI(operator_factory.OperatorFactoryABI).Methods["deployNewOperator"]

// OperatorsController manages operator contracts.
type OperatorsController struct {
	App chainlink.Application
}

// Index lists all operator contracts.
func (cc *OperatorsController) Index(c *gin.Context, size, page, offset int) {
	orm := operators.NewORM(cc.App.GetSqlxDB(), cc.App.GetLogger(), cc.App.GetConfig())
	ops, count, err := orm.FindOperators(0, size)

	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}

	var resources []presenters.OperatorResource
	for _, op := range ops {
		resources = append(resources, presenters.NewOperatorResource(op))
	}

	paginatedResponse(c, "operator", size, page, resources, count, err)
}

// DeployOperatorRequest is a JSONAPI request for creating an operator.
type DeployOperatorRequest struct {
	ChainID *utils.Big     `json:"chainID"`
	Owner   common.Address `json:"owner"`
}

// Deploy deploys a new operator on chain and return tx.
func (cc *OperatorsController) Deploy(c *gin.Context) {
	request := &DeployOperatorRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	chain, err := getChain(cc.App.GetChains().EVM, request.ChainID.String())

	switch err {
	case ErrInvalidChainID, ErrMultipleChains, ErrMissingChainID:
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	case nil:
		break
	default:
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	factoryAddr := chain.Config().OperatorFactoryAddress()
	if !strings.HasPrefix(factoryAddr, "0x") {
		jsonAPIError(c,
			http.StatusBadRequest,
			errors.New("OperatorFactoryAddress is not set for this chain"))
		return
	}

	inputs, err := deployOperatorABI.Inputs.Pack()
	if err != nil {
		// This should never happen
		jsonAPIError(c,
			http.StatusInternalServerError,
			errors.New("error packing deployOperator inputs"))
		return
	}
	tx, err := chain.TxManager().CreateEthTransaction(txmgr.NewTx{
		FromAddress:    request.Owner,
		ToAddress:      common.HexToAddress(factoryAddr),
		EncodedPayload: append(deployOperatorABI.ID, inputs...),
		GasLimit:       5000000,
		Strategy:       txmgr.NewSendEveryStrategy(),
	})
	if err != nil {
		jsonAPIError(c,
			http.StatusInternalServerError,
			err)
		return
	}

	jsonAPIResponseWithStatus(c, presenters.NewJAIDInt64(tx.ID), "eth_tx_id", http.StatusAccepted)
}

// Create creates a new operator.
func (cc *OperatorsController) Status(c *gin.Context) {
	txId, err := stringutils.ToInt64(c.Param("txID"))
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	chainId := c.Param("chainID")

	chain, err := getChain(cc.App.GetChains().EVM, chainId)
	switch err {
	case ErrInvalidChainID, ErrMultipleChains, ErrMissingChainID:
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	case nil:
		break
	default:
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ethtx, err := cc.App.TxmORM().FindEthTxWithAttempts(txId)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	parser, err := operator_factory.NewOperatorFactory(common.Address{}, chain.Client())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	cc.App.GetLogger().Criticalf("Tx: %v State: %v", ethtx, ethtx.State)
	// If tx is still unconfirmed notify back.
	if ethtx.State != txmgr.EthTxConfirmed {
		jsonAPIError(c, http.StatusNotFound, errors.New("operator not yet deployed"))
		return
	}

	// If tx is confirmed, succeed only if the there is one succeeding attempt
	// This seems like a high order loop but the reasonable bounds are [1-2]/1/1
	for _, attempt := range ethtx.EthTxAttempts {

		for _, receipt := range attempt.EthReceipts {
			// this should change.
			if receipt.Receipt.Status == types.ReceiptStatusFailed {
				continue
			}
			for _, log := range receipt.Receipt.Logs {
				cc.App.GetLogger().Criticalf("Log: %v", log)
				opcr, err := parser.ParseOperatorCreated(*log.ToGethLog())
				if err != nil {
					continue
				}
				orm := operators.NewORM(cc.App.GetSqlxDB(), cc.App.GetLogger(), cc.App.GetConfig())
				opr, err := orm.CreateOperator(opcr.Operator, utils.Big(*chain.ID()))
				if err != nil {
					jsonAPIResponseWithStatus(c, presenters.NewOperatorResource(opr), "operator", http.StatusCreated)
					return
				}

			}
		}
	}
	jsonAPIError(c, http.StatusFailedDependency, errors.New("operator deployment failed"))
}

// Delete removes an operator.
func (cc *OperatorsController) Delete(c *gin.Context) {
	id, err := stringutils.ToInt32(c.Param("operatorID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	orm := operators.NewORM(cc.App.GetSqlxDB(), cc.App.GetLogger(), cc.App.GetConfig())
	err = orm.DeleteOperator(id)

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, nil, "operator", http.StatusNoContent)
}
