// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

contract ReceiveReverter {
  receive() external payable {
    revert("Can't send funds");
  }
}
