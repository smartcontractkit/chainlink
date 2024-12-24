// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

//go:generate go run ../generation/generate/wrap_short.go shared BurnMintERC677 burn_mint_erc677
//go:generate go run ../generation/generate/wrap_short.go shared LinkToken link_token
//go:generate go run ../generation/generate/wrap_short.go shared BurnMintERC20 burn_mint_erc20
//go:generate go run ../generation/generate/wrap_short.go shared WERC20Mock werc20_mock

//go:generate go run ../generation/generate/wrap_short.go vendor ERC20 erc20
//go:generate go run ../generation/generate/wrap_short.go vendor Multicall3 multicall3
//go:generate go run ../generation/generate/wrap_short.go tests MockV3Aggregator mock_v3_aggregator_contract
