package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/sethvargo/go-retry"

	"github.com/smartcontractkit/chainlink/deployment/environment/web/sdk/client/doer"
	"github.com/smartcontractkit/chainlink/deployment/environment/web/sdk/internal/generated"
)

type Client interface {
	FetchCSAPublicKey(ctx context.Context) (*string, error)
	FetchP2PPeerID(ctx context.Context) (*string, error)
	FetchAccountAddress(ctx context.Context, chainID string) (*string, error)
	FetchKeys(ctx context.Context, chainType string) ([]string, error)
	FetchOCR2KeyBundleID(ctx context.Context, chainType string) (string, error)
	GetJob(ctx context.Context, id string) (*generated.GetJobResponse, error)
	ListJobs(ctx context.Context, offset, limit int) (*generated.ListJobsResponse, error)
	GetJobDistributor(ctx context.Context, id string) (generated.FeedsManagerParts, error)
	ListJobDistributors(ctx context.Context) (*generated.ListFeedsManagersResponse, error)
	CreateJobDistributor(ctx context.Context, cmd JobDistributorInput) (string, error)
	UpdateJobDistributor(ctx context.Context, id string, cmd JobDistributorInput) error
	CreateJobDistributorChainConfig(ctx context.Context, in JobDistributorChainConfigInput) (string, error)
	DeleteJobDistributorChainConfig(ctx context.Context, id string) error
	GetJobProposal(ctx context.Context, id string) (*generated.GetJobProposalJobProposal, error)
	ApproveJobProposalSpec(ctx context.Context, id string, force bool) (*JobProposalApprovalSuccessSpec, error)
	CancelJobProposalSpec(ctx context.Context, id string) (*generated.CancelJobProposalSpecCancelJobProposalSpecCancelJobProposalSpecSuccessSpecJobProposalSpec, error)
	RejectJobProposalSpec(ctx context.Context, id string) (*generated.RejectJobProposalSpecResponse, error)
	UpdateJobProposalSpecDefinition(ctx context.Context, id string, cmd generated.UpdateJobProposalSpecDefinitionInput) (*generated.UpdateJobProposalSpecDefinitionResponse, error)
}

type client struct {
	gqlClient   graphql.Client
	credentials Credentials
	endpoints   endpoints
	cookie      string
}

