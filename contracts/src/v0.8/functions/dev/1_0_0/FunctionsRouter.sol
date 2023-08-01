// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {RouterBase, ITypeAndVersion} from "./RouterBase.sol";
import {FulfillResult} from "./interfaces/FulfillResultCodes.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";
import {IFunctionsCoordinator} from "./interfaces/IFunctionsCoordinator.sol";
import {FunctionsSubscriptions} from "./FunctionsSubscriptions.sol";
import {IAccessController} from "../../../shared/interfaces/IAccessController.sol";
import {SafeCast} from "../../../vendor/openzeppelin-solidity/v4.8.0/contracts/utils/SafeCast.sol";

contract FunctionsRouter is RouterBase, IFunctionsRouter, FunctionsSubscriptions {
  // We limit return data to a selector plus 4 words. This is to avoid
  // malicious contracts from returning large amounts of data and causing
  // repeated out-of-gas scenarios.
  uint16 public constant MAX_CALLBACK_RETURN_BYTES = 4 + 4 * 32;

  event RequestStart(
    bytes32 indexed requestId,
    bytes32 indexed donId,
    uint64 indexed subscriptionId,
    address subscriptionOwner,
    address requestingContract,
    address requestInitiator,
    bytes data,
    uint16 dataVersion,
    uint32 callbackGasLimit
  );

  event RequestEnd(
    bytes32 indexed requestId,
    uint64 indexed subscriptionId,
    uint96 totalCostJuels,
    address transmitter,
    FulfillResult resultCode,
    bytes response,
    bytes returnData
  );

  error OnlyCallableFromCoordinator();
  error SenderMustAcceptTermsOfService(address sender);
  error InvalidGasFlagValue(uint8 value);
  error GasLimitTooBig(uint32 limit);

  struct CallbackResult {
    bool success;
    uint256 gasUsed;
    bytes returnData;
  }

  // Identifier for the route to the Terms of Service Allow List
  bytes32 private constant ALLOW_LIST_ID = keccak256("Functions Terms of Service Allow List");
  uint8 private constant MAX_CALLBACK_GAS_LIMIT_FLAGS_INDEX = 0;

  // ================================================================
  // |                    Configuration state                       |
  // ================================================================
  struct Config {
    // Flat fee (in Juels of LINK) that will be paid to the Router owner for operation of the network
    uint96 adminFee;
    // The function selector that is used when calling back to the Client contract
    bytes4 handleOracleFulfillmentSelector;
    // List of max callback gas limits used by flag with MAX_CALLBACK_GAS_LIMIT_FLAGS_INDEX
    uint32[] maxCallbackGasLimits;
  }
  Config private s_config;
  event ConfigChanged(uint96 adminFee, bytes4 handleOracleFulfillmentSelector, uint32[] maxCallbackGasLimits);

  error OnlyCallableByRoute();

  // ================================================================
  // |                       Initialization                         |
  // ================================================================
  constructor(
    uint16 timelockBlocks,
    uint16 maximumTimelockBlocks,
    address linkToken,
    bytes memory config
  ) RouterBase(msg.sender, timelockBlocks, maximumTimelockBlocks, config) FunctionsSubscriptions(linkToken) {}

  // ================================================================
  // |                          Getters                             |
  // ================================================================
  /**
   * @inheritdoc ITypeAndVersion
   */
  function typeAndVersion() public pure override returns (string memory) {
    return "Functions Router v1.0.0";
  }

  /**
   * @inheritdoc IFunctionsRouter
   */
  function getAllowListId() external pure override returns (bytes32) {
    return ALLOW_LIST_ID;
  }

  /**
   * @inheritdoc IFunctionsRouter
   */
  function getConfig() external view override returns (uint96, bytes4, uint32[] memory) {
    return (s_config.adminFee, s_config.handleOracleFulfillmentSelector, s_config.maxCallbackGasLimits);
  }

  // ================================================================
  // |                 Configuration methods                        |
  // ================================================================
  /**
   * @notice Sets the configuration for FunctionsRouter specific state
   * @param config bytes of config data to set the following:
   *  - adminFee: fee that will be paid to the Router owner for operating the network
   */
  function _updateConfig(bytes memory config) internal override {
    (uint96 adminFee, bytes4 handleOracleFulfillmentSelector, uint32[] memory maxCallbackGasLimits) = abi.decode(
      config,
      (uint96, bytes4, uint32[])
    );
    s_config = Config({
      adminFee: adminFee,
      handleOracleFulfillmentSelector: handleOracleFulfillmentSelector,
      maxCallbackGasLimits: maxCallbackGasLimits
    });
    emit ConfigChanged(adminFee, handleOracleFulfillmentSelector, maxCallbackGasLimits);
  }

  // ================================================================
  // |                      Request methods                         |
  // ================================================================

  /**
   * @inheritdoc IFunctionsRouter
   */
  function sendRequest(
    uint64 subscriptionId,
    bytes calldata data,
    uint16 dataVersion,
    uint32 callbackGasLimit,
    bytes32 donId
  ) external override returns (bytes32) {
    IFunctionsCoordinator coordinator = IFunctionsCoordinator(getContractById(donId));
    return _sendRequest(donId, coordinator, subscriptionId, data, dataVersion, callbackGasLimit);
  }

  /**
   * @inheritdoc IFunctionsRouter
   */
  function sendRequestToProposed(
    uint64 subscriptionId,
    bytes calldata data,
    uint16 dataVersion,
    uint32 callbackGasLimit,
    bytes32 donId
  ) external override returns (bytes32) {
    IFunctionsCoordinator coordinator = IFunctionsCoordinator(getProposedContractById(donId));
    return _sendRequest(donId, coordinator, subscriptionId, data, dataVersion, callbackGasLimit);
  }

  function _sendRequest(
    bytes32 donId,
    IFunctionsCoordinator coordinator,
    uint64 subscriptionId,
    bytes memory data,
    uint16 dataVersion,
    uint32 callbackGasLimit
  ) private returns (bytes32 requestId) {
    _whenNotPaused();
    _isValidSubscription(subscriptionId);
    _isValidConsumer(msg.sender, subscriptionId);
    _isValidCallbackGasLimit(subscriptionId, callbackGasLimit);

    // Forward request to DON
    uint96 estimatedCost;
    uint256 gasAfterPaymentCalculation; // Used to ensure that the transmitter supplies enough gas
    uint256 requestTimeoutSeconds;
    (
      requestId,
      estimatedCost,
      gasAfterPaymentCalculation, // Used to ensure that the transmitter supplies enough gas
      requestTimeoutSeconds
    ) = coordinator.sendRequest(
      IFunctionsCoordinator.Request({
        requestingContract: msg.sender,
        subscriptionOwner: s_subscriptions[subscriptionId].owner,
        data: data,
        subscriptionId: subscriptionId,
        dataVersion: dataVersion,
        flags: _getFlags(subscriptionId),
        callbackGasLimit: callbackGasLimit
      })
    );

    _markRequestInFlight(msg.sender, subscriptionId, estimatedCost);

    // Store a commitment about the request
    s_requestCommitments[requestId] = Commitment({
      adminFee: s_config.adminFee,
      coordinator: address(coordinator),
      client: msg.sender,
      subscriptionId: subscriptionId,
      callbackGasLimit: callbackGasLimit,
      estimatedCost: estimatedCost,
      timeoutTimestamp: uint40(block.timestamp + requestTimeoutSeconds),
      gasAfterPaymentCalculation: uint120(gasAfterPaymentCalculation)
    });

    emit RequestStart({
      requestId: requestId,
      donId: donId,
      subscriptionId: subscriptionId,
      subscriptionOwner: s_subscriptions[subscriptionId].owner,
      requestingContract: msg.sender,
      requestInitiator: tx.origin,
      data: data,
      dataVersion: dataVersion,
      callbackGasLimit: callbackGasLimit
    });

    return requestId;
  }

  /**
   * @inheritdoc IFunctionsRouter
   */
  function fulfill(
    bytes32 requestId,
    bytes memory response,
    bytes memory err,
    uint96 juelsPerGas,
    uint96 costWithoutFulfillment,
    address transmitter
  ) external override returns (FulfillResult resultCode, uint96 callbackGasCostJuels) {
    _whenNotPaused();

    Commitment memory commitment = s_requestCommitments[requestId];

    if (msg.sender != commitment.coordinator) {
      revert OnlyCallableFromCoordinator();
    }

    if (commitment.client == address(0)) {
      return (FulfillResult.INVALID_REQUEST_ID, callbackGasCostJuels);
    }

    // Check that the transmitter has supplied enough gas for the callback to succeed
    if (gasleft() < commitment.callbackGasLimit + commitment.gasAfterPaymentCalculation) {
      return (FulfillResult.INSUFFICIENT_GAS, callbackGasCostJuels);
    }

    uint96 fulfillmentCostJuels = commitment.adminFee +
      costWithoutFulfillment +
      (juelsPerGas * SafeCast.toUint96(commitment.callbackGasLimit));

    // Check that the subscription can still afford
    if (fulfillmentCostJuels > s_subscriptions[commitment.subscriptionId].balance) {
      return (FulfillResult.INSUFFICIENT_SUBSCRIPTION_BALANCE, callbackGasCostJuels);
    }

    // Check that the cost has not exceeded the quoted cost
    if (fulfillmentCostJuels > commitment.estimatedCost) {
      return (FulfillResult.COST_EXCEEDS_COMMITMENT, callbackGasCostJuels);
    }

    delete s_requestCommitments[requestId];

    CallbackResult memory result = _callback(requestId, response, err, commitment.callbackGasLimit, commitment.client);
    resultCode = result.success ? FulfillResult.USER_SUCCESS : FulfillResult.USER_ERROR;

    Receipt memory receipt = _pay(
      commitment.subscriptionId,
      commitment.estimatedCost,
      commitment.client,
      commitment.adminFee,
      juelsPerGas,
      result.gasUsed,
      costWithoutFulfillment
    );

    emit RequestEnd({
      requestId: requestId,
      subscriptionId: commitment.subscriptionId,
      totalCostJuels: receipt.totalCostJuels,
      transmitter: transmitter,
      resultCode: resultCode,
      response: result.success ? response : err,
      returnData: result.returnData
    });

    return (resultCode, receipt.callbackGasCostJuels);
  }

  function _callback(
    bytes32 requestId,
    bytes memory response,
    bytes memory err,
    uint32 callbackGasLimit,
    address client
  ) private returns (CallbackResult memory result) {
    bytes memory encodedCallback = abi.encodeWithSelector(
      s_config.handleOracleFulfillmentSelector,
      requestId,
      response,
      err
    );

    // Call with explicitly the amount of callback gas requested
    // Important to not let them exhaust the gas budget and avoid payment.
    // NOTE: that callWithExactGas will revert if we do not have sufficient gas
    // to give the callee their requested amount.

    bool success;
    uint256 gasUsed;
    // allocate return data memory ahead of time
    bytes memory returnData = new bytes(MAX_CALLBACK_RETURN_BYTES);

    // solhint-disable-next-line no-inline-assembly
    assembly {
      // solidity calls check that a contract actually exists at the destination, so we do the same
      // Note we do this check prior to measuring gas so gasForCallExactCheck (our "cushion")
      // doesn't need to account for it.
      if iszero(extcodesize(client)) {
        revert(0, 0)
      }

      let g := gas()
      // GASFORCALLEXACTCHECK = 5000
      // Compute g -= gasForCallExactCheck and check for underflow
      // The gas actually passed to the callee is _min(gasAmount, 63//64*gas available).
      // We want to ensure that we revert if gasAmount >  63//64*gas available
      // as we do not want to provide them with less, however that check itself costs
      // gas. gasForCallExactCheck ensures we have at least enough gas to be able
      // to revert if gasAmount >  63//64*gas available.
      if lt(g, 5000) {
        revert(0, 0)
      }
      g := sub(g, 5000)
      // if g - g//64 <= gasAmount, revert
      // (we subtract g//64 because of EIP-150)
      if iszero(gt(sub(g, div(g, 64)), callbackGasLimit)) {
        revert(0, 0)
      }
      // call and  whether we succeeded
      // call(gas,addr,value,argsOffset,argsLength,retOffset,retLength)
      success := call(callbackGasLimit, client, 0, add(encodedCallback, 0x20), mload(encodedCallback), 0, 0)
      gasUsed := sub(g, gas())

      // limit our copy to MAX_CALLBACK_RETURN_BYTES bytes
      let toCopy := returndatasize()
      if gt(toCopy, MAX_CALLBACK_RETURN_BYTES) {
        toCopy := MAX_CALLBACK_RETURN_BYTES
      }
      // Store the length of the copied bytes
      mstore(returnData, toCopy)
      // copy the bytes from returnData[0:_toCopy]
      returndatacopy(add(returnData, 0x20), 0, toCopy)
    }

    result = CallbackResult(success, gasUsed, returnData);
  }

  /**
   * @inheritdoc IFunctionsRouter
   */
  function isValidCallbackGasLimit(uint64 subscriptionId, uint32 callbackGasLimit) external view {
    _isValidCallbackGasLimit(subscriptionId, callbackGasLimit);
  }

  function _isValidCallbackGasLimit(uint64 subscriptionId, uint32 callbackGasLimit) internal view {
    uint8 index = uint8(_getFlags(subscriptionId)[MAX_CALLBACK_GAS_LIMIT_FLAGS_INDEX]);
    if (index >= s_config.maxCallbackGasLimits.length) {
      revert InvalidGasFlagValue(index);
    }
    if (callbackGasLimit > s_config.maxCallbackGasLimits[index]) {
      revert GasLimitTooBig(s_config.maxCallbackGasLimits[index]);
    }
  }

  // ================================================================
  // |                           Modifiers                          |
  // ================================================================

  function _whenNotPaused() internal view override {
    _requireNotPaused();
  }

  function _onlyRouterOwner() internal view override {
    _validateOwnership();
  }

  function _onlySenderThatAcceptedToS() internal view override {
    if (!IAccessController(getContractById(ALLOW_LIST_ID)).hasAccess(msg.sender, new bytes(0))) {
      revert SenderMustAcceptTermsOfService(msg.sender);
    }
  }
}
