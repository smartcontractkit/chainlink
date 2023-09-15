// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {LinkTokenInterface} from "../../../../../shared/interfaces/LinkTokenInterface.sol";
import {AggregatorV3Interface} from "../../../../../interfaces/AggregatorV3Interface.sol";
import {FunctionsBillingRegistryInterface} from "./FunctionsBillingRegistryInterface.sol";
import {FunctionsOracleInterface} from "./FunctionsOracleInterface.sol";
import {FunctionsClientInterface} from "./FunctionsClientInterface.sol";
import {TypeAndVersionInterface} from "../../../../../interfaces/TypeAndVersionInterface.sol";
import {IERC677Receiver} from "../../../../../shared/interfaces/IERC677Receiver.sol";
import {AuthorizedOriginReceiverInterface} from "./AuthorizedOriginReceiverInterface.sol";
import {ConfirmedOwnerUpgradeable} from "./ConfirmedOwnerUpgradeable.sol";
import {AuthorizedReceiver} from "./AuthorizedReceiver.sol";
import {SafeCast} from "../../../../../vendor/openzeppelin-solidity/v4.8.0/contracts/utils/math/SafeCast.sol";
import {PausableUpgradeable} from "../../../../../vendor/openzeppelin-contracts-upgradeable/v4.8.1/security/PausableUpgradeable.sol";
import {Initializable} from "../../../../../vendor/openzeppelin-contracts-upgradeable/v4.8.1/proxy/utils/Initializable.sol";

/**
 * @title Functions Billing Registry contract
 * @notice Contract that coordinates payment from users to the nodes of the Decentralized Oracle Network (DON).
 * @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
 */
