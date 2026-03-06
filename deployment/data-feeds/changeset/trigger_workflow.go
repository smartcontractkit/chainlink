package changeset

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// TriggerWorkflowChangeset is a changeset that triggers a CRE workflow via the gateway API.
// It constructs a JSON-RPC request, signs it with a JWT using the chain's signing function, and sends it to the specified gateway URL.
var TriggerWorkflowChangeset = cldf.CreateChangeSet(triggerWorkflowLogic, triggerWorkflowPrecondition)

func triggerWorkflowLogic(env cldf.Environment, c types.TriggerCREWorkflowConfig) (cldf.ChangesetOutput, error) {
	// Arbitrary EVM chainSelector can be used here since the KMS is the same across all EVM chains.
	// We just need the SignHash function and deployer key to sign the JWT.
	chain := env.BlockChains.EVMChains()[c.ChainSelector]

	var input any = struct{}{}
	if c.Input != nil {
		input = c.Input
	}

	rpcReq := jsonRPCRequest{
		ID:      uuid.New().String(),
		JSONRPC: "2.0",
		Method:  "workflows.execute",
		Params: rpcParams{
			Input:    input,
			Workflow: workflowSelector{WorkflowID: c.WorkflowID},
		},
	}

	reqBody, err := json.Marshal(rpcReq)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to marshal JSON-RPC request: %w", err)
	}

	jwt, err := createSignedJWT(chain.SignHash, chain.DeployerKey.From, reqBody)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create JWT: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.GatewayURL, strings.NewReader(string(reqBody)))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+jwt)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return cldf.ChangesetOutput{}, fmt.Errorf("gateway returned status %d: %s", resp.StatusCode, string(body))
	}

	env.Logger.Infof("Workflow triggered successfully. Response: %s", string(body))

	return cldf.ChangesetOutput{}, nil
}

func triggerWorkflowPrecondition(env cldf.Environment, c types.TriggerCREWorkflowConfig) error {
	if c.GatewayURL == "" {
		return fmt.Errorf("gatewayUrl is required")
	}

	chain, ok := env.BlockChains.EVMChains()[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain %d not found in environment", c.ChainSelector)
	}

	if chain.SignHash == nil {
		return fmt.Errorf("chain %d does not have a SignHash function configured", c.ChainSelector)
	}

	if c.WorkflowID == "" {
		return fmt.Errorf("workflowID is required")
	}

	return nil
}

type jsonRPCRequest struct {
	ID      string    `json:"id"`
	JSONRPC string    `json:"jsonrpc"`
	Method  string    `json:"method"`
	Params  rpcParams `json:"params"`
}

type rpcParams struct {
	Input    any              `json:"input"`
	Workflow workflowSelector `json:"workflow"`
}

type workflowSelector struct {
	WorkflowID string `json:"workflowID"`
}

func createSignedJWT(
	signHash func([]byte) ([]byte, error),
	signerAddr interface{ Hex() string },
	requestBody []byte,
) (string, error) {
	header := jwtHeader{Alg: "ETH", Typ: "JWT"}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JWT header: %w", err)
	}

	digest := sha256.Sum256(requestBody)
	now := time.Now().Unix()
	payload := jwtPayload{
		Digest: fmt.Sprintf("0x%x", digest),
		ISS:    signerAddr.Hex(),
		IAT:    now,
		Exp:    now + 300,
		JTI:    uuid.New().String(),
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JWT payload: %w", err)
	}

	encodedHeader := base64URLEncode(headerJSON)
	encodedPayload := base64URLEncode(payloadJSON)
	rawMessage := encodedHeader + "." + encodedPayload

	prefixed := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(rawMessage), rawMessage)
	msgHash := crypto.Keccak256([]byte(prefixed))

	ethSig, err := signHash(msgHash)
	if err != nil {
		return "", fmt.Errorf("signing failed: %w", err)
	}

	return rawMessage + "." + base64URLEncode(ethSig), nil
}

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type jwtPayload struct {
	Digest string `json:"digest"`
	ISS    string `json:"iss"`
	IAT    int64  `json:"iat"`
	Exp    int64  `json:"exp"`
	JTI    string `json:"jti"`
}

func base64URLEncode(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")
	encoded = strings.TrimRight(encoded, "=")
	return encoded
}
