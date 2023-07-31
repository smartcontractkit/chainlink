// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import "./VerifiableLoadBase.sol";

contract VerifiableLoadUpkeep is VerifiableLoadBase {
  constructor(AutomationRegistrar2_1 _registrar, bool _useArb) VerifiableLoadBase(_registrar, _useArb) {}

  function checkUpkeep(bytes calldata checkData) external returns (bool, bytes memory) {
    uint256 startGas = gasleft();
    uint256 upkeepId = abi.decode(checkData, (uint256));

    uint256 performDataSize = performDataSizes[upkeepId];
    uint256 checkGasToBurn = checkGasToBurns[upkeepId];
    bytes memory pData = abi.encode(upkeepId, new bytes(performDataSize));
    uint256 blockNum = getBlockNumber();
    bool needed = eligible(upkeepId);
    while (startGas - gasleft() + 10000 < checkGasToBurn) {
      dummyMap[blockhash(blockNum)] = false;
      blockNum--;
    }
    return (needed, pData);
  }

  function performUpkeep(bytes calldata performData) external {
    uint256 startGas = gasleft();
    (uint256 upkeepId, ) = abi.decode(performData, (uint256, bytes));
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
