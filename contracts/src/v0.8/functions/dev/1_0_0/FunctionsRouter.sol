// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {RouterBase, ITypeAndVersion} from "./RouterBase.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";
import {IFunctionsCoordinator} from "./interfaces/IFunctionsCoordinator.sol";
import {AuthorizedOriginReceiver} from "./accessControl/AuthorizedOriginReceiver.sol";
import {IFunctionsSubscriptions, FunctionsSubscriptions, IFunctionsBilling} from "./FunctionsSubscriptions.sol";

contract FunctionsRouter is RouterBase, IFunctionsRouter, AuthorizedOriginReceiver, FunctionsSubscriptions {
  event RequestStart(bytes32 indexed requestId, Request commitment);
  event RequestEnd(bytes32 indexed requestId, uint64 subscriptionId, bool success);

  error OnlyCallableFromCoordinator();

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
    bool useAllowList,
    address linkToken,
    bytes32[] memory initialJobIds,
    address[] memory initialAddresses,
    bytes memory config
  )
    RouterBase(msg.sender, timelockBlocks, maximumTimelockBlocks, initialJobIds, initialAddresses, config)
    AuthorizedOriginReceiver(useAllowList)
    FunctionsSubscriptions(linkToken)
  {}

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
    bytes32 jobId,
    bool isProposed,
    uint64 subscriptionId,
    bytes memory data,
    uint32 gasLimit
  ) internal returns (bytes32) {
    _isValidSubscription(subscriptionId);
    _isValidConsumer(msg.sender, subscriptionId);

    address route = this.getRoute(jobId, isProposed);
    IFunctionsCoordinator coordinator = IFunctionsCoordinator(route);

    (, , address owner, , ) = this.getSubscription(subscriptionId);

    (
      bytes32 requestId,
      uint96 estimatedCost,
      uint256 gasAfterPaymentCalculation,
      uint256 requestTimeoutSeconds
    ) = coordinator.sendRequest(subscriptionId, data, gasLimit, msg.sender, owner);

    _blockBalance(
      msg.sender,
      subscriptionId,
      gasLimit,
      estimatedCost,
      requestId,
      route,
      requestTimeoutSeconds,
      gasAfterPaymentCalculation,
      s_config.adminFee
    );

    emit RequestStart(requestId, s_requests[requestId]);

    return requestId;
  }

  function _smoke(bytes32 jobId, bytes calldata data) internal override onlyAuthorizedUsers returns (bytes32) {
    (uint64 subscriptionId, bytes memory reqData, uint32 gasLimit) = abi.decode(data, (uint64, bytes, uint32));
    return _sendRequest(jobId, true, subscriptionId, reqData, gasLimit);
  }

  /**
   * @inheritdoc IFunctionsRouter
   */
  function sendRequest(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 gasLimit,
    bytes32 jobId
  ) external override onlyAuthorizedUsers nonReentrant returns (bytes32) {
    return _sendRequest(jobId, false, subscriptionId, data, gasLimit);
  }

  /**
   * @inheritdoc IFunctionsRouter
   */
  function fulfill(
    bytes32 requestId,
    bytes memory response,
    bytes memory err,
    uint96 juelsPerGas,
    address transmitter,
    address[] memory to,
    uint96[] memory amount
  ) external override nonReentrant returns (FulfillResult) {
    Request memory request = s_requests[requestId];

    if (request.client == address(0)) {
      return FulfillResult.INVALID_REQUEST_ID;
    }
    if (msg.sender != request.coordinator) {
      revert OnlyCallableFromCoordinator();
    }

    _checkBalance(
      request.estimatedCost,
      request.gasAfterPaymentCalculation,
      request.adminFee,
      request.gasLimit,
      juelsPerGas,
      to,
      amount
    );

    delete s_requests[requestId];

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
    (bool success, uint256 gasUsed) = callWithExactGas(request.gasLimit, request.client, encodedCallback);
    s_reentrancyLock = false;

    _pay(
      request.subscriptionId,
      request.estimatedCost,
      request.client,
      request.adminFee,
      this.owner(),
      transmitter,
      juelsPerGas,
      gasUsed,
      to,
      amount
    );

    // TODO: payment amounts
    emit RequestEnd(requestId, request.subscriptionId, success);

    return success ? FulfillResult.USER_SUCCESS : FulfillResult.USER_ERROR;
  }

  /**
   * @dev calls target address with exactly gasAmount gas and data as calldata
   * or reverts if at least gasAmount gas is not available.
   */
  function callWithExactGas(
    uint256 gasAmount,
    address target,
    bytes memory data
  ) private returns (bool success, uint256 gasUsed) {
    // solhint-disable-next-line no-inline-assembly
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
      gasUsed := sub(g, gas())
    }
    return (success, gasUsed);
  }

  // ================================================================
  // |                           Modifiers                          |
  // ================================================================

  function _canSetAuthorizedSenders() internal view override onlyOwner returns (bool) {
    return msg.sender == owner();
  }

  modifier onlyAuthorizedUsers() override {
    _validateIsAuthorizedSender();
    _;
  }

  modifier onlyRouterOwner() override {
    _validateOwnership();
    _;
  }
}
