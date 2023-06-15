// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/BlockhashStoreInterface.sol";
import "../interfaces/AggregatorV3Interface.sol";
import "../interfaces/TypeAndVersionInterface.sol";
import "../interfaces/ERC677ReceiverInterface.sol";
import "./VRF.sol";
import "./VRFConsumerBaseV2.sol";
import "../ChainSpecificUtil.sol";
import "./SubscriptionAPI.sol";

contract VRFCoordinatorV2_5 is
  VRF,
  TypeAndVersionInterface,
  ERC677ReceiverInterface,
  SubscriptionAPI
{
  /// @dev may not be provided upon construction on some chains due to lack of availability
  AggregatorV3Interface public LINK_ETH_FEED;
  /// @dev should always be available
  BlockhashStoreInterface public immutable BLOCKHASH_STORE;

  // Set this maximum to 200 to give us a 56 block window to fulfill
  // the request before requiring the block hash feeder.
  uint16 public constant MAX_REQUEST_CONFIRMATIONS = 200;
  uint32 public constant MAX_NUM_WORDS = 500;
  // 5k is plenty for an EXTCODESIZE call (2600) + warm CALL (100)
  // and some arithmetic operations.
  uint256 private constant GAS_FOR_CALL_EXACT_CHECK = 5_000;
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
    bool nativePayment;
  }
  mapping(bytes32 => address) /* keyHash */ /* oracle */ public s_provingKeys;
  bytes32[] public s_provingKeyHashes;
  mapping(uint256 => bytes32) /* requestID */ /* commitment */ public s_requestCommitments;
  event ProvingKeyRegistered(bytes32 keyHash, address indexed oracle);
  event ProvingKeyDeregistered(bytes32 keyHash, address indexed oracle);
  event RandomWordsRequested(
    bytes32 indexed keyHash,
    uint256 requestId,
    uint256 preSeed,
    uint64 indexed subId,
    uint16 minimumRequestConfirmations,
    uint32 callbackGasLimit,
    uint32 numWords,
    bool nativePayment,
    address indexed sender
  );
  event RandomWordsFulfilled(uint256 indexed requestId, uint256 outputSeed, uint96 payment, bool nativePayment, bool success);

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
    uint32 gasAfterPaymentCalculation;
  }
  int256 public s_fallbackWeiPerUnitLink;
  Config public s_config;
  FeeConfig public s_feeConfig;
  struct FeeConfig {
    // Flat fee charged per fulfillment in millionths of link
    // So fee range is [0, 2^32/10^6].
    uint32 fulfillmentFlatFeeLinkPPMTier1;
    uint32 fulfillmentFlatFeeLinkPPMTier2;
    uint32 fulfillmentFlatFeeLinkPPMTier3;
    uint32 fulfillmentFlatFeeLinkPPMTier4;
    uint32 fulfillmentFlatFeeLinkPPMTier5;
    uint24 reqsForTier2;
    uint24 reqsForTier3;
    uint24 reqsForTier4;
    uint24 reqsForTier5;
  }
  event ConfigSet(
    uint16 minimumRequestConfirmations,
    uint32 maxGasLimit,
    uint32 stalenessSeconds,
    uint32 gasAfterPaymentCalculation,
    int256 fallbackWeiPerUnitLink,
    FeeConfig feeConfig
  );

  constructor(address blockhashStore) SubscriptionAPI() {
    BLOCKHASH_STORE = BlockhashStoreInterface(blockhashStore);
  }

  /**
    * @notice set the link eth feed to be used by this coordinator
    * @param linkEthFeed address of the link eth feed
  */
  function setLinkEthFeed(address linkEthFeed) external onlyOwner {
    LINK_ETH_FEED = AggregatorV3Interface(linkEthFeed);
  }

  /**
   * @notice Registers a proving key to an oracle.
   * @param oracle address of the oracle
   * @param publicProvingKey key that oracle can use to submit vrf fulfillments
   */
  function registerProvingKey(address oracle, uint256[2] calldata publicProvingKey) external onlyOwner {
    bytes32 kh = hashOfKey(publicProvingKey);
    if (s_provingKeys[kh] != address(0)) {
      revert ProvingKeyAlreadyRegistered(kh);
    }
    s_provingKeys[kh] = oracle;
    s_provingKeyHashes.push(kh);
    emit ProvingKeyRegistered(kh, oracle);
  }

  /**
   * @notice Deregisters a proving key to an oracle.
   * @param publicProvingKey key that oracle can use to submit vrf fulfillments
   */
  function deregisterProvingKey(uint256[2] calldata publicProvingKey) external onlyOwner {
    bytes32 kh = hashOfKey(publicProvingKey);
    address oracle = s_provingKeys[kh];
    if (oracle == address(0)) {
      revert NoSuchProvingKey(kh);
    }
    delete s_provingKeys[kh];
    for (uint256 i = 0; i < s_provingKeyHashes.length; i++) {
      if (s_provingKeyHashes[i] == kh) {
        bytes32 last = s_provingKeyHashes[s_provingKeyHashes.length - 1];
        // Copy last element and overwrite kh to be deleted with it
        s_provingKeyHashes[i] = last;
        s_provingKeyHashes.pop();
      }
    }
    emit ProvingKeyDeregistered(kh, oracle);
  }

  /**
   * @notice Returns the proving key hash key associated with this public key
   * @param publicKey the key to return the hash of
   */
  function hashOfKey(uint256[2] memory publicKey) public pure returns (bytes32) {
    return keccak256(abi.encode(publicKey));
  }

  /**
   * @notice Sets the configuration of the vrfv2 coordinator
   * @param minimumRequestConfirmations global min for request confirmations
   * @param maxGasLimit global max for request gas limit
   * @param stalenessSeconds if the eth/link feed is more stale then this, use the fallback price
   * @param gasAfterPaymentCalculation gas used in doing accounting after completing the gas measurement
   * @param fallbackWeiPerUnitLink fallback eth/link price in the case of a stale feed
   * @param feeConfig fee tier configuration
   */
  function setConfig(
    uint16 minimumRequestConfirmations,
    uint32 maxGasLimit,
    uint32 stalenessSeconds,
    uint32 gasAfterPaymentCalculation,
    int256 fallbackWeiPerUnitLink,
    FeeConfig memory feeConfig
  ) external onlyOwner {
    if (minimumRequestConfirmations > MAX_REQUEST_CONFIRMATIONS) {
      revert InvalidRequestConfirmations(
        minimumRequestConfirmations,
        minimumRequestConfirmations,
        MAX_REQUEST_CONFIRMATIONS
      );
    }
    if (fallbackWeiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(fallbackWeiPerUnitLink);
    }
    s_config = Config({
      minimumRequestConfirmations: minimumRequestConfirmations,
      maxGasLimit: maxGasLimit,
      stalenessSeconds: stalenessSeconds,
      gasAfterPaymentCalculation: gasAfterPaymentCalculation,
      reentrancyLock: false
    });
    s_feeConfig = feeConfig;
    s_fallbackWeiPerUnitLink = fallbackWeiPerUnitLink;
    emit ConfigSet(
      minimumRequestConfirmations,
      maxGasLimit,
      stalenessSeconds,
      gasAfterPaymentCalculation,
      fallbackWeiPerUnitLink,
      s_feeConfig
    );
  }

  /**
   * @notice Get configuration relevant for making requests
   * @return minimumRequestConfirmations global min for request confirmations
   * @return maxGasLimit global max for request gas limit
   * @return s_provingKeyHashes list of registered key hashes
   */
  function getRequestConfig() external view override returns (uint16, uint32, bytes32[] memory) {
    return (s_config.minimumRequestConfirmations, s_config.maxGasLimit, s_provingKeyHashes);
  }

  /**
   * @notice Request a set of random words.
   * @param keyHash - Corresponds to a particular oracle job which uses
   * that key for generating the VRF proof. Different keyHash's have different gas price
   * ceilings, so you can select a specific one to bound your maximum per request cost.
   * @param subId  - The ID of the VRF subscription. Must be funded
   * with the minimum subscription balance required for the selected keyHash.
   * @param minimumRequestConfirmations - How many blocks you'd like the
   * oracle to wait before responding to the request. See SECURITY CONSIDERATIONS
   * for why you may want to request more. The acceptable range is
   * [minimumRequestBlockConfirmations, 200].
   * @param callbackGasLimit - How much gas you'd like to receive in your
   * fulfillRandomWords callback. Note that gasleft() inside fulfillRandomWords
   * may be slightly less than this amount because of gas used calling the function
   * (argument decoding etc.), so you may need to request slightly more than you expect
   * to have inside fulfillRandomWords. The acceptable range is
   * [0, maxGasLimit]
   * @param numWords - The number of uint256 random values you'd like to receive
   * in your fulfillRandomWords callback. Note these numbers are expanded in a
   * secure way by the VRFCoordinator from a single random value supplied by the oracle.
   * @param payInEth - Whether payment should be made in ETH or LINK. Payment in LINK
   * is only available if the LINK token is available to this contract.
   * @return requestId - A unique identifier of the request. Can be used to match
   * a request to a response in fulfillRandomWords.
   */
  function requestRandomWords(
    bytes32 keyHash,
    uint64 subId,
    uint16 requestConfirmations,
    uint32 callbackGasLimit,
    uint32 numWords,
    bool payInEth
  ) external override nonReentrant returns (uint256) {
    // Input validation using the subscription storage.
    if (s_subscriptionConfigs[subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    // Its important to ensure that the consumer is in fact who they say they
    // are, otherwise they could use someone else's subscription balance.
    // A nonce of 0 indicates consumer is not allocated to the sub.
    uint64 currentNonce = s_consumers[msg.sender][subId];
    if (currentNonce == 0) {
      revert InvalidConsumer(subId, msg.sender);
    }
    // Input validation using the config storage word.
    if (
      requestConfirmations < s_config.minimumRequestConfirmations || requestConfirmations > MAX_REQUEST_CONFIRMATIONS
    ) {
      revert InvalidRequestConfirmations(
        requestConfirmations,
        s_config.minimumRequestConfirmations,
        MAX_REQUEST_CONFIRMATIONS
      );
    }
    // No lower bound on the requested gas limit. A user could request 0
    // and they would simply be billed for the proof verification and wouldn't be
    // able to do anything with the random value.
    if (callbackGasLimit > s_config.maxGasLimit) {
      revert GasLimitTooBig(callbackGasLimit, s_config.maxGasLimit);
    }
    if (numWords > MAX_NUM_WORDS) {
      revert NumWordsTooBig(numWords, MAX_NUM_WORDS);
    }
    // Note we do not check whether the keyHash is valid to save gas.
    // The consequence for users is that they can send requests
    // for invalid keyHashes which will simply not be fulfilled.
    uint64 nonce = currentNonce + 1;
    (uint256 requestId, uint256 preSeed) = computeRequestId(keyHash, msg.sender, subId, nonce);

    s_requestCommitments[requestId] = keccak256(
      abi.encode(requestId, ChainSpecificUtil.getBlockNumber(), subId, callbackGasLimit, numWords, msg.sender)
    );
    emit RandomWordsRequested(
      keyHash,
      requestId,
      preSeed,
      subId,
      requestConfirmations,
      callbackGasLimit,
      numWords,
      payInEth,
      msg.sender
    );
    s_consumers[msg.sender][subId] = nonce;

    return requestId;
  }

  function computeRequestId(
    bytes32 keyHash,
    address sender,
    uint64 subId,
    uint64 nonce
  ) private pure returns (uint256, uint256) {
    uint256 preSeed = uint256(keccak256(abi.encode(keyHash, sender, subId, nonce)));
    return (uint256(keccak256(abi.encode(keyHash, preSeed))), preSeed);
  }

  /**
   * @dev calls target address with exactly gasAmount gas and data as calldata
   * or reverts if at least gasAmount gas is not available.
   */
  function callWithExactGas(uint256 gasAmount, address target, bytes memory data) private returns (bool success) {
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

  function getRandomnessFromProof(
    Proof memory proof,
    RequestCommitment memory rc
  ) private view returns (bytes32 keyHash, uint256 requestId, uint256 randomness) {
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
    if (
      commitment != keccak256(abi.encode(requestId, rc.blockNum, rc.subId, rc.callbackGasLimit, rc.numWords, rc.sender))
    ) {
      revert IncorrectCommitment();
    }

    bytes32 blockHash = ChainSpecificUtil.getBlockhash(rc.blockNum);
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

  /*
   * @notice Compute fee based on the request count
   * @param reqCount number of requests
   * @return feePPM fee in LINK PPM
   */
  function getFeeTier(uint64 reqCount) public view returns (uint32) {
    FeeConfig memory fc = s_feeConfig;
    if (0 <= reqCount && reqCount <= fc.reqsForTier2) {
      return fc.fulfillmentFlatFeeLinkPPMTier1;
    }
    if (fc.reqsForTier2 < reqCount && reqCount <= fc.reqsForTier3) {
      return fc.fulfillmentFlatFeeLinkPPMTier2;
    }
    if (fc.reqsForTier3 < reqCount && reqCount <= fc.reqsForTier4) {
      return fc.fulfillmentFlatFeeLinkPPMTier3;
    }
    if (fc.reqsForTier4 < reqCount && reqCount <= fc.reqsForTier5) {
      return fc.fulfillmentFlatFeeLinkPPMTier4;
    }
    return fc.fulfillmentFlatFeeLinkPPMTier5;
  }

  /*
   * @notice Fulfill a randomness request
   * @param proof contains the proof and randomness
   * @param rc request commitment pre-image, committed to at request time
   * @return payment amount billed to the subscription
   * @dev simulated offchain to determine if sufficient balance is present to fulfill the request
   */
  function fulfillRandomWords(Proof memory proof, RequestCommitment memory rc) external nonReentrant returns (uint96) {
    uint256 startGas = gasleft();
    (bytes32 keyHash, uint256 requestId, uint256 randomness) = getRandomnessFromProof(proof, rc);

    uint256[] memory randomWords = new uint256[](rc.numWords);
    for (uint256 i = 0; i < rc.numWords; i++) {
      randomWords[i] = uint256(keccak256(abi.encode(randomness, i)));
    }

    delete s_requestCommitments[requestId];
    VRFConsumerBaseV2 v;
    bytes memory resp = abi.encodeWithSelector(v.rawFulfillRandomWords.selector, requestId, randomWords);
    // Call with explicitly the amount of callback gas requested
    // Important to not let them exhaust the gas budget and avoid oracle payment.
    // Do not allow any non-view/non-pure coordinator functions to be called
    // during the consumers callback code via reentrancyLock.
    // Note that callWithExactGas will revert if we do not have sufficient gas
    // to give the callee their requested amount.
    s_config.reentrancyLock = true;
    bool success = callWithExactGas(rc.callbackGasLimit, rc.sender, resp);
    s_config.reentrancyLock = false;

    // Increment the req count for fee tier selection.
    uint64 reqCount = s_subscriptions[rc.subId].reqCount;
    s_subscriptions[rc.subId].reqCount += 1;

    uint96 payment;
    if (rc.nativePayment) {
      // We want to charge users exactly for how much gas they use in their callback.
      // The gasAfterPaymentCalculation is meant to cover these additional operations where we
      // decrement the subscription balance and increment the oracles withdrawable balance.
      // We also add the flat eth fee to the payment amount.
      payment = calculatePaymentAmountNative(
        startGas,
        s_config.gasAfterPaymentCalculation,
        getFeeTier(reqCount),
        tx.gasprice
      );
      if (s_subscriptions[rc.subId].ethBalance < payment) {
        revert InsufficientBalance();
      }
      s_subscriptions[rc.subId].ethBalance -= payment;
      s_withdrawableEth[s_provingKeys[keyHash]] += payment;
    } else {
      // We want to charge users exactly for how much gas they use in their callback.
      // The gasAfterPaymentCalculation is meant to cover these additional operations where we
      // decrement the subscription balance and increment the oracles withdrawable balance.
      // We also add the flat link fee to the payment amount.
      // Its specified in millionths of link, if s_config.fulfillmentFlatFeeLinkPPM = 1
      // 1 link / 1e6 = 1e18 juels / 1e6 = 1e12 juels.
      payment = calculatePaymentAmount(
        startGas,
        s_config.gasAfterPaymentCalculation,
        getFeeTier(reqCount),
        tx.gasprice
      );
      if (s_subscriptions[rc.subId].balance < payment) {
        revert InsufficientBalance();
      }
      s_subscriptions[rc.subId].balance -= payment;
      s_withdrawableTokens[s_provingKeys[keyHash]] += payment;
    }

    // Include payment in the event for tracking costs.
    emit RandomWordsFulfilled(requestId, payment, rc.nativePayment, rc.subId, rc.sender);

    return payment;
  }

  function calculatePaymentAmount(
    uint256 startGas,
    uint256 gasAfterPaymentCalculation,
    uint32 fulfillmentFlatFeePPM,
    uint256 weiPerUnitGas,
    bool nativePayment
  ) internal view returns (uint96) {
    if (nativePayment) {
      return calculatePaymentAmountEth(
        startGas,
        gasAfterPaymentCalculation,
        fulfillmentFlatFeePPM,
        weiPerUnitGas
      );
    }
    return calculatePaymentAmountLink(
      startGas,
      gasAfterPaymentCalculation,
      fulfillmentFlatFeePPM,
      weiPerUnitGas
    );
  }

  function calculatePaymentAmountEth(
    uint256 startGas,
    uint256 gasAfterPaymentCalculation,
    uint32 fulfillmentFlatFeePPM,
    uint256 weiPerUnitGas
  ) internal view returns (uint96) {
    // Will return non-zero on chains that have this enabled
    uint256 l1CostWei = ChainSpecificUtil.getCurrentTxL1GasFees();
    // calculate the payment without the premium
    uint256 baseFeeWei = weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft());
    // calculate the flat fee in wei
    uint256 fee = 1e12 * uint256(fulfillmentFlatFeePPM);
    // return the final fee with the flat fee and l1 cost (if applicable) added
    return uint96(paymentNoFee + fee + l1CostWei);
  }

  // Get the amount of gas used for fulfillment
  function calculatePaymentAmountLink(
    uint256 startGas,
    uint256 gasAfterPaymentCalculation,
    uint32 fulfillmentFlatFeeLinkPPM,
    uint256 weiPerUnitGas
  ) internal view returns (uint96) {
    int256 weiPerUnitLink;
    weiPerUnitLink = getFeedData();
    if (weiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(weiPerUnitLink);
    }
    // Will return non-zero on chains that have this enabled
    uint256 l1CostWei = ChainSpecificUtil.getCurrentTxL1GasFees();
    // (1e18 juels/link) ((wei/gas * gas) + l1wei) / (wei/link) = juels
    uint256 paymentNoFee = (1e18 * (weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft()) + l1CostWei)) /
      uint256(weiPerUnitLink);
    uint256 fee = 1e12 * uint256(fulfillmentFlatFeeLinkPPM);
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
   * @notice Check to see if there exists a request commitment consumers
   * for all consumers and keyhashes for a given sub.
   * @param subId - ID of the subscription
   * @return true if there exists at least one unfulfilled request for the subscription, false
   * otherwise.
   * @dev Looping is bounded to MAX_CONSUMERS*(number of keyhashes).
   * @dev Used to disable subscription canceling while outstanding request are present.
   */
  function pendingRequestExists(uint64 subId) public view override returns (bool) {
    SubscriptionConfig memory subConfig = s_subscriptionConfigs[subId];
    for (uint256 i = 0; i < subConfig.consumers.length; i++) {
      for (uint256 j = 0; j < s_provingKeyHashes.length; j++) {
        (uint256 reqId, ) = computeRequestId(
          s_provingKeyHashes[j],
          subConfig.consumers[i],
          subId,
          s_consumers[subConfig.consumers[i]][subId]
        );
        if (s_requestCommitments[reqId] != 0) {
          return true;
        }
      }
    }
    return false;
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
    return "VRFCoordinatorV2 2.0.0";
  }
}
