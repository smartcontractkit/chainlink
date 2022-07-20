// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

/**
 * @notice getRevertMsg extracts a revert reason from a failed contract call
 */
function getRevertMsg(bytes memory payload) pure returns (string memory) {
  if (payload.length < 68) return "transaction reverted silently";
  assembly {
    payload := add(payload, 0x04)
  }
  return abi.decode(payload, (string));
}
