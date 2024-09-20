package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/Khan/genqlient/graphql"

	"github.com/smartcontractkit/chainlink/integration-tests/web/sdk/client/doer"
	"github.com/smartcontractkit/chainlink/integration-tests/web/sdk/internal/generated"
)

type Client interface {
	FetchCSAPublicKey(ctx context.Context) (*string, error)
	FetchP2PPeerID(ctx context.Context) (*string, error)
	FetchAccountAddress(ctx context.Context, chainID string) (*string, error)
	FetchOCR2KeyBundleID(ctx context.Context, chainType string) (string, error)
	GetJob(ctx context.Context, id string) (*generated.GetJobResponse, error)
	ListJobs(ctx context.Context, offset, limit int) (*generated.ListJobsResponse, error)
	GetJobDistributor(ctx context.Context, id string) (*generated.GetFeedsManagerResponse, error)
	ListJobDistributors(ctx context.Context) (*generated.ListFeedsManagersResponse, error)
	CreateJobDistributor(ctx context.Context, cmd JobDistributorInput) (string, error)
	UpdateJobDistributor(ctx context.Context, id string, cmd JobDistributorInput) error
	CreateJobDistributorChainConfig(ctx context.Context, in JobDistributorChainConfigInput) error
	GetJobProposal(ctx context.Context, id string) (*generated.GetJobProposalResponse, error)
	ApproveJobProposalSpec(ctx context.Context, id string, force bool) (*JobProposalApprovalSuccessSpec, error)
	CancelJobProposalSpec(ctx context.Context, id string) (*generated.CancelJobProposalSpecResponse, error)
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

	if err := c.login(); err != nil {
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
		return nil, fmt.Errorf("no CSA keys found")
	}
	return &keys.CsaKeys.GetResults()[0].PublicKey, nil
}

func (c *client) FetchP2PPeerID(ctx context.Context) (*string, error) {
	keys, err := generated.FetchP2PKeys(ctx, c.gqlClient)
	if err != nil {
		return nil, err
	}
	if keys == nil || len(keys.P2pKeys.GetResults()) == 0 {
		return nil, fmt.Errorf("no P2P keys found")
	}
	return &keys.P2pKeys.GetResults()[0].PeerID, nil
}

func (c *client) FetchOCR2KeyBundleID(ctx context.Context, chainType string) (string, error) {
	keyBundles, err := generated.FetchOCR2KeyBundles(ctx, c.gqlClient)
	if err != nil {
		return "", err
	}
	if keyBundles == nil || len(keyBundles.GetOcr2KeyBundles().Results) == 0 {
		return "", fmt.Errorf("no ocr2 keybundle found, check if ocr2 is enabled")
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
		return nil, fmt.Errorf("no accounts found")
	}
	for _, keyDetail := range keys.EthKeys.GetResults() {
		if keyDetail.GetChain().Enabled && keyDetail.GetChain().Id == chainID {
			return pointer.ToString(keyDetail.Address), nil
		}
	}
	return nil, fmt.Errorf("no account found for chain %s", chainID)
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

func (c *client) GetJobDistributor(ctx context.Context, id string) (*generated.GetFeedsManagerResponse, error) {
	return generated.GetFeedsManager(ctx, c.gqlClient, id)
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
	return "", fmt.Errorf("failed to create feeds manager")
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

func (c *client) CreateJobDistributorChainConfig(ctx context.Context, in JobDistributorChainConfigInput) error {
	var cmd generated.CreateFeedsManagerChainConfigInput
	err := DecodeInput(in, &cmd)
	if err != nil {
		return err
	}
	_, err = generated.CreateFeedsManagerChainConfig(ctx, c.gqlClient, cmd)
	return err
}

func (c *client) GetJobProposal(ctx context.Context, id string) (*generated.GetJobProposalResponse, error) {
	return generated.GetJobProposal(ctx, c.gqlClient, id)
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
	return nil, fmt.Errorf("failed to approve job proposal spec")
}

func (c *client) CancelJobProposalSpec(ctx context.Context, id string) (*generated.CancelJobProposalSpecResponse, error) {
	return generated.CancelJobProposalSpec(ctx, c.gqlClient, id)
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

	cookieHeader := res.Header.Get("Set-Cookie")
	if cookieHeader == "" {
		return fmt.Errorf("no cookie found in header")
	}

	c.cookie = strings.Split(cookieHeader, ";")[0]
	return nil
}
