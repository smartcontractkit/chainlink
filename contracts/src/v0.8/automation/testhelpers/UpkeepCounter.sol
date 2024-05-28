// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

contract UpkeepCounter {
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

  constructor(uint256 _testRange, uint256 _interval) {
    testRange = _testRange;
    interval = _interval;
    previousPerformBlock = 0;
    lastTimestamp = block.timestamp;
    initialTimestamp = 0;
    counter = 0;
  }

  function checkUpkeep(bytes calldata data) external view returns (bool, bytes memory) {
    return (eligible(), data);
  }

  function performUpkeep(bytes calldata performData) external {
    if (initialTimestamp == 0) {
      initialTimestamp = block.timestamp;
    }
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