contract FunctionsBillingRegistryMigration is
  Initializable,
  ConfirmedOwnerUpgradeable,
  PausableUpgradeable,
  FunctionsBillingRegistryInterface,
  IERC677Receiver,
  AuthorizedReceiver
{
  LinkTokenInterface public LINK;
  AggregatorV3Interface public LINK_ETH_FEED;
  AuthorizedOriginReceiverInterface private ORACLE_WITH_ALLOWLIST;

  // We need to maintain a list of consuming addresses.
  // This bound ensures we are able to loop over them as needed.
  // Should a user require more consumers, they can use multiple subscriptions.
  uint16 public constant MAX_CONSUMERS = 100;

  error TooManyConsumers();
  error InsufficientBalance();
  error InvalidConsumer(uint64 subscriptionId, address consumer);
  error InvalidSubscription();
  error OnlyCallableFromLink();
  error InvalidCalldata();
  error MustBeSubOwner(address owner);
  error PendingRequestExists();
  error MustBeRequestedOwner(address proposedOwner);
  error BalanceInvariantViolated(uint256 internalBalance, uint256 externalBalance); // Should never happen
  event FundsRecovered(address to, uint256 amount);

  struct Subscription {
    // There are only 1e9*1e18 = 1e27 juels in existence, so the balance can fit in uint96 (2^96 ~ 7e28)
    uint96 balance; // Common LINK balance that is controlled by the Registry to be used for all consumer requests.
    uint96 blockedBalance; // LINK balance that is reserved to pay for pending consumer requests.
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
  mapping(address => mapping(uint64 => uint64)) /* consumer */ /* subscriptionId */ /* nonce */ private s_consumers;
  mapping(uint64 => SubscriptionConfig) /* subscriptionId */ /* subscriptionConfig */ private s_subscriptionConfigs;
  mapping(uint64 => Subscription) /* subscriptionId */ /* subscription */ private s_subscriptions;
  // We make the sub count public so that its possible to
  // get all the current subscriptions via getSubscription.
  uint64 private s_currentsubscriptionId;
  // s_totalBalance tracks the total link sent to/from
  // this contract through onTokenTransfer, cancelSubscription and oracleWithdraw.
  // A discrepancy with this contract's link balance indicates someone
  // sent tokens using transfer and so we may need to use recoverFunds.
  uint96 private s_totalBalance;
  event SubscriptionCreated(uint64 indexed subscriptionId, address owner);
  event SubscriptionFunded(uint64 indexed subscriptionId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionConsumerAdded(uint64 indexed subscriptionId, address consumer);
  event SubscriptionConsumerRemoved(uint64 indexed subscriptionId, address consumer);
  event SubscriptionCanceled(uint64 indexed subscriptionId, address to, uint256 amount);
  event SubscriptionOwnerTransferRequested(uint64 indexed subscriptionId, address from, address to);
  event SubscriptionOwnerTransferred(uint64 indexed subscriptionId, address from, address to);

  error GasLimitTooBig(uint32 have, uint32 want);
  error InvalidLinkWeiPrice(int256 linkWei);
  error IncorrectRequestID();
  error PaymentTooLarge();
  error Reentrant();

  mapping(address => uint96) /* oracle node */ /* LINK balance */ private s_withdrawableTokens;
  struct Commitment {
    uint64 subscriptionId;
    address client;
    uint32 gasLimit;
    uint256 gasPrice;
    address don;
    uint96 donFee;
    uint96 registryFee;
    uint96 estimatedCost;
    uint256 timestamp;
  }
  mapping(bytes32 => Commitment) /* requestID */ /* Commitment */ private s_requestCommitments;
  event BillingStart(bytes32 indexed requestId, Commitment commitment);
  struct ItemizedBill {
    uint96 signerPayment;
    uint96 transmitterPayment;
    uint96 totalCost;
  }
  event BillingEnd(
    bytes32 indexed requestId,
    uint64 subscriptionId,
    uint96 signerPayment,
    uint96 transmitterPayment,
    uint96 totalCost,
    bool success
  );
  event RequestTimedOut(bytes32 indexed requestId);

  struct Config {
    // Maxiumum amount of gas that can be given to a request's client callback
    uint32 maxGasLimit;
    // Reentrancy protection.
    bool reentrancyLock;
    // stalenessSeconds is how long before we consider the feed price to be stale
    // and fallback to fallbackWeiPerUnitLink.
    uint32 stalenessSeconds;
    // Gas to cover transmitter oracle payment after we calculate the payment.
    // We make it configurable in case those operations are repriced.
    uint256 gasAfterPaymentCalculation;
    // Represents the average gas execution cost. Used in estimating cost beforehand.
    uint32 gasOverhead;
    // how many seconds it takes before we consider a request to be timed out
    uint32 requestTimeoutSeconds;
  }
  int256 private s_fallbackWeiPerUnitLink;
  Config private s_config;
  event ConfigSet(
    uint32 maxGasLimit,
    uint32 stalenessSeconds,
    uint256 gasAfterPaymentCalculation,
    int256 fallbackWeiPerUnitLink,
    uint32 gasOverhead
  );

  /**
   * @dev Initializes the contract.
   */
  function initialize(address link, address linkEthFeed, address oracle) public initializer {
    __Pausable_init();
    __ConfirmedOwner_initialize(msg.sender, address(0));
    LINK = LinkTokenInterface(link);
    LINK_ETH_FEED = AggregatorV3Interface(linkEthFeed);
    ORACLE_WITH_ALLOWLIST = AuthorizedOriginReceiverInterface(oracle);
  }

  /**
   * @notice Sets the configuration of the Chainlink Functions billing registry
   * @param maxGasLimit global max for request gas limit
   * @param stalenessSeconds if the eth/link feed is more stale then this, use the fallback price
   * @param gasAfterPaymentCalculation gas used in doing accounting after completing the gas measurement
   * @param fallbackWeiPerUnitLink fallback eth/link price in the case of a stale feed
   * @param gasOverhead average gas execution cost used in estimating total cost
   * @param requestTimeoutSeconds e2e timeout after which user won't be charged
   */
  function setConfig(
    uint32 maxGasLimit,
    uint32 stalenessSeconds,
    uint256 gasAfterPaymentCalculation,
    int256 fallbackWeiPerUnitLink,
    uint32 gasOverhead,
    uint32 requestTimeoutSeconds
  ) external onlyOwner {
    if (fallbackWeiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(fallbackWeiPerUnitLink);
    }
    s_config = Config({
      maxGasLimit: maxGasLimit,
      stalenessSeconds: stalenessSeconds,
      gasAfterPaymentCalculation: gasAfterPaymentCalculation,
      reentrancyLock: false,
      gasOverhead: gasOverhead,
      requestTimeoutSeconds: requestTimeoutSeconds
    });
    s_fallbackWeiPerUnitLink = fallbackWeiPerUnitLink;
    emit ConfigSet(maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, gasOverhead);
  }

  /**
   * @notice Gets the configuration of the Chainlink Functions billing registry
   * @return maxGasLimit global max for request gas limit
   * @return stalenessSeconds if the eth/link feed is more stale then this, use the fallback price
   * @return gasAfterPaymentCalculation gas used in doing accounting after completing the gas measurement
   * @return fallbackWeiPerUnitLink fallback eth/link price in the case of a stale feed
   * @return gasOverhead average gas execution cost used in estimating total cost
   */
  function getConfig()
    external
    view
    returns (
      uint32 maxGasLimit,
      uint32 stalenessSeconds,
      uint256 gasAfterPaymentCalculation,
      int256 fallbackWeiPerUnitLink,
      uint32 gasOverhead
    )
  {
    return (
      s_config.maxGasLimit,
      s_config.stalenessSeconds,
      s_config.gasAfterPaymentCalculation,
      s_fallbackWeiPerUnitLink,
      s_config.gasOverhead
    );
  }

  function pause() external onlyOwner {
    _pause();
  }

  function unpause() external onlyOwner {
    _unpause();
  }

  function getTotalBalance() external view returns (uint256) {
    return s_totalBalance;
  }

  /**
   * @notice Owner cancel subscription, sends remaining link directly to the subscription owner.
   * @param subscriptionId subscription id
   * @dev notably can be called even if there are pending requests, outstanding ones may fail onchain
   */
  function ownerCancelSubscription(uint64 subscriptionId) external onlyOwner {
    address owner = s_subscriptionConfigs[subscriptionId].owner;
    if (owner == address(0)) {
      revert InvalidSubscription();
    }
    cancelSubscriptionHelper(subscriptionId, owner);
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
   * @inheritdoc FunctionsBillingRegistryInterface
   */
  function getRequestConfig() external view override returns (uint32, address[] memory) {
    return (s_config.maxGasLimit, getAuthorizedSenders());
  }

  /**
   * @inheritdoc FunctionsBillingRegistryInterface
   */
  function getRequiredFee(
    bytes calldata /* data */,
    FunctionsBillingRegistryInterface.RequestBilling memory /* billing */
  ) public pure override returns (uint96) {
    return 1;
  }

  /**
   * @inheritdoc FunctionsBillingRegistryInterface
   */
  function estimateCost(
    uint32 gasLimit,
    uint256 gasPrice,
    uint96 donFee,
    uint96 registryFee
  ) public view override returns (uint96) {
    int256 weiPerUnitLink;
    weiPerUnitLink = getFeedData();
    if (weiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(weiPerUnitLink);
    }
    uint256 executionGas = s_config.gasOverhead + s_config.gasAfterPaymentCalculation + gasLimit;
    // (1e18 juels/link) (wei/gas * gas) / (wei/link) = juels
    uint256 paymentNoFee = (1e18 * gasPrice * executionGas) / uint256(weiPerUnitLink);
    uint256 fee = uint256(donFee) + uint256(registryFee);
    if (paymentNoFee > (1e27 - fee)) {
      revert PaymentTooLarge(); // Payment + fee cannot be more than all of the link in existence.
    }
    return uint96(paymentNoFee + fee);
  }

  /**
   * @inheritdoc FunctionsBillingRegistryInterface
   */
  function startBilling(
    bytes calldata data,
    RequestBilling calldata billing
  ) external override validateAuthorizedSender nonReentrant whenNotPaused returns (bytes32) {
    // Input validation using the subscription storage.
    if (s_subscriptionConfigs[billing.subscriptionId].owner == address(0)) {
      revert InvalidSubscription();
    }
    // It's important to ensure that the consumer is in fact who they say they
    // are, otherwise they could use someone else's subscription balance.
    // A nonce of 0 indicates consumer is not allocated to the sub.
    uint64 currentNonce = s_consumers[billing.client][billing.subscriptionId];
    if (currentNonce == 0) {
      revert InvalidConsumer(billing.subscriptionId, billing.client);
    }
    // No lower bound on the requested gas limit. A user could request 0
    // and they would simply be billed for the gas and computation.
    if (billing.gasLimit > s_config.maxGasLimit) {
      revert GasLimitTooBig(billing.gasLimit, s_config.maxGasLimit);
    }

    // Check that subscription can afford the estimated cost
    uint96 oracleFee = FunctionsOracleInterface(msg.sender).getRequiredFee(data, billing);
    uint96 registryFee = getRequiredFee(data, billing);
    uint96 estimatedCost = estimateCost(billing.gasLimit, billing.gasPrice, oracleFee, registryFee);
    uint96 effectiveBalance = s_subscriptions[billing.subscriptionId].balance -
      s_subscriptions[billing.subscriptionId].blockedBalance;
    if (effectiveBalance < estimatedCost) {
      revert InsufficientBalance();
    }

    uint64 nonce = currentNonce + 1;
    bytes32 requestId = computeRequestId(msg.sender, billing.client, billing.subscriptionId, nonce);

    Commitment memory commitment = Commitment(
      billing.subscriptionId,
      billing.client,
      billing.gasLimit,
      billing.gasPrice,
      msg.sender,
      oracleFee,
      registryFee,
      estimatedCost,
      block.timestamp
    );
    s_requestCommitments[requestId] = commitment;
    s_subscriptions[billing.subscriptionId].blockedBalance += estimatedCost;

    emit BillingStart(requestId, commitment);
    s_consumers[billing.client][billing.subscriptionId] = nonce;
    return requestId;
  }

  function computeRequestId(
    address don,
    address client,
    uint64 subscriptionId,
    uint64 nonce
  ) private pure returns (bytes32) {
    return keccak256(abi.encode(don, client, subscriptionId, nonce));
  }

  /**
   * @dev calls target address with exactly gasAmount gas and data as calldata
   * or reverts if at least gasAmount gas is not available.
   */
  function callWithExactGas(uint256 gasAmount, address target, bytes memory data) private returns (bool success) {
    // solhint-disable-next-line no-inline-assembly
    assembly {
      let g := gas()
      // GAS_FOR_CALL_EXACT_CHECK = 5000
      // Compute g -= GAS_FOR_CALL_EXACT_CHECK and check for underflow
      // The gas actually passed to the callee is min(gasAmount, 63//64*gas available).
      // We want to ensure that we revert if gasAmount >  63//64*gas available
      // as we do not want to provide them with less, however that check itself costs
      // gas.  GAS_FOR_CALL_EXACT_CHECK ensures we have at least enough gas to be able
      // to revert if gasAmount >  63//64*gas available.
      if lt(g, 5000) {
        revert(0, 0)
      }
      g := sub(g, 5000)
      // if g - g//64 <= gasAmount, revert
      // (we subtract g//64 because of EIP-150)
      if iszero(gt(sub(g, div(g, 64)), gasAmount)) {
        revert(0, 0)
      }
      // solidity calls check that a contract actually exists at the destination, so we do the same
      if iszero(extcodesize(target)) {
        revert(0, 0)
      }
      // call and return whether we succeeded. ignore return data
      // call(gas,addr,value,argsOffset,argsLength,retOffset,retLength)
      success := call(gasAmount, target, 0, add(data, 0x20), mload(data), 0, 0)
    }
    return success;
  }

  /**
   * @inheritdoc FunctionsBillingRegistryInterface
   */
  function fulfillAndBill(
    bytes32 requestId,
    bytes calldata response,
    bytes calldata err,
    address transmitter,
    address[31] memory signers,
    uint8 signerCount,
    uint256 reportValidationGas,
    uint256 initialGas
  ) external override validateAuthorizedSender nonReentrant whenNotPaused returns (bool success) {
    Commitment memory commitment = s_requestCommitments[requestId];
    if (commitment.don == address(0)) {
      revert IncorrectRequestID();
    }
    delete s_requestCommitments[requestId];

    bytes memory callback = abi.encodeWithSelector(
      FunctionsClientInterface.handleOracleFulfillment.selector,
      requestId,
      response,
      err
    );
    // Call with explicitly the amount of callback gas requested
    // Important to not let them exhaust the gas budget and avoid payment.
    // Do not allow any non-view/non-pure coordinator functions to be called
    // during the consumers callback code via reentrancyLock.
    // NOTE: that callWithExactGas will revert if we do not have sufficient gas
    // to give the callee their requested amount.
    s_config.reentrancyLock = true;
    success = callWithExactGas(commitment.gasLimit, commitment.client, callback);
    s_config.reentrancyLock = false;

    // We want to charge users exactly for how much gas they use in their callback.
    // The gasAfterPaymentCalculation is meant to cover these additional operations where we
    // decrement the subscription balance and increment the oracle's withdrawable balance.
    ItemizedBill memory bill = calculatePaymentAmount(
      initialGas,
      s_config.gasAfterPaymentCalculation,
      commitment.donFee,
      signerCount,
      commitment.registryFee,
      reportValidationGas,
      tx.gasprice
    );
    if (s_subscriptions[commitment.subscriptionId].balance < bill.totalCost) {
      revert InsufficientBalance();
    }
    s_subscriptions[commitment.subscriptionId].balance -= bill.totalCost;
    // Pay out signers their portion of the DON fee
    for (uint256 i = 0; i < signerCount; i++) {
      if (signers[i] != transmitter) {
        s_withdrawableTokens[signers[i]] += bill.signerPayment;
      }
    }
    // Pay out the registry fee
    s_withdrawableTokens[owner()] += commitment.registryFee;
    // Reimburse the transmitter for the execution gas cost + pay them their portion of the DON fee
    s_withdrawableTokens[transmitter] += bill.transmitterPayment;
    // Remove blocked balance
    s_subscriptions[commitment.subscriptionId].blockedBalance -= commitment.estimatedCost;
    // Include payment in the event for tracking costs.
    emit BillingEnd(
      requestId,
      commitment.subscriptionId,
      bill.signerPayment,
      bill.transmitterPayment,
      bill.totalCost,
      success
    );
  }

  // Determine the cost breakdown for payment
  function calculatePaymentAmount(
    uint256 startGas,
    uint256 gasAfterPaymentCalculation,
    uint96 donFee,
    uint8 signerCount,
    uint96 registryFee,
    uint256 reportValidationGas,
    uint256 weiPerUnitGas
  ) private view returns (ItemizedBill memory) {
    int256 weiPerUnitLink;
    weiPerUnitLink = getFeedData();
    if (weiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(weiPerUnitLink);
    }
    // (1e18 juels/link) (wei/gas * gas) / (wei/link) = juels
    uint256 paymentNoFee = (1e18 *
      weiPerUnitGas *
      (reportValidationGas + gasAfterPaymentCalculation + startGas - gasleft())) / uint256(weiPerUnitLink);
    uint256 fee = uint256(donFee) + uint256(registryFee);
    if (paymentNoFee > (1e27 - fee)) {
      revert PaymentTooLarge(); // Payment + fee cannot be more than all of the link in existence.
    }
    uint96 signerPayment = donFee / uint96(signerCount);
    uint96 transmitterPayment = uint96(paymentNoFee) + signerPayment;
    uint96 totalCost = SafeCast.toUint96(paymentNoFee + fee);
    return ItemizedBill(signerPayment, transmitterPayment, totalCost);
  }

  function getFeedData() private view returns (int256) {
    uint32 stalenessSeconds = s_config.stalenessSeconds;
    bool staleFallback = stalenessSeconds > 0;
    (, int256 weiPerUnitLink, , uint256 timestamp, ) = LINK_ETH_FEED.latestRoundData();
    // solhint-disable-next-line not-rely-on-time
    if (staleFallback && stalenessSeconds < block.timestamp - timestamp) {
      weiPerUnitLink = s_fallbackWeiPerUnitLink;
    }
    return weiPerUnitLink;
  }

  /*
   * @notice Oracle withdraw LINK earned through fulfilling requests
   * @notice If amount is 0 the full balance will be withdrawn
   * @param recipient where to send the funds
   * @param amount amount to withdraw
   */
  function oracleWithdraw(address recipient, uint96 amount) external nonReentrant whenNotPaused {
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

  function onTokenTransfer(
    address /* sender */,
    uint256 amount,
    bytes calldata data
  ) external override nonReentrant whenNotPaused {
    if (msg.sender != address(LINK)) {
      revert OnlyCallableFromLink();
    }
    if (data.length != 32) {
      revert InvalidCalldata();
    }
    uint64 subscriptionId = abi.decode(data, (uint64));
    if (s_subscriptionConfigs[subscriptionId].owner == address(0)) {
      revert InvalidSubscription();
    }
    // We do not check that the msg.sender is the subscription owner,
    // anyone can fund a subscription.
    uint256 oldBalance = s_subscriptions[subscriptionId].balance;
    s_subscriptions[subscriptionId].balance += uint96(amount);
    s_totalBalance += uint96(amount);
    emit SubscriptionFunded(subscriptionId, oldBalance, oldBalance + amount);
  }

  function getCurrentsubscriptionId() external view returns (uint64) {
    return s_currentsubscriptionId;
  }

  /**
   * @notice Get details about a subscription.
   * @param subscriptionId - ID of the subscription
   * @return balance - LINK balance of the subscription in juels.
   * @return owner - owner of the subscription.
   * @return consumers - list of consumer address which are able to use this subscription.
   */
  function getSubscription(
    uint64 subscriptionId
  ) external view returns (uint96 balance, address owner, address[] memory consumers) {
    if (s_subscriptionConfigs[subscriptionId].owner == address(0)) {
      revert InvalidSubscription();
    }
    return (
      s_subscriptions[subscriptionId].balance,
      s_subscriptionConfigs[subscriptionId].owner,
      s_subscriptionConfigs[subscriptionId].consumers
    );
  }

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
  function createSubscription() external nonReentrant whenNotPaused onlyAuthorizedUsers returns (uint64) {
    s_currentsubscriptionId++;
    uint64 currentsubscriptionId = s_currentsubscriptionId;
    address[] memory consumers = new address[](0);
    s_subscriptions[currentsubscriptionId] = Subscription({balance: 0, blockedBalance: 0});
    s_subscriptionConfigs[currentsubscriptionId] = SubscriptionConfig({
      owner: msg.sender,
      requestedOwner: address(0),
      consumers: consumers
    });

    emit SubscriptionCreated(currentsubscriptionId, msg.sender);
    return currentsubscriptionId;
  }

  /**
   * @notice Gets subscription owner.
   * @param subscriptionId - ID of the subscription
   * @return owner - owner of the subscription.
   */
  function getSubscriptionOwner(uint64 subscriptionId) external view override returns (address owner) {
    if (s_subscriptionConfigs[subscriptionId].owner == address(0)) {
      revert InvalidSubscription();
    }
    return s_subscriptionConfigs[subscriptionId].owner;
  }

  /**
   * @notice Request subscription owner transfer.
   * @param subscriptionId - ID of the subscription
   * @param newOwner - proposed new owner of the subscription
   */
  function requestSubscriptionOwnerTransfer(
    uint64 subscriptionId,
    address newOwner
  ) external onlySubOwner(subscriptionId) nonReentrant whenNotPaused {
    // Proposing to address(0) would never be claimable so don't need to check.
    if (s_subscriptionConfigs[subscriptionId].requestedOwner != newOwner) {
      s_subscriptionConfigs[subscriptionId].requestedOwner = newOwner;
      emit SubscriptionOwnerTransferRequested(subscriptionId, msg.sender, newOwner);
    }
  }

  /**
   * @notice Request subscription owner transfer.
   * @param subscriptionId - ID of the subscription
   * @dev will revert if original owner of subscriptionId has
   * not requested that msg.sender become the new owner.
   */
  function acceptSubscriptionOwnerTransfer(
    uint64 subscriptionId
  ) external nonReentrant whenNotPaused onlyAuthorizedUsers {
    if (s_subscriptionConfigs[subscriptionId].owner == address(0)) {
      revert InvalidSubscription();
    }
    if (s_subscriptionConfigs[subscriptionId].requestedOwner != msg.sender) {
      revert MustBeRequestedOwner(s_subscriptionConfigs[subscriptionId].requestedOwner);
    }
    address oldOwner = s_subscriptionConfigs[subscriptionId].owner;
    s_subscriptionConfigs[subscriptionId].owner = msg.sender;
    s_subscriptionConfigs[subscriptionId].requestedOwner = address(0);
    emit SubscriptionOwnerTransferred(subscriptionId, oldOwner, msg.sender);
  }

  /**
   * @notice Remove a consumer from a Chainlink Functions subscription.
   * @param subscriptionId - ID of the subscription
   * @param consumer - Consumer to remove from the subscription
   */
  function removeConsumer(
    uint64 subscriptionId,
    address consumer
  ) external onlySubOwner(subscriptionId) nonReentrant whenNotPaused {
    if (s_consumers[consumer][subscriptionId] == 0) {
      revert InvalidConsumer(subscriptionId, consumer);
    }
    // Note bounded by MAX_CONSUMERS
    address[] memory consumers = s_subscriptionConfigs[subscriptionId].consumers;
    uint256 lastConsumerIndex = consumers.length - 1;
    for (uint256 i = 0; i < consumers.length; i++) {
      if (consumers[i] == consumer) {
        address last = consumers[lastConsumerIndex];
        // Storage write to preserve last element
        s_subscriptionConfigs[subscriptionId].consumers[i] = last;
        // Storage remove last element
        s_subscriptionConfigs[subscriptionId].consumers.pop();
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
  function addConsumer(
    uint64 subscriptionId,
    address consumer
  ) external onlySubOwner(subscriptionId) nonReentrant whenNotPaused {
    // Already maxed, cannot add any more consumers.
    if (s_subscriptionConfigs[subscriptionId].consumers.length == MAX_CONSUMERS) {
      revert TooManyConsumers();
    }
    if (s_consumers[consumer][subscriptionId] != 0) {
      // Idempotence - do nothing if already added.
      // Ensures uniqueness in s_subscriptions[subscriptionId].consumers.
      return;
    }
    // Initialize the nonce to 1, indicating the consumer is allocated.
    s_consumers[consumer][subscriptionId] = 1;
    s_subscriptionConfigs[subscriptionId].consumers.push(consumer);

    emit SubscriptionConsumerAdded(subscriptionId, consumer);
  }

  /**
   * @notice Cancel a subscription
   * @param subscriptionId - ID of the subscription
   * @param to - Where to send the remaining LINK to
   */
  function cancelSubscription(
    uint64 subscriptionId,
    address to
  ) external onlySubOwner(subscriptionId) nonReentrant whenNotPaused {
    if (pendingRequestExists(subscriptionId)) {
      revert PendingRequestExists();
    }
    cancelSubscriptionHelper(subscriptionId, to);
  }

  function cancelSubscriptionHelper(uint64 subscriptionId, address to) private nonReentrant {
    SubscriptionConfig memory subConfig = s_subscriptionConfigs[subscriptionId];
    uint96 balance = s_subscriptions[subscriptionId].balance;
    // Note bounded by MAX_CONSUMERS;
    // If no consumers, does nothing.
    for (uint256 i = 0; i < subConfig.consumers.length; i++) {
      delete s_consumers[subConfig.consumers[i]][subscriptionId];
    }
    delete s_subscriptionConfigs[subscriptionId];
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
    address[] memory consumers = s_subscriptionConfigs[subscriptionId].consumers;
    address[] memory authorizedSendersList = getAuthorizedSenders();
    for (uint256 i = 0; i < consumers.length; i++) {
      for (uint256 j = 0; j < authorizedSendersList.length; j++) {
        bytes32 requestId = computeRequestId(
          authorizedSendersList[j],
          consumers[i],
          subscriptionId,
          s_consumers[consumers[i]][subscriptionId]
        );
        if (s_requestCommitments[requestId].don != address(0)) {
          return true;
        }
      }
    }
    return false;
  }

  /**
   * @notice Time out all expired requests: unlocks funds and removes the ability for the request to be fulfilled
   * @param requestIdsToTimeout - A list of request IDs to time out
   */

  function timeoutRequests(bytes32[] calldata requestIdsToTimeout) external whenNotPaused {
    for (uint256 i = 0; i < requestIdsToTimeout.length; i++) {
      bytes32 requestId = requestIdsToTimeout[i];
      Commitment memory commitment = s_requestCommitments[requestId];

      // Check that the message sender is the subscription owner
      if (msg.sender != s_subscriptionConfigs[commitment.subscriptionId].owner) {
        revert MustBeSubOwner(s_subscriptionConfigs[commitment.subscriptionId].owner);
      }

      if (commitment.timestamp + s_config.requestTimeoutSeconds > block.timestamp) {
        // Decrement blocked balance
        s_subscriptions[commitment.subscriptionId].blockedBalance -= commitment.estimatedCost;
        // Delete commitment
        delete s_requestCommitments[requestId];
        emit RequestTimedOut(requestId);
      }
    }
  }

  /**
   * @dev The allow list is kept on the Oracle contract. This modifier checks if a user is authorized from there.
   */
  modifier onlyAuthorizedUsers() {
    if (ORACLE_WITH_ALLOWLIST.authorizedReceiverActive() && !ORACLE_WITH_ALLOWLIST.isAuthorizedSender(msg.sender)) {
      revert UnauthorizedSender();
    }
    _;
  }

  modifier onlySubOwner(uint64 subscriptionId) {
    address owner = s_subscriptionConfigs[subscriptionId].owner;
    if (owner == address(0)) {
      revert InvalidSubscription();
    }
    if (msg.sender != owner) {
      revert MustBeSubOwner(owner);
    }
    _;
  }

  modifier nonReentrant() {
    if (s_config.reentrancyLock) {
      revert Reentrant();
    }
    _;
  }

  function _canSetAuthorizedSenders() internal view override onlyOwner returns (bool) {
    return true;
  }

  /**
   * @dev This empty reserved space is put in place to allow future versions to add new
   * variables without shifting down storage in the inheritance chain.
   * See https://docs.openzeppelin.com/contracts/4.x/upgradeable#storage_gaps
   */
  uint256[49] private __gap;
}
