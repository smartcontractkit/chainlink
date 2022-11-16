// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "../../interfaces/LinkTokenInterface.sol";
import "../../interfaces/AggregatorV3Interface.sol";
import "../interfaces/OCR2DRRegistryInterface.sol";
import "../interfaces/OCR2DRBillableInterface.sol";
import "../interfaces/OCR2DRClientInterface.sol";
import "../../interfaces/TypeAndVersionInterface.sol";
import "../../interfaces/ERC677ReceiverInterface.sol";
import "../../ConfirmedOwner.sol";

contract OCR2DRRegistry is ConfirmedOwner, TypeAndVersionInterface, OCR2DRRegistryInterface, ERC677ReceiverInterface {
  LinkTokenInterface public immutable LINK;
  AggregatorV3Interface public immutable LINK_ETH_FEED;

  // We need to maintain a list of consuming addresses.
  // This bound ensures we are able to loop over them as needed.
  // Should a user require more consumers, they can use multiple subscriptions.
  uint16 public constant MAX_CONSUMERS = 100;
  // 5k is plenty for an EXTCODESIZE call (2600) + warm CALL (100)
  // and some arithmetic operations.
  uint256 private constant GAS_FOR_CALL_EXACT_CHECK = 5_000;
  // Maximum number of oracles DON can support
  // Needs to match OCR2Abstract.sol
  uint256 internal constant MAX_NUM_ORACLES = 31;

  error TooManyConsumers();
  error InsufficientBalance();
  error InvalidConsumer(uint64 subscriptionId, address consumer);
  error InvalidSubscription();
  error OnlyCallableFromLink();
  error InvalidCalldata();
  error MustBeSubOwner(address owner);
  error MustBeAllowedDon();
  error PendingRequestExists();
  error MustBeRequestedOwner(address proposedOwner);
  error BalanceInvariantViolated(uint256 internalBalance, uint256 externalBalance); // Should never happen
  event FundsRecovered(address to, uint256 amount);

  struct Subscription {
    // There are only 1e9*1e18 = 1e27 juels in existence, so the balance can fit in uint96 (2^96 ~ 7e28)
    uint96 balance; // Common link balance used for all consumer requests.
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
  mapping(address => mapping(uint64 => uint64)) /* consumer */ /* subscriptionId */ /* nonce */
    private s_consumers;
  mapping(uint64 => SubscriptionConfig) /* subscriptionId */ /* subscriptionConfig */
    private s_subscriptionConfigs;
  mapping(uint64 => Subscription) /* subscriptionId */ /* subscription */
    private s_subscriptions;
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

  error InvalidRequestConfirmations(uint32 have, uint32 min, uint32 max);
  error GasLimitTooBig(uint32 have, uint32 want);
  error NumWordsTooBig(uint32 have, uint32 want);
  error DonAlreadyRegistered(address don);
  error NoSuchDon(address don);
  error InvalidLinkWeiPrice(int256 linkWei);
  error InsufficientGasForConsumer(uint256 have, uint256 want);
  error NoCorrespondingRequest();
  error IncorrectCommitment();
  error BlockhashNotInStore(uint256 blockNum);
  error PaymentTooLarge();
  error Reentrant();

  mapping(address => bool) /* DON */ /* is DON allowed */
    private s_allowedDons;
  address[] private s_dons;
  mapping(address => uint96) /* oracle */ /* LINK balance */
    private s_withdrawableTokens;
  struct Commitment {
    OCR2DRRegistryInterface.RequestBilling billing;
    address don;
    uint96 donFee;
    uint96 registryFee;
  }
  mapping(bytes32 => Commitment) /* requestID */ /* Commitment */
    private s_requestCommitments;
  event DonRegistered(address indexed don);
  event DonDeregistered(address indexed don);
  event BillingStart(
    address indexed don,
    bytes32 requestId,
    uint64 indexed subscriptionId,
    uint32 callbackGasLimit,
    address indexed client
  );
  event BillingEnd(uint64 subscriptionId, bytes32 indexed requestId, uint96 payment, bool success);

  struct Config {
    uint32 maxGasLimit;
    // Reentrancy protection.
    bool reentrancyLock;
    // stalenessSeconds is how long before we consider the feed price to be stale
    // and fallback to fallbackWeiPerUnitLink.
    uint32 stalenessSeconds;
    // Gas to cover oracle payment after we calculate the payment.
    // We make it configurable in case those operations are repriced.
    uint32 gasAfterPaymentCalculation;
    // Represents the average gas execution cost. Used in estimating cost beforehand.
    uint32 gasOverhead;
  }
  int256 private s_fallbackWeiPerUnitLink;
  Config private s_config;
  event ConfigSet(
    uint32 maxGasLimit,
    uint32 stalenessSeconds,
    uint32 gasAfterPaymentCalculation,
    int256 fallbackWeiPerUnitLink,
    uint32 gasOverhead
  );

  constructor(address link, address linkEthFeed) ConfirmedOwner(msg.sender) {
    LINK = LinkTokenInterface(link);
    LINK_ETH_FEED = AggregatorV3Interface(linkEthFeed);
  }

  /**
   * @notice Registers a new Decentralized Oracle Network (DON).
   * @param don address of the DON
   */
  function registerDon(address don) external onlyOwner {
    // NOTE: could validate vesion of OCR2DROracle contract here
    if (s_allowedDons[don] == true) {
      revert DonAlreadyRegistered(don);
    }
    s_allowedDons[don] = true;
    s_dons.push(don);
    emit DonRegistered(don);
  }

  /**
   * @notice Deregisters a Decentralized Oracle Network (DON).
   * @param don address of the DON
   */
  function deregisterDon(address don) external onlyOwner {
    if (s_allowedDons[don] == false) {
      revert NoSuchDon(don);
    }
    delete s_allowedDons[don];
    for (uint256 i = 0; i < s_dons.length; i++) {
      if (s_dons[i] == don) {
        address last = s_dons[s_dons.length - 1];
        // Copy last element and overwrite don to be deleted with it
        s_dons[i] = last;
        s_dons.pop();
      }
    }
    emit DonDeregistered(don);
  }

  /**
   * @notice Sets the configuration of the OCR2DR registry
   * @param maxGasLimit global max for request gas limit
   * @param stalenessSeconds if the eth/link feed is more stale then this, use the fallback price
   * @param gasAfterPaymentCalculation gas used in doing accounting after completing the gas measurement
   * @param fallbackWeiPerUnitLink fallback eth/link price in the case of a stale feed
   * @param gasOverhead fallback eth/link price in the case of a stale feed
   */
  function setConfig(
    uint32 maxGasLimit,
    uint32 stalenessSeconds,
    uint32 gasAfterPaymentCalculation,
    int256 fallbackWeiPerUnitLink,
    uint32 gasOverhead
  ) external onlyOwner {
    if (fallbackWeiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(fallbackWeiPerUnitLink);
    }
    s_config = Config({
      maxGasLimit: maxGasLimit,
      stalenessSeconds: stalenessSeconds,
      gasAfterPaymentCalculation: gasAfterPaymentCalculation,
      reentrancyLock: false,
      gasOverhead: gasOverhead
    });
    s_fallbackWeiPerUnitLink = fallbackWeiPerUnitLink;
    emit ConfigSet(maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, gasOverhead);
  }

  /**
   * @notice Gets the configuration of the OCR2DR registry
   * @return maxGasLimit global max for request gas limit
   * @return stalenessSeconds if the eth/link feed is more stale then this, use the fallback price
   * @return gasAfterPaymentCalculation gas used in doing accounting after completing the gas measurement
   * @return fallbackWeiPerUnitLink fallback eth/link price in the case of a stale feed
   * @return gasOverhead fallback eth/link price in the case of a stale feed
   */
  function getConfig()
    external
    view
    returns (
      uint32 maxGasLimit,
      uint32 stalenessSeconds,
      uint32 gasAfterPaymentCalculation,
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

  function getTotalBalance() external view returns (uint256) {
    return s_totalBalance;
  }

  function getFallbackWeiPerUnitLink() external view returns (int256) {
    return s_fallbackWeiPerUnitLink;
  }

  /**
   * @notice Owner cancel subscription, sends remaining link directly to the subscription owner.
   * @param subscriptionId subscription id
   * @dev notably can be called even if there are pending requests, outstanding ones may fail onchain
   */
  function ownerCancelSubscription(uint64 subscriptionId) external onlyOwner {
    if (s_subscriptionConfigs[subscriptionId].owner == address(0)) {
      revert InvalidSubscription();
    }
    cancelSubscriptionHelper(subscriptionId, s_subscriptionConfigs[subscriptionId].owner);
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
   * @inheritdoc OCR2DRRegistryInterface
   */
  function getRequestConfig() external view override returns (uint32, address[] memory) {
    return (s_config.maxGasLimit, s_dons);
  }

  /**
   * @inheritdoc OCR2DRRegistryInterface
   */
  function getRequiredFee(
    bytes calldata, /* data */
    OCR2DRRegistryInterface.RequestBilling calldata /* billing */
  ) public pure override returns (uint96) {
    // NOTE: Optionally, compute additional fee here
    return 0;
  }

  /**
   * @inheritdoc OCR2DRRegistryInterface
   */
  function estimateExecutionGas(OCR2DRRegistryInterface.RequestBilling calldata billing)
    public
    view
    override
    returns (uint256)
  {
    return s_config.gasOverhead + s_config.gasAfterPaymentCalculation + billing.gasLimit;
  }

  /**
   * @inheritdoc OCR2DRRegistryInterface
   */
  function estimateCost(
    bytes calldata data,
    OCR2DRRegistryInterface.RequestBilling calldata billing,
    uint96 donRequiredFee
  ) public view override returns (uint96) {
    int256 weiPerUnitLink;
    weiPerUnitLink = getFeedData();
    if (weiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(weiPerUnitLink);
    }
    uint256 executionGas = estimateExecutionGas(billing);
    // (1e18 juels/link) (wei/gas * gas) / (wei/link) = juels
    uint256 paymentNoFee = (1e18 * tx.gasprice * executionGas) / uint256(weiPerUnitLink);
    uint96 registryFee = getRequiredFee(data, billing);
    uint256 fee = uint256(donRequiredFee) + uint256(registryFee);
    if (paymentNoFee > (1e27 - fee)) {
      revert PaymentTooLarge(); // Payment + fee cannot be more than all of the link in existence.
    }
    return uint96(paymentNoFee + fee);
  }

  /**
   * @inheritdoc OCR2DRRegistryInterface
   */
  function beginBilling(bytes calldata data, RequestBilling calldata billing)
    external
    override
    onlyAllowedDons
    nonReentrant
    returns (bytes32)
  {
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
    uint256 estimatedCost = estimateCost(
      data,
      billing,
      OCR2DRBillableInterface(msg.sender).getRequiredFee(data, billing)
    );
    if (s_subscriptions[billing.subscriptionId].balance < estimatedCost) {
      revert InsufficientBalance();
    }

    uint64 nonce = currentNonce + 1;
    (bytes32 requestId, ) = computeRequestId(msg.sender, billing.client, billing.subscriptionId, nonce);

    s_requestCommitments[requestId] = Commitment(
      billing,
      msg.sender,
      OCR2DRBillableInterface(msg.sender).getRequiredFee(data, billing),
      getRequiredFee(data, billing)
    );

    emit BillingStart(msg.sender, requestId, billing.subscriptionId, billing.gasLimit, billing.client);
    s_consumers[billing.client][billing.subscriptionId] = nonce;
    return requestId;
  }

  /**
   * @inheritdoc OCR2DRRegistryInterface
   */
  function getCommitment(bytes32 requestId)
    external
    view
    override
    returns (
      address,
      uint64,
      uint32
    )
  {
    Commitment memory commitment = s_requestCommitments[requestId];
    return (commitment.billing.client, commitment.billing.subscriptionId, commitment.billing.gasLimit);
  }

  function computeRequestId(
    address don,
    address client,
    uint64 subscriptionId,
    uint64 nonce
  ) private pure returns (bytes32, uint256) {
    uint256 preSeed = uint256(keccak256(abi.encode(don, client, subscriptionId, nonce)));
    return (keccak256(abi.encode(don, preSeed)), preSeed);
  }

  /**
   * @dev calls target address with exactly gasAmount gas and data as calldata
   * or reverts if at least gasAmount gas is not available.
   */
  function callWithExactGas(
    uint256 gasAmount,
    address target,
    bytes memory data
  ) private returns (bool success) {
    // solhint-disable-next-line no-inline-assembly
    assembly {
      let g := gas()
      // Compute g -= GAS_FOR_CALL_EXACT_CHECK and check for underflow
      // The gas actually passed to the callee is min(gasAmount, 63//64*gas available).
      // We want to ensure that we revert if gasAmount >  63//64*gas available
      // as we do not want to provide them with less, however that check itself costs
      // gas.  GAS_FOR_CALL_EXACT_CHECK ensures we have at least enough gas to be able
      // to revert if gasAmount >  63//64*gas available.
      if lt(g, GAS_FOR_CALL_EXACT_CHECK) {
        revert(0, 0)
      }
      g := sub(g, GAS_FOR_CALL_EXACT_CHECK)
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
   * @inheritdoc OCR2DRRegistryInterface
   */
  function concludeBilling(
    bytes32 requestId,
    bytes calldata response,
    bytes calldata err,
    address transmitter,
    address[MAX_NUM_ORACLES] memory, /* signers */
    uint32 initialGas
  ) external onlyAllowedDons nonReentrant returns (uint96) {
    Commitment memory commitment = s_requestCommitments[requestId];
    if (commitment.billing.client == address(0)) {
      revert IncorrectCommitment();
    }
    delete s_requestCommitments[requestId];

    bytes memory callback = abi.encodeWithSelector(
      OCR2DRClientInterface.handleOracleFulfillment.selector,
      requestId,
      response,
      err
    );
    // Call with explicitly the amount of callback gas requested
    // Important to not let them exhaust the gas budget and avoid oracle payment.
    // Do not allow any non-view/non-pure coordinator functions to be called
    // during the consumers callback code via reentrancyLock.
    // NOTE: that callWithExactGas will revert if we do not have sufficient gas
    // to give the callee their requested amount.
    s_config.reentrancyLock = true;
    bool success = callWithExactGas(commitment.billing.gasLimit, commitment.billing.client, callback);
    s_config.reentrancyLock = false;

    // We want to charge users exactly for how much gas they use in their callback.
    // The gasAfterPaymentCalculation is meant to cover these additional operations where we
    // decrement the subscription balance and increment the oracles withdrawable balance.
    uint96 payment = calculatePaymentAmount(
      initialGas,
      s_config.gasAfterPaymentCalculation,
      commitment.donFee,
      commitment.registryFee,
      tx.gasprice
    );
    if (s_subscriptions[commitment.billing.subscriptionId].balance < payment) {
      revert InsufficientBalance();
    }
    /**
     * Oracle Payment *
     * Two options here:
     *   1. Reimburse the transmitter for execution cost, then split the requiredFee across all participants.
     *   2. Pay transmitter the full amount. Since the transmitter is chosen OCR, we trust the fairness of their selection algorithm.
     * Using Option 1 here.
     **/
    s_subscriptions[commitment.billing.subscriptionId].balance -= payment;
    s_withdrawableTokens[transmitter] += payment;
    // Include payment in the event for tracking costs.
    emit BillingEnd(commitment.billing.subscriptionId, requestId, payment, success);
    return payment;
  }

  // Get the amount of gas used for fulfillment
  function calculatePaymentAmount(
    uint256 startGas,
    uint32 gasAfterPaymentCalculation,
    uint96 donFee,
    uint96 registryFee,
    uint256 weiPerUnitGas
  ) internal view returns (uint96) {
    int256 weiPerUnitLink;
    weiPerUnitLink = getFeedData();
    if (weiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(weiPerUnitLink);
    }
    // (1e18 juels/link) (wei/gas * gas) / (wei/link) = juels
    uint256 paymentNoFee = (1e18 * weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft())) /
      uint256(weiPerUnitLink);
    uint256 fee = uint256(donFee) + uint256(registryFee);
    if (paymentNoFee > (1e27 - fee)) {
      revert PaymentTooLarge(); // Payment + fee cannot be more than all of the link in existence.
    }
    return uint96(paymentNoFee + fee);
  }

  function getFeedData() private view returns (int256) {
    uint32 stalenessSeconds = s_config.stalenessSeconds;
    bool staleFallback = stalenessSeconds > 0;
    uint256 timestamp;
    int256 weiPerUnitLink;
    (, weiPerUnitLink, , timestamp, ) = LINK_ETH_FEED.latestRoundData();
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
  function oracleWithdraw(address recipient, uint96 amount) external nonReentrant {
    if (amount == 0) amount = s_withdrawableTokens[msg.sender];
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
   * @inheritdoc OCR2DRRegistryInterface
   */
  function getSubscription(uint64 subscriptionId)
    external
    view
    override
    returns (
      uint96 balance,
      address owner,
      address[] memory consumers
    )
  {
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
   * @inheritdoc OCR2DRRegistryInterface
   */
  function createSubscription() external override nonReentrant returns (uint64) {
    s_currentsubscriptionId++;
    uint64 currentsubscriptionId = s_currentsubscriptionId;
    address[] memory consumers = new address[](0);
    s_subscriptions[currentsubscriptionId] = Subscription({balance: 0});
    s_subscriptionConfigs[currentsubscriptionId] = SubscriptionConfig({
      owner: msg.sender,
      requestedOwner: address(0),
      consumers: consumers
    });

    emit SubscriptionCreated(currentsubscriptionId, msg.sender);
    return currentsubscriptionId;
  }

  /**
   * @inheritdoc OCR2DRRegistryInterface
   */
  function requestSubscriptionOwnerTransfer(uint64 subscriptionId, address newOwner)
    external
    override
    onlySubOwner(subscriptionId)
    nonReentrant
  {
    // Proposing to address(0) would never be claimable so don't need to check.
    if (s_subscriptionConfigs[subscriptionId].requestedOwner != newOwner) {
      s_subscriptionConfigs[subscriptionId].requestedOwner = newOwner;
      emit SubscriptionOwnerTransferRequested(subscriptionId, msg.sender, newOwner);
    }
  }

  /**
   * @inheritdoc OCR2DRRegistryInterface
   */
  function acceptSubscriptionOwnerTransfer(uint64 subscriptionId) external override nonReentrant {
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
   * @inheritdoc OCR2DRRegistryInterface
   */
  function removeConsumer(uint64 subscriptionId, address consumer)
    external
    override
    onlySubOwner(subscriptionId)
    nonReentrant
  {
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
   * @inheritdoc OCR2DRRegistryInterface
   */
  function addConsumer(uint64 subscriptionId, address consumer)
    external
    override
    onlySubOwner(subscriptionId)
    nonReentrant
  {
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
   * @inheritdoc OCR2DRRegistryInterface
   */
  function cancelSubscription(uint64 subscriptionId, address to)
    external
    override
    onlySubOwner(subscriptionId)
    nonReentrant
  {
    if (pendingRequestExists(subscriptionId)) {
      revert PendingRequestExists();
    }
    cancelSubscriptionHelper(subscriptionId, to);
  }

  function cancelSubscriptionHelper(uint64 subscriptionId, address to) private nonReentrant {
    SubscriptionConfig memory subConfig = s_subscriptionConfigs[subscriptionId];
    Subscription memory sub = s_subscriptions[subscriptionId];
    uint96 balance = sub.balance;
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
   * @inheritdoc OCR2DRRegistryInterface
   * @dev Looping is bounded to MAX_CONSUMERS*(number of DONs).
   * @dev Used to disable subscription canceling while outstanding request are present.
   */
  function pendingRequestExists(uint64 subscriptionId) public view override returns (bool) {
    SubscriptionConfig memory subConfig = s_subscriptionConfigs[subscriptionId];
    for (uint256 i = 0; i < subConfig.consumers.length; i++) {
      for (uint256 j = 0; j < s_dons.length; j++) {
        (bytes32 reqId, ) = computeRequestId(
          s_dons[j],
          subConfig.consumers[i],
          subscriptionId,
          s_consumers[subConfig.consumers[i]][subscriptionId]
        );
        if (s_requestCommitments[reqId].don != address(0)) {
          return true;
        }
      }
    }
    return false;
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

  modifier onlyAllowedDons() {
    if (!s_allowedDons[msg.sender]) {
      revert MustBeAllowedDon();
    }
    _;
  }

  modifier nonReentrant() {
    if (s_config.reentrancyLock) {
      revert Reentrant();
    }
    _;
  }

  /**
   * @notice The type and version of this contract
   * @return Type and version string
   */
  function typeAndVersion() external pure virtual override returns (string memory) {
    return "OCR2DRRegistry 0.0.0";
  }
}
