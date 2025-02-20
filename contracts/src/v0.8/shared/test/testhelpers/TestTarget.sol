// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

///
/// @notice A simple target contract to demonstrate success/failure paths.
///
contract TestTarget {
  error CustomRevertReason();

  // Returns some small data
  function returnData() external pure returns (string memory) {
    return "Hello from TestTarget";
  }

  // Returns ~200 bytes
  function returnLargeData() external pure returns (bytes memory) {
    bytes memory out = new bytes(200);
    for (uint256 i = 0; i < 200; i++) {
      out[i] = bytes1(uint8(65 + (i % 26))); // A..Z
    }
    return out;
  }

  // Reverts with a custom reason
  function revertWithReason() external pure {
    revert CustomRevertReason();
  }

  // Reverts with no reason
  function revertNoReason() external pure {
    // solhint-disable-next-line reason-string, gas-custom-errors
    revert();
  }
}
