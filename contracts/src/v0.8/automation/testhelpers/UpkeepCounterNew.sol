// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

contract UpkeepCounterNew {
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
  bool public useMoreCheckGas;
  bool public useMorePerformGas;
  bool public useMorePerformData;
  uint256 public checkGasToBurn;
  uint256 public performGasToBurn;
  bytes public data;
  bytes public dataCopy;

  constructor() {
    testRange = 1000000;
    interval = 40;
    previousPerformBlock = 0;
    lastTimestamp = block.timestamp;
    initialTimestamp = 0;
    counter = 0;
    useMoreCheckGas = false;
    useMorePerformData = false;
    useMorePerformGas = false;
    checkGasToBurn = 5000000;
    performGasToBurn = 5000000;
  }

  function setPerformGasToBurn(uint256 value) external {
    performGasToBurn = value;
  }

  function setCheckGasToBurn(uint256 value) external {
    checkGasToBurn = value;
  }

  function setUseMoreCheckGas(bool _useMoreCheckGas) external {
    useMoreCheckGas = _useMoreCheckGas;
  }

  function setUseMorePerformGas(bool _useMorePerformGas) external {
    useMorePerformGas = _useMorePerformGas;
  }

  function setUseMorePerformData(bool _useMorePerformData) external {
    useMorePerformData = _useMorePerformData;
  }

  function setData(bytes calldata _data) external {
    data = _data;
  }

  function checkUpkeep(bytes calldata) external view returns (bool, bytes memory) {
    if (useMoreCheckGas) {
      uint256 startGas = gasleft();
      while (startGas - gasleft() < checkGasToBurn) {} // burn gas
    }

    return (eligible(), data);
  }

  function performUpkeep(bytes calldata performData) external {
    if (useMorePerformGas) {
      uint256 startGas = gasleft();
      while (startGas - gasleft() < performGasToBurn) {} // burn gas
    }

    if (initialTimestamp == 0) {
      initialTimestamp = block.timestamp;
    }
    lastTimestamp = block.timestamp;
    counter = counter + 1;
    dataCopy = performData;
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
