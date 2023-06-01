// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

contract LogUpkeepCounter {
  event PerformingUpkeep(
    address indexed from,
    uint256 initialBlock,
    uint256 lastBlock,
    uint256 previousBlock,
    uint256 counter
  );

  /**
   * @dev we include multiple event types for testing various filters, signatures, etc
   */
  event Trigger(); // 0x3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d
  event Trigger(uint256 a); // 0x57b1de35764b0939dde00771c7069cdf8d6a65d6a175623f19aa18784fd4c6da
  event Trigger(uint256 a, uint256 b); // 0x1da9f70fe932e73fba9374396c5c0b02dbd170f951874b7b4afabe4dd029a9c8
  event Trigger(uint256 a, uint256 b, uint256 c); // 0x5121119bad45ca7e58e0bdadf39045f5111e93ba4304a0f6457a3e7bc9791e71

  uint256 public testRange;
  uint256 public lastBlock;
  uint256 public previousPerformBlock;
  uint256 public initialBlock;
  uint256 public counter;

  constructor(uint256 _testRange) {
    testRange = _testRange;
    previousPerformBlock = 0;
    lastBlock = block.number;
    initialBlock = 0;
    counter = 0;
  }

  function start() public {
    // need an initial event to begin the cycle
    emit Trigger();
    emit Trigger(1);
    emit Trigger(1, 2);
    emit Trigger(1, 2, 3);
  }

  function checkUpkeep(bytes calldata logData) external view returns (bool, bytes memory) {
    require(eligible(), "not eligible");
    bytes4 sig = abi.decode(logData[:4], (bytes4));
    if (sig == 0x3d53a395 || sig == 0x57b1de35 || sig == 0x1da9f70f || sig == 0x5121119b) {
      return (true, logData);
    } else {
      revert("could not find matching event sig");
    }
  }

  function performUpkeep(bytes calldata performData) external {
    if (initialBlock == 0) {
      initialBlock = block.number;
    }
    lastBlock = block.number;
    counter = counter + 1;
    previousPerformBlock = lastBlock;
    bytes4 sig = abi.decode(performData[:4], (bytes4));
    if (sig == 0x3d53a395) {
      emit Trigger();
    } else if (sig == 0x57b1de35) {
      emit Trigger(1);
    } else if (sig == 0x1da9f70f) {
      emit Trigger(1, 2);
    } else if (sig == 0x5121119b) {
      emit Trigger(1, 2, 3);
    } else {
      revert("could not find matching sig");
    }
    emit PerformingUpkeep(tx.origin, initialBlock, lastBlock, previousPerformBlock, counter);
  }

  function eligible() public view returns (bool) {
    if (initialBlock == 0) {
      return true;
    }

    return (block.number - initialBlock) < testRange;
  }

  function setSpread(uint256 _testRange) external {
    testRange = _testRange;
    initialBlock = 0;
    counter = 0;
  }
}
