# Multi-Chain Pattern for Secure Mint Plugin

This document describes how to implement a multi-chain reading and writing pattern for the Secure Mint plugin, similar to how CCIP handles multiple chains.

## Overview

The multi-chain pattern allows the Secure Mint plugin to:
- Read from multiple blockchain networks simultaneously
- Write to multiple blockchain networks with separate transaction managers
- Manage chain-specific configurations and gas settings
- Provide health monitoring and error handling per chain

## Architecture

### 1. Multi-Chain Services Structure

```go
type MultiChainSecureMintServices struct {
    services.StateMachine
    
    // Map of chain selectors to their respective contract readers
    contractReaders map[por.ChainSelector]sm_plugin.ContractReader
    
    // Map of chain selectors to their respective contract transmitters
    contractTransmitters map[por.ChainSelector]ocr3types.ContractTransmitter[por.ChainSelector]
    
    // Map of chain selectors to their respective chain writers
    chainWriters map[por.ChainSelector]types.ContractWriter
    
    // Map of relay IDs to relayers for different chains
    relayers map[types.RelayID]loop.Relayer
    
    // Configuration for each chain
    chainConfigs map[por.ChainSelector]*ChainConfig
    
    mu sync.RWMutex
}
```

### 2. Chain Configuration

Each chain has its own configuration:

```go
type ChainConfig struct {
    ChainSelector    por.ChainSelector
    RelayID          types.RelayID
    ContractAddress  string
    GasLimit         uint64
    MaxGasPrice      string
    FromAddress      string
}
```

## Key Components

### 1. Contract Reader Per Chain

Each chain has its own contract reader that wraps the chain reader interface:

```go
type secureMintContractReader struct {
    chainReader   types.ContractReader
    chainSelector por.ChainSelector
    logger        logger.Logger
}
```

**Benefits:**
- Chain-specific error handling
- Independent health monitoring
- Isolated configuration per chain

### 2. Contract Transmitter Per Chain

Each chain has its own contract transmitter that uses a chain writer:

```go
type secureMintContractTransmitter struct {
    logger         logger.Logger
    chainWriter    types.ContractWriter
    fromAccount    ocrtypes.Account
    contractAddress string
    chainSelector  por.ChainSelector
}
```

**Benefits:**
- Chain-specific gas settings
- Independent transaction management
- Separate retry logic per chain

### 3. Chain Writer Per Chain

Each chain has its own chain writer for submitting transactions:

```go
// Each chain writer manages transactions for its specific chain
chainWriters map[por.ChainSelector]types.ContractWriter
```

**Benefits:**
- Chain-specific transaction management
- Independent gas price monitoring
- Separate nonce management

## Usage Pattern

### 1. Service Initialization

```go
// Create multi-chain services
multiChainServices := NewMultiChainSecureMintServices(
    logger.Named("MultiChain"),
    relayers,
    chainConfigs,
)

// Start all chain-specific services
if err := multiChainServices.Start(ctx); err != nil {
    return fmt.Errorf("failed to start multi-chain services: %w", err)
}
```

### 2. Chain-Specific Operations

```go
// Get contract reader for a specific chain
reader, err := multiChainServices.GetContractReader(chainSelector)
if err != nil {
    return fmt.Errorf("failed to get contract reader: %w", err)
}

// Get contract transmitter for a specific chain
transmitter, err := multiChainServices.GetContractTransmitter(chainSelector)
if err != nil {
    return fmt.Errorf("failed to get contract transmitter: %w", err)
}

// Get chain writer for a specific chain
writer, err := multiChainServices.GetChainWriter(chainSelector)
if err != nil {
    return fmt.Errorf("failed to get chain writer: %w", err)
}
```

### 3. Health Monitoring

```go
// Get health status for all chains
healthReport := multiChainServices.HealthReport()
for chainName, err := range healthReport {
    if err != nil {
        logger.Errorw("Chain health issue", "chain", chainName, "error", err)
    }
}
```

## Configuration

### 1. Multi-Chain Configuration

```go
// Example configuration for multiple chains
chainConfigs := map[por.ChainSelector]*ChainConfig{
    1: { // Ethereum Mainnet
        ChainSelector:   1,
        RelayID:         "evm-1",
        ContractAddress: "0x...",
        GasLimit:        500000,
        MaxGasPrice:     "1000000000",
        FromAddress:     "0x...",
    },
    137: { // Polygon
        ChainSelector:   137,
        RelayID:         "evm-137",
        ContractAddress: "0x...",
        GasLimit:        300000,
        MaxGasPrice:     "50000000000",
        FromAddress:     "0x...",
    },
}
```

### 2. Relayer Configuration

```go
// Map relay IDs to relayers
relayers := map[types.RelayID]loop.Relayer{
    "evm-1":   ethereumRelayer,
    "evm-137": polygonRelayer,
}
```

## Benefits of This Pattern

### 1. **Isolation**
- Each chain operates independently
- Failures on one chain don't affect others
- Separate configuration and monitoring per chain

### 2. **Scalability**
- Easy to add new chains
- Independent scaling per chain
- Chain-specific optimizations

### 3. **Maintainability**
- Clear separation of concerns
- Chain-specific error handling
- Independent testing per chain

### 4. **Reliability**
- Chain-specific health monitoring
- Independent retry logic
- Separate transaction management

## Comparison with CCIP

This pattern mirrors CCIP's multi-chain architecture:

### CCIP Pattern:
- **Relayers**: One per chain
- **Chain Writers**: One per chain with chain-specific TxMgrs
- **Contract Readers**: One per chain
- **Health Monitoring**: Per-chain health reports

### Secure Mint Pattern:
- **Relayers**: One per chain (same as CCIP)
- **Chain Writers**: One per chain (same as CCIP)
- **Contract Readers**: One per chain (same as CCIP)
- **Contract Transmitters**: One per chain (Secure Mint specific)

## Implementation Notes

### 1. **Thread Safety**
- All operations are protected by RWMutex
- Concurrent access to chain-specific services
- Safe service lifecycle management

### 2. **Error Handling**
- Chain-specific error propagation
- Graceful degradation when chains fail
- Comprehensive logging per chain

### 3. **Resource Management**
- Proper cleanup of chain-specific resources
- Independent service lifecycle per chain
- Memory-efficient chain management

### 4. **Configuration Management**
- Chain-specific configuration validation
- Dynamic configuration updates
- Configuration inheritance and overrides

## Future Enhancements

### 1. **Dynamic Chain Addition**
- Runtime chain configuration updates
- Hot-swapping of chain configurations
- Dynamic relayer management

### 2. **Advanced Monitoring**
- Chain-specific metrics
- Performance monitoring per chain
- Alerting based on chain health

### 3. **Load Balancing**
- Chain-specific load distribution
- Intelligent chain selection
- Performance-based routing

This multi-chain pattern provides a robust foundation for the Secure Mint plugin to operate across multiple blockchain networks while maintaining the reliability and scalability characteristics of the CCIP implementation. 