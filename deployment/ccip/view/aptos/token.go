package aptos

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/managed_token"
	"github.com/smartcontractkit/chainlink-aptos/relayer/codec"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	aptosCommon "github.com/smartcontractkit/chainlink/deployment/common/view/aptos"
)

type TokenView struct {
	aptosCommon.ContractMetaData

	Name       string `json:"name,omitempty"`
	Symbol     string `json:"symbol,omitempty"`
	Decimals   uint8  `json:"decimals"`
	IconURI    string `json:"iconURI,omitempty"`
	ProjectURI string `json:"projectURI,omitempty"`
	Supply     uint64 `json:"supply"`

	Burners []aptos.AccountAddress `json:"burners,omitempty"`
	Minters []aptos.AccountAddress `json:"minters,omitempty"`
}

func GenerateTokenView(chain cldf_aptos.Chain, tokenAddress aptos.AccountAddress) (TokenView, error) {
	faOwner, err := GetObjectOwner(chain.Client, tokenAddress)
	if err != nil {
		return TokenView{}, err
	}
	tokenObject, err := GetObjectOwner(chain.Client, faOwner)
	if err != nil {
		return TokenView{}, err
	}
	tokenOwner, err := GetObjectOwner(chain.Client, tokenObject)
	if err != nil {
		return TokenView{}, err
	}

	boundToken := managed_token.Bind(tokenObject, chain.Client)
	faMetadataAddress, err := boundToken.ManagedToken().TokenMetadata(nil)
	if err != nil {
		return TokenView{}, fmt.Errorf("failed to retrieve managed token metadata address: %w", err)
	}
	metadata, err := GetFungibleAssetMetadata(chain.Client, faMetadataAddress)
	if err != nil {
		return TokenView{}, fmt.Errorf("failed to retrieve fungible asset (addr %v) metadata: %w", faMetadataAddress.StringLong(), err)
	}

	supply, err := GetFungibleAssetSupply(chain.Client, faMetadataAddress)
	if err != nil {
		return TokenView{}, fmt.Errorf("failed to retrieve fungible asset (addr %v) supply: %w", faMetadataAddress.StringLong(), err)
	}

	burners, err := boundToken.ManagedToken().GetAllowedBurners(nil)
	if err != nil {
		return TokenView{}, fmt.Errorf("failed to retrieve managed token burners: %w", err)
	}
	minters, err := boundToken.ManagedToken().GetAllowedMinters(nil)
	if err != nil {
		return TokenView{}, fmt.Errorf("failed to retrieve managed token minters: %w", err)
	}

	return TokenView{
		ContractMetaData: aptosCommon.ContractMetaData{
			Address: tokenAddress.StringLong(),
			Owner:   tokenOwner.StringLong(),
		},
		Name:       metadata.Name,
		Symbol:     metadata.Symbol,
		Decimals:   metadata.Decimals,
		IconURI:    metadata.IconURI,
		ProjectURI: metadata.ProjectURI,
		Supply:     supply,
		Burners:    burners,
		Minters:    minters,
	}, nil
}

// TODO: Move this to chainlink-aptos

func GetObjectOwner(
	client aptos.AptosRpcClient,
	objectAddress aptos.AccountAddress,
) (aptos.AccountAddress, error) {
	bc := bind.NewBoundContract(
		aptos.AccountOne,
		"std",
		"object",
		client,
	)
	module, function, typeTags, args, err := bc.Encode(
		"owner",
		[]string{
			"0x1::object::ObjectCore",
		},
		[]string{
			"address",
		}, []any{
			objectAddress,
		})
	callData, err := bc.Call(nil, module, function, typeTags, args)
	if err != nil {
		return aptos.AccountAddress{}, err
	}

	var owner aptos.AccountAddress
	if err := codec.DecodeAptosJsonArray(callData, &owner); err != nil {
		return aptos.AccountAddress{}, err
	}
	return owner, nil
}

type FungibleAssetMetadata struct {
	Name       string
	Symbol     string
	Decimals   uint8
	IconURI    string
	ProjectURI string
}

func GetFungibleAssetMetadata(
	client aptos.AptosRpcClient,
	faMetadataAddress aptos.AccountAddress,
) (FungibleAssetMetadata, error) {
	bc := bind.NewBoundContract(
		aptos.AccountOne,
		"std",
		"fungible_asset",
		client,
	)
	module, function, typeTags, args, err := bc.Encode(
		"metadata",
		[]string{
			"0x1::fungible_asset::Metadata",
		},
		[]string{
			"address",
		}, []any{
			faMetadataAddress,
		})
	callData, err := bc.Call(nil, module, function, typeTags, args)
	if err != nil {
		return FungibleAssetMetadata{}, err
	}

	var metadata FungibleAssetMetadata
	if err := codec.DecodeAptosJsonArray(callData, &metadata); err != nil {
		return FungibleAssetMetadata{}, err
	}
	return metadata, nil
}

func GetFungibleAssetSupply(
	client aptos.AptosRpcClient,
	faMetadataAddress aptos.AccountAddress,
) (uint64, error) {
	bc := bind.NewBoundContract(
		aptos.AccountOne,
		"std",
		"fungible_asset",
		client,
	)
	module, function, typeTags, args, err := bc.Encode(
		"supply",
		[]string{
			"0x1::fungible_asset::Metadata",
		},
		[]string{
			"address",
		}, []any{
			faMetadataAddress,
		})
	if err != nil {
		return 0, err
	}
	callData, err := bc.Call(nil, module, function, typeTags, args)
	if err != nil {
		return 0, err
	}

	var supply bind.StdOption[uint64]
	if err := codec.DecodeAptosJsonArray(callData, &supply); err != nil {
		return 0, err
	}
	return *supply.Value(), nil
}
