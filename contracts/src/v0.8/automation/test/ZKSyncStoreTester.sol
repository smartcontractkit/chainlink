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
  bytes[] public data;

  constructor() {
    testRange = 10000;
    interval = 200;
    previousPerformBlock = 0;
    lastTimestamp = block.timestamp;
    initialTimestamp = 0;
    counter = 0;

    iterations = 100;
  }

  function storeData() public {
    data = new bytes[](iterations);
    storedData = new bytes(iterations);
    bytes1 d = 0xff;
    if (reset) {
      d = 0x00;
    }
    for (uint32 i = 0; i < iterations; i++) {
      storedData[i] = d;
    }
    reset = !reset;
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
