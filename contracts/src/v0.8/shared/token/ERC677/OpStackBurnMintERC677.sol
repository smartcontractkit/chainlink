// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IOptimismMintableERC20Minimal, IOptimismMintableERC20} from "../ERC20/IOptimismMintableERC20.sol";

import {BurnMintERC677} from "./BurnMintERC677.sol";

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/IERC165.sol";

/// @notice A basic ERC677 compatible token contract with burn and minting roles that supports
/// the native L2 bridging requirements of the Optimism Stack.
/// @dev Note: the L2 bridge contract needs to be given burn and mint privileges manually,
/// since this contract does not automatically grant them. This allows the owner to revoke
/// the bridge's privileges if necessary.
contract OpStackBurnMintERC677 is BurnMintERC677, IOptimismMintableERC20Minimal {
  /// @dev The address of the L1 token.
  address internal immutable i_l1Token;
  /// @dev The address of the L2 bridge.
  address internal immutable i_l2Bridge;

  constructor(
    string memory name,
    string memory symbol,
    uint8 decimals_,
    uint256 maxSupply_,
    address l1Token,
    address l2Bridge
  ) BurnMintERC677(name, symbol, decimals_, maxSupply_) {
    i_l1Token = l1Token;
    i_l2Bridge = l2Bridge;
  }

  function supportsInterface(bytes4 interfaceId) public pure virtual override(IERC165, BurnMintERC677) returns (bool) {
    return
      interfaceId == type(IOptimismMintableERC20).interfaceId ||
      interfaceId == type(IOptimismMintableERC20Minimal).interfaceId ||
      super.supportsInterface(interfaceId);
  }

  /// @notice Returns the address of the L1 token.
  function remoteToken() public view override returns (address) {
    return i_l1Token;
  }

  /// @notice Returns the address of the L2 bridge.
  function bridge() public view override returns (address) {
    return i_l2Bridge;
  }
}
