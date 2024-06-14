// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "../../../shared/interfaces/LinkTokenInterface.sol";
import "../../../interfaces/BlockhashStoreInterface.sol";
import "../../../interfaces/TypeAndVersionInterface.sol";
import "../../interfaces/IVRFCoordinatorV2Plus.sol";
import "../../../vrf/VRF.sol";
import "../VRFConsumerBaseV2Plus.sol";
import "../../../ChainSpecificUtil.sol";
import "../SubscriptionAPI.sol";
import "../libraries/VRFV2PlusClient.sol";
import "../../interfaces/IVRFCoordinatorV2PlusMigration.sol";
import "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";

contract VRFCoordinatorV2PlusUpgradedVersion is
  VRF,
  SubscriptionAPI,
  IVRFCoordinatorV2PlusMigration,
  IVRFCoordinatorV2Plus
{
  using EnumerableSet for EnumerableSet.UintSet;
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
  error InvalidExtraArgsTag();
  /// @notice emitted when version in the request doesn't match expected version
  error InvalidVersion(uint8 requestVersion, uint8 expectedVersion);
  /// @notice emitted when transferred balance (msg.value) does not match the metadata in V1MigrationData
  error InvalidNativeBalance(uint256 transferredValue, uint96 expectedValue);
  error SubscriptionIDCollisionFound();

  struct RequestCommitment {
    uint64 blockNum;
    uint256 subId;
    uint32 callbackGasLimit;
    uint32 numWords;
    address sender;
    bytes extraArgs;
  }

  mapping(bytes32 => address) /* keyHash */ /* oracle */ public s_provingKeys;
  bytes32[] public s_provingKeyHashes;
  mapping(uint256 => bytes32) /* requestID */ /* commitment */ public s_requestCommitments;

  event ProvingKeyRegistered(bytes32 keyHash, address indexed oracle);
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
    uint256 indexed subID,
    uint96 payment,
    bool success
  );

  int256 public s_fallbackWeiPerUnitLink;

  FeeConfig public s_feeConfig;

  struct FeeConfig {
    // Flat fee charged per fulfillment in millionths of link
    // So fee range is [0, 2^32/10^6].
    uint32 fulfillmentFlatFeeLinkPPM;
    // Flat fee charged per fulfillment in millionths of native.
    // So fee range is [0, 2^32/10^6].
    uint32 fulfillmentFlatFeeNativePPM;
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
   * @param stalenessSeconds if the native/link feed is more stale then this, use the fallback price
   * @param gasAfterPaymentCalculation gas used in doing accounting after completing the gas measurement
   * @param fallbackWeiPerUnitLink fallback native/link price in the case of a stale feed
   * @param feeConfig fee configuration
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
  function getRequestConfig() external view returns (uint16, uint32, bytes32[] memory) {
    return (s_config.minimumRequestConfirmations, s_config.maxGasLimit, s_provingKeyHashes);
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
  ) external override nonReentrant returns (uint256) {
    // Input validation using the subscription storage.
    if (s_subscriptionConfigs[req.subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    // Its important to ensure that the consumer is in fact who they say they
    // are, otherwise they could use someone else's subscription balance.
    // A nonce of 0 indicates consumer is not allocated to the sub.
    uint64 currentNonce = s_consumers[msg.sender][req.subId];
    if (currentNonce == 0) {
      revert InvalidConsumer(req.subId, msg.sender);
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
    uint64 nonce = currentNonce + 1;
    (uint256 requestId, uint256 preSeed) = computeRequestId(req.keyHash, msg.sender, req.subId, nonce);

    VRFV2PlusClient.ExtraArgsV1 memory extraArgs = _fromBytes(req.extraArgs);
    bytes memory extraArgsBytes = VRFV2PlusClient._argsToBytes(extraArgs);
    s_requestCommitments[requestId] = keccak256(
      abi.encode(
        requestId,
        ChainSpecificUtil.getBlockNumber(),
        req.subId,
        req.callbackGasLimit,
        req.numWords,
        msg.sender,
        extraArgsBytes
      )
    );
    emit RandomWordsRequested(
      req.keyHash,
      requestId,
      preSeed,
      req.subId,
      req.requestConfirmations,
      req.callbackGasLimit,
      req.numWords,
      extraArgsBytes,
      msg.sender
    );
    s_consumers[msg.sender][req.subId] = nonce;

    return requestId;
  }

  function computeRequestId(
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

  struct Output {
    bytes32 keyHash;
    uint256 requestId;
    uint256 randomness;
  }

  function getRandomnessFromProof(
    Proof memory proof,
    RequestCommitment memory rc
  ) internal view returns (Output memory) {
    bytes32 keyHash = hashOfKey(proof.pk);
    // Only registered proving keys are permitted.
    address oracle = s_provingKeys[keyHash];
    if (oracle == address(0)) {
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

    bytes32 blockHash = ChainSpecificUtil.getBlockhash(rc.blockNum);
    if (blockHash == bytes32(0)) {
      blockHash = BLOCKHASH_STORE.getBlockhash(rc.blockNum);
      if (blockHash == bytes32(0)) {
        revert BlockhashNotInStore(rc.blockNum);
      }
    }

    // The seed actually used by the VRF machinery, mixing in the blockhash
    uint256 actualSeed = uint256(keccak256(abi.encodePacked(proof.seed, blockHash)));
    uint256 randomness = VRF.randomValueFromVRFProof(proof, actualSeed); // Reverts on failure
    return Output(keyHash, requestId, randomness);
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
    Output memory output = getRandomnessFromProof(proof, rc);

    uint256[] memory randomWords = new uint256[](rc.numWords);
    for (uint256 i = 0; i < rc.numWords; i++) {
      randomWords[i] = uint256(keccak256(abi.encode(output.randomness, i)));
    }

    delete s_requestCommitments[output.requestId];
    VRFConsumerBaseV2Plus v;
    bytes memory resp = abi.encodeWithSelector(v.rawFulfillRandomWords.selector, output.requestId, randomWords);
    // Call with explicitly the amount of callback gas requested
    // Important to not let them exhaust the gas budget and avoid oracle payment.
    // Do not allow any non-view/non-pure coordinator functions to be called
    // during the consumers callback code via reentrancyLock.
    // Note that callWithExactGas will revert if we do not have sufficient gas
    // to give the callee their requested amount.
    s_config.reentrancyLock = true;
    bool success = callWithExactGas(rc.callbackGasLimit, rc.sender, resp);
    s_config.reentrancyLock = false;

    // Increment the req count for the subscription.
    uint64 reqCount = s_subscriptions[rc.subId].reqCount;
    s_subscriptions[rc.subId].reqCount = reqCount + 1;

    // stack too deep error
    {
      bool nativePayment = uint8(rc.extraArgs[rc.extraArgs.length - 1]) == 1;
      // We want to charge users exactly for how much gas they use in their callback.
      // The gasAfterPaymentCalculation is meant to cover these additional operations where we
      // decrement the subscription balance and increment the oracles withdrawable balance.
      uint96 payment = calculatePaymentAmount(
        startGas,
        s_config.gasAfterPaymentCalculation,
        tx.gasprice,
        nativePayment
      );
      if (nativePayment) {
        if (s_subscriptions[rc.subId].nativeBalance < payment) {
          revert InsufficientBalance();
        }
        s_subscriptions[rc.subId].nativeBalance -= payment;
        s_withdrawableNative[s_provingKeys[output.keyHash]] += payment;
      } else {
        if (s_subscriptions[rc.subId].balance < payment) {
          revert InsufficientBalance();
        }
        s_subscriptions[rc.subId].balance -= payment;
        s_withdrawableTokens[s_provingKeys[output.keyHash]] += payment;
      }

      // Include payment in the event for tracking costs.
      // event RandomWordsFulfilled(uint256 indexed requestId, uint256 outputSeed, uint96 payment, bytes extraArgs, bool success);
      emit RandomWordsFulfilled(output.requestId, output.randomness, rc.subId, payment, success);

      return payment;
    }
  }

  function calculatePaymentAmount(
    uint256 startGas,
    uint256 gasAfterPaymentCalculation,
    uint256 weiPerUnitGas,
    bool nativePayment
  ) internal view returns (uint96) {
    if (nativePayment) {
      return
        calculatePaymentAmountNative(
          startGas,
          gasAfterPaymentCalculation,
          s_feeConfig.fulfillmentFlatFeeNativePPM,
          weiPerUnitGas
        );
    }
    return
      calculatePaymentAmountLink(
        startGas,
        gasAfterPaymentCalculation,
        s_feeConfig.fulfillmentFlatFeeLinkPPM,
        weiPerUnitGas
      );
  }

  function calculatePaymentAmountNative(
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
    uint256 flatFeeWei = 1e12 * uint256(fulfillmentFlatFeePPM);
    // return the final fee with the flat fee and l1 cost (if applicable) added
    return uint96(baseFeeWei + flatFeeWei + l1CostWei);
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
    (, weiPerUnitLink, , timestamp, ) = LINK_NATIVE_FEED.latestRoundData();
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
  function pendingRequestExists(uint256 subId) public view override returns (bool) {
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

  /**
   * @inheritdoc IVRFSubscriptionV2Plus
   */
  function removeConsumer(uint256 subId, address consumer) external override onlySubOwner(subId) nonReentrant {
    if (pendingRequestExists(subId)) {
      revert PendingRequestExists();
    }
    if (s_consumers[consumer][subId] == 0) {
      revert InvalidConsumer(subId, consumer);
    }
    // Note bounded by MAX_CONSUMERS
    address[] memory consumers = s_subscriptionConfigs[subId].consumers;
    uint256 lastConsumerIndex = consumers.length - 1;
    for (uint256 i = 0; i < consumers.length; i++) {
      if (consumers[i] == consumer) {
        address last = consumers[lastConsumerIndex];
        // Storage write to preserve last element
        s_subscriptionConfigs[subId].consumers[i] = last;
        // Storage remove last element
        s_subscriptionConfigs[subId].consumers.pop();
        break;
      }
    }
    delete s_consumers[consumer][subId];
    emit SubscriptionConsumerRemoved(subId, consumer);
  }

  /**
   * @inheritdoc IVRFSubscriptionV2Plus
   */
  function cancelSubscription(uint256 subId, address to) external override onlySubOwner(subId) nonReentrant {
    if (pendingRequestExists(subId)) {
      revert PendingRequestExists();
    }
    cancelSubscriptionHelper(subId, to);
  }

  /***************************************************************************
   * Section: Migration
   ***************************************************************************/

  address[] internal s_migrationTargets;

  /// @dev Emitted when new coordinator is registered as migratable target
  event CoordinatorRegistered(address coordinatorAddress);

  /// @notice emitted when migration to new coordinator completes successfully
  /// @param newCoordinator coordinator address after migration
  /// @param subId subscription ID
  event MigrationCompleted(address newCoordinator, uint256 subId);

  /// @notice emitted when migrate() is called and given coordinator is not registered as migratable target
  error CoordinatorNotRegistered(address coordinatorAddress);

  /// @notice emitted when migrate() is called and given coordinator is registered as migratable target
  error CoordinatorAlreadyRegistered(address coordinatorAddress);

  /// @dev encapsulates data to be migrated from current coordinator
  struct V1MigrationData {
    uint8 fromVersion;
    uint256 subId;
    address subOwner;
    address[] consumers;
    uint96 linkBalance;
    uint96 nativeBalance;
  }

  function isTargetRegistered(address target) internal view returns (bool) {
    for (uint256 i = 0; i < s_migrationTargets.length; i++) {
      if (s_migrationTargets[i] == target) {
        return true;
      }
    }
    return false;
  }

  function registerMigratableCoordinator(address target) external onlyOwner {
    if (isTargetRegistered(target)) {
      revert CoordinatorAlreadyRegistered(target);
    }
    s_migrationTargets.push(target);
    emit CoordinatorRegistered(target);
  }

  function migrate(uint256 subId, address newCoordinator) external nonReentrant {
    if (!isTargetRegistered(newCoordinator)) {
      revert CoordinatorNotRegistered(newCoordinator);
    }
    (uint96 balance, uint96 nativeBalance, , address owner, address[] memory consumers) = getSubscription(subId);
    require(owner == msg.sender, "Not subscription owner");
    require(!pendingRequestExists(subId), "Pending request exists");

    V1MigrationData memory migrationData = V1MigrationData({
      fromVersion: migrationVersion(),
      subId: subId,
      subOwner: owner,
      consumers: consumers,
      linkBalance: balance,
      nativeBalance: nativeBalance
    });
    bytes memory encodedData = abi.encode(migrationData);
    deleteSubscription(subId);
    IVRFCoordinatorV2PlusMigration(newCoordinator).onMigration{value: nativeBalance}(encodedData);

    // Only transfer LINK if the token is active and there is a balance.
    if (address(LINK) != address(0) && balance != 0) {
      require(LINK.transfer(address(newCoordinator), balance), "insufficient funds");
    }

    // despite the fact that we follow best practices this is still probably safest
    // to prevent any re-entrancy possibilities.
    s_config.reentrancyLock = true;
    for (uint256 i = 0; i < consumers.length; i++) {
      IVRFMigratableConsumerV2Plus(consumers[i]).setCoordinator(newCoordinator);
    }
    s_config.reentrancyLock = false;

    emit MigrationCompleted(newCoordinator, subId);
  }

  function migrationVersion() public pure returns (uint8 version) {
    return 1;
  }

  /**
   * @inheritdoc IVRFCoordinatorV2PlusMigration
   */
  function onMigration(bytes calldata encodedData) external payable override {
    V1MigrationData memory migrationData = abi.decode(encodedData, (V1MigrationData));

    if (migrationData.fromVersion != 1) {
      revert InvalidVersion(migrationData.fromVersion, 1);
    }

    if (msg.value != uint256(migrationData.nativeBalance)) {
      revert InvalidNativeBalance(msg.value, migrationData.nativeBalance);
    }

    // it should be impossible to have a subscription id collision, for two reasons:
    // 1. the subscription ID is calculated using inputs that cannot be replicated under different
    // conditions.
    // 2. once a subscription is migrated it is deleted from the previous coordinator, so it cannot
    // be migrated again.
    // however, we should have this check here in case the `migrate` function on
    // future coordinators "forgets" to delete subscription data allowing re-migration of the same
    // subscription.
    if (s_subscriptionConfigs[migrationData.subId].owner != address(0)) {
      revert SubscriptionIDCollisionFound();
    }

    for (uint i = 0; i < migrationData.consumers.length; i++) {
      s_consumers[migrationData.consumers[i]][migrationData.subId] = 1;
    }

    s_subscriptions[migrationData.subId] = Subscription({
      nativeBalance: migrationData.nativeBalance,
      balance: migrationData.linkBalance,
      reqCount: 0
    });
    s_subscriptionConfigs[migrationData.subId] = SubscriptionConfig({
      owner: migrationData.subOwner,
      consumers: migrationData.consumers,
      requestedOwner: address(0)
    });

    s_totalBalance += uint96(migrationData.linkBalance);
    s_totalNativeBalance += uint96(migrationData.nativeBalance);

    s_subIds.add(migrationData.subId);
  }
}
