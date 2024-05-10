// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IBurnMintERC20} from "../../shared/token/ERC20/IBurnMintERC20.sol";

import {Pool} from "../libraries/Pool.sol";
import {LegacyPoolWrapper} from "./LegacyPoolWrapper.sol";

contract BurnMintTokenPoolAndProxy is ITypeAndVersion, LegacyPoolWrapper {
  string public constant override typeAndVersion = "BurnMintTokenPoolAndProxy 1.5.0-dev";

  constructor(
    IBurnMintERC20 token,
    address[] memory allowlist,
    address rmnProxy,
    address router
  ) LegacyPoolWrapper(token, allowlist, rmnProxy, router) {}

  /// @notice Burn the token in the pool
  /// @dev The whenNotCursed check is important to ensure that even if a ramp is compromised
  /// we're able to stop token movement via RMN.
  function lockOrBurn(Pool.LockOrBurnInV1 calldata lockOrBurnIn)
    external
    virtual
    override
    whenNotCursed(lockOrBurnIn.remoteChainSelector)
    returns (Pool.LockOrBurnOutV1 memory)
  {
    _checkAllowList(lockOrBurnIn.originalSender);
    _onlyOnRamp(lockOrBurnIn.remoteChainSelector);
    _consumeOutboundRateLimit(lockOrBurnIn.remoteChainSelector, lockOrBurnIn.amount);

    if (!_hasLegacyPool()) {
      IBurnMintERC20(address(i_token)).burn(lockOrBurnIn.amount);
    } else {
      _lockOrBurnLegacy(lockOrBurnIn);
    }

    emit Burned(msg.sender, lockOrBurnIn.amount);

    return Pool.LockOrBurnOutV1({destPoolAddress: getRemotePool(lockOrBurnIn.remoteChainSelector), destPoolData: ""});
  }

  /// @notice Mint tokens from the pool to the recipient
  /// @dev The whenNotCursed check is important to ensure that even if a ramp is compromised
  /// we're able to stop token movement via RMN.
  function releaseOrMint(Pool.ReleaseOrMintInV1 calldata releaseOrMintIn)
    external
    virtual
    override
    whenNotCursed(releaseOrMintIn.remoteChainSelector)
    returns (Pool.ReleaseOrMintOutV1 memory)
  {
    _onlyOffRamp(releaseOrMintIn.remoteChainSelector);
    _validateSourceCaller(releaseOrMintIn.remoteChainSelector, releaseOrMintIn.sourcePoolAddress);
    _consumeInboundRateLimit(releaseOrMintIn.remoteChainSelector, releaseOrMintIn.amount);

    if (!_hasLegacyPool()) {
      IBurnMintERC20(address(i_token)).mint(releaseOrMintIn.receiver, releaseOrMintIn.amount);
    } else {
      _releaseOrMintLegacy(releaseOrMintIn);
    }

    emit Minted(msg.sender, releaseOrMintIn.receiver, releaseOrMintIn.amount);

    return Pool.ReleaseOrMintOutV1({localToken: address(i_token), destinationAmount: releaseOrMintIn.amount});
  }
}
