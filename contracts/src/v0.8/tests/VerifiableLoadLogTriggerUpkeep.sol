// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import "./VerifiableLoadBase.sol";
import "../dev/automation/2_1/interfaces/ILogAutomation.sol";
import "../dev/automation/2_1/interfaces/FeedLookupCompatibleInterface.sol";

contract VerifiableLoadLogTriggerUpkeep is VerifiableLoadBase, ILogAutomation, FeedLookupCompatibleInterface {
  event PerformingLogTriggerUpkeep(
    address indexed from,
    uint256 upkeepId,
    uint256 counter,
    uint256 logBlockNum,
    uint256 blockNum,
    bytes blob,
    bytes verified
  );

  // for mercury config
  string[] public feedsHex = [
    "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
    "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"
  ];
  string public feedParamKey = "feedIdHex";
  string public timeParamKey = "blockNumber";

  constructor(address _registrar, bool _useL1BlockNumber) VerifiableLoadBase(_registrar, _useL1BlockNumber) {}

  function setFeedsHex(string[] memory newFeeds) external {
    feedsHex = newFeeds;
  }

  function checkLog(Log calldata log) external override returns (bool upkeepNeeded, bytes memory performData) {
    uint256 startGas = gasleft();
    uint256 blockNum = getBlockNumber();

    // filter by event signature
    if (log.topics[0] == emittedSig) {
      // filter by indexed parameters
      bytes memory t1 = abi.encodePacked(log.topics[1]); // bytes32 to bytes
      uint256 upkeepId = abi.decode(t1, (uint256));
      bool needed = eligible(upkeepId);
      if (needed) {
        uint256 checkGasToBurn = checkGasToBurns[upkeepId];
        while (startGas - gasleft() + 10000 < checkGasToBurn) {
          // 10K margin over gas to burn
          // Hard coded check gas to burn
          dummyMap[blockhash(blockNum)] = false; // arbitrary storage writes
        }

        revert FeedLookup(feedParamKey, feedsHex, timeParamKey, blockNum, abi.encode(upkeepId, log.blockNumber));
      }
    }
    revert("could not find matching event sig");
  }

  function performUpkeep(bytes calldata performData) external override {
    uint256 startGas = gasleft();
    (bytes[] memory values, bytes memory extraData) = abi.decode(performData, (bytes[], bytes));
    (uint256 upkeepId, uint256 logBlockNumber) = abi.decode(extraData, (uint256, uint256));

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
      KeeperRegistryBase2_1.UpkeepInfo memory info = registry.getUpkeep(upkeepId);
      uint96 minBalance = registry.getMinBalanceForUpkeep(upkeepId);
      if (info.balance < minBalanceThresholdMultiplier * minBalance) {
        this.addFunds(upkeepId, addLinkAmount);
        lastTopUpBlocks[upkeepId] = blockNum;
        emit UpkeepTopUp(upkeepId, addLinkAmount, blockNum);
      }
    }

    bytes memory verifiedResponse = VERIFIER.verify(values[0]);
    emit LogEmitted(upkeepId, logBlockNumber, blockNum);
    emit PerformingLogTriggerUpkeep(
      tx.origin,
      upkeepId,
      counter,
      logBlockNumber,
      blockNum,
      values[0],
      verifiedResponse
    );

    uint256 performGasToBurn = performGasToBurns[upkeepId];
    while (startGas - gasleft() + 10000 < performGasToBurn) {
      // 10K margin over gas to burn
      dummyMap[blockhash(blockNum)] = false; // arbitrary storage writes
    }
  }

  function checkCallback(
    bytes[] memory values,
    bytes memory extraData
  ) external view override returns (bool upkeepNeeded, bytes memory performData) {
    // do sth about the chainlinkBlob data in values and extraData
    bytes memory performData = abi.encode(values, extraData);
    return (true, performData);
  }
}
