// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

contract ReceiveEmitter {
  event FundsReceived(uint256 amount, uint256 newBalance);

  receive() external payable {
    emit FundsReceived(msg.value, address(this).balance);
  }
}
