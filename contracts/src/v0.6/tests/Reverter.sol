// SPDX-License-Identifier: MIT
pragma solidity 0.6.6;

contract Reverter {

  fallback() external payable {
    require(false, "Raised by Reverter.sol");
  }

}
