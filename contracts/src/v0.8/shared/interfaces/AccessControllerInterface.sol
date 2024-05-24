// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// solhint-disable-next-line interface-starts-with-i
interface AccessControllerInterface {
  function hasAccess(address user, bytes calldata data) external view returns (bool);
}
