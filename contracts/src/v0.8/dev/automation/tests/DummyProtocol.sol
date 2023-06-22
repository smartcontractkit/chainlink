// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

contract DummyProtocol {
  event LimitOrderSent(uint256 indexed amount, uint256 indexed price, address indexed to); // keccak256(LimitOrderSent(uint256,uint256,address)) => 0x3e9c37b3143f2eb7e9a2a0f8091b6de097b62efcfe48e1f68847a832e521750a
  event LimitOrderWithdrawn(uint256 indexed amount, uint256 indexed price, address indexed from); // keccak256(LimitOrderWithdrawn(uint256,uint256,address)) => 0x0a71b8ed921ff64d49e4d39449f8a21094f38a0aeae489c3051aedd63f2c229f
  event LimitOrderExecuted(uint256 indexed orderId, uint256 indexed amount, address indexed exchange); // keccak(LimitOrderExecuted(uint256,uint256,address)) => 0xd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd

  function sendLimitedOrder(uint256 amount, uint256 price, address to) public {
    // send an order to an exchange
    emit LimitOrderSent(amount, price, to);
  }

  function withdrawLimit(uint256 amount, uint256 price, address from) public {
    // withdraw an order from an exchange
    emit LimitOrderSent(amount, price, from);
  }

  function executeLimitOrder(uint256 orderId, uint256 amount, address exchange) public {
    // execute a limit order
    emit LimitOrderExecuted(orderId, amount, exchange);
  }
}
