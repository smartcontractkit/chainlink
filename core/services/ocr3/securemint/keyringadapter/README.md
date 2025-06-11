# PoR OCR3 OnchainKeyring Adapter

This file contains an adapter implementation that enables the use of existing OCR2 OnchainKeyring implementations with the OCR3 PoR (Proof of Reserve) plugin.

## Overview

The `OnchainKeyringAdapter` wraps an existing `types.OnchainKeyring` (OCR2) and adapts it to implement `ocr3types.OnchainKeyring[ChainSelector]` (OCR3) specifically for the PoR system.

## Key Features

- **Interface Adaptation**: Converts between OCR2 and OCR3 keyring interfaces
- **Parameter Conversion**: Automatically converts OCR3 parameters (config digest, sequence number, report with info) to OCR2 ReportContext format
- **Backward Compatibility**: Allows reuse of existing OCR2 keyring implementations
- **Type Safety**: Strongly typed for PoR ChainSelector

## Usage Example

```go
package main

import (
    "github.com/smartcontractkit/por_mock_ocr3plugin/por"
    "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
    "github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
)

func main() {
    // Create an OCR2 keyring (example using EVM keyring)
    ocr2Bundle, err := ocr2key.New(chaintype.EVM)
    if err != nil {
        panic(err)
    }

    // Wrap the OCR2 keyring with the PoR adapter
    porKeyring := por.NewOnchainKeyringAdapter(ocr2Bundle)

    // Now you can use porKeyring as an ocr3types.OnchainKeyring[por.ChainSelector]
    // for the PoR OCR3 plugin
    
    // Example usage in OCR3 context:
    configDigest := types.ConfigDigest([32]byte{1, 2, 3})
    seqNr := uint64(42)
    reportWithInfo := ocr3types.ReportWithInfo[por.ChainSelector]{
        Report: []byte("example-report"),
        Info:   por.ChainSelector(1234),
    }
    
    signature, err := porKeyring.Sign(configDigest, seqNr, reportWithInfo)
    if err != nil {
        panic(err)
    }
    
    isValid := porKeyring.Verify(
        porKeyring.PublicKey(),
        configDigest,
        seqNr,
        reportWithInfo,
        signature,
    )
    
    println("Signature valid:", isValid)
}
```

## Implementation Details

### Interface Mapping

The adapter maps OCR3 interface methods to OCR2 interface methods as follows:

| OCR3 Method | OCR2 Method | Conversion |
|-------------|-------------|------------|
| `PublicKey()` | `PublicKey()` | Direct passthrough |
| `Sign(ConfigDigest, uint64, ReportWithInfo[ChainSelector])` | `Sign(ReportContext, Report)` | Converts parameters to ReportContext |
| `Verify(OnchainPublicKey, ConfigDigest, uint64, ReportWithInfo[ChainSelector], []byte)` | `Verify(OnchainPublicKey, ReportContext, Report, []byte)` | Converts parameters to ReportContext |
| `MaxSignatureLength()` | `MaxSignatureLength()` | Direct passthrough |

### Parameter Conversion

OCR3 parameters are converted to OCR2 ReportContext as follows:

```go
reportContext := types.ReportContext{
    ReportTimestamp: types.ReportTimestamp{
        ConfigDigest: configDigest,     // From OCR3 parameter
        Epoch:        uint32(seqNr),    // OCR3 sequence number as epoch
        Round:        0,                // Fixed to 0 (OCR3 doesn't use rounds)
    },
    ExtraHash: [32]byte{},             // Empty hash
}
```

## Benefits

1. **Reusability**: Existing OCR2 keyrings can be used with OCR3 PoR without modification
2. **Simplicity**: Single adapter handles all necessary conversions
3. **Type Safety**: Generic implementation ensures compile-time type checking
4. **Testing**: Comprehensive test suite ensures correct parameter conversion

## Related Files

- `onchain_keyring_adapter.go` - Main adapter implementation
- `onchain_keyring_adapter_test.go` - Comprehensive test suite
- `types.go` - PoR-specific type definitions including ChainSelector
- `external_adapter_interface.go` - Original external adapter interface

## See Also

- Reference implementations in `/core/services/ocrcommon/adapters.go`
- OCR2 keyring implementations in `/core/services/keystore/keys/ocr2key/`
- OCR3 types in `libocr/offchainreporting2plus/ocr3types/`
