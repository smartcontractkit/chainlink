// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IBurnMintERC20} from "../../shared/token/ERC20/IBurnMintERC20.sol";

import {BurnMintTokenPoolAbstract} from "./BurnMintTokenPoolAbstract.sol";
import {TokenPool} from "./TokenPool.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @notice This pool mints and burns a 3rd-party token by sending tokens to an address which is unrecoverable.
/// @dev The pool is designed to have an immutable burn address. If the tokens at the burn address become recoverable,
/// for example, a quantum computer calculating a private key for the zero address, the pool will need to be replaced
/// with a new pool with a different burn address.
contract BurnToAddressMintTokenPool is BurnMintTokenPoolAbstract, ITypeAndVersion {
  using SafeERC20 for IERC20;

  string public constant override typeAndVersion = "BurnToAddressTokenPool 1.5.1";

  /// @notice The address where tokens are sent during a call to lockOrBurn, functionally burning but without decreasing
  /// total supply. This address is expected to have no ability to recover the tokens sent to it, and will thus be locked forever.
  /// This can be either an EOA without a corresponding private key, or a contract which does not have the ability to transfer the tokens.
  address public immutable i_burnAddress;

  /// @dev Since burnAddress is expected to make the tokens unrecoverable, no check for the zero address needs to be
  /// performed, as it is a valid input.
  constructor(
    IBurnMintERC20 token,
    uint8 localTokenDecimals,
    address[] memory allowlist,
    address rmnProxy,
    address router,
    address burnAddress
  ) TokenPool(token, localTokenDecimals, allowlist, rmnProxy, router) {
    i_burnAddress = burnAddress;
  }

  /// @inheritdoc BurnMintTokenPoolAbstract
  /// @notice Tokens are burned by sending to an address which can never transfer them,
  /// making the tokens unrecoverable without reducing the total supply.
  function _burn(
    uint256 amount
  ) internal virtual override {
    getToken().safeTransfer(i_burnAddress, amount);
  }

  /// @notice Returns the address where tokens are sent during a call to lockOrBurn
  /// @return burnAddress the address which receives the tokens.
  function getBurnAddress() public view returns (address burnAddress) {
    return i_burnAddress;
  }
}
