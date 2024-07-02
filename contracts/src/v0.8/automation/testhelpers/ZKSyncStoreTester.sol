// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

contract ZKSyncStoreTester {
  event PerformingUpkeep(
    address indexed from,
    uint256 initialTimestamp,
    uint256 lastTimestamp,
    uint256 previousBlock,
    uint256 counter
  );

  uint256 public testRange;
  uint256 public interval;
  uint256 public lastTimestamp;
  uint256 public previousPerformBlock;
  uint256 public initialTimestamp;
  uint256 public counter;

  uint32 public iterations;
  bytes public storedData;
  bool public reset;
  bytes public data0;
  bytes public data1;
  bytes public data2;
  bytes public data3;
  bytes public data4;
  bytes public data5;
  bytes public full = hex"ffffffff0f0ff0fffff0fffffff0fffffffffff0fffff0ffffff0ffffffffffff00000fffffffffffff000fffffffffffffffffffffffff0ffffffffff0fffffffffffffffffff00fffffffffff00ffffffffffffffffffffffff000000fffffffffffffffffffffffffffffffffff0ffffffffffffffff0fffff0ffffffff";
  bytes public constant empty = hex"00";

  constructor() {
    testRange = 10000;
    interval = 200;
    previousPerformBlock = 0;
    lastTimestamp = block.timestamp;
    initialTimestamp = 0;
    counter = 0;

    iterations = 1;
  }

  function storeData() public {
    for (uint32 i = 0; i < iterations; i++) {
      if (reset) {
        data0 = empty;
        data1 = empty;
        data2 = empty;
        data3 = empty;
        data4 = empty;
        data5 = empty;
      } else {
        data0 = full;
        data1 = full;
        data2 = full;
        data3 = full;
        data4 = full;
        data5 = full;
      }
    }
    reset = !reset;
  }

  function setFull(bytes calldata d) external {
    full = d;
  }

  function setIterations(uint32 _i) external {
    iterations = _i;
  }

  function checkUpkeep(bytes calldata data) external view returns (bool, bytes memory) {
    return (eligible(), data);
  }

  function performUpkeep(bytes calldata) external {
    if (initialTimestamp == 0) {
      initialTimestamp = block.timestamp;
    }
    storeData();
    lastTimestamp = block.timestamp;
    counter = counter + 1;
    emit PerformingUpkeep(tx.origin, initialTimestamp, lastTimestamp, previousPerformBlock, counter);
    previousPerformBlock = lastTimestamp;
  }

  function eligible() public view returns (bool) {
    if (initialTimestamp == 0) {
      return true;
    }

    return (block.timestamp - initialTimestamp) < testRange && (block.timestamp - lastTimestamp) >= interval;
  }

  function setSpread(uint256 _testRange, uint256 _interval) external {
    testRange = _testRange;
    interval = _interval;
    initialTimestamp = 0;
    counter = 0;
  }
}
