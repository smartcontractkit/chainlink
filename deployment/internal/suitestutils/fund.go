package suitestutils

import (
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/go-resty/resty/v2"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// FundAccount funds a Sui account with the given amount of SUI.
func FundAccount(url string, address string) error {
	r := resty.New().SetBaseURL(url)
	b := &models.FaucetRequest{
		FixedAmountRequest: &models.FaucetFixedAmountRequest{
			Recipient: address,
		},
	}
	resp, err := r.R().SetBody(b).SetHeader("Content-Type", "application/json").Post("/gas")
	if err != nil {
		return err
	}
	framework.L.Info().Any("Resp", resp).Msg("Address is funded!")

	return nil
}
