// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Pool} from "../../libraries/Pool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract TokenPoolHelper is TokenPool {
  constructor(
    IERC20 token,
    address[] memory allowlist,
    address rmnProxy,
    address router
  ) TokenPool(token, allowlist, rmnProxy, router) {}

  function lockOrBurn(
    Pool.LockOrBurnInV1 calldata lockOrBurnIn
  ) external view override returns (Pool.LockOrBurnOutV1 memory) {
    return Pool.LockOrBurnOutV1({destTokenAddress: getRemoteToken(lockOrBurnIn.remoteChainSelector), destPoolData: ""});
  }

  function releaseOrMint(
    Pool.ReleaseOrMintInV1 calldata releaseOrMintIn
  ) external pure override returns (Pool.ReleaseOrMintOutV1 memory) {
    return Pool.ReleaseOrMintOutV1({destinationAmount: releaseOrMintIn.amount});
  }

  function onlyOnRampModifier(uint64 remoteChainSelector) external view {
    _onlyOnRamp(remoteChainSelector);
  }

  function onlyOffRampModifier(uint64 remoteChainSelector) external view {
    _onlyOffRamp(remoteChainSelector);
  }
}
