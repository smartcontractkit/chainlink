// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../dev/interfaces/IAuthorizedForwarder.sol";

contract MockAuthorizedForwarder is IAuthorizedForwarder {
  event ForwardFuncCalled(address to, bytes data);

  function forward(address to, bytes calldata data) external override {
    to.call(data);
    emit ForwardFuncCalled(to, data);
  }
}
