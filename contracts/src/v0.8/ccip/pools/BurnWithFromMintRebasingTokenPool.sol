// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IBurnMintERC20} from "../../shared/token/ERC20/IBurnMintERC20.sol";

import {Pool} from "../libraries/Pool.sol";
import {BurnWithFromMintTokenPool} from "./BurnWithFromMintTokenPool.sol";

/// @notice This pool mints and burns a 3rd-party token.
/// @dev This contract is a variant of BurnMintTokenPool that uses `burn(from, amount)`.
/// @dev This contract supports minting tokens that do not mint the exact amount they are asked to mint. This can be
/// used for rebasing tokens. NOTE: for true rebasing support, the lockOrBurn method must also be updated to support
/// relaying the correct amount.
contract BurnWithFromMintRebasingTokenPool is BurnWithFromMintTokenPool {
  error NegativeMintAmount(uint256 amountBurned);

  string public constant override typeAndVersion = "BurnWithFromMintRebasingTokenPool 1.5.0";

  constructor(
    IBurnMintERC20 token,
    address[] memory allowlist,
    address rmnProxy,
    address router
  ) BurnWithFromMintTokenPool(token, allowlist, rmnProxy, router) {}

  /// @notice Mint tokens from the pool to the recipient
  /// @dev The _validateReleaseOrMint check is an essential security check
  function releaseOrMint(
    Pool.ReleaseOrMintInV1 calldata releaseOrMintIn
  ) external virtual override returns (Pool.ReleaseOrMintOutV1 memory) {
    _validateReleaseOrMint(releaseOrMintIn);

    uint256 balancePre = IBurnMintERC20(address(i_token)).balanceOf(releaseOrMintIn.receiver);

    // Mint to the receiver
    IBurnMintERC20(address(i_token)).mint(releaseOrMintIn.receiver, releaseOrMintIn.amount);

    uint256 balancePost = IBurnMintERC20(address(i_token)).balanceOf(releaseOrMintIn.receiver);

    // Mint should not reduce the number of tokens in the receiver, if it does it will revert the call.
    if (balancePost < balancePre) {
      revert NegativeMintAmount(balancePre - balancePost);
    }

    emit Minted(msg.sender, releaseOrMintIn.receiver, balancePost - balancePre);

    return Pool.ReleaseOrMintOutV1({destinationAmount: balancePost - balancePre});
  }
}
