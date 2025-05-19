// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";
import {IFunctionsCoordinator} from "./interfaces/IFunctionsCoordinator.sol";
import {IAccessController} from "../../../shared/interfaces/IAccessController.sol";

import {FunctionsSubscriptions} from "./FunctionsSubscriptions.sol";
import {FunctionsResponse} from "./libraries/FunctionsResponse.sol";
import {ConfirmedOwner} from "../../../shared/access/ConfirmedOwner.sol";

import {SafeCast} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/math/SafeCast.sol";
import {Pausable} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/security/Pausable.sol";

contract FunctionsRouter is IFunctionsRouter, FunctionsSubscriptions, Pausable, ITypeAndVersion, ConfirmedOwner {
  using FunctionsResponse for FunctionsResponse.RequestMeta;
  using FunctionsResponse for FunctionsResponse.Commitment;
  using FunctionsResponse for FunctionsResponse.FulfillResult;

  string public constant override typeAndVersion = "Functions Router v2.0.0";

  // We limit return data to a selector plus 4 words. This is to avoid
  // malicious contracts from returning large amounts of data and causing
  // repeated out-of-gas scenarios.
  uint16 public constant MAX_CALLBACK_RETURN_BYTES = 4 + 4 * 32;
  uint8 private constant MAX_CALLBACK_GAS_LIMIT_FLAGS_INDEX = 0;

  event RequestStart(
    bytes32 indexed requestId,
    bytes32 indexed donId,
    uint64 indexed subscriptionId,
    address subscriptionOwner,
    address requestingContract,
    address requestInitiator,
    bytes data,
    uint16 dataVersion,
    uint32 callbackGasLimit,
    uint96 estimatedTotalCostJuels
  );

  event RequestProcessed(
    bytes32 indexed requestId,
    uint64 indexed subscriptionId,
    uint96 totalCostJuels,
    address transmitter,
    FunctionsResponse.FulfillResult resultCode,
    bytes response,
    bytes err,
    bytes callbackReturnData
  );

  event RequestNotProcessed(
    bytes32 indexed requestId,
    address coordinator,
    address transmitter,
    FunctionsResponse.FulfillResult resultCode
  );

  error EmptyRequestData();
  error OnlyCallableFromCoordinator();
  error SenderMustAcceptTermsOfService(address sender);
  error InvalidGasFlagValue(uint8 value);
  error GasLimitTooBig(uint32 limit);
  error DuplicateRequestId(bytes32 requestId);

  struct CallbackResult {
    bool success; // ══════╸ Whether the callback succeeded or not
    uint256 gasUsed; // ═══╸ The amount of gas consumed during the callback
    bytes returnData; // ══╸ The return of the callback function
  }

  // ================================================================
  // |                    Route state                       |
  // ================================================================

  mapping(bytes32 id => address routableContract) private s_route;

  error RouteNotFound(bytes32 id);

  // Identifier for the route to the Terms of Service Allow List
  bytes32 private s_allowListId;

  // ================================================================
  // |                    Configuration state                       |
  // ================================================================
  // solhint-disable-next-line gas-struct-packing
  struct Config {
    uint16 maxConsumersPerSubscription; // ═════════╗ Maximum number of consumers which can be added to a single subscription. This bound ensures we are able to loop over all subscription consumers as needed, without exceeding gas limits. Should a user require more consumers, they can use multiple subscriptions.
    uint72 adminFee; //                             ║ Flat fee (in Juels of LINK) that will be paid to the Router owner for operation of the network
    bytes4 handleOracleFulfillmentSelector; //      ║ The function selector that is used when calling back to the Client contract
    uint16 gasForCallExactCheck; // ════════════════╝ Used during calling back to the client. Ensures we have at least enough gas to be able to revert if gasAmount >  63//64*gas available.
    uint32[] maxCallbackGasLimits; // ══════════════╸ List of max callback gas limits used by flag with MAX_CALLBACK_GAS_LIMIT_FLAGS_INDEX
    uint16 subscriptionDepositMinimumRequests; //═══╗ Amount of requests that must be completed before the full subscription balance will be released when closing a subscription account.
    uint72 subscriptionDepositJuels; // ════════════╝ Amount of subscription funds that are held as a deposit until Config.subscriptionDepositMinimumRequests are made using the subscription.
  }

  Config private s_config;

  event ConfigUpdated(Config);

  // ================================================================
  // |                         Proposal state                       |
  // ================================================================

  uint8 private constant MAX_PROPOSAL_SET_LENGTH = 8;

  struct ContractProposalSet {
    bytes32[] ids; // ══╸ The IDs that key into the routes that will be modified if the update is applied
    address[] to; // ═══╸ The address of the contracts that the route will point to if the updated is applied
  }
  ContractProposalSet private s_proposedContractSet;

  event ContractProposed(
    bytes32 proposedContractSetId,
    address proposedContractSetFromAddress,
    address proposedContractSetToAddress
  );

  event ContractUpdated(bytes32 id, address from, address to);

  error InvalidProposal();
  error IdentifierIsReserved(bytes32 id);

  // ================================================================
  // |                       Initialization                         |
  // ================================================================

  constructor(
    address linkToken,
    Config memory config
  ) FunctionsSubscriptions(linkToken) ConfirmedOwner(msg.sender) Pausable() {
    // Set the intial configuration
    updateConfig(config);
  }

  // ================================================================
  // |                        Configuration                         |
  // ================================================================

  /// @notice The identifier of the route to retrieve the address of the access control contract
  // The access control contract controls which accounts can manage subscriptions
  /// @return id - bytes32 id that can be passed to the "getContractById" of the Router
  function getConfig() external view returns (Config memory) {
    return s_config;
  }

  /// @notice The router configuration
  function updateConfig(Config memory config) public onlyOwner {
    s_config = config;
    emit ConfigUpdated(config);
  }

  /// @inheritdoc IFunctionsRouter
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

  /// @inheritdoc IFunctionsRouter
  function getAdminFee() external view override returns (uint72) {
    return s_config.adminFee;
  }

  /// @inheritdoc IFunctionsRouter
  function getAllowListId() external view override returns (bytes32) {
    return s_allowListId;
  }

  /// @inheritdoc IFunctionsRouter
  function setAllowListId(bytes32 allowListId) external override onlyOwner {
    s_allowListId = allowListId;
  }

  /// @dev Used within FunctionsSubscriptions.sol
  function _getMaxConsumers() internal view override returns (uint16) {
    return s_config.maxConsumersPerSubscription;
  }

  /// @dev Used within FunctionsSubscriptions.sol
  function _getSubscriptionDepositDetails() internal view override returns (uint16, uint72) {
    return (s_config.subscriptionDepositMinimumRequests, s_config.subscriptionDepositJuels);
  }

  // ================================================================
  // |                           Requests                           |
  // ================================================================

  /// @inheritdoc IFunctionsRouter
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

  /// @inheritdoc IFunctionsRouter
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
    _isExistingSubscription(subscriptionId);
    _isAllowedConsumer(msg.sender, subscriptionId);
    isValidCallbackGasLimit(subscriptionId, callbackGasLimit);

    if (data.length == 0) {
      revert EmptyRequestData();
    }

    Subscription memory subscription = getSubscription(subscriptionId);
    Consumer memory consumer = getConsumer(msg.sender, subscriptionId);
    uint72 adminFee = s_config.adminFee;

    // Forward request to DON
    FunctionsResponse.Commitment memory commitment = coordinator.startRequest(
      FunctionsResponse.RequestMeta({
        requestingContract: msg.sender,
        data: data,
        subscriptionId: subscriptionId,
        dataVersion: dataVersion,
        flags: getFlags(subscriptionId),
        callbackGasLimit: callbackGasLimit,
        adminFee: adminFee,
        initiatedRequests: consumer.initiatedRequests,
        completedRequests: consumer.completedRequests,
        availableBalance: subscription.balance - subscription.blockedBalance,
        subscriptionOwner: subscription.owner
      })
    );

    // Do not allow setting a comittment for a requestId that already exists
    if (s_requestCommitments[commitment.requestId] != bytes32(0)) {
      revert DuplicateRequestId(commitment.requestId);
    }

    // Store a commitment about the request
    s_requestCommitments[commitment.requestId] = keccak256(
      abi.encode(
        FunctionsResponse.Commitment({
          adminFee: adminFee,
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
      subscriptionOwner: subscription.owner,
      requestingContract: msg.sender,
      // solhint-disable-next-line avoid-tx-origin
      requestInitiator: tx.origin,
      data: data,
      dataVersion: dataVersion,
      callbackGasLimit: callbackGasLimit,
      estimatedTotalCostJuels: commitment.estimatedTotalCostJuels
    });

    return commitment.requestId;
  }

  // ================================================================
  // |                           Responses                          |
  // ================================================================

  /// @inheritdoc IFunctionsRouter
  function fulfill(
    bytes memory response,
    bytes memory err,
    uint96 juelsPerGas,
    uint96 costWithoutFulfillment,
    address transmitter,
    FunctionsResponse.Commitment memory commitment
  ) external override returns (FunctionsResponse.FulfillResult resultCode, uint96) {
    _whenNotPaused();

    if (msg.sender != commitment.coordinator) {
      revert OnlyCallableFromCoordinator();
    }

    {
      bytes32 commitmentHash = s_requestCommitments[commitment.requestId];

      if (commitmentHash == bytes32(0)) {
        resultCode = FunctionsResponse.FulfillResult.INVALID_REQUEST_ID;
        emit RequestNotProcessed(commitment.requestId, commitment.coordinator, transmitter, resultCode);
        return (resultCode, 0);
      }

      if (keccak256(abi.encode(commitment)) != commitmentHash) {
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
    }

    {
      uint96 callbackCost = juelsPerGas * SafeCast.toUint96(commitment.callbackGasLimit);
      uint96 totalCostJuels = commitment.adminFee + costWithoutFulfillment + callbackCost;

      // Check that the subscription can still afford to fulfill the request
      if (totalCostJuels > getSubscription(commitment.subscriptionId).balance) {
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
      ? FunctionsResponse.FulfillResult.FULFILLED
      : FunctionsResponse.FulfillResult.USER_CALLBACK_ERROR;

    Receipt memory receipt = _pay(
      commitment.subscriptionId,
      commitment.estimatedTotalCostJuels,
      commitment.client,
      commitment.adminFee,
      juelsPerGas,
      SafeCast.toUint96(result.gasUsed),
      costWithoutFulfillment
    );

    emit RequestProcessed({
      requestId: commitment.requestId,
      subscriptionId: commitment.subscriptionId,
      totalCostJuels: receipt.totalCostJuels,
      transmitter: transmitter,
      resultCode: resultCode,
      response: response,
      err: err,
      callbackReturnData: result.returnData
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
    bool destinationNoLongerExists;
    assembly {
      // solidity calls check that a contract actually exists at the destination, so we do the same
      destinationNoLongerExists := iszero(extcodesize(client))
    }
    if (destinationNoLongerExists) {
      // Return without attempting callback
      // The subscription will still be charged to reimburse transmitter's gas overhead
      return CallbackResult({success: false, gasUsed: 0, returnData: new bytes(0)});
    }

    bytes memory encodedCallback = abi.encodeWithSelector(
      s_config.handleOracleFulfillmentSelector,
      requestId,
      response,
      err
    );

    uint16 gasForCallExactCheck = s_config.gasForCallExactCheck;

    // Call with explicitly the amount of callback gas requested
    // Important to not let them exhaust the gas budget and avoid payment.
    // NOTE: that callWithExactGas will revert if we do not have sufficient gas
    // to give the callee their requested amount.

    bool success;
    uint256 gasUsed;
    // allocate return data memory ahead of time
    bytes memory returnData = new bytes(MAX_CALLBACK_RETURN_BYTES);

    assembly {
      let g := gas()
      // Compute g -= gasForCallExactCheck and check for underflow
      // The gas actually passed to the callee is _min(gasAmount, 63//64*gas available).
      // We want to ensure that we revert if gasAmount >  63//64*gas available
      // as we do not want to provide them with less, however that check itself costs
      // gas. gasForCallExactCheck ensures we have at least enough gas to be able
      // to revert if gasAmount >  63//64*gas available.
      if lt(g, gasForCallExactCheck) {
        revert(0, 0)
      }
      g := sub(g, gasForCallExactCheck)
      // if g - g//64 <= gasAmount, revert
      // (we subtract g//64 because of EIP-150)
      if iszero(gt(sub(g, div(g, 64)), callbackGasLimit)) {
        revert(0, 0)
      }
      // call and report whether we succeeded
      // call(gas,addr,value,argsOffset,argsLength,retOffset,retLength)
      let gasBeforeCall := gas()
      success := call(callbackGasLimit, client, 0, add(encodedCallback, 0x20), mload(encodedCallback), 0, 0)
      gasUsed := sub(gasBeforeCall, gas())

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

  // ================================================================
  // |                        Route methods                         |
  // ================================================================

  /// @inheritdoc IFunctionsRouter
  function getContractById(bytes32 id) public view override returns (address) {
    address currentImplementation = s_route[id];
    if (currentImplementation == address(0)) {
      revert RouteNotFound(id);
    }
    return currentImplementation;
  }

  /// @inheritdoc IFunctionsRouter
  function getProposedContractById(bytes32 id) public view override returns (address) {
    // Iterations will not exceed MAX_PROPOSAL_SET_LENGTH
    for (uint8 i = 0; i < s_proposedContractSet.ids.length; ++i) {
      if (id == s_proposedContractSet.ids[i]) {
        return s_proposedContractSet.to[i];
      }
    }
    revert RouteNotFound(id);
  }

  // ================================================================
  // |                 Contract Proposal methods                    |
  // ================================================================

  /// @inheritdoc IFunctionsRouter
  function getProposedContractSet() external view override returns (bytes32[] memory, address[] memory) {
    return (s_proposedContractSet.ids, s_proposedContractSet.to);
  }

  /// @inheritdoc IFunctionsRouter
  function proposeContractsUpdate(
    bytes32[] memory proposedContractSetIds,
    address[] memory proposedContractSetAddresses
  ) external override onlyOwner {
    // IDs and addresses arrays must be of equal length and must not exceed the max proposal length
    uint256 idsArrayLength = proposedContractSetIds.length;
    if (idsArrayLength != proposedContractSetAddresses.length || idsArrayLength > MAX_PROPOSAL_SET_LENGTH) {
      revert InvalidProposal();
    }

    // NOTE: iterations of this loop will not exceed MAX_PROPOSAL_SET_LENGTH
    for (uint256 i = 0; i < idsArrayLength; ++i) {
      bytes32 id = proposedContractSetIds[i];
      address proposedContract = proposedContractSetAddresses[i];
      if (
        proposedContract == address(0) || // The Proposed address must be a valid address
        s_route[id] == proposedContract // The Proposed address must point to a different address than what is currently set
      ) {
        revert InvalidProposal();
      }

      emit ContractProposed({
        proposedContractSetId: id,
        proposedContractSetFromAddress: s_route[id],
        proposedContractSetToAddress: proposedContract
      });
    }

    s_proposedContractSet = ContractProposalSet({ids: proposedContractSetIds, to: proposedContractSetAddresses});
  }

  /// @inheritdoc IFunctionsRouter
  function updateContracts() external override onlyOwner {
    // Iterations will not exceed MAX_PROPOSAL_SET_LENGTH
    for (uint256 i = 0; i < s_proposedContractSet.ids.length; ++i) {
      bytes32 id = s_proposedContractSet.ids[i];
      address to = s_proposedContractSet.to[i];
      emit ContractUpdated({id: id, from: s_route[id], to: to});
      s_route[id] = to;
    }

    delete s_proposedContractSet;
  }

  // ================================================================
  // |                           Modifiers                          |
  // ================================================================
  // Favoring internal functions over actual modifiers to reduce contract size

  /// @dev Used within FunctionsSubscriptions.sol
  function _whenNotPaused() internal view override {
    _requireNotPaused();
  }

  /// @dev Used within FunctionsSubscriptions.sol
  function _onlyRouterOwner() internal view override {
    _validateOwnership();
  }

  /// @dev Used within FunctionsSubscriptions.sol
  function _onlySenderThatAcceptedToS() internal view override {
    address currentImplementation = s_route[s_allowListId];
    if (currentImplementation == address(0)) {
      // If not set, ignore this check, allow all access
      return;
    }
    if (!IAccessController(currentImplementation).hasAccess(msg.sender, new bytes(0))) {
      revert SenderMustAcceptTermsOfService(msg.sender);
    }
  }

  /// @inheritdoc IFunctionsRouter
  function pause() external override onlyOwner {
    _pause();
  }

  /// @inheritdoc IFunctionsRouter
  function unpause() external override onlyOwner {
    _unpause();
  }
}
