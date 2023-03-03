/*
Package ocr2keepers provides an implementation of the OCR2 oracle plugin.
Sub-packages include chain specific configurations starting with EVM based
chains. To create a new Delegate that can start and stop a Keepers OCR2 Oracle
plugin, run the following:

	del, err := ocr2keepers.NewDelegate(config)

# Multi-chain Support

Chain specific supported can be added by providing implementations for both
Registry and ReportEncoder. A registry is used to collect all upkeeps registered
and to check the perform status of each. A report encoder produces a byte array
from a set of upkeeps to be performed such that the resulting bytes can be
transacted on chain.

# Types

Most types are wrappers for byte arrays, but their internal structure is
important when creating new Registry and ReportEncoder implementations. It is
assumed that a block on any chain has an identifier, be it numeric or textual.
A BlockKey wraps a byte array for the reason that each implementation handle
encoding of the block identifier internally.

Likewise, an UpkeepKey is a wrapper for a byte array such that the encoded data
should be the combination of both a block and upkeep id. In most chains, an
upkeep id will be numeric (*big.Big), but is not required to be. It is up to
the implementation for detail. The main idea is that the plugin assumes an
UpkeepKey to be an upkeep id at a specific point in time on a block chain.
*/
package ocr2keepers
