// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IBridgeAdapter} from "./interfaces/IBridge.sol";
import {IRebalancer} from "./interfaces/IRebalancer.sol";
import {ILiquidityContainer} from "./interfaces/ILiquidityContainer.sol";

import {OCR3Base} from "./ocr/OCR3Base.sol";

import {IERC20} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @notice Rebalancer for a single token over multiple chains.
/// @dev This contract is designed to be used with the LockReleaseTokenPool contract but
/// isn't constraint to it. It can be used with any contract that implements the ILiquidityContainer
/// interface.
/// @dev The OCR3 DON should only be able to transfer funds to other pre-approved contracts
/// on other chains. Under no circumstances should it be able to transfer funds to arbitrary
/// addresses. The owner is therefore in full control of the funds in this contract, not the DON.
/// This is a security feature. The worst that can happen is that the DON can lock up funds in
/// bridges, but it can't steal them.
/// @dev References to local mean logic on the same chain as this contract is deployed on.
/// References to remote mean logic on other chains.
contract Rebalancer is IRebalancer, OCR3Base {
  using SafeERC20 for IERC20;

  error ZeroAddress();
  error InvalidRemoteChain(uint64 chainSelector);
  error ZeroChainSelector();
  error InsufficientLiquidity(uint256 requested, uint256 available);

  event LiquidityTransferred(
    uint64 indexed ocrSeqNum,
    uint64 indexed fromChainSelector,
    uint64 indexed toChainSelector,
    address to,
    uint256 amount,
    bytes bridgeSpecificData,
    bytes bridgeReturnData
  );
  event LiquidityAdded(address indexed provider, uint256 indexed amount);
  event LiquidityRemoved(address indexed remover, uint256 indexed amount);

  event FinalizationFailed(
    uint64 indexed ocrSeqNum,
    uint64 indexed remoteChainSelector,
    bytes bridgeSpecificData,
    bytes reason
  );

  struct CrossChainRebalancer {
    address remoteRebalancer;
    IBridgeAdapter localBridge;
    address remoteToken;
    bool enabled;
  }

  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  string public constant override typeAndVersion = "Rebalancer 1.0.0-dev";

  /// @notice The token that this pool manages liquidity for.
  IERC20 public immutable i_localToken;

  /// @notice The chain selector belonging to the chain this pool is deployed on.
  uint64 internal immutable i_localChainSelector;

  /// @notice Mapping of chain selector to liquidity container on other chains
  mapping(uint64 chainSelector => CrossChainRebalancer) private s_crossChainRebalancer;

  uint64[] private s_supportedDestChains;

  /// @notice The liquidity container on the local chain
  /// @dev In the case of CCIP, this would be the token pool.
  ILiquidityContainer private s_localLiquidityContainer;

  constructor(IERC20 token, uint64 localChainSelector, ILiquidityContainer localLiquidityContainer) OCR3Base() {
    if (localChainSelector == 0) {
      revert ZeroChainSelector();
    }

    if (address(token) == address(0)) {
      revert ZeroAddress();
    }
    i_localToken = token;
    i_localChainSelector = localChainSelector;
    s_localLiquidityContainer = localLiquidityContainer;
  }

  receive() external payable {}

  // ================================================================
  // │                    Liquidity management                      │
  // ================================================================

  /// @inheritdoc IRebalancer
  function getLiquidity() public view returns (uint256 currentLiquidity) {
    return i_localToken.balanceOf(address(s_localLiquidityContainer));
  }

  /// @notice Adds liquidity to the multi-chain system.
  /// @dev Anyone can call this function, but anyone other than the owner should regard
  /// adding liquidity as a donation to the system, as there is no way to get it out.
  /// This function is open to anyone to be able to quickly add funds to the system
  /// without having to go through potentially complicated multisig schemes to do it from
  /// the owner address.
  function addLiquidity(uint256 amount) external {
    i_localToken.safeTransferFrom(msg.sender, address(this), amount);

    // Make sure this is tether compatible, as they have strange approval requirements
    // Should be good since all approvals are always immediately used.
    i_localToken.approve(address(s_localLiquidityContainer), amount);
    s_localLiquidityContainer.provideLiquidity(amount);

    emit LiquidityAdded(msg.sender, amount);
  }

  /// @notice Removes liquidity from the system and sends it to the caller, so the owner.
  /// @dev Only the owner can call this function.
  function removeLiquidity(uint256 amount) external onlyOwner {
    uint256 currentBalance = i_localToken.balanceOf(address(s_localLiquidityContainer));
    if (currentBalance < amount) {
      revert InsufficientLiquidity(amount, currentBalance);
    }

    s_localLiquidityContainer.withdrawLiquidity(amount);
    i_localToken.safeTransfer(msg.sender, amount);

    emit LiquidityRemoved(msg.sender, amount);
  }

  /// @notice Transfers liquidity to another chain.
  /// @dev This function is a public version of the internal _rebalanceLiquidity function.
  /// to allow the owner to also initiate a rebalancing when needed.
  function rebalanceLiquidity(
    uint64 chainSelector,
    uint256 amount,
    uint256 nativeBridgeFee,
    bytes calldata bridgeSpecificPayload
  ) external onlyOwner {
    _rebalanceLiquidity(chainSelector, amount, nativeBridgeFee, type(uint64).max, bridgeSpecificPayload);
  }

  /// @notice Finalizes liquidity from another chain.
  /// @dev This function is a public version of the internal _receiveLiquidity function.
  /// to allow the owner to also initiate a finalization when needed.
  function receiveLiquidity(
    uint64 remoteChainSelector,
    uint256 amount,
    bytes calldata bridgeSpecificPayload
  ) external onlyOwner {
    _receiveLiquidity(remoteChainSelector, amount, bridgeSpecificPayload, type(uint64).max);
  }

  /// @notice Transfers liquidity to another chain.
  /// @dev Called by both the owner and the DON.
  function _rebalanceLiquidity(
    uint64 chainSelector,
    uint256 tokenAmount,
    uint256 nativeBridgeFee,
    uint64 ocrSeqNum,
    bytes memory bridgeSpecificPayload
  ) internal {
    uint256 currentBalance = getLiquidity();
    if (currentBalance < tokenAmount) {
      revert InsufficientLiquidity(tokenAmount, currentBalance);
    }

    CrossChainRebalancer memory remoteLiqManager = s_crossChainRebalancer[chainSelector];

    if (!remoteLiqManager.enabled) {
      revert InvalidRemoteChain(chainSelector);
    }

    // XXX: Could be optimized by withdrawing once and then sending to all destinations
    s_localLiquidityContainer.withdrawLiquidity(tokenAmount);
    i_localToken.approve(address(remoteLiqManager.localBridge), tokenAmount);

    bytes memory bridgeReturnData = remoteLiqManager.localBridge.sendERC20{value: nativeBridgeFee}(
      address(i_localToken),
      remoteLiqManager.remoteToken,
      remoteLiqManager.remoteRebalancer,
      tokenAmount,
      bridgeSpecificPayload
    );

    emit LiquidityTransferred(
      ocrSeqNum,
      i_localChainSelector,
      chainSelector,
      remoteLiqManager.remoteRebalancer,
      tokenAmount,
      bridgeSpecificPayload,
      bridgeReturnData
    );
  }

  function _receiveLiquidity(
    uint64 remoteChainSelector,
    uint256 amount,
    bytes memory bridgeSpecificPayload,
    uint64 ocrSeqNum
  ) internal {
    // check if the remote chain is supported
    CrossChainRebalancer memory remoteRebalancer = s_crossChainRebalancer[remoteChainSelector];
    if (!remoteRebalancer.enabled) {
      revert InvalidRemoteChain(remoteChainSelector);
    }

    // finalize the withdrawal through the bridge adapter
    try
      remoteRebalancer.localBridge.finalizeWithdrawERC20(
        remoteRebalancer.remoteRebalancer, // remoteSender: the remote rebalancer
        address(this), // localReceiver: us
        bridgeSpecificPayload
      )
    {
      // successfully finalized the withdrawal
    } catch (bytes memory lowLevelData) {
      // failed to finalize the withdrawal.
      // this could mean that the withdrawal was already finalized
      // or that the withdrawal failed.
      // we assume the former and continue
      emit FinalizationFailed(ocrSeqNum, remoteChainSelector, bridgeSpecificPayload, lowLevelData);
    }

    // inject liquidity into the liquidity container
    // approve and liquidity container should transferFrom
    i_localToken.safeIncreaseAllowance(address(s_localLiquidityContainer), amount);
    s_localLiquidityContainer.provideLiquidity(amount);

    emit LiquidityTransferred(
      ocrSeqNum,
      remoteChainSelector,
      i_localChainSelector,
      address(this),
      amount,
      bridgeSpecificPayload,
      bytes("") // no bridge return data when receiving
    );
  }

  function _report(bytes calldata report, uint64 ocrSeqNum) internal override {
    IRebalancer.LiquidityInstructions memory instructions = abi.decode(report, (IRebalancer.LiquidityInstructions));

    uint256 sendInstructions = instructions.sendLiquidityParams.length;
    for (uint256 i = 0; i < sendInstructions; ++i) {
      _rebalanceLiquidity(
        instructions.sendLiquidityParams[i].remoteChainSelector,
        instructions.sendLiquidityParams[i].amount,
        instructions.sendLiquidityParams[i].nativeBridgeFee,
        ocrSeqNum,
        instructions.sendLiquidityParams[i].bridgeData
      );
    }

    uint256 receiveInstructions = instructions.receiveLiquidityParams.length;
    for (uint256 i = 0; i < receiveInstructions; ++i) {
      _receiveLiquidity(
        instructions.receiveLiquidityParams[i].remoteChainSelector,
        instructions.receiveLiquidityParams[i].amount,
        instructions.receiveLiquidityParams[i].bridgeData,
        ocrSeqNum
      );
    }

    // todo emit?
  }

  // ================================================================
  // │                           Config                             │
  // ================================================================

  function getSupportedDestChains() external view returns (uint64[] memory) {
    return s_supportedDestChains;
  }

  /// @notice Gets the cross chain liquidity manager
  function getCrossChainRebalancer(uint64 chainSelector) external view returns (CrossChainRebalancer memory) {
    return s_crossChainRebalancer[chainSelector];
  }

  /// @notice Gets all cross chain liquidity managers
  /// @dev We don't care too much about gas since this function is intended for offchain usage.
  function getAllCrossChainRebalancers() external view returns (CrossChainRebalancerArgs[] memory) {
    uint256 numChains = s_supportedDestChains.length;
    CrossChainRebalancerArgs[] memory managers = new CrossChainRebalancerArgs[](numChains);
    for (uint256 i = 0; i < numChains; ++i) {
      uint64 chainSelector = s_supportedDestChains[i];
      CrossChainRebalancer memory currentManager = s_crossChainRebalancer[chainSelector];
      managers[i] = CrossChainRebalancerArgs({
        remoteRebalancer: currentManager.remoteRebalancer,
        localBridge: currentManager.localBridge,
        remoteToken: currentManager.remoteToken,
        remoteChainSelector: chainSelector,
        enabled: currentManager.enabled
      });
    }

    return managers;
  }

  /// @notice Sets a list of cross chain liquidity managers.
  /// @dev Will update the list of supported dest chains if the chain is new.
  function setCrossChainRebalancer(CrossChainRebalancerArgs[] calldata crossChainRebalancers) external onlyOwner {
    for (uint256 i = 0; i < crossChainRebalancers.length; ++i) {
      setCrossChainRebalancer(crossChainRebalancers[i]);
    }
  }

  /// @notice Sets a single cross chain liquidity manager.
  /// @dev Will update the list of supported dest chains if the chain is new.
  function setCrossChainRebalancer(CrossChainRebalancerArgs calldata crossChainLiqManager) public onlyOwner {
    if (crossChainLiqManager.remoteChainSelector == 0) {
      revert ZeroChainSelector();
    }

    if (
      crossChainLiqManager.remoteRebalancer == address(0) ||
      address(crossChainLiqManager.localBridge) == address(0) ||
      crossChainLiqManager.remoteToken == address(0)
    ) {
      revert ZeroAddress();
    }

    // If the destination chain is new, add it to the list of supported chains
    if (s_crossChainRebalancer[crossChainLiqManager.remoteChainSelector].remoteToken == address(0)) {
      s_supportedDestChains.push(crossChainLiqManager.remoteChainSelector);
    }

    s_crossChainRebalancer[crossChainLiqManager.remoteChainSelector] = CrossChainRebalancer({
      remoteRebalancer: crossChainLiqManager.remoteRebalancer,
      localBridge: crossChainLiqManager.localBridge,
      remoteToken: crossChainLiqManager.remoteToken,
      enabled: crossChainLiqManager.enabled
    });
  }

  /// @notice Gets the local liquidity container.
  function getLocalLiquidityContainer() external view returns (address) {
    return address(s_localLiquidityContainer);
  }

  /// @notice Sets the local liquidity container.
  /// @dev Only the owner can call this function.
  function setLocalLiquidityContainer(ILiquidityContainer localLiquidityContainer) external onlyOwner {
    s_localLiquidityContainer = localLiquidityContainer;
  }
}
