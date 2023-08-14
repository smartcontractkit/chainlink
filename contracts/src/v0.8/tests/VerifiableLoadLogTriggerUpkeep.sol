// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import "./VerifiableLoadBase.sol";
import "../dev/automation/2_1/interfaces/ILogAutomation.sol";
import "../dev/automation/2_1/interfaces/FeedLookupCompatibleInterface.sol";

contract VerifiableLoadLogTriggerUpkeep is VerifiableLoadBase, FeedLookupCompatibleInterface, ILogAutomation {
  string[] public feedsHex = [
    "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
    "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"
  ];
  string public feedParamKey = "feedIdHex";
  string public timeParamKey = "blockNumber";
  bool public autoLog;
  bool public useMercury;

  /**
   * @param _registrar a automation registrar 2.1 address
   * @param _useArb if this contract will use arbitrum block number
   * @param _autoLog if the upkeep will emit logs to trigger its next log trigger process
   * @param _useMercury if the log trigger upkeeps will use mercury lookup
   */
  constructor(
    AutomationRegistrar2_1 _registrar,
    bool _useArb,
    bool _autoLog,
    bool _useMercury
  ) VerifiableLoadBase(_registrar, _useArb) {
    autoLog = _autoLog;
    useMercury = _useMercury;
  }

  function setAutoLog(bool _autoLog) external {
    autoLog = _autoLog;
  }

  function setUseMercury(bool _useMercury) external {
    useMercury = _useMercury;
  }

  function setFeedsHex(string[] memory newFeeds) external {
    feedsHex = newFeeds;
  }

  function checkLog(Log calldata log, bytes memory checkData) external returns (bool, bytes memory) {
    uint256 startGas = gasleft();
    uint256 blockNum = getBlockNumber();

    // filter by event signature
    if (log.topics[0] == emittedSig) {
      bytes memory t1 = abi.encodePacked(log.topics[1]); // bytes32 to bytes
      uint256 upkeepId = abi.decode(t1, (uint256));
      bytes memory t2 = abi.encodePacked(log.topics[2]);
      uint256 blockNum = abi.decode(t2, (uint256));

      uint256 checkGasToBurn = checkGasToBurns[upkeepId];
      while (startGas - gasleft() + 15000 < checkGasToBurn) {
        dummyMap[blockhash(blockNum)] = false;
      }

      if (useMercury) {
        revert FeedLookup(feedParamKey, feedsHex, timeParamKey, blockNum, abi.encode(upkeepId, blockNum));
      }

      // if we don't use mercury, create a perform data which resembles the output of checkCallback
      bytes[] memory values = new bytes[](1);
      bytes memory extraData = abi.encode(upkeepId, blockNum);
      return (true, abi.encode(values, extraData));
    }
    revert("could not find matching event sig");
  }

  function performUpkeep(bytes calldata performData) external {
    uint256 startGas = gasleft();
    (bytes[] memory values, bytes memory extraData) = abi.decode(performData, (bytes[], bytes));
    (uint256 upkeepId, uint256 logBlockNumber) = abi.decode(extraData, (uint256, uint256));

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
    if (autoLog) {
      emit LogEmitted(upkeepId, currentBlockNum, address(this));
    }
    burnPerformGas(upkeepId, startGas, currentBlockNum);
  }

  function checkCallback(
    bytes[] memory values,
    bytes memory extraData
  ) external pure override returns (bool, bytes memory) {
    bytes memory performData = abi.encode(values, extraData);
    return (true, performData);
  }
}
