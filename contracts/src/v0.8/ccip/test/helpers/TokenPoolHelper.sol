// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {Pool} from "../../libraries/Pool.sol";
import "../../pools/TokenPool.sol";

contract TokenPoolHelper is TokenPool {
  event LockOrBurn(uint256 amount);
  event ReleaseOrMint(address indexed recipient, uint256 amount);
  event AssertionPassed();

  constructor(
    IERC20 token,
    address[] memory allowlist,
    address armProxy,
    address router
  ) TokenPool(token, allowlist, armProxy, router) {}

  function lockOrBurn(
    address,
    bytes calldata,
    uint256 amount,
    uint64 remoteChainSelector,
    bytes calldata
  ) external override returns (bytes memory) {
    emit LockOrBurn(amount);
    return Pool._generatePoolReturnDataV1(getRemotePool(remoteChainSelector), "");
  }

  function releaseOrMint(
    bytes memory,
    address receiver,
    uint256 amount,
    uint64,
    IPool.SourceTokenData memory,
    bytes memory
  ) external override returns (address) {
    emit ReleaseOrMint(receiver, amount);
    return address(i_token);
  }

  function onlyOnRampModifier(uint64 remoteChainSelector) external onlyOnRamp(remoteChainSelector) {}

  function onlyOffRampModifier(uint64 remoteChainSelector) external onlyOffRamp(remoteChainSelector) {}
}
