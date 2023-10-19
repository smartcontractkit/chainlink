// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import "./VerifiableLoadBase.sol";
import "../automation/interfaces/StreamsLookupCompatibleInterface.sol";

contract VerifiableLoadStreamsLookupUpkeep is VerifiableLoadBase, StreamsLookupCompatibleInterface {
  constructor(AutomationRegistrar2_1 _registrar, bool _useArb) VerifiableLoadBase(_registrar, _useArb) {}

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

    uint256 timeParam;
    if (keccak256(abi.encodePacked(feedParamKey)) == keccak256(abi.encodePacked("feedIdHex"))) {
      timeParam = blockNum;
    } else {
      // assume this will be feedIDs for v0.3
      timeParam = block.timestamp;
    }

    revert StreamsLookup(feedParamKey, feedsHex, timeParamKey, timeParam, abi.encode(upkeepId));
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
    } else {
      uint256 delay = blockNum - previousPerformBlock - intervals[upkeepId];
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
    previousPerformBlocks[upkeepId] = blockNum;

    topUpFund(upkeepId, blockNum);
    burnPerformGas(upkeepId, startGas, blockNum);
  }
}
