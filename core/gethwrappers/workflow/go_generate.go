// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Workflow

//go:generate go run ../generation/wrap.go workflow WorkflowRegistry workflow_registry_wrapper
