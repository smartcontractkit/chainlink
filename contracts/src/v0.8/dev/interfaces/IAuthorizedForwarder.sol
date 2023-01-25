// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

interface IAuthorizedForwarder {
  function forward(address to, bytes calldata data) external;
}