type endpoints struct {
	Sessions string
	Query    string
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func New(baseURI string, creds Credentials) (Client, error) {
	ep := endpoints{
		Sessions: baseURI + "/sessions",
		Query:    baseURI + "/query",
	}
	c := &client{
		endpoints:   ep,
		credentials: creds,
	}

	err := retry.Do(context.Background(), retry.WithMaxDuration(10*time.Second, retry.NewFibonacci(2*time.Second)), func(ctx context.Context) error {
		err := c.login()
		if err != nil {
			return retry.RetryableError(fmt.Errorf("retrying login to node: %w", err))
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to login to node: %w", err)
	}

	c.gqlClient = graphql.NewClient(
		c.endpoints.Query,
		doer.NewAuthed(c.cookie),
	)

	return c, nil
}

func (c *client) FetchCSAPublicKey(ctx context.Context) (*string, error) {
	keys, err := generated.FetchCSAKeys(ctx, c.gqlClient)
	if err != nil {
		return nil, err
	}
	if keys == nil || len(keys.CsaKeys.GetResults()) == 0 {
		return nil, errors.New("no CSA keys found")
	}
	return &keys.CsaKeys.GetResults()[0].PublicKey, nil
}

func (c *client) FetchP2PPeerID(ctx context.Context) (*string, error) {
	keys, err := generated.FetchP2PKeys(ctx, c.gqlClient)
	if err != nil {
		return nil, err
	}
	if keys == nil || len(keys.P2pKeys.GetResults()) == 0 {
		return nil, errors.New("no P2P keys found")
	}
	return &keys.P2pKeys.GetResults()[0].PeerID, nil
}

func (c *client) FetchOCR2KeyBundleID(ctx context.Context, chainType string) (string, error) {
	keyBundles, err := generated.FetchOCR2KeyBundles(ctx, c.gqlClient)
	if err != nil {
		return "", err
	}
	if keyBundles == nil || len(keyBundles.GetOcr2KeyBundles().Results) == 0 {
		return "", errors.New("no ocr2 keybundle found, check if ocr2 is enabled")
	}
	for _, keyBundle := range keyBundles.GetOcr2KeyBundles().Results {
		if keyBundle.ChainType == generated.OCR2ChainType(chainType) {
			return keyBundle.GetId(), nil
		}
	}
	return "", fmt.Errorf("no ocr2 keybundle found for chain type %s", chainType)
}

func (c *client) FetchAccountAddress(ctx context.Context, chainID string) (*string, error) {
	keys, err := generated.FetchAccounts(ctx, c.gqlClient)
	if err != nil {
		return nil, err
	}
	if keys == nil || len(keys.EthKeys.GetResults()) == 0 {
		return nil, errors.New("no accounts found")
	}
	for _, keyDetail := range keys.EthKeys.GetResults() {
		if keyDetail.GetChain().Enabled && keyDetail.GetChain().Id == chainID {
			return &keyDetail.Address, nil
		}
	}
	return nil, fmt.Errorf("no account found for chain %s", chainID)
}

func (c *client) FetchKeys(ctx context.Context, chainType string) ([]string, error) {
	keys, err := generated.FetchKeys(ctx, c.gqlClient)
	if err != nil {
		return nil, err
	}
	if keys == nil {
		return nil, errors.New("no accounts found")
	}
	switch generated.OCR2ChainType(chainType) {
	case generated.OCR2ChainTypeAptos:
		var accounts []string
		for _, key := range keys.AptosKeys.GetResults() {
			accounts = append(accounts, key.Account)
		}
		return accounts, nil
	case generated.OCR2ChainTypeSolana:
		var accounts []string
		for _, key := range keys.SolanaKeys.GetResults() {
			accounts = append(accounts, key.Id)
		}
		return accounts, nil
	default:
		return nil, fmt.Errorf("unsupported chainType %v", chainType)
	}
}

func (c *client) GetJob(ctx context.Context, id string) (*generated.GetJobResponse, error) {
	return generated.GetJob(ctx, c.gqlClient, id)
}

func (c *client) ListJobs(ctx context.Context, offset, limit int) (*generated.ListJobsResponse, error) {
	return generated.ListJobs(ctx, c.gqlClient, offset, limit)
}

func (c *client) GetBridge(ctx context.Context, id string) (*generated.GetBridgeResponse, error) {
	return generated.GetBridge(ctx, c.gqlClient, id)
}

func (c *client) ListBridges(ctx context.Context, offset, limit int) (*generated.ListBridgesResponse, error) {
	return generated.ListBridges(ctx, c.gqlClient, offset, limit)
}

func (c *client) GetJobDistributor(ctx context.Context, id string) (generated.FeedsManagerParts, error) {
	res, err := generated.GetFeedsManager(ctx, c.gqlClient, id)
	if err != nil {
		return generated.FeedsManagerParts{}, err
	}
	if res == nil {
		return generated.FeedsManagerParts{}, errors.New("no feeds manager found")
	}
	if success, ok := res.GetFeedsManager().(*generated.GetFeedsManagerFeedsManager); ok {
		return success.FeedsManagerParts, nil
	}
	return generated.FeedsManagerParts{}, errors.New("failed to get feeds manager")
}

func (c *client) ListJobDistributors(ctx context.Context) (*generated.ListFeedsManagersResponse, error) {
	return generated.ListFeedsManagers(ctx, c.gqlClient)
}

func (c *client) CreateJobDistributor(ctx context.Context, in JobDistributorInput) (string, error) {
	var cmd generated.CreateFeedsManagerInput
	err := DecodeInput(in, &cmd)
	if err != nil {
		return "", err
	}
	response, err := generated.CreateFeedsManager(ctx, c.gqlClient, cmd)
	if err != nil {
		return "", err
	}
	// Access the FeedsManager ID
	if success, ok := response.GetCreateFeedsManager().(*generated.CreateFeedsManagerCreateFeedsManagerCreateFeedsManagerSuccess); ok {
		feedsManager := success.GetFeedsManager()
		return feedsManager.GetId(), nil
	}
	if err, ok := response.GetCreateFeedsManager().(*generated.CreateFeedsManagerCreateFeedsManagerSingleFeedsManagerError); ok {
		msg := err.GetMessage()
		return "", fmt.Errorf("failed to create feeds manager: %v", msg)
	}
	return "", fmt.Errorf("failed to create feeds manager: %v", response.GetCreateFeedsManager().GetTypename())
}

func (c *client) UpdateJobDistributor(ctx context.Context, id string, in JobDistributorInput) error {
	var cmd generated.UpdateFeedsManagerInput
	err := DecodeInput(in, &cmd)
	if err != nil {
		return err
	}
	_, err = generated.UpdateFeedsManager(ctx, c.gqlClient, id, cmd)
	return err
}

func (c *client) CreateJobDistributorChainConfig(ctx context.Context, in JobDistributorChainConfigInput) (string, error) {
	var cmd generated.CreateFeedsManagerChainConfigInput
	err := DecodeInput(in, &cmd)
	if err != nil {
		return "", err
	}
	res, err := generated.CreateFeedsManagerChainConfig(ctx, c.gqlClient, cmd)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("failed to create feeds manager chain config")
	}
	if success, ok := res.GetCreateFeedsManagerChainConfig().(*generated.CreateFeedsManagerChainConfigCreateFeedsManagerChainConfigCreateFeedsManagerChainConfigSuccess); ok {
		return success.ChainConfig.Id, nil
	}
	return "", errors.New("failed to create feeds manager chain config")
}

func (c *client) DeleteJobDistributorChainConfig(ctx context.Context, id string) error {
	res, err := generated.DeleteFeedsManagerChainConfig(ctx, c.gqlClient, id)
	if err != nil {
		return err
	}
	if res == nil {
		return errors.New("failed to delete feeds manager chain config")
	}
	if _, ok := res.GetDeleteFeedsManagerChainConfig().(*generated.DeleteFeedsManagerChainConfigDeleteFeedsManagerChainConfigDeleteFeedsManagerChainConfigSuccess); ok {
		return nil
	}
	return errors.New("failed to delete feeds manager chain config")
}

func (c *client) GetJobProposal(ctx context.Context, id string) (*generated.GetJobProposalJobProposal, error) {
	proposal, err := generated.GetJobProposal(ctx, c.gqlClient, id)
	if err != nil {
		return nil, err
	}
	if proposal == nil {
		return nil, errors.New("no job proposal found")
	}
	if success, ok := proposal.GetJobProposal().(*generated.GetJobProposalJobProposal); ok {
		return success, nil
	}
	return nil, errors.New("failed to get job proposal")
}

func (c *client) ApproveJobProposalSpec(ctx context.Context, id string, force bool) (*JobProposalApprovalSuccessSpec, error) {
	res, err := generated.ApproveJobProposalSpec(ctx, c.gqlClient, id, force)
	if err != nil {
		return nil, err
	}
	if success, ok := res.GetApproveJobProposalSpec().(*generated.ApproveJobProposalSpecApproveJobProposalSpecApproveJobProposalSpecSuccess); ok {
		var cmd JobProposalApprovalSuccessSpec
		if success.Spec.Status == generated.SpecStatusApproved {
			err := DecodeInput(success.Spec, &cmd)
			if err != nil {
				return nil, fmt.Errorf("failed to decode job proposal spec: %w ; and job proposal spec not approved", err)
			}
			return &cmd, nil
		}
	}
	return nil, errors.New("failed to approve job proposal spec")
}

func (c *client) CancelJobProposalSpec(ctx context.Context, id string) (*generated.CancelJobProposalSpecCancelJobProposalSpecCancelJobProposalSpecSuccessSpecJobProposalSpec, error) {
	res, err := generated.CancelJobProposalSpec(ctx, c.gqlClient, id)
	if err != nil {
		return nil, err
	}
	if success, ok := res.GetCancelJobProposalSpec().(*generated.CancelJobProposalSpecCancelJobProposalSpecCancelJobProposalSpecSuccess); ok {
		var cmd generated.CancelJobProposalSpecCancelJobProposalSpecCancelJobProposalSpecSuccessSpecJobProposalSpec
		err := DecodeInput(success.Spec, &cmd)
		if err != nil {
			return nil, fmt.Errorf("failed to decode job proposal spec: %w ; and job proposal spec not cancelled", err)
		}
		if cmd.Status != generated.SpecStatusCancelled {
			return nil, errors.New("job proposal spec not cancelled")
		}
		return &cmd, nil
	}
	return nil, errors.New("failed to cancel job proposal spec")
}

func (c *client) RejectJobProposalSpec(ctx context.Context, id string) (*generated.RejectJobProposalSpecResponse, error) {
	return generated.RejectJobProposalSpec(ctx, c.gqlClient, id)
}

func (c *client) UpdateJobProposalSpecDefinition(ctx context.Context, id string, cmd generated.UpdateJobProposalSpecDefinitionInput) (*generated.UpdateJobProposalSpecDefinitionResponse, error) {
	return generated.UpdateJobProposalSpecDefinition(ctx, c.gqlClient, id, cmd)
}

func (c *client) login() error {
	b, err := json.Marshal(c.credentials)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	payload := strings.NewReader(string(b))

	req, err := http.NewRequest("POST", c.endpoints.Sessions, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status code: %d", res.StatusCode)
	}

	cookieHeader := res.Header.Get("Set-Cookie")
	if cookieHeader != "" {
		c.cookie = strings.Split(cookieHeader, ";")[0]
		return nil
	}

	return fmt.Errorf("no set-cookie found in header. Check credentials and scheme. Response code was: %d", res.StatusCode)
}
