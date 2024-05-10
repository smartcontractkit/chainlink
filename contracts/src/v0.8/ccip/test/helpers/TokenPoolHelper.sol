// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {Pool} from "../../libraries/Pool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract TokenPoolHelper is TokenPool {
  event LockOrBurn(uint256 amount);
  event ReleaseOrMint(address indexed recipient, uint256 amount);
  event AssertionPassed();

  constructor(
    IERC20 token,
    address[] memory allowlist,
    address rmnProxy,
    address router
  ) TokenPool(token, allowlist, rmnProxy, router) {}

  function lockOrBurn(Pool.LockOrBurnInV1 calldata lockOrBurnIn)
    external
    override
    returns (Pool.LockOrBurnOutV1 memory)
  {
    emit LockOrBurn(lockOrBurnIn.amount);
    return Pool.LockOrBurnOutV1({destPoolAddress: getRemotePool(lockOrBurnIn.remoteChainSelector), destPoolData: ""});
  }

  function releaseOrMint(Pool.ReleaseOrMintInV1 calldata releaseOrMintIn)
    external
    override
    returns (Pool.ReleaseOrMintOutV1 memory)
  {
    emit ReleaseOrMint(releaseOrMintIn.receiver, releaseOrMintIn.amount);
    return Pool.ReleaseOrMintOutV1({localToken: address(i_token), destinationAmount: releaseOrMintIn.amount});
  }

  function onlyOnRampModifier(uint64 remoteChainSelector) external view {
    _onlyOnRamp(remoteChainSelector);
  }

  function onlyOffRampModifier(uint64 remoteChainSelector) external view {
    _onlyOffRamp(remoteChainSelector);
  }
}
