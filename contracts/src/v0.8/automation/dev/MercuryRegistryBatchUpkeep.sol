pragma solidity 0.8.6;

import "../../shared/access/ConfirmedOwner.sol";
import "../interfaces/AutomationCompatibleInterface.sol";
import "../interfaces/StreamsLookupCompatibleInterface.sol";
import "./MercuryRegistry.sol";

contract MercuryRegistryBatchUpkeep is ConfirmedOwner, AutomationCompatibleInterface, StreamsLookupCompatibleInterface {
  error BatchSizeTooLarge(uint256 batchsize, uint256 maxBatchSize);
  // Use a reasonable maximum batch size. Every Mercury report is ~750 bytes, too many reports
  // passed into a single batch could exceed the calldata or transaction size limit for some blockchains.
  uint256 constant MAX_BATCH_SIZE = 50;

  MercuryRegistry immutable i_registry; // master registry, where feed data is stored

  uint256 s_batchStart; // starting index of upkeep batch on the MercuryRegistry's s_feeds array, inclusive
  uint256 s_batchEnd; // ending index of upkeep batch on the MercuryRegistry's s_feeds array, exclusive

  constructor(address mercuryRegistry, uint256 batchStart, uint256 batchEnd) ConfirmedOwner(msg.sender) {
    i_registry = MercuryRegistry(mercuryRegistry);

    updateBatchingWindow(batchStart, batchEnd);
  }

  // Invoke a feed lookup for the feeds this upkeep is responsible for.
  function checkUpkeep(bytes calldata /* data */) external view override returns (bool, bytes memory) {
    uint256 start = s_batchStart;
    uint256 end = s_batchEnd;
    string[] memory feeds = new string[](end - start);
    uint256 count = 0;
    for (uint256 i = start; i < end; i++) {
      string memory feedId;

      // If the feed doesn't exist, then the batching window exceeds the underlying registry length.
      // So, the batch will be partially empty.
      try i_registry.s_feeds(i) returns (string memory f) {
        feedId = f;
      } catch (bytes memory /* data */) {
        break;
      }

      // Assign feed.
      feeds[i - start] = feedId;
      count++;
    }

    // Adjusts the length of the batch to `count` such that it does not
    // contain any empty feed Ids.
    assembly {
      mstore(feeds, count)
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

  function updateBatchingWindow(uint256 batchStart, uint256 batchEnd) public onlyOwner {
    // Do not allow a batched mercury registry to use an excessive batch size, as to avoid
    // calldata size limits. If more feeds need to be updated than allowed by the batch size,
    // deploy another `MercuryRegistryBatchUpkeep` contract and register another upkeep job.
    if (batchEnd - batchStart > MAX_BATCH_SIZE) {
      revert BatchSizeTooLarge(batchEnd - batchStart, MAX_BATCH_SIZE);
    }

    s_batchStart = batchStart;
    s_batchEnd = batchEnd;
  }
}
