// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";

import {Ownable2StepMsgSender} from "../../shared/access/Ownable2StepMsgSender.sol";
import {Pool} from "../libraries/Pool.sol";
import {TokenPool} from "./TokenPool.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";

/// @notice A variation on Lock Release token pools where liquidity is shared among some chains, and stored independently
/// for others. Chains which do not share liquidity are known as siloed chains.
/// @dev One token per LockReleaseTokenPool.
contract SiloedLockReleaseTokenPool is TokenPool, ITypeAndVersion {
  using SafeERC20 for IERC20;
  using EnumerableSet for EnumerableSet.UintSet;

  error InsufficientLiquidity();
  error LiquidityNotAccepted();

  event LiquidityTransferred(uint64 remoteChainSelector, address indexed from, uint256 amount);
  event LiquidityAdded(uint64 remoteChainSelector, address indexed provider, uint256 indexed amount);
  event LiquidityRemoved(uint64 remoteChainSelector, address indexed provider, uint256 indexed amount);
  event ChainSiloeDesignationUpdated(uint64 remoteChainSelector, bool isSiloed);

  event RebalancerSet(uint64 indexed remoteChainSelector, address oldRebalancer, address newRebalancer);

  string public constant override typeAndVersion = "SiloedLockReleaseTokenPool 1.5.10-dev";

  mapping(uint64 chainSelector => uint256 lockedBalance) internal s_lockedTokensByChainSelector;
  mapping(uint64 remoteChainSelector => address rebalancer) internal s_rebalancerByChain;
  mapping(uint64 remoteChainSelector => bool shouldSiloFunds) internal s_siloedChainSelectors;

  EnumerableSet.UintSet internal s_siloedChains;

  constructor(
    IERC20 token,
    uint8 localTokenDecimals,
    address[] memory allowlist,
    address rmnProxy,
    address router
  ) TokenPool(token, localTokenDecimals, allowlist, rmnProxy, router) {}

  /// @notice Locks the token in the pool
  /// @dev The _validateLockOrBurn check is an essential security check
  function lockOrBurn(
    Pool.LockOrBurnInV1 calldata lockOrBurnIn
  ) external virtual override returns (Pool.LockOrBurnOutV1 memory) {
    _validateLockOrBurn(lockOrBurnIn);

    // If funds need to be siloed, update internal accounting;
    if (s_siloedChainSelectors[lockOrBurnIn.remoteChainSelector]) {
      s_lockedTokensByChainSelector[lockOrBurnIn.remoteChainSelector] += lockOrBurnIn.amount;
    }

    emit Locked(msg.sender, lockOrBurnIn.amount);

    return Pool.LockOrBurnOutV1({
      destTokenAddress: getRemoteToken(lockOrBurnIn.remoteChainSelector),
      destPoolData: _encodeLocalDecimals()
    });
  }

  /// @notice Release tokens from the pool to the recipient
  /// @dev The _validateReleaseOrMint check is an essential security check
  function releaseOrMint(
    Pool.ReleaseOrMintInV1 calldata releaseOrMintIn
  ) external virtual override returns (Pool.ReleaseOrMintOutV1 memory) {
    _validateReleaseOrMint(releaseOrMintIn);

    // Calculate the local amount
    uint256 localAmount =
      _calculateLocalAmount(releaseOrMintIn.amount, _parseRemoteDecimals(releaseOrMintIn.sourcePoolData));

    // If the message comes from a chain with siloed funds, update internal accounting before releasing tokens.
    if (s_siloedChainSelectors[releaseOrMintIn.remoteChainSelector]) {
      s_lockedTokensByChainSelector[releaseOrMintIn.remoteChainSelector] -= localAmount;
    }

    // Release to the recipient
    getToken().safeTransfer(releaseOrMintIn.receiver, localAmount);

    emit Released(msg.sender, releaseOrMintIn.receiver, localAmount);

    return Pool.ReleaseOrMintOutV1({destinationAmount: localAmount});
  }

  /// @notice Returns whether the tokens locked for a given remote chain should be siloed independently
  /// from all other remote chains.
  /// @param remoteChainSelector the CCIP specific selector for the remote chain being interacted with.
  /// @return isSiloed Whether the funds should be isolated from all the others.
  function chainFundsAreSiloed(
    uint64 remoteChainSelector
  ) external view returns (bool isSiloed) {
    return s_siloedChainSelectors[remoteChainSelector];
  }

  /// @notice Returns the amount of tokens in the token pool that were locked for a specific chain selector.
  /// @param remoteChainSelector the CCIP specific selector for the remote chain being interacted with.
  /// @return lockedTokens The tokens locked into this token pool for the given selector, zero if the remote chain
  /// funds are not siloed.
  function getLockedTokensByChain(
    uint64 remoteChainSelector
  ) external view returns (uint256 lockedTokens) {
    return s_lockedTokensByChainSelector[remoteChainSelector];
  }

  /// @notice Updates designations for chains on whether to mark funds as Siloed or not
  /// @param removes A list of chain selectors to disable Siloeing
  /// @param adds A list of chain selectors to enable LR instead of BM. These chains must not have been migrated
  /// to CCTP yet or the transaction will revert
  function updateChainSelectorMechanisms(uint64[] calldata removes, uint64[] calldata adds) external onlyOwner {
    // If removing, change designation and update tokens as being available for all chains to use.
    for (uint256 i = 0; i < removes.length; ++i) {
      delete s_siloedChainSelectors[removes[i]];
      delete s_lockedTokensByChainSelector[removes[i]];

      emit ChainSiloeDesignationUpdated(removes[i], false);
    }

    for (uint256 i = 0; i < adds.length; ++i) {
      s_siloedChainSelectors[adds[i]] = true;
      emit ChainSiloeDesignationUpdated(adds[i], true);
    }
  }

  /// @notice Gets the rebalancer able to provide liquidity for a remote chain selector
  /// @param remoteChainSelector The CCIP specific selector for the remote chain being interacted with.
  /// @return The current liquidity manager, owner if the chain's funds are not siloed.
  function getRebalancerByChain(
    uint64 remoteChainSelector
  ) external view returns (address) {
    if (s_siloedChainSelectors[remoteChainSelector]) return s_rebalancerByChain[remoteChainSelector];
    else return owner();
  }

  /// @notice Sets the LiquidityManager address.
  /// @dev Only callable by the owner.
  /// @param remoteChainSelector the remote chain to set.
  /// @param rebalancer the address allowed to add liquidity for the given siloed chain.
  function setRebalancer(uint64 remoteChainSelector, address rebalancer) external onlyOwner {
    address oldRebalancer = s_rebalancerByChain[remoteChainSelector];

    s_rebalancerByChain[remoteChainSelector] = rebalancer;

    emit RebalancerSet(remoteChainSelector, rebalancer, oldRebalancer);
  }

  /// @notice Adds liquidity to the pool. The tokens should be approved first.
  /// @param remoteChainSelector the remote chain to set.
  /// @param amount The amount of liquidity to provide.
  function provideLiquidity(uint64 remoteChainSelector, uint256 amount) external {
    // Save gas by performing the enumerable set query once for both authorization and internal accounting
    if (s_siloedChainSelectors[remoteChainSelector]) {
      if (msg.sender != s_rebalancerByChain[remoteChainSelector]) revert Unauthorized(msg.sender);

      s_lockedTokensByChainSelector[remoteChainSelector] += amount;
    } else if (msg.sender != owner()) {
      revert Unauthorized(msg.sender);
    }

    i_token.safeTransferFrom(msg.sender, address(this), amount);
    emit LiquidityAdded(remoteChainSelector, msg.sender, amount);
  }

  /// @notice Removed liquidity to the pool. The tokens will be sent to msg.sender.
  /// @param remoteChainSelector the remote chain to set. If the chain is not siloed, then no accounting will be updated,
  /// which can be considered the liquidity for all non-siloed chains sharing liquidity.
  /// @param amount The amount of liquidity to remove.
  function withdrawLiquidity(uint64 remoteChainSelector, uint256 amount) external onlyOwner {
    // If funds are siloed by chain, prevent more than has been locked from being removed from the token pool.
    if (s_siloedChainSelectors[remoteChainSelector]) {
      s_lockedTokensByChainSelector[remoteChainSelector] -= amount;
    }

    i_token.safeTransfer(msg.sender, amount);
    emit LiquidityRemoved(remoteChainSelector, msg.sender, amount);
  }

  /// @notice This function can be used to transfer liquidity from an older version of the pool to this pool. To do so
  /// this pool will have to be set as the rebalancer in the older version of the pool. This allows it to transfer the
  /// funds in the old pool to the new pool.
  /// @dev When upgrading a LockRelease pool, this function can be called at the same time as the pool is changed in the
  /// TokenAdminRegistry. This allows for a smooth transition of both liquidity and transactions to the new pool.
  /// Alternatively, when no multicall is available, a portion of the funds can be transferred to the new pool before
  /// changing which pool CCIP uses, to ensure both pools can operate. Then the pool should be changed in the
  /// TokenAdminRegistry, which will activate the new pool. All new transactions will use the new pool and its
  /// liquidity. Finally, the remaining liquidity can be transferred to the new pool using this function one more time.
  /// @param remoteChainSelector the remote chain to set. If the chain is not siloed, then no accounting will be updated,
  /// which can be considered the liquidity for all non-siloed chains sharing liquidity.
  /// @param from The address of the old pool.
  /// @param amount The amount of liquidity to transfer.
  function transferLiquidity(uint64 remoteChainSelector, address from, uint256 amount) external onlyOwner {
    // If The ownership has already been accepted, do not attempt to accept again as it would fail.
    if (Ownable2StepMsgSender(from).owner() != address(this)) Ownable2StepMsgSender(from).acceptOwnership();

    SiloedLockReleaseTokenPool(from).withdrawLiquidity(remoteChainSelector, amount);

    // Since both siloed and non-siloed token liquidity can be transferred, allow transfers from both, but only
    // update internal accounting for siloed chains.
    if (s_siloedChainSelectors[remoteChainSelector]) {
      s_lockedTokensByChainSelector[remoteChainSelector] += amount;
    }

    emit LiquidityTransferred(remoteChainSelector, from, amount);
  }
}
