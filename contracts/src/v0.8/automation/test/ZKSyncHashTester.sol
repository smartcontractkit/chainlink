// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

contract ZKSyncHashTester {
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
  bytes public data;
  bytes32 public storedHash;

  constructor() {
    testRange = 10000;
    interval = 200;
    previousPerformBlock = 0;
    lastTimestamp = block.timestamp;
    initialTimestamp = 0;
    counter = 0;
    iterations = 100;
    data = "0xff";
  }

  function storeHash() public {
    bytes32 h = keccak256(data);
    for (uint32 i = 0; i < iterations - 1; i++) {
      h = keccak256(abi.encode(h));
    }
    storedHash = h;
  }

  function setIterations(uint32 _i) external {
    iterations = _i;
  }

  function setData(bytes calldata _data) external {
    data = _data;
  }

  function checkUpkeep(bytes calldata _data) external view returns (bool, bytes memory) {
    return (eligible(), _data);
  }

  function performUpkeep(bytes calldata performData) external {
    if (initialTimestamp == 0) {
      initialTimestamp = block.timestamp;
    }
    storeHash();
    lastTimestamp = block.timestamp;
    counter = counter + 1;
    performData;
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
