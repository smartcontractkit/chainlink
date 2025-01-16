// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IBurnMintERC20} from "../../shared/token/ERC20/IBurnMintERC20.sol";

import {Pool} from "../libraries/Pool.sol";
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

  event OutstandingTokensSet(uint256 newMintedTokenAmount, uint256 oldMintedTokenAmount);

  error InsufficientOutstandingTokens();

  string public constant override typeAndVersion = "BurnToAddressTokenPool 1.5.1";

  /// @notice The address where tokens are sent during a call to lockOrBurn, functionally burning but without decreasing
  /// total supply. This address is expected to have no ability to recover the tokens sent to it, and will thus be locked forever.
  /// This can be either an EOA without a corresponding private key, or a contract which does not have the ability to transfer the tokens.
  address public immutable i_burnAddress;

  /// @notice Minted Tokens is a safety mechanism to ensure that more tokens cannot be sent out of the bridge
  /// than were originally sent in via CCIP. On incoming messages the value is increased, and on outgoing messages,
  /// the value is decreased. For pools with existing tokens in circulation, the value may not be known at deployment
  /// time, and thus should be set later using the setoutstandingTokens() function.
  uint256 internal s_outstandingTokens;

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

  /// @notice Mint tokens from the pool to the recipient, updating the internal accounting for an outflow of tokens.
  /// @dev If the amount of tokens to be
  function releaseOrMint(
    Pool.ReleaseOrMintInV1 calldata releaseOrMintIn
  ) public virtual override returns (Pool.ReleaseOrMintOutV1 memory) {
    // When minting tokens, the local outstanding supply increases. These tokens will be burned
    // when they are sent back to the pool on an outgoing message.
    s_outstandingTokens += releaseOrMintIn.amount;

    return super.releaseOrMint(releaseOrMintIn);
  }

  /// @inheritdoc BurnMintTokenPoolAbstract
  /// @notice Tokens are burned by sending to an address which can never transfer them,
  /// making the tokens unrecoverable without reducing the total supply.
  function _burn(
    uint256 amount
  ) internal virtual override {
    if (amount > s_outstandingTokens) {
      revert InsufficientOutstandingTokens();
    }

    // When tokens are burned, the amount outstanding decreases. This ensures that more tokens cannot be sent out
    // of the bridge than were originally sent in via CCIP.
    s_outstandingTokens -= amount;

    getToken().safeTransfer(i_burnAddress, amount);
  }

  /// @notice Returns the address where tokens are sent during a call to lockOrBurn
  /// @return burnAddress the address which receives the tokens.
  function getBurnAddress() public view returns (address burnAddress) {
    return i_burnAddress;
  }

  /// @notice Return the amount of tokens which were minted by this contract and not yet burned.
  /// @return outstandingTokens The amount of tokens which were minted by this token pool and not yet burned.
  function getOutstandingTokens() public view returns (uint256 outstandingTokens) {
    return s_outstandingTokens;
  }

  /// @notice Set the amount of tokens which were minted by this contract and not yet burned.
  /// @param amount The new amount of tokens which were minted by this token pool and not yet burned.
  function setOutstandingTokens(
    uint256 amount
  ) external onlyOwner {
    uint256 currentOutstandingTokens = s_outstandingTokens;

    s_outstandingTokens = amount;

    emit OutstandingTokensSet(amount, currentOutstandingTokens);
  }
}
