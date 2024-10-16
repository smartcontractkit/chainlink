// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";

import {Pool} from "../../libraries/Pool.sol";
import {BurnMintTokenPool} from "../../pools/BurnMintTokenPool.sol";

contract MaybeRevertingBurnMintTokenPool is BurnMintTokenPool {
  bytes public s_revertReason = "";
  bytes public s_sourceTokenData = "";
  uint256 public s_releaseOrMintMultiplier = 1;

  constructor(
    IBurnMintERC20 token,
    address[] memory allowlist,
    address rmnProxy,
    address router
  ) BurnMintTokenPool(token, allowlist, rmnProxy, router) {}

  function setShouldRevert(
    bytes calldata revertReason
  ) external {
    s_revertReason = revertReason;
  }

  function setSourceTokenData(
    bytes calldata sourceTokenData
  ) external {
    s_sourceTokenData = sourceTokenData;
  }

  function setReleaseOrMintMultiplier(
    uint256 multiplier
  ) external {
    s_releaseOrMintMultiplier = multiplier;
  }

  function lockOrBurn(
    Pool.LockOrBurnInV1 calldata lockOrBurnIn
  ) external virtual override returns (Pool.LockOrBurnOutV1 memory) {
    _validateLockOrBurn(lockOrBurnIn);

    bytes memory revertReason = s_revertReason;
    if (revertReason.length != 0) {
      assembly {
        revert(add(32, revertReason), mload(revertReason))
      }
    }

    IBurnMintERC20(address(i_token)).burn(lockOrBurnIn.amount);
    emit Burned(msg.sender, lockOrBurnIn.amount);
    return Pool.LockOrBurnOutV1({
      destTokenAddress: getRemoteToken(lockOrBurnIn.remoteChainSelector),
      destPoolData: s_sourceTokenData
    });
  }

  /// @notice Reverts depending on the value of `s_revertReason`
  function releaseOrMint(
    Pool.ReleaseOrMintInV1 calldata releaseOrMintIn
  ) external virtual override returns (Pool.ReleaseOrMintOutV1 memory) {
    _validateReleaseOrMint(releaseOrMintIn);

    bytes memory revertReason = s_revertReason;
    if (revertReason.length != 0) {
      assembly {
        revert(add(32, revertReason), mload(revertReason))
      }
    }
    uint256 amount = releaseOrMintIn.amount * s_releaseOrMintMultiplier;
    IBurnMintERC20(address(i_token)).mint(releaseOrMintIn.receiver, amount);

    emit Minted(msg.sender, releaseOrMintIn.receiver, amount);
    return Pool.ReleaseOrMintOutV1({destinationAmount: amount});
  }
}
