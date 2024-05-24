// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

import {ILogAutomation, Log} from "../interfaces/ILogAutomation.sol";
import "../interfaces/StreamsLookupCompatibleInterface.sol";

struct CheckData {
  uint256 checkBurnAmount;
  uint256 performBurnAmount;
  bytes32 eventSig;
  string[] feeds;
}

contract SimpleLogUpkeepCounter is ILogAutomation, StreamsLookupCompatibleInterface {
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
  bool internal isRecovered;
  bool public isStreamsLookup;
  bool public shouldRetryOnError;
  string public feedParamKey = "feedIDs";
  string public timeParamKey = "timestamp";

  constructor(bool _isStreamsLookup) {
    previousPerformBlock = 0;
    lastBlock = block.number;
    initialBlock = 0;
    counter = 0;
    isStreamsLookup = _isStreamsLookup;
  }

  function _checkDataConfig(CheckData memory) external {}

  function setTimeParamKey(string memory timeParam) external {
    timeParamKey = timeParam;
  }

  function setFeedParamKey(string memory feedParam) external {
    feedParamKey = feedParam;
  }

  function setShouldRetryOnErrorBool(bool value) public {
    shouldRetryOnError = value;
  }

  function checkLog(Log calldata log, bytes calldata checkData) external view override returns (bool, bytes memory) {
    CheckData memory _checkData = abi.decode(checkData, (CheckData));
    uint256 startGas = gasleft();
    bytes32 dummyIndex = blockhash(block.number - 1);
    bool dummy;
    // burn gas
    if (_checkData.checkBurnAmount > 0) {
      while (startGas - gasleft() < _checkData.checkBurnAmount) {
        dummy = dummy && dummyMap[dummyIndex]; // arbitrary storage reads
        dummyIndex = keccak256(abi.encode(dummyIndex, address(this)));
      }
    }
    bytes[] memory values = new bytes[](2);
    values[0] = abi.encode(0x00);
    values[1] = abi.encode(0x00);
    bytes memory extraData = abi.encode(log, block.number, checkData);
    if (log.topics[2] == _checkData.eventSig) {
      if (isStreamsLookup) {
        revert StreamsLookup(feedParamKey, _checkData.feeds, timeParamKey, block.timestamp, extraData);
      }
      return (true, abi.encode(values, extraData));
    }
    return (false, abi.encode(values, extraData));
  }

  function checkCallback(
    bytes[] memory values,
    bytes memory extraData
  ) external view override returns (bool, bytes memory) {
    // do sth about the chainlinkBlob data in values and extraData
    bytes memory performData = abi.encode(values, extraData);
    return (true, performData);
  }

  function checkErrorHandler(
    uint256 errCode,
    bytes memory extraData
  ) external view override returns (bool upkeepNeeded, bytes memory performData) {
    bytes[] memory values = new bytes[](2);
    values[0] = abi.encode(errCode);
    values[1] = abi.encode(extraData);
    bytes memory returnData = abi.encode(values, extraData);
    return (shouldRetryOnError, returnData);
  }

  function performUpkeep(bytes calldata performData) external override {
    if (initialBlock == 0) {
      initialBlock = block.number;
    }
    lastBlock = block.number;
    counter = counter + 1;
    previousPerformBlock = lastBlock;
    (, bytes memory extraData) = abi.decode(performData, (bytes[], bytes));
    (Log memory log, uint256 checkBlock, bytes memory checkData) = abi.decode(extraData, (Log, uint256, bytes));
    timeToPerform = block.timestamp - log.timestamp;
    isRecovered = false;
    if (checkBlock != log.blockNumber) {
      isRecovered = true;
    }
    CheckData memory _checkData = abi.decode(checkData, (CheckData));
    uint256 startGas = gasleft();
    bytes32 dummyIndex = blockhash(block.number - 1);
    bool dummy;
    if (log.topics[2] != _checkData.eventSig) {
      revert("Invalid event signature");
    }
    // burn gas
    if (_checkData.performBurnAmount > 0) {
      while (startGas - gasleft() < _checkData.performBurnAmount) {
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
  }
}
