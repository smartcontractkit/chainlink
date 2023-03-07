pragma solidity ^0.8.6;

contract UpkeepCounterStats {
  event PerformingUpkeep(
    uint256 initialBlock,
    uint256 lastBlock,
    uint256 previousBlock,
    uint256 counter,
    bytes performData
  );

  mapping(uint256 => uint256) public upkeepIdsToIntervals;
  mapping(uint256 => uint256) public upkeepIdsToLastBlock;
  mapping(uint256 => uint256) public upkeepIdsToPreviousPerformBlock;
  mapping(uint256 => uint256) public upkeepIdsToInitialBlock;
  mapping(uint256 => uint256) public upkeepIdsToCounter;
  mapping(uint256 => uint256) public upkeepIdsToPerformGasToBurn;
  mapping(uint256 => uint256) public upkeepIdsToCheckGasToBurn;
  mapping(uint256 => uint256) public upkeepIdsToPerformDataSize;
  uint256 public interval;
  mapping(bytes32 => bool) public dummyMap; // used to force storage lookup

  mapping(uint256 => uint256[]) private upkeepIdsToDelay;
  uint256[] private delays;

  constructor(uint256 _interval) {
    interval = _interval;
  }

  function checkUpkeep(bytes calldata checkData) external returns (bool, bytes memory) {
    uint256 startGas = gasleft();

    (uint256 upkeepId, uint256 interval, uint256 checkGasToBurn, uint256 performGasToBurn, uint256 performDataSize) = abi.decode(
      checkData,
      (uint256, uint256, uint256, uint256, uint256)
    );

    upkeepIdsToIntervals[upkeepId] = interval;
    upkeepIdsToCheckGasToBurn[upkeepId] = checkGasToBurn;
    upkeepIdsToPerformGasToBurn[upkeepId] = performGasToBurn;
    upkeepIdsToPerformDataSize[upkeepId] = performDataSize;
    bytes memory pData = bytes.concat(abi.encode(upkeepId, performGasToBurn), new bytes(performDataSize));
    uint256 blockNum = block.number;
    bool needed = eligible(upkeepId);
    while (startGas - gasleft() + 10000 < checkGasToBurn) {
      // 10K margin over gas to burn
      // Hard coded check gas to burn
      dummyMap[blockhash(blockNum)] = false; // arbitrary storage writes
      blockNum--;
    }
    return (needed, pData);
  }

  function performUpkeep(bytes calldata performData) external {
    uint256 startGas = gasleft();
    (uint256 upkeepId, uint256 performGasToBurn, bytes memory performDataPlaceHolder) = abi.decode(
      performData,
      (uint256, uint256, bytes)
    );
    uint256 initialBlock = upkeepIdsToInitialBlock[upkeepId];
    uint256 blockNum = block.number;
    uint256 interval = upkeepIdsToIntervals[upkeepId];
    if (initialBlock == 0) {
      upkeepIdsToInitialBlock[upkeepId] = blockNum;
      initialBlock = blockNum;
    } else {
      // Calculate and append delay
      uint256 delay = block.number - upkeepIdsToPreviousPerformBlock[upkeepId] - interval;
      upkeepIdsToDelay[upkeepId].push(delay);
      upkeepIdsToDelay[upkeepId] = delays;
    }

    upkeepIdsToLastBlock[upkeepId] = blockNum;
    uint256 counter = upkeepIdsToCounter[upkeepId] + 1;
    upkeepIdsToCounter[upkeepId] = counter;
    emit PerformingUpkeep(initialBlock, blockNum, upkeepIdsToPreviousPerformBlock[upkeepId], counter, performData);
    upkeepIdsToPreviousPerformBlock[upkeepId] = blockNum;

    while (startGas - gasleft() + 10000 < performGasToBurn) {
      // 10K margin over gas to burn
      dummyMap[blockhash(blockNum)] = false; // arbitrary storage writes
      blockNum--;
    }
  }

  function eligible(uint256 upkeepId) public view returns (bool) {
    if (upkeepIdsToInitialBlock[upkeepId] == 0) {
      return true;
    }
    return (block.number - upkeepIdsToLastBlock[upkeepId]) >= upkeepIdsToIntervals[upkeepId];
  }

  function setPerformGasToBurn(uint256 upkeepId, uint256 value) public {
    upkeepIdsToPerformGasToBurn[upkeepId] = value;
  }

  function setCheckGasToBurn(uint256 upkeepId, uint256 value) public {
    upkeepIdsToCheckGasToBurn[upkeepId] = value;
  }

  function setPerformDataSize(uint256 upkeepId, uint256 value) public {
    upkeepIdsToPerformDataSize[upkeepId] = value;
  }

  function setSpread(uint256 upkeepId, uint256 _interval) external {
    upkeepIdsToIntervals[upkeepId] = _interval;
    upkeepIdsToInitialBlock[upkeepId] = 0;
    upkeepIdsToCounter[upkeepId] = 0;

    delete upkeepIdsToDelay[upkeepId];
  }

  function getDelaysLength(uint256 upkeepId) public view returns (uint256) {
    return upkeepIdsToDelay[upkeepId].length;
  }

  function getDelays(uint256 upkeepId) public view returns (uint256[] memory) {
    return upkeepIdsToDelay[upkeepId];
  }

  function getSumDelayLastNPerforms(uint256 upkeepId, uint256 n) public view returns (uint256, uint256) {
    uint256[] memory delays = upkeepIdsToDelay[upkeepId];
    uint256 i;
    uint256 len = delays.length;
    if (n == 0 || n >= len) {
      n = len;
    }
    uint256 sum = 0;

    for (i = 0; i < n; i++) sum = sum + delays[len - i - 1];
    return (sum, n);
  }

  function getPxDelayLastNPerforms(uint256 upkeepId, uint256 p, uint256 n) public view returns (uint256) {
    uint256[] memory delays = upkeepIdsToDelay[upkeepId];
    uint256 i;
    uint256 len = delays.length;
    if (n == 0 || n >= len) {
      n = len;
    }
    uint256[] memory subArr = new uint256[](n);

    for (i = 0; i < n; i++) subArr[i] = (delays[len - i - 1]);
    quickSort(subArr, int256(0), int256(subArr.length - 1));

    uint256 index = (p * subArr.length) / 100;
    return subArr[index];
  }

  function quickSort(
    uint256[] memory arr,
    int256 left,
    int256 right
  ) private pure {
    int256 i = left;
    int256 j = right;
    if (i == j) return;
    uint256 pivot = arr[uint256(left + (right - left) / 2)];
    while (i <= j) {
      while (arr[uint256(i)] < pivot) i++;
      while (pivot < arr[uint256(j)]) j--;
      if (i <= j) {
        (arr[uint256(i)], arr[uint256(j)]) = (arr[uint256(j)], arr[uint256(i)]);
        i++;
        j--;
      }
    }
    if (left < j) quickSort(arr, left, j);
    if (i < right) quickSort(arr, i, right);
  }
}
