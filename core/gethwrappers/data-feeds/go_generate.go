// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Chainlink Data Feeds

//go:generate go run ../generation/wrap.go data-feeds BundleAggregatorProxy bundle_aggregator_proxy
//go:generate go run ../generation/wrap.go data-feeds DataFeedsCache data_feeds_cache
