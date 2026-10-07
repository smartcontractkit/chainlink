# OCR2 Delegate Registry-Driven Launch

## Goal

Decouple OCR2 delegate service creation from the job-spec lifecycle so DONTime (and other OCR3-based capabilities) can be launched from Capabilities Registry changes, without a job spec.

## Background

The standard capabilities delegate was partially refactored to support registry-driven launch (PR #22982). The `LocalCapabilityManager` calls `NewServicesFn` when the registry changes, threading in the capability ID, DON ID, and OCR3 config. The OCR2 delegate needs the same treatment.

## Current State

`ocr2/delegate.go` `ServicesForSpec` is tightly coupled to `jb.OCR2OracleSpec`:
- Transmitter ID, key bundle ID, bootstrap peers, signing strategy → all from job spec
- Plugin type switch (`newDonTimePlugin`, `newServicesLLO`, etc.) → all read from spec
- Relay config, contract ID, external job ID → all from spec

## Proposed Changes

### 1. Extract `NewServices` method

Add a method to `ocr2.Delegate` that works without a job spec:

```go
func (d *Delegate) NewServices(
    ctx context.Context,
    capabilityID string,
    donID uint32,
    pluginType types.PluginType,
    configJSON string,
    registryOCRConfig *ocrtypes.ContractConfig,
) ([]job.ServiceCtx, error)
```

This method:
- Resolves OCR key bundle from keystore (or CapReg OCR config)
- Resolves transmitter from CapReg OCR config (or keystore fallback)
- Resolves bootstrap peers from `P2P.V2.DefaultBootstrappers`
- Resolves signing strategy from keystore
- Calls the appropriate `newDonTimePlugin` / `newServicesGenericPlugin` etc.

### 2. Wire `LocalCapabilityManager` for DONTime

In `cre.go`, add `dontime@1.0.0` to the `RegistryBasedLaunchAllowlist` and wire the OCR2 delegate:

```go
s.SetDelegatesDeps = func(ocr2Delegate *ocr.Delegate) (commonsrv.Service, error) {
    newServicesFn := func(ctx context.Context, capID string, donID uint32, command string, configJSON string, ocr3Config *ocrtypes.ContractConfig) ([]job.ServiceCtx, error) {
        return ocr2Delegate.NewServices(ctx, capID, donID, types.DonTimePlugin, configJSON, ocr3Config)
    }
    localCapMgr, lcmErr := localcapmgr.NewLocalCapabilityManager(lggr, localCfg, newServicesFn)
    ...
}
```

### 3. Resolve config from CapReg

Reuse `ResolveOracleFactoryConfig` from `oracle_factory_config.go` to fill missing fields:
- Contract address → Capabilities Registry address
- Chain ID → Capabilities Registry chain ID
- Key bundle → node keystore
- Transmitter → CapReg OCR config (paired with signer) or keystore fallback
- Signing strategy → `multi-chain` with EVM key bundle

### 4. Update `LocalCapabilityManager`

The `NewServicesFn` signature may need to be extended to pass `pluginType` or the OCR2 delegate needs to infer it from the capability ID.

## Files to Change

- `core/services/ocr2/delegate.go` — extract `NewServices`, refactor `newDonTimePlugin` to accept resolved config
- `core/services/cre/cre.go` — wire `LocalCapabilityManager` for DONTime
- `core/services/chainlink/application.go` — pass `defaultBootstrappers`, `capRegistryAddress`, `capRegistryChainID` to OCR2 delegate
- `core/capabilities/localcapmgr/manager.go` — possibly extend `NewServicesFn` signature

## Testing

- Unit tests for `NewServices` with nil job spec (registry-driven path)
- Unit tests for config resolution from CapReg
- Integration tests via CRE E2E (DONTime without job spec)
