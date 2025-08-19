# Shared Peer Abstraction

This package provides a unified peer abstraction that can work with both OCR (Off-Chain Reporting) and non-OCR P2P networking contexts, leveraging the `PeerGroupFactory` from `SingletonPeerWrapper`.

## Overview

The `SharedPeerAbstraction` provides a single interface that can:

1. **OCR Mode**: Use the `SingletonPeerWrapper` for OCR-enabled applications, providing access to `PeerGroupFactory` for creating OCR peer groups
2. **External Peer Mode**: Use the direct P2P peer implementation for non-OCR applications
3. **Flexible Configuration**: Allow runtime selection between OCR and external peer modes based on configuration

## Architecture

```
SharedPeerAbstraction
├── OCR Mode (SingletonPeerWrapper)
│   ├── PeerGroupFactory (for OCR operations)
│   ├── OCR1/2/3 adapters
│   └── Cryptographic signing
└── External Peer Mode
    ├── Direct P2P peer
    ├── Stream management
    └── Cryptographic signing
```

## Key Components

### SharedPeerAbstraction Interface

```go
type SharedPeerAbstraction interface {
    services.Service
    
    // GetPeer returns the native P2P peer for direct communication
    GetPeer() p2ptypes.Peer
    
    // GetPeerGroupFactory returns the OCR PeerGroupFactory if available
    GetPeerGroupFactory() ocrnetworking.PeerGroupFactory
    
    // HasOCRCapability indicates if this abstraction supports OCR functionality
    HasOCRCapability() bool
    
    // GetSingletonWrapper returns the underlying SingletonPeerWrapper if available
    GetSingletonWrapper() *ocrcommon.SingletonPeerWrapper
    
    // Sign provides cryptographic signing capability
    Sign(data []byte) ([]byte, error)
}
```

### Configuration

```go
type SharedPeerConfig struct {
    KeyStore      keystore.Master
    KeystoreP2P   keystore.P2P
    P2PConfig     config.P2P
    OCRConfig     ocrcommon.PeerWrapperOCRConfig
    DataSource    sqlutil.DataSource
    Logger        logger.Logger
    
    // EnableOCR determines if OCR functionality should be enabled
    EnableOCR     bool
    
    // ForceExternal forces the use of external peer even when OCR is available
    ForceExternal bool
}
```

## Usage Examples

### Basic Usage

```go
// Create a shared peer abstraction
sharedPeer := p2p.NewSharedPeerAbstraction(p2p.SharedPeerConfig{
    KeyStore:      keyStore,
    KeystoreP2P:   keystoreP2P,
    P2PConfig:     p2pConfig,
    OCRConfig:     ocrConfig,
    DataSource:    ds,
    Logger:        lggr,
    EnableOCR:     true,  // Enable OCR capabilities when available
    ForceExternal: false, // Allow automatic selection
})

// Start the abstraction
ctx := context.Background()
if err := sharedPeer.Start(ctx); err != nil {
    return fmt.Errorf("failed to start shared peer: %w", err)
}
defer sharedPeer.Close()

// Check what capabilities are available
if sharedPeer.HasOCRCapability() {
    // Use OCR functionality
    peerGroupFactory := sharedPeer.GetPeerGroupFactory()
    // Use peerGroupFactory for OCR operations...
} else {
    // Use direct P2P functionality
    peer := sharedPeer.GetPeer()
    // Use peer for direct P2P operations...
}
```

### OCR-Only Mode

```go
sharedPeer := p2p.NewSharedPeerAbstraction(p2p.SharedPeerConfig{
    KeyStore:      keyStore,
    KeystoreP2P:   keyStore.P2P(),
    P2PConfig:     p2pConfig,
    OCRConfig:     ocrConfig,
    DataSource:    ds,
    Logger:        lggr,
    EnableOCR:     true,
    ForceExternal: false,
})

if err := sharedPeer.Start(ctx); err != nil {
    return fmt.Errorf("failed to start OCR peer: %w", err)
}

if !sharedPeer.HasOCRCapability() {
    return fmt.Errorf("OCR capability not available")
}

// Access OCR-specific functionality
peerGroupFactory := sharedPeer.GetPeerGroupFactory()
singletonWrapper := sharedPeer.GetSingletonWrapper()
```

### External Peer Only Mode

```go
sharedPeer := p2p.NewSharedPeerAbstraction(p2p.SharedPeerConfig{
    KeyStore:      nil, // Not needed for external peer only
    KeystoreP2P:   keystoreP2P,
    P2PConfig:     p2pConfig,
    OCRConfig:     nil, // Not needed for external peer only
    DataSource:    ds,
    Logger:        lggr,
    EnableOCR:     false, // Disable OCR
    ForceExternal: false,
})

if err := sharedPeer.Start(ctx); err != nil {
    return fmt.Errorf("failed to start external peer: %w", err)
}

peer := sharedPeer.GetPeer()
if peer == nil {
    return fmt.Errorf("external peer not available")
}

// Use direct P2P operations
peerID := peer.ID()
```

## Mode Selection Logic

The abstraction automatically selects the appropriate mode based on configuration:

1. **OCR Mode**: Selected when `EnableOCR=true` and `ForceExternal=false`
   - Uses `SingletonPeerWrapper`
   - Provides `PeerGroupFactory` access
   - Supports OCR1, OCR2, and OCR3 protocols

2. **External Peer Mode**: Selected when `EnableOCR=false` or `ForceExternal=true`
   - Uses direct P2P peer implementation
   - Provides direct peer access for messaging
   - Independent of OCR infrastructure

## Features

### Common Features (Both Modes)
- Cryptographic signing capabilities
- Service lifecycle management (Start/Stop)
- Health reporting and readiness checks
- Logging and monitoring

### OCR Mode Features
- Access to `PeerGroupFactory` for creating OCR peer groups
- Integration with existing OCR infrastructure
- Support for multiple OCR protocol versions
- Automatic peer discovery and management

### External Peer Mode Features
- Direct P2P messaging capabilities
- Stream management for peer connections
- Custom peer discovery configuration
- Lower overhead for non-OCR applications

## Testing

The package includes comprehensive tests:

- Unit tests for basic functionality and configuration
- Integration tests for both OCR and external peer modes
- Mock-based tests for isolated component testing
- Error handling and edge case testing

Run tests with:
```bash
go test ./core/services/p2p/ -v
```

## Integration with Existing Code

The shared peer abstraction is designed to integrate seamlessly with existing Chainlink components:

1. **Replaces direct usage** of `SingletonPeerWrapper` where both OCR and non-OCR functionality might be needed
2. **Provides backward compatibility** through access to underlying implementations
3. **Simplifies configuration** by providing a single interface for both modes
4. **Enables runtime mode selection** based on application requirements

## Future Enhancements

Potential future improvements include:

1. **Dynamic mode switching** - Allow runtime switching between OCR and external peer modes
2. **Multi-peer support** - Support multiple peer instances for different networks
3. **Enhanced metrics** - Provide mode-specific metrics and monitoring
4. **Configuration validation** - Enhanced validation for different mode requirements
