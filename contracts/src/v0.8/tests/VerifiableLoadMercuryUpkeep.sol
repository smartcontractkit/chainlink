// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./VerifiableLoadBase.sol";
import "../dev/automation/2_1/interfaces/FeedLookupCompatibleInterface.sol";

contract VerifiableLoadMercuryUpkeep is VerifiableLoadBase, FeedLookupCompatibleInterface {
  string[] public feedsHex = [
    "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
    "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000",
    "0x555344432d5553442d415242495452554d2d544553544e455400000000000000"
  ];
  string public constant feedParamKey = "feedIDHex";
  string public constant timeParamKey = "blockNumber";

  event MercuryPerformEvent(
    address indexed origin,
    uint256 indexed upkeepId,
    uint256 indexed blockNumber,
    bytes v0,
    bytes ed
  );

  constructor(address registrarAddress, bool useArb) VerifiableLoadBase(registrarAddress, useArb) {}

  function setFeedsHex(string[] memory newFeeds) external {
    feedsHex = newFeeds;
  }

  function checkCallback(
    bytes[] memory values,
    bytes memory extraData
  ) external pure override returns (bool, bytes memory) {
    // do sth about the chainlinkBlob data in values and extraData
    bytes memory performData = abi.encode(values, extraData);
    return (true, performData);
  }

  function checkUpkeep(bytes calldata checkData) external returns (bool, bytes memory) {
    uint256 startGas = gasleft();
    uint256 upkeepId = abi.decode(checkData, (uint256));

    uint256 performDataSize = performDataSizes[upkeepId];
    uint256 checkGasToBurn = checkGasToBurns[upkeepId];
    bytes memory pData = abi.encode(upkeepId, new bytes(performDataSize));
    uint256 blockNum = getBlockNumber();
    bool needed = eligible(upkeepId);
    while (startGas - gasleft() + 10000 < checkGasToBurn) {
      // 10K margin over gas to burn
      // Hard coded check gas to burn
      dummyMap[blockhash(blockNum)] = false; // arbitrary storage writes
    }
    if (!needed) {
      return (false, pData);
    }

    revert FeedLookup(feedParamKey, feedsHex, timeParamKey, blockNum, abi.encode(upkeepId));
  }

  function performUpkeep(bytes calldata performData) external {
    uint256 startGas = gasleft();
    (bytes[] memory values, bytes memory extraData) = abi.decode(performData, (bytes[], bytes));
    uint256 upkeepId = abi.decode(extraData, (uint256));
    uint256 firstPerformBlock = firstPerformBlocks[upkeepId];
    uint256 previousPerformBlock = previousPerformBlocks[upkeepId];
    uint256 blockNum = getBlockNumber();

    if (firstPerformBlock == 0) {
      firstPerformBlocks[upkeepId] = blockNum;
      firstPerformBlock = blockNum;
      timestamps[upkeepId].push(block.timestamp);
    } else {
      // Calculate and append delay
      uint256 delay = blockNum - previousPerformBlock - intervals[upkeepId];

      uint16 timestampBucket = timestampBuckets[upkeepId];
      if (block.timestamp - TIMESTAMP_INTERVAL > timestamps[upkeepId][timestampBucket]) {
        timestamps[upkeepId].push(block.timestamp);
        timestampBucket++;
        timestampBuckets[upkeepId] = timestampBucket;
      }

      uint16 bucket = buckets[upkeepId];
      uint256[] memory bucketDelays = bucketedDelays[upkeepId][bucket];
      if (bucketDelays.length == BUCKET_SIZE) {
        bucket++;
        buckets[upkeepId] = bucket;
      }
      bucketedDelays[upkeepId][bucket].push(delay);
      timestampDelays[upkeepId][timestampBucket].push(delay);
      delays[upkeepId].push(delay);
    }

    uint256 counter = counters[upkeepId] + 1;
    counters[upkeepId] = counter;
    emit PerformingUpkeep(firstPerformBlock, blockNum, previousPerformBlock, counter);
    previousPerformBlocks[upkeepId] = blockNum;

    // for every upkeepTopUpCheckInterval (5), check if the upkeep balance is at least
    // minBalanceThresholdMultiplier (20) * min balance. If not, add addLinkAmount (0.2) to the upkeep
    // upkeepTopUpCheckInterval, minBalanceThresholdMultiplier, and addLinkAmount are configurable
    if (blockNum - lastTopUpBlocks[upkeepId] > upkeepTopUpCheckInterval) {
      UpkeepInfo memory info = registry.getUpkeep(upkeepId);
      uint96 minBalance = registry.getMinBalanceForUpkeep(upkeepId);
      if (info.balance < minBalanceThresholdMultiplier * minBalance) {
        this.addFunds(upkeepId, addLinkAmount);
        lastTopUpBlocks[upkeepId] = blockNum;
        emit UpkeepTopUp(upkeepId, addLinkAmount, blockNum);
      }
    }

    uint256 performGasToBurn = performGasToBurns[upkeepId];
    while (startGas - gasleft() + 10000 < performGasToBurn) {
      // 10K margin over gas to burn
      dummyMap[blockhash(blockNum)] = false; // arbitrary storage writes
    }

    //    bytes memory v0 = VERIFIER.verify(values[0]);
    //    bytes memory v1 = VERIFIER.verify(values[1]);
    emit MercuryPerformEvent(tx.origin, upkeepId, blockNum, values[0], extraData);
  }
}
