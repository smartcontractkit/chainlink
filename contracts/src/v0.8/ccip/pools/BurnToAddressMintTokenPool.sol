// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IBurnMintERC20} from "../../shared/token/ERC20/IBurnMintERC20.sol";

import {Pool} from "../libraries/Pool.sol";
import {BurnMintTokenPoolAbstract} from "./BurnMintTokenPoolAbstract.sol";
import {TokenPool} from "./TokenPool.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @notice This pool mints and burns a 3rd-party token by sending tokens to an address which is unrecoverable
/// @dev Pool whitelisting mode is set in the constructor and cannot be modified later.
/// It either accepts any address as originalSender, or only accepts whitelisted originalSender.
/// The only way to change whitelisting mode is to deploy a new pool.
/// If that is expected, please make sure the token's burner/minter roles are adjustable.
/// @dev This contract is a variant of BurnMintTokenPool that uses `burn(amount)`.
contract BurnToAddressMintTokenPool is BurnMintTokenPoolAbstract, ITypeAndVersion {
  using SafeERC20 for IERC20;

  string public constant override typeAndVersion = "BurnToAddressTokenPool 1.5.1";

  address public immutable i_burnAddress;

  /// @notice Locked Tokens is a safety mechanism to ensure that more tokens cannot be sent out of the bridge
  /// than were originally sent in via CCIP.
  uint256 internal s_lockedTokens;

  /// @dev burnAddress is expected to be an address of which has no corresponding private key. Therefore the zero
  /// address is a valid input, and no check for non-zero address is performed.
  constructor(
    IBurnMintERC20 token,
    uint8 localTokenDecimals,
    address[] memory allowlist,
    address rmnProxy,
    address router,
    address burnAddress,
    uint256 initialLockedTokens
  ) TokenPool(token, localTokenDecimals, allowlist, rmnProxy, router) {
    i_burnAddress = burnAddress;
    s_lockedTokens = initialLockedTokens;
  }

  /// @notice Mint tokens from the pool to the recipient
  /// @dev The _validateReleaseOrMint check is an essential security check
  function releaseOrMint(
    Pool.ReleaseOrMintInV1 calldata releaseOrMintIn
  ) public virtual override returns (Pool.ReleaseOrMintOutV1 memory) {
    s_lockedTokens += releaseOrMintIn.amount;

    return super.releaseOrMint(releaseOrMintIn);
  }

  /// @inheritdoc BurnMintTokenPoolAbstract
  /// @notice Tokens are burned by sending to an address which is expected to have no corresponding private key,
  /// which makes the tokens unrecoverable without reducing the total supply.
  function _burn(
    uint256 amount
  ) internal virtual override {
    s_lockedTokens -= amount;

    getToken().safeTransfer(i_burnAddress, amount);
  }

  /// @notice returns the address where tokens are sent during a call to lockOrBurn
  /// @return burnAddress the address which receives the tokens.
  function getBurnAddress() public view returns (address burnAddress) {
    return i_burnAddress;
  }

  /// @notice Return the amount of tokens which were minted by this contract and not yet burned.
  /// @return lockedTokens The amount of tokens which were minted by this token pool and not yet burned.
  function getLockedTokens() public view returns (uint256 lockedTokens) {
    return s_lockedTokens;
  }
}
