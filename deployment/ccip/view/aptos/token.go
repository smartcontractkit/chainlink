package aptos

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-aptos/bindings/helpers"
	"github.com/smartcontractkit/chainlink-aptos/bindings/managed_token"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	aptosCommon "github.com/smartcontractkit/chainlink/deployment/common/view/aptos"
)

type TokenView struct {
	aptosCommon.ContractMetaData

	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
	Decimals   uint8  `json:"decimals"`
	IconURI    string `json:"iconURI,omitempty"`
	ProjectURI string `json:"projectURI,omitempty"`
	Supply     uint64 `json:"supply"`

	Burners                   []aptos.AccountAddress `json:"burners"`
	Minters                   []aptos.AccountAddress `json:"minters"`
	ManagedTokenObjectAddress string                 `json:"managedTokenObjectAddress,omitempty"`
}

// GenerateTokenView generates a token view for a given managed token.
// The provided address must be a `managed_token` deployment.
func GenerateTokenView(chain cldf_aptos.Chain, managedTokenObjectAddress aptos.AccountAddress) (TokenView, error) {
	objectOwner, err := helpers.GetObjectOwner(chain.Client, managedTokenObjectAddress)
	if err != nil {
		return TokenView{}, err
	}

	boundToken := managed_token.Bind(managedTokenObjectAddress, chain.Client)
	typeAndVersion, err := boundToken.ManagedToken().TypeAndVersion(nil)
	if err != nil {
		return TokenView{}, fmt.Errorf("failed to get typeAndVersion of managedToken %s: %w", managedTokenObjectAddress.StringLong(), err)
	}
	faMetadataAddress, err := boundToken.ManagedToken().TokenMetadata(nil)
	if err != nil {
		return TokenView{}, fmt.Errorf("failed to get tokenMetadata of managedToken %s: %w", managedTokenObjectAddress.StringLong(), err)
	}
	metadata, err := helpers.GetFungibleAssetMetadata(chain.Client, faMetadataAddress)
	if err != nil {
		return TokenView{}, fmt.Errorf("failed to get fungible asset metadata of fungibleAsset %s: %w", faMetadataAddress.StringLong(), err)
	}

	supply, err := helpers.GetFungibleAssetSupply(chain.Client, faMetadataAddress)
	if err != nil {
		return TokenView{}, fmt.Errorf("failed to get fungible asset supply of fungibleAsset %s: %w", faMetadataAddress.StringLong(), err)
	}

	burners, err := boundToken.ManagedToken().GetAllowedBurners(nil)
	if err != nil {
		return TokenView{}, fmt.Errorf("failed to get burners of managedToken %s: %w", managedTokenObjectAddress.StringLong(), err)
	}
	minters, err := boundToken.ManagedToken().GetAllowedMinters(nil)
	if err != nil {
		return TokenView{}, fmt.Errorf("failed to get minters of managedToken %s: %w", managedTokenObjectAddress.StringLong(), err)
	}

	return TokenView{
		ContractMetaData: aptosCommon.ContractMetaData{
			Address:        faMetadataAddress.StringLong(),
			Owner:          objectOwner.StringLong(),
			TypeAndVersion: typeAndVersion,
		},
		Name:                      metadata.Name,
		Symbol:                    metadata.Symbol,
		Decimals:                  metadata.Decimals,
		IconURI:                   metadata.IconURI,
		ProjectURI:                metadata.ProjectURI,
		Supply:                    supply,
		Burners:                   burners,
		Minters:                   minters,
		ManagedTokenObjectAddress: managedTokenObjectAddress.StringLong(),
	}, nil
}
