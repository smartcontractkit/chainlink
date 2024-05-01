// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IBurnMintERC20} from "../../shared/token/ERC20/IBurnMintERC20.sol";

import {Pool} from "../libraries/Pool.sol";
import {TokenPool} from "./TokenPool.sol";

abstract contract BurnMintTokenPoolAbstract is TokenPool {
  /// @notice Contains the specific burn call for a pool.
  /// @dev overriding this method allows us to create pools with different burn signatures
  /// without duplicating the underlying logic.
  function _burn(uint256 amount) internal virtual;

  /// @notice Burn the token in the pool
  /// @dev The whenHealthy check is important to ensure that even if a ramp is compromised
  /// we're able to stop token movement via ARM.
  function lockOrBurn(Pool.LockOrBurnInV1 calldata lockOrBurnIn)
    external
    virtual
    override
    whenHealthy
    returns (Pool.LockOrBurnOutV1 memory)
  {
    _checkAllowList(lockOrBurnIn.originalSender);
    _onlyOnRamp(lockOrBurnIn.remoteChainSelector);
    _consumeOutboundRateLimit(lockOrBurnIn.remoteChainSelector, lockOrBurnIn.amount);

    _burn(lockOrBurnIn.amount);

    emit Burned(msg.sender, lockOrBurnIn.amount);

    return Pool.LockOrBurnOutV1({destPoolAddress: getRemotePool(lockOrBurnIn.remoteChainSelector), destPoolData: ""});
  }

  /// @notice Mint tokens from the pool to the recipient
  /// @dev The whenHealthy check is important to ensure that even if a ramp is compromised
  /// we're able to stop token movement via ARM.
  function releaseOrMint(Pool.ReleaseOrMintInV1 calldata releaseOrMintIn)
    external
    virtual
    override
    whenHealthy
    returns (Pool.ReleaseOrMintOutV1 memory)
  {
    _onlyOffRamp(releaseOrMintIn.remoteChainSelector);
    _validateSourceCaller(releaseOrMintIn.remoteChainSelector, releaseOrMintIn.sourcePoolAddress);
    _consumeInboundRateLimit(releaseOrMintIn.remoteChainSelector, releaseOrMintIn.amount);

    IBurnMintERC20(address(i_token)).mint(releaseOrMintIn.receiver, releaseOrMintIn.amount);

    emit Minted(msg.sender, releaseOrMintIn.receiver, releaseOrMintIn.amount);

    return Pool.ReleaseOrMintOutV1({localToken: address(i_token), destinationAmount: releaseOrMintIn.amount});
  }
}
