// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IBridgeAdapter} from "./interfaces/IBridge.sol";
import {ILiquidityManager} from "./interfaces/ILiquidityManager.sol";
import {ILiquidityContainer} from "./interfaces/ILiquidityContainer.sol";
import {IWrappedNative} from "../ccip/interfaces/IWrappedNative.sol";

import {OCR3Base} from "./ocr/OCR3Base.sol";

import {IERC20} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @notice LiquidityManager for a single token over multiple chains.
/// @dev This contract is designed to be used with the LockReleaseTokenPool contract but
/// isn't constrained to it. It can be used with any contract that implements the ILiquidityContainer
/// interface.
/// @dev The OCR3 DON should only be able to transfer funds to other pre-approved contracts
/// on other chains. Under no circumstances should it be able to transfer funds to arbitrary
/// addresses. The owner is therefore in full control of the funds in this contract, not the DON.
/// This is a security feature. The worst that can happen is that the DON can lock up funds in
/// bridges, but it can't steal them.
/// @dev References to local mean logic on the same chain as this contract is deployed on.
/// References to remote mean logic on other chains.
contract LiquidityManager is ILiquidityManager, OCR3Base {
  using SafeERC20 for IERC20;

  error ZeroAddress();
  error InvalidRemoteChain(uint64 chainSelector);
  error ZeroChainSelector();
  error InsufficientLiquidity(uint256 requested, uint256 available, uint256 reserve);
  error EmptyReport();
  error TransferFailed();
  error OnlyFinanceRole();

  /// @notice Emitted when a finalization step is completed without funds being available.
  /// @param ocrSeqNum The OCR sequence number of the report.
  /// @param remoteChainSelector The chain selector of the remote chain funds are coming from.
  /// @param bridgeSpecificData The bridge specific data that was used to finalize the transfer.
  event FinalizationStepCompleted(
    uint64 indexed ocrSeqNum,
    uint64 indexed remoteChainSelector,
    bytes bridgeSpecificData
  );

  /// @notice Emitted when the CLL finance role is set.
  /// @param financeRole The address of the new finance role.
  event FinanceRoleSet(address financeRole);

  /// @notice Emitted when liquidity is transferred to another chain, or received from another chain.
  /// @param ocrSeqNum The OCR sequence number of the report.
  /// @param fromChainSelector The chain selector of the chain the funds are coming from.
  /// In the event fromChainSelector == i_localChainSelector, this is an outgoing transfer.
  /// Otherwise, it is an incoming transfer.
  /// @param toChainSelector The chain selector of the chain the funds are going to.
  /// In the event toChainSelector == i_localChainSelector, this is an incoming transfer.
  /// Otherwise, it is an outgoing transfer.
  /// @param to The address the funds are going to.
  /// If this is address(this), the funds are arriving in this contract.
  /// @param amount The amount of tokens being transferred.
  /// @param bridgeSpecificData The bridge specific data that was passed to the local bridge adapter
  /// when transferring the funds.
  /// @param bridgeReturnData The return data from the local bridge adapter when transferring the funds.
  event LiquidityTransferred(
    uint64 indexed ocrSeqNum,
    uint64 indexed fromChainSelector,
    uint64 indexed toChainSelector,
    address to,
    uint256 amount,
    bytes bridgeSpecificData,
    bytes bridgeReturnData
  );

  /// @notice Emitted when liquidity is added to the local liquidity container.
  /// @param provider The address of the provider that added the liquidity.
  /// @param amount The amount of liquidity that was added.
  event LiquidityAddedToContainer(address indexed provider, uint256 indexed amount);

  /// @notice Emitted when liquidity is removed from the local liquidity container.
  /// @param remover The address of the remover that removed the liquidity.
  /// @param amount The amount of liquidity that was removed.
  event LiquidityRemovedFromContainer(address indexed remover, uint256 indexed amount);

  /// @notice Emitted when the local liquidity container is set.
  /// @param newLiquidityContainer The address of the new liquidity container.
  event LiquidityContainerSet(address indexed newLiquidityContainer);

  /// @notice Emitted when the minimum liquidity is set.
  /// @param oldBalance The old minimum liquidity.
  /// @param newBalance The new minimum liquidity.
  event MinimumLiquiditySet(uint256 oldBalance, uint256 newBalance);

  /// @notice Emitted when someone sends native to this contract
  /// @param amount The amount of native deposited
  /// @param depositor The address that deposited the native
  event NativeDeposited(uint256 amount, address depositor);

  /// @notice Emitted when native balance is withdrawn by contract owner
  /// @param amount The amount of native withdrawn
  /// @param destination The address the native is sent to
  event NativeWithdrawn(uint256 amount, address destination);

  /// @notice Emitted when a cross chain rebalancer is set.
  /// @param remoteChainSelector The chain selector of the remote chain.
  /// @param localBridge The local bridge adapter that will be used to transfer funds.
  /// @param remoteToken The address of the token on the remote chain.
  /// @param remoteRebalancer The address of the remote rebalancer contract.
  /// @param enabled Whether the rebalancer is enabled.
  event CrossChainRebalancerSet(
    uint64 indexed remoteChainSelector,
    IBridgeAdapter localBridge,
    address remoteToken,
    address remoteRebalancer,
    bool enabled
  );

  /// @notice Emitted when a finalization step fails.
  /// @param ocrSeqNum The OCR sequence number of the report.
  /// @param remoteChainSelector The chain selector of the remote chain funds are coming from.
  /// @param bridgeSpecificData The bridge specific data that was used to finalize the transfer.
  /// @param reason The reason the finalization failed.
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

  string public constant override typeAndVersion = "LiquidityManager 1.0.0-dev";

  /// @notice The token that this pool manages liquidity for.
  IERC20 public immutable i_localToken;

  /// @notice The chain selector belonging to the chain this pool is deployed on.
  uint64 internal immutable i_localChainSelector;

  /// @notice The target balance defines the expected amount of tokens for this network.
  /// Setting the balance to 0 will disable any automated rebalancing operations.
  uint256 internal s_minimumLiquidity;

  /// @notice Mapping of chain selector to liquidity container on other chains
  mapping(uint64 chainSelector => CrossChainRebalancer) private s_crossChainRebalancer;

  uint64[] private s_supportedDestChains;

  /// @notice The liquidity container on the local chain
  /// @dev In the case of CCIP, this would be the token pool.
  ILiquidityContainer private s_localLiquidityContainer;

  /// @notice The CLL finance team multisig
  address private s_finance;

  constructor(
    IERC20 token,
    uint64 localChainSelector,
    ILiquidityContainer localLiquidityContainer,
    uint256 minimumLiquidity,
    address finance
  ) OCR3Base() {
    if (localChainSelector == 0) {
      revert ZeroChainSelector();
    }

    if (address(token) == address(0) || address(localLiquidityContainer) == address(0)) {
      revert ZeroAddress();
    }
    i_localToken = token;
    i_localChainSelector = localChainSelector;
    s_localLiquidityContainer = localLiquidityContainer;
    s_minimumLiquidity = minimumLiquidity;
    s_finance = finance;
  }

  // ================================================================
  // │                      Native Management                       │
  // ================================================================

  receive() external payable {
    emit NativeDeposited(msg.value, msg.sender);
  }

  /// @notice withdraw native balance
  function withdrawNative(uint256 amount, address payable destination) external onlyFinance {
    (bool success, ) = destination.call{value: amount}("");
    if (!success) revert TransferFailed();

    emit NativeWithdrawn(amount, destination);
  }

  // ================================================================
  // │                     Liquidity Management                     │
  // ================================================================

  /// @inheritdoc ILiquidityManager
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
    i_localToken.safeApprove(address(s_localLiquidityContainer), amount);
    s_localLiquidityContainer.provideLiquidity(amount);

    emit LiquidityAddedToContainer(msg.sender, amount);
  }

  /// @notice Removes liquidity from the system and sends it to the caller, so the owner.
  /// @dev Only the owner can call this function.
  function removeLiquidity(uint256 amount) external onlyFinance {
    uint256 currentBalance = getLiquidity();
    if (currentBalance < amount) {
      revert InsufficientLiquidity(amount, currentBalance, 0);
    }

    s_localLiquidityContainer.withdrawLiquidity(amount);
    i_localToken.safeTransfer(msg.sender, amount);

    emit LiquidityRemovedFromContainer(msg.sender, amount);
  }

  /// @notice escape hatch to manually withdraw any ERC20 token from the LM contract
  /// @param token The address of the token to withdraw
  /// @param amount The amount of tokens to withdraw
  /// @param destination The address to send the tokens to
  function withdrawERC20(address token, uint256 amount, address destination) external onlyFinance {
    IERC20(token).safeTransfer(destination, amount);
  }

  /// @notice Transfers liquidity to another chain.
  /// @dev This function is a public version of the internal _rebalanceLiquidity function.
  /// to allow the owner to also initiate a rebalancing when needed.
  function rebalanceLiquidity(
    uint64 chainSelector,
    uint256 amount,
    uint256 nativeBridgeFee,
    bytes calldata bridgeSpecificPayload
  ) external onlyFinance {
    _rebalanceLiquidity(chainSelector, amount, nativeBridgeFee, type(uint64).max, bridgeSpecificPayload);
  }

  /// @notice Finalizes liquidity from another chain.
  /// @dev This function is a public version of the internal _receiveLiquidity function.
  /// to allow the owner to also initiate a finalization when needed.
  function receiveLiquidity(
    uint64 remoteChainSelector,
    uint256 amount,
    bool shouldWrapNative,
    bytes calldata bridgeSpecificPayload
  ) external onlyFinance {
    _receiveLiquidity(remoteChainSelector, amount, bridgeSpecificPayload, shouldWrapNative, type(uint64).max);
  }

  /// @notice Transfers liquidity to another chain.
  /// @dev Called by both the owner and the DON.
  /// @param chainSelector The chain selector of the chain to transfer liquidity to.
  /// @param tokenAmount The amount of tokens to transfer.
  /// @param nativeBridgeFee The fee to pay to the bridge.
  /// @param ocrSeqNum The OCR sequence number of the report.
  /// @param bridgeSpecificPayload The bridge specific data to pass to the bridge adapter.
  function _rebalanceLiquidity(
    uint64 chainSelector,
    uint256 tokenAmount,
    uint256 nativeBridgeFee,
    uint64 ocrSeqNum,
    bytes memory bridgeSpecificPayload
  ) internal {
    uint256 currentBalance = getLiquidity();
    uint256 minBalance = s_minimumLiquidity;
    if (currentBalance < minBalance || currentBalance - minBalance < tokenAmount) {
      revert InsufficientLiquidity(tokenAmount, currentBalance, minBalance);
    }

    CrossChainRebalancer memory remoteLiqManager = s_crossChainRebalancer[chainSelector];

    if (!remoteLiqManager.enabled) {
      revert InvalidRemoteChain(chainSelector);
    }

    // XXX: Could be optimized by withdrawing once and then sending to all destinations
    s_localLiquidityContainer.withdrawLiquidity(tokenAmount);
    i_localToken.safeApprove(address(remoteLiqManager.localBridge), tokenAmount);

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

  /// @notice Receives liquidity from another chain.
  /// @dev Called by both the owner and the DON.
  /// @param remoteChainSelector The chain selector of the chain to receive liquidity from.
  /// @param amount The amount of tokens to receive.
  /// @param bridgeSpecificPayload The bridge specific data to pass to the bridge adapter finalizeWithdrawERC20 call.
  /// @param shouldWrapNative Whether the token should be wrapped before injecting it into the liquidity container.
  /// This only applies to native tokens wrapper contracts, e.g WETH.
  /// @param ocrSeqNum The OCR sequence number of the report.
  function _receiveLiquidity(
    uint64 remoteChainSelector,
    uint256 amount,
    bytes memory bridgeSpecificPayload,
    bool shouldWrapNative,
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
        address(this), // localReceiver: this contract
        bridgeSpecificPayload
      )
    returns (bool fundsAvailable) {
      if (fundsAvailable) {
        // finalization was successful and we can inject the liquidity into the container.
        // approve and liquidity container should transferFrom.
        _injectLiquidity(amount, ocrSeqNum, remoteChainSelector, bridgeSpecificPayload, shouldWrapNative);
      } else {
        // a finalization step was completed, but funds are not available.
        // hence, we cannot inject any liquidity yet.
        emit FinalizationStepCompleted(ocrSeqNum, remoteChainSelector, bridgeSpecificPayload);
      }

      // return here on the happy path.
      // sad path is when finalizeWithdrawERC20 reverts, which is handled after the catch block.
      return;
    } catch (bytes memory lowLevelData) {
      // failed to finalize the withdrawal.
      // this could mean that the withdrawal was already finalized
      // or that the withdrawal failed.
      // we assume the former and continue
      emit FinalizationFailed(ocrSeqNum, remoteChainSelector, bridgeSpecificPayload, lowLevelData);
    }

    // if we reach this point, the finalization failed.
    // since we don't have enough information to know why it failed,
    // we assume that it failed because the withdrawal was already finalized,
    // and that the funds are available.
    _injectLiquidity(amount, ocrSeqNum, remoteChainSelector, bridgeSpecificPayload, shouldWrapNative);
  }

  /// @notice Injects liquidity into the local liquidity container.
  /// @param amount The amount of tokens to inject.
  /// @param ocrSeqNum The OCR sequence number of the report.
  /// @param remoteChainSelector The chain selector of the remote chain.
  /// @param bridgeSpecificPayload The bridge specific data passed to the bridge adapter finalizeWithdrawERC20 call.
  /// @param shouldWrapNative Whether the token should be wrapped before injecting it into the liquidity container.
  function _injectLiquidity(
    uint256 amount,
    uint64 ocrSeqNum,
    uint64 remoteChainSelector,
    bytes memory bridgeSpecificPayload,
    bool shouldWrapNative
  ) private {
    // We trust the DON or the owner (the only two actors who can end up calling this function)
    // to correctly set the shouldWrapNative flag.
    // Some bridges only bridge native and not wrapped native.
    // In such a case we need to re-wrap the native in order to inject it into the liquidity container.
    // TODO: escape hatch in case of bug?
    if (shouldWrapNative) {
      IWrappedNative(address(i_localToken)).deposit{value: amount}();
    }

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

  /// @notice Process the OCR report.
  /// @dev Called by OCR3Base's transmit() function.
  function _report(bytes calldata report, uint64 ocrSeqNum) internal override {
    ILiquidityManager.LiquidityInstructions memory instructions = abi.decode(
      report,
      (ILiquidityManager.LiquidityInstructions)
    );

    uint256 sendInstructions = instructions.sendLiquidityParams.length;
    uint256 receiveInstructions = instructions.receiveLiquidityParams.length;

    // There should always be instructions to send or receive, if not, the report is invalid
    // and we revert to save the gas of the signature validation of OCR.
    if (sendInstructions == 0 && receiveInstructions == 0) {
      revert EmptyReport();
    }

    for (uint256 i = 0; i < sendInstructions; ++i) {
      _rebalanceLiquidity(
        instructions.sendLiquidityParams[i].remoteChainSelector,
        instructions.sendLiquidityParams[i].amount,
        instructions.sendLiquidityParams[i].nativeBridgeFee,
        ocrSeqNum,
        instructions.sendLiquidityParams[i].bridgeData
      );
    }

    for (uint256 i = 0; i < receiveInstructions; ++i) {
      _receiveLiquidity(
        instructions.receiveLiquidityParams[i].remoteChainSelector,
        instructions.receiveLiquidityParams[i].amount,
        instructions.receiveLiquidityParams[i].bridgeData,
        instructions.receiveLiquidityParams[i].shouldWrapNative,
        ocrSeqNum
      );
    }
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
  function setCrossChainRebalancers(CrossChainRebalancerArgs[] calldata crossChainRebalancers) external onlyOwner {
    for (uint256 i = 0; i < crossChainRebalancers.length; ++i) {
      _setCrossChainRebalancer(crossChainRebalancers[i]);
    }
  }

  function setCrossChainRebalancer(CrossChainRebalancerArgs calldata crossChainLiqManager) external onlyOwner {
    _setCrossChainRebalancer(crossChainLiqManager);
  }

  /// @notice Sets a single cross chain liquidity manager.
  /// @dev Will update the list of supported dest chains if the chain is new.
  function _setCrossChainRebalancer(CrossChainRebalancerArgs calldata crossChainLiqManager) internal {
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

    emit CrossChainRebalancerSet(
      crossChainLiqManager.remoteChainSelector,
      crossChainLiqManager.localBridge,
      crossChainLiqManager.remoteToken,
      crossChainLiqManager.remoteRebalancer,
      crossChainLiqManager.enabled
    );
  }

  /// @notice Gets the local liquidity container.
  function getLocalLiquidityContainer() external view returns (address) {
    return address(s_localLiquidityContainer);
  }

  /// @notice Sets the local liquidity container.
  /// @dev Only the owner can call this function.
  function setLocalLiquidityContainer(ILiquidityContainer localLiquidityContainer) external onlyOwner {
    if (address(localLiquidityContainer) == address(0)) {
      revert ZeroAddress();
    }
    s_localLiquidityContainer = localLiquidityContainer;

    emit LiquidityContainerSet(address(localLiquidityContainer));
  }

  /// @notice Gets the target tokens balance.
  function getMinimumLiquidity() external view returns (uint256) {
    return s_minimumLiquidity;
  }

  /// @notice Sets the target tokens balance.
  /// @dev Only the owner can call this function.
  function setMinimumLiquidity(uint256 minimumLiquidity) external onlyOwner {
    uint256 oldLiquidity = s_minimumLiquidity;
    s_minimumLiquidity = minimumLiquidity;
    emit MinimumLiquiditySet(oldLiquidity, s_minimumLiquidity);
  }

  /// @notice Gets the CLL finance team multisig address
  function getFinanceRole() external view returns (address) {
    return s_finance;
  }

  /// @notice Sets the finance team multisig address
  /// @dev Only the owner can call this function.
  function setFinanceRole(address finance) external onlyOwner {
    s_finance = finance;
    emit FinanceRoleSet(finance);
  }

  modifier onlyFinance() {
    if (msg.sender != s_finance) revert OnlyFinanceRole();
    _;
  }
}
