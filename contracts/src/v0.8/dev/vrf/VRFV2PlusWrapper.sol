// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../../shared/access/ConfirmedOwner.sol";
import "../../interfaces/TypeAndVersionInterface.sol";
import "./VRFConsumerBaseV2Plus.sol";
import "../../shared/interfaces/LinkTokenInterface.sol";
import "../../interfaces/AggregatorV3Interface.sol";
import "../interfaces/IVRFCoordinatorV2Plus.sol";
import "../interfaces/VRFV2PlusWrapperInterface.sol";
import "./VRFV2PlusWrapperConsumerBase.sol";
import "../../ChainSpecificUtil.sol";

/**
 * @notice A wrapper for VRFCoordinatorV2 that provides an interface better suited to one-off
 * @notice requests for randomness.
 */
contract VRFV2PlusWrapper is ConfirmedOwner, TypeAndVersionInterface, VRFConsumerBaseV2Plus, VRFV2PlusWrapperInterface {
  event WrapperFulfillmentFailed(uint256 indexed requestId, address indexed consumer);

  error LinkAlreadySet();
  error FailedToTransferLink();

  LinkTokenInterface public s_link;
  AggregatorV3Interface public s_linkNativeFeed;
  ExtendedVRFCoordinatorV2PlusInterface public immutable COORDINATOR;
  uint256 public immutable SUBSCRIPTION_ID;
  /// @dev this is the size of a VRF v2 fulfillment's calldata abi-encoded in bytes.
  /// @dev proofSize = 13 words = 13 * 256 = 3328 bits
  /// @dev commitmentSize = 5 words = 5 * 256 = 1280 bits
  /// @dev dataSize = proofSize + commitmentSize = 4608 bits
  /// @dev selector = 32 bits
  /// @dev total data size = 4608 bits + 32 bits = 4640 bits = 580 bytes
  uint32 public s_fulfillmentTxSizeBytes = 580;

  // 5k is plenty for an EXTCODESIZE call (2600) + warm CALL (100)
  // and some arithmetic operations.
  uint256 private constant GAS_FOR_CALL_EXACT_CHECK = 5_000;

  // lastRequestId is the request ID of the most recent VRF V2 request made by this wrapper. This
  // should only be relied on within the same transaction the request was made.
  uint256 public override lastRequestId;

  // Configuration fetched from VRFCoordinatorV2

  // s_configured tracks whether this contract has been configured. If not configured, randomness
  // requests cannot be made.
  bool public s_configured;

  // s_disabled disables the contract when true. When disabled, new VRF requests cannot be made
  // but existing ones can still be fulfilled.
  bool public s_disabled;

  // s_fallbackWeiPerUnitLink is the backup LINK exchange rate used when the LINK/NATIVE feed is
  // stale.
  int256 private s_fallbackWeiPerUnitLink;

  // s_stalenessSeconds is the number of seconds before we consider the feed price to be stale and
  // fallback to fallbackWeiPerUnitLink.
  uint32 private s_stalenessSeconds;

  // s_fulfillmentFlatFeeLinkPPM is the flat fee in millionths of LINK that VRFCoordinatorV2
  // charges.
  uint32 private s_fulfillmentFlatFeeLinkPPM;

  // s_fulfillmentFlatFeeLinkPPM is the flat fee in millionths of LINK that VRFCoordinatorV2
  // charges.
  uint32 private s_fulfillmentFlatFeeNativePPM;

  // Other configuration

  // s_wrapperGasOverhead reflects the gas overhead of the wrapper's fulfillRandomWords
  // function. The cost for this gas is passed to the user.
  uint32 private s_wrapperGasOverhead;

  // s_coordinatorGasOverhead reflects the gas overhead of the coordinator's fulfillRandomWords
  // function. The cost for this gas is billed to the subscription, and must therefor be included
  // in the pricing for wrapped requests. This includes the gas costs of proof verification and
  // payment calculation in the coordinator.
  uint32 private s_coordinatorGasOverhead;

  // s_wrapperPremiumPercentage is the premium ratio in percentage. For example, a value of 0
  // indicates no premium. A value of 15 indicates a 15 percent premium.
  uint8 private s_wrapperPremiumPercentage;

  // s_keyHash is the key hash to use when requesting randomness. Fees are paid based on current gas
  // fees, so this should be set to the highest gas lane on the network.
  bytes32 s_keyHash;

  // s_maxNumWords is the max number of words that can be requested in a single wrapped VRF request.
  uint8 s_maxNumWords;

  struct Callback {
    address callbackAddress;
    uint32 callbackGasLimit;
    uint256 requestGasPrice;
  }
  mapping(uint256 => Callback) /* requestID */ /* callback */ public s_callbacks;

  constructor(address _link, address _linkNativeFeed, address _coordinator) VRFConsumerBaseV2Plus(_coordinator) {
    if (_link != address(0)) {
      s_link = LinkTokenInterface(_link);
    }
    if (_linkNativeFeed != address(0)) {
      s_linkNativeFeed = AggregatorV3Interface(_linkNativeFeed);
    }
    COORDINATOR = ExtendedVRFCoordinatorV2PlusInterface(_coordinator);

    // Create this wrapper's subscription and add itself as a consumer.
    uint256 subId = ExtendedVRFCoordinatorV2PlusInterface(_coordinator).createSubscription();
    SUBSCRIPTION_ID = subId;
    ExtendedVRFCoordinatorV2PlusInterface(_coordinator).addConsumer(subId, address(this));
  }

  /**
   * @notice set the link token to be used by this wrapper
   * @param link address of the link token
   */
  function setLINK(address link) external onlyOwner {
    // Disallow re-setting link token because the logic wouldn't really make sense
    if (address(s_link) != address(0)) {
      revert LinkAlreadySet();
    }
    s_link = LinkTokenInterface(link);
  }

  /**
   * @notice set the link native feed to be used by this wrapper
   * @param linkNativeFeed address of the link native feed
   */
  function setLinkNativeFeed(address linkNativeFeed) external onlyOwner {
    s_linkNativeFeed = AggregatorV3Interface(linkNativeFeed);
  }

  /**
   * @notice setFulfillmentTxSize sets the size of the fulfillment transaction in bytes.
   * @param size is the size of the fulfillment transaction in bytes.
   */
  function setFulfillmentTxSize(uint32 size) external onlyOwner {
    s_fulfillmentTxSizeBytes = size;
  }

  /**
   * @notice setConfig configures VRFV2Wrapper.
   *
   * @dev Sets wrapper-specific configuration based on the given parameters, and fetches any needed
   * @dev VRFCoordinatorV2 configuration from the coordinator.
   *
   * @param _wrapperGasOverhead reflects the gas overhead of the wrapper's fulfillRandomWords
   *        function.
   *
   * @param _coordinatorGasOverhead reflects the gas overhead of the coordinator's
   *        fulfillRandomWords function.
   *
   * @param _wrapperPremiumPercentage is the premium ratio in percentage for wrapper requests.
   *
   * @param _keyHash to use for requesting randomness.
   */
  function setConfig(
    uint32 _wrapperGasOverhead,
    uint32 _coordinatorGasOverhead,
    uint8 _wrapperPremiumPercentage,
    bytes32 _keyHash,
    uint8 _maxNumWords
  ) external onlyOwner {
    s_wrapperGasOverhead = _wrapperGasOverhead;
    s_coordinatorGasOverhead = _coordinatorGasOverhead;
    s_wrapperPremiumPercentage = _wrapperPremiumPercentage;
    s_keyHash = _keyHash;
    s_maxNumWords = _maxNumWords;
    s_configured = true;

    // Get other configuration from coordinator
    (, , s_stalenessSeconds, ) = COORDINATOR.s_config();
    s_fallbackWeiPerUnitLink = COORDINATOR.s_fallbackWeiPerUnitLink();
    (s_fulfillmentFlatFeeLinkPPM, s_fulfillmentFlatFeeNativePPM) = COORDINATOR.s_feeConfig();
  }

  /**
   * @notice getConfig returns the current VRFV2Wrapper configuration.
   *
   * @return fallbackWeiPerUnitLink is the backup LINK exchange rate used when the LINK/NATIVE feed
   *         is stale.
   *
   * @return stalenessSeconds is the number of seconds before we consider the feed price to be stale
   *         and fallback to fallbackWeiPerUnitLink.
   *
   * @return fulfillmentFlatFeeLinkPPM is the flat fee in millionths of LINK that VRFCoordinatorV2Plus
   *         charges.
   *
   * @return wrapperGasOverhead reflects the gas overhead of the wrapper's fulfillRandomWords
   *         function. The cost for this gas is passed to the user.
   *
   * @return coordinatorGasOverhead reflects the gas overhead of the coordinator's
   *         fulfillRandomWords function.
   *
   * @return wrapperPremiumPercentage is the premium ratio in percentage. For example, a value of 0
   *         indicates no premium. A value of 15 indicates a 15 percent premium.
   *
   * @return keyHash is the key hash to use when requesting randomness. Fees are paid based on
   *         current gas fees, so this should be set to the highest gas lane on the network.
   *
   * @return maxNumWords is the max number of words that can be requested in a single wrapped VRF
   *         request.
   */
  function getConfig()
    external
    view
    returns (
      int256 fallbackWeiPerUnitLink,
      uint32 stalenessSeconds,
      uint32 fulfillmentFlatFeeLinkPPM,
      uint32 wrapperGasOverhead,
      uint32 coordinatorGasOverhead,
      uint8 wrapperPremiumPercentage,
      bytes32 keyHash,
      uint8 maxNumWords
    )
  {
    return (
      s_fallbackWeiPerUnitLink,
      s_stalenessSeconds,
      s_fulfillmentFlatFeeLinkPPM,
      s_wrapperGasOverhead,
      s_coordinatorGasOverhead,
      s_wrapperPremiumPercentage,
      s_keyHash,
      s_maxNumWords
    );
  }

  /**
   * @notice Calculates the price of a VRF request with the given callbackGasLimit at the current
   * @notice block.
   *
   * @dev This function relies on the transaction gas price which is not automatically set during
   * @dev simulation. To estimate the price at a specific gas price, use the estimatePrice function.
   *
   * @param _callbackGasLimit is the gas limit used to estimate the price.
   */
  function calculateRequestPrice(
    uint32 _callbackGasLimit
  ) external view override onlyConfiguredNotDisabled returns (uint256) {
    int256 weiPerUnitLink = getFeedData();
    return calculateRequestPriceInternal(_callbackGasLimit, tx.gasprice, weiPerUnitLink);
  }

  function calculateRequestPriceNative(
    uint32 _callbackGasLimit
  ) external view override onlyConfiguredNotDisabled returns (uint256) {
    return calculateRequestPriceNativeInternal(_callbackGasLimit, tx.gasprice);
  }

  /**
   * @notice Estimates the price of a VRF request with a specific gas limit and gas price.
   *
   * @dev This is a convenience function that can be called in simulation to better understand
   * @dev pricing.
   *
   * @param _callbackGasLimit is the gas limit used to estimate the price.
   * @param _requestGasPriceWei is the gas price in wei used for the estimation.
   */
  function estimateRequestPrice(
    uint32 _callbackGasLimit,
    uint256 _requestGasPriceWei
  ) external view override onlyConfiguredNotDisabled returns (uint256) {
    int256 weiPerUnitLink = getFeedData();
    return calculateRequestPriceInternal(_callbackGasLimit, _requestGasPriceWei, weiPerUnitLink);
  }

  function estimateRequestPriceNative(
    uint32 _callbackGasLimit,
    uint256 _requestGasPriceWei
  ) external view override onlyConfiguredNotDisabled returns (uint256) {
    return calculateRequestPriceNativeInternal(_callbackGasLimit, _requestGasPriceWei);
  }

  function calculateRequestPriceNativeInternal(uint256 _gas, uint256 _requestGasPrice) internal view returns (uint256) {
    // costWei is the base fee denominated in wei (native)
    // costWei takes into account the L1 posting costs of the VRF fulfillment
    // transaction, if we are on an L2.
    uint256 costWei = (_requestGasPrice *
      (_gas + s_wrapperGasOverhead + s_coordinatorGasOverhead) +
      ChainSpecificUtil.getL1CalldataGasCost(s_fulfillmentTxSizeBytes));
    // ((wei/gas * (gas)) + l1wei)
    // baseFee is the base fee denominated in wei
    uint256 baseFee = costWei;
    // feeWithPremium is the fee after the percentage premium is applied
    uint256 feeWithPremium = (baseFee * (s_wrapperPremiumPercentage + 100)) / 100;
    // feeWithFlatFee is the fee after the flat fee is applied on top of the premium
    uint256 feeWithFlatFee = feeWithPremium + (1e12 * uint256(s_fulfillmentFlatFeeNativePPM));

    return feeWithFlatFee;
  }

  function calculateRequestPriceInternal(
    uint256 _gas,
    uint256 _requestGasPrice,
    int256 _weiPerUnitLink
  ) internal view returns (uint256) {
    // costWei is the base fee denominated in wei (native)
    // costWei takes into account the L1 posting costs of the VRF fulfillment
    // transaction, if we are on an L2.
    uint256 costWei = (_requestGasPrice *
      (_gas + s_wrapperGasOverhead + s_coordinatorGasOverhead) +
      ChainSpecificUtil.getL1CalldataGasCost(s_fulfillmentTxSizeBytes));
    // (1e18 juels/link) * ((wei/gas * (gas)) + l1wei) / (wei/link) == 1e18 juels * wei/link / (wei/link) == 1e18 juels * wei/link * link/wei == juels
    // baseFee is the base fee denominated in juels (link)
    uint256 baseFee = (1e18 * costWei) / uint256(_weiPerUnitLink);
    // feeWithPremium is the fee after the percentage premium is applied
    uint256 feeWithPremium = (baseFee * (s_wrapperPremiumPercentage + 100)) / 100;
    // feeWithFlatFee is the fee after the flat fee is applied on top of the premium
    uint256 feeWithFlatFee = feeWithPremium + (1e12 * uint256(s_fulfillmentFlatFeeLinkPPM));

    return feeWithFlatFee;
  }

  /**
   * @notice onTokenTransfer is called by LinkToken upon payment for a VRF request.
   *
   * @dev Reverts if payment is too low.
   *
   * @param _sender is the sender of the payment, and the address that will receive a VRF callback
   *        upon fulfillment.
   *
   * @param _amount is the amount of LINK paid in Juels.
   *
   * @param _data is the abi-encoded VRF request parameters: uint32 callbackGasLimit,
   *        uint16 requestConfirmations, and uint32 numWords.
   */
  function onTokenTransfer(address _sender, uint256 _amount, bytes calldata _data) external onlyConfiguredNotDisabled {
    require(msg.sender == address(s_link), "only callable from LINK");

    (uint32 callbackGasLimit, uint16 requestConfirmations, uint32 numWords) = abi.decode(
      _data,
      (uint32, uint16, uint32)
    );
    uint32 eip150Overhead = getEIP150Overhead(callbackGasLimit);
    int256 weiPerUnitLink = getFeedData();
    uint256 price = calculateRequestPriceInternal(callbackGasLimit, tx.gasprice, weiPerUnitLink);
    require(_amount >= price, "fee too low");
    require(numWords <= s_maxNumWords, "numWords too high");
    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: s_keyHash,
      subId: SUBSCRIPTION_ID,
      requestConfirmations: requestConfirmations,
      callbackGasLimit: callbackGasLimit + eip150Overhead + s_wrapperGasOverhead,
      numWords: numWords,
      extraArgs: "" // empty extraArgs defaults to link payment
    });
    uint256 requestId = COORDINATOR.requestRandomWords(req);
    s_callbacks[requestId] = Callback({
      callbackAddress: _sender,
      callbackGasLimit: callbackGasLimit,
      requestGasPrice: tx.gasprice
    });
    lastRequestId = requestId;
  }

  function requestRandomWordsInNative(
    uint32 _callbackGasLimit,
    uint16 _requestConfirmations,
    uint32 _numWords
  ) external payable override returns (uint256 requestId) {
    uint32 eip150Overhead = getEIP150Overhead(_callbackGasLimit);
    uint256 price = calculateRequestPriceNativeInternal(_callbackGasLimit, tx.gasprice);
    require(msg.value >= price, "fee too low");
    require(_numWords <= s_maxNumWords, "numWords too high");
    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: s_keyHash,
      subId: SUBSCRIPTION_ID,
      requestConfirmations: _requestConfirmations,
      callbackGasLimit: _callbackGasLimit + eip150Overhead + s_wrapperGasOverhead,
      numWords: _numWords,
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: true}))
    });
    requestId = COORDINATOR.requestRandomWords(req);
    s_callbacks[requestId] = Callback({
      callbackAddress: msg.sender,
      callbackGasLimit: _callbackGasLimit,
      requestGasPrice: tx.gasprice
    });

    return requestId;
  }

  /**
   * @notice withdraw is used by the VRFV2Wrapper's owner to withdraw LINK revenue.
   *
   * @param _recipient is the address that should receive the LINK funds.
   *
   * @param _amount is the amount of LINK in Juels that should be withdrawn.
   */
  function withdraw(address _recipient, uint256 _amount) external onlyOwner {
    if (!s_link.transfer(_recipient, _amount)) {
      revert FailedToTransferLink();
    }
  }

  /**
   * @notice withdraw is used by the VRFV2Wrapper's owner to withdraw native revenue.
   *
   * @param _recipient is the address that should receive the native funds.
   *
   * @param _amount is the amount of native in Wei that should be withdrawn.
   */
  function withdrawNative(address _recipient, uint256 _amount) external onlyOwner {
    (bool success, ) = payable(_recipient).call{value: _amount}("");
    require(success, "failed to withdraw native");
  }

  /**
   * @notice enable this contract so that new requests can be accepted.
   */
  function enable() external onlyOwner {
    s_disabled = false;
  }

  /**
   * @notice disable this contract so that new requests will be rejected. When disabled, new requests
   * @notice will revert but existing requests can still be fulfilled.
   */
  function disable() external onlyOwner {
    s_disabled = true;
  }

  function fulfillRandomWords(uint256 _requestId, uint256[] memory _randomWords) internal override {
    Callback memory callback = s_callbacks[_requestId];
    delete s_callbacks[_requestId];
    require(callback.callbackAddress != address(0), "request not found"); // This should never happen

    VRFV2PlusWrapperConsumerBase c;
    bytes memory resp = abi.encodeWithSelector(c.rawFulfillRandomWords.selector, _requestId, _randomWords);

    bool success = callWithExactGas(callback.callbackGasLimit, callback.callbackAddress, resp);
    if (!success) {
      emit WrapperFulfillmentFailed(_requestId, callback.callbackAddress);
    }
  }

  function getFeedData() private view returns (int256) {
    bool staleFallback = s_stalenessSeconds > 0;
    uint256 timestamp;
    int256 weiPerUnitLink;
    (, weiPerUnitLink, , timestamp, ) = s_linkNativeFeed.latestRoundData();
    // solhint-disable-next-line not-rely-on-time
    if (staleFallback && s_stalenessSeconds < block.timestamp - timestamp) {
      weiPerUnitLink = s_fallbackWeiPerUnitLink;
    }
    require(weiPerUnitLink >= 0, "Invalid LINK wei price");
    return weiPerUnitLink;
  }

  /**
   * @dev Calculates extra amount of gas required for running an assembly call() post-EIP150.
   */
  function getEIP150Overhead(uint32 gas) private pure returns (uint32) {
    return gas / 63 + 1;
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

  function typeAndVersion() external pure virtual override returns (string memory) {
    return "VRFV2Wrapper 1.0.0";
  }

  modifier onlyConfiguredNotDisabled() {
    require(s_configured, "wrapper is not configured");
    require(!s_disabled, "wrapper is disabled");
    _;
  }
}

interface ExtendedVRFCoordinatorV2PlusInterface is IVRFCoordinatorV2Plus {
  function s_config()
    external
    view
    returns (
      uint16 minimumRequestConfirmations,
      uint32 maxGasLimit,
      uint32 stalenessSeconds,
      uint32 gasAfterPaymentCalculation
    );

  function s_fallbackWeiPerUnitLink() external view returns (int256);

  function s_feeConfig() external view returns (uint32 fulfillmentFlatFeeLinkPPM, uint32 fulfillmentFlatFeeNativePPM);
}
