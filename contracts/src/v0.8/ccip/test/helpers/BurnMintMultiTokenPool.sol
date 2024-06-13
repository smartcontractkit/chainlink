// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";

import {Pool} from "../../libraries/Pool.sol";
import {MultiTokenPool} from "./MultiTokenPool.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract BurnMintMultiTokenPool is MultiTokenPool {
  constructor(
    IERC20[] memory tokens,
    address[] memory allowlist,
    address rmnProxy,
    address router
  ) MultiTokenPool(tokens, allowlist, rmnProxy, router) {}

  /// @notice Burn the token in the pool
  /// @dev The whenNotCursed check is important to ensure that even if a ramp is compromised
  /// we're able to stop token movement via RMN.
  function lockOrBurn(Pool.LockOrBurnInV1 calldata lockOrBurnIn)
    external
    virtual
    override
    returns (Pool.LockOrBurnOutV1 memory)
  {
    _validateLockOrBurn(lockOrBurnIn);

    IBurnMintERC20(lockOrBurnIn.localToken).burn(lockOrBurnIn.amount);

    emit Burned(msg.sender, lockOrBurnIn.amount);

    return Pool.LockOrBurnOutV1({
      destTokenAddress: getRemoteToken(lockOrBurnIn.localToken, lockOrBurnIn.remoteChainSelector),
      destPoolData: ""
    });
  }

  /// @notice Mint tokens from the pool to the recipient
  /// @dev The whenNotCursed check is important to ensure that even if a ramp is compromised
  /// we're able to stop token movement via RMN.
  function releaseOrMint(Pool.ReleaseOrMintInV1 calldata releaseOrMintIn)
    external
    virtual
    override
    returns (Pool.ReleaseOrMintOutV1 memory)
  {
    _validateReleaseOrMint(releaseOrMintIn);

    // Mint to the offRamp, which forwards it to the recipient
    IBurnMintERC20(releaseOrMintIn.localToken).mint(msg.sender, releaseOrMintIn.amount);

    emit Minted(msg.sender, releaseOrMintIn.receiver, releaseOrMintIn.amount);

    return Pool.ReleaseOrMintOutV1({destinationAmount: releaseOrMintIn.amount});
  }
}
