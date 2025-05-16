package main

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/verification/evm"
)

func main() {
	/*
		verifier, err := evm.NewEtherscanContractVerifier(
			"https://api-sepolia.etherscan.io/api",
			"BUNTKQQYYN1AMDCNNMV8SFN57R6CQGQHVS",
			"0x1B4284a86Cc0f3ac975980DD5D951b8456fa7C36",
			changeset.BurnMintTokenPool,
			deployment.Version1_5_1,
			3*time.Second,
		)
		if err != nil {
			panic(err)
		}
		verified, err := verifier.Verify(context.Background())
		if err != nil {
			panic(err)
		}
		if !verified {
			panic("contract not verified")
		}
	*/
	lggr, err := logger.New()
	if err != nil {
		panic(err)
	}
	verifier, err := evm.NewEtherscanContractVerifier(
		"https://api.sonicscan.org/api",
		"NT5Q2Y87EB9BTP4K6RITT3P1866DKUJMNT",
		"0x006bC1F599a10B73C88cc3cD19a92829C4AC1E83",
		shared.FeeQuoter,
		deployment.Version1_6_0,
		5*time.Second,
		lggr,
	)
	if err != nil {
		panic(err)
	}
	err = verifier.Verify(context.Background())
	if err != nil {
		panic(err)
	}
	verified, err := verifier.IsVerified(context.Background())
	if err != nil {
		panic(err)
	}
	if !verified {
		panic("contract not verified")
	} else {
		println("contract verified")
	}
}

/*
result=gpqnm8ddz3svmjhf1niyuqm2kbrthnz7dnrx6ash16plnqumvf
result=Contract source code already verified

FeeQuoter issues
Router issues
*/
