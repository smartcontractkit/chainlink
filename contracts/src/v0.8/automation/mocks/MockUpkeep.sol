// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MockUpkeep {
  bool public shouldCheckRevert;
  bool public shouldPerformRevert;
  bool public checkResult = true;
  bytes public performData;
  uint256 public checkGasToBurn;
  uint256 public performGasToBurn;

  event UpkeepPerformedWith(bytes upkeepData);
  error CheckRevert();
  error PerformRevert();

  function setShouldCheckRevert(bool value) public {
    shouldCheckRevert = value;
  }

  function setShouldPerformRevert(bool value) public {
    shouldPerformRevert = value;
  }

  function setCheckResult(bool value) public {
    checkResult = value;
  }

  function setPerformData(bytes calldata data) public {
    performData = data;
  }

  function setCheckGasToBurn(uint256 value) public {
    checkGasToBurn = value;
  }

  function setPerformGasToBurn(uint256 value) public {
    performGasToBurn = value;
  }

  function checkUpkeep(bytes calldata) external view returns (bool callable, bytes memory executedata) {
    if (shouldCheckRevert) revert CheckRevert();
    uint256 startGas = gasleft();
    while (startGas - gasleft() < checkGasToBurn) {} // burn gas
    return (checkResult, performData);
  }

  function performUpkeep(bytes calldata data) external {
    if (shouldPerformRevert) revert PerformRevert();
    uint256 startGas = gasleft();
    while (startGas - gasleft() < performGasToBurn) {} // burn gas
    emit UpkeepPerformedWith(data);
  }
}
