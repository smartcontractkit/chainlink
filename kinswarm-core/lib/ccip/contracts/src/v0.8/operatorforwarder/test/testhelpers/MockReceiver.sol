// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

contract MockReceiver {
  uint256 private s_value;

  function receiveData(uint256 _value) public {
    s_value = _value;
  }

  function revertMessage() public pure {
    // solhint-disable-next-line gas-custom-errors
    revert("test revert message");
  }

  function getValue() external view returns (uint256) {
    return s_value;
  }
}
