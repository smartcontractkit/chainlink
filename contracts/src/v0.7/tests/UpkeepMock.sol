// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "../KeeperCompatible.sol";

contract UpkeepMock is KeeperCompatible {
  bool public shouldRevertCheck;
  bool public canCheck;
  bool public canPerform;
  bytes public performData;
  uint256 public checkGasToBurn;
  uint256 public performGasToBurn;

  uint256 constant checkGasBuffer = 6000; // use all but this amount in gas burn loops
  uint256 constant performGasBuffer = 1000; // use all but this amount in gas burn loops

  event UpkeepPerformedWith(bytes upkeepData);

  function setShouldRevertCheck(bool value) public {
    shouldRevertCheck = value;
  }

  function setPerformData(bytes calldata data) public {
    performData = data;
  }

  function setCanCheck(bool value) public {
    canCheck = value;
  }

  function setCanPerform(bool value) public {
    canPerform = value;
  }

  function setCheckGasToBurn(uint256 value) public {
    require(value > checkGasBuffer || value == 0, "checkGasToBurn must be 0 (disabled) or greater than buffer");
    if (value > 0) {
      checkGasToBurn = value - checkGasBuffer;
    } else {
      checkGasToBurn = 0;
    }
  }

  function setPerformGasToBurn(uint256 value) public {
    require(value > performGasBuffer || value == 0, "performGasToBurn must be 0 (disabled) or greater than buffer");
    if (value > 0) {
      performGasToBurn = value - performGasBuffer;
    } else {
      performGasToBurn = 0;
    }
  }

  function checkUpkeep(bytes calldata data)
    external
    override
    cannotExecute
    returns (bool callable, bytes memory executedata)
  {
    require(!shouldRevertCheck, "shouldRevertCheck should be false");
    uint256 startGas = gasleft();
    bool couldCheck = canCheck;

    setCanCheck(false); // test that state modifications don't stick

    while (startGas - gasleft() < checkGasToBurn) {} // burn gas

    return (couldCheck, performData);
  }

  function performUpkeep(bytes calldata data) external override {
    uint256 startGas = gasleft();

    require(canPerform, "Cannot perform");

    emit UpkeepPerformedWith(data);

    while (startGas - gasleft() < performGasToBurn) {} // burn gas
  }
}
