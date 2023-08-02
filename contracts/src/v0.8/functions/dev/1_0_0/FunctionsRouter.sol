// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {RouterBase, ITypeAndVersion} from "./RouterBase.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";
import {IFunctionsCoordinator} from "./interfaces/IFunctionsCoordinator.sol";
import {FunctionsSubscriptions} from "./FunctionsSubscriptions.sol";
import {FulfillResult} from "./interfaces/FulfillResultCodes.sol";
import {IAccessController} from "../../../shared/interfaces/IAccessController.sol";
import {IFunctionsRequest} from "./interfaces/IFunctionsRequest.sol";
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

  event RequestProcessed(
    bytes32 indexed requestId,
    uint64 indexed subscriptionId,
    uint96 totalCostJuels,
    address transmitter,
    uint8 resultCode,
    bytes response,
    bytes returnData
  );

  event RequestNotProcessed(bytes32 indexed requestId, address coordinator, address transmitter, uint8 resultCode);

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
  uint8 private constant GAS_FLAG_INDEX = 0;

  // ================================================================
  // |                    Configuration state                       |
  // ================================================================
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
  string public constant override typeAndVersion = "Functions Router v1.0.0";

  /**
   * @inheritdoc IFunctionsRouter
   */
  function getAllowListId() external pure override returns (bytes32) {
    return ALLOW_LIST_ID;
  }

  /**
   * @inheritdoc IFunctionsRouter
   */
  function getConfig() external view override returns (Config memory) {
    return s_config;
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
    (
      uint16 maxConsumers,
      uint96 adminFee,
      bytes4 handleOracleFulfillmentSelector,
      uint32[] memory maxCallbackGasLimits
    ) = abi.decode(config, (uint16, uint96, bytes4, uint32[]));
    s_config = Config({
      maxConsumers: maxConsumers,
      adminFee: adminFee,
      handleOracleFulfillmentSelector: handleOracleFulfillmentSelector,
      maxCallbackGasLimits: maxCallbackGasLimits
    });
    emit ConfigChanged(adminFee, handleOracleFulfillmentSelector, maxCallbackGasLimits);
  }

  // ================================================================
  // |                      Request methods                         |
  // ================================================================

  function _sendRequest(
    bytes32 donId,
    bool useProposed,
    uint64 subscriptionId,
    bytes memory data,
    uint16 dataVersion,
    uint32 callbackGasLimit
  ) private returns (bytes32 requestId) {
    _whenNotPaused();
    _isValidSubscription(subscriptionId);
    _isValidConsumer(msg.sender, subscriptionId);
    isValidCallbackGasLimit(subscriptionId, callbackGasLimit);

    address coordinatorAddress = _getContractById(donId, useProposed);

    // Forward request to DON
    IFunctionsRequest.Commitment memory commitment = IFunctionsCoordinator(coordinatorAddress).sendRequest(
      IFunctionsCoordinator.Request(
        msg.sender,
        s_subscriptions[subscriptionId].owner,
        data,
        subscriptionId,
        dataVersion,
        _getFlags(subscriptionId),
        callbackGasLimit,
        s_config.adminFee
      )
    );

    _markRequestInFlight(msg.sender, subscriptionId, commitment.estimatedTotalCostJuels);

    // Store a commitment about the request
    s_requestCommitments[commitment.requestId] = keccak256(
      abi.encode(
        IFunctionsRequest.Commitment({
          adminFee: s_config.adminFee,
          coordinator: coordinatorAddress,
          client: msg.sender,
          subscriptionId: subscriptionId,
          callbackGasLimit: callbackGasLimit,
          estimatedTotalCostJuels: commitment.estimatedTotalCostJuels,
          timeoutTimestamp: commitment.timeoutTimestamp,
          requestId: commitment.requestId,
          donFee: commitment.donFee,
          gasOverheadBeforeCallback: commitment.gasOverheadBeforeCallback,
          gasOverheadAfterCallback: commitment.gasOverheadAfterCallback
        })
      )
    );

    emit RequestStart(
      commitment.requestId,
      donId,
      subscriptionId,
      s_subscriptions[subscriptionId].owner,
      msg.sender,
      tx.origin,
      data,
      dataVersion,
      callbackGasLimit
    );

    return commitment.requestId;
  }

  function _validateProposedContracts(
    bytes32 donId,
    bytes calldata data
  ) internal override returns (bytes memory output) {
    (uint64 subscriptionId, bytes memory reqData, uint16 reqDataVersion, uint32 callbackGasLimit) = abi.decode(
      data,
      (uint64, bytes, uint16, uint32)
    );
    bytes32 requestId = _sendRequest(donId, true, subscriptionId, reqData, reqDataVersion, callbackGasLimit);
    // Convert to bytes as a more generic return
    output = new bytes(32);
    // Bounded by 32
    for (uint256 i; i < 32; ++i) {
      output[i] = requestId[i];
    }
  }

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
    return _sendRequest(donId, false, subscriptionId, data, dataVersion, callbackGasLimit);
  }

  /**
   * @inheritdoc IFunctionsRouter
   */
  function fulfill(
    bytes memory response,
    bytes memory err,
    uint96 juelsPerGas,
    uint96 costWithoutFulfillment,
    address transmitter,
    IFunctionsRequest.Commitment memory commitment
  ) external override returns (uint8 resultCode, uint96 callbackGasCostJuels) {
    _whenNotPaused();

    if (msg.sender != commitment.coordinator) {
      revert OnlyCallableFromCoordinator();
    }

    if (s_requestCommitments[commitment.requestId] == bytes32(0)) {
      resultCode = uint8(FulfillResult.INVALID_REQUEST_ID);
      emit RequestNotProcessed(commitment.requestId, commitment.coordinator, transmitter, resultCode);
      return (resultCode, callbackGasCostJuels);
    }

    if (keccak256(abi.encode(commitment)) != s_requestCommitments[commitment.requestId]) {
      resultCode = uint8(FulfillResult.INVALID_COMMITMENT);
      emit RequestNotProcessed(commitment.requestId, commitment.coordinator, transmitter, resultCode);
      return (resultCode, callbackGasCostJuels);
    }

    // Check that the transmitter has supplied enough gas for the callback to succeed
    if (gasleft() < commitment.callbackGasLimit + commitment.gasOverheadAfterCallback) {
      resultCode = uint8(FulfillResult.INSUFFICIENT_GAS_PROVIDED);
      emit RequestNotProcessed(commitment.requestId, commitment.coordinator, transmitter, resultCode);
      return (resultCode, callbackGasCostJuels);
    }

    // Check that the subscription can still afford
    if (
      commitment.adminFee + costWithoutFulfillment + (juelsPerGas * SafeCast.toUint96(commitment.callbackGasLimit)) >
      s_subscriptions[commitment.subscriptionId].balance
    ) {
      resultCode = uint8(FulfillResult.SUBSCRIPTION_BALANCE_INVARIANT_VIOLATION);
      delete s_requestCommitments[commitment.requestId];
      emit RequestNotProcessed(commitment.requestId, commitment.coordinator, transmitter, resultCode);
      return (resultCode, callbackGasCostJuels);
    }

    // Check that the cost has not exceeded the quoted cost
    if (
      commitment.adminFee + costWithoutFulfillment + (juelsPerGas * SafeCast.toUint96(commitment.callbackGasLimit)) >
      commitment.estimatedTotalCostJuels
    ) {
      resultCode = uint8(FulfillResult.COST_EXCEEDS_COMMITMENT);
      emit RequestNotProcessed(commitment.requestId, commitment.coordinator, transmitter, resultCode);
      return (resultCode, callbackGasCostJuels);
    }

    // If checks pass, continue as default, 0 = USER_SUCCESS;
    resultCode = 0;

    delete s_requestCommitments[commitment.requestId];

    CallbackResult memory result = _callback(
      commitment.requestId,
      response,
      err,
      commitment.callbackGasLimit,
      commitment.client
    );
    resultCode = result.success
      ? 0 // FulfillResult.USER_SUCCESS
      : 1; // FulfillResult.USER_ERROR

    Receipt memory receipt = _pay(
      commitment.subscriptionId,
      commitment.estimatedTotalCostJuels,
      commitment.client,
      commitment.adminFee,
      juelsPerGas,
      result.gasUsed,
      costWithoutFulfillment
    );

    // resultCode must be 0 or 1 here
    emit RequestProcessed(
      commitment.requestId,
      commitment.subscriptionId,
      receipt.totalCostJuels,
      transmitter,
      resultCode,
      result.success ? response : err,
      result.returnData
    );

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
  function isValidCallbackGasLimit(uint64 subscriptionId, uint32 callbackGasLimit) public view {
    uint8 index = uint8(_getFlags(subscriptionId)[GAS_FLAG_INDEX]);
    if (index >= s_config.maxCallbackGasLimits.length) {
      revert InvalidGasFlagValue(index);
    }
    if (callbackGasLimit > s_config.maxCallbackGasLimits[index]) {
      revert GasLimitTooBig(s_config.maxCallbackGasLimits[index]);
    }
  }

  function _getMaxConsumers() internal view override returns (uint16) {
    return s_config.maxConsumers;
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
    if (!IAccessController(_getContractById(ALLOW_LIST_ID, false)).hasAccess(msg.sender, new bytes(0))) {
      revert SenderMustAcceptTermsOfService(msg.sender);
    }
  }
}
