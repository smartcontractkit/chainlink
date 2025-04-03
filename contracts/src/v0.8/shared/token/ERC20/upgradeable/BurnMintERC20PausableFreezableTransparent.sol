// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {BurnMintERC20PausableTransparent} from "./BurnMintERC20PausableTransparent.sol";

contract BurnMintERC20PausableFreezableTransparent is BurnMintERC20PausableTransparent {
  event AccountFrozen(address indexed account);
  event AccountUnfrozen(address indexed account);

  error BurnMintERC20PausableFreezableTransparent__InvalidRecipient(address recipient);
  error BurnMintERC20PausableFreezableTransparent__AccountFrozen(address account);
  error BurnMintERC20PausableFreezableTransparent__AccountNotFrozen(address account);

  bytes32 public constant FREEZER_ROLE = keccak256("FREEZER_ROLE");

  // ================================================================
  // │                          Storage                             │
  // ================================================================

  /// @custom:storage-location erc7201:chainlink.storage.BurnMintERC20PausableFreezableTransparent
  struct BurnMintERC20PausableFreezableTransparentStorage {
    /// @dev Mapping to keep track of the frozen status of an address
    mapping(address => bool) isFrozen;
  }

  // keccak256(abi.encode(uint256(keccak256("chainlink.storage.BurnMintERC20PausableFreezableTransparent")) - 1)) & ~bytes32(uint256(0xff));
  bytes32 private constant BURN_MINT_ERC20_PAUSABLE_FREEZABLE_TRANSPARENT_STORAGE_LOCATION =
    0xe4a0d511ce93f7d3bf378a3a2c82dfeda12e9faf72c0533ddcd2be06e2d60f00;

  // solhint-disable-next-line chainlink-solidity/explicit-returns
  function _getBurnMintERC20PausableFreezableTransparentStorage()
    private
    pure
    returns (BurnMintERC20PausableFreezableTransparentStorage storage $)
  {
    assembly {
      $.slot := BURN_MINT_ERC20_PAUSABLE_FREEZABLE_TRANSPARENT_STORAGE_LOCATION
    }
  }

  // ================================================================
  // │                         Freezing                             │
  // ================================================================

  /// @notice Freezes an account, disallowing transfers, minting and burning from/to it.
  /// @dev Requires the caller to have the FREEZER_ROLE.
  /// @dev Can be called even if the contract is paused.
  function freeze(address account) public onlyRole(FREEZER_ROLE) {
    if (account == address(0)) revert BurnMintERC20PausableFreezableTransparent__InvalidRecipient(account);
    if (account == address(this)) revert BurnMintERC20PausableFreezableTransparent__InvalidRecipient(account);

    BurnMintERC20PausableFreezableTransparentStorage storage $ = _getBurnMintERC20PausableFreezableTransparentStorage();
    if ($.isFrozen[account]) revert BurnMintERC20PausableFreezableTransparent__AccountFrozen(account);

    $.isFrozen[account] = true;

    emit AccountFrozen(account);
  }

  /// @notice Unfreezes an account
  /// @dev Requires the caller to have the FREEZER_ROLE.
  /// @dev Can be called even if the contract is paused.
  function unfreeze(address account) public onlyRole(FREEZER_ROLE) {
    BurnMintERC20PausableFreezableTransparentStorage storage $ = _getBurnMintERC20PausableFreezableTransparentStorage();
    if (!$.isFrozen[account]) revert BurnMintERC20PausableFreezableTransparent__AccountNotFrozen(account);

    $.isFrozen[account] = false;

    emit AccountUnfrozen(account);
  }

  function isFrozen(address account) public view returns (bool) {
    return _getBurnMintERC20PausableFreezableTransparentStorage().isFrozen[account];
  }

  // ================================================================
  // │                            ERC20                             │
  // ================================================================

  /// @dev Uses BurnMintERC20PausableTransparent _update hook to disallow transfers, minting and burning from/to frozen addresses.
  function _update(address from, address to, uint256 value) internal virtual override {
    BurnMintERC20PausableFreezableTransparentStorage storage $ = _getBurnMintERC20PausableFreezableTransparentStorage();
    if ($.isFrozen[from]) revert BurnMintERC20PausableFreezableTransparent__AccountFrozen(from);
    if ($.isFrozen[to]) revert BurnMintERC20PausableFreezableTransparent__AccountFrozen(to);

    super._update(from, to, value);
  }

  /// @dev Uses BurnMintERC20PausableTransparent _approve to disallow approving from and to frozen addresses.
  function _approve(address owner, address spender, uint256 value, bool emitEvent) internal virtual override {
    BurnMintERC20PausableFreezableTransparentStorage storage $ = _getBurnMintERC20PausableFreezableTransparentStorage();
    if ($.isFrozen[owner]) revert BurnMintERC20PausableFreezableTransparent__AccountFrozen(owner);
    if ($.isFrozen[spender]) revert BurnMintERC20PausableFreezableTransparent__AccountFrozen(spender);

    super._approve(owner, spender, value, emitEvent);
  }
}
