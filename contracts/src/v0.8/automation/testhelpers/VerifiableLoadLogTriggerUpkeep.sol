// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import "./VerifiableLoadBase.sol";
import "../interfaces/ILogAutomation.sol";
import "../interfaces/StreamsLookupCompatibleInterface.sol";

contract VerifiableLoadLogTriggerUpkeep is VerifiableLoadBase, StreamsLookupCompatibleInterface, ILogAutomation {
  bool public useMercury;
  uint8 public logNum;

  /**
   * @param _registrar a automation registrar 2.1 address
   * @param _useArb if this contract will use arbitrum block number
   * @param _useMercury if the log trigger upkeeps will use mercury lookup
   */
  constructor(
    AutomationRegistrar2_1 _registrar,
    bool _useArb,
    bool _useMercury
  ) VerifiableLoadBase(_registrar, _useArb) {
    useMercury = _useMercury;
    logNum = 0;
  }

  function setLog(uint8 _log) external {
    logNum = _log;
  }

  function checkLog(Log calldata log, bytes memory checkData) external returns (bool, bytes memory) {
    uint256 startGas = gasleft();
    uint256 blockNum = getBlockNumber();
    uint256 uid = abi.decode(checkData, (uint256));

    bytes32 sig = emittedSig;
    if (logNum != 0) {
      sig = emittedAgainSig;
    }
    // filter by event signature
    if (log.topics[0] == sig) {
      bytes memory t1 = abi.encodePacked(log.topics[1]); // bytes32 to bytes
      uint256 upkeepId = abi.decode(t1, (uint256));
      if (upkeepId != uid) {
        revert("upkeep ids don't match");
      }
      bytes memory t2 = abi.encodePacked(log.topics[2]);
      uint256 blockNum = abi.decode(t2, (uint256));

      bytes memory t3 = abi.encodePacked(log.topics[3]);
      address addr = abi.decode(t3, (address));

      uint256 checkGasToBurn = checkGasToBurns[upkeepId];
      while (startGas - gasleft() + 15000 < checkGasToBurn) {
        dummyMap[blockhash(blockNum)] = false;
      }

      uint256 timeParam;
      if (keccak256(abi.encodePacked(feedParamKey)) == keccak256(abi.encodePacked("feedIdHex"))) {
        timeParam = blockNum;
      } else {
        // assume this will be feedIDs for v0.3
        timeParam = block.timestamp;
      }

      if (useMercury) {
        revert StreamsLookup(feedParamKey, feedsHex, timeParamKey, timeParam, abi.encode(upkeepId, blockNum, addr));
      }

      // if we don't use mercury, create a perform data which resembles the output of checkCallback
      bytes[] memory values = new bytes[](feedsHex.length);
      bytes memory extraData = abi.encode(upkeepId, blockNum, addr);
      return (true, abi.encode(values, extraData));
    }
    revert("unexpected event sig");
  }

  function performUpkeep(bytes calldata performData) external {
    uint256 startGas = gasleft();
    (bytes[] memory values, bytes memory extraData) = abi.decode(performData, (bytes[], bytes));
    (uint256 upkeepId, uint256 logBlockNumber, address addr) = abi.decode(extraData, (uint256, uint256, address));

    uint256 firstPerformBlock = firstPerformBlocks[upkeepId];
    uint256 previousPerformBlock = previousPerformBlocks[upkeepId];
    uint256 currentBlockNum = getBlockNumber();

    if (firstPerformBlock == 0) {
      firstPerformBlocks[upkeepId] = currentBlockNum;
    } else {
      uint256 delay = currentBlockNum - logBlockNumber;
      uint16 bucket = buckets[upkeepId];
      uint256[] memory bucketDelays = bucketedDelays[upkeepId][bucket];
      if (bucketDelays.length == BUCKET_SIZE) {
        bucket++;
        buckets[upkeepId] = bucket;
      }
      bucketedDelays[upkeepId][bucket].push(delay);
      delays[upkeepId].push(delay);
    }

    uint256 counter = counters[upkeepId] + 1;
    counters[upkeepId] = counter;
    previousPerformBlocks[upkeepId] = currentBlockNum;

    // for every upkeepTopUpCheckInterval (5), check if the upkeep balance is at least
    // minBalanceThresholdMultiplier (20) * min balance. If not, add addLinkAmount (0.2) to the upkeep
    // upkeepTopUpCheckInterval, minBalanceThresholdMultiplier, and addLinkAmount are configurable
    topUpFund(upkeepId, currentBlockNum);
    emit LogEmitted(upkeepId, currentBlockNum, address(this));
    burnPerformGas(upkeepId, startGas, currentBlockNum);
  }

  function checkCallback(
    bytes[] memory values,
    bytes memory extraData
  ) external pure override returns (bool, bytes memory) {
    bytes memory performData = abi.encode(values, extraData);
    return (true, performData);
  }

  function checkErrorHandler(
    uint256 errCode,
    bytes memory extraData
  ) external view override returns (bool upkeepNeeded, bytes memory performData) {
    // dummy function with default values
    return (false, new bytes(0));
  }
}
