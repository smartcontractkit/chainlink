// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IBurnMintERC20} from "../../shared/token/ERC20/IBurnMintERC20.sol";

import {BurnMintTokenPoolAbstract} from "./BurnMintTokenPoolAbstract.sol";
import {TokenPool} from "./TokenPool.sol";

/// @notice This pool mints and burns a 3rd-party token.
/// @dev Pool whitelisting mode is set in the constructor and cannot be modified later.
/// It either accepts any address as originalSender, or only accepts whitelisted originalSender.
/// The only way to change whitelisting mode is to deploy a new pool.
/// If that is expected, please make sure the token's burner/minter roles are adjustable.
/// @dev This contract is a variant of BurnMintTokenPool that uses `burn(amount)`.
contract BurnMintTokenPool is BurnMintTokenPoolAbstract, ITypeAndVersion {
  string public constant override typeAndVersion = "BurnMintTokenPool 1.5.0";

  constructor(
    IBurnMintERC20 token,
    address[] memory allowlist,
    address rmnProxy,
    address router
  ) TokenPool(token, allowlist, rmnProxy, router) {}

  /// @inheritdoc BurnMintTokenPoolAbstract
  function _burn(
    uint256 amount
  ) internal virtual override {
    IBurnMintERC20(address(i_token)).burn(amount);
  }
}
