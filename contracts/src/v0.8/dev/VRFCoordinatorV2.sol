// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/BlockHashStoreInterface.sol";
import "../interfaces/AggregatorV3Interface.sol";
import "../interfaces/TypeAndVersionInterface.sol";

import "./VRF.sol";
import "./ConfirmedOwner.sol";
import "../interfaces/VRFConsumerV2Interface.sol";

contract VRFCoordinatorV2 is VRF, ConfirmedOwner, TypeAndVersionInterface {

  LinkTokenInterface public immutable LINK;
  AggregatorV3Interface public immutable LINK_ETH_FEED;
  BlockHashStoreInterface public immutable BLOCKHASH_STORE;

  error InsufficientBalance();
  error InvalidConsumer(address consumer);
  error InvalidSubscription();
  error MustBeSubOwner(address owner);
  struct Subscription {
    uint96 balance; // Common link balance used for all consumer requests.
    address owner; // Owner can fund/withdraw/cancel the sub
    address[] consumers; // List of addresses which can consume using this subscription.
  }
  mapping(uint64 /* subId */ => Subscription /* subscription */) private s_subscriptions;
  uint64 private currentSubId;
  event SubscriptionCreated(uint64 subId, address owner, address[] consumers);
  event SubscriptionFundsAdded(uint64 subId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionConsumersUpdated(uint64 subId, address[] oldConsumers, address[] newConsumers);
  event SubscriptionFundsWithdrawn(uint64 subId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionCanceled(uint64 subId, address to, uint256 amount);

  error RequestBlockConfsTooLow(uint64 have, uint64 want);
  error KeyHashAlreadyRegistered(bytes32 keyHash);
  error InvalidFeedResponse(int256 linkWei);
  error InsufficientGasForConsumer(uint256 have, uint256 want);
  error InvalidProofLength(uint256 have, uint256 want);
  error NoCorrespondingRequest();
  error IncorrectCommitment();
  error BlockHashNotInStore(uint256 blockNum);
  // Just to relieve stack pressure
  struct FulfillmentParams {
    uint64 subId;
    uint64 callbackGasLimit;
    uint64 numWords;
    address sender;
  }
  mapping(bytes32 /* keyHash */ => address /* oracle */) private s_serviceAgreements;
  mapping(address /* oracle */ => uint96 /* LINK balance */) private s_withdrawableTokens;
  mapping(bytes32 /* keyHash */ => mapping(address /* consumer */ => uint256 /* nonce */)) public s_nonces;
  mapping(uint256 /* requestID */ => bytes32) private s_callbacks;
  event NewServiceAgreement(bytes32 keyHash, address oracle);
  event RandomWordsRequested(
    bytes32 indexed keyHash,
    uint256 preSeedAndRequestId,
    uint64 subId,
    uint64 minimumRequestConfirmations,
    uint64 callbackGasLimit,
    uint64 numWords,
    address sender
  );
  event RandomWordsFulfilled(
    uint256 requestId,
    uint256[] output,
    bool success
  );

  struct Config {
    // Gas to cover oracle payment after we calculate the payment.
    // We make it configurable in case those operations are repriced.
    uint32 gasAfterPaymentCalculation;
    uint32 stalenessSeconds;
    uint16 minimumRequestBlockConfirmations;
    uint96 minimumSubscriptionBalance;
  }
  Config private s_config;
  int256 private s_fallbackLinkPrice;
  event ConfigSet(
    uint16 minimumRequestBlockConfirmations,
    uint32 stalenessSeconds,
    uint32 gasAfterPaymentCalculation,
    int256 fallbackLinkPrice,
    uint256 minimumSubscriptionBalance
  );

  constructor(
    address link,
    address blockHashStore,
    address linkEthFeed
  )
    ConfirmedOwner(msg.sender)
  {
    LINK = LinkTokenInterface(link);
    LINK_ETH_FEED = AggregatorV3Interface(linkEthFeed);
    BLOCKHASH_STORE = BlockHashStoreInterface(blockHashStore);
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
    uint32 stalenessSeconds,
    uint32 gasAfterPaymentCalculation,
    int256 fallbackLinkPrice,
    uint96 minimumSubscriptionBalance
  )
    external
    onlyOwner()
  {
    s_config = Config({
      minimumRequestBlockConfirmations: minimumRequestBlockConfirmations,
      stalenessSeconds: stalenessSeconds,
      gasAfterPaymentCalculation: gasAfterPaymentCalculation,
      minimumSubscriptionBalance: minimumSubscriptionBalance
    });
    s_fallbackLinkPrice = fallbackLinkPrice;
    emit ConfigSet(minimumRequestBlockConfirmations,
      stalenessSeconds,
      gasAfterPaymentCalculation,
      fallbackLinkPrice,
      minimumSubscriptionBalance
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
      uint32 stalenessSeconds,
      uint32 gasAfterPaymentCalculation,
      int256 fallbackLinkPrice
    )
  {
    Config memory config = s_config;
    return (
      config.minimumRequestBlockConfirmations,
      config.stalenessSeconds,
      config.gasAfterPaymentCalculation,
      s_fallbackLinkPrice
    );
  }

  function requestRandomWords(
    bytes32 keyHash,  // Corresponds to a particular offchain job which uses that key for the proofs
    uint64  subId,
    uint64  minimumRequestConfirmations,
    uint64  callbackGasLimit,
    uint64  numWords,  // Desired number of random words
    uint16  consumerID // Index into consumers to avoid SLOADing all the consumers
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
    // We use this consumer index to ensure that the cost remains constant no matter
    // the number of consumers.
    // Its important to ensure that the consumer is in fact who they say they
    // are, otherwise they could use someone else's subscription balance.
    if (s_subscriptions[subId].consumers[consumerID] != msg.sender) {
      revert InvalidConsumer(msg.sender);
    }
    // Input validation using the config storage word.
    if (minimumRequestConfirmations < s_config.minimumRequestBlockConfirmations) {
      revert RequestBlockConfsTooLow(minimumRequestConfirmations, s_config.minimumRequestBlockConfirmations);
    }
    if (s_subscriptions[subId].balance < s_config.minimumSubscriptionBalance) {
      revert InsufficientBalance();
    }
    // We could additionally check s_serviceAgreements[keyHash] != address(0)
    // but that would require reading another word of storage. To save gas
    // we leave that out.
    uint256 nonce = s_nonces[keyHash][msg.sender] + 1;
    uint256 preSeedAndRequestId = uint256(keccak256(abi.encode(keyHash, msg.sender, nonce)));

    s_callbacks[preSeedAndRequestId] = keccak256(abi.encodePacked(
        preSeedAndRequestId,
        block.number,
        subId,
        callbackGasLimit,
        numWords,
        msg.sender));
    emit RandomWordsRequested(keyHash, preSeedAndRequestId, subId, minimumRequestConfirmations, callbackGasLimit, numWords, msg.sender);
    s_nonces[keyHash][msg.sender] = nonce;

    return preSeedAndRequestId;
  }

  function getCallback(
      uint256 requestId
  )
    external
    view
    returns (
      bytes32
    )
  {
    return s_callbacks[requestId];
  }

  // Offsets into fulfillRandomnessRequest's proof of various values
  //
  // Public key. Skips byte array's length prefix.
  uint256 public constant PUBLIC_KEY_OFFSET = 0x20;
  // Seed is 7th word in proof, plus word for length, (6+1)*0x20=0xe0
  uint256 public constant PRESEED_OFFSET = 0xe0;

  function fulfillRandomWords(
    bytes memory proof
  )
    external
  {
    uint256 startGas = gasleft();
    (bytes32 keyHash, uint256 requestId,
    uint256 randomness, FulfillmentParams memory fp) = getRandomnessFromProof(proof);

    uint256[] memory randomWords = new uint256[](fp.numWords);
    for (uint256 i = 0; i < fp.numWords; i++) {
      randomWords[i] = uint256(keccak256(abi.encode(randomness, i)));
    }

    // Prevent re-entrancy. The user callback cannot call fulfillRandomWords again
    // with the same proof because this getRandomnessFromProof will revert because the requestId
    // is gone.
    delete s_callbacks[requestId];
    VRFConsumerV2Interface v;
    bytes memory resp = abi.encodeWithSelector(v.fulfillRandomWords.selector, requestId, randomWords);
    uint256 gasPreCallback = gasleft();
    if (gasPreCallback < fp.callbackGasLimit) {
      revert InsufficientGasForConsumer(gasPreCallback, fp.callbackGasLimit);
    }
    (bool success,) = fp.sender.call(resp);
    emit RandomWordsFulfilled(requestId, randomWords, success);

    // We want to charge users exactly for how much gas they use in their callback.
    // The gasAfterPaymentCalculation is meant to cover these additional operations where we
    // decrement the subscription balance and increment the oracles withdrawable balance.
    uint96 payment = calculatePaymentAmount(startGas, s_config.gasAfterPaymentCalculation, tx.gasprice);
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
    // (1e18 jules/link) (wei/gas * gas) / (wei/link) = jules
    return uint96(1e18*weiPerUnitGas*(gasAfterPaymentCalculation + startGas - gasleft()) / uint256(weiPerUnitLink));
  }

  function getRandomnessFromProof(
    bytes memory proof
  )
    public 
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
    bytes32 callback = s_callbacks[preSeed];
    requestId = preSeed;
    if (callback == 0) {
      revert NoCorrespondingRequest();
    }
    if (callback != keccak256(abi.encodePacked(
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
        revert BlockHashNotInStore(blockNum);
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
    int256 linkEth;
    (,linkEth,,timestamp,) = LINK_ETH_FEED.latestRoundData();
    if (staleFallback && stalenessSeconds < block.timestamp - timestamp) {
      linkEth = s_fallbackLinkPrice;
    }
    return linkEth;
  }

  function withdraw(
    address recipient,
    uint96 amount
  )
    external
  {
    if (s_withdrawableTokens[msg.sender] < amount) {
      revert InsufficientBalance();
    }
    s_withdrawableTokens[msg.sender] -= amount;
    assert(LINK.transfer(recipient, amount));
  }

  function getSubscription(
    uint64 subId
  )
    external
    view
    returns (
      Subscription memory
    )
  {
    return s_subscriptions[subId];
  }

  function createSubscription(
    address[] memory consumers // permitted consumers of the subscription
  )
    external
    returns (
      uint64
    )
  {
    currentSubId++;
    s_subscriptions[currentSubId] = Subscription({
      owner: msg.sender,
      consumers: consumers,
      balance: 0
    });
    emit SubscriptionCreated(currentSubId, msg.sender, consumers);
    return currentSubId;
  }

  function updateSubscription(
    uint64 subId,
    address[] memory consumers // permitted consumers of the subscription
  )
    external
    onlySubOwner(subId)
  {
    address[] memory oldConsumers = s_subscriptions[subId].consumers;
    s_subscriptions[subId].consumers = consumers;
    emit SubscriptionConsumersUpdated(subId, oldConsumers, consumers);
  }

  function fundSubscription(
    uint64 subId,
    uint96 amount
  )
    external
    onlySubOwner(subId)
  {
    if (s_subscriptions[subId].owner == address(0))  {
      revert InvalidSubscription();
    }
    uint256 oldBalance = s_subscriptions[subId].balance;
    s_subscriptions[subId].balance += amount;
    LINK.transferFrom(msg.sender, address(this), amount);
    emit SubscriptionFundsAdded(subId, oldBalance, oldBalance+amount);
  }

  function withdrawFromSubscription(
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
    LINK.transfer(to, amount);
    emit SubscriptionFundsWithdrawn(subId, oldBalance, s_subscriptions[subId].balance);
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
    uint256 balance = s_subscriptions[subId].balance;
    delete s_subscriptions[subId];
    LINK.transfer(to, balance);
    emit SubscriptionCanceled(subId, to, balance);
  }

  modifier onlySubOwner(uint64 subId) {
    address owner = s_subscriptions[subId].owner;
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
