// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/BlockhashStoreInterface.sol";
import "../interfaces/AggregatorV3Interface.sol";
import "../interfaces/TypeAndVersionInterface.sol";

import "./VRF.sol";
import "../ConfirmedOwner.sol";
import "../interfaces/VRFConsumerV2Interface.sol";

contract VRFCoordinatorV2 is VRF, ConfirmedOwner, TypeAndVersionInterface {

  LinkTokenInterface public immutable LINK;
  AggregatorV3Interface public immutable LINK_ETH_FEED;
  BlockhashStoreInterface public immutable BLOCKHASH_STORE;

  error InsufficientBalance();
  error InvalidConsumer(uint64 subId, address consumer);
  error InvalidSubscription();
  error AlreadySubscribed(uint64 subId, address consumer);
  error OnlyCallableFromLink();
  error InvalidCalldata();
  error MustBeSubOwner(address owner);
  error MustBeRequestedOwner(address proposedOwner);
  error BalanceInvariantViolated(uint256 internalBalance, uint256 externalBalance); // Should never happen
  event FundsRecovered(address to, uint256 amount);
  // There are only 1e9*1e18 = 1e27 juels in existence, so the balance can fit in uint96 (2^96 ~ 7e28)
  struct Subscription {
    uint96 balance; // Common link balance used for all consumer requests.
    address owner; // Owner can fund/withdraw/cancel the sub.
    address requestedOwner; // For safe transfering sub ownership.
  }
  // We make this public so people can lookup which subscription a given consumer is using (if any).
  mapping(address /* consumer */ => uint64 /* subId */) public s_consumers;
  mapping(uint64 /* subId */ => Subscription /* subscription */) private s_subscriptions;
  uint64 private s_currentSubId;
  // s_totalBalance tracks the total link sent to/from
  // this contract through onTokenTransfer, defundSubscription, cancelSubscription and oracleWithdraw.
  // A discrepancy with this contracts link balance indicates someone
  // sent tokens using transfer and so we may need to use recoverFunds.
  uint96 public s_totalBalance;
  event SubscriptionCreated(uint64 indexed subId, address owner, address[] consumers);
  event SubscriptionFunded(uint64 indexed subId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionConsumerAdded(uint64 indexed subId, address consumer);
  event SubscriptionConsumerRemoved(uint64 indexed subId, address consumer);
  event SubscriptionDefunded(uint64 indexed subId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionCanceled(uint64 indexed subId, address to, uint256 amount);
  event SubscriptionOwnerTransferRequested(uint64 indexed subId, address from, address to);
  event SubscriptionOwnerTransferred(uint64 indexed subId, address from, address to);

  // Set this maximum to 200 to give us a 56 block window to fulfill
  // the request before requiring the block hash feeder.
  uint16 constant MAX_REQUEST_CONFIRMATIONS = 200;
  error InvalidRequestBlockConfs(uint16 have, uint16 min, uint16 max);
  error GasLimitTooBig(uint32 have, uint32 want);
  error KeyHashAlreadyRegistered(bytes32 keyHash);
  error InvalidFeedResponse(int256 linkWei);
  error InsufficientGasForConsumer(uint256 have, uint256 want);
  error InvalidProofLength(uint256 have, uint256 want);
  error NoCorrespondingRequest();
  error IncorrectCommitment();
  error BlockhashNotInStore(uint256 blockNum);
  error PaymentTooLarge();
  //
  uint256 constant private CUSHION = 5_000;
  // Just to relieve stack pressure
  struct FulfillmentParams {
    uint64 subId;
    uint32 callbackGasLimit;
    uint32 numWords;
    address sender;
  }
  mapping(bytes32 /* keyHash */ => address /* oracle */) private s_serviceAgreements;
  mapping(address /* oracle */ => uint96 /* LINK balance */) private s_withdrawableTokens;
  mapping(bytes32 /* keyHash */ => mapping(address /* consumer */ => uint256 /* nonce */)) private s_nonces;
  mapping(uint256 /* requestID */ => bytes32 /* commitment */) private s_commitments;
  event NewServiceAgreement(bytes32 keyHash, address indexed oracle);
  event RandomWordsRequested(
    bytes32 indexed keyHash,
    uint256 preSeedAndRequestId,
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
    uint16 minimumRequestBlockConfirmations;
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
  }
  int256 s_fallbackWeiPerUnitLink;
  Config private s_config;
  event ConfigSet(
    uint16 minimumRequestBlockConfirmations,
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

  function registerProvingKey(
    address oracle,
    uint256[2] calldata publicProvingKey
  )
    external
    onlyOwner()
  {
    bytes32 kh = hashOfKey(publicProvingKey);
    if (s_serviceAgreements[kh] != address(0)) {
      revert KeyHashAlreadyRegistered(kh);
    }
    s_serviceAgreements[kh] = oracle;
    emit NewServiceAgreement(kh, oracle);
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
    return keccak256(abi.encodePacked(publicKey));
  }

  function setConfig(
    uint16 minimumRequestBlockConfirmations,
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
    s_config = Config({
      minimumRequestBlockConfirmations: minimumRequestBlockConfirmations,
      fulfillmentFlatFeeLinkPPM: fulfillmentFlatFeeLinkPPM,
      maxGasLimit: maxGasLimit,
      stalenessSeconds: stalenessSeconds,
      gasAfterPaymentCalculation: gasAfterPaymentCalculation,
      minimumSubscriptionBalance: minimumSubscriptionBalance
    });
  s_fallbackWeiPerUnitLink = fallbackWeiPerUnitLink;
  emit ConfigSet(
      minimumRequestBlockConfirmations,
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
      uint16 minimumRequestBlockConfirmations,
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
      config.minimumRequestBlockConfirmations,
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
  // so in the worse case where the consuming contract has to read all of them
  // from storage, it only has to read 2 words.
  function requestRandomWords(
    bytes32 keyHash,  // Corresponds to a particular offchain job which uses that key for the proofs
    uint64  subId,
    uint16  requestConfirmations,
    uint32  callbackGasLimit,
    uint32  numWords  // Desired number of random words
  )
    external
    returns (
      uint256 requestId
    )
  {
    // Input validation using the subscription storage.
    if (s_subscriptions[subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    // Its important to ensure that the consumer is in fact who they say they
    // are, otherwise they could use someone else's subscription balance.
    if (s_consumers[msg.sender] != subId) {
      revert InvalidConsumer(subId, msg.sender);
    }
    // Input validation using the config storage word.
    if (requestConfirmations < s_config.minimumRequestBlockConfirmations || requestConfirmations > MAX_REQUEST_CONFIRMATIONS) {
      revert InvalidRequestBlockConfs(requestConfirmations, s_config.minimumRequestBlockConfirmations, MAX_REQUEST_CONFIRMATIONS);
    }
    if (s_subscriptions[subId].balance < s_config.minimumSubscriptionBalance) {
      revert InsufficientBalance();
    }
    if (callbackGasLimit > s_config.maxGasLimit) {
      revert GasLimitTooBig(callbackGasLimit, s_config.maxGasLimit);
    }
    // We could additionally check s_serviceAgreements[keyHash] != address(0)
    // but that would require reading another word of storage. To save gas
    // we leave that out.
    uint256 nonce = s_nonces[keyHash][msg.sender] + 1;
    uint256 preSeedAndRequestId = uint256(keccak256(abi.encode(keyHash, msg.sender, nonce)));

    s_commitments[preSeedAndRequestId] = keccak256(abi.encodePacked(
        preSeedAndRequestId,
        block.number,
        subId,
        callbackGasLimit,
        numWords,
        msg.sender));
    emit RandomWordsRequested(keyHash, preSeedAndRequestId, subId, requestConfirmations, callbackGasLimit, numWords, msg.sender);
    s_nonces[keyHash][msg.sender] = nonce;

    return preSeedAndRequestId;
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
    return s_commitments[requestId];
  }

  /**
   * @dev calls target address with exactly gasAmount gas and data as calldata
   * or reverts if at least gasAmount gas is not available
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
    assembly{
      let g := gas()
      // Compute g -= CUSHION and check for underflow
      if lt(g, CUSHION) { revert(0, 0) }
      g := sub(g, CUSHION)
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

  // Offsets into fulfillRandomnessRequest's proof of various values
  //
  // Public key. Skips byte array's length prefix.
  uint256 private constant PUBLIC_KEY_OFFSET = 0x20;
  // Seed is 7th word in proof, plus word for length, (6+1)*0x20=0xe0
  uint256 private constant PRESEED_OFFSET = 7*0x20;

  function fulfillRandomWords(
    bytes memory proof
  )
    external
  {
    uint256 startGas = gasleft();
    (
      bytes32 keyHash,
      uint256 requestId,
      uint256 randomness,
      FulfillmentParams memory fp
    ) = getRandomnessFromProof(proof);


    uint256[] memory randomWords = new uint256[](fp.numWords);
    for (uint256 i = 0; i < fp.numWords; i++) {
      randomWords[i] = uint256(keccak256(abi.encode(randomness, i)));
    }

    // Prevent re-entrancy. The user callback cannot call fulfillRandomWords again
    // with the same proof because this getRandomnessFromProof will revert because the requestId
    // is gone.
    delete s_commitments[requestId];
    VRFConsumerV2Interface v;
    bytes memory resp = abi.encodeWithSelector(v.fulfillRandomWords.selector, requestId, randomWords);
    uint256 gasPreCallback = gasleft();
    if (gasPreCallback < fp.callbackGasLimit) {
      revert InsufficientGasForConsumer(gasPreCallback, fp.callbackGasLimit);
    }
    // Call with explicitly the amount of callback gas requested
    // Important to not let them exhaust the gas budget and avoid oracle payment.
    bool success = callWithExactGas(fp.callbackGasLimit, fp.sender, resp);
    emit RandomWordsFulfilled(requestId, randomWords, success);

    // We want to charge users exactly for how much gas they use in their callback.
    // The gasAfterPaymentCalculation is meant to cover these additional operations where we
    // decrement the subscription balance and increment the oracles withdrawable balance.
    // We also add the flat link fee to the payment amount.
    // Its specified in millionths of link, if s_config.fulfillmentFlatFeeLinkPPM = 1
    // 1 link / 1e6 = 1e18 juels / 1e6 = 1e12 juels.
    uint96 payment = calculatePaymentAmount(startGas, s_config.gasAfterPaymentCalculation, s_config.fulfillmentFlatFeeLinkPPM, tx.gasprice);
    if (s_subscriptions[fp.subId].balance < payment) {
      revert InsufficientBalance();
    }
    s_subscriptions[fp.subId].balance -= payment;
    s_withdrawableTokens[s_serviceAgreements[keyHash]] += payment;
  }

  // Get the amount of gas used for fulfillment
  function calculatePaymentAmount(
      uint256 startGas,
      uint256 gasAfterPaymentCalculation,
      uint32  fulfillmentFlatFeeLinkPPM,
      uint256 weiPerUnitGas
  )
    private
    view
    returns (
      uint96
    )
  {
    int256 weiPerUnitLink;
    weiPerUnitLink = getFeedData();
    if (weiPerUnitLink < 0) {
      revert InvalidFeedResponse(weiPerUnitLink);
    }
    // (1e18 juels/link) (wei/gas * gas) / (wei/link) = jules
    uint256 paymentNoFee = 1e18*weiPerUnitGas*(gasAfterPaymentCalculation + startGas - gasleft()) / uint256(weiPerUnitLink);
    uint256 fee = 1e12*uint256(fulfillmentFlatFeeLinkPPM);
    if (paymentNoFee > (1e27-fee)) {
      revert PaymentTooLarge(); // Payment + fee cannot be more than all of the link in existence.
    }
    return uint96(paymentNoFee+fee);
  }

  function getRandomnessFromProof(
    bytes memory proof
  )
    private
    view 
    returns (
      bytes32 currentKeyHash,
      uint256 requestId, 
      uint256 randomness, 
      FulfillmentParams memory fp
    ) 
  {
    // blockNum follows proof, which follows length word (only direct-number
    // constants are allowed in assembly, so have to compute this in code)
    uint256 blockNumOffset = 0x20 + PROOF_LENGTH;
    // Note that proof.length skips the initial length word.
    // We expect the total length to be proof + 6 words
    // (blocknum, subId, callbackLimit, nw, sender)
    if (proof.length != PROOF_LENGTH + 0x20*5) {
      revert InvalidProofLength(proof.length, PROOF_LENGTH + 0x20*5);
    }
    uint256[2] memory publicKey;
    uint256 preSeed;
    uint256 blockNum;
    address sender;
    assembly { // solhint-disable-line no-inline-assembly
      publicKey := add(proof, PUBLIC_KEY_OFFSET)
      preSeed := mload(add(proof, PRESEED_OFFSET))
      blockNum := mload(add(proof, blockNumOffset))
      // We use a struct to limit local variables to avoid stack depth errors.
      mstore(fp, mload(add(add(proof, blockNumOffset), 0x20))) // subId
      mstore(add(fp, 0x20), mload(add(add(proof, blockNumOffset), 0x40))) // callbackGasLimit
      mstore(add(fp, 0x40), mload(add(add(proof, blockNumOffset), 0x60))) // numWords
      sender := mload(add(add(proof, blockNumOffset), 0x80))
    }
    currentKeyHash = hashOfKey(publicKey);
    bytes32 commitment = s_commitments[preSeed];
    requestId = preSeed;
    if (commitment == 0) {
      revert NoCorrespondingRequest();
    }
    if (commitment != keccak256(abi.encodePacked(
        requestId,
        blockNum,
        fp.subId,
        fp.callbackGasLimit,
        fp.numWords,
        sender)))
    {
      revert IncorrectCommitment();
    }
    fp.sender = sender;

    bytes32 blockHash = blockhash(blockNum);
    if (blockHash == bytes32(0)) {
      blockHash = BLOCKHASH_STORE.getBlockhash(blockNum);
      if (blockHash == bytes32(0)) {
        revert BlockhashNotInStore(blockNum);
      }
    }
    // The seed actually used by the VRF machinery, mixing in the blockhash
    uint256 actualSeed = uint256(keccak256(abi.encodePacked(preSeed, blockHash)));
    // solhint-disable-next-line no-inline-assembly
    assembly { // Construct the actual proof from the remains of proof
      mstore(add(proof, PRESEED_OFFSET), actualSeed)
      mstore(proof, PROOF_LENGTH)
    }
    randomness = VRF.randomValueFromVRFProof(proof); // Reverts on failure
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
      address owner
    )
  {
    return (
      s_subscriptions[subId].balance,
      s_subscriptions[subId].owner
    );
  }

  function createSubscription(
    address[] memory consumers // permitted consumers of the subscription
  )
    external
    returns (
      uint64
    )
  {
    s_currentSubId++;
    uint64 currentSubId = s_currentSubId;
    s_subscriptions[currentSubId] = Subscription({
      balance: 0,
      owner: msg.sender,
      requestedOwner: address(0)
    });
    // This is ok to run out of gas, simply limits the number of
    // consumers that can be added at subscription time.
    // (more can be added with addConsumer).
    for (uint256 i; i < consumers.length; i++) {
      s_consumers[consumers[i]] = currentSubId;
    }
    emit SubscriptionCreated(currentSubId, msg.sender, consumers);
    return currentSubId;
  }

  function requestSubscriptionOwnerTransfer(
    uint64 subId,
    address newOwner
  )
    external
    onlySubOwner(subId)
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
  {
    if (s_consumers[consumer] != subId) {
      revert InvalidConsumer(subId, consumer);
    }
    delete s_consumers[consumer];
    emit SubscriptionConsumerRemoved(subId, consumer);
  }

  function addConsumer(
    uint64 subId,
    address consumer
  )
    external
    onlySubOwner(subId)
  {
    // Must explicitly remove a consumer before changing its subscription.
    if (s_consumers[consumer] != 0) {
      revert AlreadySubscribed(subId, consumer);
    }
    s_consumers[consumer] = subId;
    emit SubscriptionConsumerAdded(subId, consumer);
  }

  function defundSubscription(
    uint64 subId,
    address to,
    uint96 amount
  )
    external
    onlySubOwner(subId)
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
  {
    uint96 balance = s_subscriptions[subId].balance;
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
