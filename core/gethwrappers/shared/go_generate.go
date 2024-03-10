// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/BurnMintERC677/BurnMintERC677.abi ../../../contracts/solc/v0.8.19/BurnMintERC677/BurnMintERC677.bin BurnMintERC677 burn_mint_erc677
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/LinkToken/LinkToken.abi ../../../contracts/solc/v0.8.19/LinkToken/LinkToken.bin LinkToken link_token
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/ERC20/ERC20.abi ../../../contracts/solc/v0.8.19/ERC20/ERC20.bin ERC20 erc20
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/WERC20Mock/WERC20Mock.abi ../../../contracts/solc/v0.8.19/WERC20Mock/WERC20Mock.bin WERC20Mock werc20_mock
