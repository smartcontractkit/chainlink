// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";
import "../../shared/interfaces/LinkTokenInterface.sol";
import "../../shared/access/ConfirmedOwner.sol";
import "../../interfaces/AggregatorV3Interface.sol";
import "../../shared/interfaces/IERC677Receiver.sol";
import "../interfaces/IVRFSubscriptionV2Plus.sol";

abstract contract SubscriptionAPI is ConfirmedOwner, IERC677Receiver, IVRFSubscriptionV2Plus {
  using EnumerableSet for EnumerableSet.UintSet;

  /// @dev may not be provided upon construction on some chains due to lack of availability
  LinkTokenInterface public LINK;
  /// @dev may not be provided upon construction on some chains due to lack of availability
  AggregatorV3Interface public LINK_NATIVE_FEED;

  // We need to maintain a list of consuming addresses.
  // This bound ensures we are able to loop over them as needed.
  // Should a user require more consumers, they can use multiple subscriptions.
  uint16 public constant MAX_CONSUMERS = 100;
  error TooManyConsumers();
  error InsufficientBalance();
  error InvalidConsumer(uint256 subId, address consumer);
  error InvalidSubscription();
  error OnlyCallableFromLink();
  error InvalidCalldata();
  error MustBeSubOwner(address owner);
  error PendingRequestExists();
  error MustBeRequestedOwner(address proposedOwner);
  error BalanceInvariantViolated(uint256 internalBalance, uint256 externalBalance); // Should never happen
  event FundsRecovered(address to, uint256 amount);
  event NativeFundsRecovered(address to, uint256 amount);
  error LinkAlreadySet();
  error FailedToSendNative();
  error FailedToTransferLink();
  error IndexOutOfRange();
  error LinkNotSet();

  // We use the subscription struct (1 word)
  // at fulfillment time.
  struct Subscription {
    // There are only 1e9*1e18 = 1e27 juels in existence, so the balance can fit in uint96 (2^96 ~ 7e28)
    uint96 balance; // Common link balance used for all consumer requests.
    // a uint96 is large enough to hold around ~8e28 wei, or 80 billion ether.
    // That should be enough to cover most (if not all) subscriptions.
    uint96 nativeBalance; // Common native balance used for all consumer requests.
    uint64 reqCount;
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
  mapping(address => mapping(uint256 => uint64)) /* consumer */ /* subId */ /* nonce */ internal s_consumers;
  mapping(uint256 => SubscriptionConfig) /* subId */ /* subscriptionConfig */ internal s_subscriptionConfigs;
  mapping(uint256 => Subscription) /* subId */ /* subscription */ internal s_subscriptions;
  // subscription nonce used to construct subId. Rises monotonically
  uint64 public s_currentSubNonce;
  // track all subscription id's that were created by this contract
  // note: access should be through the getActiveSubscriptionIds() view function
  // which takes a starting index and a max number to fetch in order to allow
  // "pagination" of the subscription ids. in the event a very large number of
  // subscription id's are stored in this set, they cannot be retrieved in a
  // single RPC call without violating various size limits.
  EnumerableSet.UintSet internal s_subIds;
  // s_totalBalance tracks the total link sent to/from
  // this contract through onTokenTransfer, cancelSubscription and oracleWithdraw.
  // A discrepancy with this contract's link balance indicates someone
  // sent tokens using transfer and so we may need to use recoverFunds.
  uint96 public s_totalBalance;
  // s_totalNativeBalance tracks the total native sent to/from
  // this contract through fundSubscription, cancelSubscription and oracleWithdrawNative.
  // A discrepancy with this contract's native balance indicates someone
  // sent native using transfer and so we may need to use recoverNativeFunds.
  uint96 public s_totalNativeBalance;
  mapping(address => uint96) /* oracle */ /* LINK balance */ internal s_withdrawableTokens;
  mapping(address => uint96) /* oracle */ /* native balance */ internal s_withdrawableNative;

  event SubscriptionCreated(uint256 indexed subId, address owner);
  event SubscriptionFunded(uint256 indexed subId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionFundedWithNative(uint256 indexed subId, uint256 oldNativeBalance, uint256 newNativeBalance);
  event SubscriptionConsumerAdded(uint256 indexed subId, address consumer);
  event SubscriptionConsumerRemoved(uint256 indexed subId, address consumer);
  event SubscriptionCanceled(uint256 indexed subId, address to, uint256 amountLink, uint256 amountNative);
  event SubscriptionOwnerTransferRequested(uint256 indexed subId, address from, address to);
  event SubscriptionOwnerTransferred(uint256 indexed subId, address from, address to);

  struct Config {
    uint16 minimumRequestConfirmations;
    uint32 maxGasLimit;
    // Reentrancy protection.
    bool reentrancyLock;
    // stalenessSeconds is how long before we consider the feed price to be stale
    // and fallback to fallbackWeiPerUnitLink.
    uint32 stalenessSeconds;
    // Gas to cover oracle payment after we calculate the payment.
    // We make it configurable in case those operations are repriced.
    // The recommended number is below, though it may vary slightly
    // if certain chains do not implement certain EIP's.
    // 21000 + // base cost of the transaction
    // 100 + 5000 + // warm subscription balance read and update. See https://eips.ethereum.org/EIPS/eip-2929
    // 2*2100 + 5000 - // cold read oracle address and oracle balance and first time oracle balance update, note first time will be 20k, but 5k subsequently
    // 4800 + // request delete refund (refunds happen after execution), note pre-london fork was 15k. See https://eips.ethereum.org/EIPS/eip-3529
    // 6685 + // Positive static costs of argument encoding etc. note that it varies by +/- x*12 for every x bytes of non-zero data in the proof.
    // Total: 37,185 gas.
    uint32 gasAfterPaymentCalculation;
  }
  Config public s_config;

  error Reentrant();
  modifier nonReentrant() {
    if (s_config.reentrancyLock) {
      revert Reentrant();
    }
    _;
  }

  constructor() ConfirmedOwner(msg.sender) {}

  /**
   * @notice set the LINK token contract and link native feed to be
   * used by this coordinator
   * @param link - address of link token
   * @param linkNativeFeed address of the link native feed
   */
  function setLINKAndLINKNativeFeed(address link, address linkNativeFeed) external onlyOwner {
    // Disallow re-setting link token because the logic wouldn't really make sense
    if (address(LINK) != address(0)) {
      revert LinkAlreadySet();
    }
    LINK = LinkTokenInterface(link);
    LINK_NATIVE_FEED = AggregatorV3Interface(linkNativeFeed);
  }

  /**
   * @notice Owner cancel subscription, sends remaining link directly to the subscription owner.
   * @param subId subscription id
   * @dev notably can be called even if there are pending requests, outstanding ones may fail onchain
   */
  function ownerCancelSubscription(uint256 subId) external onlyOwner {
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
    // If LINK is not set, we cannot recover funds.
    // It is possible that this coordinator address was funded with LINK
    // by accident by a user but the LINK token needs to be set first
    // before we can recover it.
    if (address(LINK) == address(0)) {
      revert LinkNotSet();
    }

    uint256 externalBalance = LINK.balanceOf(address(this));
    uint256 internalBalance = uint256(s_totalBalance);
    if (internalBalance > externalBalance) {
      revert BalanceInvariantViolated(internalBalance, externalBalance);
    }
    if (internalBalance < externalBalance) {
      uint256 amount = externalBalance - internalBalance;
      if (!LINK.transfer(to, amount)) {
        revert FailedToTransferLink();
      }
      emit FundsRecovered(to, amount);
    }
    // If the balances are equal, nothing to be done.
  }

  /**
   * @notice Recover native sent with transfer/call/send instead of fundSubscription.
   * @param to address to send native to
   */
  function recoverNativeFunds(address payable to) external onlyOwner {
    uint256 externalBalance = address(this).balance;
    uint256 internalBalance = uint256(s_totalNativeBalance);
    if (internalBalance > externalBalance) {
      revert BalanceInvariantViolated(internalBalance, externalBalance);
    }
    if (internalBalance < externalBalance) {
      uint256 amount = externalBalance - internalBalance;
      (bool sent, ) = to.call{value: amount}("");
      if (!sent) {
        revert FailedToSendNative();
      }
      emit NativeFundsRecovered(to, amount);
    }
    // If the balances are equal, nothing to be done.
  }

  /*
   * @notice Oracle withdraw LINK earned through fulfilling requests
   * @param recipient where to send the funds
   * @param amount amount to withdraw
   */
  function oracleWithdraw(address recipient, uint96 amount) external nonReentrant {
    if (address(LINK) == address(0)) {
      revert LinkNotSet();
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

  /*
   * @notice Oracle withdraw native earned through fulfilling requests
   * @param recipient where to send the funds
   * @param amount amount to withdraw
   */
  function oracleWithdrawNative(address payable recipient, uint96 amount) external nonReentrant {
    if (s_withdrawableNative[msg.sender] < amount) {
      revert InsufficientBalance();
    }
    // Prevent re-entrancy by updating state before transfer.
    s_withdrawableNative[msg.sender] -= amount;
    s_totalNativeBalance -= amount;
    (bool sent, ) = recipient.call{value: amount}("");
    if (!sent) {
      revert FailedToSendNative();
    }
  }

  function onTokenTransfer(address /* sender */, uint256 amount, bytes calldata data) external override nonReentrant {
    if (msg.sender != address(LINK)) {
      revert OnlyCallableFromLink();
    }
    if (data.length != 32) {
      revert InvalidCalldata();
    }
    uint256 subId = abi.decode(data, (uint256));
    if (s_subscriptionConfigs[subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    // We do not check that the sender is the subscription owner,
    // anyone can fund a subscription.
    uint256 oldBalance = s_subscriptions[subId].balance;
    s_subscriptions[subId].balance += uint96(amount);
    s_totalBalance += uint96(amount);
    emit SubscriptionFunded(subId, oldBalance, oldBalance + amount);
  }

  /**
   * @inheritdoc IVRFSubscriptionV2Plus
   */
  function fundSubscriptionWithNative(uint256 subId) external payable override nonReentrant {
    if (s_subscriptionConfigs[subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    // We do not check that the msg.sender is the subscription owner,
    // anyone can fund a subscription.
    // We also do not check that msg.value > 0, since that's just a no-op
    // and would be a waste of gas on the caller's part.
    uint256 oldNativeBalance = s_subscriptions[subId].nativeBalance;
    s_subscriptions[subId].nativeBalance += uint96(msg.value);
    s_totalNativeBalance += uint96(msg.value);
    emit SubscriptionFundedWithNative(subId, oldNativeBalance, oldNativeBalance + msg.value);
  }

  /**
   * @inheritdoc IVRFSubscriptionV2Plus
   */
  function getSubscription(
    uint256 subId
  )
    public
    view
    override
    returns (uint96 balance, uint96 nativeBalance, uint64 reqCount, address owner, address[] memory consumers)
  {
    if (s_subscriptionConfigs[subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    return (
      s_subscriptions[subId].balance,
      s_subscriptions[subId].nativeBalance,
      s_subscriptions[subId].reqCount,
      s_subscriptionConfigs[subId].owner,
      s_subscriptionConfigs[subId].consumers
    );
  }

  /**
   * @inheritdoc IVRFSubscriptionV2Plus
   */
  function getActiveSubscriptionIds(
    uint256 startIndex,
    uint256 maxCount
  ) external view override returns (uint256[] memory) {
    uint256 numSubs = s_subIds.length();
    if (startIndex >= numSubs) revert IndexOutOfRange();
    uint256 endIndex = startIndex + maxCount;
    endIndex = endIndex > numSubs || maxCount == 0 ? numSubs : endIndex;
    uint256[] memory ids = new uint256[](endIndex - startIndex);
    for (uint256 idx = 0; idx < ids.length; idx++) {
      ids[idx] = s_subIds.at(idx + startIndex);
    }
    return ids;
  }

  /**
   * @inheritdoc IVRFSubscriptionV2Plus
   */
  function createSubscription() external override nonReentrant returns (uint256) {
    // Generate a subscription id that is globally unique.
    uint256 subId = uint256(
      keccak256(abi.encodePacked(msg.sender, blockhash(block.number - 1), address(this), s_currentSubNonce))
    );
    // Increment the subscription nonce counter.
    s_currentSubNonce++;
    // Initialize storage variables.
    address[] memory consumers = new address[](0);
    s_subscriptions[subId] = Subscription({balance: 0, nativeBalance: 0, reqCount: 0});
    s_subscriptionConfigs[subId] = SubscriptionConfig({
      owner: msg.sender,
      requestedOwner: address(0),
      consumers: consumers
    });
    // Update the s_subIds set, which tracks all subscription ids created in this contract.
    s_subIds.add(subId);

    emit SubscriptionCreated(subId, msg.sender);
    return subId;
  }

  /**
   * @inheritdoc IVRFSubscriptionV2Plus
   */
  function requestSubscriptionOwnerTransfer(
    uint256 subId,
    address newOwner
  ) external override onlySubOwner(subId) nonReentrant {
    // Proposing to address(0) would never be claimable so don't need to check.
    if (s_subscriptionConfigs[subId].requestedOwner != newOwner) {
      s_subscriptionConfigs[subId].requestedOwner = newOwner;
      emit SubscriptionOwnerTransferRequested(subId, msg.sender, newOwner);
    }
  }

  /**
   * @inheritdoc IVRFSubscriptionV2Plus
   */
  function acceptSubscriptionOwnerTransfer(uint256 subId) external override nonReentrant {
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
   * @inheritdoc IVRFSubscriptionV2Plus
   */
  function addConsumer(uint256 subId, address consumer) external override onlySubOwner(subId) nonReentrant {
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

  function deleteSubscription(uint256 subId) internal returns (uint96 balance, uint96 nativeBalance) {
    SubscriptionConfig memory subConfig = s_subscriptionConfigs[subId];
    Subscription memory sub = s_subscriptions[subId];
    balance = sub.balance;
    nativeBalance = sub.nativeBalance;
    // Note bounded by MAX_CONSUMERS;
    // If no consumers, does nothing.
    for (uint256 i = 0; i < subConfig.consumers.length; i++) {
      delete s_consumers[subConfig.consumers[i]][subId];
    }
    delete s_subscriptionConfigs[subId];
    delete s_subscriptions[subId];
    s_subIds.remove(subId);
    s_totalBalance -= balance;
    s_totalNativeBalance -= nativeBalance;
    return (balance, nativeBalance);
  }

  function cancelSubscriptionHelper(uint256 subId, address to) internal {
    (uint96 balance, uint96 nativeBalance) = deleteSubscription(subId);

    // Only withdraw LINK if the token is active and there is a balance.
    if (address(LINK) != address(0) && balance != 0) {
      if (!LINK.transfer(to, uint256(balance))) {
        revert InsufficientBalance();
      }
    }

    // send native to the "to" address using call
    (bool success, ) = to.call{value: uint256(nativeBalance)}("");
    if (!success) {
      revert FailedToSendNative();
    }
    emit SubscriptionCanceled(subId, to, balance, nativeBalance);
  }

  modifier onlySubOwner(uint256 subId) {
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
