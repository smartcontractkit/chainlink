// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {ILogAutomation, Log} from "../interfaces/ILogAutomation.sol";
import "../interfaces/StreamsLookupCompatibleInterface.sol";

contract LogTriggeredLoad is ILogAutomation, StreamsLookupCompatibleInterface {
  event PerformingFeedLookup(
    address indexed from,
    bytes blob
  );
  event PerformingUpkeep(
    address indexed from,
    uint256 initialBlock,
    uint256 lastBlock,
    uint256 previousBlock,
    uint256 counter,
    uint256 timeToPerform,
    bool isRecovered
  );

  mapping(bytes32 => bool) public dummyMap; // used to force storage lookup
  uint256 public lastBlock;
  uint256 public previousPerformBlock;
  uint256 public initialBlock;
  uint256 public counter;
  uint256 public timeToPerform;
  bool public isRecovered;

  // for feed lookup config
  bool public isFeedsLookup;
  bool public shouldRetryOnError;
  string[] public feedsHex = ["0x000200"];
  string public feedParamKey = "feedIDs";
  string public timeParamKey = "timestamp";
  uint256 public counter;

  constructor(bool _isFeedLookup) {
    isFeedsLookup = _isFeedLookup;
    previousPerformBlock = 0;
    lastBlock = block.number;
    initialBlock = 0;
    counter = 0;
  }

  function setTimeParamKey(string memory timeParam) external {
    timeParamKey = timeParam;
  }

  function setFeedParamKey(string memory feedParam) external {
    feedParamKey = feedParam;
  }

  function setFeedsHex(string[] memory newFeeds) external {
    feedsHex = newFeeds;
  }

  function setShouldRetryOnErrorBool(bool value) public {
    shouldRetryOnError = value;
  }

  function checkLog(
    Log calldata log,
    bytes calldata checkData
  ) external override returns (bool upkeepNeeded, bytes memory performData) {

    (uint256 checkBurnAmount, uint256 performBurnAmount, bytes32 eventSig) = abi.decode(
      checkData,
      (uint256, uint256, bytes32)
    );
    uint256 startGas = gasleft();
    bytes32 dummyIndex = blockhash(block.number - 1);
    bool dummy;
    // burn gas
    if (checkBurnAmount > 0) {
      while (startGas - gasleft() < checkBurnAmount) {
        dummy = dummy && dummyMap[dummyIndex]; // arbitrary storage reads
        dummyIndex = keccak256(abi.encode(dummyIndex, address(this)));
      }
    }
    if (log.topics[2] == eventSig) {
      revert StreamsLookup(
        feedParamKey,
        feedsHex,
        timeParamKey,
        block.number,
        abi.encode(log, block.number, checkData)
      );
    }
    revert("could not find matching event sig");
  }

  function performUpkeep(bytes calldata performData) external override {
    if (initialBlock == 0) {
      initialBlock = block.number;
    }
    lastBlock = block.number;
    counter = counter + 1;
    previousPerformBlock = lastBlock;
    (bytes[] memory values, bytes memory extraData) = abi.decode(performData, (bytes[], bytes));

    (Log memory log, uint256 checkBlock, bytes memory extraData) = abi.decode(extraData, (Log, uint256, bytes));
    timeToPerform = block.timestamp - log.timestamp;
    isRecovered = false;
    if (checkBlock != log.blockNumber) {
      isRecovered = true;
    }
    (uint256 checkBurnAmount, uint256 performBurnAmount, bytes32 eventSig) = abi.decode(
      extraData,
      (uint256, uint256, bytes32)
    );
    uint256 startGas = gasleft();
    bytes32 dummyIndex = blockhash(block.number - 1);
    bool dummy;
    // burn gas
    if (performBurnAmount > 0) {
      while (startGas - gasleft() < performBurnAmount) {
        dummy = dummy && dummyMap[dummyIndex]; // arbitrary storage reads
        dummyIndex = keccak256(abi.encode(dummyIndex, address(this)));
      }
    }
    emit PerformingUpkeep(
      tx.origin,
      initialBlock,
      lastBlock,
      previousPerformBlock,
      counter,
      timeToPerform,
      isRecovered
    );
    emit PerformingFeedLookup(
      values[0]
    );
  }

  function checkCallback(
    bytes[] memory values,
    bytes memory extraData
  ) external view override returns (bool, bytes memory) {
    bytes memory performData = abi.encode(values, extraData);
    return (true, performData);
  }

  function checkErrorHandler(
    uint256 errCode,
    bytes memory extraData
  ) external view returns (bool upkeepNeeded, bytes memory performData) {
    bytes[] memory values = new bytes[](2);
    values[0] = abi.encode(errCode);
    uint256 tmp = 123;
    bytes32 tmpBytes = 0x0;
    values[1] = abi.encode(tmp, tmp, tmpBytes);
    bytes memory performData = abi.encode(values, extraData);
    return (shouldRetryOnError, performData);
  }
}
