// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";
import {IFunctionsCoordinator} from "./interfaces/IFunctionsCoordinator.sol";
import {IAccessController} from "../../../shared/interfaces/IAccessController.sol";

import {RouterBase} from "./RouterBase.sol";
import {FunctionsSubscriptions} from "./FunctionsSubscriptions.sol";
import {FunctionsResponse} from "./libraries/FunctionsResponse.sol";

import {SafeCast} from "../../../vendor/openzeppelin-solidity/v4.8.0/contracts/utils/SafeCast.sol";

contract FunctionsRouter is RouterBase, IFunctionsRouter, FunctionsSubscriptions {
  using FunctionsResponse for FunctionsResponse.Commitment;
  using FunctionsResponse for FunctionsResponse.FulfillResult;

  // @inheritdoc ITypeAndVersion
  string public constant override typeAndVersion = "Functions Router v1.0.0";

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
    FunctionsResponse.FulfillResult resultCode,
    bytes response,
    bytes returnData
  );

  event RequestNotProcessed(
    bytes32 indexed requestId,
    address coordinator,
    address transmitter,
    FunctionsResponse.FulfillResult resultCode
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

  // @inheritdoc IFunctionsRouter
  function getAllowListId() external pure override returns (bytes32) {
    return ALLOW_LIST_ID;
  }

  // @inheritdoc IFunctionsRouter
  function getConfig()
    external
    view
    override
    returns (
      uint16 maxConsumersPerSubscription,
      uint96 adminFee,
      bytes4 handleOracleFulfillmentSelector,
      uint32[] memory maxCallbackGasLimits
    )
  {
    maxConsumersPerSubscription = s_config.maxConsumersPerSubscription;
    adminFee = s_config.adminFee;
    handleOracleFulfillmentSelector = s_config.handleOracleFulfillmentSelector;
    maxCallbackGasLimits = s_config.maxCallbackGasLimits;

    return (maxConsumersPerSubscription, adminFee, handleOracleFulfillmentSelector, maxCallbackGasLimits);
  }

  // ================================================================
  // |                 Configuration methods                        |
  // ================================================================

  // @notice Sets the configuration for FunctionsRouter specific state
  // @param config bytes of config data to set the following:
  // - adminFee: fee that will be paid to the Router owner for operating the network
  function _updateConfig(bytes memory config) internal override {
    (
      uint16 maxConsumersPerSubscription,
      uint96 adminFee,
      bytes4 handleOracleFulfillmentSelector,
      uint32[] memory maxCallbackGasLimits
    ) = abi.decode(config, (uint16, uint96, bytes4, uint32[]));
    s_config = Config({
      maxConsumersPerSubscription: maxConsumersPerSubscription,
      adminFee: adminFee,
      handleOracleFulfillmentSelector: handleOracleFulfillmentSelector,
      maxCallbackGasLimits: maxCallbackGasLimits
    });
    emit ConfigChanged(adminFee, handleOracleFulfillmentSelector, maxCallbackGasLimits);
  }

  // ================================================================
  // |                      Request methods                         |
  // ================================================================

  // @inheritdoc IFunctionsRouter
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

  // @inheritdoc IFunctionsRouter
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
  ) private returns (bytes32) {
    _whenNotPaused();
    _isValidSubscription(subscriptionId);
    _isValidConsumer(msg.sender, subscriptionId);
    isValidCallbackGasLimit(subscriptionId, callbackGasLimit);

    // Forward request to DON
    FunctionsResponse.Commitment memory commitment = coordinator.sendRequest(
      IFunctionsCoordinator.Request({
        requestingContract: msg.sender,
        subscriptionOwner: s_subscriptions[subscriptionId].owner,
        data: data,
        subscriptionId: subscriptionId,
        dataVersion: dataVersion,
        flags: getFlags(subscriptionId),
        callbackGasLimit: callbackGasLimit,
        adminFee: s_config.adminFee
      })
    );

    // Store a commitment about the request
    s_requestCommitments[commitment.requestId] = keccak256(
      abi.encode(
        FunctionsResponse.Commitment({
          adminFee: s_config.adminFee,
          coordinator: address(coordinator),
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

    _markRequestInFlight(msg.sender, subscriptionId, commitment.estimatedTotalCostJuels);

    emit RequestStart({
      requestId: commitment.requestId,
      donId: donId,
      subscriptionId: subscriptionId,
      subscriptionOwner: s_subscriptions[subscriptionId].owner,
      requestingContract: msg.sender,
      requestInitiator: tx.origin,
      data: data,
      dataVersion: dataVersion,
      callbackGasLimit: callbackGasLimit
    });

    return commitment.requestId;
  }

  // @inheritdoc IFunctionsRouter
  function fulfill(
    bytes memory response,
    bytes memory err,
    uint96 juelsPerGas,
    uint96 costWithoutCallback,
    address transmitter,
    FunctionsResponse.Commitment memory commitment
  ) external override returns (FunctionsResponse.FulfillResult resultCode, uint96) {
    _whenNotPaused();

    if (msg.sender != commitment.coordinator) {
      revert OnlyCallableFromCoordinator();
    }

    if (s_requestCommitments[commitment.requestId] == bytes32(0)) {
      resultCode = FunctionsResponse.FulfillResult.INVALID_REQUEST_ID;
      emit RequestNotProcessed(commitment.requestId, commitment.coordinator, transmitter, resultCode);
      return (resultCode, 0);
    }

    if (keccak256(abi.encode(commitment)) != s_requestCommitments[commitment.requestId]) {
      resultCode = FunctionsResponse.FulfillResult.INVALID_COMMITMENT;
      emit RequestNotProcessed(commitment.requestId, commitment.coordinator, transmitter, resultCode);
      return (resultCode, 0);
    }

    // Check that the transmitter has supplied enough gas for the callback to succeed
    if (gasleft() < commitment.callbackGasLimit + commitment.gasOverheadAfterCallback) {
      resultCode = FunctionsResponse.FulfillResult.INSUFFICIENT_GAS_PROVIDED;
      emit RequestNotProcessed(commitment.requestId, commitment.coordinator, transmitter, resultCode);
      return (resultCode, 0);
    }

    {
      uint96 callbackCost = juelsPerGas * SafeCast.toUint96(commitment.callbackGasLimit);
      uint96 totalCostJuels = commitment.adminFee + costWithoutCallback + callbackCost;

      // Check that the subscription can still afford
      if (totalCostJuels > s_subscriptions[commitment.subscriptionId].balance) {
        resultCode = FunctionsResponse.FulfillResult.SUBSCRIPTION_BALANCE_INVARIANT_VIOLATION;
        emit RequestNotProcessed(commitment.requestId, commitment.coordinator, transmitter, resultCode);
        return (resultCode, 0);
      }

      // Check that the cost has not exceeded the quoted cost
      if (totalCostJuels > commitment.estimatedTotalCostJuels) {
        resultCode = FunctionsResponse.FulfillResult.COST_EXCEEDS_COMMITMENT;
        emit RequestNotProcessed(commitment.requestId, commitment.coordinator, transmitter, resultCode);
        return (resultCode, 0);
      }
    }

    delete s_requestCommitments[commitment.requestId];

    CallbackResult memory result = _callback(
      commitment.requestId,
      response,
      err,
      commitment.callbackGasLimit,
      commitment.client
    );

    resultCode = result.success
      ? FunctionsResponse.FulfillResult.USER_SUCCESS
      : FunctionsResponse.FulfillResult.USER_ERROR;

    Receipt memory receipt = _pay(
      commitment.subscriptionId,
      commitment.estimatedTotalCostJuels,
      commitment.client,
      commitment.adminFee,
      juelsPerGas,
      SafeCast.toUint96(result.gasUsed),
      costWithoutCallback
    );

    emit RequestProcessed({
      requestId: commitment.requestId,
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
  ) private returns (CallbackResult memory) {
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

    return CallbackResult({success: success, gasUsed: gasUsed, returnData: returnData});
  }

  // @inheritdoc IFunctionsRouter
  function isValidCallbackGasLimit(uint64 subscriptionId, uint32 callbackGasLimit) public view {
    uint8 callbackGasLimitsIndexSelector = uint8(getFlags(subscriptionId)[MAX_CALLBACK_GAS_LIMIT_FLAGS_INDEX]);
    if (callbackGasLimitsIndexSelector >= s_config.maxCallbackGasLimits.length) {
      revert InvalidGasFlagValue(callbackGasLimitsIndexSelector);
    }
    uint32 maxCallbackGasLimit = s_config.maxCallbackGasLimits[callbackGasLimitsIndexSelector];
    if (callbackGasLimit > maxCallbackGasLimit) {
      revert GasLimitTooBig(maxCallbackGasLimit);
    }
  }

  function _getMaxConsumers() internal view override returns (uint16) {
    return s_config.maxConsumersPerSubscription;
  }

  // ================================================================
  // |                           Modifiers                          |
  // ================================================================
  // Favoring internal functions over actual modifiers to reduce contract size

  // Used within FunctionsSubscriptions.sol
  function _whenNotPaused() internal view override {
    _requireNotPaused();
  }

  // Used within FunctionsSubscriptions.sol
  function _onlyRouterOwner() internal view override {
    _validateOwnership();
  }

  // Used within FunctionsSubscriptions.sol
  function _onlySenderThatAcceptedToS() internal view override {
    if (!IAccessController(getContractById(ALLOW_LIST_ID)).hasAccess(msg.sender, new bytes(0))) {
      revert SenderMustAcceptTermsOfService(msg.sender);
    }
  }
}
