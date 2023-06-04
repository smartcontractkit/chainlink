// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {IFunctionsSubscriptions} from "./interfaces/IFunctionsSubscriptions.sol";
import {ERC677ReceiverInterface} from "../../../interfaces/ERC677ReceiverInterface.sol";
import {LinkTokenInterface} from "../../../interfaces/LinkTokenInterface.sol";

/**
 * @title Functions Subscriptions contract
 * @notice Contract that coordinates payment from users to the nodes of the Decentralized Oracle Network (DON).
 * @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
 */
abstract contract FunctionsSubscriptions is IFunctionsSubscriptions, ERC677ReceiverInterface {
  LinkTokenInterface private LINK;

  // We need to maintain a list of consuming addresses.
  // This bound ensures we are able to loop over them as needed.
  // Should a user require more consumers, they can use multiple subscriptions.
  uint16 public constant MAX_CONSUMERS = 100;

  // ================================================================
  // |                      Subscription state                      |
  // ================================================================

  // We make the sub count public so that its possible to
  // get all the current subscriptions via getSubscription.
  uint64 private s_currentsubscriptionId;

  // s_totalBalance tracks the total link sent to/from
  // this contract through onTokenTransfer, cancelSubscription and oracleWithdraw.
  // A discrepancy with this contract's link balance indicates someone
  // sent tokens using transfer and so we may need to use recoverFunds.
  uint96 private s_totalBalance;

  mapping(uint64 => IFunctionsSubscriptions.Subscription) /* subscriptionId */ /* subscription */
    private s_subscriptions;

  mapping(address => mapping(uint64 => IFunctionsSubscriptions.Consumer)) /* consumer */ /* subscriptionId */ /* Consumer data */
    private s_consumers;

  event SubscriptionCreated(uint64 indexed subscriptionId, address owner);
  event SubscriptionFunded(uint64 indexed subscriptionId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionConsumerAdded(uint64 indexed subscriptionId, address consumer);
  event SubscriptionConsumerRemoved(uint64 indexed subscriptionId, address consumer);
  event SubscriptionCanceled(uint64 indexed subscriptionId, address to, uint256 amount);
  event SubscriptionOwnerTransferRequested(uint64 indexed subscriptionId, address from, address to);
  event SubscriptionOwnerTransferred(uint64 indexed subscriptionId, address from, address to);

  error Reentrant();
  error TooManyConsumers();
  error InsufficientBalance();
  error InvalidConsumer(uint64 subscriptionId, address consumer);
  error ConsumerRequestsInFlight();
  error InvalidSubscription();
  error OnlyCallableFromLink();
  error InvalidCalldata();
  error MustBeSubOwner(address owner);
  error PendingRequestExists();
  error MustBeRequestedOwner(address proposedOwner);
  error BalanceInvariantViolated(uint256 internalBalance, uint256 externalBalance); // Should never happen
  event FundsRecovered(address to, uint256 amount);

  mapping(address => uint96) /* oracle node */ /* LINK balance */
    private s_withdrawableTokens;

  struct Request {
    address coordinator;
    uint64 subscriptionId;
    uint96 estimatedCost;
    //timestamp?
  }

  mapping(bytes32 => Request) /* request ID */ /* Request data */
    private s_requests;

  // ================================================================
  // |                       Initialization                         |
  // ================================================================
  constructor(address link) {
    LINK = LinkTokenInterface(link);
  }

  function getTotalBalance() external view returns (uint256) {
    return s_totalBalance;
  }

  // ================================================================
  // |                      Getter methods                          |
  // ================================================================
  function _isValidSubscription(uint64 subscriptionId) internal view {
    if (s_subscriptions[subscriptionId].owner == address(0)) {
      revert InvalidSubscription();
    }
  }

  function _isValidConsumer(address client, uint64 subscriptionId) internal view {
    if (s_consumers[client][subscriptionId].allowed == false) {
      revert InvalidConsumer(subscriptionId, client);
    }
  }

  function getCurrentsubscriptionId() external view returns (uint64) {
    return s_currentsubscriptionId;
  }

  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function getSubscription(uint64 subscriptionId)
    external
    view
    returns (
      uint96 balance,
      uint96 blockedBalance,
      address owner,
      address requestedOwner,
      address[] memory consumers
    )
  {
    _isValidSubscription(subscriptionId);
    return (
      s_subscriptions[subscriptionId].balance,
      s_subscriptions[subscriptionId].blockedBalance,
      s_subscriptions[subscriptionId].owner,
      s_subscriptions[subscriptionId].requestedOwner,
      s_subscriptions[subscriptionId].consumers
    );
  }

  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function getConsumer(address client, uint64 subscriptionId)
    external
    view
    onlyRoute
    returns (
      bool allowed,
      uint64 initiatedRequests,
      uint64 completedRequests
    )
  {
    return (
      s_consumers[client][subscriptionId].allowed,
      s_consumers[client][subscriptionId].initiatedRequests,
      s_consumers[client][subscriptionId].completedRequests
    );
  }

  // ================================================================
  // |                 Internal Payment methods                     |
  // ================================================================
  /**
   * @notice Blocks funds on a subscription account to be used for an in flight request
   * @dev Only callable within the Router
   * @param client - the consumer contract that initiated the request
   * @param subscriptionId - The subscription ID to block funds for
   * @param amount - The amount to transfer
   */
  function _blockBalance(
    address client,
    uint64 subscriptionId,
    uint96 amount,
    bytes32 requestId,
    address coordinator
  ) internal {
    s_subscriptions[subscriptionId].blockedBalance += amount;
    s_consumers[client][subscriptionId].initiatedRequests += 1;
    s_requests[requestId] = Request(coordinator, subscriptionId, amount);
  }

  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function unblockBalance(
    address client,
    uint64 subscriptionId,
    uint96 amount
  ) external onlyRoute {
    s_subscriptions[subscriptionId].blockedBalance -= amount;
    s_consumers[client][subscriptionId].completedRequests += 1;
  }

  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function pay(
    uint64 from,
    address to,
    uint96 amount
  ) external onlyRoute {
    // TODO: optimize into one call to pay every operator
    s_subscriptions[from].balance -= amount;
    s_withdrawableTokens[to] += amount;
  }

  // ================================================================
  // |                      Owner methods                           |
  // ================================================================
  /**
   * @notice Owner cancel subscription, sends remaining link directly to the subscription owner.
   * @param subscriptionId subscription id
   * @dev notably can be called even if there are pending requests, outstanding ones may fail onchain
   */
  function ownerCancelSubscription(uint64 subscriptionId) external onlyRouterOwner {
    address owner = s_subscriptions[subscriptionId].owner;
    if (owner == address(0)) {
      revert InvalidSubscription();
    }
    cancelSubscriptionHelper(subscriptionId, owner);
  }

  /**
   * @notice Recover link sent with transfer instead of transferAndCall.
   * @param to address to send link to
   */
  function recoverFunds(address to) external onlyRouterOwner {
    uint256 externalBalance = LINK.balanceOf(address(this));
    uint256 internalBalance = uint256(s_totalBalance);
    if (internalBalance > externalBalance) {
      revert BalanceInvariantViolated(internalBalance, externalBalance);
    }
    if (internalBalance < externalBalance) {
      uint256 amount = externalBalance - internalBalance;
      LINK.transfer(to, amount);
      emit FundsRecovered(to, amount);
    }
    // If the balances are equal, nothing to be done.
  }

  // ================================================================
  // |                   Node Operator methods                      |
  // ================================================================
  /*
   * @notice Oracle withdraw LINK earned through fulfilling requests
   * @notice If amount is 0 the full balance will be withdrawn
   * @notice Both signing and transmitting wallets will have a balance to withdraw
   * @param recipient where to send the funds
   * @param amount amount to withdraw
   */
  function oracleWithdraw(address recipient, uint96 amount) external nonReentrant {
    if (amount == 0) {
      amount = s_withdrawableTokens[msg.sender];
    }
    if (s_withdrawableTokens[msg.sender] < amount) {
      revert InsufficientBalance();
    }
    s_withdrawableTokens[msg.sender] -= amount;
    s_totalBalance -= amount;
    if (!LINK.transfer(recipient, amount)) {
      revert InsufficientBalance();
    }
  }

  // ================================================================
  // |                   Deposit helper method                      |
  // ================================================================
  function onTokenTransfer(
    address, /* sender */
    uint256 amount,
    bytes calldata data
  ) external override nonReentrant {
    if (msg.sender != address(LINK)) {
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
    s_totalBalance += uint96(amount);
    emit SubscriptionFunded(subscriptionId, oldBalance, oldBalance + amount);
  }

  // ================================================================
  // |                    Subscription methods                      |
  // ================================================================
  /**
   * @notice Create a new subscription.
   * @return subscriptionId - A unique subscription id.
   * @dev You can manage the consumer set dynamically with addConsumer/removeConsumer.
   * @dev Note to fund the subscription, use transferAndCall. For example
   * @dev  LINKTOKEN.transferAndCall(
   * @dev    address(REGISTRY),
   * @dev    amount,
   * @dev    abi.encode(subscriptionId));
   */
  function createSubscription() external nonReentrant onlyAuthorizedUsers returns (uint64) {
    s_currentsubscriptionId++;
    uint64 currentsubscriptionId = s_currentsubscriptionId;
    address[] memory consumers = new address[](0);
    s_subscriptions[currentsubscriptionId] = Subscription({
      balance: 0,
      blockedBalance: 0,
      owner: msg.sender,
      requestedOwner: address(0),
      consumers: consumers
    });

    emit SubscriptionCreated(currentsubscriptionId, msg.sender);
    return currentsubscriptionId;
  }

  /**
   * @notice Request subscription owner transfer.
   * @param subscriptionId - ID of the subscription
   * @param newOwner - proposed new owner of the subscription
   */
  function requestSubscriptionOwnerTransfer(uint64 subscriptionId, address newOwner)
    external
    onlySubOwner(subscriptionId)
    nonReentrant
  {
    // Proposing to address(0) would never be claimable so don't need to check.
    if (s_subscriptions[subscriptionId].requestedOwner != newOwner) {
      s_subscriptions[subscriptionId].requestedOwner = newOwner;
      emit SubscriptionOwnerTransferRequested(subscriptionId, msg.sender, newOwner);
    }
  }

  /**
   * @notice Request subscription owner transfer.
   * @param subscriptionId - ID of the subscription
   * @dev will revert if original owner of subscriptionId has
   * not requested that msg.sender become the new owner.
   */
  function acceptSubscriptionOwnerTransfer(uint64 subscriptionId) external nonReentrant onlyAuthorizedUsers {
    if (s_subscriptions[subscriptionId].owner == address(0)) {
      revert InvalidSubscription();
    }
    if (s_subscriptions[subscriptionId].requestedOwner != msg.sender) {
      revert MustBeRequestedOwner(s_subscriptions[subscriptionId].requestedOwner);
    }
    address oldOwner = s_subscriptions[subscriptionId].owner;
    s_subscriptions[subscriptionId].owner = msg.sender;
    s_subscriptions[subscriptionId].requestedOwner = address(0);
    emit SubscriptionOwnerTransferred(subscriptionId, oldOwner, msg.sender);
  }

  /**
   * @notice Remove a consumer from a Chainlink Functions subscription.
   * @param subscriptionId - ID of the subscription
   * @param consumer - Consumer to remove from the subscription
   */
  function removeConsumer(uint64 subscriptionId, address consumer) external onlySubOwner(subscriptionId) nonReentrant {
    Consumer memory consumerData = s_consumers[consumer][subscriptionId];
    if (consumerData.allowed == false) {
      revert InvalidConsumer(subscriptionId, consumer);
    }
    if (consumerData.initiatedRequests != consumerData.completedRequests) {
      revert ConsumerRequestsInFlight();
    }
    // Note bounded by MAX_CONSUMERS
    address[] memory consumers = s_subscriptions[subscriptionId].consumers;
    uint256 lastConsumerIndex = consumers.length - 1;
    for (uint256 i = 0; i < consumers.length; i++) {
      if (consumers[i] == consumer) {
        address last = consumers[lastConsumerIndex];
        // Storage write to preserve last element
        s_subscriptions[subscriptionId].consumers[i] = last;
        // Storage remove last element
        s_subscriptions[subscriptionId].consumers.pop();
        break;
      }
    }
    delete s_consumers[consumer][subscriptionId];
    emit SubscriptionConsumerRemoved(subscriptionId, consumer);
  }

  /**
   * @notice Add a consumer to a Chainlink Functions subscription.
   * @param subscriptionId - ID of the subscription
   * @param consumer - New consumer which can use the subscription
   */
  function addConsumer(uint64 subscriptionId, address consumer) external onlySubOwner(subscriptionId) nonReentrant {
    // Already maxed, cannot add any more consumers.
    if (s_subscriptions[subscriptionId].consumers.length == MAX_CONSUMERS) {
      revert TooManyConsumers();
    }
    if (s_consumers[consumer][subscriptionId].allowed == true) {
      // Idempotence - do nothing if already added.
      // Ensures uniqueness in s_subscriptions[subscriptionId].consumers.
      return;
    }
    s_consumers[consumer][subscriptionId].allowed = true;
    s_subscriptions[subscriptionId].consumers.push(consumer);

    emit SubscriptionConsumerAdded(subscriptionId, consumer);
  }

  /**
   * @notice Cancel a subscription
   * @param subscriptionId - ID of the subscription
   * @param to - Where to send the remaining LINK to
   */
  function cancelSubscription(uint64 subscriptionId, address to) external onlySubOwner(subscriptionId) nonReentrant {
    if (pendingRequestExists(subscriptionId)) {
      revert PendingRequestExists();
    }
    cancelSubscriptionHelper(subscriptionId, to);
  }

  function cancelSubscriptionHelper(uint64 subscriptionId, address to) private nonReentrant {
    Subscription memory sub = s_subscriptions[subscriptionId];
    uint96 balance = s_subscriptions[subscriptionId].balance;
    // Note bounded by MAX_CONSUMERS;
    // If no consumers, does nothing.
    for (uint256 i = 0; i < sub.consumers.length; i++) {
      delete s_consumers[sub.consumers[i]][subscriptionId];
    }
    delete s_subscriptions[subscriptionId];
    s_totalBalance -= balance;
    if (!LINK.transfer(to, uint256(balance))) {
      revert InsufficientBalance();
    }
    emit SubscriptionCanceled(subscriptionId, to, balance);
  }

  /**
   * @notice Check to see if there exists a request commitment for all consumers for a given sub.
   * @param subscriptionId - ID of the subscription
   * @return true if there exists at least one unfulfilled request for the subscription, false
   * otherwise.
   * @dev Looping is bounded to MAX_CONSUMERS*(number of DONs).
   * @dev Used to disable subscription canceling while outstanding request are present.
   */

  function pendingRequestExists(uint64 subscriptionId) public view returns (bool) {
    address[] memory consumers = s_subscriptions[subscriptionId].consumers;
    for (uint256 i = 0; i < consumers.length; i++) {
      Consumer memory consumer = s_consumers[consumers[i]][subscriptionId];
      if (consumer.initiatedRequests != consumer.completedRequests) {
        return true;
      }
    }
    return false;
  }

  modifier onlySubOwner(uint64 subscriptionId) {
    address owner = s_subscriptions[subscriptionId].owner;
    if (owner == address(0)) {
      revert InvalidSubscription();
    }
    if (msg.sender != owner) {
      revert MustBeSubOwner(owner);
    }
    _;
  }

  /**
   * @dev The allow list is kept on the Router contract. This modifier checks if a user is authorized from there.
   */
  modifier onlyAuthorizedUsers() virtual {
    _;
  }
  modifier nonReentrant() virtual {
    _;
  }
  modifier onlyRouterOwner() virtual {
    _;
  }
  modifier onlyRoute() virtual {
    _;
  }
}
