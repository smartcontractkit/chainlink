package framework

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering"
)

type fakeBillingClient struct {
}

func NewFakeBillingClient() metering.BillingClient {
	return &fakeBillingClient{}
}

func (f fakeBillingClient) GetOrganizationCreditsByWorkflow(ctx context.Context, req *billing.GetOrganizationCreditsByWorkflowRequest) (*billing.GetOrganizationCreditsByWorkflowResponse, error) {
	return &billing.GetOrganizationCreditsByWorkflowResponse{OrganizationId: "", Credits: &billing.OrganizationCredits{CreditsReserved: "", Credits: ""}}, nil
}

func (f fakeBillingClient) GetWorkflowExecutionRates(context.Context, *billing.GetWorkflowExecutionRatesRequest) (*billing.GetWorkflowExecutionRatesResponse, error) {
	return &billing.GetWorkflowExecutionRatesResponse{RateCards: []*billing.RateCard{{ResourceType: billing.ResourceType_RESOURCE_TYPE_COMPUTE, MeasurementUnit: billing.MeasurementUnit_MEASUREMENT_UNIT_MILLISECONDS, UnitsPerCredit: "0.0001"}}}, nil
}

func (f fakeBillingClient) SubmitWorkflowReceipt(ctx context.Context, request *billing.SubmitWorkflowReceiptRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (f fakeBillingClient) ReserveCredits(ctx context.Context, request *billing.ReserveCreditsRequest) (*billing.ReserveCreditsResponse, error) {
	return &billing.ReserveCreditsResponse{Success: true, RateCards: []*billing.RateCard{{ResourceType: billing.ResourceType_RESOURCE_TYPE_COMPUTE, MeasurementUnit: billing.MeasurementUnit_MEASUREMENT_UNIT_MILLISECONDS, UnitsPerCredit: "0.0001"}}, Credits: "10000"}, nil
}
