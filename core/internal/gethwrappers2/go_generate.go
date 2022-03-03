// Package gethwrappers keeps track of the golang wrappers of the solidity contracts
package main

//go:generate ./compile.sh 15000 ../../../../libocr-internal/contract2/AccessControlledOCR2Aggregator.sol
//go:generate ./compile.sh 15000 ../../../../libocr-internal/contract2/OCR2Aggregator.sol
//go:generate ./compile.sh 15000 ../../../../libocr-internal/contract2/ExposedOCR2Aggregator.sol

//go:generate ./compile.sh 1000 ../../../../libocr-internal/contract2/TestOCR2Aggregator.sol
//go:generate ./compile.sh 1000 ../../../../libocr-internal/contract2/TestValidator.sol
//go:generate ./compile.sh 1000 ../../../../libocr-internal/contract2/AccessControlTestHelper.sol
