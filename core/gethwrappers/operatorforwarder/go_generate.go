// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Chainlink Operator Forwarder contracts
//go:generate go run ../generation/wrap.go operatorforwarder AuthorizedForwarder authorized_forwarder
//go:generate go run ../generation/wrap.go operatorforwarder AuthorizedReceiver authorized_receiver
//go:generate go run ../generation/wrap.go operatorforwarder LinkTokenReceiver link_token_receiver
//go:generate go run ../generation/wrap.go operatorforwarder Operator operator
//go:generate go run ../generation/wrap.go operatorforwarder OperatorFactory operator_factory
