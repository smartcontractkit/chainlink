// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

contract ReceiveFallbackEmitter {
  event FundsReceived(uint256 amount, uint256 newBalance);

  fallback() external payable {
    emit FundsReceived(msg.value, address(this).balance);
  }
}
