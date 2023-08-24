pragma solidity 0.8.6;

import "../../../automation/interfaces/AutomationCompatibleInterface.sol";
import "../2_1/interfaces/FeedLookupCompatibleInterface.sol";
import "./MercuryRegistry.sol";

contract MercuryRegistryBatchUpkeep is AutomationCompatibleInterface, FeedLookupCompatibleInterface {
  // Use a reasonable maximum batch size. Every Mercury report is ~750 bytes, too many reports
  // passed into a single batch could exceed the calldata or transaction size limit for some blockchains.
  uint256 constant MAX_BATCH_SIZE = 50;

  MercuryRegistry immutable i_registry; // master registry, where feed data is stored

  uint256 s_batchStart; // starting index of upkeep batch, inclusive
  uint256 s_batchEnd; // ending index of upkeep batch, exclusive

  constructor(address mercuryRegistry, uint256 batchStart, uint256 batchEnd) {
    i_registry = MercuryRegistry(mercuryRegistry);

    // Do not allow a batched mercury registry to use an excessive batch size, as to avoid
    // calldata size limits. If more feeds need to be updated than allowed by the batch size,
    // deploy another `MercuryRegistryBatchUpkeep` contract and register another upkeep job.
    require(batchEnd - batchStart <= MAX_BATCH_SIZE, "batch size is too large");

    s_batchStart = batchStart;
    s_batchEnd = batchEnd;
  }

  // Invoke a feed lookup for the feeds this upkeep is responsible for.
  function checkUpkeep(bytes calldata /* data */) external view override returns (bool, bytes memory) {
    uint256 start = s_batchStart;
    uint256 end = s_batchEnd;
    string[] memory feeds = new string[](end - start);
    for (uint256 i = start; i < end; i++) {
      feeds[i - start] = i_registry.s_feeds(i);
    }
    return i_registry.revertForFeedLookup(feeds);
  }

  // Use the master registry to assess deviations.
  function checkCallback(
    bytes[] memory values,
    bytes memory lookupData
  ) external view override returns (bool, bytes memory) {
    return i_registry.checkCallback(values, lookupData);
  }

  // Use the master registry to update state.
  function performUpkeep(bytes calldata performData) external override {
    i_registry.performUpkeep(performData);
  }
}
