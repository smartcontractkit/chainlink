// SPDX-License-Identifier: MIT
// A mock for testing code that relies on VRFCoordinatorV2_5.
pragma solidity ^0.8.19;

// solhint-disable-next-line no-unused-import
import {IVRFCoordinatorV2Plus, IVRFSubscriptionV2Plus} from "../dev/interfaces/IVRFCoordinatorV2Plus.sol";
import {VRFV2PlusClient} from "../dev/libraries/VRFV2PlusClient.sol";
import {SubscriptionAPI} from "../dev/SubscriptionAPI.sol";
import {VRFConsumerBaseV2Plus} from "../dev/VRFConsumerBaseV2Plus.sol";

contract VRFCoordinatorV2_5Mock is SubscriptionAPI, IVRFCoordinatorV2Plus {
  uint96 public immutable i_base_fee;
  uint96 public immutable i_gas_price;
  int256 public immutable i_wei_per_unit_link;

  error InvalidRequest();
  error InvalidRandomWords();
  error InvalidExtraArgsTag();
  error NotImplemented();

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
  event ConfigSet();

  uint64 internal s_currentSubId;
  uint256 internal s_nextRequestId = 1;
  uint256 internal s_nextPreSeed = 100;

  struct Request {
    uint256 subId;
    uint32 callbackGasLimit;
    uint32 numWords;
    bytes extraArgs;
  }
  mapping(uint256 => Request) internal s_requests; /* requestId */ /* request */

  constructor(uint96 _baseFee, uint96 _gasPrice, int256 _weiPerUnitLink) SubscriptionAPI() {
    i_base_fee = _baseFee;
    i_gas_price = _gasPrice;
    i_wei_per_unit_link = _weiPerUnitLink;
    setConfig();
  }

  /**
   * @notice Sets the configuration of the vrfv2 mock coordinator
   */
  function setConfig() public onlyOwner {
    s_config = Config({
      minimumRequestConfirmations: 0,
      maxGasLimit: 0,
      stalenessSeconds: 0,
      gasAfterPaymentCalculation: 0,
      reentrancyLock: false,
      fulfillmentFlatFeeNativePPM: 0,
      fulfillmentFlatFeeLinkDiscountPPM: 0,
      nativePremiumPercentage: 0,
      linkPremiumPercentage: 0
    });
    emit ConfigSet();
  }

  function consumerIsAdded(uint256 _subId, address _consumer) public view returns (bool) {
    return s_consumers[_consumer][_subId].active;
  }

  modifier onlyValidConsumer(uint256 _subId, address _consumer) {
    if (!consumerIsAdded(_subId, _consumer)) {
      revert InvalidConsumer(_subId, _consumer);
    }
    _;
  }

  /**
   * @notice fulfillRandomWords fulfills the given request, sending the random words to the supplied
   * @notice consumer.
   *
   * @dev This mock uses a simplified formula for calculating payment amount and gas usage, and does
   * @dev not account for all edge cases handled in the real VRF coordinator. When making requests
   * @dev against the real coordinator a small amount of additional LINK is required.
   *
   * @param _requestId the request to fulfill
   * @param _consumer the VRF randomness consumer to send the result to
   */
  function fulfillRandomWords(uint256 _requestId, address _consumer) external nonReentrant {
    fulfillRandomWordsWithOverride(_requestId, _consumer, new uint256[](0));
  }

  /**
   * @notice fulfillRandomWordsWithOverride allows the user to pass in their own random words.
   *
   * @param _requestId the request to fulfill
   * @param _consumer the VRF randomness consumer to send the result to
   * @param _words user-provided random words
   */
  function fulfillRandomWordsWithOverride(uint256 _requestId, address _consumer, uint256[] memory _words) public {
    uint256 startGas = gasleft();
    if (s_requests[_requestId].subId == 0) {
      revert InvalidRequest();
    }
    Request memory req = s_requests[_requestId];

    if (_words.length == 0) {
      _words = new uint256[](req.numWords);
      for (uint256 i = 0; i < req.numWords; i++) {
        _words[i] = uint256(keccak256(abi.encode(_requestId, i)));
      }
    } else if (_words.length != req.numWords) {
      revert InvalidRandomWords();
    }

    VRFConsumerBaseV2Plus v;
    bytes memory callReq = abi.encodeWithSelector(v.rawFulfillRandomWords.selector, _requestId, _words);
    s_config.reentrancyLock = true;
    // solhint-disable-next-line avoid-low-level-calls, no-unused-vars
    (bool success, ) = _consumer.call{gas: req.callbackGasLimit}(callReq);
    s_config.reentrancyLock = false;

    bool nativePayment = uint8(req.extraArgs[req.extraArgs.length - 1]) == 1;

    uint256 rawPayment = i_base_fee + ((startGas - gasleft()) * i_gas_price);
    if (!nativePayment) {
      rawPayment = (1e18 * rawPayment) / uint256(i_wei_per_unit_link);
    }
    uint96 payment = uint96(rawPayment);

    _chargePayment(payment, nativePayment, req.subId);

    delete (s_requests[_requestId]);
    emit RandomWordsFulfilled(_requestId, _requestId, req.subId, payment, nativePayment, success, false);
  }

  function _chargePayment(uint96 payment, bool nativePayment, uint256 subId) internal {
    Subscription storage subcription = s_subscriptions[subId];
    if (nativePayment) {
      uint96 prevBal = subcription.nativeBalance;
      if (prevBal < payment) {
        revert InsufficientBalance();
      }
      subcription.nativeBalance = prevBal - payment;
      s_withdrawableNative += payment;
    } else {
      uint96 prevBal = subcription.balance;
      if (prevBal < payment) {
        revert InsufficientBalance();
      }
      subcription.balance = prevBal - payment;
      s_withdrawableTokens += payment;
    }
  }

  /**
   * @notice fundSubscription allows funding a subscription with an arbitrary amount for testing.
   *
   * @param _subId the subscription to fund
   * @param _amount the amount to fund
   */
  function fundSubscription(uint256 _subId, uint256 _amount) public {
    if (s_subscriptionConfigs[_subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    uint256 oldBalance = s_subscriptions[_subId].balance;
    s_subscriptions[_subId].balance += uint96(_amount);
    s_totalBalance += uint96(_amount);
    emit SubscriptionFunded(_subId, oldBalance, oldBalance + _amount);
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

  function requestRandomWords(
    VRFV2PlusClient.RandomWordsRequest calldata _req
  ) external override nonReentrant onlyValidConsumer(_req.subId, msg.sender) returns (uint256) {
    uint256 subId = _req.subId;
    if (s_subscriptionConfigs[subId].owner == address(0)) {
      revert InvalidSubscription();
    }

    uint256 requestId = s_nextRequestId++;
    uint256 preSeed = s_nextPreSeed++;

    bytes memory extraArgsBytes = VRFV2PlusClient._argsToBytes(_fromBytes(_req.extraArgs));
    s_requests[requestId] = Request({
      subId: _req.subId,
      callbackGasLimit: _req.callbackGasLimit,
      numWords: _req.numWords,
      extraArgs: _req.extraArgs
    });

    emit RandomWordsRequested(
      _req.keyHash,
      requestId,
      preSeed,
      _req.subId,
      _req.requestConfirmations,
      _req.callbackGasLimit,
      _req.numWords,
      extraArgsBytes,
      msg.sender
    );
    return requestId;
  }

  /**
   * @inheritdoc IVRFSubscriptionV2Plus
   */
  function removeConsumer(
    uint256 _subId,
    address _consumer
  ) external override onlySubOwner(_subId) onlyValidConsumer(_subId, _consumer) nonReentrant {
    if (!s_consumers[_consumer][_subId].active) {
      revert InvalidConsumer(_subId, _consumer);
    }
    address[] memory consumers = s_subscriptionConfigs[_subId].consumers;
    uint256 lastConsumerIndex = consumers.length - 1;
    for (uint256 i = 0; i < consumers.length; ++i) {
      if (consumers[i] == _consumer) {
        address last = consumers[lastConsumerIndex];
        s_subscriptionConfigs[_subId].consumers[i] = last;
        s_subscriptionConfigs[_subId].consumers.pop();
        break;
      }
    }
    s_consumers[_consumer][_subId].active = false;
    emit SubscriptionConsumerRemoved(_subId, _consumer);
  }

  /**
   * @inheritdoc IVRFSubscriptionV2Plus
   */
  function cancelSubscription(uint256 _subId, address _to) external override onlySubOwner(_subId) nonReentrant {
    (uint96 balance, uint96 nativeBalance) = _deleteSubscription(_subId);

    (bool success, ) = _to.call{value: uint256(nativeBalance)}("");
    if (!success) {
      revert FailedToSendNative();
    }
    emit SubscriptionCanceled(_subId, _to, balance, nativeBalance);
  }

  function pendingRequestExists(uint256 /*subId*/) public pure override returns (bool) {
    revert NotImplemented();
  }
}
