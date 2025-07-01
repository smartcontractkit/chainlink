package framework

import (
	"context"

	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering"
)

type fakeBillingClient struct {
}

func NewFakeBillingClient() metering.BillingClient {
	return &fakeBillingClient{}
}

func (f fakeBillingClient) SubmitWorkflowReceipt(ctx context.Context, request *billing.SubmitWorkflowReceiptRequest) (*billing.SubmitWorkflowReceiptResponse, error) {
	return &billing.SubmitWorkflowReceiptResponse{Success: true}, nil
}

func (f fakeBillingClient) ReserveCredits(ctx context.Context, request *billing.ReserveCreditsRequest) (*billing.ReserveCreditsResponse, error) {
	return &billing.ReserveCreditsResponse{Success: true, Rates: []*billing.ResourceUnitRate{{ResourceUnit: metering.ComputeResourceDimension, ConversionRate: "0.0001"}}, Credits: 10000}, nil
}
