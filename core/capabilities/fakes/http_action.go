package fakes

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commonCap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	http "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/http"
	httpserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/http/server"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var _ httpserver.ClientCapability = (*FakeHttpAction)(nil)
var _ services.Service = (*FakeHttpAction)(nil)
var _ commonCap.ExecutableCapability = (*FakeHttpAction)(nil)

const HTTPActionID = "http-action@1.0.0"
const HTTPActionServiceName = "HttpActionService"

var fakeHttpActionInfo = capabilities.MustNewCapabilityInfo(
	HTTPActionID,
	capabilities.CapabilityTypeAction,
	"An action that uses an HTTP request to run periodically at fixed times, dates, or intervals.",
)

type FakeHttpAction struct {
	commonCap.CapabilityInfo
	services.Service
	eng *services.Engine

	lggr logger.Logger
}

func (fh *FakeHttpAction) SendRequest(ctx context.Context, metadata capabilities.RequestMetadata, input *http.Request) (*http.Response, error) {
	fh.eng.Infow("Fake Http Action SendRequest Started", "input", input)
	return nil, nil
}

func (fh *FakeHttpAction) Start(ctx context.Context) error {
	fh.eng.Infow("Fake Http Action Start Started")
	return nil
}

func (fh *FakeHttpAction) Close() error {
	fh.eng.Infow("Fake Http Action Close Started")
	return nil
}

func (fh *FakeHttpAction) Name() string {
	return HTTPActionServiceName
}

func (fh *FakeHttpAction) Description() string {
	return fakeHttpActionInfo.Description
}

func (fh *FakeHttpAction) Ready() error {
	return nil
}

func (fh *FakeHttpAction) Initialise(ctx context.Context, config string, _ core.TelemetryService,
	_ core.KeyValueStore,
	_ core.ErrorLog,
	_ core.PipelineRunnerService,
	_ core.RelayerSet,
	_ core.OracleFactory) error {

	// TODO: do validation of config here

	err := fh.Start(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (fh *FakeHttpAction) Execute(ctx context.Context, request commonCap.CapabilityRequest) (commonCap.CapabilityResponse, error) {
	fh.eng.Infow("Fake Http Action Execute Started", "request", request)
	return commonCap.CapabilityResponse{}, nil
}

func (fh *FakeHttpAction) RegisterToWorkflow(ctx context.Context, request commonCap.RegisterToWorkflowRequest) error {
	fh.eng.Infow("Registered to Fake Http Action", "workflowID", request.Metadata.WorkflowID)
	return nil
}

func (fh *FakeHttpAction) UnregisterFromWorkflow(ctx context.Context, request commonCap.UnregisterFromWorkflowRequest) error {
	fh.eng.Infow("Unregistered from Fake Http Action", "workflowID", request.Metadata.WorkflowID)
	return nil
}
