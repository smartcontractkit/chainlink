// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MockAuthorizedForwarder {
  event ForwardFuncCalled(address to, bytes data);

  function forward(address to, bytes calldata data) external {
    to.call(data);
    emit ForwardFuncCalled(to, data);
  }
}
