// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/LinkTokenInterface.sol";
import "../ConfirmedOwner.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "../interfaces/ERC677ReceiverInterface.sol";

contract SubscriptionAPI is ConfirmedOwner, ReentrancyGuard, ERC677ReceiverInterface {
  /// @dev may not be provided upon construction on some chains due to lack of availability
  LinkTokenInterface public LINK;

  // We need to maintain a list of consuming addresses.
  // This bound ensures we are able to loop over them as needed.
  // Should a user require more consumers, they can use multiple subscriptions.
  uint16 public constant MAX_CONSUMERS = 100;
  error TooManyConsumers();
  error InsufficientBalance();
  error InvalidConsumer(uint64 subId, address consumer);
  error InvalidSubscription();
  error OnlyCallableFromLink();
  error InvalidCalldata();
  error MustBeSubOwner(address owner);
  error PendingRequestExists();
  error MustBeRequestedOwner(address proposedOwner);
  error BalanceInvariantViolated(uint256 internalBalance, uint256 externalBalance); // Should never happen
  event FundsRecovered(address to, uint256 amount);
  event EthFundsRecovered(address to, uint256 amount);
  error LinkAlreadySet();
  error FailedToSendEther();

  // We use the subscription struct (1 word)
  // at fulfillment time.
  struct Subscription {
    // There are only 1e9*1e18 = 1e27 juels in existence, so the balance can fit in uint96 (2^96 ~ 7e28)
    uint96 balance; // Common link balance used for all consumer requests.
    uint96 ethBalance; // Common eth balance used for all consumer requests.
    uint64 reqCount; // For fee tiers
  }
  // We use the config for the mgmt APIs
  struct SubscriptionConfig {
    address owner; // Owner can fund/withdraw/cancel the sub.
    address requestedOwner; // For safely transferring sub ownership.
    // Maintains the list of keys in s_consumers.
    // We do this for 2 reasons:
    // 1. To be able to clean up all keys from s_consumers when canceling a subscription.
    // 2. To be able to return the list of all consumers in getSubscription.
    // Note that we need the s_consumers map to be able to directly check if a
    // consumer is valid without reading all the consumers from storage.
    address[] consumers;
  }
  // Note a nonce of 0 indicates an the consumer is not assigned to that subscription.
  mapping(address => mapping(uint64 => uint64)) /* consumer */ /* subId */ /* nonce */ internal s_consumers;
  mapping(uint64 => SubscriptionConfig) /* subId */ /* subscriptionConfig */ internal s_subscriptionConfigs;
  mapping(uint64 => Subscription) /* subId */ /* subscription */ internal s_subscriptions;
  // We make the sub count public so that its possible to
  // get all the current subscriptions via getSubscription.
  uint64 public s_currentSubId;
  // s_totalBalance tracks the total link sent to/from
  // this contract through onTokenTransfer, cancelSubscription and oracleWithdraw.
  // A discrepancy with this contract's link balance indicates someone
  // sent tokens using transfer and so we may need to use recoverFunds.
  uint96 public s_totalBalance;
  // s_totalEthBalance tracks the total eth sent to/from
  // this contract through fundSubscription, cancelSubscription and oracleWithdrawEth.
  // A discrepancy with this contract's eth balance indicates someone
  // sent eth using transfer and so we may need to use recoverEthFunds.
  uint96 public s_totalEthBalance;
  mapping(address => uint96) /* oracle */ /* LINK balance */ internal s_withdrawableTokens;
  mapping(address => uint96) /* oracle */ /* ETH balance */ internal s_withdrawableEth;

  event SubscriptionCreated(uint64 indexed subId, address owner);
  event SubscriptionFunded(uint64 indexed subId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionFundedWithEth(uint64 indexed subId, uint256 oldEthBalance, uint256 newEthBalance);
  event SubscriptionConsumerAdded(uint64 indexed subId, address consumer);
  event SubscriptionConsumerRemoved(uint64 indexed subId, address consumer);
  event SubscriptionCanceled(uint64 indexed subId, address to, uint256 amountLink, uint256 amountEth);
  event SubscriptionOwnerTransferRequested(uint64 indexed subId, address from, address to);
  event SubscriptionOwnerTransferred(uint64 indexed subId, address from, address to);

  constructor() ConfirmedOwner(msg.sender) {
  }

  function setLINK(address link) external onlyOwner {
    // Disallow re-setting link token because the logic wouldn't really make sense
    if (address(LINK) != address(0)) {
      revert LinkAlreadySet();
    }
    LINK = LinkTokenInterface(link);
  }

  /**
   * @notice Owner cancel subscription, sends remaining link directly to the subscription owner.
   * @param subId subscription id
   * @dev notably can be called even if there are pending requests, outstanding ones may fail onchain
   */
  function ownerCancelSubscription(uint64 subId) external onlyOwner {
    if (s_subscriptionConfigs[subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    cancelSubscriptionHelper(subId, s_subscriptionConfigs[subId].owner);
  }

  /**
   * @notice Recover link sent with transfer instead of transferAndCall.
   * @param to address to send link to
   */
  function recoverFunds(address to) external onlyOwner {
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

  /**
    * @notice Recover eth sent with transfer/call/send instead of fundSubscription.
    * @param to address to send eth to
  */
  function recoverEthFunds(address payable to) external onlyOwner {
    uint256 externalBalance = address(this).balance;
    uint256 internalBalance = uint256(s_totalEthBalance);
    if (internalBalance > externalBalance) {
      revert BalanceInvariantViolated(internalBalance, externalBalance);
    }
    if (internalBalance < externalBalance) {
      uint256 amount = externalBalance - internalBalance;
      (bool sent, ) = to.call{value: amount}("");
      if (!sent) {
        revert FailedToSendEther();
      }
      emit EthFundsRecovered(to, amount);
    }
    // If the balances are equal, nothing to be done.
  }

  /*
   * @notice Oracle withdraw LINK earned through fulfilling requests
   * @param recipient where to send the funds
   * @param amount amount to withdraw
   */
  function oracleWithdraw(address recipient, uint96 amount) external nonReentrant {
    if (s_withdrawableTokens[msg.sender] < amount) {
      revert InsufficientBalance();
    }
    s_withdrawableTokens[msg.sender] -= amount;
    s_totalBalance -= amount;
    if (!LINK.transfer(recipient, amount)) {
      revert InsufficientBalance();
    }
  }

  function oracleWithdrawEth(address payable recipient, uint96 amount) external nonReentrant {
    if (s_withdrawableEth[msg.sender] < amount) {
      revert InsufficientBalance();
    }
    // Prevent re-entrancy by updating state before transfer.
    s_withdrawableEth[msg.sender] -= amount;
    s_totalEthBalance -= amount;
    (bool sent, ) = recipient.call{value: amount}("");
    if (!sent) {
      revert FailedToSendEther();
    }
  }

  function onTokenTransfer(address /* sender */, uint256 amount, bytes calldata data) external override nonReentrant {
    if (msg.sender != address(LINK)) {
      revert OnlyCallableFromLink();
    }
    if (data.length != 32) {
      revert InvalidCalldata();
    }
    uint64 subId = abi.decode(data, (uint64));
    if (s_subscriptionConfigs[subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    // We do not check that the msg.sender is the subscription owner,
    // anyone can fund a subscription.
    uint256 oldBalance = s_subscriptions[subId].balance;
    s_subscriptions[subId].balance += uint96(amount);
    s_totalBalance += uint96(amount);
    emit SubscriptionFunded(subId, oldBalance, oldBalance + amount);
  }

  /**
    * @notice Fund a subscription with ETH.
    * @param subId - ID of the subscription
  */
  function fundSubscriptionWithEth(uint64 subId) external payable nonReentrant {
    if (s_subscriptionConfigs[subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    // We do not check that the msg.sender is the subscription owner,
    // anyone can fund a subscription.
    uint256 oldEthBalance = s_subscriptions[subId].ethBalance;
    s_subscriptions[subId].ethBalance += uint96(msg.value);
    s_totalEthBalance += uint96(msg.value);
    emit SubscriptionFundedWithEth(subId, oldEthBalance, oldEthBalance + msg.value);
  }

  /**
   * @notice Get a VRF subscription.
   * @param subId - ID of the subscription
   * @return balance - LINK balance of the subscription in juels.
   * @return ethBalance - ETH balance of the subscription in wei.
   * @return reqCount - number of requests for this subscription, determines fee tier.
   * @return owner - owner of the subscription.
   * @return consumers - list of consumer address which are able to use this subscription.
   */
  function getSubscription(
    uint64 subId
  ) external view returns (uint96 balance, uint96 ethBalance, uint64 reqCount, address owner, address[] memory consumers) {
    if (s_subscriptionConfigs[subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    return (
      s_subscriptions[subId].balance,
      s_subscriptions[subId].ethBalance,
      s_subscriptions[subId].reqCount,
      s_subscriptionConfigs[subId].owner,
      s_subscriptionConfigs[subId].consumers
    );
  }

  /**
   * @notice Create a VRF subscription.
   * @return subId - A unique subscription id.
   * @dev You can manage the consumer set dynamically with addConsumer/removeConsumer.
   * @dev Note to fund the subscription with LINK, use transferAndCall. For example
   * @dev  LINKTOKEN.transferAndCall(
   * @dev    address(COORDINATOR),
   * @dev    amount,
   * @dev    abi.encode(subId));
   */
  function createSubscription() external nonReentrant returns (uint64) {
    s_currentSubId++;
    uint64 currentSubId = s_currentSubId;
    address[] memory consumers = new address[](0);
    s_subscriptions[currentSubId] = Subscription({balance: 0, ethBalance : 0, reqCount: 0});
    s_subscriptionConfigs[currentSubId] = SubscriptionConfig({
      owner: msg.sender,
      requestedOwner: address(0),
      consumers: consumers
    });

    emit SubscriptionCreated(currentSubId, msg.sender);
    return currentSubId;
  }

  /**
   * @notice Request subscription owner transfer.
   * @param subId - ID of the subscription
   * @param newOwner - proposed new owner of the subscription
   */
  function requestSubscriptionOwnerTransfer(
    uint64 subId,
    address newOwner
  ) external onlySubOwner(subId) nonReentrant {
    // Proposing to address(0) would never be claimable so don't need to check.
    if (s_subscriptionConfigs[subId].requestedOwner != newOwner) {
      s_subscriptionConfigs[subId].requestedOwner = newOwner;
      emit SubscriptionOwnerTransferRequested(subId, msg.sender, newOwner);
    }
  }

  /**
   * @notice Request subscription owner transfer.
   * @param subId - ID of the subscription
   * @dev will revert if original owner of subId has
   * not requested that msg.sender become the new owner.
   */
  function acceptSubscriptionOwnerTransfer(uint64 subId) external nonReentrant {
    if (s_subscriptionConfigs[subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    if (s_subscriptionConfigs[subId].requestedOwner != msg.sender) {
      revert MustBeRequestedOwner(s_subscriptionConfigs[subId].requestedOwner);
    }
    address oldOwner = s_subscriptionConfigs[subId].owner;
    s_subscriptionConfigs[subId].owner = msg.sender;
    s_subscriptionConfigs[subId].requestedOwner = address(0);
    emit SubscriptionOwnerTransferred(subId, oldOwner, msg.sender);
  }

  /**
   * @notice Add a consumer to a VRF subscription.
   * @param subId - ID of the subscription
   * @param consumer - New consumer which can use the subscription
   */
  function addConsumer(uint64 subId, address consumer) external onlySubOwner(subId) nonReentrant {
    // Already maxed, cannot add any more consumers.
    if (s_subscriptionConfigs[subId].consumers.length == MAX_CONSUMERS) {
      revert TooManyConsumers();
    }
    if (s_consumers[consumer][subId] != 0) {
      // Idempotence - do nothing if already added.
      // Ensures uniqueness in s_subscriptions[subId].consumers.
      return;
    }
    // Initialize the nonce to 1, indicating the consumer is allocated.
    s_consumers[consumer][subId] = 1;
    s_subscriptionConfigs[subId].consumers.push(consumer);

    emit SubscriptionConsumerAdded(subId, consumer);
  }

  function cancelSubscriptionHelper(uint64 subId, address to) internal nonReentrant {
    SubscriptionConfig memory subConfig = s_subscriptionConfigs[subId];
    Subscription memory sub = s_subscriptions[subId];
    uint96 balance = sub.balance;
    uint96 ethBalance = sub.ethBalance;
    // Note bounded by MAX_CONSUMERS;
    // If no consumers, does nothing.
    for (uint256 i = 0; i < subConfig.consumers.length; i++) {
      delete s_consumers[subConfig.consumers[i]][subId];
    }
    delete s_subscriptionConfigs[subId];
    delete s_subscriptions[subId];
    s_totalBalance -= balance;
    s_totalEthBalance -= ethBalance;
    if (!LINK.transfer(to, uint256(balance))) {
      revert InsufficientBalance();
    }
    // send eth to the "to" address using call
    (bool success, ) = to.call{value: uint256(ethBalance)}("");
    if (!success) {
      revert FailedToSendEther();
    }
    emit SubscriptionCanceled(subId, to, balance, ethBalance);
  }

  modifier onlySubOwner(uint64 subId) {
    address owner = s_subscriptionConfigs[subId].owner;
    if (owner == address(0)) {
      revert InvalidSubscription();
    }
    if (msg.sender != owner) {
      revert MustBeSubOwner(owner);
    }
    _;
  }
}
