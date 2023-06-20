// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {ILogAutomation, Log} from "../2_1/interfaces/ILogAutomation.sol";
import "../2_1/interfaces/FeedLookupCompatibleInterface.sol";
import {ArbSys} from "../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol";
import "../../../automation/2_0/KeeperRegistrar2_0.sol";
import "../../../tests/VerifiableLoadBase.sol";

contract VerifiableLoadLogTriggeredFeedLookup is VerifiableLoadBase, ILogAutomation, FeedLookupCompatibleInterface {
  event TriggerMercury(uint256 indexed upkeepId, uint256 indexed logBlockNumber); // keccak256(TriggerMercury(uint256,uint256)) => 0xcd89a1cdede3e128a8e92d77495b16cc12f0fc7564a712113f006adaf640a4a6
  event DoNotTriggerMercury(uint256 indexed bn, uint256 indexed ts); // keccak256(DoNotTriggerMercury(uint256,uint256)) => 0x3338104ccab396f091b767e17a2a863f70d777261982306eeb72053c94b9cd47
  event PerformingLogTriggerUpkeep(
    address indexed from,
    uint256 upkeepId,
    uint256 counter,
    uint256 logBlockNumber,
    uint256 blockNumber
  );

  // for log trigger
  bytes32 constant triggerSig = 0xcd89a1cdede3e128a8e92d77495b16cc12f0fc7564a712113f006adaf640a4a6;
  bytes32 constant NoTriggerSig = 0x3338104ccab396f091b767e17a2a863f70d777261982306eeb72053c94b9cd47;

  // for mercury config
  string[] public feedsHex = [
    "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
    "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000",
    "0x555344432d5553442d415242495452554d2d544553544e455400000000000000"
  ];
  string public constant feedParamKey = "feedIDHex";
  string public constant timeParamKey = "blockNumber";

  constructor(address _registrar, bool _useL1BlockNumber) VerifiableLoadBase(_registrar, _useL1BlockNumber) {}

  // ???
  function registerUpkeep() external returns (bytes memory logTrigger) {
    KeeperRegistryBase2_1.LogTriggerConfig memory cfg = KeeperRegistryBase2_1.LogTriggerConfig({
      contractAddress: address(this),
      filterSelector: 0, // ??
      topic0: triggerSig, // only triggerSig will be available at `checkLog`?
      topic1: 0x000000000000000000000000000000000000000000000000000000000000000, // ??
      topic2: 0x000000000000000000000000000000000000000000000000000000000000000,
      topic3: 0x000000000000000000000000000000000000000000000000000000000000000
    });
    bytes memory logTriggerConfig = abi.encode(cfg);

    return logTriggerConfig;
  }

  function startLogs(uint256 upkeepId) external {
    emit TriggerMercury(upkeepId, getBlockNumber());
    emit DoNotTriggerMercury(getBlockNumber(), block.timestamp);
  }

  function setFeedsHex(string[] memory newFeeds) external {
    feedsHex = newFeeds;
  }

  function checkLog(Log calldata log) external override returns (bool upkeepNeeded, bytes memory performData) {
    uint256 startGas = gasleft();
    uint256 blockNum = getBlockNumber();

    // filter by event signature
    if (log.topics[0] == triggerSig) {
      // filter by indexed parameters
      bytes memory p1 = abi.encodePacked(log.topics[1]); // bytes32 to bytes
      uint256 upkeepId = abi.decode(p1, (uint256));
      bytes memory p2 = abi.encodePacked(log.topics[2]);
      uint256 logBlockNumber = abi.decode(p2, (uint256));

      uint256 checkGasToBurn = checkGasToBurns[upkeepId];
      bool needed = eligible(upkeepId);

      if (needed) {
        while (startGas - gasleft() + 10000 < checkGasToBurn) {
          // 10K margin over gas to burn
          // Hard coded check gas to burn
          dummyMap[blockhash(blockNum)] = false; // arbitrary storage writes
          blockNum--;
        }
        revert FeedLookup(feedParamKey, feedsHex, timeParamKey, blockNum, abi.encodePacked(upkeepId, logBlockNumber));
      }
      revert("upkeep not needed");
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

    uint256 performGasToBurn = performGasToBurns[upkeepId];
    while (startGas - gasleft() + 10000 < performGasToBurn) {
      // 10K margin over gas to burn
      dummyMap[blockhash(blockNum)] = false; // arbitrary storage writes
    }

    emit TriggerMercury(upkeepId, blockNum);
    emit DoNotTriggerMercury(blockNum, block.timestamp);
    emit PerformingLogTriggerUpkeep(tx.origin, upkeepId, counter, logBlockNumber, blockNum);
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
