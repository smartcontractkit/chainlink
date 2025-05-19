// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {BlockhashStoreInterface} from "../interfaces/BlockhashStoreInterface.sol";
import {VRF} from "../../vrf/VRF.sol";
import {VRFTypes} from "../VRFTypes.sol";
import {VRFConsumerBaseV2Plus, IVRFMigratableConsumerV2Plus} from "./VRFConsumerBaseV2Plus.sol";
import {SubscriptionAPI} from "./SubscriptionAPI.sol";
import {VRFV2PlusClient} from "./libraries/VRFV2PlusClient.sol";
import {IVRFCoordinatorV2PlusMigration} from "./interfaces/IVRFCoordinatorV2PlusMigration.sol";
// solhint-disable-next-line no-unused-import
import {IVRFCoordinatorV2Plus, IVRFSubscriptionV2Plus} from "./interfaces/IVRFCoordinatorV2Plus.sol";

contract VRFCoordinatorV2_5 is VRF, SubscriptionAPI, IVRFCoordinatorV2Plus {
  /// @dev should always be available
  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  BlockhashStoreInterface public immutable BLOCKHASH_STORE;

  // Set this maximum to 200 to give us a 56 block window to fulfill
  // the request before requiring the block hash feeder.
  uint16 public constant MAX_REQUEST_CONFIRMATIONS = 200;
  uint32 public constant MAX_NUM_WORDS = 500;
  // 5k is plenty for an EXTCODESIZE call (2600) + warm CALL (100)
  // and some arithmetic operations.
  uint256 private constant GAS_FOR_CALL_EXACT_CHECK = 5_000;
  // upper bound limit for premium percentages to make sure fee calculations don't overflow
  uint8 private constant PREMIUM_PERCENTAGE_MAX = 155;
  error InvalidRequestConfirmations(uint16 have, uint16 min, uint16 max);
  error GasLimitTooBig(uint32 have, uint32 want);
  error NumWordsTooBig(uint32 have, uint32 want);
  error MsgDataTooBig(uint256 have, uint32 max);
  error ProvingKeyAlreadyRegistered(bytes32 keyHash);
  error NoSuchProvingKey(bytes32 keyHash);
  error InvalidLinkWeiPrice(int256 linkWei);
  error LinkDiscountTooHigh(uint32 flatFeeLinkDiscountPPM, uint32 flatFeeNativePPM);
  error InvalidPremiumPercentage(uint8 premiumPercentage, uint8 max);
  error NoCorrespondingRequest();
  error IncorrectCommitment();
  error BlockhashNotInStore(uint256 blockNum);
  error PaymentTooLarge();
  error InvalidExtraArgsTag();
  error GasPriceExceeded(uint256 gasPrice, uint256 maxGas);

  struct ProvingKey {
    bool exists; // proving key exists
    uint64 maxGas; // gas lane max gas price for fulfilling requests
  }

  mapping(bytes32 => ProvingKey) /* keyHash */ /* provingKey */ public s_provingKeys;
  bytes32[] public s_provingKeyHashes;
  mapping(uint256 => bytes32) /* requestID */ /* commitment */ public s_requestCommitments;
  event ProvingKeyRegistered(bytes32 keyHash, uint64 maxGas);
  event ProvingKeyDeregistered(bytes32 keyHash, uint64 maxGas);

  event RandomWordsRequested(
    bytes32 indexed keyHash,
    uint256 requestId,
    uint256 preSeed,
    uint256 indexed subId,
    uint16 minimumRequestConfirmations,
    uint32 callbackGasLimit,
    uint32 numWords,
    bytes extraArgs,
    address indexed sender
  );

  event RandomWordsFulfilled(
    uint256 indexed requestId,
    uint256 outputSeed,
    uint256 indexed subId,
    uint96 payment,
    bool nativePayment,
    bool success,
    bool onlyPremium
  );

  event L1GasFee(uint256 fee);

  int256 public s_fallbackWeiPerUnitLink;

  event ConfigSet(
    uint16 minimumRequestConfirmations,
    uint32 maxGasLimit,
    uint32 stalenessSeconds,
    uint32 gasAfterPaymentCalculation,
    int256 fallbackWeiPerUnitLink,
    uint32 fulfillmentFlatFeeNativePPM,
    uint32 fulfillmentFlatFeeLinkDiscountPPM,
    uint8 nativePremiumPercentage,
    uint8 linkPremiumPercentage
  );

  event FallbackWeiPerUnitLinkUsed(uint256 requestId, int256 fallbackWeiPerUnitLink);

  constructor(address blockhashStore) SubscriptionAPI() {
    BLOCKHASH_STORE = BlockhashStoreInterface(blockhashStore);
  }

  /**
   * @notice Registers a proving key to.
   * @param publicProvingKey key that oracle can use to submit vrf fulfillments
   */
  function registerProvingKey(uint256[2] calldata publicProvingKey, uint64 maxGas) external onlyOwner {
    bytes32 kh = hashOfKey(publicProvingKey);
    if (s_provingKeys[kh].exists) {
      revert ProvingKeyAlreadyRegistered(kh);
    }
    s_provingKeys[kh] = ProvingKey({exists: true, maxGas: maxGas});
    s_provingKeyHashes.push(kh);
    emit ProvingKeyRegistered(kh, maxGas);
  }

  /**
   * @notice Deregisters a proving key.
   * @param publicProvingKey key that oracle can use to submit vrf fulfillments
   */
  function deregisterProvingKey(uint256[2] calldata publicProvingKey) external onlyOwner {
    bytes32 kh = hashOfKey(publicProvingKey);
    ProvingKey memory key = s_provingKeys[kh];
    if (!key.exists) {
      revert NoSuchProvingKey(kh);
    }
    delete s_provingKeys[kh];
    uint256 s_provingKeyHashesLength = s_provingKeyHashes.length;
    for (uint256 i = 0; i < s_provingKeyHashesLength; ++i) {
      if (s_provingKeyHashes[i] == kh) {
        // Copy last element and overwrite kh to be deleted with it
        s_provingKeyHashes[i] = s_provingKeyHashes[s_provingKeyHashesLength - 1];
        s_provingKeyHashes.pop();
        break;
      }
    }
    emit ProvingKeyDeregistered(kh, key.maxGas);
  }

  /**
   * @notice Returns the proving key hash key associated with this public key
   * @param publicKey the key to return the hash of
   */
  function hashOfKey(uint256[2] calldata publicKey) public pure returns (bytes32) {
    return keccak256(abi.encode(publicKey));
  }

  /**
   * @notice Sets the configuration of the vrfv2 coordinator
   * @param minimumRequestConfirmations global min for request confirmations
   * @param maxGasLimit global max for request gas limit
   * @param stalenessSeconds if the native/link feed is more stale then this, use the fallback price
   * @param gasAfterPaymentCalculation gas used in doing accounting after completing the gas measurement
   * @param fallbackWeiPerUnitLink fallback native/link price in the case of a stale feed
   * @param fulfillmentFlatFeeNativePPM flat fee in native for native payment
   * @param fulfillmentFlatFeeLinkDiscountPPM flat fee discount for link payment in native
   * @param nativePremiumPercentage native premium percentage
   * @param linkPremiumPercentage link premium percentage
   */
  function setConfig(
    uint16 minimumRequestConfirmations,
    uint32 maxGasLimit,
    uint32 stalenessSeconds,
    uint32 gasAfterPaymentCalculation,
    int256 fallbackWeiPerUnitLink,
    uint32 fulfillmentFlatFeeNativePPM,
    uint32 fulfillmentFlatFeeLinkDiscountPPM,
    uint8 nativePremiumPercentage,
    uint8 linkPremiumPercentage
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
    if (fulfillmentFlatFeeLinkDiscountPPM > fulfillmentFlatFeeNativePPM) {
      revert LinkDiscountTooHigh(fulfillmentFlatFeeLinkDiscountPPM, fulfillmentFlatFeeNativePPM);
    }
    if (nativePremiumPercentage > PREMIUM_PERCENTAGE_MAX) {
      revert InvalidPremiumPercentage(nativePremiumPercentage, PREMIUM_PERCENTAGE_MAX);
    }
    if (linkPremiumPercentage > PREMIUM_PERCENTAGE_MAX) {
      revert InvalidPremiumPercentage(linkPremiumPercentage, PREMIUM_PERCENTAGE_MAX);
    }
    s_config = Config({
      minimumRequestConfirmations: minimumRequestConfirmations,
      maxGasLimit: maxGasLimit,
      stalenessSeconds: stalenessSeconds,
      gasAfterPaymentCalculation: gasAfterPaymentCalculation,
      reentrancyLock: false,
      fulfillmentFlatFeeNativePPM: fulfillmentFlatFeeNativePPM,
      fulfillmentFlatFeeLinkDiscountPPM: fulfillmentFlatFeeLinkDiscountPPM,
      nativePremiumPercentage: nativePremiumPercentage,
      linkPremiumPercentage: linkPremiumPercentage
    });
    s_fallbackWeiPerUnitLink = fallbackWeiPerUnitLink;
    emit ConfigSet(
      minimumRequestConfirmations,
      maxGasLimit,
      stalenessSeconds,
      gasAfterPaymentCalculation,
      fallbackWeiPerUnitLink,
      fulfillmentFlatFeeNativePPM,
      fulfillmentFlatFeeLinkDiscountPPM,
      nativePremiumPercentage,
      linkPremiumPercentage
    );
  }

  /// @dev Convert the extra args bytes into a struct
  /// @param extraArgs The extra args bytes
  /// @return The extra args struct
  function _fromBytes(bytes calldata extraArgs) internal pure returns (VRFV2PlusClient.ExtraArgsV1 memory) {
    if (extraArgs.length == 0) {
      return VRFV2PlusClient.ExtraArgsV1({nativePayment: false});
    }
    if (bytes4(extraArgs) != VRFV2PlusClient.EXTRA_ARGS_V1_TAG) revert InvalidExtraArgsTag();
    return abi.decode(extraArgs[4:], (VRFV2PlusClient.ExtraArgsV1));
  }

  /**
   * @notice Request a set of random words.
   * @param req - a struct containing following fiels for randomness request:
   * keyHash - Corresponds to a particular oracle job which uses
   * that key for generating the VRF proof. Different keyHash's have different gas price
   * ceilings, so you can select a specific one to bound your maximum per request cost.
   * subId  - The ID of the VRF subscription. Must be funded
   * with the minimum subscription balance required for the selected keyHash.
   * requestConfirmations - How many blocks you'd like the
   * oracle to wait before responding to the request. See SECURITY CONSIDERATIONS
   * for why you may want to request more. The acceptable range is
   * [minimumRequestBlockConfirmations, 200].
   * callbackGasLimit - How much gas you'd like to receive in your
   * fulfillRandomWords callback. Note that gasleft() inside fulfillRandomWords
   * may be slightly less than this amount because of gas used calling the function
   * (argument decoding etc.), so you may need to request slightly more than you expect
   * to have inside fulfillRandomWords. The acceptable range is
   * [0, maxGasLimit]
   * numWords - The number of uint256 random values you'd like to receive
   * in your fulfillRandomWords callback. Note these numbers are expanded in a
   * secure way by the VRFCoordinator from a single random value supplied by the oracle.
   * extraArgs - Encoded extra arguments that has a boolean flag for whether payment
   * should be made in native or LINK. Payment in LINK is only available if the LINK token is available to this contract.
   * @return requestId - A unique identifier of the request. Can be used to match
   * a request to a response in fulfillRandomWords.
   */
  function requestRandomWords(
    VRFV2PlusClient.RandomWordsRequest calldata req
  ) external override nonReentrant returns (uint256 requestId) {
    // Input validation using the subscription storage.
    uint256 subId = req.subId;
    _requireValidSubscription(s_subscriptionConfigs[subId].owner);
    // Its important to ensure that the consumer is in fact who they say they
    // are, otherwise they could use someone else's subscription balance.
    mapping(uint256 => ConsumerConfig) storage consumerConfigs = s_consumers[msg.sender];
    ConsumerConfig memory consumerConfig = consumerConfigs[subId];
    if (!consumerConfig.active) {
      revert InvalidConsumer(subId, msg.sender);
    }
    // Input validation using the config storage word.
    if (
      req.requestConfirmations < s_config.minimumRequestConfirmations ||
      req.requestConfirmations > MAX_REQUEST_CONFIRMATIONS
    ) {
      revert InvalidRequestConfirmations(
        req.requestConfirmations,
        s_config.minimumRequestConfirmations,
        MAX_REQUEST_CONFIRMATIONS
      );
    }
    // No lower bound on the requested gas limit. A user could request 0
    // and they would simply be billed for the proof verification and wouldn't be
    // able to do anything with the random value.
    if (req.callbackGasLimit > s_config.maxGasLimit) {
      revert GasLimitTooBig(req.callbackGasLimit, s_config.maxGasLimit);
    }
    if (req.numWords > MAX_NUM_WORDS) {
      revert NumWordsTooBig(req.numWords, MAX_NUM_WORDS);
    }

    // Note we do not check whether the keyHash is valid to save gas.
    // The consequence for users is that they can send requests
    // for invalid keyHashes which will simply not be fulfilled.
    ++consumerConfig.nonce;
    ++consumerConfig.pendingReqCount;
    uint256 preSeed;
    (requestId, preSeed) = _computeRequestId(req.keyHash, msg.sender, subId, consumerConfig.nonce);

    bytes memory extraArgsBytes = VRFV2PlusClient._argsToBytes(_fromBytes(req.extraArgs));
    s_requestCommitments[requestId] = keccak256(
      abi.encode(requestId, _getBlockNumber(), subId, req.callbackGasLimit, req.numWords, msg.sender, extraArgsBytes)
    );
    emit RandomWordsRequested(
      req.keyHash,
      requestId,
      preSeed,
      subId,
      req.requestConfirmations,
      req.callbackGasLimit,
      req.numWords,
      extraArgsBytes,
      msg.sender
    );
    consumerConfigs[subId] = consumerConfig;

    return requestId;
  }

  function _computeRequestId(
    bytes32 keyHash,
    address sender,
    uint256 subId,
    uint64 nonce
  ) internal pure returns (uint256, uint256) {
    uint256 preSeed = uint256(keccak256(abi.encode(keyHash, sender, subId, nonce)));
    return (uint256(keccak256(abi.encode(keyHash, preSeed))), preSeed);
  }

  /**
   * @dev calls target address with exactly gasAmount gas and data as calldata
   * or reverts if at least gasAmount gas is not available.
   */
  function _callWithExactGas(uint256 gasAmount, address target, bytes memory data) private returns (bool success) {
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

  struct Output {
    ProvingKey provingKey;
    uint256 requestId;
    uint256 randomness;
  }

  function _getRandomnessFromProof(
    Proof calldata proof,
    VRFTypes.RequestCommitmentV2Plus calldata rc
  ) internal view returns (Output memory) {
    bytes32 keyHash = hashOfKey(proof.pk);
    ProvingKey memory key = s_provingKeys[keyHash];
    // Only registered proving keys are permitted.
    if (!key.exists) {
      revert NoSuchProvingKey(keyHash);
    }
    uint256 requestId = uint256(keccak256(abi.encode(keyHash, proof.seed)));
    bytes32 commitment = s_requestCommitments[requestId];
    if (commitment == 0) {
      revert NoCorrespondingRequest();
    }
    if (
      commitment !=
      keccak256(abi.encode(requestId, rc.blockNum, rc.subId, rc.callbackGasLimit, rc.numWords, rc.sender, rc.extraArgs))
    ) {
      revert IncorrectCommitment();
    }

    bytes32 blockHash = _getBlockhash(rc.blockNum);
    if (blockHash == bytes32(0)) {
      blockHash = BLOCKHASH_STORE.getBlockhash(rc.blockNum);
      if (blockHash == bytes32(0)) {
        revert BlockhashNotInStore(rc.blockNum);
      }
    }

    // The seed actually used by the VRF machinery, mixing in the blockhash
    uint256 actualSeed = uint256(keccak256(abi.encodePacked(proof.seed, blockHash)));
    uint256 randomness = VRF._randomValueFromVRFProof(proof, actualSeed); // Reverts on failure
    return Output(key, requestId, randomness);
  }

  function _getValidatedGasPrice(bool onlyPremium, uint64 gasLaneMaxGas) internal view returns (uint256 gasPrice) {
    if (tx.gasprice > gasLaneMaxGas) {
      if (onlyPremium) {
        // if only the premium amount needs to be billed, then the premium is capped by the gas lane max
        return uint256(gasLaneMaxGas);
      } else {
        // Ensure gas price does not exceed the gas lane max gas price
        revert GasPriceExceeded(tx.gasprice, gasLaneMaxGas);
      }
    }
    return tx.gasprice;
  }

  function _deliverRandomness(
    uint256 requestId,
    VRFTypes.RequestCommitmentV2Plus calldata rc,
    uint256[] memory randomWords
  ) internal returns (bool success) {
    VRFConsumerBaseV2Plus v;
    bytes memory resp = abi.encodeWithSelector(v.rawFulfillRandomWords.selector, requestId, randomWords);
    // Call with explicitly the amount of callback gas requested
    // Important to not let them exhaust the gas budget and avoid oracle payment.
    // Do not allow any non-view/non-pure coordinator functions to be called
    // during the consumers callback code via reentrancyLock.
    // Note that _callWithExactGas will revert if we do not have sufficient gas
    // to give the callee their requested amount.
    s_config.reentrancyLock = true;
    success = _callWithExactGas(rc.callbackGasLimit, rc.sender, resp);
    s_config.reentrancyLock = false;
    return success;
  }

  /*
   * @notice Fulfill a randomness request.
   * @param proof contains the proof and randomness
   * @param rc request commitment pre-image, committed to at request time
   * @param onlyPremium only charge premium
   * @return payment amount billed to the subscription
   * @dev simulated offchain to determine if sufficient balance is present to fulfill the request
   */
  function fulfillRandomWords(
    Proof calldata proof,
    VRFTypes.RequestCommitmentV2Plus calldata rc,
    bool onlyPremium
  ) external nonReentrant returns (uint96 payment) {
    uint256 startGas = gasleft();
    // fulfillRandomWords msg.data has 772 bytes and with an additional
    // buffer of 32 bytes, we get 804 bytes.
    /* Data size split:
     * fulfillRandomWords function signature - 4 bytes
     * proof - 416 bytes
     *   pk - 64 bytes
     *   gamma - 64 bytes
     *   c - 32 bytes
     *   s - 32 bytes
     *   seed - 32 bytes
     *   uWitness - 32 bytes
     *   cGammaWitness - 64 bytes
     *   sHashWitness - 64 bytes
     *   zInv - 32 bytes
     * requestCommitment - 320 bytes
     *   blockNum - 32 bytes
     *   subId - 32 bytes
     *   callbackGasLimit - 32 bytes
     *   numWords - 32 bytes
     *   sender - 32 bytes
     *   extraArgs - 128 bytes
     * onlyPremium - 32 bytes
     */
    if (msg.data.length > 804) {
      revert MsgDataTooBig(msg.data.length, 804);
    }
    Output memory output = _getRandomnessFromProof(proof, rc);
    uint256 gasPrice = _getValidatedGasPrice(onlyPremium, output.provingKey.maxGas);

    uint256[] memory randomWords;
    uint256 randomness = output.randomness;
    // stack too deep error
    {
      uint256 numWords = rc.numWords;
      randomWords = new uint256[](numWords);
      for (uint256 i = 0; i < numWords; ++i) {
        randomWords[i] = uint256(keccak256(abi.encode(randomness, i)));
      }
    }

    delete s_requestCommitments[output.requestId];
    bool success = _deliverRandomness(output.requestId, rc, randomWords);

    // Increment the req count for the subscription.
    ++s_subscriptions[rc.subId].reqCount;
    // Decrement the pending req count for the consumer.
    --s_consumers[rc.sender][rc.subId].pendingReqCount;

    bool nativePayment = uint8(rc.extraArgs[rc.extraArgs.length - 1]) == 1;

    // stack too deep error
    {
      // We want to charge users exactly for how much gas they use in their callback with
      // an additional premium. If onlyPremium is true, only premium is charged without
      // the gas cost. The gasAfterPaymentCalculation is meant to cover these additional
      // operations where we decrement the subscription balance and increment the
      // withdrawable balance.
      bool isFeedStale;
      (payment, isFeedStale) = _calculatePaymentAmount(startGas, gasPrice, nativePayment, onlyPremium);
      if (isFeedStale) {
        emit FallbackWeiPerUnitLinkUsed(output.requestId, s_fallbackWeiPerUnitLink);
      }
    }

    _chargePayment(payment, nativePayment, rc.subId);

    // Include payment in the event for tracking costs.
    emit RandomWordsFulfilled(output.requestId, randomness, rc.subId, payment, nativePayment, success, onlyPremium);

    return payment;
  }

  function _chargePayment(uint96 payment, bool nativePayment, uint256 subId) internal {
    Subscription storage subcription = s_subscriptions[subId];
    if (nativePayment) {
      uint96 prevBal = subcription.nativeBalance;
      _requireSufficientBalance(prevBal >= payment);
      subcription.nativeBalance = prevBal - payment;
      s_withdrawableNative += payment;
    } else {
      uint96 prevBal = subcription.balance;
      _requireSufficientBalance(prevBal >= payment);
      subcription.balance = prevBal - payment;
      s_withdrawableTokens += payment;
    }
  }

  function _calculatePaymentAmount(
    uint256 startGas,
    uint256 weiPerUnitGas,
    bool nativePayment,
    bool onlyPremium
  ) internal returns (uint96, bool) {
    if (nativePayment) {
      return (_calculatePaymentAmountNative(startGas, weiPerUnitGas, onlyPremium), false);
    }
    return _calculatePaymentAmountLink(startGas, weiPerUnitGas, onlyPremium);
  }

  function _calculatePaymentAmountNative(
    uint256 startGas,
    uint256 weiPerUnitGas,
    bool onlyPremium
  ) internal returns (uint96) {
    // Will return non-zero on chains that have this enabled
    uint256 l1CostWei = _getL1CostWei(msg.data);
    // calculate the payment without the premium
    uint256 baseFeeWei = weiPerUnitGas * (s_config.gasAfterPaymentCalculation + startGas - gasleft());
    // calculate flat fee in native
    uint256 flatFeeWei = 1e12 * uint256(s_config.fulfillmentFlatFeeNativePPM);
    // emit this event only if this is an L2 chain that needs to cover for L1 gas fees
    if (l1CostWei > 0) {
      emit L1GasFee(l1CostWei);
    }
    if (onlyPremium) {
      return uint96((((l1CostWei + baseFeeWei) * (s_config.nativePremiumPercentage)) / 100) + flatFeeWei);
    } else {
      return uint96((((l1CostWei + baseFeeWei) * (100 + s_config.nativePremiumPercentage)) / 100) + flatFeeWei);
    }
  }

  // Get the amount of gas used for fulfillment
  function _calculatePaymentAmountLink(
    uint256 startGas,
    uint256 weiPerUnitGas,
    bool onlyPremium
  ) internal returns (uint96, bool) {
    (int256 weiPerUnitLink, bool isFeedStale) = _getFeedData();
    if (weiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(weiPerUnitLink);
    }
    // Will return non-zero on chains that have this enabled
    uint256 l1CostWei = _getL1CostWei(msg.data);
    // (1e18 juels/link) ((wei/gas * gas) + l1wei) / (wei/link) = juels
    uint256 paymentNoFee = (1e18 *
      (weiPerUnitGas * (s_config.gasAfterPaymentCalculation + startGas - gasleft()) + l1CostWei)) /
      uint256(weiPerUnitLink);
    // calculate the flat fee in wei
    uint256 flatFeeWei = 1e12 *
      uint256(s_config.fulfillmentFlatFeeNativePPM - s_config.fulfillmentFlatFeeLinkDiscountPPM);
    uint256 flatFeeJuels = (1e18 * flatFeeWei) / uint256(weiPerUnitLink);
    // emit this event only if this is an L2 chain that needs to cover for L1 gas fees
    if (l1CostWei > 0) {
      emit L1GasFee(l1CostWei);
    }
    uint256 payment;
    if (onlyPremium) {
      payment = ((paymentNoFee * (s_config.linkPremiumPercentage)) / 100 + flatFeeJuels);
    } else {
      payment = ((paymentNoFee * (100 + s_config.linkPremiumPercentage)) / 100 + flatFeeJuels);
    }
    if (payment > 1e27) {
      revert PaymentTooLarge(); // Payment + fee cannot be more than all of the link in existence.
    }
    return (uint96(payment), isFeedStale);
  }

  function _getFeedData() private view returns (int256 weiPerUnitLink, bool isFeedStale) {
    uint32 stalenessSeconds = s_config.stalenessSeconds;
    uint256 timestamp;
    (, weiPerUnitLink, , timestamp, ) = LINK_NATIVE_FEED.latestRoundData();
    // solhint-disable-next-line not-rely-on-time
    isFeedStale = stalenessSeconds > 0 && stalenessSeconds < block.timestamp - timestamp;
    if (isFeedStale) {
      weiPerUnitLink = s_fallbackWeiPerUnitLink;
    }
    return (weiPerUnitLink, isFeedStale);
  }

  /**
   * @notice Returns the block number of the current block by using specific opcode.
   * @notice Override this function in chain specific way if needed (L2 chains).
   */
  function _getBlockNumber() internal view virtual returns (uint256) {
    return block.number;
  }

  /**
   * @notice Returns the blockhash for the given blockNumber by using specific opcode.
   * @notice If the blockNumber is more than 256 blocks in the past, returns the empty string.
   * @notice Override this function in chain specific way if needed (L2 chains).
   */
  function _getBlockhash(uint64 blockNumber) internal view virtual returns (bytes32) {
    return blockhash(blockNumber);
  }

  /**
   * @notice Returns the L1 fee for the calldata payload (always return 0 on L1 chains).
   * @notice Override this function in chain specific way for L2 chains.
   */
  function _getL1CostWei(bytes calldata /* data */) internal view virtual returns (uint256) {
    return 0;
  }

  /**
   * @inheritdoc IVRFSubscriptionV2Plus
   */
  function pendingRequestExists(uint256 subId) public view override returns (bool) {
    address[] storage consumers = s_subscriptionConfigs[subId].consumers;
    uint256 consumersLength = consumers.length;
    for (uint256 i = 0; i < consumersLength; ++i) {
      if (s_consumers[consumers[i]][subId].pendingReqCount > 0) {
        return true;
      }
    }
    return false;
  }

  /**
   * @inheritdoc IVRFSubscriptionV2Plus
   */
  function removeConsumer(uint256 subId, address consumer) external override onlySubOwner(subId) nonReentrant {
    if (pendingRequestExists(subId)) {
      revert PendingRequestExists();
    }
    if (!s_consumers[consumer][subId].active) {
      revert InvalidConsumer(subId, consumer);
    }
    // Note bounded by MAX_CONSUMERS
    address[] storage s_subscriptionConsumers = s_subscriptionConfigs[subId].consumers;
    uint256 consumersLength = s_subscriptionConsumers.length;
    for (uint256 i = 0; i < consumersLength; ++i) {
      if (s_subscriptionConsumers[i] == consumer) {
        // Storage write to preserve last element
        s_subscriptionConsumers[i] = s_subscriptionConsumers[consumersLength - 1];
        // Storage remove last element
        s_subscriptionConsumers.pop();
        break;
      }
    }
    s_consumers[consumer][subId].active = false;
    emit SubscriptionConsumerRemoved(subId, consumer);
  }

  /**
   * @inheritdoc IVRFSubscriptionV2Plus
   */
  function cancelSubscription(uint256 subId, address to) external override onlySubOwner(subId) nonReentrant {
    if (pendingRequestExists(subId)) {
      revert PendingRequestExists();
    }
    _cancelSubscriptionHelper(subId, to);
  }

  /***************************************************************************
   * Section: Migration
   ***************************************************************************/

  address[] internal s_migrationTargets;

  /// @dev Emitted when new coordinator is registered as migratable target
  event CoordinatorRegistered(address coordinatorAddress);

  /// @dev Emitted when new coordinator is deregistered
  event CoordinatorDeregistered(address coordinatorAddress);

  /// @notice emitted when migration to new coordinator completes successfully
  /// @param newCoordinator coordinator address after migration
  /// @param subId subscription ID
  event MigrationCompleted(address newCoordinator, uint256 subId);

  /// @notice emitted when migrate() is called and given coordinator is not registered as migratable target
  error CoordinatorNotRegistered(address coordinatorAddress);

  /// @notice emitted when migrate() is called and given coordinator is registered as migratable target
  error CoordinatorAlreadyRegistered(address coordinatorAddress);

  /// @dev encapsulates data to be migrated from current coordinator
  // solhint-disable-next-line gas-struct-packing
  struct V1MigrationData {
    uint8 fromVersion;
    uint256 subId;
    address subOwner;
    address[] consumers;
    uint96 linkBalance;
    uint96 nativeBalance;
  }

  function _isTargetRegistered(address target) internal view returns (bool) {
    uint256 migrationTargetsLength = s_migrationTargets.length;
    for (uint256 i = 0; i < migrationTargetsLength; ++i) {
      if (s_migrationTargets[i] == target) {
        return true;
      }
    }
    return false;
  }

  function registerMigratableCoordinator(address target) external onlyOwner {
    if (_isTargetRegistered(target)) {
      revert CoordinatorAlreadyRegistered(target);
    }
    s_migrationTargets.push(target);
    emit CoordinatorRegistered(target);
  }

  function deregisterMigratableCoordinator(address target) external onlyOwner {
    uint256 nTargets = s_migrationTargets.length;
    for (uint256 i = 0; i < nTargets; ++i) {
      if (s_migrationTargets[i] == target) {
        s_migrationTargets[i] = s_migrationTargets[nTargets - 1];
        s_migrationTargets.pop();
        emit CoordinatorDeregistered(target);
        return;
      }
    }
    revert CoordinatorNotRegistered(target);
  }

  function migrate(uint256 subId, address newCoordinator) external nonReentrant {
    if (!_isTargetRegistered(newCoordinator)) {
      revert CoordinatorNotRegistered(newCoordinator);
    }
    (uint96 balance, uint96 nativeBalance, , address subOwner, address[] memory consumers) = getSubscription(subId);
    if (subOwner != msg.sender) {
      revert MustBeSubOwner(subOwner);
    }
    if (pendingRequestExists(subId)) {
      revert PendingRequestExists();
    }

    V1MigrationData memory migrationData = V1MigrationData({
      fromVersion: 1,
      subId: subId,
      subOwner: subOwner,
      consumers: consumers,
      linkBalance: balance,
      nativeBalance: nativeBalance
    });
    bytes memory encodedData = abi.encode(migrationData);
    _deleteSubscription(subId);
    IVRFCoordinatorV2PlusMigration(newCoordinator).onMigration{value: nativeBalance}(encodedData);

    // Only transfer LINK if the token is active and there is a balance.
    if (address(LINK) != address(0) && balance != 0) {
      _requireSufficientBalance(LINK.transfer(address(newCoordinator), balance));
    }

    // despite the fact that we follow best practices this is still probably safest
    // to prevent any re-entrancy possibilities.
    s_config.reentrancyLock = true;
    uint256 consumersLength = consumers.length;
    for (uint256 i = 0; i < consumersLength; ++i) {
      IVRFMigratableConsumerV2Plus(consumers[i]).setCoordinator(newCoordinator);
    }
    s_config.reentrancyLock = false;

    emit MigrationCompleted(newCoordinator, subId);
  }
}
