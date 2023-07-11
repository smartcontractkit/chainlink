// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {RouterBase, ITypeAndVersion} from "./RouterBase.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";
import {IFunctionsCoordinator} from "./interfaces/IFunctionsCoordinator.sol";
import {IFunctionsSubscriptions, FunctionsSubscriptions, IFunctionsBilling} from "./FunctionsSubscriptions.sol";

contract FunctionsRouter is RouterBase, IFunctionsRouter, FunctionsSubscriptions {
  event RequestStart(bytes32 indexed requestId, Request request);
  event RequestEnd(
    bytes32 indexed requestId,
    uint64 indexed subscriptionId,
    uint96 totalCostJuels,
    address transmitter,
    uint8 resultCode,
    bytes response
  );

  error OnlyCallableFromCoordinator();

  struct CallbackResult {
    bool success;
    uint256 gasUsed;
  }

  // ================================================================
  // |                    Configuration state                       |
  // ================================================================
  struct Config {
    // Flat fee (in Juels of LINK) that will be paid to the Router owner for operation of the network
    uint96 adminFee;
    // The function selector that is used when calling back to the Client contract
    bytes4 handleOracleFulfillmentSelector;
  }
  Config private s_config;
  event ConfigSet(uint96 adminFee, bytes4 handleOracleFulfillmentSelector);

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

  /**
   * @inheritdoc ITypeAndVersion
   */
  function typeAndVersion() public pure override returns (string memory) {
    return "Functions Router v1";
  }

  // ================================================================
  // |                 Configuration methods                        |
  // ================================================================
  /**
   * @notice Sets the configuration for FunctionsRouter specific state
   * @param config bytes of config data to set the following:
   *  - adminFee: fee that will be paid to the Router owner for operating the network
   */
  function _setConfig(bytes memory config) internal override {
    (uint96 adminFee, bytes4 handleOracleFulfillmentSelector) = abi.decode(config, (uint96, bytes4));
    s_config = Config({adminFee: adminFee, handleOracleFulfillmentSelector: handleOracleFulfillmentSelector});
    emit ConfigSet(adminFee, handleOracleFulfillmentSelector);
  }

  /**
   * @inheritdoc IFunctionsRouter
   */
  function getAdminFee() external view override returns (uint96) {
    return s_config.adminFee;
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
  ) internal returns (bytes32) {
    _isValidSubscription(subscriptionId);
    _isValidConsumer(msg.sender, subscriptionId);

    address route = this.getRoute(donId, useProposed);
    IFunctionsCoordinator coordinator = IFunctionsCoordinator(route);

    (, , address owner, , ) = this.getSubscription(subscriptionId);

    (
      bytes32 requestId,
      uint96 estimatedCost,
      uint256 gasAfterPaymentCalculation, // Used to ensure that the transmitter supplies enough gas
      uint256 requestTimeoutSeconds
    ) = coordinator.sendRequest(
        IFunctionsCoordinator.Request(subscriptionId, data, dataVersion, callbackGasLimit, msg.sender, owner)
      );

    _blockBalance(msg.sender, subscriptionId, estimatedCost);

    // Store a commitment about the request
    s_requests[requestId] = Request(
      route,
      msg.sender,
      subscriptionId,
      callbackGasLimit,
      estimatedCost,
      block.timestamp + requestTimeoutSeconds,
      gasAfterPaymentCalculation,
      s_config.adminFee
    );

    emit RequestStart(requestId, s_requests[requestId]);

    return requestId;
  }

  function _validateProposal(bytes32 donId, bytes calldata data) internal override returns (bytes32) {
    (uint64 subscriptionId, bytes memory reqData, uint16 reqDataVersion, uint32 callbackGasLimit) = abi.decode(
      data,
      (uint64, bytes, uint16, uint32)
    );
    return _sendRequest(donId, true, subscriptionId, reqData, reqDataVersion, callbackGasLimit);
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
  ) external override nonReentrant returns (bytes32) {
    return _sendRequest(donId, false, subscriptionId, data, dataVersion, callbackGasLimit);
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
  ) external override nonReentrant returns (uint8 resultCode, uint96 callbackGasCostJuels) {
    Request memory request = s_requests[requestId];

    if (request.client == address(0)) {
      resultCode = 2; // FulfillResult.INVALID_REQUEST_ID
      return (resultCode, callbackGasCostJuels);
    }
    if (msg.sender != request.coordinator) {
      revert OnlyCallableFromCoordinator();
    }

    resultCode = _checkBalance(
      request.estimatedCost,
      request.gasAfterPaymentCalculation,
      request.adminFee,
      request.callbackGasLimit,
      juelsPerGas,
      costWithoutFulfillment
    );

    delete s_requests[requestId];

    CallbackResult memory result = _callback(requestId, response, err, request.callbackGasLimit, request.client);
    resultCode = result.success
      ? 0 // FulfillResult.USER_SUCCESS
      : 1; // FulfillResult.USER_ERROR

    Receipt memory receipt = _pay(
      request.subscriptionId,
      request.estimatedCost,
      request.client,
      request.adminFee,
      juelsPerGas,
      result.gasUsed,
      costWithoutFulfillment
    );

    emit RequestEnd(
      requestId,
      request.subscriptionId,
      receipt.totalCostJuels,
      transmitter,
      resultCode,
      result.success ? response : err // TODO: handle more response data scenarios
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
    // Do not allow any non-view/non-pure coordinator functions to be called
    // during the consumers callback code via reentrancyLock.
    // NOTE: that callWithExactGas will revert if we do not have sufficient gas
    // to give the callee their requested amount.
    s_reentrancyLock = true;
    bool success;
    uint256 gasUsed;
    assembly {
      let g := gas()
      // GAS_FOR_CALL_EXACT_CHECK = 5000
      // Compute g -= GAS_FOR_CALL_EXACT_CHECK and check for underflow
      // The gas actually passed to the callee is min(gasAmount, 63//64*gas available).
      // We want to ensure that we revert if gasAmount >  63//64*gas available
      // as we do not want to provide them with less, however that check itself costs
      // gas.  GAS_FOR_CALL_EXACT_CHECK ensures we have at least enough gas to be able
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
      // solidity calls check that a contract actually exists at the destination, so we do the same
      if iszero(extcodesize(client)) {
        revert(0, 0)
      }
      // call and return whether we succeeded. ignore return data
      // call(gas,addr,value,argsOffset,argsLength,retOffset,retLength)
      success := call(callbackGasLimit, client, 0, add(encodedCallback, 0x20), mload(encodedCallback), 0, 0)
      gasUsed := sub(g, gas())
    }
    s_reentrancyLock = false;

    result = CallbackResult(success, gasUsed);
  }

  // ================================================================
  // |                           Modifiers                          |
  // ================================================================

  modifier onlyRouterOwner() override {
    _validateOwnership();
    _;
  }
}
