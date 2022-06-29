// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IAccessController {
  function hasAccess(address user, bytes calldata data) external view returns (bool);
}
