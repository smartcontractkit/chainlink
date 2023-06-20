// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {ILogAutomation, Log} from "../2_1/interfaces/ILogAutomation.sol";
import "../2_1/interfaces/FeedLookupCompatibleInterface.sol";
import {ArbSys} from "../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol";
import "../../../automation/2_0/KeeperRegistrar2_0.sol";
import "../../../tests/VerifiableLoadBase.sol";

contract LogTriggeredFeedLookup is ILogAutomation, FeedLookupCompatibleInterface {
  event TriggerMercury(uint256 indexed upkeepId, uint256 indexed logBlockNumber); // keccak256(TriggerMercury(uint256,uint256)) => 0xcd89a1cdede3e128a8e92d77495b16cc12f0fc7564a712113f006adaf640a4a6
  event DoNotTriggerMercury(uint256 indexed bn, uint256 indexed ts); // keccak256(DoNotTriggerMercury(uint256,uint256)) => 0x3338104ccab396f091b767e17a2a863f70d777261982306eeb72053c94b9cd47
  event PerformingLogTriggerUpkeep(
    address indexed from,
    uint256 upkeepId,
    uint256 counter,
    uint256 logBlockNumber,
    uint256 blockNumber
  );

  ArbSys internal constant ARB_SYS = ArbSys(0x0000000000000000000000000000000000000064);

  // for log trigger
  bytes32 constant triggerSig = 0xcd89a1cdede3e128a8e92d77495b16cc12f0fc7564a712113f006adaf640a4a6;
  bytes32 constant NoTriggerSig = 0x3338104ccab396f091b767e17a2a863f70d777261982306eeb72053c94b9cd47;

  // for mercury config
  string[] public feedsHex = ["0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"];
  string public constant feedParamKey = "feedIDHex";
  string public constant timeParamKey = "blockNumber";

  uint256 public testRange;
  uint256 public interval;
  uint256 public previousPerformBlock;
  uint256 public initialBlock;
  uint256 public counter;
  bool public useArbitrumBlockNum;

  constructor(uint256 _testRange, uint256 _interval, bool _useArbitrumBlockNum) {
    testRange = _testRange;
    interval = _interval;
    useArbitrumBlockNum = _useArbitrumBlockNum;
    previousPerformBlock = 0;
    initialBlock = 0;
    counter = 0;
  }

  function registerUpkeep() external returns (bytes memory logTrigger) {
    KeeperRegistryBase2_1.LogTriggerConfig memory cfg = KeeperRegistryBase2_1.LogTriggerConfig({
      contractAddress: address(this),
      filterSelector: 0, // not applying any filter, if 1, it will check topic1, if 3, it will check both topic1 and topic2, if 5, it will check both topic1 and topic3
      topic0: triggerSig, // only triggerSig will be available at `checkLog` this will always apply
      topic1: 0x000000000000000000000000000000000000000000000000000000000000000,
      topic2: 0x000000000000000000000000000000000000000000000000000000000000000,
      topic3: 0x000000000000000000000000000000000000000000000000000000000000000
    });
    bytes memory logTriggerConfig = abi.encode(cfg);

    return logTriggerConfig;
  }

  function startLogs(uint256 upkeepId) external {
    uint256 blockNumber = getBlockNumber();
    emit TriggerMercury(upkeepId, blockNumber);
    emit DoNotTriggerMercury(blockNumber, block.timestamp);
  }

  function setFeedsHex(string[] memory newFeeds) external {
    feedsHex = newFeeds;
  }

  function checkLog(Log calldata log) external override returns (bool upkeepNeeded, bytes memory performData) {
    uint256 blockNum = getBlockNumber();

    // filter by event signature
    if (log.topics[0] == triggerSig) {
      // filter by indexed parameters
      bytes memory p1 = abi.encodePacked(log.topics[1]); // bytes32 to bytes
      uint256 upkeepId = abi.decode(p1, (uint256));
      bytes memory p2 = abi.encodePacked(log.topics[2]);
      uint256 logBlockNumber = abi.decode(p2, (uint256));

      bool needed = eligible();
      if (needed) {
        revert FeedLookup(feedParamKey, feedsHex, timeParamKey, blockNum, abi.encodePacked(upkeepId, logBlockNumber));
      }
      revert("upkeep not needed");
    }
    revert("could not find matching event sig");
  }

  function performUpkeep(bytes calldata performData) external override {
    uint256 blockNumber = getBlockNumber();
    if (initialBlock == 0) {
      initialBlock = blockNumber;
    }
    counter = counter + 1;
    previousPerformBlock = blockNumber;

    (bytes[] memory values, bytes memory extraData) = abi.decode(performData, (bytes[], bytes));
    (uint256 upkeepId, uint256 logBlockNumber) = abi.decode(extraData, (uint256, uint256));

    emit TriggerMercury(upkeepId, blockNumber);
    emit DoNotTriggerMercury(blockNumber, block.timestamp);
    emit PerformingLogTriggerUpkeep(tx.origin, upkeepId, counter, logBlockNumber, blockNumber);
  }

  function checkCallback(
    bytes[] memory values,
    bytes memory extraData
  ) external view override returns (bool upkeepNeeded, bytes memory performData) {
    // do sth about the chainlinkBlob data in values and extraData
    bytes memory performData = abi.encode(values, extraData);
    return (true, performData);
  }

  function getBlockNumber() internal view returns (uint256) {
    if (useArbitrumBlockNum) {
      return ARB_SYS.arbBlockNumber();
    } else {
      return block.number;
    }
  }

  function eligible() public view returns (bool) {
    if (initialBlock == 0) {
      return true;
    }

    uint256 blockNumber = getBlockNumber();
    return (blockNumber - initialBlock) < testRange && (blockNumber - previousPerformBlock) >= interval;
  }

  function setSpread(uint256 _testRange, uint256 _interval) external {
    testRange = _testRange;
    interval = _interval;
    initialBlock = 0;
    counter = 0;
  }
}
