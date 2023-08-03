// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";
import {IFunctionsCoordinator} from "./interfaces/IFunctionsCoordinator.sol";
import {IAccessController} from "../../../shared/interfaces/IAccessController.sol";
import {IConfigurable} from "./interfaces/IConfigurable.sol";

import {FunctionsSubscriptions} from "./FunctionsSubscriptions.sol";
import {FunctionsResponse} from "./libraries/FunctionsResponse.sol";
import {ConfirmedOwner} from "../../../ConfirmedOwner.sol";

import {SafeCast} from "../../../vendor/openzeppelin-solidity/v4.8.0/contracts/utils/SafeCast.sol";
import {Pausable} from "../../../vendor/openzeppelin-solidity/v4.8.0/contracts/security/Pausable.sol";

contract FunctionsRouter is IFunctionsRouter, FunctionsSubscriptions, Pausable, ITypeAndVersion, ConfirmedOwner {
  using FunctionsResponse for FunctionsResponse.Commitment;
  using FunctionsResponse for FunctionsResponse.FulfillResult;

  string public constant override typeAndVersion = "Functions Router v1.0.0";

  // We limit return data to a selector plus 4 words. This is to avoid
  // malicious contracts from returning large amounts of data and causing
  // repeated out-of-gas scenarios.
  uint16 public constant MAX_CALLBACK_RETURN_BYTES = 4 + 4 * 32;

  mapping(bytes32 id => address routableContract) private s_route;

  error RouteNotFound(bytes32 id);

  // Use empty bytes to self-identify, since it does not have an id
  bytes32 private constant ROUTER_ID = bytes32(0);

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

  event ConfigChanged(Config);

  error OnlyCallableByRoute();

  // ================================================================
  // |                          Timelock state                      |
  // ================================================================
  uint16 private immutable s_maximumTimelockBlocks;
  uint16 private s_timelockBlocks;

  struct TimeLockProposal {
    uint16 from;
    uint16 to;
    uint64 timelockEndBlock;
  }

  TimeLockProposal private s_timelockProposal;

  event TimeLockProposed(uint16 from, uint16 to);
  event TimeLockUpdated(uint16 from, uint16 to);
  error ProposedTimelockAboveMaximum();
  error TimelockInEffect();

  // ================================================================
  // |                         Proposal state                       |
  // ================================================================

  uint8 private constant MAX_PROPOSAL_SET_LENGTH = 8;

  struct ContractProposalSet {
    bytes32[] ids;
    address[] to;
    uint64 timelockEndBlock;
  }
  ContractProposalSet private s_proposedContractSet;

  event ContractProposed(
    bytes32 proposedContractSetId,
    address proposedContractSetFromAddress,
    address proposedContractSetToAddress,
    uint64 timelockEndBlock
  );

  event ContractUpdated(bytes32 id, address from, address to);

  struct ConfigProposal {
    bytes to;
    uint64 timelockEndBlock;
  }
  mapping(bytes32 id => ConfigProposal) private s_proposedConfig;
  event ConfigProposed(bytes32 id, bytes toBytes);
  event ConfigUpdated(bytes32 id, bytes toBytes);
  error InvalidProposal();
  error IdentifierIsReserved(bytes32 id);

  constructor(
    uint16 timelockBlocks,
    uint16 maximumTimelockBlocks,
    address linkToken,
    Config memory config
  ) FunctionsSubscriptions(linkToken) ConfirmedOwner(msg.sender) Pausable() {
    // Set initial value for the number of blocks of the timelock
    s_timelockBlocks = timelockBlocks;
    // Set maximum number of blocks that the timelock can be
    s_maximumTimelockBlocks = maximumTimelockBlocks;
    // Set the initial configuration for the Router
    s_route[ROUTER_ID] = address(this);
    _updateConfig(config);
  }

  // @inheritdoc IFunctionsRouter
  function getAllowListId() external pure override returns (bytes32) {
    return ALLOW_LIST_ID;
  }

  // ================================================================
  // |                        Configuration                         |
  // ================================================================

  // @inheritdoc IFunctionsRouter
  function getConfig() external view override returns (Config memory) {
    return s_config;
  }

  // @inheritdoc IRouterBase
  function proposeConfigUpdateSelf(bytes calldata config) external override onlyOwner {
    s_proposedConfig[ROUTER_ID] = ConfigProposal({
      to: config,
      timelockEndBlock: uint64(block.number + s_timelockBlocks)
    });
    emit ConfigProposed({id: ROUTER_ID, toBytes: config});
  }

  // @inheritdoc IRouterBase
  function updateConfigSelf() external override onlyOwner {
    ConfigProposal memory proposal = s_proposedConfig[ROUTER_ID];
    if (block.number < proposal.timelockEndBlock) {
      revert TimelockInEffect();
    }
    _updateConfig(abi.decode(proposal.to, (Config)));
    emit ConfigUpdated({id: ROUTER_ID, toBytes: proposal.to});
  }

  // @inheritdoc IRouterBase
  function proposeConfigUpdate(bytes32 id, bytes calldata config) external override onlyOwner {
    s_proposedConfig[id] = ConfigProposal({to: config, timelockEndBlock: uint64(block.number + s_timelockBlocks)});
    emit ConfigProposed({id: id, toBytes: config});
  }

  // @inheritdoc IRouterBase
  function updateConfig(bytes32 id) external override onlyOwner {
    ConfigProposal memory proposal = s_proposedConfig[id];

    if (block.number < proposal.timelockEndBlock) {
      revert TimelockInEffect();
    }

    IConfigurable(getContractById(id)).updateConfig(proposal.to);

    emit ConfigUpdated({id: id, toBytes: proposal.to});
  }

  // @notice Sets the configuration for FunctionsRouter specific state
  // @param config bytes of config data to set the following:
  // - adminFee: fee that will be paid to the Router owner for operating the network
  function _updateConfig(Config memory config) internal {
    s_config = config;
    emit ConfigChanged(config);
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

  // Used within FunctionsSubscriptions.sol
  function _getMaxConsumers() internal view override returns (uint16) {
    return s_config.maxConsumersPerSubscription;
  }

  // ================================================================
  // |                           Requests                           |
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

  // ================================================================
  // |                           Responses                          |
  // ================================================================

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

  // ================================================================
  // |                        Route methods                         |
  // ================================================================

  // @inheritdoc IRouterBase
  function getContractById(bytes32 id) public view override returns (address) {
    address currentImplementation = s_route[id];
    if (currentImplementation == address(0)) {
      revert RouteNotFound(id);
    }
    return currentImplementation;
  }

  // @inheritdoc IRouterBase
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
  // |                            Timelock                          |
  // ================================================================

  // @inheritdoc IRouterBase
  function proposeTimelockBlocks(uint16 blocks) external override onlyOwner {
    if (s_timelockBlocks == blocks) {
      revert InvalidProposal();
    }
    if (blocks > s_maximumTimelockBlocks) {
      revert ProposedTimelockAboveMaximum();
    }
    s_timelockProposal = TimeLockProposal({
      from: s_timelockBlocks,
      to: blocks,
      timelockEndBlock: uint64(block.number + s_timelockBlocks)
    });
  }

  // @inheritdoc IRouterBase
  function updateTimelockBlocks() external override onlyOwner {
    if (block.number < s_timelockProposal.timelockEndBlock) {
      revert TimelockInEffect();
    }
    s_timelockBlocks = s_timelockProposal.to;
  }

  // ================================================================
  // |                 Contract Proposal methods                    |
  // ================================================================

  // @inheritdoc IRouterBase
  function getProposedContractSet()
    external
    view
    override
    returns (uint256 timelockEndBlock, bytes32[] memory ids, address[] memory to)
  {
    timelockEndBlock = s_proposedContractSet.timelockEndBlock;
    ids = s_proposedContractSet.ids;
    to = s_proposedContractSet.to;
    return (timelockEndBlock, ids, to);
  }

  // @inheritdoc IRouterBase
  function proposeContractsUpdate(
    bytes32[] memory proposedContractSetIds,
    address[] memory proposedContractSetAddresses
  ) external override onlyOwner {
    // All arrays must be of equal length and not must not exceed the max length
    uint256 idsArrayLength = proposedContractSetIds.length;
    if (idsArrayLength != proposedContractSetAddresses.length || idsArrayLength > MAX_PROPOSAL_SET_LENGTH) {
      revert InvalidProposal();
    }
    // Iterations will not exceed MAX_PROPOSAL_SET_LENGTH
    for (uint256 i = 0; i < idsArrayLength; ++i) {
      bytes32 id = proposedContractSetIds[i];
      address proposedContract = proposedContractSetAddresses[i];
      if (
        proposedContract == address(0) || // The Proposed address must be a valid address
        s_route[id] == proposedContract // The Proposed address must point to a different address than what is currently set
      ) {
        revert InvalidProposal();
      }
      // Reserved ids cannot be set
      if (id == ROUTER_ID) {
        revert IdentifierIsReserved(id);
      }
    }

    uint64 timelockEndBlock = uint64(block.number + s_timelockBlocks);

    s_proposedContractSet = ContractProposalSet({
      ids: proposedContractSetIds,
      to: proposedContractSetAddresses,
      timelockEndBlock: timelockEndBlock
    });

    // Iterations will not exceed MAX_PROPOSAL_SET_LENGTH
    for (uint256 i = 0; i < proposedContractSetIds.length; ++i) {
      emit ContractProposed({
        proposedContractSetId: proposedContractSetIds[i],
        proposedContractSetFromAddress: s_route[proposedContractSetIds[i]],
        proposedContractSetToAddress: proposedContractSetAddresses[i],
        timelockEndBlock: timelockEndBlock
      });
    }
  }

  // @inheritdoc IRouterBase
  function updateContracts() external override onlyOwner {
    if (block.number < s_proposedContractSet.timelockEndBlock) {
      revert TimelockInEffect();
    }
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

  // @inheritdoc IRouterBase
  function pause() external override onlyOwner {
    _pause();
  }

  // @inheritdoc IRouterBase
  function unpause() external override onlyOwner {
    _unpause();
  }
}
