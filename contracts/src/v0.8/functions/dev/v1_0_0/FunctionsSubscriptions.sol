// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IFunctionsSubscriptions} from "./interfaces/IFunctionsSubscriptions.sol";
import {IERC677Receiver} from "../../../shared/interfaces/IERC677Receiver.sol";
import {LinkTokenInterface} from "../../../shared/interfaces/LinkTokenInterface.sol";
import {IFunctionsBilling} from "./interfaces/IFunctionsBilling.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";

import {FunctionsResponse} from "./libraries/FunctionsResponse.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.0/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../../vendor/openzeppelin-solidity/v4.8.0/contracts/token/ERC20/utils/SafeERC20.sol";
import {SafeCast} from "../../../vendor/openzeppelin-solidity/v4.8.0/contracts/utils/math/SafeCast.sol";

/// @title Functions Subscriptions contract
/// @notice Contract that coordinates payment from users to the nodes of the Decentralized Oracle Network (DON).
/// @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
abstract contract FunctionsSubscriptions is IFunctionsSubscriptions, IERC677Receiver {
  using SafeERC20 for IERC20;
  using FunctionsResponse for FunctionsResponse.Commitment;

  // ================================================================
  // |                         Balance state                        |
  // ================================================================
  // link token address
  IERC20 internal immutable i_linkToken;

  // s_totalLinkBalance tracks the total LINK sent to/from
  // this contract through onTokenTransfer, cancelSubscription and oracleWithdraw.
  // A discrepancy with this contract's LINK balance indicates that someone
  // sent tokens using transfer and so we may need to use recoverFunds.
  uint96 private s_totalLinkBalance;

  /// @dev NOP balances are held as a single amount. The breakdown is held by the Coordinator.
  mapping(address coordinator => uint96 balanceJuelsLink) private s_withdrawableTokens;

  // ================================================================
  // |                      Subscription state                      |
  // ================================================================
  // Keep a count of the number of subscriptions so that its possible to
  // loop through all the current subscriptions via .getSubscription().
  uint64 private s_currentSubscriptionId;

  mapping(uint64 subscriptionId => Subscription) private s_subscriptions;

  // Maintains the list of keys in s_consumers.
  // We do this for 2 reasons:
  // 1. To be able to clean up all keys from s_consumers when canceling a subscription.
  // 2. To be able to return the list of all consumers in getSubscription.
  // Note that we need the s_consumers map to be able to directly check if a
  // consumer is valid without reading all the consumers from storage.
  mapping(address consumer => mapping(uint64 subscriptionId => Consumer)) private s_consumers;

  event SubscriptionCreated(uint64 indexed subscriptionId, address owner);
  event SubscriptionFunded(uint64 indexed subscriptionId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionConsumerAdded(uint64 indexed subscriptionId, address consumer);
  event SubscriptionConsumerRemoved(uint64 indexed subscriptionId, address consumer);
  event SubscriptionCanceled(uint64 indexed subscriptionId, address fundsRecipient, uint256 fundsAmount);
  event SubscriptionOwnerTransferRequested(uint64 indexed subscriptionId, address from, address to);
  event SubscriptionOwnerTransferred(uint64 indexed subscriptionId, address from, address to);

  error TooManyConsumers(uint16 maximumConsumers);
  error InsufficientBalance(uint96 currentBalanceJuels);
  error InvalidConsumer();
  error CannotRemoveWithPendingRequests();
  error InvalidSubscription();
  error OnlyCallableFromLink();
  error InvalidCalldata();
  error MustBeSubscriptionOwner();
  error TimeoutNotExceeded();
  error MustBeProposedOwner(address proposedOwner);
  event FundsRecovered(address to, uint256 amount);

  // ================================================================
  // |                       Request state                          |
  // ================================================================

  mapping(bytes32 requestId => bytes32 commitmentHash) internal s_requestCommitments;

  struct Receipt {
    uint96 callbackGasCostJuels;
    uint96 totalCostJuels;
  }

  event RequestTimedOut(bytes32 indexed requestId);

  // ================================================================
  // |                       Initialization                         |
  // ================================================================
  constructor(address link) {
    i_linkToken = IERC20(link);
  }

  // ================================================================
  // |                      Request/Response                        |
  // ================================================================

  /// @notice Sets a request as in-flight
  /// @dev Only callable within the Router
  function _markRequestInFlight(address client, uint64 subscriptionId, uint96 estimatedTotalCostJuels) internal {
    // Earmark subscription funds
    s_subscriptions[subscriptionId].blockedBalance += estimatedTotalCostJuels;

    // Increment sent requests
    s_consumers[client][subscriptionId].initiatedRequests += 1;
  }

  /// @notice Moves funds from one subscription account to another.
  /// @dev Only callable by the Coordinator contract that is saved in the request commitment
  function _pay(
    uint64 subscriptionId,
    uint96 estimatedTotalCostJuels,
    address client,
    uint96 adminFee,
    uint96 juelsPerGas,
    uint96 gasUsed,
    uint96 costWithoutCallbackJuels
  ) internal returns (Receipt memory) {
    uint96 callbackGasCostJuels = juelsPerGas * gasUsed;
    uint96 totalCostJuels = costWithoutCallbackJuels + adminFee + callbackGasCostJuels;

    if (
      s_subscriptions[subscriptionId].balance < totalCostJuels ||
      s_subscriptions[subscriptionId].blockedBalance < estimatedTotalCostJuels
    ) {
      revert InsufficientBalance(s_subscriptions[subscriptionId].balance);
    }

    // Charge the subscription
    s_subscriptions[subscriptionId].balance -= totalCostJuels;

    // Unblock earmarked funds
    s_subscriptions[subscriptionId].blockedBalance -= estimatedTotalCostJuels;

    // Pay the DON's fees and gas reimbursement
    s_withdrawableTokens[msg.sender] += costWithoutCallbackJuels + callbackGasCostJuels;

    // Pay out the administration fee
    s_withdrawableTokens[address(this)] += adminFee;

    // Increment finished requests
    s_consumers[client][subscriptionId].completedRequests += 1;

    return Receipt({callbackGasCostJuels: callbackGasCostJuels, totalCostJuels: totalCostJuels});
  }

  // ================================================================
  // |                      Owner methods                           |
  // ================================================================

  /// @inheritdoc IFunctionsSubscriptions
  function ownerCancelSubscription(uint64 subscriptionId) external override {
    _onlyRouterOwner();
    _isExistingSubscription(subscriptionId);
    _cancelSubscriptionHelper(subscriptionId, s_subscriptions[subscriptionId].owner, false);
  }

  /// @inheritdoc IFunctionsSubscriptions
  function recoverFunds(address to) external override {
    _onlyRouterOwner();
    uint256 externalBalance = i_linkToken.balanceOf(address(this));
    uint256 internalBalance = uint256(s_totalLinkBalance);
    if (internalBalance < externalBalance) {
      uint256 amount = externalBalance - internalBalance;
      i_linkToken.safeTransfer(to, amount);
      emit FundsRecovered(to, amount);
    }
    // If the balances are equal, nothing to be done.
  }

  // ================================================================
  // |                      Fund withdrawal                         |
  // ================================================================

  /// @inheritdoc IFunctionsSubscriptions
  function oracleWithdraw(address recipient, uint96 amount) external override {
    _whenNotPaused();

    if (amount == 0) {
      revert InvalidCalldata();
    }
    uint96 currentBalance = s_withdrawableTokens[msg.sender];
    if (currentBalance < amount) {
      revert InsufficientBalance(currentBalance);
    }
    s_withdrawableTokens[msg.sender] -= amount;
    s_totalLinkBalance -= amount;
    i_linkToken.safeTransfer(recipient, amount);
  }

  /// @notice Owner withdraw LINK earned through admin fees
  /// @notice If amount is 0 the full balance will be withdrawn
  /// @param recipient where to send the funds
  /// @param amount amount to withdraw
  function ownerWithdraw(address recipient, uint96 amount) external {
    _onlyRouterOwner();
    if (amount == 0) {
      amount = s_withdrawableTokens[address(this)];
    }
    uint96 currentBalance = s_withdrawableTokens[address(this)];
    if (currentBalance < amount) {
      revert InsufficientBalance(currentBalance);
    }
    s_withdrawableTokens[address(this)] -= amount;
    s_totalLinkBalance -= amount;

    i_linkToken.safeTransfer(recipient, amount);
  }

  // ================================================================
  // |                TransferAndCall Deposit helper                |
  // ================================================================

  // This function is to be invoked when using LINK.transferAndCall
  /// @dev Note to fund the subscription, use transferAndCall. For example
  /// @dev  LINKTOKEN.transferAndCall(
  /// @dev    address(ROUTER),
  /// @dev    amount,
  /// @dev    abi.encode(subscriptionId));
  function onTokenTransfer(address /* sender */, uint256 amount, bytes calldata data) external override {
    _whenNotPaused();
    if (msg.sender != address(i_linkToken)) {
      revert OnlyCallableFromLink();
    }
    if (data.length != 32) {
      revert InvalidCalldata();
    }
    uint64 subscriptionId = abi.decode(data, (uint64));
    if (s_subscriptions[subscriptionId].owner == address(0)) {
      revert InvalidSubscription();
    }
    // We do not check that the msg.sender is the subscription owner,
    // anyone can fund a subscription.
    uint256 oldBalance = s_subscriptions[subscriptionId].balance;
    s_subscriptions[subscriptionId].balance += uint96(amount);
    s_totalLinkBalance += uint96(amount);
    emit SubscriptionFunded(subscriptionId, oldBalance, oldBalance + amount);
  }

  // ================================================================
  // |                   Subscription management                   |
  // ================================================================

  /// @inheritdoc IFunctionsSubscriptions
  function getTotalBalance() external view override returns (uint96) {
    return s_totalLinkBalance;
  }

  /// @inheritdoc IFunctionsSubscriptions
  function getSubscriptionCount() external view override returns (uint64) {
    return s_currentSubscriptionId;
  }

  /// @inheritdoc IFunctionsSubscriptions
  function getSubscription(uint64 subscriptionId) public view override returns (Subscription memory) {
    _isExistingSubscription(subscriptionId);
    return s_subscriptions[subscriptionId];
  }

  /// @inheritdoc IFunctionsSubscriptions
  function getSubscriptionsInRange(
    uint64 subscriptionIdStart,
    uint64 subscriptionIdEnd
  ) external view override returns (Subscription[] memory subscriptions) {
    if (
      subscriptionIdStart > subscriptionIdEnd ||
      subscriptionIdEnd > s_currentSubscriptionId ||
      s_currentSubscriptionId == 0
    ) {
      revert InvalidCalldata();
    }

    subscriptions = new Subscription[]((subscriptionIdEnd - subscriptionIdStart) + 1);
    for (uint256 i = 0; i <= subscriptionIdEnd - subscriptionIdStart; ++i) {
      subscriptions[i] = s_subscriptions[uint64(subscriptionIdStart + i)];
    }

    return subscriptions;
  }

  /// @inheritdoc IFunctionsSubscriptions
  function getConsumer(address client, uint64 subscriptionId) public view override returns (Consumer memory) {
    return s_consumers[client][subscriptionId];
  }

  /// @dev Used within this file & FunctionsRouter.sol
  function _isExistingSubscription(uint64 subscriptionId) internal view {
    if (s_subscriptions[subscriptionId].owner == address(0)) {
      revert InvalidSubscription();
    }
  }

  /// @dev Used within FunctionsRouter.sol
  function _isAllowedConsumer(address client, uint64 subscriptionId) internal view {
    if (!s_consumers[client][subscriptionId].allowed) {
      revert InvalidConsumer();
    }
  }

  /// @inheritdoc IFunctionsSubscriptions
  function createSubscription() external override returns (uint64 subscriptionId) {
    _whenNotPaused();
    _onlySenderThatAcceptedToS();

    subscriptionId = ++s_currentSubscriptionId;
    s_subscriptions[subscriptionId] = Subscription({
      balance: 0,
      blockedBalance: 0,
      owner: msg.sender,
      proposedOwner: address(0),
      consumers: new address[](0),
      flags: bytes32(0)
    });

    emit SubscriptionCreated(subscriptionId, msg.sender);

    return subscriptionId;
  }

  /// @inheritdoc IFunctionsSubscriptions
  function createSubscriptionWithConsumer(address consumer) external override returns (uint64 subscriptionId) {
    _whenNotPaused();
    _onlySenderThatAcceptedToS();

    subscriptionId = ++s_currentSubscriptionId;
    s_subscriptions[subscriptionId] = Subscription({
      balance: 0,
      blockedBalance: 0,
      owner: msg.sender,
      proposedOwner: address(0),
      consumers: new address[](0),
      flags: bytes32(0)
    });

    s_subscriptions[subscriptionId].consumers.push(consumer);
    s_consumers[consumer][subscriptionId].allowed = true;

    emit SubscriptionCreated(subscriptionId, msg.sender);
    emit SubscriptionConsumerAdded(subscriptionId, consumer);

    return subscriptionId;
  }

  /// @inheritdoc IFunctionsSubscriptions
  function proposeSubscriptionOwnerTransfer(uint64 subscriptionId, address newOwner) external override {
    _whenNotPaused();
    _onlySubscriptionOwner(subscriptionId);
    _onlySenderThatAcceptedToS();

    if (newOwner == address(0) || s_subscriptions[subscriptionId].proposedOwner == newOwner) {
      revert InvalidCalldata();
    }

    s_subscriptions[subscriptionId].proposedOwner = newOwner;
    emit SubscriptionOwnerTransferRequested(subscriptionId, msg.sender, newOwner);
  }

  /// @inheritdoc IFunctionsSubscriptions
  function acceptSubscriptionOwnerTransfer(uint64 subscriptionId) external override {
    _whenNotPaused();
    _onlySenderThatAcceptedToS();

    address previousOwner = s_subscriptions[subscriptionId].owner;
    address proposedOwner = s_subscriptions[subscriptionId].proposedOwner;
    if (proposedOwner != msg.sender) {
      revert MustBeProposedOwner(proposedOwner);
    }
    s_subscriptions[subscriptionId].owner = msg.sender;
    s_subscriptions[subscriptionId].proposedOwner = address(0);
    emit SubscriptionOwnerTransferred(subscriptionId, previousOwner, msg.sender);
  }

  /// @inheritdoc IFunctionsSubscriptions
  function removeConsumer(uint64 subscriptionId, address consumer) external override {
    _whenNotPaused();
    _onlySubscriptionOwner(subscriptionId);
    _onlySenderThatAcceptedToS();

    Consumer memory consumerData = s_consumers[consumer][subscriptionId];
    _isAllowedConsumer(consumer, subscriptionId);
    if (consumerData.initiatedRequests != consumerData.completedRequests) {
      revert CannotRemoveWithPendingRequests();
    }
    // Note bounded by config.maxConsumers
    address[] memory consumers = s_subscriptions[subscriptionId].consumers;
    for (uint256 i = 0; i < consumers.length; ++i) {
      if (consumers[i] == consumer) {
        // Storage write to preserve last element
        s_subscriptions[subscriptionId].consumers[i] = consumers[consumers.length - 1];
        // Storage remove last element
        s_subscriptions[subscriptionId].consumers.pop();
        break;
      }
    }
    delete s_consumers[consumer][subscriptionId];
    emit SubscriptionConsumerRemoved(subscriptionId, consumer);
  }

  /// @dev Overriden in FunctionsRouter.sol
  function _getMaxConsumers() internal view virtual returns (uint16);

  /// @inheritdoc IFunctionsSubscriptions
  function addConsumer(uint64 subscriptionId, address consumer) external override {
    _whenNotPaused();
    _onlySubscriptionOwner(subscriptionId);
    _onlySenderThatAcceptedToS();

    // Already maxed, cannot add any more consumers.
    uint16 maximumConsumers = _getMaxConsumers();
    if (s_subscriptions[subscriptionId].consumers.length >= maximumConsumers) {
      revert TooManyConsumers(maximumConsumers);
    }
    if (s_consumers[consumer][subscriptionId].allowed) {
      // Idempotence - do nothing if already added.
      // Ensures uniqueness in s_subscriptions[subscriptionId].consumers.
      return;
    }

    s_consumers[consumer][subscriptionId].allowed = true;
    s_subscriptions[subscriptionId].consumers.push(consumer);

    emit SubscriptionConsumerAdded(subscriptionId, consumer);
  }

  /// @dev Overriden in FunctionsRouter.sol
  function _getSubscriptionDepositDetails() internal virtual returns (uint16, uint72);

  function _cancelSubscriptionHelper(uint64 subscriptionId, address toAddress, bool checkDepositRefundability) private {
    Subscription memory subscription = s_subscriptions[subscriptionId];
    uint96 balance = subscription.balance;
    uint64 completedRequests = 0;

    // NOTE: loop iterations are bounded by config.maxConsumers
    // If no consumers, does nothing.
    for (uint256 i = 0; i < subscription.consumers.length; ++i) {
      address consumer = subscription.consumers[i];
      completedRequests += s_consumers[consumer][subscriptionId].completedRequests;
      delete s_consumers[consumer][subscriptionId];
    }
    delete s_subscriptions[subscriptionId];

    (uint16 subscriptionDepositMinimumRequests, uint72 subscriptionDepositJuels) = _getSubscriptionDepositDetails();

    // If subscription has not made enough requests, deposit will be forfeited
    if (checkDepositRefundability && completedRequests < subscriptionDepositMinimumRequests) {
      uint96 deposit = subscriptionDepositJuels > balance ? balance : subscriptionDepositJuels;
      if (deposit > 0) {
        s_withdrawableTokens[address(this)] += deposit;
        balance -= deposit;
      }
    }

    if (balance > 0) {
      s_totalLinkBalance -= balance;
      i_linkToken.safeTransfer(toAddress, uint256(balance));
    }
    emit SubscriptionCanceled(subscriptionId, toAddress, balance);
  }

  /// @inheritdoc IFunctionsSubscriptions
  function cancelSubscription(uint64 subscriptionId, address to) external override {
    _whenNotPaused();
    _onlySubscriptionOwner(subscriptionId);
    _onlySenderThatAcceptedToS();

    if (pendingRequestExists(subscriptionId)) {
      revert CannotRemoveWithPendingRequests();
    }

    _cancelSubscriptionHelper(subscriptionId, to, true);
  }

  /// @inheritdoc IFunctionsSubscriptions
  function pendingRequestExists(uint64 subscriptionId) public view override returns (bool) {
    address[] memory consumers = s_subscriptions[subscriptionId].consumers;
    // NOTE: loop iterations are bounded by config.maxConsumers
    for (uint256 i = 0; i < consumers.length; ++i) {
      Consumer memory consumer = s_consumers[consumers[i]][subscriptionId];
      if (consumer.initiatedRequests != consumer.completedRequests) {
        return true;
      }
    }
    return false;
  }

  /// @inheritdoc IFunctionsSubscriptions
  function setFlags(uint64 subscriptionId, bytes32 flags) external override {
    _onlyRouterOwner();
    _isExistingSubscription(subscriptionId);
    s_subscriptions[subscriptionId].flags = flags;
  }

  /// @inheritdoc IFunctionsSubscriptions
  function getFlags(uint64 subscriptionId) public view returns (bytes32) {
    return s_subscriptions[subscriptionId].flags;
  }

  // ================================================================
  // |                        Request Timeout                       |
  // ================================================================

  /// @inheritdoc IFunctionsSubscriptions
  function timeoutRequests(FunctionsResponse.Commitment[] calldata requestsToTimeoutByCommitment) external override {
    _whenNotPaused();

    for (uint256 i = 0; i < requestsToTimeoutByCommitment.length; ++i) {
      FunctionsResponse.Commitment memory request = requestsToTimeoutByCommitment[i];
      bytes32 requestId = request.requestId;
      uint64 subscriptionId = request.subscriptionId;

      // Check that request ID is valid
      if (keccak256(abi.encode(request)) != s_requestCommitments[requestId]) {
        revert InvalidCalldata();
      }

      // Check that request has exceeded allowed request time
      if (block.timestamp < request.timeoutTimestamp) {
        revert TimeoutNotExceeded();
      }

      // Notify the Coordinator that the request should no longer be fulfilled
      IFunctionsBilling(request.coordinator).deleteCommitment(requestId);
      // Release the subscription's balance that had been earmarked for the request
      s_subscriptions[subscriptionId].blockedBalance -= request.estimatedTotalCostJuels;
      s_consumers[request.client][subscriptionId].completedRequests += 1;
      // Delete commitment within Router state
      delete s_requestCommitments[requestId];

      emit RequestTimedOut(requestId);
    }
  }

  // ================================================================
  // |                         Modifiers                            |
  // ================================================================

  function _onlySubscriptionOwner(uint64 subscriptionId) internal view {
    address owner = s_subscriptions[subscriptionId].owner;
    if (owner == address(0)) {
      revert InvalidSubscription();
    }
    if (msg.sender != owner) {
      revert MustBeSubscriptionOwner();
    }
  }

  /// @dev Overriden in FunctionsRouter.sol
  function _onlySenderThatAcceptedToS() internal virtual;

  /// @dev Overriden in FunctionsRouter.sol
  function _onlyRouterOwner() internal virtual;

  /// @dev Overriden in FunctionsRouter.sol
  function _whenNotPaused() internal virtual;
}
