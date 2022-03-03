// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "../KeeperCompatible.sol";

contract UpkeepMock is KeeperCompatible {
  bool public canCheck;
  bool public canPerform;
  uint256 public checkGasToBurn;
  uint256 public performGasToBurn;

  uint256 constant gasBuffer = 1000; // use all but this amount in gas burn loops

  event UpkeepPerformedWith(bytes upkeepData);

  function setCanCheck(bool value) public {
    canCheck = value;
  }

  function setCanPerform(bool value) public {
    canPerform = value;
  }

  function setCheckGasToBurn(uint256 value) public {
    require(value > gasBuffer || value == 0, "checkGasToBurn must be 0 (disabled) or greater than buffer");
    checkGasToBurn = value - gasBuffer;
  }

  function setPerformGasToBurn(uint256 value) public {
    require(value > gasBuffer || value == 0, "performGasToBurn must be 0 (disabled) or greater than buffer");
    performGasToBurn = value - gasBuffer;
  }

  function checkUpkeep(bytes calldata data)
    external
    override
    cannotExecute
    returns (bool callable, bytes calldata executedata)
  {
    uint256 startGas = gasleft();
    bool couldCheck = canCheck;

    setCanCheck(false); // test that state modifications don't stick

    while (startGas - gasleft() < checkGasToBurn) {} // burn gas

    return (couldCheck, data);
  }

  function performUpkeep(bytes calldata data) external override {
    uint256 startGas = gasleft();

    require(canPerform, "Cannot perform");

    setCanPerform(false);

    emit UpkeepPerformedWith(data);

    while (startGas - gasleft() < performGasToBurn) {} // burn gas
  }
}
