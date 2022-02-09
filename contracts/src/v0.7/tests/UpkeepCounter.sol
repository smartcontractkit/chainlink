pragma solidity ^0.7.6;

contract UpkeepCounter {
  uint256 public counter;
  uint256 public testRange;
  uint256 public interval;
  uint256 public lastBlock;

  constructor(uint256 _testRange, uint256 _interval) {
    testRange = _testRange;
    interval = _interval;
    lastBlock = block.number;
    counter = 0;
  }

  function checkUpkeep(bytes calldata data) external view returns (bool, bytes memory) {
    return (counter < testRange && (block.number - lastBlock) > interval, data);
  }

  function performUpkeep(bytes calldata performData) external {
    lastBlock = block.number;
    counter = counter + 1;
    performData;
  }

  function reset(uint256 _testRange, uint256 _interval) external {
    testRange = _testRange;
    interval = _interval;
    counter = 0;
  }
}
