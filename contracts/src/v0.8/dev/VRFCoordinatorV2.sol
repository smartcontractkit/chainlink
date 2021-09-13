// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/BlockhashStoreInterface.sol";
import "../interfaces/AggregatorV3Interface.sol";
import "../interfaces/TypeAndVersionInterface.sol";

import "./VRF.sol";
import "../ConfirmedOwner.sol";
import "./VRFConsumerBaseV2.sol";

contract VRFCoordinatorV2 is VRF, ConfirmedOwner, TypeAndVersionInterface {

  LinkTokenInterface public immutable LINK;
  AggregatorV3Interface public immutable LINK_ETH_FEED;
  BlockhashStoreInterface public immutable BLOCKHASH_STORE;

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
  error MustBeRequestedOwner(address proposedOwner);
  error BalanceInvariantViolated(uint256 internalBalance, uint256 externalBalance); // Should never happen
  event FundsRecovered(address to, uint256 amount);
  struct Subscription {
    // There are only 1e9*1e18 = 1e27 juels in existence, so the balance can fit in uint96 (2^96 ~ 7e28)
    uint96 balance; // Common link balance used for all consumer requests.
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
  struct Consumer {
    uint64  subId;
    uint64  nonce;
  }
  mapping(address /* consumer */ => mapping(uint64 /* subId */ => Consumer )) private s_consumers;
  mapping(uint64 /* subId */ => Subscription /* subscription */) private s_subscriptions;
  uint64 private s_currentSubId;
  // s_totalBalance tracks the total link sent to/from
  // this contract through onTokenTransfer, defundSubscription, cancelSubscription and oracleWithdraw.
  // A discrepancy with this contract's link balance indicates someone
  // sent tokens using transfer and so we may need to use recoverFunds.
  uint96 public s_totalBalance;
  event SubscriptionCreated(uint64 indexed subId, address owner);
  event SubscriptionFunded(uint64 indexed subId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionConsumerAdded(uint64 indexed subId, address consumer);
  event SubscriptionConsumerRemoved(uint64 indexed subId, address consumer);
  event SubscriptionDefunded(uint64 indexed subId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionCanceled(uint64 indexed subId, address to, uint256 amount);
  event SubscriptionOwnerTransferRequested(uint64 indexed subId, address from, address to);
  event SubscriptionOwnerTransferred(uint64 indexed subId, address from, address to);

  // Set this maximum to 200 to give us a 56 block window to fulfill
  // the request before requiring the block hash feeder.
  uint16 public constant MAX_REQUEST_CONFIRMATIONS = 200;
  uint32 public constant MAX_NUM_WORDS = 500;
  // The minimum gas limit that could be requested for a callback.
  // Set to 5k to ensure plenty of room to make the call itself.
  uint256 public constant MIN_GAS_LIMIT = 5_000;
  error InvalidRequestConfirmations(uint16 have, uint16 min, uint16 max);
  error GasLimitTooBig(uint32 have, uint32 want);
  error NumWordsTooBig(uint32 have, uint32 want);
  error ProvingKeyAlreadyRegistered(bytes32 keyHash);
  error NoSuchProvingKey(bytes32 keyHash);
  error InvalidLinkWeiPrice(int256 linkWei);
  error InsufficientGasForConsumer(uint256 have, uint256 want);
  error NoCorrespondingRequest();
  error IncorrectCommitment();
  error BlockhashNotInStore(uint256 blockNum);
  error PaymentTooLarge();
  error Reentrant();
  struct RequestCommitment {
    uint64 blockNum;
    uint64 subId;
    uint32 callbackGasLimit;
    uint32 numWords;
    address sender;
  }
  mapping(bytes32 /* keyHash */ => address /* oracle */) private s_provingKeys;
  mapping(address /* oracle */ => uint96 /* LINK balance */) private s_withdrawableTokens;
  mapping(uint256 /* requestID */ => bytes32 /* commitment */) private s_requestCommitments;
  event ProvingKeyRegistered(bytes32 keyHash, address indexed oracle);
  event ProvingKeyDeregistered(bytes32 keyHash, address indexed oracle);
  event RandomWordsRequested(
    bytes32 indexed keyHash,
    uint256 requestId,
    uint256 preSeed,
    uint64 subId,
    uint16 minimumRequestConfirmations,
    uint32 callbackGasLimit,
    uint32 numWords,
    address indexed sender
  );
  event RandomWordsFulfilled(
    uint256 indexed requestId,
    uint256[] output,
    bool success
  );

  struct Config {
    uint16 minimumRequestConfirmations;
    // Flat fee charged per fulfillment in millionths of link
    // So fee range is [0, 2^32/10^6].
    uint32 fulfillmentFlatFeeLinkPPM;
    uint32 maxGasLimit;
    // stalenessSeconds is how long before we consider the feed price to be stale
    // and fallback to fallbackWeiPerUnitLink.
    uint32 stalenessSeconds;
    // Gas to cover oracle payment after we calculate the payment.
    // We make it configurable in case those operations are repriced.
    uint32 gasAfterPaymentCalculation;
    uint96 minimumSubscriptionBalance;
    // Re-entrancy protection.
    bool reentrancyLock;
  }
  int256 internal s_fallbackWeiPerUnitLink;
  Config private s_config;
  event ConfigSet(
    uint16 minimumRequestConfirmations,
    uint32 fulfillmentFlatFeeLinkPPM,
    uint32 maxGasLimit,
    uint32 stalenessSeconds,
    uint32 gasAfterPaymentCalculation,
    uint96 minimumSubscriptionBalance,
    int256 fallbackWeiPerUnitLink
  );

  constructor(
    address link,
    address blockhashStore,
    address linkEthFeed
  )
    ConfirmedOwner(msg.sender)
  {
    LINK = LinkTokenInterface(link);
    LINK_ETH_FEED = AggregatorV3Interface(linkEthFeed);
    BLOCKHASH_STORE = BlockhashStoreInterface(blockhashStore);
  }

  /**
   * @notice Registers a proving key to an oracle.
   * @param oracle address of the oracle
   * @param publicProvingKey key that oracle can use to submit vrf fulfillments
   */
  function registerProvingKey(
    address oracle,
    uint256[2] calldata publicProvingKey
  )
    external
    onlyOwner()
  {
    bytes32 kh = hashOfKey(publicProvingKey);
    if (s_provingKeys[kh] != address(0)) {
      revert ProvingKeyAlreadyRegistered(kh);
    }
    s_provingKeys[kh] = oracle;
    emit ProvingKeyRegistered(kh, oracle);
  }

  /**
   * @notice Deregisters a proving key to an oracle.
   * @param publicProvingKey key that oracle can use to submit vrf fulfillments
   */
  function deregisterProvingKey(
    uint256[2] calldata publicProvingKey
  )
    external
    onlyOwner()
  {
    bytes32 kh = hashOfKey(publicProvingKey);
    address oracle = s_provingKeys[kh];
    if (oracle == address(0)) {
      revert NoSuchProvingKey(kh);
    }
    delete s_provingKeys[kh];
    emit ProvingKeyDeregistered(kh, oracle);
  }

  /**
   * @notice Returns the serviceAgreements key associated with this public key
   * @param publicKey the key to return the address for
   */
  function hashOfKey(
    uint256[2] memory publicKey
  )
    public
    pure
    returns (
      bytes32
    )
  {
    return keccak256(abi.encode(publicKey));
  }

  function setConfig(
    uint16 minimumRequestConfirmations,
    uint32 fulfillmentFlatFeeLinkPPM,
    uint32 maxGasLimit,
    uint32 stalenessSeconds,
    uint32 gasAfterPaymentCalculation,
    uint96 minimumSubscriptionBalance,
    int256 fallbackWeiPerUnitLink
  )
    external
    onlyOwner()
  {
    if (minimumRequestConfirmations > MAX_REQUEST_CONFIRMATIONS) {
      revert InvalidRequestConfirmations(minimumRequestConfirmations, minimumRequestConfirmations, MAX_REQUEST_CONFIRMATIONS);
    }
    if (fallbackWeiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(fallbackWeiPerUnitLink);
    }
    s_config = Config({
      minimumRequestConfirmations: minimumRequestConfirmations,
      fulfillmentFlatFeeLinkPPM: fulfillmentFlatFeeLinkPPM,
      maxGasLimit: maxGasLimit,
      stalenessSeconds: stalenessSeconds,
      gasAfterPaymentCalculation: gasAfterPaymentCalculation,
      minimumSubscriptionBalance: minimumSubscriptionBalance,
      reentrancyLock: false
    });
    s_fallbackWeiPerUnitLink = fallbackWeiPerUnitLink;
    emit ConfigSet(
      minimumRequestConfirmations,
      fulfillmentFlatFeeLinkPPM,
      maxGasLimit,
      stalenessSeconds,
      gasAfterPaymentCalculation,
      minimumSubscriptionBalance,
      fallbackWeiPerUnitLink
    );
  }

  /**
   * @notice read the current configuration of the coordinator.
   */
  function getConfig()
    external
    view
    returns (
      uint16 minimumRequestConfirmations,
      uint32 fulfillmentFlatFeeLinkPPM,
      uint32 maxGasLimit,
      uint32 stalenessSeconds,
      uint32 gasAfterPaymentCalculation,
      uint96 minimumSubscriptionBalance,
      int256 fallbackWeiPerUnitLink
    )
  {
    Config memory config = s_config;
    return (
      config.minimumRequestConfirmations,
      config.fulfillmentFlatFeeLinkPPM,
      config.maxGasLimit,
      config.stalenessSeconds,
      config.gasAfterPaymentCalculation,
      config.minimumSubscriptionBalance,
      s_fallbackWeiPerUnitLink
    );
  }

  function recoverFunds(
    address to
  )
    external
    onlyOwner()
  {
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

  // Want to ensure these arguments can fit inside of 2 words
  // so in the worst case where the consuming contract has to read all of them
  // from storage, it only has to read 2 words.
  function requestRandomWords(
    bytes32 keyHash,  // Corresponds to a particular offchain job which uses that key for the proofs
    uint64  subId,
    uint16  requestConfirmations,
    uint32  callbackGasLimit,
    uint32  numWords  // Desired number of random words
  )
    external
    nonReentrant()
    returns (
      uint256
    )
  {
    // Input validation using the subscription storage.
    if (s_subscriptions[subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    // Its important to ensure that the consumer is in fact who they say they
    // are, otherwise they could use someone else's subscription balance.
    Consumer memory consumer = s_consumers[msg.sender][subId];
    if (consumer.subId == 0) {
      revert InvalidConsumer(subId, msg.sender);
    }
    // Input validation using the config storage word.
    if (requestConfirmations < s_config.minimumRequestConfirmations || requestConfirmations > MAX_REQUEST_CONFIRMATIONS) {
      revert InvalidRequestConfirmations(requestConfirmations, s_config.minimumRequestConfirmations, MAX_REQUEST_CONFIRMATIONS);
    }
    if (s_subscriptions[subId].balance < s_config.minimumSubscriptionBalance) {
      revert InsufficientBalance();
    }
    if (callbackGasLimit > s_config.maxGasLimit) {
      revert GasLimitTooBig(callbackGasLimit, s_config.maxGasLimit);
    }
    if (numWords > MAX_NUM_WORDS) {
      revert NumWordsTooBig(numWords, MAX_NUM_WORDS);
    }
    // Note we do not check whether the keyHash is valid to save gas.
    // The consequence for users is that they can send requests
    // for invalid keyHashes which will simply not be fulfilled.
    uint64 nonce = consumer.nonce + 1;
    uint256 preSeed = uint256(keccak256(abi.encode(keyHash, msg.sender, subId, nonce)));
    uint256 requestId = uint256(keccak256(abi.encode(keyHash, preSeed)));

    s_requestCommitments[requestId] = keccak256(abi.encode(
        requestId,
        block.number,
        subId,
        callbackGasLimit,
        numWords,
        msg.sender));
    emit RandomWordsRequested(keyHash, requestId, preSeed, subId, requestConfirmations, callbackGasLimit, numWords, msg.sender);
    s_consumers[msg.sender][subId].nonce = nonce;

    return requestId;
  }

  function getCommitment(
    uint256 requestId
  )
    external
    view
    returns (
      bytes32
    )
  {
    return s_requestCommitments[requestId];
  }

  /**
   * @dev calls target address with exactly gasAmount gas and data as calldata
   * or reverts if at least gasAmount gas is not available.
   * The maximum amount of gasAmount is all gas available but 1/64th.
   * The minimum amount of gasAmount is MIN_GAS_LIMIT.
   */
  function callWithExactGas(
    uint256 gasAmount,
    address target,
    bytes memory data
  )
    private
    returns (
      bool success
    )
  {
    // solhint-disable-next-line no-inline-assembly
    assembly{
      let g := gas()
      // Compute g -= MIN_GAS_LIMIT and check for underflow
      if lt(g, MIN_GAS_LIMIT) { revert(0, 0) }
      g := sub(g, MIN_GAS_LIMIT)
      // if g - g//64 <= gasAmount, revert
      // (we subtract g//64 because of EIP-150)
      if iszero(gt(sub(g, div(g, 64)), gasAmount)) { revert(0, 0) }
      // solidity calls check that a contract actually exists at the destination, so we do the same
      if iszero(extcodesize(target)) { revert(0, 0) }
      // call and return whether we succeeded. ignore return data
      success := call(gasAmount, target, 0, add(data, 0x20), mload(data), 0, 0)
    }
    return success;
  }

  function getRandomnessFromProof(
    Proof memory proof,
    RequestCommitment memory rc
  )
    private
    view
      returns (
      bytes32 keyHash,
      uint256 requestId,
      uint256 randomness
    ) {
      keyHash = hashOfKey(proof.pk);
      // Only registered proving keys are permitted.
      address oracle = s_provingKeys[keyHash];
      if (oracle == address(0)) {
        revert NoSuchProvingKey(keyHash);
      }
      requestId = uint256(keccak256(abi.encode(keyHash, proof.seed)));
      bytes32 commitment = s_requestCommitments[requestId];
      if (commitment == 0) {
        revert NoCorrespondingRequest();
      }
      if (commitment != keccak256(abi.encode(
        requestId,
        rc.blockNum,
        rc.subId,
        rc.callbackGasLimit,
        rc.numWords,
        rc.sender)))
      {
        revert IncorrectCommitment();
      }

      bytes32 blockHash = blockhash(rc.blockNum);
        if (blockHash == bytes32(0)) {
          blockHash = BLOCKHASH_STORE.getBlockhash(rc.blockNum);
          if (blockHash == bytes32(0)) {
            revert BlockhashNotInStore(rc.blockNum);
          }
      }

      // The seed actually used by the VRF machinery, mixing in the blockhash
      uint256 actualSeed = uint256(keccak256(abi.encodePacked(proof.seed, blockHash)));
      randomness = VRF.randomValueFromVRFProof(proof, actualSeed); // Reverts on failure
  }

  function fulfillRandomWords(
    Proof memory proof,
    RequestCommitment memory rc
  )
    external
    nonReentrant()
  {
    uint256 startGas = gasleft();
    (bytes32 keyHash,
     uint256 requestId,
     uint256 randomness) = getRandomnessFromProof(proof, rc);

    uint256[] memory randomWords = new uint256[](rc.numWords);
    for (uint256 i = 0; i < rc.numWords; i++) {
      randomWords[i] = uint256(keccak256(abi.encode(randomness, i)));
    }

    delete s_requestCommitments[requestId];
    VRFConsumerBaseV2 v;
    bytes memory resp = abi.encodeWithSelector(v.rawFulfillRandomWords.selector, proof.seed, randomWords);
    uint256 gasPreCallback = gasleft();
    if (gasPreCallback < rc.callbackGasLimit) {
      revert InsufficientGasForConsumer(gasPreCallback, rc.callbackGasLimit);
    }
    // Call with explicitly the amount of callback gas requested
    // Important to not let them exhaust the gas budget and avoid oracle payment.
    // Do not allow any non-view/non-pure coordinator functions to be called
    // during the consumers callback code via reentrancyLock.
    s_config.reentrancyLock = true;
    bool success = callWithExactGas(rc.callbackGasLimit, rc.sender, resp);
    emit RandomWordsFulfilled(requestId, randomWords, success);
    s_config.reentrancyLock = false;

    // We want to charge users exactly for how much gas they use in their callback.
    // The gasAfterPaymentCalculation is meant to cover these additional operations where we
    // decrement the subscription balance and increment the oracles withdrawable balance.
    // We also add the flat link fee to the payment amount.
    // Its specified in millionths of link, if s_config.fulfillmentFlatFeeLinkPPM = 1
    // 1 link / 1e6 = 1e18 juels / 1e6 = 1e12 juels.
    uint96 payment = calculatePaymentAmount(startGas, s_config.gasAfterPaymentCalculation, s_config.fulfillmentFlatFeeLinkPPM, tx.gasprice);
    if (s_subscriptions[rc.subId].balance < payment) {
      revert InsufficientBalance();
    }
    s_subscriptions[rc.subId].balance -= payment;
    s_withdrawableTokens[s_provingKeys[keyHash]] += payment;
  }

  // Get the amount of gas used for fulfillment
  function calculatePaymentAmount(
      uint256 startGas,
      uint256 gasAfterPaymentCalculation,
      uint32  fulfillmentFlatFeeLinkPPM,
      uint256 weiPerUnitGas
  )
    internal
    view
    returns (
      uint96
    )
  {
    int256 weiPerUnitLink;
    weiPerUnitLink = getFeedData();
    if (weiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(weiPerUnitLink);
    }
    // (1e18 juels/link) (wei/gas * gas) / (wei/link) = juels
    uint256 paymentNoFee = 1e18*weiPerUnitGas*(gasAfterPaymentCalculation + startGas - gasleft()) / uint256(weiPerUnitLink);
    uint256 fee = 1e12*uint256(fulfillmentFlatFeeLinkPPM);
    if (paymentNoFee > (1e27-fee)) {
      revert PaymentTooLarge(); // Payment + fee cannot be more than all of the link in existence.
    }
    return uint96(paymentNoFee+fee);
  }

  function getFeedData()
    private
    view
    returns (
      int256
    )
  {
    uint32 stalenessSeconds = s_config.stalenessSeconds;
    bool staleFallback = stalenessSeconds > 0;
    uint256 timestamp;
    int256 weiPerUnitLink;
    (,weiPerUnitLink,,timestamp,) = LINK_ETH_FEED.latestRoundData();
    // solhint-disable-next-line not-rely-on-time
    if (staleFallback && stalenessSeconds < block.timestamp - timestamp) {
      weiPerUnitLink = s_fallbackWeiPerUnitLink;
    }
    return weiPerUnitLink;
  }

  function oracleWithdraw(
    address recipient,
    uint96 amount
  )
    external
    nonReentrant()
  {
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
    address sender,
    uint256 amount,
    bytes calldata data
  )
    external
    nonReentrant()
  {
    if (msg.sender != address(LINK)) {
      revert OnlyCallableFromLink();
    }
    if (data.length != 32) {
      revert InvalidCalldata();
    }
    uint64 subId = abi.decode(data, (uint64));
    if (s_subscriptions[subId].owner == address(0))  {
      revert InvalidSubscription();
    }
    address owner = s_subscriptions[subId].owner;
    if (owner != sender) {
      revert MustBeSubOwner(owner);
    }
    uint256 oldBalance = s_subscriptions[subId].balance;
    s_subscriptions[subId].balance += uint96(amount);
    s_totalBalance += uint96(amount);
    emit SubscriptionFunded(subId, oldBalance, oldBalance+amount);
  }

  function getSubscription(
    uint64 subId
  )
    external
    view
    returns (
      uint96 balance,
      address owner,
      address[] memory consumers
    )
  {
    if (s_subscriptions[subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    return (
      s_subscriptions[subId].balance,
      s_subscriptions[subId].owner,
      s_subscriptions[subId].consumers
    );
  }

  function createSubscription()
    external
    nonReentrant()
    returns (
      uint64
    )
  {
    s_currentSubId++;
    uint64 currentSubId = s_currentSubId;
    address[] memory consumers = new address[](0);
    s_subscriptions[currentSubId] = Subscription({
      balance: 0,
      owner: msg.sender,
      requestedOwner: address(0),
      consumers: consumers
    });

    emit SubscriptionCreated(currentSubId, msg.sender);
    return currentSubId;
  }

  function requestSubscriptionOwnerTransfer(
    uint64 subId,
    address newOwner
  )
    external
    onlySubOwner(subId)
    nonReentrant()
  {
    // Proposing to address(0) would never be claimable so don't need to check.
    if (s_subscriptions[subId].requestedOwner != newOwner) {
      s_subscriptions[subId].requestedOwner = newOwner;
      emit SubscriptionOwnerTransferRequested(subId, msg.sender, newOwner);
    }
  }

  function acceptSubscriptionOwnerTransfer(
    uint64 subId
  )
    external
    nonReentrant()
  {
    if (s_subscriptions[subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    if (s_subscriptions[subId].requestedOwner != msg.sender) {
      revert MustBeRequestedOwner(s_subscriptions[subId].requestedOwner);
    }
    address oldOwner = s_subscriptions[subId].owner;
    s_subscriptions[subId].owner = msg.sender;
    s_subscriptions[subId].requestedOwner = address(0);
    emit SubscriptionOwnerTransferred(subId, oldOwner, msg.sender);
  }

  function removeConsumer(
    uint64 subId,
    address consumer
  )
    external
    onlySubOwner(subId)
    nonReentrant()
  {
    if (s_consumers[consumer][subId].subId == 0) {
      revert InvalidConsumer(subId, consumer);
    }
    // Note bounded by MAX_CONSUMERS
    address[] memory consumers = s_subscriptions[subId].consumers;
    uint256 lastConsumerIndex = consumers.length-1;
    for (uint256 i = 0; i < consumers.length; i++) {
      if (consumers[i] == consumer) {
        address last = consumers[lastConsumerIndex];
        // Storage write to preserve last element
        s_subscriptions[subId].consumers[i] = last;
        // Storage remove last element
        s_subscriptions[subId].consumers.pop();
        break;
      }
    }
    delete s_consumers[consumer][subId];
    emit SubscriptionConsumerRemoved(subId, consumer);
  }

  function addConsumer(
    uint64 subId,
    address consumer
  )
    external
    onlySubOwner(subId)
    nonReentrant()
  {
    // Already maxed, cannot add any more consumers.
    if (s_subscriptions[subId].consumers.length == MAX_CONSUMERS) {
      revert TooManyConsumers();
    }
    if (s_consumers[consumer][subId].subId != 0) {
      // Idempotence - do nothing if already added.
      // Ensures uniqueness in s_subscriptions[subId].consumers.
      return;
    }
    s_consumers[consumer][subId] = Consumer({
      subId: subId,
      nonce: 0
    });
    s_subscriptions[subId].consumers.push(consumer);

    emit SubscriptionConsumerAdded(subId, consumer);
  }

  function defundSubscription(
    uint64 subId,
    address to,
    uint96 amount
  )
    external
    onlySubOwner(subId)
    nonReentrant()
  {
    if (s_subscriptions[subId].balance < amount) {
      revert InsufficientBalance();
    }
    uint256 oldBalance = s_subscriptions[subId].balance;
    s_subscriptions[subId].balance -= amount;
    s_totalBalance -= amount;
    if (!LINK.transfer(to, amount)) {
      revert InsufficientBalance();
    }
    emit SubscriptionDefunded(subId, oldBalance, s_subscriptions[subId].balance);
  }

  // Keep this separate from zeroing, perhaps there is a use case where consumers
  // want to keep the subId, but withdraw all the link.
  function cancelSubscription(
    uint64 subId,
    address to
  )
    external
    onlySubOwner(subId)
    nonReentrant()
  {
    Subscription memory sub = s_subscriptions[subId];
    uint96 balance = sub.balance;
    // Note bounded by MAX_CONSUMERS;
    // If no consumers, does nothing.
    for (uint256 i = 0; i < sub.consumers.length; i++) {
      delete s_consumers[sub.consumers[i]][subId];
    }
    delete s_subscriptions[subId];
    s_totalBalance -= balance;
    if (!LINK.transfer(to, uint256(balance))) {
      revert InsufficientBalance();
    }
    emit SubscriptionCanceled(subId, to, balance);
  }

  modifier onlySubOwner(uint64 subId) {
    address owner = s_subscriptions[subId].owner;
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

  /**
   * @notice The type and version of this contract
   * @return Type and version string
   */
  function typeAndVersion()
    external
    pure
    virtual
    override
    returns (
      string memory
    )
  {
    return "VRFCoordinatorV2 1.0.0";
  }
}
