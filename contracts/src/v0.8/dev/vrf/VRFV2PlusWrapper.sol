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
import "../interfaces/IVRFV2PlusPriceRegistry.sol";

/**
 * @notice A wrapper for VRFCoordinatorV2 that provides an interface better suited to one-off
 * @notice requests for randomness.
 */
contract VRFV2PlusWrapper is ConfirmedOwner, TypeAndVersionInterface, VRFConsumerBaseV2Plus, VRFV2PlusWrapperInterface {
  event WrapperFulfillmentFailed(uint256 indexed requestId, address indexed consumer);

  error LinkAlreadySet();

  LinkTokenInterface public s_link;
  ExtendedVRFCoordinatorV2PlusInterface public immutable COORDINATOR;
  IVRFV2PlusPriceRegistry public immutable PRICE_REGISTRY;
  uint256 public immutable SUBSCRIPTION_ID;

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

  // Other configuration
  // s_wrapperGasOverhead reflects the gas overhead of the wrapper's fulfillRandomWords
  // function. The cost for this gas is passed to the user.
  uint32 s_wrapperGasOverhead;

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

  constructor(address _link, address _coordinator) VRFConsumerBaseV2Plus(_coordinator) {
    if (_link != address(0)) {
      s_link = LinkTokenInterface(_link);
    }

    COORDINATOR = ExtendedVRFCoordinatorV2PlusInterface(_coordinator);
    PRICE_REGISTRY = IVRFV2PlusPriceRegistry(ExtendedVRFCoordinatorV2PlusInterface(_coordinator).PRICE_REGISTRY());

    // Create this wrapper's subscription and add itself as a consumer.
    uint256 subId = ExtendedVRFCoordinatorV2PlusInterface(_coordinator).createSubscription();
    SUBSCRIPTION_ID = subId;
    ExtendedVRFCoordinatorV2PlusInterface(_coordinator).addConsumer(subId, address(this));
  }

  /**
   * @inheritdoc VRFV2PlusWrapperInterface
   */
  function getPriceRegistry() external view override returns (address) {
    return address(PRICE_REGISTRY);
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
   * @notice setConfig configures VRFV2Wrapper.
   *
   * @dev Sets wrapper-specific configuration based on the given parameters, and fetches any needed
   * @dev VRFCoordinatorV2 configuration from the coordinator.
   *
   * @param _wrapperGasOverhead reflects the gas overhead of the wrapper's fulfillRandomWords
   * @param _keyHash to use for requesting randomness.
   * @param _maxNumWords is the max number of words that can be requested in a single wrapped VRF request
   */
  function setConfig(uint32 _wrapperGasOverhead, bytes32 _keyHash, uint8 _maxNumWords) external onlyOwner {
    s_wrapperGasOverhead = _wrapperGasOverhead;
    s_keyHash = _keyHash;
    s_maxNumWords = _maxNumWords;
    s_configured = true;
  }

  /**
   * @notice getConfig returns the current VRFV2Wrapper configuration.
   *
   * @return keyHash is the key hash to use when requesting randomness. Fees are paid based on
   *         current gas fees, so this should be set to the highest gas lane on the network.
   *
   * @return maxNumWords is the max number of words that can be requested in a single wrapped VRF
   *         request.
   */
  function getConfig() external view returns (bytes32 keyHash, uint8 maxNumWords) {
    return (s_keyHash, s_maxNumWords);
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
    uint256 price = PRICE_REGISTRY.calculateRequestPriceWrapper(callbackGasLimit);
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
    uint256 price = PRICE_REGISTRY.calculateRequestPriceNativeWrapper(_callbackGasLimit);
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
    s_link.transfer(_recipient, _amount);
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
  function PRICE_REGISTRY() external view returns (address);
}
