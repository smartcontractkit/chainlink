package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/smartcontractkit/ccip/integration-tests/web/sdk/client/internal/doer"
	"github.com/smartcontractkit/ccip/integration-tests/web/sdk/internal/generated"
)

type Client interface {
	GetCSAKeys(ctx context.Context) (*generated.GetCSAKeysResponse, error)
	GetJob(ctx context.Context, id string) (*generated.GetJobResponse, error)
	ListJobs(ctx context.Context, offset, limit int) (*generated.ListJobsResponse, error)
	GetBridge(ctx context.Context, id string) (*generated.GetBridgeResponse, error)
	ListBridges(ctx context.Context, offset, limit int) (*generated.ListBridgesResponse, error)
	GetFeedsManager(ctx context.Context, id string) (*generated.GetFeedsManagerResponse, error)
	ListFeedsManagers(ctx context.Context) (*generated.ListFeedsManagersResponse, error)
	CreateFeedsManager(ctx context.Context, cmd generated.CreateFeedsManagerInput) (*generated.CreateFeedsManagerResponse, error)
	UpdateFeedsManager(ctx context.Context, id string, cmd generated.UpdateFeedsManagerInput) (*generated.UpdateFeedsManagerResponse, error)
	GetJobProposal(ctx context.Context, id string) (*generated.GetJobProposalResponse, error)
	ApproveJobProposalSpec(ctx context.Context, id string, force bool) (*generated.ApproveJobProposalSpecResponse, error)
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

func (c *client) GetCSAKeys(ctx context.Context) (*generated.GetCSAKeysResponse, error) {
	return generated.GetCSAKeys(ctx, c.gqlClient)
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

func (c *client) GetFeedsManager(ctx context.Context, id string) (*generated.GetFeedsManagerResponse, error) {
	return generated.GetFeedsManager(ctx, c.gqlClient, id)
}

func (c *client) ListFeedsManagers(ctx context.Context) (*generated.ListFeedsManagersResponse, error) {
	return generated.ListFeedsManagers(ctx, c.gqlClient)
}

func (c *client) CreateFeedsManager(ctx context.Context, cmd generated.CreateFeedsManagerInput) (*generated.CreateFeedsManagerResponse, error) {
	return generated.CreateFeedsManager(ctx, c.gqlClient, cmd)
}

func (c *client) UpdateFeedsManager(ctx context.Context, id string, cmd generated.UpdateFeedsManagerInput) (*generated.UpdateFeedsManagerResponse, error) {
	return generated.UpdateFeedsManager(ctx, c.gqlClient, id, cmd)
}

func (c *client) GetJobProposal(ctx context.Context, id string) (*generated.GetJobProposalResponse, error) {
	return generated.GetJobProposal(ctx, c.gqlClient, id)
}

func (c *client) ApproveJobProposalSpec(ctx context.Context, id string, force bool) (*generated.ApproveJobProposalSpecResponse, error) {
	return generated.ApproveJobProposalSpec(ctx, c.gqlClient, id, force)
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
