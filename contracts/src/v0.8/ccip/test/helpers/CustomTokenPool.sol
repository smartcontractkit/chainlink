// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Pool} from "../../libraries/Pool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract CustomTokenPool is TokenPool {
  event SynthBurned(uint256 amount);
  event SynthMinted(uint256 amount);

  constructor(IERC20 token, address rmnProxy, address router) TokenPool(token, new address[](0), rmnProxy, router) {}

  /// @notice Locks the token in the pool
  function lockOrBurn(Pool.LockOrBurnInV1 calldata lockOrBurnIn)
    external
    virtual
    override
    whenNotCursed(lockOrBurnIn.remoteChainSelector)
    returns (Pool.LockOrBurnOutV1 memory)
  {
    _onlyOnRamp(lockOrBurnIn.remoteChainSelector);
    emit SynthBurned(lockOrBurnIn.amount);
    return Pool.LockOrBurnOutV1({destPoolAddress: getRemotePool(lockOrBurnIn.remoteChainSelector), destPoolData: ""});
  }

  /// @notice Release tokens from the pool to the recipient
  function releaseOrMint(Pool.ReleaseOrMintInV1 calldata releaseOrMintIn)
    external
    override
    whenNotCursed(releaseOrMintIn.remoteChainSelector)
    returns (Pool.ReleaseOrMintOutV1 memory)
  {
    _onlyOffRamp(releaseOrMintIn.remoteChainSelector);
    emit SynthMinted(releaseOrMintIn.amount);
    return Pool.ReleaseOrMintOutV1({localToken: address(i_token), destinationAmount: releaseOrMintIn.amount});
  }
}
