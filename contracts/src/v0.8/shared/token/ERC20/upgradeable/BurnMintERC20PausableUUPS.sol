// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {PausableUpgradeable} from "../../../../vendor/openzeppelin-solidity-upgradeable/v5.0.2/contracts/utils/PausableUpgradeable.sol";
import {BurnMintERC20UUPS} from "./BurnMintERC20UUPS.sol";

contract BurnMintERC20PausableUUPS is BurnMintERC20UUPS, PausableUpgradeable {
  bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");

  // ================================================================
  // │                          Pausing                             │
  // ================================================================

  /// @notice Pauses the implementation.
  /// @dev Requires the caller to have the PAUSER_ROLE.
  function pause() public onlyRole(PAUSER_ROLE) {
    _pause();

    emit Paused(msg.sender);
  }

  /// @notice Unpauses the implementation.
  /// @dev Requires the caller to have the DEFAULT_ADMIN_ROLE.
  function unpause() public onlyRole(DEFAULT_ADMIN_ROLE) {
    _unpause();

    emit Unpaused(msg.sender);
  }

  // ================================================================
  // │                            ERC20                             │
  // ================================================================

  /// @dev Disallows sending, minting and burning if implementation is paused.
  function _update(address from, address to, uint256 value) internal virtual override {
    _requireNotPaused();

    super._update(from, to, value);
  }

  /// @dev Disallows approving if implementation is paused.
  function _approve(address owner, address spender, uint256 value, bool emitEvent) internal virtual override {
    _requireNotPaused();

    super._approve(owner, spender, value, emitEvent);
  }
}
