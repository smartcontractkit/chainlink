# Migration Plan: Secure Mint Workflow t## Implementation Steps
1. ✅ **COMPLETED**: Add DataFeedsCache imports alongside existing imports
2. ✅ **COMPLETED**: Create new `SetupDataFeedsCacheContract` function alongside existing setup functions
3. ✅ **COMPLETED**: Create new `setupSecureMintDons` function for secure mint test only
4. ✅ **COMPLETED**: Update workflow ABI encoding for new field order
5. ✅ **COMPLETED**: Update workflow target method from `receive(report bytes)` to `onReport(bytes metadata, bytes rawReport)`
6. ✅ **COMPLETED**: Create new event handler for DataFeedsCache events (secure mint test only)
7. ✅ **COMPLETED**: Update `Test_runSecureMintWorkflow` to use new setup and event handling
8. ✅ **COMPLETED**: Ensure other tests remain unchanged and functionaledsCache

## Overview
Update `Test_runSecureMintWorkflow` to use DataFeedsCache contract instead of KeystoneFeedsConsumer, while keeping other integration tests in the package unchanged (they will continue using KeystoneFeedsConsumer).

## Key Changes Required

### 1. Contract Import and Type Updates
- Add `data_feeds_cache` import alongside existing `feeds_consumer` import
- Update contract wrapper types for secure mint test only
- Keep existing KeystoneFeedsConsumer imports for other tests

### 2. Contract Deployment and Configuration
- Create new `SetupDataFeedsCacheContract` function alongside existing `SetupConsumerContract`
- Use `setDecimalFeedConfigs` for DataFeedsCache setup (per user choice #1)
- Keep existing `setupKeystoneDons` function unchanged for other tests
- Create separate setup function for secure mint test: `setupSecureMintDons`

### 3. Method Signature Changes
- **Current**: `receive(report bytes)` with `(bytes32 FeedID, uint224 Price, uint32 Timestamp)[]`
- **New**: `onReport(bytes metadata, bytes rawReport)` with `(bytes32 RemappedID, uint32 Timestamp, uint224 Price)[]`
- **Key Difference**: Field order changed - timestamp and price swapped positions
- **Workflow Update**: Need to update `secureMintWorkflow` ABI from `"(bytes32 FeedID, uint224 Price, uint32 Timestamp)[] Reports"` to `"(bytes32 RemappedID, uint32 Timestamp, uint224 Price)[] Reports"`

### 4. Event Handling Updates
- Create new event handler for DataFeedsCache `DecimalReportUpdated` events (secure mint test only)
- Keep existing `FeedReceived` event handling for other tests
- Update secure mint handler to process different event structure

### 5. Test Isolation
- Modify only `Test_runSecureMintWorkflow` to use DataFeedsCache
- Keep all other integration tests unchanged (using KeystoneFeedsConsumer)
- Ensure no breaking changes to existing test infrastructure

### 5. Source of Field Order Information
The field order `(bytes32 RemappedID, uint32 Timestamp, uint224 Price)[]` comes from:
- **DataFeedsCache.sol line ~490**: `ReceivedDecimalReport` struct definition
- **DataFeedsCache.sol line ~570**: ABI decode logic in `onReport` method
- Confirmed by contract source: timestamp field comes before answer/price field

## Implementation Steps
1. Add DataFeedsCache imports alongside existing imports
2. Create new `SetupDataFeedsCacheContract` function alongside existing setup functions
3. Create new `setupSecureMintDons` function for secure mint test only
4. Update workflow ABI encoding for new field order
5. Create new event handler for DataFeedsCache events (secure mint test only)
6. Update `Test_runSecureMintWorkflow` to use new setup and event handling
7. Ensure other tests remain unchanged and functional

## Implementation Summary

### Changes Made:

1. **Contract Imports**: Added DataFeedsCache import alongside existing KeystoneFeedsConsumer imports
2. **New Setup Functions**:
   - `setupSecureMintDons()`: New setup function specifically for secure mint test
   - `SetupDataFeedsCacheContract()`: Deploys and configures DataFeedsCache contract with `setDecimalFeedConfigs()`
3. **Workflow Updates**:
   - Updated encoder ABI: `"(bytes32 RemappedID, uint32 Timestamp, uint224 Price)[] Reports"`
   - Updated target method: `"onReport(bytes metadata, bytes rawReport)"`
4. **Event Handling**:
   - New `dataFeedsCacheHandler` interface for DataFeedsCache events
   - New `waitForDataFeedsCacheReports()` function for `DecimalReportUpdated` events
   - Updated `secureMintHandler` to implement both interfaces
5. **Test Migration**:
   - `Test_runSecureMintWorkflow` now uses DataFeedsCache instead of KeystoneFeedsConsumer
   - All other tests continue using KeystoneFeedsConsumer (no breaking changes)

### Key Differences Handled:
- **Contract Method**: `receive(report bytes)` → `onReport(bytes metadata, bytes rawReport)`
- **Field Order**: `(FeedID, Price, Timestamp)` → `(RemappedID, Timestamp, Price)`
- **Event Type**: `FeedReceived` → `DecimalReportUpdated`
- **Data ID Format**: `bytes32 feedId` → `bytes16 dataId`
- **Configuration**: `SetConfig()` → `SetFeedAdmin()` + `SetDecimalFeedConfigs()`

**Status**: ✅ Implementation Complete - All files compile successfully
