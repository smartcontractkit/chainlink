// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Workflow

//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/WorkflowRegistry/WorkflowRegistry.abi ../../../contracts/solc/v0.8.24/WorkflowRegistry/WorkflowRegistry.bin WorkflowRegistry workflow_registry_wrapper
