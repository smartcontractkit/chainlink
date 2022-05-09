package web

import (
	"net/http"

	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/client"

	solanaGo "github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	solanamodels "github.com/smartcontractkit/chainlink/core/store/models/solana"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// SolanaTransfersController can send LINK tokens to another address
type SolanaTransfersController struct {
	App chainlink.Application
}

// Create sends Luna and other native coins from the Chainlink's account to a specified address.
func (tc *SolanaTransfersController) Create(c *gin.Context) {
	solanaChains := tc.App.GetChains().Solana
	if solanaChains == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrSolanaNotEnabled)
		return
	}

	var tr solanamodels.SendRequest
	if err := c.ShouldBindJSON(&tr); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	if tr.SolanaChainID == "" {
		jsonAPIError(c, http.StatusBadRequest, errors.New("missing solanaChainID"))
		return
	}
	chain, err := solanaChains.Chain(c.Request.Context(), tr.SolanaChainID)
	switch err {
	case chains.ErrChainIDInvalid, chains.ErrChainIDEmpty:
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	case nil:
		break
	default:
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	if tr.From.IsZero() {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("source address is missing: %v", tr.From))
		return
	}

	if tr.Amount == 0 {
		jsonAPIError(c, http.StatusBadRequest, errors.New("amount must be greater than zero"))
		return
	}

	fromKey, err := tc.App.GetKeyStore().Solana().Get(tr.From.String())
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("fail to get key: %v", err))
		return
	}

	txm := chain.TxManager()
	var reader client.Reader
	reader, err = chain.Reader()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, errors.Errorf("chain unreachable: %v", err))
		return
	}

	blockhash, err := reader.LatestBlockhash()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, errors.Errorf("failed to get latest block hash: %v", err))
		return
	}

	tx, err := solanaGo.NewTransaction(
		[]solanaGo.Instruction{
			system.NewTransferInstruction(
				tr.Amount,
				tr.From,
				tr.To,
			).Build(),
		},
		blockhash.Value.Blockhash,
		solanaGo.TransactionPayer(tr.From),
	)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, errors.Errorf("failed to create tx: %v", err))
		return
	}

	if !tr.AllowHigherAmounts {
		if err := solanaValidateBalance(reader, tr.From, tr.Amount, tx.Message.ToBase64()); err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("failed to validate balance: %v", err))
			return
		}
	}

	// marshal transaction
	msg, err := tx.Message.MarshalBinary()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, errors.Errorf("failed to marshal tx: %v", err))
		return
	}

	// sign tx
	sigBytes, err := fromKey.Sign(msg)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, errors.Errorf("failed to sign tx: %v", err))
		return
	}
	var finalSig [64]byte
	copy(finalSig[:], sigBytes)
	tx.Signatures = append(tx.Signatures, finalSig)

	err = txm.Enqueue("", tx)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, errors.Errorf("transaction failed: %v", err))
		return
	}

	resource := presenters.NewSolanaMsgResource("sol_transfer_"+uuid.New().String(), tr.SolanaChainID)
	resource.Amount = tr.Amount
	resource.From = tr.From.String()
	resource.To = tr.To.String()

	jsonAPIResponse(c, resource, "solana_tx")
}

func solanaValidateBalance(reader client.Reader, from solanaGo.PublicKey, amount uint64, msg string) error {
	balance, err := reader.Balance(from)
	if err != nil {
		return err
	}

	fee, err := reader.GetFeeForMessage(msg)
	if err != nil {
		return err
	}

	if balance < (amount + fee) {
		return errors.Errorf("balance %d is too low for this transaction to be executed: amount %d + fee %d", balance, amount, fee)
	}
	return nil
}
