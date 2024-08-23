// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Chainlink Operator Forwarder contracts
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/AuthorizedForwarder/AuthorizedForwarder.abi ../../../contracts/solc/v0.8.19/AuthorizedForwarder/AuthorizedForwarder.bin AuthorizedForwarder authorized_forwarder
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/AuthorizedReceiver/AuthorizedReceiver.abi ../../../contracts/solc/v0.8.19/AuthorizedReceiver/AuthorizedReceiver.bin AuthorizedReceiver authorized_receiver
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/LinkTokenReceiver/LinkTokenReceiver.abi ../../../contracts/solc/v0.8.19/LinkTokenReceiver/LinkTokenReceiver.bin LinkTokenReceiver link_token_receiver
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/Operator/Operator.abi ../../../contracts/solc/v0.8.19/Operator/Operator.bin Operator operator
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/OperatorFactory/OperatorFactory.abi ../../../contracts/solc/v0.8.19/OperatorFactory/OperatorFactory.bin OperatorFactory operator_factory
