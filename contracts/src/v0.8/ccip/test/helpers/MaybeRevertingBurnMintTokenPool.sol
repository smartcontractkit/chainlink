// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";

import {Pool} from "../../libraries/Pool.sol";
import {BurnMintTokenPool} from "../../pools/BurnMintTokenPool.sol";

contract MaybeRevertingBurnMintTokenPool is BurnMintTokenPool {
  bytes public s_revertReason = "";
  bytes public s_sourceTokenData = "";

  constructor(
    IBurnMintERC20 token,
    address[] memory allowlist,
    address rmnProxy,
    address router
  ) BurnMintTokenPool(token, allowlist, rmnProxy, router) {}

  function setShouldRevert(bytes calldata revertReason) external {
    s_revertReason = revertReason;
  }

  function setSourceTokenData(bytes calldata sourceTokenData) external {
    s_sourceTokenData = sourceTokenData;
  }

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

    bytes memory revertReason = s_revertReason;
    if (revertReason.length != 0) {
      assembly {
        revert(add(32, revertReason), mload(revertReason))
      }
    }

    IBurnMintERC20(address(i_token)).burn(lockOrBurnIn.amount);
    emit Burned(msg.sender, lockOrBurnIn.amount);
    return Pool.LockOrBurnOutV1({
      destPoolAddress: getRemotePool(lockOrBurnIn.remoteChainSelector),
      destPoolData: s_sourceTokenData
    });
  }

  /// @notice Reverts depending on the value of `s_revertReason`
  function releaseOrMint(Pool.ReleaseOrMintInV1 calldata releaseOrMintIn)
    external
    virtual
    override
    whenNotCursed(releaseOrMintIn.remoteChainSelector)
    returns (Pool.ReleaseOrMintOutV1 memory)
  {
    _onlyOffRamp(releaseOrMintIn.remoteChainSelector);
    _validateSourceCaller(releaseOrMintIn.remoteChainSelector, releaseOrMintIn.sourcePoolAddress);
    bytes memory revertReason = s_revertReason;
    if (revertReason.length != 0) {
      assembly {
        revert(add(32, revertReason), mload(revertReason))
      }
    }
    _consumeInboundRateLimit(releaseOrMintIn.remoteChainSelector, releaseOrMintIn.amount);
    IBurnMintERC20(address(i_token)).mint(releaseOrMintIn.receiver, releaseOrMintIn.amount);
    emit Minted(msg.sender, releaseOrMintIn.receiver, releaseOrMintIn.amount);
    return Pool.ReleaseOrMintOutV1({localToken: address(i_token), destinationAmount: releaseOrMintIn.amount});
  }
}
