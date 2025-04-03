// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {BurnMintERC20PausableUUPS} from "./BurnMintERC20PausableUUPS.sol";

contract BurnMintERC20PausableFreezableUUPS is BurnMintERC20PausableUUPS {
  event AccountFrozen(address indexed account);
  event AccountUnfrozen(address indexed account);

  error BurnMintERC20PausableFreezableUUPS__InvalidRecipient(address recipient);
  error BurnMintERC20PausableFreezableUUPS__AccountFrozen(address account);
  error BurnMintERC20PausableFreezableUUPS__AccountNotFrozen(address account);

  bytes32 public constant FREEZER_ROLE = keccak256("FREEZER_ROLE");

  // ================================================================
  // │                          Storage                             │
  // ================================================================

  /// @custom:storage-location erc7201:chainlink.storage.BurnMintERC20PausableFreezableUUPS
  struct BurnMintERC20PausableFreezableUUPSStorage {
    /// @dev Mapping to keep track of the frozen status of an address
    mapping(address => bool) isFrozen;
  }

  // keccak256(abi.encode(uint256(keccak256("chainlink.storage.BurnMintERC20PausableFreezableUUPS")) - 1)) & ~bytes32(uint256(0xff));
  bytes32 private constant BURN_MINT_ERC20_PAUSABLE_FREEZABLE_UUPS_STORAGE_LOCATION =
    0x36a30f686feb055c8d90421e230dafb8f47433e358189345608518a408badc00;

  // solhint-disable-next-line chainlink-solidity/explicit-returns
  function _getBurnMintERC20PausableFreezableUUPSStorage()
    private
    pure
    returns (BurnMintERC20PausableFreezableUUPSStorage storage $)
  {
    assembly {
      $.slot := BURN_MINT_ERC20_PAUSABLE_FREEZABLE_UUPS_STORAGE_LOCATION
    }
  }

  // ================================================================
  // │                         Freezing                             │
  // ================================================================

  /// @notice Freezes an account, disallowing transfers, minting and burning from/to it.
  /// @dev Requires the caller to have the FREEZER_ROLE.
  /// @dev Can be called even if the contract is paused.
  function freeze(address account) public onlyRole(FREEZER_ROLE) {
    if (account == address(0)) revert BurnMintERC20PausableFreezableUUPS__InvalidRecipient(account);
    if (account == address(this)) revert BurnMintERC20PausableFreezableUUPS__InvalidRecipient(account);

    BurnMintERC20PausableFreezableUUPSStorage storage $ = _getBurnMintERC20PausableFreezableUUPSStorage();
    if ($.isFrozen[account]) revert BurnMintERC20PausableFreezableUUPS__AccountFrozen(account);

    $.isFrozen[account] = true;

    emit AccountFrozen(account);
  }

  /// @notice Unfreezes an account
  /// @dev Requires the caller to have the FREEZER_ROLE.
  /// @dev Can be called even if the contract is paused.
  function unfreeze(address account) public onlyRole(FREEZER_ROLE) {
    BurnMintERC20PausableFreezableUUPSStorage storage $ = _getBurnMintERC20PausableFreezableUUPSStorage();
    if (!$.isFrozen[account]) revert BurnMintERC20PausableFreezableUUPS__AccountNotFrozen(account);

    $.isFrozen[account] = false;

    emit AccountUnfrozen(account);
  }

  function isFrozen(address account) public view returns (bool) {
    return _getBurnMintERC20PausableFreezableUUPSStorage().isFrozen[account];
  }

  // ================================================================
  // │                            ERC20                             │
  // ================================================================

  /// @dev Uses BurnMintERC20PausableUUPS _update hook to disallow transfers, minting and burning from/to frozen addresses.
  function _update(address from, address to, uint256 value) internal virtual override {
    BurnMintERC20PausableFreezableUUPSStorage storage $ = _getBurnMintERC20PausableFreezableUUPSStorage();
    if ($.isFrozen[from]) revert BurnMintERC20PausableFreezableUUPS__AccountFrozen(from);
    if ($.isFrozen[to]) revert BurnMintERC20PausableFreezableUUPS__AccountFrozen(to);

    super._update(from, to, value);
  }

  /// @dev Uses BurnMintERC20PausableUUPS _approve to disallow approving from and to frozen addresses.
  function _approve(address owner, address spender, uint256 value, bool emitEvent) internal virtual override {
    BurnMintERC20PausableFreezableUUPSStorage storage $ = _getBurnMintERC20PausableFreezableUUPSStorage();
    if ($.isFrozen[owner]) revert BurnMintERC20PausableFreezableUUPS__AccountFrozen(owner);
    if ($.isFrozen[spender]) revert BurnMintERC20PausableFreezableUUPS__AccountFrozen(spender);

    super._approve(owner, spender, value, emitEvent);
  }
}
