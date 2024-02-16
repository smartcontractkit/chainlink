// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IBurnMintERC20} from "../../shared/token/ERC20/IBurnMintERC20.sol";

import {TokenPool} from "./TokenPool.sol";
import {BurnMintTokenPoolAbstract} from "./BurnMintTokenPoolAbstract.sol";

import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @notice This pool mints and burns a 3rd-party token.
/// @dev Pool whitelisting mode is set in the constructor and cannot be modified later.
/// It either accepts any address as originalSender, or only accepts whitelisted originalSender.
/// The only way to change whitelisting mode is to deploy a new pool.
/// If that is expected, please make sure the token's burner/minter roles are adjustable.
contract BurnWithFromMintTokenPool is BurnMintTokenPoolAbstract, ITypeAndVersion {
  using SafeERC20 for IBurnMintERC20;

  string public constant override typeAndVersion = "BurnWithFromMintTokenPool 1.4.0";

  constructor(
    IBurnMintERC20 token,
    address[] memory allowlist,
    address armProxy,
    address router
  ) TokenPool(token, allowlist, armProxy, router) {
    // Some tokens allow burning from the sender without approval, but not all do.
    // To be safe, we approve the pool to burn from the pool.
    token.safeIncreaseAllowance(address(this), type(uint256).max);
  }

  /// @inheritdoc BurnMintTokenPoolAbstract
  function _burn(uint256 amount) internal virtual override {
    IBurnMintERC20(address(i_token)).burn(address(this), amount);
  }
}
