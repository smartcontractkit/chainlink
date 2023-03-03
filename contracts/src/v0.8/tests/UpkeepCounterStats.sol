pragma solidity ^0.8.6;

contract UpkeepCounterStats {
  event PerformingUpkeep(
    uint256 initialBlock,
    uint256 lastBlock,
    uint256 previousBlock,
    uint256 counter,
    bytes performData
  );

  uint256 public interval;
  uint256 public lastBlock;
  uint256 public previousPerformBlock;
  uint256 public initialBlock;
  uint256 public counter;
  uint256 public performGasToBurn;
  uint256 public checkGasToBurn;
  mapping(bytes32 => bool) public dummyMap; // used to force storage lookup
  uint256 public performDataSize;

  uint256[] private delays;

  constructor(uint256 _interval) {
    interval = _interval;
    previousPerformBlock = 0;
    lastBlock = block.number;
    initialBlock = 0;
    counter = 0;
    performGasToBurn = 0;
    checkGasToBurn = 0;
  }

  function checkUpkeep(bytes calldata) external returns (bool, bytes memory) {
    uint256 startGas = gasleft();
    bytes memory pData = new bytes(performDataSize);
    uint256 blockNum = block.number;
    bool needed = eligible();
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
    if (initialBlock == 0) {
      initialBlock = block.number;
    } else {
      // Calculate and append delay
      uint256 delay = block.number - previousPerformBlock - interval;
      delays.push(delay);
    }

    lastBlock = block.number;
    counter = counter + 1;
    emit PerformingUpkeep(initialBlock, lastBlock, previousPerformBlock, counter, performData);
    previousPerformBlock = lastBlock;

    uint256 blockNum = block.number;
    while (startGas - gasleft() + 10000 < performGasToBurn) {
      // 10K margin over gas to burn
      dummyMap[blockhash(blockNum)] = false; // arbitrary storage writes
      blockNum--;
    }
  }

  function eligible() public view returns (bool) {
    if (initialBlock == 0) {
      return true;
    }
    return (block.number - lastBlock) >= interval;
  }

  function setPerformGasToBurn(uint256 value) public {
    performGasToBurn = value;
  }

  function setCheckGasToBurn(uint256 value) public {
    checkGasToBurn = value;
  }

  function setPerformDataSize(uint256 value) public {
    performDataSize = value;
  }

  function setSpread(uint256 _interval) external {
    interval = _interval;
    initialBlock = 0;
    counter = 0;

    uint256 n = delays.length;
    uint256 i;
    for (i = 0; i < n; i++) delays.pop();
  }

  function getDelaysLength() public view returns (uint256) {
    return delays.length;
  }

  function getDelays() public view returns (uint256[] memory) {
    return delays;
  }

  function getSumDelayLastNPerforms(uint256 n) public view returns (uint256, uint256) {
    uint256 i;
    uint256 len = delays.length;
    if (n == 0 || n >= len) {
      n = len;
    }
    uint256 sum = 0;

    for (i = 0; i < n; i++) sum = sum + delays[len - i - 1];
    return (sum, n);
  }

  function getPxDelayLastNPerforms(uint256 p, uint256 n) public view returns (uint256) {
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
